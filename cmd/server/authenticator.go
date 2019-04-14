package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rlongo/ictf-gradings-backend/app"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

// NewAuthMiddleware creates a new authentication middleware
// to inject into our webstack
func NewAuthMiddleware(aud, iss string) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: getKeyValidator(aud, iss),
		SigningMethod:       jwt.SigningMethodRS256,
	})
}

func RoleParser(r *http.Request) app.Role {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]
	return getRole(token)
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func getRole(tokenString string) app.Role {
	role := app.RoleNone

	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)
	claims, _ := token.Claims.(*CustomClaims)
	result := strings.Split(claims.Scope, " ")

	var modTest, modStudents bool

	for i := range result {
		switch result[i] {
		case "modify:students":
			modStudents = true
		case "modify:tests":
			modTest = true
		case "read:tests":
			role |= app.RoleInstructor
		case "admin":
			role |= app.RoleAdmin
		}
	}

	if modStudents && modTest && true {
		role |= app.RoleSupervisor
	}

	return app.RoleNone
}

func getKeyValidator(aud, iss string) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		// Verify 'aud' claim
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("jwt: iinvalid audience")
		}
		// Verify 'iss' claim
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("jwt: invalid issuer")
		}

		cert, err := getPemCert(token, iss)
		if err != nil {
			panic(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	}
}

func getPemCert(token *jwt.Token, iss string) (string, error) {
	cert := ""
	resp, err := http.Get(iss + ".well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
