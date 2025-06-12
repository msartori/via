package guide_web_provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	http_client "via/internal/client/http"
	custom_error "via/internal/error"
	"via/internal/model"
	"via/internal/secret"
)

type GuideWebProvider struct {
	client *http_client.HttpClient
}

func New(cfg http_client.HttpClientCfg) *GuideWebProvider {
	cfg.AuthorizationHeaderSecret = secret.ReadSecret(cfg.AuthorizationHeaderSecret)
	return (&GuideWebProvider{
		client: http_client.New(cfg),
	})
}

func (p *GuideWebProvider) GetGuide(ctx context.Context, id string) (model.Guide, error) {
	params := url.Values{"nenvio": {id}, "pagina": {"1"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/atencion_cliente/historico/consulta_historico_resultado.do", p.client.BaseURL),
		strings.NewReader(params.Encode()))

	if err != nil {
		return model.Guide{}, fmt.Errorf("creating request: %w", err)
	}

	// Set content-type for form encoding
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if p.client.AuthorizationHeader != "" {
		req.Header.Set("Authorization", p.client.AuthorizationHeader)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return model.Guide{}, fmt.Errorf("making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Guide{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	guide, err := ParseHistoricalQueryResponse(resp.Body)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadResponse):
		case errors.Is(err, ErrMalformedHTML):
		case errors.Is(err, ErrMissingColumn):
			return guide, custom_error.NewHttpError(http.StatusInternalServerError, "Error interno de servidor")
		case errors.Is(err, ErrNoResultRow):
			return guide, custom_error.NewHttpError(http.StatusNotFound, "Guía no econtrada")
		default:
			return guide, err
		}
	}
	return guide, nil

}
