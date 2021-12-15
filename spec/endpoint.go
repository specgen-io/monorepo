package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"regexp"
	"strings"
)

type UrlPart struct {
	Part  string
	Param *NamedParam
}

type Endpoint struct {
	Method    string
	Url       string
	UrlParams UrlParams
	UrlParts  []UrlPart
}

func (value Endpoint) MarshalYAML() (interface{}, error) {
	url := value.Url
	for _, param := range value.UrlParams {
		paramStr := fmt.Sprintf(`{%s:%s}`, param.Name.Source, param.Type.Definition.String())
		url = strings.Replace(url, UrlParamStr(&param), paramStr, 1)
	}
	return fmt.Sprintf(`%s %s`, value.Method, url), nil
}

func (value *Endpoint) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return yamlError(node, "operation endpoint should be string")
	}
	endpoint, err := parseEndpoint(node.Value, node)
	if err != nil {
		return yamlError(node, err.Error())
	}
	*value = *endpoint
	return nil
}

func parseEndpoint(endpoint string, node *yaml.Node) (*Endpoint, error) {
	spaces_count := strings.Count(endpoint, " ")
	if spaces_count != 1 {
		return nil, errors.New("endpoint should be in format 'METHOD url'")
	}
	endpointParts := strings.SplitN(endpoint, " ", 2)
	method := endpointParts[0]
	err := HttpMethod.Check(method)
	if err != nil {
		return nil, err
	}
	url := endpointParts[1]
	cleanUrl, params, err := getCleanUrl(url, node)
	if err != nil {
		return nil, err
	}
	urlParts := getUrlParts(cleanUrl, params)
	return &Endpoint{method, cleanUrl, params, urlParts}, nil
}

func parseUrlParam(paramStr string, node *yaml.Node) (*NamedParam, error) {
	paramStr = strings.Replace(paramStr, "{", "", -1)
	paramStr = strings.Replace(paramStr, "}", "", -1)
	paramParts := strings.Split(paramStr, ":")
	paramName := strings.TrimSpace(paramParts[0])
	paramType := strings.TrimSpace(paramParts[1])

	typ, err := parseType(paramType)
	if err != nil {
		return nil, err
	}

	return &NamedParam{
		Name: Name{Source: paramName, Location: node},
		DefinitionDefault: DefinitionDefault{
			Type:     Type{*typ, node},
			Location: node,
		},
	}, nil
}

func getUrlParts(url string, params UrlParams) []UrlPart {
	urlParts := []UrlPart{}
	reminder := url
	for i, param := range params {
		paramStr := UrlParamStr(&param)
		parts := strings.SplitN(reminder, paramStr, 2)
		urlParts = append(urlParts, UrlPart{Part: parts[0]})
		urlParts = append(urlParts, UrlPart{Part: paramStr, Param: &params[i]})
		reminder = parts[1]
	}
	if reminder != "" {
		urlParts = append(urlParts, UrlPart{Part: reminder})
	}
	return urlParts
}

func getCleanUrl(url string, node *yaml.Node) (string, UrlParams, error) {
	re := regexp.MustCompile(`\{[a-zA-Z][a-zA-Z0-9-_]*:[a-zA-Z][a-zA-Z0-9-_]*\}`)
	matches := re.FindAllStringIndex(url, -1)
	params := UrlParams{}
	cleanUrl := url
	for _, match := range matches {
		paramStr := url[match[0]:match[1]]
		param, err := parseUrlParam(paramStr, node)
		if err != nil {
			return "", nil, err
		}
		params = append(params, *param)
		cleanUrl = strings.Replace(cleanUrl, paramStr, UrlParamStr(param), 1)
	}
	return cleanUrl, params, nil
}

func UrlParamStr(param *NamedParam) string {
	return "{" + param.Name.Source + "}"
}
