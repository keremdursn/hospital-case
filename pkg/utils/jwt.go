package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/keremdursn/hospital-case/internal/config"
)

type Claims struct {
	AuthorityID uint   `json:"authority_id"`
	HospitalID  uint   `json:"hospital_id"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(authorityID, hospitalID uint, role string, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		AuthorityID: authorityID,
		HospitalID:  hospitalID,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}
