// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/matrix/create": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "matrix"
                ],
                "parameters": [
                    {
                        "description": "Matrix create",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/swagger.MatrixBase"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully created matrix",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/matrix/get_difference": {
            "get": {
                "description": "Retrieves the differences between two matrices identified by their names.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "matrix"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the first matrix",
                        "name": "from_name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the second matrix",
                        "name": "to_name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Found matrices differences",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.MatrixDifference"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input, missing matrix names",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/matrix/get_history": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "matrix"
                ],
                "parameters": [
                    {
                        "description": "Get data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/swagger.GetHistoryMatrix"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Found matrixes",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/swagger.ResponseHistoryMatrix"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.MatrixDifference": {
            "type": "object",
            "properties": {
                "added": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MatrixNode"
                    }
                },
                "deleted": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MatrixNode"
                    }
                },
                "updated": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/models.MatrixNode"
                        }
                    }
                }
            }
        },
        "models.MatrixNode": {
            "type": "object",
            "properties": {
                "microcategory_id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "region_id": {
                    "type": "integer"
                }
            }
        },
        "swagger.GetHistoryMatrix": {
            "type": "object",
            "properties": {
                "is_baseline": {
                    "type": "boolean"
                },
                "time_end": {
                    "type": "string"
                },
                "time_start": {
                    "type": "string"
                }
            }
        },
        "swagger.MatrixBase": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.MatrixNode"
                    }
                },
                "is_baseline": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "parent_name": {
                    "type": "string"
                }
            }
        },
        "swagger.ResponseHistoryMatrix": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "parent_name": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
