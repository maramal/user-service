package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// ErrorResponse devuelve una respuesta JSON con un mensaje de error
func ErrorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

// SuccessResponse devuelve una respuesta JSON con un mensaje de éxito
func SuccessResponse(data interface{}) gin.H {
	return gin.H{
		"data": data,
	}
}

// Convierte un modelo en una documento BSON
func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
