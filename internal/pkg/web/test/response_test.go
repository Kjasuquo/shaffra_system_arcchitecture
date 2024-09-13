package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shaffra_assessment/internal/pkg/web"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"shaffra_assessment/internal/models"
)

func TestJSON(t *testing.T) {
	// Create a new Gin engine and set to test mode
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	// Create a mock Gin context with the test HTTP response writer
	c, _ := gin.CreateTestContext(w)

	// Example data for testing
	message := "Success"
	status := http.StatusOK
	data := models.User{
		Models: models.Models{
			ID:        "d2164f8f-8331-4a2e-9aaf-525a5bd0a0ff",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		Name:  "JA Asuquo",
		Email: "okoasuquo@gmail.com",
		Age:   16,
	}

	// Call the JSON function
	web.JSON(c, message, status, data, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body into a map
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Extract the "data" from the response
	responseData := response["data"]

	// Convert the data to JSON
	jsonData, err := json.Marshal(responseData)
	assert.NoError(t, err)

	// Convert JSON to models.User struct
	var user models.User
	err = json.Unmarshal(jsonData, &user)
	assert.NoError(t, err)

	// Validate the JSON response
	assert.Equal(t, "Success", response["message"])
	assert.Equal(t, data, user)
	assert.Equal(t, "", response["errors"])
	assert.Equal(t, "OK", response["status"])
}
