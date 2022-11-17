//go:generate specgen-golang service-go --spec-file spec.yaml --module-name {{project.value}} --generate-path ./spec --services-path ./services {{#swagger.value}}--swagger-path docs/swagger.yaml {{/swagger.value}}

package main

import (
	"flag"
	"github.com/husobee/vestigo"
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

	router := vestigo.NewRouter()

	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin: []string{"*", "*"},
	})

	sampleService := &services.SampleService{}

	spec.AddRoutes(router, sampleService)

	router.Get("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))).ServeHTTP)

	log.Infof("Starting service on port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
