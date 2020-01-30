package spec

import (
	"fmt"

	"github.com/anothrnick/machinable/dsi/models"
)

func injectProjectSchema(spec *ProjectSpec, resources []*models.ResourceDefinition) {
	// update spec info
	for _, resource := range resources {
		tag := Tag{
			Name:        resource.Title,
			Description: resource.Title,
		}
		spec.Tags = append(spec.Tags, tag)
		groupTags := spec.XTagGroups[1].Tags
		spec.XTagGroups[1].Tags = append(groupTags, tag.Name)
		injectPaths(spec, resource)
		injectComponents(spec, resource)
	}
}

func injectComponents(spec *ProjectSpec, resource *models.ResourceDefinition) {
	schema, _ := resource.GetSchemaMap()
	// TODO: Title could have spaces?
	spec.Components.Schemas[resource.Title] = schema
	spec.Components.Responses[fmt.Sprintf("%sList", resource.Title)] = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"count": map[string]interface{}{
				"description": "Total item count",
				"type":        "integer",
			},
			"links": map[string]interface{}{
				"description": "Absolute pagination links",
				"type":        "object",
				"properties": map[string]interface{}{
					"self": map[string]interface{}{
						"type": "string",
					},
					"next": map[string]interface{}{
						"type": "string",
					},
					"prev": map[string]interface{}{
						"type": "string",
					},
				},
			},
			"items": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"$ref": fmt.Sprintf("#/components/schemas/%s", resource.Title),
				},
			},
		},
	}
}

func injectPaths(spec *ProjectSpec, resource *models.ResourceDefinition) {
	componentLink := fmt.Sprintf("#/components/schemas/%s", resource.Title)
	listLink := fmt.Sprintf("#/components/responses/%sList", resource.Title)
	paths := map[string]map[string]Verb{
		fmt.Sprintf("/api/%s/{%sId}", resource.PathName, resource.PathName): {
			"get": {
				Tags:        []string{resource.Title},
				Summary:     fmt.Sprintf("Get %s", resource.Title),
				OperationID: fmt.Sprintf("Get%s", resource.Title),
				Security: []map[string][]interface{}{
					{
						"JWT": []interface{}{},
					},
				},
				Responses: map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Resource retrieved successfully",
						"headers":     map[string]interface{}{},
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": componentLink,
								},
							},
						},
					},
					"400": map[string]interface{}{
						"$ref": "#/components/responses/BadRequest",
					},
					"401": map[string]interface{}{
						"$ref": "#/components/responses/UnauthorizedError",
					},
					"404": map[string]interface{}{
						"$ref": "#/components/responses/NotFound",
					},
				},
			},
			"put": {
				Tags:        []string{resource.Title},
				Summary:     fmt.Sprintf("Update %s", resource.Title),
				OperationID: fmt.Sprintf("Update%s", resource.Title),
				Security: []map[string][]interface{}{
					{
						"JWT": []interface{}{},
					},
				},
				RequestBody: map[string]interface{}{
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": componentLink,
							},
						},
					},
				},
				Responses: map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Resource created successfully",
						"headers":     map[string]interface{}{},
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": componentLink,
								},
							},
						},
					},
					"400": map[string]interface{}{
						"$ref": "#/components/responses/BadRequest",
					},
					"401": map[string]interface{}{
						"$ref": "#/components/responses/UnauthorizedError",
					},
					"404": map[string]interface{}{
						"$ref": "#/components/responses/NotFound",
					},
				},
			},
			"delete": {
				Tags:        []string{resource.Title},
				Summary:     fmt.Sprintf("Delete %s", resource.Title),
				OperationID: fmt.Sprintf("Delete%s", resource.Title),
				Security:    []map[string][]interface{}{},
				Responses: map[string]interface{}{
					"204": map[string]interface{}{
						"description": "Resource was deleted successfully",
						"headers":     map[string]interface{}{},
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
								},
							},
						},
					},
					"401": map[string]interface{}{
						"$ref": "#/components/responses/UnauthorizedError",
					},
					"404": map[string]interface{}{
						"$ref": "#/components/responses/NotFound",
					},
				},
			},
		},
		fmt.Sprintf("/api/%s", resource.PathName): {
			"get": {
				Tags:        []string{resource.Title},
				Summary:     fmt.Sprintf("List %s", resource.Title),
				OperationID: fmt.Sprintf("List%s", resource.Title),
				Security: []map[string][]interface{}{
					{
						"JWT": []interface{}{},
					},
				},
				Responses: map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Resource list retrieved successfully",
						"headers":     map[string]interface{}{},
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": listLink,
								},
							},
						},
					},
					"400": map[string]interface{}{
						"$ref": "#/components/responses/BadRequest",
					},
					"401": map[string]interface{}{
						"$ref": "#/components/responses/UnauthorizedError",
					},
					"404": map[string]interface{}{
						"$ref": "#/components/responses/NotFound",
					},
				},
			},
			"post": {
				Tags:        []string{resource.Title},
				Summary:     fmt.Sprintf("Create %s", resource.Title),
				OperationID: fmt.Sprintf("Create%s", resource.Title),
				Security: []map[string][]interface{}{
					{
						"JWT": []interface{}{},
					},
				},
				RequestBody: map[string]interface{}{
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": componentLink,
							},
						},
					},
				},
				Responses: map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Resource created successfully",
						"headers":     map[string]interface{}{},
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": componentLink,
								},
							},
						},
					},
					"400": map[string]interface{}{
						"$ref": "#/components/responses/BadRequest",
					},
					"401": map[string]interface{}{
						"$ref": "#/components/responses/UnauthorizedError",
					},
					"404": map[string]interface{}{
						"$ref": "#/components/responses/NotFound",
					},
				},
			},
		},
	}

	for key, val := range paths {
		spec.Paths[key] = val
	}
}

func baseSpec(projectPath string) *ProjectSpec {
	return &ProjectSpec{
		Openapi: "3.0.0",
		Info: Info{
			Title: "Sample Machinable Project API",
			Contact: Contact{
				Name:  "Machinable Support",
				URL:   "https://www.machinable.io/",
				Email: "support@machinable.io",
			},
			XLogo: XLogo{
				URL:             "test",
				BackgroundColor: "#fafafa",
				AltText:         projectPath,
			},
			Description: "Sample OpenAPI specification for Machinable project.",
		},
		Servers: []Server{
			{
				URL:         fmt.Sprintf("https://%s.machinable.io", projectPath),
				Description: "Live Server",
			},
		},
		Tags: []Tag{
			{
				Name:        "JWT Session",
				Description: "User session with JWT",
			},
		},
		XTagGroups: []XTagGroup{
			{
				Name: "Security",
				Tags: []string{
					"JWT Session",
				},
			},
			{
				Name: "API Resources",
				Tags: []string{},
			},
		},
		Security: []map[string][]interface{}{
			{
				"APIKey": []interface{}{},
			},
			{
				"JWT": []interface{}{},
			},
		},
		Paths: map[string]map[string]Verb{
			"/sessions": {
				"post": Verb{
					Tags:        []string{"JWT Session"},
					Summary:     "Create new session",
					OperationID: "CreateSession",
					Security: []map[string][]interface{}{
						{
							"BasicAuth": []interface{}{},
						},
					},
					Responses: map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Session was created successfully",
							"headers":     map[string]interface{}{},
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"access_token": map[string]interface{}{
												"description": "Access token which is used to authenticate future requests. The `access_token` has an expiration of 5 minutes.",
												"type":        "string",
											},
											"refresh_token": map[string]interface{}{
												"description": "Refresh token which is used to retrieve a new `access_token`.",
												"type":        "string",
											},
											"session_id": map[string]interface{}{
												"description": "The `ID` of this session.",
												"type":        "string",
											},
										},
									},
								},
							},
						},
						"401": map[string]interface{}{
							"$ref": "#/components/responses/UnauthorizedError",
						},
						"404": map[string]interface{}{
							"$ref": "#/components/responses/NotFound",
						},
					},
					CodeSamples: []CodeSample{
						{
							Lang:   "bash",
							Source: fmt.Sprintf("# base64 encode username|password to make HTTP Basic authn request\n$ echo \"testUser:hunter2\" | base64\ndGVzdFVzZXI6aHVudGVyMgo=\n\n# POST credentials to /sessions/ endpoint to recieve access token\n$ curl -X POST \\\n\thttps://%s.machinable.io/sessions/ \\\n\t-H 'authorization: Basic dGVzdFVzZXI6aHVudGVyMg=='", projectPath),
						},
					},
				},
			},
			"/sessions/refresh": {
				"post": Verb{
					Tags:        []string{"JWT Session"},
					Summary:     "Refresh current session",
					OperationID: "RefreshSession",
					Security: []map[string][]interface{}{
						{
							"JWT": []interface{}{},
						},
					},
					Responses: map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Session was refreshed successfully",
							"headers":     map[string]interface{}{},
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"access_token": map[string]interface{}{
												"description": "The new access token which is used to authenticate future requests. The `access_token` has an expiration of 5 minutes.",
												"type":        "string",
											},
										},
									},
								},
							},
						},
						"401": map[string]interface{}{
							"$ref": "#/components/responses/UnauthorizedError",
						},
						"404": map[string]interface{}{
							"$ref": "#/components/responses/NotFound",
						},
					},
					CodeSamples: []CodeSample{
						{
							Lang:   "bash",
							Source: fmt.Sprintf("curl -X POST \\\n https://%s.machinable.io/sessions/refresh/ \\\n -H 'authorization: Bearer {refresh_token}'", projectPath),
						},
					},
				},
			},
			"/sessions/{sessionId}": {
				"delete": Verb{
					Tags:        []string{"JWT Session"},
					Summary:     "Delete a session",
					OperationID: "DeleteSession",
					Security:    []map[string][]interface{}{},
					Responses: map[string]interface{}{
						"204": map[string]interface{}{
							"description": "Session was deleted successfully",
							"headers":     map[string]interface{}{},
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
									},
								},
							},
						},
						"401": map[string]interface{}{
							"$ref": "#/components/responses/UnauthorizedError",
						},
						"404": map[string]interface{}{
							"$ref": "#/components/responses/NotFound",
						},
					},
				},
			},
		},
		Components: Components{
			Schemas: map[string]interface{}{
				"Error": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"error": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
			Responses: map[string]interface{}{
				"NotFound": map[string]interface{}{
					"description": "Resource was not found",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
				"UnauthorizedError": map[string]interface{}{
					"description": "API Key or JWT is missing or invalid",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
				"BadRequest": map[string]interface{}{
					"description": "Invalid parameters",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
				"ServerError": map[string]interface{}{
					"description": "Unknown server error occurred",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/Error",
							},
						},
					},
				},
			},
			SecuritySchemes: map[string]SecurityScheme{
				"JWT": SecurityScheme{
					Description:  "You can create a JSON Web Token (JWT) via our [JWT Session resource](#tag/JWT-Session).\nUsage format: `Bearer <JWT>`\n",
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
				"ApiKey": SecurityScheme{
					Description: "API Keys can be created from the [project dashboard](https://www.machinable.io/documentation/projects/access/#api-keys).\nUsage format: `apikey <API Key>`\n",
					Name:        "Authorization",
					Type:        "apiKey",
					In:          "header",
				},
				"BasicAuth": SecurityScheme{
					Description: "Basic authentication is used to acquire a new JWT session.",
					Type:        "http",
					Scheme:      "basic",
				},
			},
		},
	}
}

type ProjectSpec struct {
	Openapi    string                     `json:"openapi"`
	Info       Info                       `json:"info"`
	Servers    []Server                   `json:"servers"`
	Tags       []Tag                      `json:"tags"`
	XTagGroups []XTagGroup                `json:"x-tagGroups"`
	Security   []map[string][]interface{} `json:"security"`
	Paths      map[string]map[string]Verb `json:"paths"`
	Components Components                 `json:"components"`
}
type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}
type XLogo struct {
	URL             string `json:"url"`
	BackgroundColor string `json:"backgroundColor"`
	AltText         string `json:"altText"`
}
type Info struct {
	Title       string  `json:"title"`
	Contact     Contact `json:"contact"`
	XLogo       XLogo   `json:"x-logo"`
	Description string  `json:"description"`
}
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type XTagGroup struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
type Verb struct {
	Tags        []string                   `json:"tags"`
	Summary     string                     `json:"summary"`
	OperationID string                     `json:"operationId"`
	Security    []map[string][]interface{} `json:"security"`
	RequestBody map[string]interface{}     `json:"requestBody"`
	Responses   map[string]interface{}     `json:"responses,omitempty"`
	CodeSamples []CodeSample               `json:"x-code-samples,omitempty"`
}
type CodeSample struct {
	Lang   string `json:"lang"`
	Source string `json:"source"`
}
type Schema struct {
}
type Response struct {
}
type SecurityScheme struct {
	Description  string `json:"description,omitempty"`
	Type         string `json:"type,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	In           string `json:"in,omitempty"`
	Name         string `json:"name,omitempty"`
}
type Components struct {
	Schemas         map[string]interface{}    `json:"schemas"`
	Responses       map[string]interface{}    `json:"responses"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
}
