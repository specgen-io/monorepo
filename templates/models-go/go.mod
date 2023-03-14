module {{project.value}}

go {{versions.go.value}}

require (
	cloud.google.com/go v{{versions.civil.value}}
	github.com/google/go-cmp v{{versions.go-cmp.value}} // indirect
	github.com/google/uuid v{{versions.uuid.value}}
	github.com/shopspring/decimal v{{versions.decimal.value}}
)
