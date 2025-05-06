package routes

import (
	"net/http"
	"rest-api/app/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event ID"})
		return
	}

	event, err := models.GetEventByID(eventID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event"})
		return
	}

	// Check if user is already registered
	isRegistered, err := event.IsUserRegistered(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check registration status"})
		return
	}
	if isRegistered {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User is already registered for this event"})
		return
	}

	err = event.Register(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event: " + err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

func cancelRegistration(context *gin.Context) {
	userID := context.GetInt64("userID")
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)

	var event models.Event
	event.ID = eventID

	err = event.CancelRegistration(userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel the registation"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "You have canceled registartion"})
}
