package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DocsHandler 구조체
type DocsHandler struct{}

// NewDocsHandler 생성자
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// GetOpenAPISchema OpenAPI 스키마를 반환
func (h *DocsHandler) GetOpenAPISchema(c *gin.Context) {
	schema := h.generateOpenAPISchema()
	c.JSON(http.StatusOK, schema)
}

// GetSwaggerUI Swagger UI HTML을 반환
func (h *DocsHandler) GetSwaggerUI(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Resort Management API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css">
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/docs/schema",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// generateOpenAPISchema OpenAPI 스키마 생성
func (h *DocsHandler) generateOpenAPISchema() map[string]interface{} {
	return map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Resort Management System API",
			"description": "리조트 예약 관리 시스템 API",
			"version":     "1.0.0",
			"contact": map[string]interface{}{
				"name":  "Resort Management Team",
				"email": "support@resort.com",
			},
		},
		"servers": []map[string]interface{}{
			{
				"url":         "http://localhost:8080/api/v1",
				"description": "Local development server",
			},
			{
				"url":         "https://api.resort.com/api/v1",
				"description": "Production server",
			},
		},
		"security": []map[string]interface{}{
			{"bearerAuth": []string{}},
		},
		"tags": []map[string]interface{}{
			{"name": "Health", "description": "헬스체크 API"},
			{"name": "Authentication", "description": "인증 관련 API"},
			{"name": "Users", "description": "사용자 관리 API"},
			{"name": "Rooms", "description": "객실 관리 API"},
			{"name": "Room Groups", "description": "객실 그룹 관리 API"},
			{"name": "Reservations", "description": "예약 관리 API"},
			{"name": "Payment Methods", "description": "결제 수단 관리 API"},
			{"name": "Development", "description": "개발용 API (프로덕션 환경에서는 사용 불가)"},
		},
		"paths": map[string]interface{}{
			"/actuator/health": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"Health"},
					"summary":     "헬스체크",
					"description": "애플리케이션의 전반적인 상태를 확인합니다",
					"operationId": "getHealth",
					"security":    []map[string]interface{}{},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "정상 상태",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/HealthStatus",
									},
								},
							},
						},
						"503": map[string]interface{}{
							"description": "서비스 이용 불가",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/HealthStatus",
									},
								},
							},
						},
					},
				},
			},
			"/auth/login": map[string]interface{}{
				"post": map[string]interface{}{
					"tags":        []string{"Authentication"},
					"summary":     "로그인",
					"description": "사용자 인증 후 JWT 토큰을 발급합니다",
					"operationId": "login",
					"security":    []map[string]interface{}{},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/LoginRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "로그인 성공",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/AuthResponse",
									},
								},
							},
						},
						"401": map[string]interface{}{
							"description": "인증 실패",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
			"/reservations": map[string]interface{}{
				"get": map[string]interface{}{
					"tags":        []string{"Reservations"},
					"summary":     "예약 목록 조회",
					"description": "예약 목록을 조회합니다",
					"operationId": "listReservations",
					"parameters": []map[string]interface{}{
						{
							"name":        "page",
							"in":          "query",
							"description": "페이지 번호 (0부터 시작)",
							"schema": map[string]interface{}{
								"type":    "integer",
								"minimum": 0,
								"default": 0,
							},
						},
						{
							"name":        "size",
							"in":          "query",
							"description": "페이지 크기",
							"schema": map[string]interface{}{
								"type":    "integer",
								"minimum": 1,
								"maximum": 2000,
								"default": 20,
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "조회 성공",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"data": map[string]interface{}{
												"type": "array",
												"items": map[string]interface{}{
													"$ref": "#/components/schemas/ReservationResponse",
												},
											},
											"pagination": map[string]interface{}{
												"$ref": "#/components/schemas/PaginationMeta",
											},
										},
									},
								},
							},
						},
						"401": map[string]interface{}{
							"description": "인증 필요",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": "#/components/schemas/ErrorResponse",
									},
								},
							},
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"bearerAuth": map[string]interface{}{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
					"description":  "JWT 인증을 위한 Bearer 토큰",
				},
			},
			"schemas": map[string]interface{}{
				"ErrorResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message": map[string]interface{}{
							"type":        "string",
							"description": "에러 메시지",
						},
						"errors": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "상세 에러 목록",
						},
						"fieldErrors": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "필드 유효성 검사 에러",
						},
					},
				},
				"HealthStatus": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"UP", "DOWN"},
							"description": "헬스 상태",
						},
						"components": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"db": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"status": map[string]interface{}{
											"type": "string",
											"enum": []string{"UP", "DOWN"},
										},
									},
								},
								"redis": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"status": map[string]interface{}{
											"type": "string",
											"enum": []string{"UP", "DOWN"},
										},
									},
								},
							},
						},
					},
				},
				"PaginationMeta": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "현재 페이지 번호",
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "페이지 크기",
						},
						"totalElements": map[string]interface{}{
							"type":        "integer",
							"description": "전체 요소 수",
						},
						"totalPages": map[string]interface{}{
							"type":        "integer",
							"description": "전체 페이지 수",
						},
					},
				},
				"LoginRequest": map[string]interface{}{
					"type":     "object",
					"required": []string{"username", "password"},
					"properties": map[string]interface{}{
						"username": map[string]interface{}{
							"type":        "string",
							"description": "사용자명",
						},
						"password": map[string]interface{}{
							"type":        "string",
							"description": "비밀번호",
						},
					},
				},
				"AuthResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"accessToken": map[string]interface{}{
							"type":        "string",
							"description": "액세스 토큰",
						},
						"refreshToken": map[string]interface{}{
							"type":        "string",
							"description": "리프레시 토큰",
						},
						"expiresIn": map[string]interface{}{
							"type":        "integer",
							"description": "액세스 토큰 만료 시간 (초)",
						},
					},
				},
				"ReservationResponse": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "integer",
							"description": "예약 ID",
						},
						"roomId": map[string]interface{}{
							"type":        "integer",
							"description": "객실 ID",
						},
						"guestName": map[string]interface{}{
							"type":        "string",
							"description": "투숙객 이름",
						},
						"guestPhone": map[string]interface{}{
							"type":        "string",
							"description": "투숙객 전화번호",
						},
						"checkInAt": map[string]interface{}{
							"type":        "string",
							"format":      "date",
							"description": "체크인 날짜",
						},
						"checkOutAt": map[string]interface{}{
							"type":        "string",
							"format":      "date",
							"description": "체크아웃 날짜",
						},
						"adultCount": map[string]interface{}{
							"type":        "integer",
							"description": "성인 수",
						},
						"childCount": map[string]interface{}{
							"type":        "integer",
							"description": "어린이 수",
						},
						"paymentAmount": map[string]interface{}{
							"type":        "number",
							"description": "결제 금액",
						},
						"memo": map[string]interface{}{
							"type":        "string",
							"description": "메모",
						},
						"canceledAt": map[string]interface{}{
							"type":        "string",
							"format":      "date-time",
							"description": "취소일시",
						},
						"createdAt": map[string]interface{}{
							"type":        "string",
							"format":      "date-time",
							"description": "생성일시",
						},
						"updatedAt": map[string]interface{}{
							"type":        "string",
							"format":      "date-time",
							"description": "수정일시",
						},
					},
				},
			},
		},
	}
}
