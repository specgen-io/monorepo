//go:generate specgen-golang service-go --jsonmode strict --server vestigo --spec-file spec.yaml --module-name the-service --generate-path ./spec --services-path ./services --swagger-path docs/swagger.yaml

package main

import (
	"flag"
	"net/http"
	"the-service/services"
	"the-service/spec"

	"github.com/husobee/vestigo"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := flag.String("port", "8081", "port number")
	flag.Parse()

	decimal.MarshalJSONWithoutQuotes = true

	router := vestigo.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin: []string{"*", "*"},
	})

	spec.AddRoutes(router, services.Create())

	router.Get("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))).ServeHTTP)

	log.Infof("Starting service on port: %s", *port)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
