package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"main/config"
	"time"
)

type TokenClaim struct {
	UserId       int64  `json:"user_id"`
	TenantId     int64  `json:"tenant_id"`
	UserType     int8   `json:"user_type"`
	UserRole     int64  `json:"user_role"`
	PersonalCode string `json:"personal_code"`
	AuthTitle    string `json:"auth_title"`
	UserStatus   int8   `json:"user_status"`
	DefOrgId     int64  `json:"def_org_id"`
	DefLang      string `json:"def_lang"`
	IsRoot       bool   `json:"is_root"`
	HasOrg       bool   `json:"has_org"`
	Exp          int64  `json:"exp"`
	Iat          int64  `json:"iat"`
	jwt.RegisteredClaims
}

func GenerateToken(cfg *config.Config, claims TokenClaim) (string, error) {
	signingKey := []byte(cfg.Server.APP_SECRET)
	claims.Iat = time.Now().Unix()
	claims.Exp = time.Now().Local().Add(time.Hour * time.Duration(cfg.Server.JWT_TOKEN_EXPIRE_TIME)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}
