package handlers

import (
	"encoding/json"
	"github.com/go-chi/jwtauth"
	"github.com/leobelini-studies/go_expert_api/internal/dto"
	"github.com/leobelini-studies/go_expert_api/internal/entity"
	"github.com/leobelini-studies/go_expert_api/internal/infra/database"
	"net/http"
	"time"
)

type UserHandler struct {
	UserDb        database.UserInterface
	Jwt           *jwtauth.JWTAuth
	JwtExperiesIn int
}

func NewUserHandler(db database.UserInterface, Jwt *jwtauth.JWTAuth, JwtExperiesIn int) *UserHandler {
	return &UserHandler{
		UserDb:        db,
		Jwt:           Jwt,
		JwtExperiesIn: JwtExperiesIn,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.UserDb.Create(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.UserDb.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	clains := map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExperiesIn)).Unix(),
	}
	_, token, _ := h.Jwt.Encode(clains)

	accessToken :=
		struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: token,
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
