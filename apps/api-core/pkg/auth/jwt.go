package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type Claims struct {
	Username    string `json:"username"`
	Authorities string `json:"authorities"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	DeviceFingerprint string `json:"deviceFingerprint,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:          secretKey,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (j *JWTService) GenerateTokenPair(userID uint, username, role string) (string, string, error) {
	accessToken, err := j.GenerateAccessToken(userID, username, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateRefreshToken(userID, "")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) GenerateTokenPairWithDevice(userID uint, username, role, deviceFingerprint string) (string, string, error) {
	accessToken, err := j.GenerateAccessToken(userID, username, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateRefreshToken(userID, deviceFingerprint)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) GenerateAccessToken(userID uint, username, role string) (string, error) {
	claims := &Claims{
		Username:    username,
		Authorities: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) GenerateRefreshToken(userID uint, deviceFingerprint string) (string, error) {
	claims := &RefreshClaims{
		DeviceFingerprint: deviceFingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return j.ValidateToken(tokenString)
}

func (j *JWTService) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

func (j *JWTService) GetAccessTokenExpiry() time.Duration {
	return j.accessTokenExpiry
}

func (j *JWTService) GetRefreshTokenExpiry() time.Duration {
	return j.refreshTokenExpiry
}

func (j *JWTService) GetAccessTokenExpiryMillis() int64 {
	return time.Now().Add(j.accessTokenExpiry).UnixMilli()
}

func (j *JWTService) GetUserIDFromRefreshClaims(claims *RefreshClaims) (uint, error) {
	if claims.Subject == "" {
		return 0, errors.New("subject not found in refresh token")
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in refresh token: %v", err)
	}

	return uint(userID), nil
}
