package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shaffra_assessment/internal/models"
	"shaffra_assessment/internal/pkg/web"
)

// CreateUser gets payload from user for creating user
func (h *Handler) CreateUser(c *gin.Context) {

	start := time.Now()
	h.Wg.Add(1)
	defer h.trackRequestTime(c, h.Wg, start)

	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		web.JSON(c, "invalid input", http.StatusBadRequest, nil, err)
		return
	}

	result := make(chan web.Response)

	go func() {
		defer close(result)

		id, err := h.UserService.CreateUser(c, user)
		result <- web.Response{Data: id, Err: err}

	}()

	res := <-result

	if res.Err != nil {
		web.JSON(c, "cannot create user", http.StatusInternalServerError, nil, res.Err)
		return
	}
	web.JSON(c, "user created", http.StatusCreated, res.Data, nil)
}

// GetUser gets the id for retrieving the user
func (h *Handler) GetUser(c *gin.Context) {

	start := time.Now()
	h.Wg.Add(1)

	defer h.trackRequestTime(c, h.Wg, start)

	id := c.Param("id")

	result := make(chan web.Response)

	go func() {

		defer close(result)

		user, err := h.UserService.GetUserByID(c, id)
		result <- web.Response{Data: user, Err: err}

	}()

	res := <-result

	if res.Err != nil {
		web.JSON(c, "user not found", http.StatusInternalServerError, nil, res.Err)
		return
	}

	web.JSON(c, "user found successfully", http.StatusOK, res.Data, nil)

}

// UpdateUser gets the id and payload for updating the user
func (h *Handler) UpdateUser(c *gin.Context) {

	start := time.Now()
	h.Wg.Add(1)

	defer h.trackRequestTime(c, h.Wg, start)

	id := c.Param("id")

	user := &models.User{}

	err := c.ShouldBindJSON(user)
	if err != nil {
		web.JSON(c, "invalid input", http.StatusBadRequest, nil, err)
		return
	}

	result := make(chan web.Response)

	go func() {
		defer close(result)

		err = h.UserService.UpdateUserByID(c, id, user)
		result <- web.Response{Data: nil, Err: err}
	}()

	res := <-result

	if res.Err != nil {
		web.JSON(c, "cannot update user", http.StatusInternalServerError, nil, res.Err)
		return
	}

	web.JSON(c, "user updated successfully", http.StatusOK, user, nil)

}

// DeleteUser gets the id for deleting the user
func (h *Handler) DeleteUser(c *gin.Context) {

	start := time.Now()
	h.Wg.Add(1)

	defer h.trackRequestTime(c, h.Wg, start)

	id := c.Param("id")

	result := make(chan web.Response)

	go func() {
		defer close(result)

		err := h.UserService.DeleteUserID(c, id)
		result <- web.Response{Data: nil, Err: err}

	}()

	res := <-result

	if res.Err != nil {
		web.JSON(c, "cannot delete user", http.StatusInternalServerError, nil, res.Err)
		return
	}

	web.JSON(c, "user deleted successfully", http.StatusOK, nil, nil)
}
