package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateCheckService(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]]

import (
	"cloud.google.com/go/civil"
	"github.com/shopspring/decimal"
)

type CheckService struct{}

func (service *CheckService) CheckQuery(
	pString string,
	pStringOpt *string,
	pStringArray []string,
	pDate civil.Date,
	pDateArray []civil.Date,
	pDatetime civil.DateTime,
	pInt int,
	pLong int64,
	pDecimal decimal.Decimal,
	pEnum Choice,
	pStringDefaulted string) (*CheckQueryResponse, error) {

	return &CheckQueryResponse{Ok: &Empty}, nil
}

func (service *CheckService) CheckForbidden() (*CheckForbiddenResponse, error) {
	return &CheckForbiddenResponse{Forbidden: &Empty}, nil
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
