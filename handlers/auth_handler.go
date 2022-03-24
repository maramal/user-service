package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maramal/user-service/models"
	"github.com/maramal/user-service/services"
	"github.com/maramal/user-service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userResponse struct {
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	ProfileImage      string    `json:"profile_image"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             primitive.ObjectID `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  userResponse       `json:"user"`
}

func newUserResponse(user models.User) userResponse {
	return userResponse{
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: user.ProfileImage,
		CreatedAt:    user.CreatedAt,
	}
}

// @Summary Ingresa un usuario
// @ID 		login-user
// @Accept 	json
// @Produce	json
// @Param   loginUserRequest body loginUserRequest true 	"Datos del usuario"
// @Success 200 {object} loginUserResponse "Respuesta del login"
// @Failure 400 {object} gin.H	"Error en la solicitud"
// @Router 	/login [post]
func (server *Server) handleLoginUser(userService services.IUserService, AuthService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req loginUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}

		resp, err := userService.GetUserByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		user := resp.User

		err = utils.CheckPassword(req.Password, user.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		accessToken, accessPayload, err := server.TokenMaker.CreateToken(
			user.FirstName,
			user.LastName,
			user.Email,
			user.Type,
			user.ProfileImage,
			server.Config.AccessTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		refreshToken, refreshPayload, err := server.TokenMaker.CreateToken(
			user.FirstName,
			user.LastName,
			user.Email,
			user.Type,
			user.ProfileImage,
			server.Config.RefreshTokenDuration,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		session, err := AuthService.CreateSession(services.CreateSessionParams{
			Email:        user.Email,
			RefreshToken: refreshToken,
			UserAgent:    ctx.Request.UserAgent(),
			ClientIp:     ctx.ClientIP(),
			IsBlocked:    false,
			ExpiresAt:    refreshPayload.ExpiredAt,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		response := loginUserResponse{
			SessionID:             session.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
			User:                  newUserResponse(user),
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func newAuthHandler(group *gin.RouterGroup, userService services.IUserService, authService services.IAuthService, server *Server) *gin.RouterGroup {
	group.POST("/login", server.handleLoginUser(userService, authService))

	return group
}
