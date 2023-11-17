package jwt

import (
	"errors"
	"go/goProyect/models"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var Email string
var IdUser string

func ProcesoToken(token string, JWTSign string) (*models.Claim, bool, string, error) {
	miClave := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return &claims, false, string(""), errors.New("formato de token invalido")
	}

	//obtenemos el token sin espacios vacios y sin la palabra "Bearer"
	token = strings.TrimSpace(splitToken[1])

	//decodificar miClave -> (JWTSIGN) y procesa el token con la clave decodificada para ver si es valido
	tokenProccessed, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return miClave, nil //return de funcion anonima -> el codigo continua
	})
	if err == nil {
		//Rutina que checkea contra la DB
	}

	if !tokenProccessed.Valid {
		return &claims, false, string(""), errors.New("Token Invalido")
	}

	return &claims, false, string(""), nil
}
