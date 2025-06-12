package guide_web_provider

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"via/internal/model"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNoResultRow   = errors.New("no result row found")
	ErrMissingColumn = errors.New("missing expected column")
	ErrMalformedHTML = errors.New("malformed HTML")
	ErrBadResponse   = errors.New("malformed response")
)

// ViaCargo/atencion_cliente/historico/consulta_historico_resultado.do
func ParseHistoricalQueryResponse(r io.Reader) (model.Guide, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return model.Guide{}, fmt.Errorf("%w: %v", ErrBadResponse, err)
	}
	htmlStr := string(raw)

	// remove comments from <td> elements
	htmlStr = uncommentTDs(htmlStr)
	// new reader with cleaned HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return model.Guide{}, fmt.Errorf("%w: %v", ErrMalformedHTML, err)
	}

	// extract header row
	headerRow := doc.Find("#listado tr.cabecera").First()
	headerMap := make(map[string]int)

	headerRow.Find("td").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Find("span.cabecera").Text())
		if text != "" {
			headerMap[text] = i
		}
	})

	if len(headerMap) == 0 {
		return model.Guide{}, ErrNoResultRow
	}

	requiredHeaders := []string{
		"Envío", "Estado", "Acciones",
	}
	for _, header := range requiredHeaders {
		if _, ok := headerMap[header]; !ok {
			return model.Guide{}, fmt.Errorf("%w: %s", ErrMissingColumn, header)
		}
	}

	// Get data
	row := doc.Find("#listado tbody.cuerpo_tabla tr").First()

	cells := row.Find("td")
	getCell := func(header string) string {
		idx := headerMap[header]
		return strings.TrimSpace(cells.Eq(idx).Text())
	}

	parseInt := func(header string) int {
		v, _ := strconv.Atoi(getCell(header))
		return v
	}

	parseFloat := func(header string) float64 {
		v, _ := strconv.ParseFloat(strings.ReplaceAll(getCell(header), ",", "."), 64)
		return v
	}

	parseDestination := func(header string) model.Destination {
		parts := strings.Fields(getCell(header))
		if len(parts) < 2 {
			return model.Destination{}
		}

		description := strings.Join(parts[1:], " ")
		return model.Destination{
			ID:          parts[0],
			Description: description,
		}

	}

	guide := model.Guide{
		ID:          getCell("Envío"),
		Reference:   getCell("Referencia"),
		Status:      getCell("Estado"),
		Packages:    parseInt("BUL"),
		Weight:      parseFloat("PR Kg"),
		Shipping:    getCell("Portes"),
		Route:       getCell("Ruta"),
		Date:        getCell("Fecha"),
		Sender:      getCell("Remitente"),
		Recipient:   getCell("Destinatario"),
		Destination: parseDestination("Acciones"),
	}

	return guide, nil
}

func uncommentTDs(html string) string {
	re := regexp.MustCompile(`<!--\s*(<td[^>]*>.*?</td>)\s*-->`)
	return re.ReplaceAllString(html, "$1")
}
