package controllers

import (
	"crud_go_gin/initailizer"
	"crud_go_gin/middleware"
	"crud_go_gin/models"
	"crud_go_gin/util"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type signUpPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

var jwtSecret = os.Getenv("JWT_SECRET")

func signUp(c *gin.Context) {
	// get data from request body
	var body signUpPayload

	if !util.BindAndValidate(c, &body) {
		return
	}

	// hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		util.ResponseWithError(c, http.StatusBadRequest, "Failed to hash password")
		return
	}

	// create a user
	user := models.User{Email: body.Email, Password: string(hashPassword)}
	if result := initailizer.DB.Create(&user); result.Error != nil {
		util.ResponseWithError(c, http.StatusBadRequest, "Failed to create user")
		return
	}

	// response it
	util.ResponseWithJson(c, http.StatusOK, "")
}

func login(c *gin.Context) {
	// get data from req body
	var body signUpPayload

	if !util.BindAndValidate(c, &body) {
		return
	}

	// compare the password with the database hash
	var user models.User
	initailizer.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		util.ResponseWithError(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		util.ResponseWithError(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	// create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		util.ResponseWithError(c, http.StatusBadRequest, "failed to create token")
		return
	}

	// response it
	// 酷東西
	//? https://medium.com/%E7%A8%8B%E5%BC%8F%E7%8C%BF%E5%90%83%E9%A6%99%E8%95%89/%E5%86%8D%E6%8E%A2%E5%90%8C%E6%BA%90%E6%94%BF%E7%AD%96-%E8%AB%87-samesite-%E8%A8%AD%E5%AE%9A%E5%B0%8D-cookie-%E7%9A%84%E5%BD%B1%E9%9F%BF%E8%88%87%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A0%85-6195d10d4441

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*24*30, "/", "", false, true)

	util.ResponseWithJson(c, http.StatusOK, "")
}

func validate(c *gin.Context) {
	user, _ := c.Get("user")

	util.ResponseWithJson(c, http.StatusOK, user)
}

func SetupAuthRouter(r *gin.Engine) {
	authRouters := r.Group("/")
	{
		authRouters.POST("/signup", signUp)
		authRouters.POST("/login", login)
		authRouters.GET("/validate", middleware.RequireAuth, validate)
	}
}
