package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"time"
)

type JwtClaims struct {
	UserID       string `json:"userId"`
	RoleID       string `json:"roleId"`
	LoginVersion string `json:"loginVersion"`
	Device       string `json:"device"`
	jwt.RegisteredClaims
}

func GetJwtExpireTime() time.Duration {
	return time.Duration(config.C.Jwt.ExpireHours) * time.Hour
}

// GenerateJwtToken 生成 token
func GenerateJwtToken(claims JwtClaims) (string, error) {

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(GetJwtExpireTime()))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.Jwt.Secret))
}

// ParseJwtToken 解析 token
func ParseJwtToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.C.Jwt.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
