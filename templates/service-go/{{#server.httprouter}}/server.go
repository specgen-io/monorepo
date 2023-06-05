//go:generate specgen-golang service-go --jsonmode {{jsonmode.value}} --server httprouter --spec-file spec.yaml --module-name {{project.value}} --generate-path ./spec --services-path ./services {{#swagger.value}}--swagger-path docs/swagger.yaml {{/swagger.value}}

package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"net/http"
	"{{project.value}}/services"
	"{{project.value}}/spec"
)

func main() {
	port := flag.String("port", "8081", "port number")
	flag.Parse()

	decimal.MarshalJSONWithoutQuotes = true

	router := httprouter.New()
	
	{{#cors.value}}
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Origin", "*")
		}
	})
	{{/cors.value}}

	sampleService := &services.SampleService{}

	spec.AddRoutes(router, sampleService)

	{{#swagger.value}}
	router.ServeFiles("/docs/*filepath", http.Dir("./docs"))
	{{/swagger.value}}

	log.Infof("Starting service on port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
