package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"shaffra_assessment/cmd/server"
	"shaffra_assessment/internal/controller"
	"shaffra_assessment/internal/models"
	mock_database "shaffra_assessment/internal/repository/mocks"
	"shaffra_assessment/internal/service"
)

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_database.NewMockUserRepository(ctrl)
	serv := service.NewService(mockDB)
	h := &controller.Handler{
		UserService:     serv,
		Wg:              &sync.WaitGroup{},
		ReqDurationChan: make(chan string),
	}

	route := server.SetupRouter(h)

	// Payload sent in the request
	reqPayload := models.User{
		Name:  "Joseph Asuquo",
		Email: "okoasuquo@gmail.com",
		Age:   16,
	}

	// The user returned from the repository after the update
	returnedUser := models.User{
		Models: models.Models{
			ID:        "d2164f8f-8331-4a2e-9aaf-525a5bd0a000", // Make sure the ID is correctly set
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:  "Joseph Asuquo",
		Email: "okoasuquo@gmail.com",
		Age:   16,
	}

	bodyJSON, err := json.Marshal(reqPayload)
	assert.NoError(t, err)

	// Expected response (should be correctly formatted as JSON)
	expectedResponse := fmt.Sprintf(`{
		"data": {
			"ID": "%s",
			"name": "%s",
			"email": "%s",
			"age": %d,
			"created_at": "%s",
			"updated_at": "%s"
		},
		"errors": "",
		"message": "user updated successfully",
		"status": "OK"
	}`,
		reqPayload.ID, reqPayload.Name, reqPayload.Email, reqPayload.Age,
		reqPayload.CreatedAt.Format(time.RFC3339), reqPayload.UpdatedAt.Format(time.RFC3339))

	t.Run("update user successful", func(t *testing.T) {
		// Mock GetUserByID call
		mockDB.EXPECT().GetUserByID(gomock.Any(), returnedUser.ID).Return(&returnedUser, nil).Times(1)

		// Mock UpdateUserByID call
		mockDB.EXPECT().UpdateUserByID(gomock.Any(), returnedUser.ID, gomock.Any()).Return(nil).Times(1)

		// Perform request
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPut, "/users/"+returnedUser.ID, strings.NewReader(string(bodyJSON)))
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		// Assert response status
		assert.Equal(t, http.StatusOK, rw.Code)

		// Assert JSON response
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})

	t.Run("update user error occurs in get", func(t *testing.T) {
		// Mock GetUserByID call
		mockDB.EXPECT().GetUserByID(gomock.Any(), returnedUser.ID).Return(nil, errors.New("error exists")).Times(1)

		// Perform request
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPut, "/users/"+returnedUser.ID, strings.NewReader(string(bodyJSON)))
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		// Assert response status
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse = `{
			"data": null,
			"errors": "service---get user: error exists",
			"message": "cannot update user",
			"status": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})

	t.Run("update user error occurs in update", func(t *testing.T) {
		// Mock GetUserByID call
		mockDB.EXPECT().GetUserByID(gomock.Any(), returnedUser.ID).Return(&returnedUser, nil).Times(1)

		// Mock UpdateUserByID call
		mockDB.EXPECT().UpdateUserByID(gomock.Any(), returnedUser.ID, gomock.Any()).Return(errors.New("error exists")).Times(1)

		// Perform request
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPut, "/users/"+returnedUser.ID, strings.NewReader(string(bodyJSON)))
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		// Assert response status
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse = `{
			"data": null,
			"errors": "service---update user: error exists",
			"message": "cannot update user",
			"status": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})
}
