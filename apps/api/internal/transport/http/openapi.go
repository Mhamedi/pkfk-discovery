package http

import "github.com/getkin/kin-openapi/openapi3"

func GenerateOpenAPISpec() *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.1.0",
		Info: &openapi3.Info{
			Title:       "Adapter Studio API",
			Description: "Enterprise-grade Adapter Studio + Registry + Engine Admin REST API",
			Version:     "1.0.0",
		},
	}
}
