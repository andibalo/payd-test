package pkg

import (
	"github.com/andibalo/payd-test/backend/internal/model"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func GenerateToken(user *model.User) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"role":       user.Role,
	})

	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Println(err)

		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to sign JWT")
	}

	return tokenString, nil
}
