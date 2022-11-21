package controllers

import (
	"fmt"
	"jwt-authentication-golang/auth"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/models"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}
	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	t := time.Now()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	ip, err := net.LookupIP("ispycode.com")
	if err != nil {
		fmt.Println("host is unknown")
	}
	context.JSON(http.StatusOK, gin.H{"ip": ip, "timestamp": t.Format("2006-01-02 15:04:05"), "token": tokenString})
}
