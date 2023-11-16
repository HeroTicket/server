// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users/create-tba": {
            "post": {
                "responses": {}
            }
        },
        "/users/issued-tickets": {
            "get": {
                "responses": {}
            }
        },
        "/users/login-callback": {
            "post": {
                "description": "processes login callback",
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "processes login callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session id",
                        "name": "sessionId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    }
                }
            }
        },
        "/users/login-qr": {
            "get": {
                "description": "returns login qr code",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "returns login qr code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session id",
                        "name": "sessionId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/protocol.AuthorizationRequestMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    }
                }
            }
        },
        "/users/logout": {
            "post": {
                "description": "logs out user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "logs out user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    }
                }
            }
        },
        "/users/profile": {
            "get": {
                "description": "returns user profile",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "returns user profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_heroticket_internal_user.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    }
                }
            }
        },
        "/users/purchased-tickets": {
            "get": {
                "responses": {}
            }
        },
        "/users/refresh-token": {
            "post": {
                "description": "refreshes token pair",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "refreshes token pair",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/rest.CommonResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_heroticket_internal_user.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "did": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "tba_address": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "wallet_address": {
                    "type": "string"
                }
            }
        },
        "protocol.AuthorizationRequestMessage": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/protocol.AuthorizationRequestMessageBody"
                },
                "from": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "thid": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                },
                "typ": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "protocol.AuthorizationRequestMessageBody": {
            "type": "object",
            "properties": {
                "callbackUrl": {
                    "type": "string"
                },
                "did_doc": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "message": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "scope": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/protocol.ZeroKnowledgeProofRequest"
                    }
                }
            }
        },
        "protocol.ZeroKnowledgeProofRequest": {
            "type": "object",
            "properties": {
                "circuitId": {
                    "type": "string"
                },
                "id": {
                    "description": "unique request id",
                    "type": "integer"
                },
                "optional": {
                    "type": "boolean"
                },
                "query": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "rest.CommonResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Hero Ticket API",
	Description:      "This is Hero Ticket API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
