package via_guide_web_provider

import (
	"os"
	"testing"
	"via/internal/log"
	"via/internal/model"

	app_log "via/internal/log/app"

	"github.com/stretchr/testify/require"
)

func TestParseHistoricalQueryResponse(t *testing.T) {

	log.Set(app_log.New(app_log.LogCfg{Level: "debug", DefaultWriter: app_log.DefaultWriterCfg{Enabled: true}}))
	t.Run("success", func(t *testing.T) {
		f, err := os.ReadFile("testdata/consulta_historico_resultado_success.html")
		require.NoError(t, err)
		var guide model.ViaGuide
		err = HistoricalQueryResponseParser{}.Parse(f, &guide)
		require.NoError(t, err)

		expected := model.ViaGuide{
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
			Destination: model.ViaDestination{ID: "7600", Description: "MAR DEL PLATA"},
		}
		require.Equal(t, expected, guide)
	})

	t.Run("error no rows", func(t *testing.T) {
		f, err := os.ReadFile("testdata/consulta_historico_resultado_empty.html")
		require.NoError(t, err)
		var guide model.ViaGuide
		err = HistoricalQueryResponseParser{}.Parse(f, &guide)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no result row found")
	})

	t.Run("error insufficient columns", func(t *testing.T) {
		f, err := os.ReadFile("testdata/consulta_historico_resultado_missing_cols.html")
		require.NoError(t, err)
		var guide model.ViaGuide

		err = HistoricalQueryResponseParser{}.Parse(f, &guide)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrMissingColumn)
	})
}
