package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create User godoc
// @ID register_user
// @Router /register [POST]
// @Summary Register
// @Description Register
// @Tags Register
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "CreateUserRequest"
// @Success 201 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) RegisterUser(c *gin.Context) {

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

	hash, _ := h.HashPassword(createUser.Password)
	createUser.Password = hash

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{
		Login: createUser.Login,
	})
	fmt.Println("Hello")
	fmt.Println(resp)
	if err == nil {
		h.handlerResponse(c, "register user", http.StatusBadRequest, "user already exists")
		return
	} else {
		if err.Error() == "no rows in result set" {
		} else {
			h.handlerResponse(c, "get user by login", http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := h.storages.User().Create(context.Background(), &createUser)
	if err != nil {
		h.handlerResponse(c, "storage.user.register", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err = h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.user.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @ID login_user
// @Router /login [POST]
// @Summary Login
// @Description Login
// @Tags Login
// @Accept json
// @Produce json
// @Param user body models.Login true "LoginRequest"
// @Success 201 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server Error"
func (h *Handler) LoginUser(c *gin.Context) {

	var logPass models.Login

	err := c.ShouldBindJSON(&logPass) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "login user", http.StatusBadRequest, err.Error())
		return
	}
	fmt.Printf("%+v", logPass)
	fmt.Println(len(logPass.Login), len(logPass.Password))

	if len(logPass.Login) < 6 || len(logPass.Password) < 6 {
		fmt.Println("modmcso")
		h.handlerResponse(c, "login user", http.StatusBadRequest, "Login and Password length must be longer than 6")
		return
	}

	resp, err := h.storages.User().GetByID(context.Background(), &models.UserPrimaryKey{
		Login: logPass.Login,
	})
	if err != nil {
		if err.Error() == "no rows in result set" {
			h.handlerResponse(c, "login user", http.StatusBadRequest, "this user does not exist")
			return
		} else {
			h.handlerResponse(c, "login user", http.StatusBadRequest, err.Error())
			return
		}

	}

	// if resp.Password != logPass.Password {
	// 	h.handlerResponse(c, "login user", http.StatusBadRequest, "Incorrect Password")
	// 	return
	// }
	if !h.CheckPasswordHash(logPass.Password, resp.Password) {
		h.handlerResponse(c, "login user", http.StatusBadRequest, "Incorrect Password")
		return
	}



	m := make(map[string]interface{})

	m["user_id"] = resp.Id

	token, err := helper.GenerateJWT(m, time.Hour*24, "secret")

	if err != nil {
		h.handlerResponse(c, "token response", http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, token)
}
