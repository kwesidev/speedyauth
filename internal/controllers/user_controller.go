package controllers

import (
	"database/sql"
	"net/http"

	"github.com/kwesidev/authserver/internal/services"
	"github.com/kwesidev/authserver/internal/utilities"
)

type UserController struct {
	// Registered Services
	db          *sql.DB
	userService services.UserService
}

// Creates a new User Controller for passing requests
func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		db:          db,
		userService: *services.NewUserService(db),
	}
}

// Index  Welcome user
func (this *UserController) Index(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(map[string]interface{})
	userDetails := this.userService.Get(claims["userId"].(int))
	utilities.JSONResponse(w, userDetails)
}

func (c *AuthController) Update(w http.ResponseWriter, r *http.Request) {
}
