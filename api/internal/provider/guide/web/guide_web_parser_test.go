package guide_web_provider

import (
	"fmt"
	"os"
	"testing"
	"via/internal/log"
	"via/internal/model"

	app_log "via/internal/log/app"

	"github.com/stretchr/testify/require"
)

type brokenReader struct {
}

func (b brokenReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("boom")
}

func TestParseHistoricalQueryResponse(t *testing.T) {

	log.Set(app_log.New(app_log.LogCfg{Level: "debug", ConsoleWriter: app_log.ConsoleWriterCfg{Enabled: true}}))
	t.Run("success", func(t *testing.T) {
		f, err := os.Open("testdata/consulta_historico_resultado_success.html")
		require.NoError(t, err)
		defer f.Close()

		guide, err := ParseHistoricalQueryResponse(f)
		require.NoError(t, err)

		expected := model.Guide{
			ID:          "999025862539",
			Reference:   "999025862539",
			Status:      "ENT",
			Packages:    1,
			Weight:      1.714,
			Shipping:    "D",
			Route:       "TUC-MDQ",
			Date:        "13/03/25",
			Sender:      "ROLDAN PAULA",
			Recipient:   "DANIELA AGUIRRE @",
			Destination: model.Destination{ID: "7600", Description: "MAR DEL PLATA"},
		}
		require.Equal(t, expected, guide)
	})

	t.Run("error invalid input", func(t *testing.T) {
		r := brokenReader{}
		// Simulate a broken reader that returns an error
		_, err := ParseHistoricalQueryResponse(r)
		//require.Error(t, err)
		require.ErrorIs(t, err, ErrBadResponse)
	})

	t.Run("error no rows", func(t *testing.T) {
		f, err := os.Open("testdata/consulta_historico_resultado_empty.html")
		require.NoError(t, err)
		defer f.Close()

		_, err = ParseHistoricalQueryResponse(f)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no result row found")
	})

	t.Run("error insufficient columns", func(t *testing.T) {
		f, err := os.Open("testdata/consulta_historico_resultado_missing_cols.html")
		require.NoError(t, err)
		defer f.Close()

		_, err = ParseHistoricalQueryResponse(f)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrMissingColumn)
	})
}
