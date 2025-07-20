package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yiheyistm/go-restful-api/internal/database"
	"github.com/gin-gonic/gin"
)

// GetEvents returns all events
//
//	@Summary		Returns all events
//	@Description	Returns all events
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]database.Event
//	@Router			/api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.Model.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetEvent returns a single event
//
//	@Summary		Returns a single event
//	@Description	Returns a single event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]
func (app *application) getEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	event, err := app.Model.Events.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusOK, event)
}

// CreateEvent creates a new event
//
//	@Summary		Creates a new event
//	@Description	Creates a new event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		database.Event	true	"Event"
//	@Success		201		{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
func (app *application) createEvent(c *gin.Context) {
	var newEvent database.Event
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user := app.GetUserFromContext(c)
	newEvent.OwnerId = user.ID
	err := app.Model.Events.Insert(&newEvent)
	if err != nil {
		log.Println("Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, newEvent)
}

// UpdateEvent updates an existing event
//
//	@Summary		Updates an existing event
//	@Description	Updates an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success		200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
func (app *application) updateEvent(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	user := app.GetUserFromContext(c)
	existedEvent, err := app.Model.Events.GetByID(id)
	if err != nil || existedEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existedEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this event"})
		return
	}
	updatedEvent := &database.Event{}
	updatedEvent.ID = id
	updatedEvent.OwnerId = user.ID
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = app.Model.Events.Update(updatedEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an existing event
//
//	@Summary		Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	user := app.GetUserFromContext(c)
	existedEvent, err := app.Model.Events.GetByID(id)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if existedEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this event"})
		return
	}
	err = app.Model.Events.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// AddAttendeeToEvent adds an attendee to an event
//
//	@Summary		Adds an attendee to an event
//	@Description	Adds an attendee
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		200		{object}	database.Attendee
//	@Router			/api/v1/events/{id}/attendees/{userId} [post]
//	@Security		BearerAuth
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user := app.GetUserFromContext(c)
	event, err := app.Model.Events.GetByID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if user.ID != event.OwnerId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add attendees to this event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user, err = app.Model.Users.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existedAttendee, err := app.Model.Attendees.GetByEventAndUserId(event.ID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check attendee"})
		return
	}
	if existedAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already an attendee"})
		return
	}
	attendee := &database.Attendee{
		EventID: event.ID,
		UserID:  user.ID,
	}
	err = app.Model.Attendees.Insert(attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendee added successfully", "attendee": attendee})
}

// GetAttendeesForEvent retrieves all attendees for a specific event
//
//	@Summary		Retrieves all attendees for a specific event
//	@Description	Retrieves all attendees for a specific event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{array}	database.Attendee
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	attendees, err := app.Model.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees"})
		return
	}
	c.JSON(http.StatusOK, attendees)
}

// DeleteAttendeeFromEvent removes an attendee from an event
//
//	@Summary		Removes an attendee from an event
//	@Description	Removes an attendee from an event
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			userId	path		int	true	"User ID"
//	@Success		204
//	@Router			/api/v1/events/{id}/attendees/{userId} [delete]
//	@Security		BearerAuth
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user := app.GetUserFromContext(c)
	existedEvent, err := app.Model.Events.GetByID(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to retrieve Event"})
		return
	}

	if user.ID != existedEvent.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this event"})
		return
	}

	existedAttendee, err := app.Model.Attendees.GetByEventAndUserId(eventID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check attendee"})
		return
	}
	if existedAttendee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
		return
	}

	err = app.Model.Attendees.Delete(existedAttendee.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetEventsByAttendee retrieves all events for a specific attendee
//
//	@Summary		Retrieves all events for a specific attendee
//	@Description	Retrieves all events for a specific attendee
//	@Tags			attendees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Attendee ID"
//	@Success		200	{array}	database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}
	attendee, err := app.Model.Attendees.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}
	if attendee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendee not found"})
		return
	}

	events, err := app.Model.Events.GetByAttendeeId(attendee.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events for attendee"})
		return
	}
	c.JSON(http.StatusOK, events)
}
