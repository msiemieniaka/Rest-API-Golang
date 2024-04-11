package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"rest-api/app/models"
	"rest-api/app/utils"
)

func signup(context *gin.Context) {
	var user models.User
	
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data input"})
		return
	}
	
	err = user.Save()
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user"})
	}
	
	context.JSON(http.StatusCreated, gin.H{"message":"User has been created!"})
}

func login(context *gin.Context) {
var user models.User

err := context.ShouldBindJSON(&user)
if err != nil {
	context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse user data input"})
	return
}

err = user.ValidateCredentials()
if err != nil {
	context.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user data"})
	return
}

token, err := utils.GenerateToken(user.Email, user.ID)

if err != nil {
	context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token for user"})
	return
}

context.JSON(http.StatusOK, gin.H{"message":"Login successful!", "token": token})

}