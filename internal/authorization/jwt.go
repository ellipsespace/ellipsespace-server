package authorization

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/qwuiemme/ellipsespace-server/config"
)

var secrets = []byte(config.New().JwtSecrets)

type JWT struct {
	SessionBase
	jwt.StandardClaims
}

func GenerateJWT(s SessionBase) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWT{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(720 * time.Hour).Unix(),
		},
		SessionBase: SessionBase{
			Id:          s.Id,
			SessionName: s.SessionName,
			AccessLevel: s.AccessLevel,
		},
	})

	return token.SignedString(secrets)
}

func ParseJWT(rawToken string) (SessionBase, error) {
	token, err := jwt.ParseWithClaims(rawToken, &JWT{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted signing method: %s", t.Header["Alg"])
		}

		return secrets, nil
	})

	if err != nil {
		return SessionBase{}, err
	}

	if claims, ok := token.Claims.(*JWT); ok && token.Valid {
		return claims.SessionBase, nil
	}

	return SessionBase{}, errors.New("invalid JWT token")
}
