package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maramal/user-service/database"
	_ "github.com/maramal/user-service/docs"
	"github.com/maramal/user-service/middlewares"
	"github.com/maramal/user-service/services"
	"github.com/maramal/user-service/token"
	"github.com/maramal/user-service/utils"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gindump "github.com/tpkeeper/gin-dump"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Config     utils.Config
	TokenMaker token.IMaker
	Client     *mongo.Client
	Database   *mongo.Database
	Router     *gin.Engine
	APMApp     *newrelic.Application
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

	client, err := database.Open(context.TODO(), config.MongoURI)
	if err != nil {
		return nil, fmt.Errorf("Error al conectar con la base de datos: %s", utils.ErrorResponse(err))
	}

	server := &Server{
		Config:     config,
		TokenMaker: tokenMaker,
		Client:     client,
		Database:   client.Database("users-dev"),
	}

	if config.APMAppName != "" && config.APMLicense != "" {
		app, err := configAPM(config)
		if err != nil {
			return nil, fmt.Errorf("Error al conectar la aplicación de monitorización: %s", utils.ErrorResponse(err))
		}

		server.APMApp = app
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()

	// Middlewares
	if server.APMApp != nil {
		router.Use(nrgin.Middleware(server.APMApp))
	}

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
	adminRouter := apiRouter.Group("/admin")
	authRouter := apiRouter.Group("/")

	adminRouter.Use(middlewares.AuthMiddleware(server.TokenMaker)).Use(middlewares.AdminMiddleware())
	authRouter.Use(middlewares.AuthMiddleware(server.TokenMaker))

	// Usuarios
	userRoutes := adminRouter.Group("/users")
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

func configAPM(config utils.Config) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(config.APMAppName),
		newrelic.ConfigLicense(config.APMLicense),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
}
