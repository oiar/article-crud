package controller

import (
	"errors"

	ginjwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

var (
	errActive          = errors.New("Admin active is wrong")
	errUserIDNotExists = errors.New("Get Admin ID is wrong")
)

func (c *Controller) ExtendJWTMiddleWare(JWTMiddleware *ginjwt.GinJWTMiddleware) func(ctx *gin.Context) (uint32, error) {
	JWTMiddleware.Authenticator = func(ctx *gin.Context) (interface{}, error) {
		return c.Login(ctx)
	}

	JWTMiddleware.PayloadFunc = func(data interface{}) ginjwt.MapClaims {
		return ginjwt.MapClaims{
			"userID": data,
		}
	}

	JWTMiddleware.IdentityHandler = func(claims jwt.MapClaims) interface{} {
		return claims["userID"]
	}

	return func(ctx *gin.Context) (uint32, error) {
		id, ok := ctx.Get("userID")
		if !ok {
			return 0, errUserIDNotExists
		}

		v := id.(float64)
		return uint32(v), nil
	}
}
