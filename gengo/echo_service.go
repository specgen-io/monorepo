package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateEchoService(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]]

type EchoService struct{}

func (service *EchoService) EchoBody(body *Message) (*EchoBodyResponse, error) {
	return &EchoBodyResponse{Ok: body}, nil
}

func (service *EchoService) EchoQuery(intQuery int, stringQuery string) (*EchoQueryResponse, error) {
	return &EchoQueryResponse{Ok: &Message{IntField: intQuery, StringField: stringQuery}}, nil
}

func (service *EchoService) EchoHeader(intHeader int, stringHeader string) (*EchoHeaderResponse, error) {
	return &EchoHeaderResponse{Ok: &Message{IntField: intHeader, StringField: stringHeader}}, nil
}

func (service *EchoService) EchoUrlParams(intUrl int, stringUrl string) (*EchoUrlParamsResponse, error) {
	return &EchoUrlParamsResponse{Ok: &Message{IntField: intUrl, StringField: stringUrl}}, nil
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}

func generateEchoServiceV2(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]]

type EchoService struct{}

func (service *EchoService) EchoBody(body *Message) (*EchoBodyResponse, error) {
	return &EchoBodyResponse{Ok: body}, nil
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
