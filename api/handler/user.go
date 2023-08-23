package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// Create User godoc
// @ID create_user
// @Router /v1/user [POST]
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "CreateUserRequest"
// @Success 201 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CreateUser(c *gin.Context) {

	var createUser models.CreateUser

	err := c.ShouldBindJSON(&createUser) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "create user", http.StatusBadRequest, err.Error())
		return
	}
	if len(createUser.Login) < 6 || len(createUser.Password) < 6 {
		h.handlerResponse(c, "register user", http.StatusBadRequest, "Login and Password length must be longer than 6")
		return
	}

	id, err := h.storages.User().Create(context.Background(), &createUser)
	if err != nil {
		h.handlerResponse(c, "storage.user.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// Get By ID User godoc
// @ID get_by_id_user
// @Router /v1/user/{id} [GET]
// @Summary Get By ID User
// @Description Get By ID User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetByIdUser(c *gin.Context) {

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)

	var id string
	if len(userData.UserID) > 0 {
		id = userData.UserID
	} else {
		id = c.Param("id")
	}

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get user by id", http.StatusCreated, resp)
}



// @Security ApiKeyAuth
// Get By Name User godoc
// @ID get_by_name_user
// @Router /v1/user/{name} [GET]
// @Summary Get By Name User
// @Description Get By Name User
// @Tags User
// @Accept json
// @Produce json
// @Param name path string false "name"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetByNameUser(c *gin.Context) {

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get by name in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)
	
	var name string = c.Param("name")

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{
		Name: name,
		UserID: userData.UserID,
	})
	if err != nil {
		h.handlerResponse(c, "storage.user.getByName", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get user by Name", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// Get List User godoc
// @ID get_list_user
// @Router /v1/user [GET]
// @Summary Get List User
// @Description Get List User
// @Tags User
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetListUser(c *gin.Context) {

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get by name in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list user", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list user", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.storages.User().GetList(context.Background(), &models.GetListUserRequest{
		UserID: userData.UserID,
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.user.getlist", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list user response", http.StatusOK, resp)
}


// @Security ApiKeyAuth
// Update User godoc
// @ID update_user
// @Router /v1/user/{id} [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "id"
// @Param user body models.UpdateUser true "UpdateUserRequest"
// @Success 202 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) UpdateUser(c *gin.Context) {
	
	var updateUser models.UpdateUser

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)


	var id string
	if len(userData.UserID) > 0 {
		id = userData.UserID
	} else {
		id = c.Param("id")
	}


	err := c.ShouldBindJSON(&updateUser)
	if err != nil {
		h.handlerResponse(c, "update user", http.StatusBadRequest, err.Error())
		return
	}

	updateUser.Id = id

	if len(updateUser.Login) < 6 || len(updateUser.Password) < 6 {
		h.handlerResponse(c, "register user", http.StatusBadRequest, "Login and Password length must be longer than 6")
		return
	}

	rowsAffected, err := h.storages.User().Update(context.Background(), &updateUser)
	if err != nil {
		h.handlerResponse(c, "storage.user.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.user.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "update user", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// DELETE User godoc
// @ID delete_user
// @Router /v1/user/{id} [DELETE]
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept json
// @Produce json
// @Param id path string false "id"
// @Param user body models.UserPrimaryKey true "DeleteUserRequest"
// @Success 204 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) DeleteUser(c *gin.Context) {

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)


	var id string
	if len(userData.UserID) > 0 {
		id = userData.UserID
	} else {
		id = c.Param("id")
	}

	rowsAffected, err := h.storages.User().Delete(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.delete", http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.user.delete", http.StatusBadRequest, "now rows affected")
		return
	}

	h.handlerResponse(c, "delete user", http.StatusNoContent, nil)
}
