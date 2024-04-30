// Package v1 Code generated by swaggo/swag. DO NOT EDIT
package v1

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Курмыза Павел",
            "email": "tmrrwnxtsn@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/payment/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Платежи"
                ],
                "summary": "Создать запрос на пополнение баланса",
                "parameters": [
                    {
                        "description": "Тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.paymentCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "$ref": "#/definitions/v1.paymentCreateResponse"
                        }
                    },
                    "default": {
                        "description": "Ответ с ошибкой",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        },
        "/payment/methods": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Платежи"
                ],
                "summary": "Получить список способов для пополнения баланса",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор клиента",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Валюта платежа в соответствии со стандартом ISO 4217",
                        "name": "currency",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Код языка, обозначение по RFC 5646",
                        "name": "lang_code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "$ref": "#/definitions/v1.paymentMethodsResponse"
                        }
                    },
                    "default": {
                        "description": "Ответ с ошибкой",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        },
        "/payout/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Выплаты"
                ],
                "summary": "Создать запрос на вывод средств",
                "parameters": [
                    {
                        "description": "Тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.payoutCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "$ref": "#/definitions/v1.payoutCreateResponse"
                        }
                    },
                    "default": {
                        "description": "Ответ с ошибкой",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        },
        "/payout/methods": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Выплаты"
                ],
                "summary": "Получить список способов для вывода средств",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор клиента",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Валюта выплаты в соответствии со стандартом ISO 4217",
                        "name": "currency",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Код языка, обозначение по RFC 5646",
                        "name": "lang_code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "$ref": "#/definitions/v1.payoutMethodsResponse"
                        }
                    },
                    "default": {
                        "description": "Ответ с ошибкой",
                        "schema": {
                            "$ref": "#/definitions/v1.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.commission": {
            "type": "object",
            "required": [
                "type"
            ],
            "properties": {
                "absolute": {
                    "description": "Значение комиссии (передается, если тип комиссии, \"type\", равен \"fixed\" или \"combined\")",
                    "type": "number",
                    "example": 10
                },
                "caption": {
                    "description": "Текстовая репрезентация комиссии (передается, если тип комиссии, \"type\", равен \"text\")",
                    "type": "string",
                    "example": "Комиссия взимается провайдером"
                },
                "currency": {
                    "description": "Код валюты комиссии (передается, если тип комиссии, \"type\", равен \"fixed\" или \"combined\")",
                    "type": "string",
                    "example": "RUB"
                },
                "percent": {
                    "description": "Значение комиссии (передается, если тип комиссии, \"type\", равен \"percent\" или \"combined\")",
                    "type": "number",
                    "example": 11.99
                },
                "type": {
                    "description": "Тип комиссии:\n* \"percent\" - комиссия в процентах\n* \"fixed\" - фиксированная комиссия\n* \"combined\" - комбинированная комиссия, из процентов и фиксированной суммы\n* \"text\" - текстовая коммиссия, например, \"Взимается провайдером\"",
                    "type": "string",
                    "example": "combined"
                }
            }
        },
        "v1.errorContent": {
            "type": "object",
            "required": [
                "code",
                "description"
            ],
            "properties": {
                "code": {
                    "description": "Код ошибки",
                    "type": "string",
                    "example": "InvalidRequest"
                },
                "description": {
                    "description": "Описание ошибки для разработки",
                    "type": "string",
                    "example": "user_id param is required"
                },
                "message": {
                    "description": "Сообщение об ошибке для клиента",
                    "type": "string",
                    "example": "Internal server error occurred. Please try again later."
                }
            }
        },
        "v1.errorResponse": {
            "type": "object",
            "required": [
                "error",
                "success"
            ],
            "properties": {
                "error": {
                    "description": "Развернутая информация об ошибке",
                    "allOf": [
                        {
                            "$ref": "#/definitions/v1.errorContent"
                        }
                    ]
                },
                "success": {
                    "description": "Результат обработки запроса (всегда false)",
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "v1.limits": {
            "type": "object",
            "required": [
                "currency",
                "max_amount",
                "min_amount"
            ],
            "properties": {
                "currency": {
                    "description": "Код валюты лимита в соответствии со стандартом ISO 4217",
                    "type": "string",
                    "example": "RUB"
                },
                "max_amount": {
                    "description": "Максимальное значение суммы (в дробных единицах)",
                    "type": "number",
                    "example": 60000
                },
                "min_amount": {
                    "description": "Минимальное значение суммы (в дробных единицах)",
                    "type": "number",
                    "example": 100
                }
            }
        },
        "v1.method": {
            "type": "object",
            "required": [
                "commission",
                "external_method",
                "external_system",
                "favorite",
                "id",
                "limits",
                "name"
            ],
            "properties": {
                "commission": {
                    "description": "Объект, содержащий данные о комиссии",
                    "allOf": [
                        {
                            "$ref": "#/definitions/v1.commission"
                        }
                    ]
                },
                "external_method": {
                    "description": "Внутренний код платежного метода платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa_bank_card"
                },
                "external_system": {
                    "description": "Внутренний код платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa"
                },
                "favorite": {
                    "description": "Флаг о том, что платежная система добавлена в избранное",
                    "type": "boolean",
                    "example": true
                },
                "id": {
                    "description": "Идентификатор ПС из внутреннего справочника",
                    "type": "string",
                    "example": "CARD"
                },
                "limits": {
                    "description": "Массив объектов, содержащих данные о лимитах",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.limits"
                    }
                },
                "name": {
                    "description": "Название платежной системы",
                    "type": "string",
                    "example": "Банковская карта"
                },
                "tools": {
                    "description": "Массив объектов, содержащих данные о сохраненных платежных инструментах",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.tool"
                    }
                }
            }
        },
        "v1.paymentCreateRequest": {
            "type": "object",
            "required": [
                "amount",
                "currency",
                "external_method",
                "external_system",
                "lang_code",
                "return_urls",
                "user_id"
            ],
            "properties": {
                "additional_data": {
                    "description": "Дополнительная информация, специфичная для платежной системы, к которой направляется целевой запрос",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "ip": "127.0.0.1",
                        "phone_number": "+71234567890"
                    }
                },
                "amount": {
                    "description": "Сумма платежа в минорных единицах валюты (копейки, центы и т.п.)",
                    "type": "integer",
                    "minimum": 100,
                    "example": 10000
                },
                "currency": {
                    "description": "Валюта платежа в соответствии со стандартом ISO 4217",
                    "type": "string",
                    "example": "RUB"
                },
                "external_method": {
                    "description": "Внутренний код платежного метода платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa_bank_card"
                },
                "external_system": {
                    "description": "Внутренний код платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa"
                },
                "lang_code": {
                    "description": "Код языка, обозначение по RFC 5646",
                    "type": "string",
                    "example": "en"
                },
                "return_urls": {
                    "description": "Объект, содержащий ссылки для возврата пользователя для каждого из возможных результатов проведения платежа",
                    "allOf": [
                        {
                            "$ref": "#/definitions/v1.paymentReturnURLs"
                        }
                    ]
                },
                "tool_id": {
                    "description": "Идентификатор сохраненного платежного средства",
                    "type": "integer",
                    "example": 1
                },
                "user_id": {
                    "description": "Идентификатор клиента",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "v1.paymentCreateResponse": {
            "type": "object",
            "required": [
                "operation_id",
                "success",
                "type"
            ],
            "properties": {
                "message": {
                    "description": "Сообщение, которое необходимо показать клиенту",
                    "type": "string",
                    "example": "Баланс пополнен!"
                },
                "operation_id": {
                    "description": "Идентификатор созданного платежа",
                    "type": "integer",
                    "example": 1
                },
                "redirect_url": {
                    "description": "URL платежной страницы, на которую необходимо перенаправить клиента",
                    "type": "string",
                    "example": "https://securepayments.example.com"
                },
                "success": {
                    "description": "Результат обработки запроса (всегда true)",
                    "type": "boolean",
                    "example": true
                },
                "type": {
                    "description": "Тип ответа:\n* Перенаправление клиента на платежную страницу - \"redirect\"\n* Текстовое сообщение - \"message\"",
                    "type": "string",
                    "example": "redirect"
                }
            }
        },
        "v1.paymentMethodsResponse": {
            "type": "object",
            "required": [
                "payment_methods",
                "success"
            ],
            "properties": {
                "payment_methods": {
                    "description": "Массив платежных методов, доступных для пополнения баланса",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.method"
                    }
                },
                "success": {
                    "description": "Результат обработки запроса (всегда true)",
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "v1.paymentReturnURLs": {
            "type": "object",
            "required": [
                "common"
            ],
            "properties": {
                "common": {
                    "description": "URL для возврата пользователя, используемый когда результат платежа неизвестен или по умолчанию",
                    "type": "string",
                    "example": "https://example.com"
                },
                "fail": {
                    "description": "URL для возврата пользователя, используемый при неуспешном осуществлении платежа",
                    "type": "string",
                    "example": "https://example.com/failed"
                },
                "success": {
                    "description": "URL для возврата пользователя, используемый при успешном осуществлении платежа",
                    "type": "string",
                    "example": "https://example.com/success"
                }
            }
        },
        "v1.payoutCreateRequest": {
            "type": "object",
            "required": [
                "amount",
                "currency",
                "external_method",
                "external_system",
                "lang_code",
                "tool_id",
                "user_id"
            ],
            "properties": {
                "additional_data": {
                    "description": "Дополнительная информация, специфичная для платежной системы, к которой направляется целевой запрос",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "ip": "127.0.0.1",
                        "phone_number": "+71234567890"
                    }
                },
                "amount": {
                    "description": "Сумма выплаты в минорных единицах валюты (копейки, центы и т.п.)",
                    "type": "integer",
                    "minimum": 100,
                    "example": 10000
                },
                "currency": {
                    "description": "Валюта выплаты в соответствии со стандартом ISO 4217",
                    "type": "string",
                    "example": "RUB"
                },
                "external_method": {
                    "description": "Внутренний код платежного метода платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa_bank_card"
                },
                "external_system": {
                    "description": "Внутренний код платежной системы, к которой направляется целевой запрос",
                    "type": "string",
                    "example": "yookassa"
                },
                "lang_code": {
                    "description": "Код языка, обозначение по RFC 5646",
                    "type": "string",
                    "example": "en"
                },
                "tool_id": {
                    "description": "Идентификатор сохраненного платежного средства",
                    "type": "integer",
                    "example": 1
                },
                "user_id": {
                    "description": "Идентификатор клиента",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "v1.payoutCreateResponse": {
            "type": "object",
            "required": [
                "operation_id",
                "success"
            ],
            "properties": {
                "operation_id": {
                    "description": "Идентификатор созданной выплаты",
                    "type": "integer",
                    "example": 1
                },
                "success": {
                    "description": "Результат обработки запроса (всегда true)",
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "v1.payoutMethodsResponse": {
            "type": "object",
            "required": [
                "payout_methods",
                "success"
            ],
            "properties": {
                "payout_methods": {
                    "description": "Массив платежных методов, доступных для вывода средств",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.method"
                    }
                },
                "success": {
                    "description": "Результат обработки запроса (всегда true)",
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "v1.tool": {
            "type": "object",
            "required": [
                "caption",
                "id",
                "type"
            ],
            "properties": {
                "caption": {
                    "description": "Значение платежного инструмента, например:\n* Маскированная банковская карта\n* Номер электронного кошелька\n* Адрес электронной почты\n* и т.д.",
                    "type": "string",
                    "example": "444444******4444"
                },
                "details": {
                    "description": "Дополнительная информация о платежном инструменте",
                    "allOf": [
                        {
                            "$ref": "#/definitions/v1.toolDetails"
                        }
                    ]
                },
                "id": {
                    "description": "Идентификатор платежного инструмента",
                    "type": "integer",
                    "example": 14124
                },
                "type": {
                    "description": "Тип платежного инструмента:\n* Банковская карта - \"card\"\n* Электронный кошелек - \"wallet\"",
                    "type": "string",
                    "example": "card"
                }
            }
        },
        "v1.toolDetails": {
            "type": "object",
            "properties": {
                "bank_name": {
                    "description": "Название банка, выпустившего банковскую карту",
                    "type": "string",
                    "example": "Sberbank"
                },
                "card_holder": {
                    "description": "Владелец банковской карты",
                    "type": "string",
                    "example": "Ivanov Ivan"
                },
                "card_type": {
                    "description": "Тип банковской карты",
                    "type": "string",
                    "example": "Visa"
                },
                "expiry_month": {
                    "description": "Срок действия банковской карты (месяц, MM)",
                    "type": "integer",
                    "example": 10
                },
                "expiry_year": {
                    "description": "Срок действия банковской карты (год, YYYY)",
                    "type": "integer",
                    "example": 2023
                },
                "wallet_number": {
                    "description": "Номер электронного кошелька",
                    "type": "string",
                    "example": "410011758831136"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Секретный ключ",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Платежный шлюз для E-commerce системы",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
