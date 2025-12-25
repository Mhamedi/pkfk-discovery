package http

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
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
		Paths: &openapi3.Paths{},
	}

	// Define common schemas using fluent API
	spec.Components.Schemas["Error"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("error", openapi3.NewStringSchema()),
	}

	spec.Components.Schemas["UUID"] = &openapi3.SchemaRef{
		Value: openapi3.NewStringSchema().WithFormat("uuid"),
	}

	spec.Components.Schemas["Timestamp"] = &openapi3.SchemaRef{
		Value: openapi3.NewStringSchema().WithFormat("date-time"),
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
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("name", openapi3.NewStringSchema()).
			WithProperty("vendor", openapi3.NewStringSchema()).
			WithProperty("db_family", openapi3.NewStringSchema()).
			WithProperty("version", openapi3.NewStringSchema()).
			WithProperty("maturity_level", openapi3.NewStringSchema().WithEnum("L0", "L1", "L2", "L3", "L4")).
			WithProperty("bundle_path", openapi3.NewStringSchema()).
			WithProperty("signature", openapi3.NewStringSchema()).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithRequired("id", "name", "vendor", "db_family", "version", "maturity_level"),
	}
}

func defineAdapterDraftSchema(spec *openapi3.T) {
	spec.Components.Schemas["AdapterDraft"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("name", openapi3.NewStringSchema()).
			WithProperty("status", openapi3.NewStringSchema()).
			WithProperty("adapter_id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("created_by", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}),
	}
}

func defineConnectionSchema(spec *openapi3.T) {
	spec.Components.Schemas["Connection"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("name", openapi3.NewStringSchema()).
			WithProperty("db_type", openapi3.NewStringSchema()).
			WithProperty("host", openapi3.NewStringSchema()).
			WithProperty("port", openapi3.NewInt64Schema()).
			WithProperty("database", openapi3.NewStringSchema()).
			WithProperty("username", openapi3.NewStringSchema()).
			WithProperty("created_by", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithRequired("id", "name", "db_type", "host", "port", "database", "username"),
	}
}

func defineScanSchema(spec *openapi3.T) {
	spec.Components.Schemas["Scan"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("connection_id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("adapter_id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("status", openapi3.NewStringSchema().WithEnum("pending", "running", "completed", "failed", "cancelled")).
			WithProperty("policy", openapi3.NewObjectSchema()).
			WithProperty("results", openapi3.NewObjectSchema()).
			WithProperty("created_by", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}),
	}
}

func defineUserSchema(spec *openapi3.T) {
	spec.Components.Schemas["User"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("username", openapi3.NewStringSchema()).
			WithProperty("email", openapi3.NewStringSchema().WithFormat("email")).
			WithProperty("role", openapi3.NewStringSchema().WithEnum("admin", "editor", "viewer")).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}),
	}
}

func defineAIProviderSchema(spec *openapi3.T) {
	spec.Components.Schemas["AIProvider"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("name", openapi3.NewStringSchema()).
			WithProperty("type", openapi3.NewStringSchema().WithEnum("local", "cloud")).
			WithProperty("endpoint", openapi3.NewStringSchema()).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}).
			WithProperty("updated_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}),
	}
}

func defineAuditLogSchema(spec *openapi3.T) {
	spec.Components.Schemas["AuditLog"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("user_id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("action", openapi3.NewStringSchema()).
			WithProperty("resource_type", openapi3.NewStringSchema()).
			WithProperty("resource_id", &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"}).
			WithProperty("details", openapi3.NewObjectSchema()).
			WithProperty("ip_address", openapi3.NewStringSchema()).
			WithProperty("created_at", &openapi3.SchemaRef{Ref: "#/components/schemas/Timestamp"}),
	}
}

func defineSettingsSchema(spec *openapi3.T) {
	spec.Components.Schemas["Settings"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("key", openapi3.NewStringSchema()).
			WithProperty("value", openapi3.NewObjectSchema()),
	}
}

func defineLoginRequestSchema(spec *openapi3.T) {
	spec.Components.Schemas["LoginRequest"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("username", openapi3.NewStringSchema()).
			WithProperty("password", openapi3.NewStringSchema().WithFormat("password")).
			WithRequired("username", "password"),
	}
}

func defineLoginResponseSchema(spec *openapi3.T) {
	spec.Components.Schemas["LoginResponse"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("access_token", openapi3.NewStringSchema()).
			WithProperty("refresh_token", openapi3.NewStringSchema()).
			WithProperty("expires_in", openapi3.NewInt64Schema()),
	}
}

func defineValidationResultSchema(spec *openapi3.T) {
	spec.Components.Schemas["ValidationResult"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("level", openapi3.NewStringSchema()).
			WithProperty("passed", openapi3.NewBoolSchema()).
			WithProperty("tests", openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema())).
			WithProperty("errors", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
			WithProperty("warnings", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
	}
}

func defineOptimizationResponseSchema(spec *openapi3.T) {
	spec.Components.Schemas["OptimizationResponse"] = &openapi3.SchemaRef{
		Value: openapi3.NewObjectSchema().
			WithProperty("diff_patch", openapi3.NewStringSchema()).
			WithProperty("explanation", openapi3.NewStringSchema()),
	}
}

func defineAuthPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/auth/login"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:     "Login",
			Description: "Authenticate user and receive JWT tokens",
			Tags:        []string{"Authentication"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Ref: "#/components/schemas/LoginRequest"},
						},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/LoginResponse"}),
				}),
				openapi3.WithStatus(401, &openapi3.Response{
					Description: openapi3.Ptr("Unauthorized"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Error"}),
				}),
			),
		},
	}
}

func defineAdapterPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/adapters"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Adapters",
			Tags:     []string{"Adapters"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
					}),
				}),
			),
		},
		Post: &openapi3.Operation{
			Summary:  "Create Adapter",
			Tags:     []string{"Adapters"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/adapters/{id}"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "Get Adapter",
			Tags:     []string{"Adapters"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				}),
			),
		},
		Put: &openapi3.Operation{
			Summary:  "Update Adapter",
			Tags:     []string{"Adapters"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				}),
			),
		},
		Delete: &openapi3.Operation{
			Summary:  "Delete Adapter",
			Tags:     []string{"Adapters"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(204, &openapi3.Response{
					Description: openapi3.Ptr("No Content"),
				}),
			),
		},
	}
}

func defineStudioPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/studio/drafts"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:  "Create Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/studio/drafts/{id}"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "Get Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"}),
				}),
			),
		},
		Put: &openapi3.Operation{
			Summary:  "Update Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AdapterDraft"}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/studio/drafts/{id}/validate"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:  "Validate Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/ValidationResult"}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/studio/drafts/{id}/optimize"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:  "Optimize Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/OptimizationResponse"}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/studio/drafts/{id}/publish"] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Summary:  "Publish Draft",
			Tags:     []string{"Studio"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     "id",
						In:       "path",
						Required: true,
						Schema:   &openapi3.SchemaRef{Ref: "#/components/schemas/UUID"},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Adapter"}),
				}),
			),
		},
	}
}

func defineConnectionPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/connections"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Connections",
			Tags:     []string{"Connections"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/Connection"}),
					}),
				}),
			),
		},
		Post: &openapi3.Operation{
			Summary:  "Create Connection",
			Tags:     []string{"Connections"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Connection"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Connection"}),
				}),
			),
		},
	}
}

func defineScanPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/scans"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Scans",
			Tags:     []string{"Scans"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/Scan"}),
					}),
				}),
			),
		},
		Post: &openapi3.Operation{
			Summary:  "Create Scan",
			Tags:     []string{"Scans"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Scan"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Scan"}),
				}),
			),
		},
	}
}

func defineAIProviderPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/ai/providers"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List AI Providers",
			Tags:     []string{"AI"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"}),
					}),
				}),
			),
		},
		Post: &openapi3.Operation{
			Summary:  "Create AI Provider",
			Tags:     []string{"AI"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(201, &openapi3.Response{
					Description: openapi3.Ptr("Created"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/AIProvider"}),
				}),
			),
		},
	}
}

func defineAdminPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/admin/users"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Users",
			Tags:     []string{"Admin"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/User"}),
					}),
				}),
			),
		},
	}

	spec.Paths.Value["/api/v1/admin/audit"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "List Audit Logs",
			Tags:     []string{"Admin"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{
						Value: openapi3.NewArraySchema().WithItems(&openapi3.SchemaRef{Ref: "#/components/schemas/AuditLog"}),
					}),
				}),
			),
		},
	}
}

func defineSettingsPaths(spec *openapi3.T) {
	spec.Paths.Value["/api/v1/settings"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:  "Get Settings",
			Tags:     []string{"Settings"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Settings"}),
				}),
			),
		},
		Put: &openapi3.Operation{
			Summary:  "Update Settings",
			Tags:     []string{"Settings"},
			Security: &openapi3.SecurityRequirements{{"BearerAuth": []string{}}},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: openapi3.BoolPtr(true),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Settings"}),
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("Success"),
					Content: openapi3.NewContentWithJSONSchemaRef(&openapi3.SchemaRef{Ref: "#/components/schemas/Settings"}),
				}),
			),
		},
	}
}

func defineHealthPaths(spec *openapi3.T) {
	spec.Paths.Value["/healthz"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Health Check",
			Description: "Liveness probe endpoint",
			Tags:        []string{"Health"},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("OK"),
				}),
			),
		},
	}

	spec.Paths.Value["/readyz"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Summary:     "Readiness Check",
			Description: "Readiness probe endpoint",
			Tags:        []string{"Health"},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.Response{
					Description: openapi3.Ptr("OK"),
				}),
			),
		},
	}
}

