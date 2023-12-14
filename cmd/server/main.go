package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/leobelini-studies/go_expert_api/internal/infra/webserver/handlers"
	"net/http"

	"github.com/leobelini-studies/go_expert_api/configs"
	"github.com/leobelini-studies/go_expert_api/internal/entity"
	"github.com/leobelini-studies/go_expert_api/internal/infra/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/leobelini-studies/go_expert_api/docs"
)

// @title                      Go Expert API Example
// @version                    1.0
// @description                Product API with authentication
// @termsOfService             http://swagger.op/terms/

// @contact.name               Leonardo Siervo Belini
// @contact.url                https://lsbelini.dev
// @contact.email              leobelini96@gmail.com

// @license.name               Full Cycle License
// @license.url                http://www.fullcycle.com.br

// @host                       localhost:8081
// @basePath                   /
// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Products
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.API.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	// Users
	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.API.TokenAuth, config.API.JWTExperesIn)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*",httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/docs/doc.json")))

	println("Starting server on port " + config.API.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.API.Port), r)
}

//func LogRequest(next http.Handler) http.Handler{
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		log.Printf("Request %s %s",r.Method,r.URL.Path)
//		next.ServeHTTP(w,r)
//	})
//}
