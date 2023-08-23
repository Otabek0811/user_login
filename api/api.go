package api

import (
	_ "app/api/docs"
	"app/api/handler"

	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func NewApi(r *gin.Engine, cfg *config.Config, store storage.StorageI, logger logger.LoggerI) {

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	handler := handler.NewHandler(cfg, store, logger)

	r.Use(customCORSMiddleware())

	v1 := r.Group("/v1")

	

	// register
	r.POST("/register", handler.RegisterUser)

	// login
	r.POST("/login", handler.LoginUser)

	// user api
	v1.Use(handler.AuthMiddleware())
	v1.POST("/user", handler.CreateUser)
	v1.GET("/user/:id", handler.GetByIdUser)
	v1.GET("/user", handler.GetListUser)
	v1.PUT("/user/:id", handler.UpdateUser)
	v1.DELETE("/user/:id", handler.DeleteUser)

	// phone api
	v1.POST("/user/phone", handler.CreatePhone)
	v1.GET("/user/phone/:id", handler.GetByIdPhone)
	v1.GET("/user/phone", handler.GetListPhone)
	v1.PUT("/user/phone/:id", handler.UpdatePhone)
	v1.DELETE("/v1/user/phone/:id", handler.DeletePhone)



	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Accesp-Encoding, Authorization, Cache-Control")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

