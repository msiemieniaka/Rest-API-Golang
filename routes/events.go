package routes

import (
	"net/http"
	"rest-api/app/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the data"})
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse single event"})
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch single event"})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not post data"})
		return
	}

	userID := context.GetInt64("userID")
	event.UserID = userID

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save the data to table"})
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func updateEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse single event"})
		return
	}

	_, err = models.GetEventByID(eventID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch single event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not post data"})
		return
	}
	updatedEvent.ID = eventID
	err = updatedEvent.Update()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not update event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Sucesfully updated an event"})
}

func deleteEvents(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse single event"})
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch single event"})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not delete event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Sucesfully deleted an event"})
}
