package database

import (
	"errors"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

/** Abre una conexión con la base de datos
 *
 * @param settings db.ConnectionURL "La configuración de la conexión"
 * @return db.Session "La sesión con la base de datos"
 * @return error "El error de conexión"
 */
func Open(settings db.ConnectionURL) (db.Session, error) {
	db, err := mysql.Open(settings)
	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		return nil, errors.New("error al conectar la base de datos")
	}

	return db.WithContext(db.Context()), nil
}
