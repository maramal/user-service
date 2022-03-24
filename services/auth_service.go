package services

import (
	"time"

	"github.com/maramal/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateSessionParams struct {
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type IAuthService interface {
	CreateSession(params CreateSessionParams) (models.Session, error)
	GetSession(sessionId string) (models.Session, error)
}

type AuthService struct {
	db *mongo.Database
}

/** Crea una sesión para un usuario
 *
 * @param params CreateSessionParams "Los parámetros de la sesión a crear"
 * @return models.Session "La sesión creada"
 * @return error "El error que ocurrió al crear la sesión"
 */
func (service *AuthService) CreateSession(params CreateSessionParams) (models.Session, error) {
	var collection = service.db.Collection("sessions")

	session := models.Session{
		Email:        params.Email,
		RefreshToken: params.RefreshToken,
		UserAgent:    params.UserAgent,
		ClientIP:     params.ClientIp,
		IsBlocked:    params.IsBlocked,
		CreatedAt:    time.Now(),
		ExpiresAt:    params.ExpiresAt,
	}

	result, err := collection.InsertOne(ctx, &session)
	if err != nil {
		return models.Session{}, err
	}

	session.ID = result.InsertedID.(primitive.ObjectID)

	return session, nil
}

/** Obtiene una sesión
 *
 * @param sessionId string "El id de la sesión a obtener"
 * @return models.Session "La sesión obtenida"
 * @return error "El error que ocurrió al obtener la sesión"
 */
func (service *AuthService) GetSession(sessionId string) (models.Session, error) {
	var collection = service.db.Collection("sessions")
	var session models.Session

	var filter = bson.M{"id": sessionId}
	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func NewAuthService(db *mongo.Database) IAuthService {
	return &AuthService{db: db}
}
