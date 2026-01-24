---
id: api-core-middleware
title: "api-core 미들웨어"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-core 미들웨어

> HTTP 요청 처리 미들웨어 (인증, 에러, 감사, 개발)

---

## 1. 개요

### 1.1 미들웨어 목록

| 미들웨어 | 파일 | 역할 |
|----------|------|------|
| ErrorHandler | `error.go` | 에러 처리, GORM 에러 변환 |
| AuthMiddleware | `auth.go` | JWT 토큰 검증 |
| RoleMiddleware | `auth.go` | 역할 기반 접근 제어 |
| AuditMiddleware | `audit.go` | 감사 컨텍스트 설정 |
| DevelopmentOnlyMiddleware | `development.go` | 개발 환경 전용 |

### 1.2 미들웨어 체인

```go
// 공개 엔드포인트
router.GET("/actuator/health", ...)

// 인증 필요
authGroup := router.Group("/api/v1")
authGroup.Use(middleware.AuthMiddleware(jwtService))
authGroup.Use(middleware.AuditMiddleware())

// 관리자 전용
adminGroup := authGroup.Group("/")
adminGroup.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
```

---

## 2. ErrorHandler

### 2.1 역할

- Gin 에러 수집 및 처리
- GORM 에러를 HTTP 상태 코드로 변환
- 검증 에러 한글 메시지 변환

### 2.2 구현

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            // GORM 에러 처리
            if errors.Is(err, gorm.ErrRecordNotFound) {
                response.NotFound(c, "리소스를 찾을 수 없습니다")
                return
            }
            if errors.Is(err, gorm.ErrDuplicatedKey) {
                response.Conflict(c, "이미 존재하는 데이터입니다")
                return
            }
            
            // 검증 에러 처리
            var validationErrs validator.ValidationErrors
            if errors.As(err, &validationErrs) {
                messages := translateValidationErrors(validationErrs)
                response.BadRequest(c, "입력값 검증 실패", messages)
                return
            }
            
            response.InternalError(c, "서버 오류가 발생했습니다")
        }
    }
}
```

### 2.3 검증 에러 한글화

| 태그 | 한글 메시지 |
|------|-------------|
| required | 필수 입력 항목입니다 |
| min | {n}자 이상 입력해주세요 |
| max | {n}자 이하로 입력해주세요 |
| email | 올바른 이메일 형식이 아닙니다 |
| oneof | 허용된 값 중 하나를 입력해주세요 |

---

## 3. AuthMiddleware

### 3.1 역할

- Authorization 헤더에서 Bearer 토큰 추출
- JWT 토큰 검증
- 사용자 정보 Context에 설정

### 3.2 구현

```go
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Authorization 헤더 확인
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "인증 토큰이 필요합니다")
            c.Abort()
            return
        }
        
        // 2. Bearer 토큰 추출
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.Unauthorized(c, "잘못된 인증 헤더 형식입니다")
            c.Abort()
            return
        }
        
        // 3. JWT 검증
        claims, err := jwtService.ValidateToken(parts[1])
        if err != nil {
            response.Unauthorized(c, "유효하지 않은 토큰입니다")
            c.Abort()
            return
        }
        
        // 4. Context에 사용자 정보 설정
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("userRole", claims.Role)
        
        c.Next()
    }
}
```

### 3.3 Context 키

| 키 | 타입 | 설명 |
|----|------|------|
| userID | uint | 사용자 ID |
| username | string | 사용자명 |
| userRole | string | 역할 (NORMAL, ADMIN, SUPER_ADMIN) |

---

## 4. RoleMiddleware

### 4.1 역할

- Context에서 사용자 역할 확인
- 허용된 역할이 아니면 403 응답

### 4.2 구현

```go
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")
        if !exists {
            response.Forbidden(c, "권한 정보가 없습니다")
            c.Abort()
            return
        }
        
        role := userRole.(string)
        for _, allowedRole := range allowedRoles {
            if role == allowedRole {
                c.Next()
                return
            }
        }
        
        response.Forbidden(c, "접근 권한이 없습니다")
        c.Abort()
    }
}
```

### 4.3 사용 예시

```go
// ADMIN 또는 SUPER_ADMIN만 허용
adminGroup.Use(middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))

// SUPER_ADMIN만 허용
superAdminGroup.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
```

---

## 5. AuditMiddleware

### 5.1 역할

- 인증된 사용자 정보를 감사 컨텍스트에 설정
- 엔티티 생성/수정 시 사용자 추적

### 5.2 구현

```go
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("userID")
        username, _ := c.Get("username")
        
        ctx := audit.SetContext(c.Request.Context(), &audit.Info{
            UserID:   userID.(uint),
            Username: username.(string),
        })
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

### 5.3 감사 정보 사용

```go
// 서비스/리포지토리에서 감사 정보 조회
auditInfo := audit.GetFromContext(ctx)
model.CreatedBy = auditInfo.UserID
model.UpdatedBy = auditInfo.UserID
```

---

## 6. DevelopmentOnlyMiddleware

### 6.1 역할

- 개발/스테이징 환경에서만 접근 허용
- 프로덕션 환경에서 403 응답

### 6.2 구현

```go
func DevelopmentOnlyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        profile := viper.GetString("api.profile")
        
        allowedProfiles := []string{"local", "development", "staging"}
        for _, allowed := range allowedProfiles {
            if profile == allowed {
                c.Next()
                return
            }
        }
        
        response.Forbidden(c, "개발 환경에서만 사용 가능합니다")
        c.Abort()
    }
}
```

### 6.3 사용처

```go
// 개발 도구 엔드포인트
devGroup := router.Group("/api/v1/dev")
devGroup.Use(middleware.AuthMiddleware(jwtService))
devGroup.Use(middleware.RoleMiddleware("SUPER_ADMIN"))
devGroup.Use(middleware.DevelopmentOnlyMiddleware())
```

---

## 7. Helper 함수

### 7.1 Context에서 사용자 정보 조회

```go
// middleware/auth.go

func GetUserID(c *gin.Context) (uint, bool) {
    id, exists := c.Get("userID")
    if !exists {
        return 0, false
    }
    return id.(uint), true
}

func GetUsername(c *gin.Context) (string, bool) {
    username, exists := c.Get("username")
    if !exists {
        return "", false
    }
    return username.(string), true
}

func GetUserRole(c *gin.Context) (string, bool) {
    role, exists := c.Get("userRole")
    if !exists {
        return "", false
    }
    return role.(string), true
}
```

---

## 8. 테스트

- `auth_test.go`
- `error_test.go`

### 테스트 케이스

- 토큰 없이 접근 시 401
- 잘못된 토큰으로 접근 시 401
- 만료된 토큰으로 접근 시 401
- 권한 없는 역할로 접근 시 403
- GORM NotFound 에러 → 404
- GORM DuplicatedKey 에러 → 409
- 검증 에러 한글 메시지
