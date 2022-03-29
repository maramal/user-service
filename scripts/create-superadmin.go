package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/maramal/user-service/database"
	"github.com/maramal/user-service/models"
	"github.com/maramal/user-service/utils"
)

var ctx = context.Background()

func main() {
	adminFirstName := flag.String("firstname", "Super", "Nombre de administrador")
	adminLastName := flag.String("lastname", "Admin", "Apellido del administrador")
	adminEmail := flag.String("email", "admin@mafer.dev", "El correo electrónico del administrador")
	adminPassword := flag.String("password", "123456", "Contraseña del administrador")

	// Obtiene los valores de los argumentos pasados por CLI
	flag.Parse()

	// Carga la configuración de app.env
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error al leer la configuración: %v", err)
	}

	// Encripta el password
	password, err := utils.HashPassword(*adminPassword)
	if err != nil {
		log.Fatalf("Error al encriptar la contraseña del administrador: %v", err)
	}

	// Crea el cliente MongoDB
	client, err := database.Open(ctx, config.MongoURI)
	if err != nil {
		log.Fatalf("Error al crear el cliente de la base de datos: %v", err)
	}

	// Asigna la colección de usuarios
	collection := client.Database("users-dev").Collection("users")

	// Crea el usuario
	admin := models.User{
		FirstName:         *adminFirstName,
		LastName:          *adminLastName,
		Email:             *adminEmail,
		Password:          password,
		Type:              "superadmin",
		Status:            "active",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		PasswordChangedAt: time.Now(),
	}

	// Inserta el registro en la base de dats
	_, err = collection.InsertOne(ctx, admin)
	if err != nil {
		log.Fatalf("Error al crear el superadministrador: %v", err)
	}

	log.Println("El usuario ha sido creado.")
}
