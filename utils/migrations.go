package utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/maramal/user-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReadTags() {
	user := models.User{
		ID:                primitive.NewObjectID(),
		FirstName:         "John",
		LastName:          "Doe",
		Email:             "john@doe.com",
		Password:          "password",
		Type:              "user",
		Status:            "active",
		ProfileImage:      "",
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	t := reflect.TypeOf(user)

	fmt.Println(t.Name(), t.Kind())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mg")

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
