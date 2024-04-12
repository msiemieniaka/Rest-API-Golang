package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/app/models"
	"strconv"
)

func registerForEvent (context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse single event"})
		return
	}
	
	event, err := models.GetEventByID(eventID)
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Could not fetch data"})
		return
	}
	
	
	err = event.Register(userID)
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Could not register user for event"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message":"User Registered!"})
}

func cancelRegistration(context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	
	var event models.Event
	event.ID = eventID
	
	err = event.CancelRegistration(userID)
	
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message":"Could not cancel the registation"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message":"You have canceled registartion"})
}