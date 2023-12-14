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

// CreateUser Create user godoc
// @Summary     Create user
// @Description Create user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       resquest body dto.CreateUserInput true "user request"
// @Success     201
// @Failure     500 {object} dto.ErrorOutput
// @Router      /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	if err := h.UserDb.Create(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetJWT Get JWT godoc
// @Summary     Get a user JWT
// @Description Get a user JWT
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       resquest body dto.GetJWTInput true "user credentials"
// @Success     200  {object} dto.GetJWTOutput
// @Failure     404 {object} dto.ErrorOutput
// @Failure     500 {object} dto.ErrorOutput
// @Router      /users/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var user dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	u, err := h.UserDb.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := dto.ErrorOutput{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	clains := map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExperiesIn)).Unix(),
	}
	_, token, _ := h.Jwt.Encode(clains)

	accessToken := dto.GetJWTOutput{
		AccessToken: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
