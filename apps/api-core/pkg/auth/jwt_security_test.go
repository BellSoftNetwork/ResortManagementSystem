package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_SecurityTests(t *testing.T) {
	secretKey := "test-secret-key"
	jwtService := NewJWTService(secretKey, 15*time.Minute, 7*24*time.Hour)

	t.Run("다른 시크릿으로 서명된 토큰은 검증에 실패한다", func(t *testing.T) {
		// Given: 다른 시크릿으로 서명된 토큰 생성
		differentSecret := "different-secret-key"
		differentService := NewJWTService(differentSecret, 15*time.Minute, 7*24*time.Hour)

		invalidToken, err := differentService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		// When: 원래 서비스로 토큰 검증
		claims, err := jwtService.ValidateToken(invalidToken)

		// Then: 검증 실패
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("만료된 토큰은 검증에 실패한다", func(t *testing.T) {
		// Given: 만료된 토큰 생성
		expiredService := NewJWTService(secretKey, -1*time.Hour, 7*24*time.Hour) // 1시간 전 만료
		expiredToken, err := expiredService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		// When: 토큰 검증
		claims, err := jwtService.ValidateToken(expiredToken)

		// Then: 검증 실패
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("잘못된 형식의 토큰은 검증에 실패한다", func(t *testing.T) {
		testCases := []struct {
			name  string
			token string
		}{
			{
				name:  "완전히 잘못된 형식",
				token: "invalid.jwt.token",
			},
			{
				name:  "부분적으로 누락된 형식",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0", // 서명 누락
			},
			{
				name:  "서명 없는 토큰",
				token: "eyJhbGciOiJub25lIn0.eyJzdWIiOiIxMjM0NTY3ODkwIn0",
			},
			{
				name:  "빈 문자열",
				token: "",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// When: 잘못된 형식의 토큰 검증
				claims, err := jwtService.ValidateToken(tc.token)

				// Then: 검증 실패
				assert.Error(t, err)
				assert.Nil(t, claims)
			})
		}
	})

	t.Run("변조된 페이로드를 가진 토큰은 검증에 실패한다", func(t *testing.T) {
		// Given: 유효한 토큰을 생성하고 수동으로 페이로드를 변조
		validToken, err := jwtService.GenerateAccessToken(1, "testuser", "USER")
		assert.NoError(t, err)

		// 토큰을 분리하여 페이로드 부분을 변조
		parts := strings.Split(validToken, ".")
		assert.Len(t, parts, 3)

		// 변조된 페이로드 (권한을 ADMIN으로 변경)
		tamperedPayload := "eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiYXV0aG9yaXRpZXMiOiJBRE1JTiIsInN1YiI6IjEiLCJleHAiOjE3NTAwMDAwMDAsImlhdCI6MTcwMDAwMDAwMH0"
		tamperedToken := parts[0] + "." + tamperedPayload + "." + parts[2]

		// When: 변조된 토큰 검증
		claims, err := jwtService.ValidateToken(tamperedToken)

		// Then: 변조된 토큰은 검증 실패
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("리프레시 토큰의 서명 검증", func(t *testing.T) {
		// Given: 다른 시크릿으로 서명된 리프레시 토큰
		differentSecret := "different-refresh-secret"
		differentService := NewJWTService(differentSecret, 15*time.Minute, 7*24*time.Hour)

		invalidRefreshToken, err := differentService.GenerateRefreshToken(1, "device123")
		assert.NoError(t, err)

		// When: 원래 서비스로 리프레시 토큰 검증
		claims, err := jwtService.ValidateRefreshToken(invalidRefreshToken)

		// Then: 검증 실패
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "signature is invalid")
	})

	t.Run("지원하지 않는 서명 알고리즘은 거부된다", func(t *testing.T) {
		// Given: RS256 알고리즘으로 생성된 토큰 (HMAC 대신 RSA 사용)
		// 유효한 RS256 헤더와 페이로드를 만들고 가짜 서명을 추가
		rsaHeader := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9" // {"alg":"RS256","typ":"JWT"}
		validPayload := "eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiYXV0aG9yaXRpZXMiOiJVU0VSIiwic3ViIjoiMSIsImV4cCI6MTc1MDAwMDAwMCwiaWF0IjoxNzAwMDAwMDAwfQ"
		fakeSignature := "dGVzdC1zaWduYXR1cmU" // base64("test-signature")

		unsupportedAlgToken := rsaHeader + "." + validPayload + "." + fakeSignature

		// When: 토큰 검증
		claims, err := jwtService.ValidateToken(unsupportedAlgToken)

		// Then: 검증 실패
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("빈 클레임이나 필수 필드 누락 토큰 검증", func(t *testing.T) {
		// Given: subject가 없는 토큰 생성
		claims := &Claims{
			Username:    "testuser",
			Authorities: "USER",
			RegisteredClaims: jwt.RegisteredClaims{
				// Subject 누락
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		assert.NoError(t, err)

		// When: 토큰 검증
		validatedClaims, err := jwtService.ValidateToken(tokenString)

		// Then: 유효한 토큰이지만 Subject가 비어있음을 확인
		assert.NoError(t, err)
		assert.NotNil(t, validatedClaims)
		assert.Equal(t, "", validatedClaims.Subject)
		assert.Equal(t, "testuser", validatedClaims.Username)
	})
}
