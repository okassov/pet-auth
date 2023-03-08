package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/entity"
	"github.com/okassov/pet-auth/internal/usecase"
)

type SignInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authRoutes struct {
	a usecase.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth) {
	r := &authRoutes{a}

	h := handler.Group("/auth")
	{
		h.POST("register", r.SignUp)
		h.POST("token", r.SignIn)
	}
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

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *authRoutes) SignIn(c *gin.Context) {

	inp := new(SignInput)

	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	tokenPair, err := r.a.SignIn(
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

	c.JSON(http.StatusOK, tokenPair)
}
