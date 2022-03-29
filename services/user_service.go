package services

import (
	"context"
	"errors"
	"time"

	"github.com/maramal/user-service/models"
	"github.com/maramal/user-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateUserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfileImage string `json:"profile_image"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type UpdateUserRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	ProfileImage string `json:"profile_image"`
	Type         string `json:"type"`
	Status       string `json:"status"`
}

type ChangePasswordRequest struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type GetUsersResponse struct {
	Users []models.User `json:"users"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type GetUserResponse struct {
	User models.User `json:"user"`
}

type UpdateUserResponse struct {
	User models.User `json:"user"`
}

type IUserService interface {
	GetUsers() (response GetUsersResponse, err error)
	CreateUser(req CreateUserRequest) (response CreateUserResponse, err error)

	GetUser(id string) (response GetUserResponse, err error)
	UpdateUser(id string, req UpdateUserRequest) (response UpdateUserResponse, err error)
	DeleteUser(id string) (err error)

	ChangePassword(id string, req ChangePasswordRequest) (err error)
	SetSuperadmin(id string, enable bool) (err error)

	GetUserByEmail(email string) (response GetUserResponse, err error)
}

type UserService struct {
	db *mongo.Database
}

var ctx = context.Background()

/** Obtiene todos los usuarios
 *
 * @return GetUsersResponse "Los usuarios"
 * @return err error "El error de la operación"
 */
func (service *UserService) GetUsers() (response GetUsersResponse, err error) {
	var users []models.User
	collection := service.db.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return
	}

	if err = cursor.All(ctx, &users); err != nil {
		return
	}

	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}

	response.Users = users
	return
}

/** Crea un usuario
 *
 * @param req CreateUserRequest "Los valores del usuario a crear"
 * @return CreateUserResponse "El id del usuario creado"
 * @return err error "El error de la operación"
 */
func (service *UserService) CreateUser(req CreateUserRequest) (response CreateUserResponse, err error) {
	collection := service.db.Collection("users")

	filter := bson.M{"email": req.Email}
	existingUser, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	if existingUser > 0 {
		err = errors.New("el correo electrónico ya está ingresado en la base de datos")
		return
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return
	}

	user := models.User{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		Password:          password,
		Type:              req.Type,
		Status:            req.Status,
		ProfileImage:      req.ProfileImage,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return
	}

	response.UserID = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

/** Obtiene un usuario
 *
 * @param id string "El ID del usuario"
 * @return GetUserResponse "El usuario"
 * @return err error "El error de la operación"
 */
func (service *UserService) GetUser(userId string) (response GetUserResponse, err error) {
	collection := service.db.Collection("users")
	var user models.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}

	filter := bson.M{"_id": id}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("no se encontró el usuario")
		return
	}

	result := collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		return
	}

	err = result.Decode(&user)
	if err != nil {
		return
	}

	user.Password = ""

	response.User = user
	return

}

/** Actualiza un usuario
 *
 * @param req UpdateUserRequest "Los valores del usuario a actualizar"
 * @param id string "El id del usuario"
 * @return UpdateUserResponse "Los datos del usuario actualizado"
 * @return err error "El error de la operación"
 */
func (service *UserService) UpdateUser(userId string, req UpdateUserRequest) (response UpdateUserResponse, err error) {
	collection := service.db.Collection("users")
	var user models.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}

	filter := bson.M{"_id": id}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("no se encontró el usuario")
		return
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return
	}

	if err = cursor.Decode(&user); err != nil {
		return
	}

	user.UpdatedAt = time.Now()

	updatedValues, err := utils.ToDoc(user)
	if err != nil {
		return
	}

	if req.ProfileImage != "" && req.ProfileImage != user.ProfileImage {
		// guardar imagen, etc
	}

	update := bson.D{{"$set", updatedValues}}
	if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
		return
	}

	user.Password = ""

	response.User = user
	return
}

/** Elimina un usuario
 *
 * @param id string "El id del usuario"
 * @return err error "El error de la operación"
 */
func (service *UserService) DeleteUser(userId string) (err error) {
	collection := service.db.Collection("users")

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}

	filter := bson.M{"id": id}
	result := collection.FindOneAndDelete(ctx, filter)
	if err = result.Err(); err != nil {
		return
	}

	return
}

/** Cambia la contraseña de un usuario
 *
 * @param id string "El id del usuario"
 * @param req ChangePasswordRequest "Los valores de la contraseña"
 * @return err error "El error de la operación"
 */
func (service *UserService) ChangePassword(userId string, req ChangePasswordRequest) (err error) {
	collection := service.db.Collection("users")
	var user models.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}

	if req.Password == "" {
		err = errors.New("la contraseña no puede estar vacía")
		return
	}
	if req.Password != req.PasswordConfirmation {
		err = errors.New("las contraseñas no coinciden")
		return
	}

	filter := bson.M{"_id": id}

	result := collection.FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		return
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return
	}

	user.Password = password
	user.PasswordChangedAt = time.Now()

	updatedValues, err := utils.ToDoc(user)
	if err != nil {
		return
	}

	update := bson.D{{"$set", updatedValues}}
	if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
		return
	}

	return
}

/** Configura un usuario como super usuario
 *
 * @param id string "El id del usuario"
 * @param enable bool "Si se desea habilitar o deshabilitar"
 * @return err error "El error de la operación"
 */
func (service *UserService) SetSuperadmin(userId string, enable bool) (err error) {
	collection := service.db.Collection("users")
	var user models.User

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}

	filter := bson.M{"_id": id}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return
	}

	if enable {
		user.Type = "superadmin"
	} else {
		user.Type = "user"
	}

	user.UpdatedAt = time.Now()

	updatedValues, err := utils.ToDoc(user)
	if err != nil {
		return
	}
	update := bson.D{{"$set", updatedValues}}
	_, err = collection.UpdateOne(ctx, filter, update)

	return err
}

/** Obtiene un usuario por su email
 *
 * @param email string "El email del usuario"
 * @return GetUserResponse "El usuario"
 * @return err error "El error de la operación"
 */
func (service *UserService) GetUserByEmail(email string) (response GetUserResponse, err error) {
	collection := service.db.Collection("users")
	var user models.User

	filter := bson.M{"email": email}
	if err = collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return
	}

	response.User = user
	return
}

func NewUserService(db *mongo.Database) IUserService {
	return &UserService{db: db}
}
