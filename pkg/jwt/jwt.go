package jwt

import (
	"log"
	"time"

	"github.com/Shelex/split-specs/internal/users"
	"github.com/dgrijalva/jwt-go"
)

// secret key being used to sign tokens
// TODO switch to secret management solution
var (
	SecretKey = []byte("secretsecret")
)

// //data we save in each token
// type Claims struct {
// 	username string //nolint
// 	jwt.StandardClaims
// }

//GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(user users.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = user.Username
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

//ParseToken parses a jwt token and returns the username it it's claims
func ParseToken(tokenStr string) (users.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return users.User{
			Username: claims["username"].(string),
			ID:       claims["id"].(string),
		}, nil
	} else {
		return users.User{}, err
	}
}
