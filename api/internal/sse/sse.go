package sse

import (
	"errors"
	"net/http"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/pubsub"
	"via/internal/response"
)

type Loader[T any] func(r *http.Request) response.Response[T]

func HandleSSE(w http.ResponseWriter, r *http.Request, loader Loader[any], eventSource ...string) {
	r.Header.Set("Accept-Language", r.URL.Query().Get("lang"))
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Get().Error(r.Context(), errors.New("unable to get flush from writer"),
			"msg", "unable to create SSE")
		res := response.Response[any]{Message: i18n.Get(r, i18n.MsgInternalServerError)}
		response.WriteJSONEvent(w, r, res)
		return
	}

	sub, err := pubsub.Get().Subscribe(r.Context(), eventSource...)
	if err != nil {
		log.Get().Error(r.Context(), err, "msg", "unable to subscribe", "channel", eventSource)
		res := response.Response[any]{Message: i18n.Get(r, i18n.MsgInternalServerError)}
		response.WriteJSONEvent(w, r, res)
		return
	}
	defer sub.Close()

	//initial connection -> respond original data
	log.Get().Info(r.Context(), "event", "initial connection")
	response.WriteJSONEvent(w, r, loader(r))
	flusher.Flush()

	for {
		select {
		case event, ok := <-sub.Channel():
			if !ok {
				log.Get().Warn(r.Context(), "msg", "Subscription channel closed")
				res := response.Response[any]{Message: i18n.Get(r, i18n.MsgInternalServerError)}
				response.WriteJSONEvent(w, r, res)
				return
			}
			log.Get().Info(r.Context(), "event", event)
			response.WriteJSONEvent(w, r, loader(r))
			flusher.Flush()
		case <-r.Context().Done():
			log.Get().Info(r.Context(), "msg", "Client disconnected")
			return
		}
	}
}
