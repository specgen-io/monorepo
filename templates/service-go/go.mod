module {{project.value}}

go {{versions.go.value}}

require (
	cloud.google.com/go v{{versions.civil.value}}
	github.com/google/uuid v{{versions.uuid.value}}
	github.com/husobee/vestigo v{{versions.vestigo.value}}
	github.com/shopspring/decimal v{{versions.decimal.value}}
	github.com/sirupsen/logrus v{{versions.logrus.value}}
	github.com/stretchr/testify v{{versions.testify.value}} // indirect
)
