package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/entity"
	"github.com/okassov/pet-auth/internal/usecase"
)

type authRoutes struct {
	a usecase.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth) {
	r := &authRoutes{a}

	h := handler.Group("/auth")
	{
		h.POST("register", r.SignUp)
		h.POST("login", r.SignIn)
	}
}

func (r *authRoutes) generateToken(c *gin.Context) {

}

type SignInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *authRoutes) SignUp(c *gin.Context) {

	inp := new(SignInput)

	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err := r.a.SignUp(
		c.Request.Context(),
		entity.User{
			Name:     inp.Name,
			Username: inp.Username,
			Email:    inp.Email,
			Password: inp.Password,
		},
	)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "{ message: User signed up }")

}

func (r *authRoutes) SignIn(c *gin.Context) {

	inp := new(SignInput)

	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	user, err := r.a.SignIn(
		c.Request.Context(),
		entity.User{
			Name:     inp.Name,
			Username: inp.Username,
			Email:    inp.Email,
			Password: inp.Password,
		},
	)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	fmt.Println(user)

	c.JSON(http.StatusOK, "{ message: User login success }")
}
