package http

import (
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
)

func GenerateOpenAPISpec() *openapi3.T {
	spec := &openapi3.T{
		OpenAPI: "3.1.0",
		Info: &openapi3.Info{
			Title:       "Adapter Studio API",
			Description: "Enterprise-grade Adapter Studio + Registry + Engine Admin REST API",
			Version:     "1.0.0",
			Contact: &openapi3.Contact{
				Name: "PK-FK Discovery",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL:         "http://localhost:8080",
				Description: "Local development server",
			},
		},
		Components: &openapi3.Components{
			Schemas: make(openapi3.Schemas),
			SecuritySchemes: openapi3.SecuritySchemes{
				"BearerAuth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
				},
			},
		},
		Paths: make(openapi3.Paths),
	}

	// Define common schemas
	spec.Components.Schemas["Error"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"error": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
			},
		},
	}

	spec.Components.Schemas["UUID"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:   "string",
			Format: "uuid",
		},
	}

	spec.Components.Schemas["Timestamp"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:   "string",
			Format: "date-time",
		},
	}

	// Define domain schemas
	defineAdapterSchema(spec)
	defineAdapterDraftSchema(spec)
	defineConnectionSchema(spec)
	defineScanSchema(spec)
	defineUserSchema(spec)
	defineAIProviderSchema(spec)
	defineAuditLogSchema(spec)
	defineSettingsSchema(spec)
	defineLoginRequestSchema(spec)
	defineLoginResponseSchema(spec)
	defineValidationResultSchema(spec)
	defineOptimizationResponseSchema(spec)

	// Define paths
	defineAuthPaths(spec)
	defineAdapterPaths(spec)
	defineStudioPaths(spec)
	defineConnectionPaths(spec)
	defineScanPaths(spec)
	defineAIProviderPaths(spec)
	defineAdminPaths(spec)
	defineSettingsPaths(spec)
	defineHealthPaths(spec)

	return spec
}

func defineAdapterSchema(spec *openapi3.T) {
	spec.Components.Schemas["Adapter"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":            &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"name":          &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"vendor":        &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"db_family":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"version":       &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"maturity_level": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"L0", "L1", "L2", "L3", "L4"}}},
				"bundle_path":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"signature":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"created_at":    &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at":    &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
			Required: []string{"id", "name", "vendor", "db_family", "version", "maturity_level"},
		},
	}
}

func defineAdapterDraftSchema(spec *openapi3.T) {
	spec.Components.Schemas["AdapterDraft"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":          &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"name":        &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"status":      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"adapter_id":  &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"created_by":  &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"created_at":  &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at":  &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
		},
	}
}

func defineConnectionSchema(spec *openapi3.T) {
	spec.Components.Schemas["Connection"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":       &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"name":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"db_type":  &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"host":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"port":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
				"database": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"username": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"created_by": &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"created_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
			Required: []string{"id", "name", "db_type", "host", "port", "database", "username"},
		},
	}
}

func defineScanSchema(spec *openapi3.T) {
	spec.Components.Schemas["Scan"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":           &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"connection_id": &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"adapter_id":   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"status":       &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"pending", "running", "completed", "failed", "cancelled"}}},
				"policy":       &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object"}},
				"results":      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object"}},
				"created_by":   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"created_at":   &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at":   &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
		},
	}
}

func defineUserSchema(spec *openapi3.T) {
	spec.Components.Schemas["User"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":        &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"username":  &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"email":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Format: "email"}},
				"role":      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"admin", "editor", "viewer"}}},
				"created_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
		},
	}
}

func defineAIProviderSchema(spec *openapi3.T) {
	spec.Components.Schemas["AIProvider"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":       &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"name":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"type":     &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"local", "cloud"}}},
				"endpoint": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"created_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
				"updated_at": &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
		},
	}
}

func defineAuditLogSchema(spec *openapi3.T) {
	spec.Components.Schemas["AuditLog"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"id":            &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"user_id":       &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"action":        &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"resource_type": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"resource_id":   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
				"details":       &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object"}},
				"ip_address":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"created_at":    &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"},
			},
		},
	}
}

func defineSettingsSchema(spec *openapi3.T) {
	spec.Components.Schemas["Settings"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"key":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"value": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object"}},
			},
		},
	}
}

func defineLoginRequestSchema(spec *openapi3.T) {
	spec.Components.Schemas["LoginRequest"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"username": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"password": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string", Format: "password"}},
			},
			Required: []string{"username", "password"},
		},
	}
}

func defineLoginResponseSchema(spec *openapi3.T) {
	spec.Components.Schemas["LoginResponse"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"access_token":  &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"refresh_token": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"expires_in":    &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
			},
		},
	}
}

func defineValidationResultSchema(spec *openapi3.T) {
	spec.Components.Schemas["ValidationResult"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"level":    &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"passed":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "boolean"}},
				"tests":    &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "array", Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "object"}}}},
				"errors":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "array", Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
				"warnings": &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "array", Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}}},
			},
		},
	}
}

func defineOptimizationResponseSchema(spec *openapi3.T) {
	spec.Components.Schemas["OptimizationResponse"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "object",
			Properties: openapi3.Schemas{
				"diff_patch":   &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				"explanation":  &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
			},
		},
	}
}

func defineAuthPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/auth/login"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Login",
			Description: "Authenticate user and receive JWT tokens",
			Tags:        []string{"Authentication"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/LoginRequest"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/LoginResponse"},
							},
						},
					},
				},
				"401": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Unauthorized"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Error"},
							},
						},
					},
				},
			},
		},
	}
}

func defineAdapterPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/adapters"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "List Adapters",
			Description: "List all adapters in the registry",
			Tags:        []string{"Adapters"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &openapi3.Operation{
			Summary:     "Create Adapter",
			Description: "Create a new adapter in the registry",
			Tags:        []string{"Adapters"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/adapters/{id}"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Get Adapter",
			Description: "Get adapter by ID",
			Tags:        []string{"Adapters"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
							},
						},
					},
				},
			},
		},
		Put: &openapi3.Operation{
			Summary:     "Update Adapter",
			Description: "Update an existing adapter",
			Tags:        []string{"Adapters"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
							},
						},
					},
				},
			},
		},
		Delete: &openapi3.Operation{
			Summary:     "Delete Adapter",
			Description: "Delete an adapter",
			Tags:        []string{"Adapters"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"204": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"No Content"}[0],
					},
				},
			},
		},
	}
}

func defineStudioPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/studio/drafts"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Create Draft",
			Description: "Create a new adapter draft",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/studio/drafts/{id}"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Get Draft",
			Description: "Get adapter draft by ID",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"},
							},
						},
					},
				},
			},
		},
		Put: &openapi3.Operation{
			Summary:     "Update Draft",
			Description: "Update an adapter draft",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/studio/drafts/{id}/validate"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Validate Draft",
			Description: "Run validation tests on an adapter draft",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/ValidationResult"},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/studio/drafts/{id}/optimize"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Optimize Draft",
			Description: "Use AI to optimize SQL templates in an adapter draft",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/OptimizationResponse"},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/studio/drafts/{id}/publish"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Publish Draft",
			Description: "Package and publish an adapter draft to the registry",
			Tags:        []string{"Studio"},
			Security:    &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: []*openapi3.ParameterRef{
				{Value: &openapi3.Parameter{Name: "id", In: "path", Required: true, Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}}},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"},
							},
						},
					},
				},
			},
		},
	}
}

func defineConnectionPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/connections"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Connections",
			Tags:     []string{"Connections"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/Connection"},
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &openapi3.Operation{
			Summary:  "Create Connection",
			Tags:     []string{"Connections"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Connection"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Connection"},
							},
						},
					},
				},
			},
		},
	}
}

func defineScanPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/scans"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Scans",
			Tags:     []string{"Scans"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/Scan"},
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &openapi3.Operation{
			Summary:  "Create Scan",
			Tags:     []string{"Scans"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Scan"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Scan"},
							},
						},
					},
				},
			},
		},
	}
}

func defineAIProviderPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/ai/providers"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List AI Providers",
			Tags:     []string{"AI"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"},
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &openapi3.Operation{
			Summary:  "Create AI Provider",
			Tags:     []string{"AI"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"201": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Created"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"},
							},
						},
					},
				},
			},
		},
	}
}

func defineAdminPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/admin/users"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Users",
			Tags:     []string{"Admin"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/User"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	spec.Paths["/api/v1/admin/audit"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Audit Logs",
			Tags:     []string{"Admin"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "array",
										Items: &openapi3.SchemaRef{Ref: "#/components/schemas/AuditLog"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func defineSettingsPaths(spec *openapi3.T) {
	spec.Paths["/api/v1/settings"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "Get Settings",
			Tags:     []string{"Settings"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Settings"},
							},
						},
					},
				},
			},
		},
		Put: &openapi3.Operation{
			Summary:  "Update Settings",
			Tags:     []string{"Settings"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Settings"},
						},
					},
				},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/Settings"},
							},
						},
					},
				},
			},
		},
	}
}

func defineHealthPaths(spec *openapi3.T) {
	spec.Paths["/healthz"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Health Check",
			Description: "Liveness probe endpoint",
			Tags:        []string{"Health"},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"OK"}[0],
					},
				},
			},
		},
	}

	spec.Paths["/readyz"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Readiness Check",
			Description: "Readiness probe endpoint",
			Tags:        []string{"Health"},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"OK"}[0],
					},
				},
			},
		},
	}
}

