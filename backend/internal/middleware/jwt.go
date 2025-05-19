package middleware

import (
	"fmt"
	"github.com/andibalo/payd-test/backend/internal/config"
	"github.com/andibalo/payd-test/backend/internal/constants"
	"github.com/andibalo/payd-test/backend/internal/response"
	"github.com/andibalo/payd-test/backend/pkg/httpclient"
	"github.com/andibalo/payd-test/backend/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// TokenClaims : struct for validate token claims
type TokenClaims struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Token     string `json:"token"`
	jwt.RegisteredClaims
}

// contextClaimKey key value store/get token on context
const ContextClaimKey = "ctx.mw.auth.claim"

// JwtMiddleware : check jwt token header bearer scheme
func JwtMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		secretKey := cfg.GetAuthCfg().JWTSecret
		staticToken := cfg.GetAuthCfg().JWTStaticToken

		// token claims
		claims := &TokenClaims{}
		headerToken, err := ParseTokenFromHeader(ctx)
		if err != nil {
			httpresp.HttpRespError(ctx, err)
			return
		}

		if headerToken == staticToken {
			ctx.Set(httpclient.XUserEmail, constants.EMAIL_ADMIN_RMS)
			ctx.Set(ContextClaimKey, &TokenClaims{
				ID:        int64(0),
				Email:     constants.EMAIL_ADMIN_RMS,
				FirstName: "admin",
				LastName:  "rms",
				Role:      constants.ADMIN_ROLE,
				Token:     staticToken,
			})

			ctx.Next()
			return
		}

		token, err := jwt.ParseWithClaims(headerToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // check signing method
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		// check parse token error
		if err != nil {
			cfg.Logger().ErrorWithContext(ctx, "[JWTMiddleware] User unauthorized", zap.Error(err))
			httpresp.HttpRespError(ctx, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(err.Error()))
			return
		}

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			claims.Token = headerToken
			ctx.Set(httpclient.XUserEmail, claims.Email)
			ctx.Set(ContextClaimKey, claims)
			ctx.Next()
		} else {
			cfg.Logger().ErrorWithContext(ctx, "[JWTMiddleware] User unauthorized", zap.Error(err))
			httpresp.HttpRespError(ctx, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(err.Error()))
			return
		}
	}
}

func IsAdminMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		claims := ParseToken(ctx)

		if claims.Email == constants.EMAIL_ADMIN_RMS {
			ctx.Next()
			return
		}

		if claims.Role != constants.ADMIN_ROLE {
			cfg.Logger().ErrorWithContext(ctx, "[IsAdminMiddleware] User unauthorized")
			httpresp.HttpRespError(ctx, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf("User unauthorized"))
			return
		}
	}
}

func ParseTokenFromHeader(ctx *gin.Context) (string, error) {
	var (
		headerToken = ctx.Request.Header.Get("Authorization")
		splitToken  []string
	)

	splitToken = strings.Split(headerToken, "Bearer ")

	// check valid bearer token
	if len(splitToken) <= 1 {
		return "", oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf("Invalid Token")
	}

	return splitToken[1], nil
}

func ParseToken(c *gin.Context) *TokenClaims {

	v := c.Value(ContextClaimKey)
	token := new(TokenClaims)
	if v == nil {
		return token
	}
	out, ok := v.(*TokenClaims)
	if !ok {
		return token
	}

	return out
}

func GetToken(c *gin.Context) string {
	authorization := c.Request.Header.Get("Authorization")
	tokens := strings.Split(authorization, "Bearer ")

	return tokens[1]
}
