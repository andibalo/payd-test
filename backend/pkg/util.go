package pkg

import (
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ToPointer[K any](val K) *K {
	return &val
}

func GetCursorData(cursor string) (string, string) {

	splitCursor := strings.Split(cursor, "_")

	return splitCursor[0], splitCursor[1]
}

func NullStrToStr(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func GetIntParam(c *gin.Context, key string) (int64, error) {
	if c.Param(key) == "" {
		return 0, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("Param not found")
	}

	val, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, oops.Code(response.BadRequest.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusBadRequest).Errorf("%s should be integer, got error: %v", key, err)
	}

	return int64(val), nil
}
