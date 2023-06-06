//go:generate specgen-golang service-go --jsonmode {{jsonmode.value}} --server chi --spec-file spec.yaml --module-name {{project.value}} --generate-path ./spec --services-path ./services {{#swagger.value}}--swagger-path docs/swagger.yaml {{/swagger.value}}

package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	{{#cors.value}}
	"github.com/go-chi/cors"
	{{/cors.value}}
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

	router := chi.NewRouter()
	
	{{#cors.value}}
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*", "*"},
	}))
	{{/cors.value}}

	spec.AddRoutes(router, services.Create())

	{{#swagger.value}}
	router.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	{{/swagger.value}}

	log.Infof("Starting service on port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
