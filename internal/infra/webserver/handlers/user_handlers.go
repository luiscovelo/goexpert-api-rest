package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/luiscovelo/goexpert-api-rest/internal/dto"
	"github.com/luiscovelo/goexpert-api-rest/internal/entity"
	"github.com/luiscovelo/goexpert-api-rest/internal/infra/database"
)

type UserHandler struct {
	Response
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserhandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	if jwt == nil {
		panic("param jwtauth.JWTAuth cannot be nil")
	}

	if jwtExpiresIn <= 0 {
		panic("param jwtExpiresIn cannot be zero")
	}

	return &UserHandler{
		UserDB:       db,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

// LoginUser godoc
//	@Summary		Get User Token
//	@Description	Get User Token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GetJwtInput	true	"user credentials"
//	@Success		200		{object}	dto.GetJwtOutput
//	@Failure		401		{object}	handlers.Response
//	@Failure		500		{object}	handlers.Response
//	@Router			/users/login [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, req *http.Request) {
	var user dto.GetJwtInput
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		h.JSON(w, http.StatusNotFound, err)
		return
	}

	if !u.ValidatePassword(user.Password) {
		h.JSON(w, http.StatusUnauthorized, "email or password is wrong")
		return
	}

	_, token, err := h.Jwt.Encode(
		map[string]interface{}{
			"sub": u.ID.String(),
			"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
		},
	)

	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}

	accessToken := dto.GetJwtOutput{
		AccessToken: token,
	}

	h.JSON(w, http.StatusOK, accessToken)
}

// CreateUser godoc
//	@Summary		Create User
//	@Description	Create User
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.CreateUserInput	true	"user request"
//	@Success		201
//	@Failure		500	{object}	handlers.Response
//	@Router			/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		h.JSON(w, http.StatusBadRequest, err)
		return
	}
	if err := h.UserDB.Create(u); err != nil {
		h.JSON(w, http.StatusInternalServerError, err)
		return
	}

	newUser, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		h.JSON(w, http.StatusNotFound, err)
		return
	}

	h.JSON(w, http.StatusCreated, newUser)
}
