package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maramal/user-service/token"
	"github.com/maramal/user-service/utils"
)

// Este middleware permite el acceso a las rutas donde el usuario es de tipo "superadmin" o "admin"
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("sesion no iniciada")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		}

		payload := _payload.(*token.Payload)
		if payload.UserType != "superadmin" && payload.UserType != "admin" {
			err := errors.New("acceso sólo para administradores")
			ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse(err))
		}

		ctx.Next()
	}
}
