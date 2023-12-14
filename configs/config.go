package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type db struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

type api struct {
	Port         string `mapstructure:"API_PORT"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExperesIn int    `mapstructure:"JWT_EXPIRES_IN"`
	TokenAuth    *jwtauth.JWTAuth
}

type conf struct {
	DB  db
	API api
}

func LoadConfig(path string) (*conf, error) {
	var cfg conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg.DB); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg.API); err != nil {
		panic(err)
	}

	cfg.API.TokenAuth = jwtauth.New("HS256", []byte(cfg.API.JWTSecret), nil)

	return &cfg, nil
}
