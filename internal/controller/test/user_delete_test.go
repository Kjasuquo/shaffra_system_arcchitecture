package test

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"shaffra_assessment/cmd/server"
	"shaffra_assessment/internal/controller"
	mock_database "shaffra_assessment/internal/repository/mocks"
	"shaffra_assessment/internal/service"
)

func TestDeleteUser(t *testing.T) {
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

	id := "d2164f8f-8331-4a2e-9aaf-525a5bd0a123"

	successfulResponse := fmt.Sprintf(`{
		"data": null,
		"errors": "",
		"message": "user deleted successfully",
		"status": "OK"
	}`)

	t.Run("delete user by ID successful", func(t *testing.T) {
		mockDB.EXPECT().DeleteUserID(gomock.Any(), id).Return(nil).Times(1)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)

		// Ensure that the response matches the expected JSON
		assert.JSONEq(t, successfulResponse, rw.Body.String())
	})

	t.Run("delete user by ID error", func(t *testing.T) {
		mockDB.EXPECT().DeleteUserID(gomock.Any(), id).Return(errors.New("error exists")).Times(1)

		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/users/"+id, nil)
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		// Ensure correct status code and error message
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse := `{
			"data": null,
			"errors": "service---delete user: error exists",
			"message": "cannot delete user",
			"status": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})
}
