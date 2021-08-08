package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Shelex/split-specs/internal/users"
	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//data we save in each token
type Claims struct {
	email string //nolint
	id    string //nolint
	jwt.StandardClaims
}

//GenerateToken generates a jwt token and assign an email to it's claims and return it
func GenerateToken(user users.User) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

//ParseToken parses a jwt token and returns the email it it's claims
func ParseToken(tokenStr string) (users.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return users.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return users.User{
			Email: claims["email"].(string),
			ID:    claims["id"].(string),
		}, nil
	}
	return users.User{}, fmt.Errorf("could not parse claims from jwt token")
}
