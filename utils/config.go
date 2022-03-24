package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config guarda toda la configuración de la aplicación
type Config struct {
	Port                 string        `mapstructure:"APP_PORT"`
	MongoURI             string        `mapstructure:"MONGO_URI"`
	SecretKey            string        `mapstructure:"SESSION_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

/** Lee la configuración del archivo o de las variables de entorno
 *
 * @param path string "Ruta del archivo de configuración"
 * @return config Config "Configuración leída"
 * @return err error "Error de lectura de la configuración"
 */
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
