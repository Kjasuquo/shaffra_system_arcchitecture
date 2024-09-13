package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"shaffra_assessment/cmd/server"
	"shaffra_assessment/internal/controller"
	"shaffra_assessment/internal/models"
	mock_database "shaffra_assessment/internal/repository/mocks"
	"shaffra_assessment/internal/service"
)

func TestGetUser(t *testing.T) {
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

	user := models.User{
		Models: models.Models{
			ID:        "d2164f8f-8331-4a2e-9aaf-525a5bd0a123",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		Name:  "JA Asuquo",
		Email: "okoasuquo@gmail.com",
		Age:   16,
	}

	expectedUserJSON, err := json.Marshal(user)
	assert.NoError(t, err)

	successfulResponse := fmt.Sprintf(`{
		"data": %v,
		"errors": "",
		"message": "user found successfully",
		"status": "OK"
	}`, string(expectedUserJSON))

	t.Run("getting user by ID successful", func(t *testing.T) {
		mockDB.EXPECT().GetUserByID(gomock.Any(), user.ID).Return(&user, nil).Times(1)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/users/"+user.ID, nil)
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)

		assert.JSONEq(t, successfulResponse, rw.Body.String())
	})

	t.Run("getting user by ID error", func(t *testing.T) {
		mockDB.EXPECT().GetUserByID(gomock.Any(), "d2164f8f-8331-4a2e-9aaf-525a5bd0a564").Return(nil, errors.New("error exists")).Times(1)

		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/users/d2164f8f-8331-4a2e-9aaf-525a5bd0a564", nil)
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		// Ensure correct status code and error message
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse := `{
			"data": null,
			"errors": "service---get user: error exists",
			"message": "user not found",
			"status": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})
}
