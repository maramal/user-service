package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maramal/user-service/services"
	"github.com/maramal/user-service/utils"
)

// @Summary	Obtiene todos los usuarios
// @ID 		get-users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} services.GetUsersResponse
// @Failure 400 {object} services.GetUsersResponse
// @Router 	/admin/users [get]
func handleGetUsers(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := service.GetUsers()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(users))
	}
}

// @Summary Crea un usuario
// @ID 		create-user
// @Accept 	json
// @Produce json
// @Security ApiKeyAuth
// @Param   CreateUserRequest body services.CreateUserRequest true "Datos del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users [post]
func handleCreateUser(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req services.CreateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}

		userID, err := service.CreateUser(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(userID))
	}
}

// @Summary Obtiene un usuario
// @ID 		get-user
// @Produce json
// @Security ApiKeyAuth
// @Param 	id path string true "ID del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id} [get]
func handleGetUser(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		user, err := service.GetUser(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(user))
	}
}

// @Summary Actualiza un usuario
// @ID 		update-user
// @Accept 	json
// @Produce json
// @Security ApiKeyAuth
// @Param   id 					path string						true "ID del usuario"
// @Param 	UpdateUserRequest 	body services.UpdateUserRequest true "Datos del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id} [put]
func handleUpdateUser(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req services.UpdateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}

		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		user, err := service.UpdateUser(id, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(user))
	}
}

// @Summary Elimina un usuario
// @ID 		delete-user
// @Produce json
// @Security ApiKeyAuth
// @Param 	id path string true "ID del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id} [delete]
func handleDeleteUser(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		err := service.DeleteUser(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(nil))
	}
}

// @Summary Cambia la contraseña de un usuario
// @ID 		change-password
// @Accept 	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	id						path string							true "ID del usuario"
// @Param 	ChangePasswordRequest	body services.ChangePasswordRequest true "Datos del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id}/password [put]
func handleChangePassword(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req services.ChangePasswordRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}

		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		err := service.ChangePassword(id, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(nil))
	}
}

// @Summary Configura un usuario como super administrador
// @ID 		set-super-admin
// @Produce json
// @Security ApiKeyAuth
// @Param 	id path string true "ID del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id}/set-super-admin [post]
func handleSetSuperadmin(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		err := service.SetSuperadmin(id, true)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(nil))
	}
}

// @Summary Configura un super administrador como usuario
// @ID 		unset-super-admin
// @Produce json
// @Security ApiKeyAuth
// @Param 	id path string true "ID del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/{id}/unset-super-admin [post]
func handleUnsetSuperadmin(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el id es requerido")))
			return
		}

		err := service.SetSuperadmin(id, false)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(nil))
	}
}

// @Summary Obtiene un usuario por su correo electrónico
// @ID 		get-user-by-email
// @Produce json
// @Security ApiKeyAuth
// @Param 	email path string true "Correo electrónico del usuario"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router 	/admin/users/email/{email} [get]
func handleGetUserByEmail(service services.IUserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Param("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(errors.New("el email es requerido")))
			return
		}

		user, err := service.GetUserByEmail(email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, utils.SuccessResponse(user))
	}
}

/** Crea un nuevo grupo de endpoints
 *
 * @param group *gin.RouterGroup "El grupo de endpoints padre"
 * @param service services.IUserService "El servicio de usuarios"
 * @return *gin.RouterGroup "El grupo de endpoints creado"
 */
func newUserHandler(group gin.IRoutes, userService services.IUserService) *gin.IRoutes {
	group.GET("/", handleGetUsers(userService))
	group.POST("/", handleCreateUser(userService))

	group.GET("/:id", handleGetUser(userService))
	group.PUT("/:id", handleUpdateUser(userService))
	group.DELETE("/:id", handleDeleteUser(userService))

	group.POST("/:id/password", handleChangePassword(userService))
	group.POST("/:id/set-superadmin", handleSetSuperadmin(userService))
	group.POST("/:id/unset-superadmin", handleUnsetSuperadmin(userService))

	group.GET("/email/:email", handleGetUserByEmail(userService))

	return &group
}
