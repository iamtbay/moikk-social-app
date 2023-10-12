package helpers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// CREATE JWT
func CreateJWT(ID, name, email string) (string, error) {
	claims := JWTClaims{
		ID,
		name,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// PARSE JWT
func JWTParse(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return &JWTClaims{}, err
	}
}

//jwt error sned

func GetUserIDFromJWT(c *gin.Context) (string, error) {
	//get token for jwt
	tokenString, err := c.Cookie("accessToken")
	if err != nil {
		return "",err
	}

	//check jwt token
	claims, err := JWTParse(tokenString)
	if err != nil {
		return "",err
	}
	return claims.ID,nil
}
