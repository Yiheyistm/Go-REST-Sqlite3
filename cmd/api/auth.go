package main

import (
	"fmt"
	"net/http"

	"github.com/Yiheyistm/go-restful-api/internal/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserResponse struct {
	Token string `json:"token"`
}

// Login logs in a user
//
//	@Summary		Logs in a user
//	@Description	Logs in a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	loginUserRequest	true	"User"
//	@Success		200	{object}	loginUserRequest
//	@Router			/api/v1/auth/login [post]
func (app *application) loginUser(c *gin.Context) {
	var login loginUserRequest

	if err := c.ShouldBindBodyWithJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	user, err := app.Model.Users.GetByEmail(login.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid email or password"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		fmt.Printf("Failed to compare password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something went wrong"})
		return
	}

	token, err := app.Model.Users.GenerateToken(user.ID, app.JwtSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate token"})
		return
	}
	fmt.Println(user)
	c.JSON(http.StatusOK, loginUserResponse{Token: token})
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body	registerUserRequest	true	"User"
// @Success		201	{object}	database.User
// @Router			/api/v1/auth/register [post]
func (app *application) registerUser(c *gin.Context) {
	var register registerUserRequest

	if err := c.ShouldBindBodyWithJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something went wrong"})
	}
	registerPassword := string(hasedPassword)
	user := database.User{
		Username: register.Name,
		Password: registerPassword,
		Email:    register.Email,
	}
	err = app.Model.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Couldn't create user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
