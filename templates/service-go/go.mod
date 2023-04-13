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
	github.com/shopspring/decimal v{{versions.decimal.value}}
	github.com/sirupsen/logrus v{{versions.logrus.value}}
	github.com/stretchr/testify v{{versions.testify.value}} // indirect
)
