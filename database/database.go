package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/** Establece conexión con la base de datos
 *
 * @param ctx context.Context El contexto de la base de datos
 * @param mongodbUri string La URI de conexión con la base de datos
 * @return *mongo.Client El cliente de base de datos
 * @return error El error al conectarse
 */
func Open(ctx context.Context, mongoDBUri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
