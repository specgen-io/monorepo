module {{project.value}}

go {{versions.go.value}}

require (
	cloud.google.com/go v{{versions.civil.value}}
	github.com/google/uuid v{{versions.uuid.value}}
	{{#server.vestigo}}
	github.com/husobee/vestigo v{{versions.vestigo.value}}
	{{/server.vestigo}}
	{{#server.httprouter}}
	github.com/julienschmidt/httprouter v{{versions.httprouter.value}}
	{{/server.httprouter}}
	{{#server.chi}}
	github.com/go-chi/chi/v5 v{{versions.chi.value}}
	{{/server.chi}}
	github.com/shopspring/decimal v{{versions.decimal.value}}
	github.com/sirupsen/logrus v{{versions.logrus.value}}
	github.com/stretchr/testify v{{versions.testify.value}} // indirect
	{{#cors.value}}
	{{#server.chi}}
	github.com/go-chi/cors v{{versions.chi-cors.value}} // indirect
	{{/server.chi}}
	{{/cors.value}}
)
