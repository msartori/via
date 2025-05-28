package provider_guide

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TableRow struct {
	Envio        string `json:"envio"`
	Referencia   string `json:"referencia"`
	Estado       string `json:"estado"`
	Bul          string `json:"bul"`
	PRKG         string `json:"prgk"`
	PORTES       string `json:"portes"`
	RUTA         string `json:"ruta"`
	FECHA        string `json:"fecha"`
	REMITENTE    string `json:"remitente"`
	DESTINATARIO string `json:"destinatario"`
}

func extractText(selection *goquery.Selection) string {
	return strings.TrimSpace(selection.Text())
}

func main() {
	/*
			html := `<html> <body><table>
			 <tr class="cabecera">
		         <td nowrap><span class="cabecera">Envío</span></td>
		         <td nowrap><span class="cabecera">Referencia</span></td>
		         <td nowrap><span class="cabecera">Estado</span></td>
		         <td nowrap align="right"><span class="cabecera">BUL</span></td>
		         <td nowrap align="right"><span class="cabecera">PR Kg</span></td>
		         <td nowrap><span class="cabecera">Portes</span></td>
		         <td nowrap><span class="cabecera">Ruta</span></td>
		         <!--<td nowrap><span class="cabecera">Destino</span></td>-->
		         <td nowrap><span class="cabecera">Fecha </span></td>
		         <td><span class="cabecera">Remitente</span></td>
		         <td><span class="cabecera">Destinatario</span></td>
		         <!--<td nowrap><span class="cabecera">Reem</span></td>-->
				 <td nowrap><span class="cabecera">Acciones</span></td>
				 <!--<td nowrap><span class="cabecera">Población destino</span></td>-->
		        </tr>

			<tbody class="cuerpo_tabla">
			<tr class="color_impar">
			 <td><a href="javascript:verExpedicion('27033228');">999025836102</a></td>
			 <td>999025836102</td>
			 <td>CRRG</td>
			 <td align="right">
			 	<a href="javascript:imprimirRangoEtiquetas('27033228')" title='Imprimir rótulos'>
			 		1
			 	</a>
			 </td>
			 <td align="right">5,900</td>
			 <td align="center">P</td>
			 <td nowrap>RTO-SF006</td>
			 <td>13/03/25</td>
			 <td>TELLO MONICA</td>
			 <td><span title="TOMADIN ROBERTO CARLOS">TOMADIN ROBERTO CA</span></td>
			 <td>
				<table class="acciones">
					<tr>
						<td>
							<a class="btn2" href="javascript:detalleParcial(49,'27033228');" title="Detalle">DE</a>
						</td>
					</tr>
				</table>
			 </td>
			</tr>

			<tr class="color_par">
			 <td><a href="javascript:verExpedicion('27033227');">999025836101</a></td>
			 <td><span title="HORMIGONERA 1HP">HORMIGONERA </span></td>
			 <td>CRRG</td>
			 <td align="right">
			 	<a href="javascript:imprimirRangoEtiquetas('27033227')" title='Imprimir rótulos'>
			 		1
			 	</a>
			 </td>
			 <td align="right">25,000</td>
			 <td align="center">P</td>
			 <td nowrap>SF015-EN015</td>
			 <td>13/03/25</td>
			 <td><span title="METALURGICACDG S.A">METALURGICACDG </span></td>
			 <td>MARCELO LUGRIN</td>
			 <td>
				<table class="acciones">
					<tr>
						<td>
							<a class="btn2" href="javascript:detalleParcial(50,'27033227');" title="Detalle">DE</a>
						</td>
					</tr>
				</table>
			 </td>
			</tr>
		</tbody> </table></body></html>`
			/*
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
				if err != nil {
					log.Fatal(err)
				}*/

	doc := queryAPI()

	var jsonRows []TableRow
	doc.Find("tbody.cuerpo_tabla").Each(func(i int, s *goquery.Selection) {
		rows := s.Find("tr")
		rows.Each(func(j int, s *goquery.Selection) {
			cols := s.Find("td")
			if cols.Length() < 2 {
				return
			}
			row := TableRow{}
			cols.Each(func(k int, s *goquery.Selection) {

				if k > 9 {
					return
				}
				switch k {
				case 0:
					row.Envio, _ = s.Html()
					row.Envio = getExpedicionNumber(row.Envio)
					if row.Envio == "expedicion no encontrada" {
						return
					}
				case 1:
					row.Referencia = s.Text()
				case 2:
					row.Estado = s.Text()
				case 3:
					row.Bul = regexp.MustCompile(`\s+`).ReplaceAllString(s.Text(), "")
				case 4:
					row.PRKG = s.Text()
				case 5:
					row.PORTES = s.Text()
				case 6:
					row.RUTA = s.Text()
				case 7:
					row.FECHA = s.Text()
				case 8:
					row.REMITENTE = s.Text()
				case 9:
					row.DESTINATARIO = s.Text()
				}
			})
			if row.Envio != "expedicion no encontrada" {
				jsonRows = append(jsonRows, row)
			}
		})
	})

	jsonData, err := json.MarshalIndent(jsonRows, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}

func getExpedicionNumber(str string) string {
	decodedStr := html.UnescapeString(str)
	re := regexp.MustCompile(`verExpedicion\('(\d+)'\)`)
	match := re.FindStringSubmatch(decodedStr)
	if len(match) > 1 {
		return match[1] // Output: 27033227
	} else {
		return "expedicion no encontrada"
	}
}

func queryAPI() *goquery.Document {
	url := "https://viacargo.alertran.net/ViaCargo/atencion_cliente/historico/consulta_historico_resultado.do?pagina=1&devolucion_conforme=false&sabados=false&retorno=false&formato=&tipo_sel=&c_refcli=&c_nomcons=&c_dircons=&c_pobcons=&c_cpcons=&c_nombreRemitente=&c_direccionRemitente=&c_poblacionRemitente=&c_cpRemitente=&c_bultos=&c_peso=&c_vol=&c_portes=&c_observ1=&c_observ2=&c_fgrab=&c_fsal=&c_flle=&c_fent=&c_estexp=&c_nomrecep=&c_rutrecep=&c_otrosrecep=&c_incexp=&c_importeReembolso=&c_localizador=&c_tipoServicio=&c_totalPagados=&c_totalDebidos=&c_fecha=&c_expedicion=&c_delegacion_bulto=&c_ruta=&delegacion=000&delegacioncentral=000&soloRef=&etiquetador_codigo=&nenvio=&referencia=&referenciaBulto=&codigoViaje=&trafico=S&producto_codigo=&producto_descripcion=&fechaDesde=10%2F03%2F2025&fechaHasta=13%2F03%2F2025&tipo_remitente=CLIR&remitente_codigo=&remitente_descripcion=&remitenteAgrupador_codigo=&remitenteAgrupador_descripcion=&tipo_consignatario=CLIC&consignatario_codigo=&consignatario_descripcion=&consignatarioAgrupador_codigo=&consignatarioAgrupador_descripcion=&tipoClienteRem_codigo=&tipoClienteRem_descripcion=&tipoClienteCons_codigo=&tipoClienteCons_descripcion=&origen_codigo=&origen_descripcion=&destino_codigo=&destino_descripcion=&delegacionRemitente_codigo=&delegacionRemitente_descripcion=&codigoSac=&situacion=&codigoCun=&codigoPostal_codigo=&codigoPostal_descripcion=&codigoPostalRemi_codigo=&codigoPostalRemi_descripcion=&recogidaCodigo=&tipo_evento_codigo=&tipo_evento_descripcion=&ficheroNombre=&identificacionRemitente=&identificacionConsignatario=&facturaRemito=factura&facturaRemitoTexto="

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add custom headers
	req.Header.Set("Authorization", "Basic QUxFUlRSQU46Q09OU1VMVEFTLjEyMw==N")

	// Send the request using http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

//\u003ca href=\"javascript:verExpedicion(\u0026#39;27033228\u0026#39;);\"\u003e999025836102\u003c/a\u003e
