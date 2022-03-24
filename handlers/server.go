package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/maramal/user-service/docs"
	"github.com/maramal/user-service/middlewares"
	"github.com/maramal/user-service/services"
	"github.com/maramal/user-service/token"
	"github.com/maramal/user-service/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gindump "github.com/tpkeeper/gin-dump"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Config     utils.Config
	TokenMaker token.IMaker
	Client     *mongo.Client
	Database   *mongo.Database
	Router     *gin.Engine
}

/** Crea un nuevo servidor HTTP y configura el router de la API
 *
 * @param config utils.Config "Configuración de la aplicación"
 * @return *Server "Instancia del servidor HTTP"
 * @return error "Error al crear el servidor HTTP"
 */
func NewServer(config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error al crear el token maker: %s", utils.ErrorResponse(err))
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("Error al conectar con la base de datos: %s", utils.ErrorResponse(err))
	}

	server := &Server{
		Config:     config,
		TokenMaker: tokenMaker,
		Client:     client,
		Database:   client.Database("users-dev"),
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()

	// Middlewares
	router.Use(
		gin.Recovery(),
		middlewares.Logger(),
		gindump.Dump(),
		middlewares.CorsConfig(),
	)

	// Documentación
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userService := services.NewUserService(server.Database)
	authService := services.NewAuthService(server.Database)

	// Rutas API
	apiRouter := router.Group("/api")

	// Usuarios
	userRoutes := apiRouter.Group("/users").Use(middlewares.AuthMiddleware(server.TokenMaker))
	newUserHandler(userRoutes, userService)

	// Autenticación
	newAuthHandler(
		apiRouter,
		userService,
		authService,
		server,
	)

	server.Router = router
}
