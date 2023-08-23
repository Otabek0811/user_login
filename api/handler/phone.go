package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// Create Phone godoc
// @ID create_phone
// @Router /v1/user/phone [POST]
// @Summary Create Phone
// @Description Create Phone
// @Tags Phone
// @Accept json
// @Produce json
// @Param phone body models.CreatePhone true "CreatePhoneRequest"
// @Success 201 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) CreatePhone(c *gin.Context) {

	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)

	var user_id string
	if len(userData.UserID) > 0 {
		user_id = userData.UserID
	}

	var createPhone models.CreatePhone

	err := c.ShouldBindJSON(&createPhone) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "create phone", http.StatusBadRequest, err.Error())
		return
	}

	createPhone.UserID = user_id

	if len(createPhone.Phone) > 12 {
		h.handlerResponse(c, "create phone", http.StatusBadRequest, "Invalid phone number")
		return
	}

	id, err := h.storages.Phone().Create(context.Background(), &createPhone)
	if err != nil {
		h.handlerResponse(c, "storage.phone.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storages.Phone().GetByID(context.Background(), &models.PhonePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.phone.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// Get By ID Phone godoc
// @ID get_by_id_phone
// @Router /v1/user/phone/{id} [GET]
// @Summary Get By ID Phone
// @Description Get By ID Phone
// @Tags Phone
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetByIdPhone(c *gin.Context) {

	_, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	
	var id string = c.Param("id")

	resp, err := h.storages.Phone().GetByID(context.Background(), &models.PhonePrimaryKey{
		Id: id,
	})
	if err != nil {
		h.handlerResponse(c, "storage.phone.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get phone by id", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// Get List Phone godoc
// @ID get_list_phone
// @Router /v1/user/phone [GET]
// @Summary Get List Phone
// @Description Get List Phone
// @Tags Phone
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) GetListPhone(c *gin.Context) {
	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)

	var user_id string
	if len(userData.UserID) > 0 {
		user_id = userData.UserID
	}

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list phone", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list phone", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.storages.Phone().GetList(context.Background(), &models.GetListPhoneRequest{
		UserID: user_id,
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	fmt.Println()
	fmt.Println(c.Query("search"))
	if err != nil {
		h.handlerResponse(c, "storage.phone.getlist", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list phone response", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update Phone godoc
// @ID update_phone
// @Router /v1/user/phone/{id} [PUT]
// @Summary Update Phone
// @Description Update Phone
// @Tags Phone
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param phone body models.UpdatePhone true "UpdatePhoneRequest"
// @Success 202 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) UpdatePhone(c *gin.Context) {
	val, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	userData := val.(helper.TokenInfo)

	var user_id string
	if len(userData.UserID) > 0 {
		user_id = userData.UserID
	}

	var updatePhone models.UpdatePhone

	id := c.Param("id")

	err := c.ShouldBindJSON(&updatePhone)
	if err != nil {
		h.handlerResponse(c, "update phone", http.StatusBadRequest, err.Error())
		return
	}

	updatePhone.Id = id
	updatePhone.UserID = user_id

	rowsAffected, err := h.storages.Phone().Update(context.Background(), &updatePhone)
	if err != nil {
		h.handlerResponse(c, "storage.phone.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.phone.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.storages.Phone().GetByID(context.Background(), &models.PhonePrimaryKey{
		Id:     id,
	})
	if err != nil {
		h.handlerResponse(c, "storage.phone.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "update phone", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// DELETE Phone godoc
// @ID delete_phone
// @Router /v1/user/phone/{id} [DELETE]
// @Summary Delete Phone
// @Description Delete Phone
// @Tags Phone
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param phone body models.PhonePrimaryKey true "DeletePhoneRequest"
// @Success 204 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) DeletePhone(c *gin.Context) {

	_, exists := c.Get("Auth")

	if !exists {
		h.handlerResponse(c, "get id in token", http.StatusInternalServerError, "invalid token")
		return
	}
	

	id := c.Param("id")

	rowsAffected, err := h.storages.Phone().Delete(context.Background(), &models.PhonePrimaryKey{
		Id:     id,
	})
	if err != nil {
		h.handlerResponse(c, "storage.phone.delete", http.StatusInternalServerError, err.Error())
		return
	}
	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.phone.delete", http.StatusBadRequest, "now rows affected")
		return
	}

	h.handlerResponse(c, "delete phone", http.StatusNoContent, nil)
}
