package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okassov/pet-auth/internal/entity"
	"github.com/okassov/pet-auth/internal/usecase"
	"github.com/okassov/pet-auth/pkg/logger"
)

type SignInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authRoutes struct {
	a usecase.Auth
	l logger.LoggerInterface
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.Auth, l logger.LoggerInterface) {

	r := &authRoutes{a, l}

	h := handler.Group("/auth")
	{
		h.POST("register", r.SignUp)
		h.POST("token", r.SignIn)
	}
}

type SignUpResponse struct {
	Message string `json:"message"`
}

//	@Summary		Sign UP
//	@Description	Register User
//	@ID				SignUp
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	SignUpResponse
//	@Failure		500	{object}	response
//	@Router			/auth/register [get]
func (r *authRoutes) SignUp(c *gin.Context) {

	inp := new(SignInput)

	if err := c.ShouldBindJSON(&inp); err != nil {
		r.l.Error("http - v1 - SignUp", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
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
		r.l.Error("http - v1 - SignUp", err)
		errorResponse(c, http.StatusUnauthorized, "database problems")
	}

	// c.JSON(http.StatusOK, "{ message: User signed up }")
	c.JSON(http.StatusOK, SignUpResponse{Message: "User registered"})

}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//	@Summary		Sign In
//	@Description	Authorized User and return token
//	@ID				SignIn
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SignInResponse	true	"Generate Tokens"
//	@Success		200		{object}	SignInResponse
//	@Failure		400		{object}	response
//	@Failure		500		{object}	response
//	@Router			/auth/token [post]
func (r *authRoutes) SignIn(c *gin.Context) {

	inp := new(SignInput)

	if err := c.ShouldBindJSON(&inp); err != nil {
		r.l.Error("http - v1 - SignIn", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
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
		r.l.Error("http - v1 - SignIn", err)
		errorResponse(c, http.StatusUnauthorized, "database problems")
	}

	// c.JSON(http.StatusOK, tokenPair)
	c.JSON(http.StatusOK, SignInResponse{
		AccessToken:  tokenPair["access_token"],
		RefreshToken: tokenPair["refresh_token"],
	})
}
