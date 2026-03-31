package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractBearerToken(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func ParseUserIDFromRequest(r *http.Request, secret string) (*int64, error) {
	tokenString := ExtractBearerToken(r.Header.Get("Authorization"))
	if tokenString == "" {
		return nil, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if sub, ok := claims["sub"]; ok {
		switch value := sub.(type) {
		case float64:
			id := int64(value)
			return &id, nil
		case string:
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			return &parsed, nil
		}
	}

	return nil, fmt.Errorf("sub claim not found")
}
