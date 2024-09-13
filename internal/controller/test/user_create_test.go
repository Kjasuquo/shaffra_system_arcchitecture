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

	"shaffra_assessment/cmd/server"
	"shaffra_assessment/internal/controller"
	"shaffra_assessment/internal/models"
	mock_database "shaffra_assessment/internal/repository/mocks"
	"shaffra_assessment/internal/service"
)

func TestCreateUser(t *testing.T) {
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

	reqPayload := models.User{
		Name:  "JA Asuquo",
		Email: "okoasuquo@gmail.com",
		Age:   16,
	}

	bodyJSON, err := json.Marshal(reqPayload)
	assert.NoError(t, err)

	newID := "d2164f8f-8331-4a2e-9aaf-525a5bd0a000"

	successfulResponse := fmt.Sprintf(`{
		"data": "%s",
		"errors": "",
		"message": "user created",
		"status": "Created"
	}`, newID)

	t.Run("create user successful", func(t *testing.T) {

		mockDB.EXPECT().CreateUser(gomock.Any(), &reqPayload).Return(newID, nil).Times(1)

		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/users", strings.NewReader(string(bodyJSON)))
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusCreated, rw.Code)
		assert.JSONEq(t, successfulResponse, rw.Body.String())
	})

	t.Run("create user error", func(t *testing.T) {

		mockDB.EXPECT().CreateUser(gomock.Any(), &reqPayload).Return("", errors.New("error exists")).Times(1)

		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/users", strings.NewReader(string(bodyJSON)))
		assert.NoError(t, err)

		route.ServeHTTP(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse := `{
			"data": null,
			"errors": "service---create user: error exists",
			"message": "cannot create user",
			"status": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
	})
}
