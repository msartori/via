package via_guide_web_provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	http_client "via/internal/client/http"
	"via/internal/log"
	"via/internal/model"
	"via/internal/secret"
)

type ViaGuideWebProvider struct {
	client      *http_client.HttpClient
	guideParser ViaResponseParser
}

func New(cfg http_client.HttpClientCfg, guideParser ViaResponseParser) *ViaGuideWebProvider {
	cfg.AuthorizationHeaderSecret = secret.ReadSecret(cfg.AuthorizationHeaderSecret)
	return (&ViaGuideWebProvider{
		client:      http_client.New(cfg),
		guideParser: guideParser,
	})
}

func (p *ViaGuideWebProvider) GetGuide(ctx context.Context, id string) (model.ViaGuide, error) {
	logger := log.Get()
	params := url.Values{"nenvio": {id}, "pagina": {"1"}}
	reqRes := fmt.Sprintf("%s/atencion_cliente/historico/consulta_historico_resultado.do", p.client.BaseURL)
	logger.Info(ctx, "msg", "http request", "resource", reqRes, "params", params)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		reqRes,
		strings.NewReader(params.Encode()))

	if err != nil {
		logger.Error(ctx, err, "msg", "error creatig request")
		return model.ViaGuide{}, fmt.Errorf("creating request: %w", err)
	}

	// Set content-type for form encoding
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if p.client.AuthorizationHeader != "" {
		req.Header.Set("Authorization", p.client.AuthorizationHeader)
	}

	resp, err := p.client.Requester.Do(req)
	if err != nil {
		logger.Error(ctx, err, "msg", "error making HTTP request")
		return model.ViaGuide{}, fmt.Errorf("making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(ctx, err, "msg", "unexpected status code")
		return model.ViaGuide{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err, "msg", "error reading response")
		return model.ViaGuide{}, fmt.Errorf("error reading response: %w", err)
	}
	var guide model.ViaGuide
	err = p.guideParser.Parse(bodyBytes, &guide)

	if err != nil {
		if !errors.Is(err, ErrNoResultRow) {
			return guide, err
		}
	}
	return guide, nil
}
