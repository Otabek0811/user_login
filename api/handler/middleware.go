package handler

import (
	"app/pkg/helper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			c.String(http.StatusNotFound, "Cookie not found")
			return
		}

		// value := c.GetHeader("Authorization")

		info, err := helper.ParseClaims(value, h.cfg.AuthSecretKey)

		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}
		fmt.Println(info)

		c.Set("Auth", info)
		c.Next()
	}
}
