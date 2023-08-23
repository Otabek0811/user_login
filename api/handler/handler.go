package handler

import (
	"app/config"
	"app/pkg/logger"
	"app/storage"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	cfg      *config.Config
	logger   logger.LoggerI
	storages storage.StorageI
}

type Response struct {
	Status      int
	Description string
	Data        interface{}
}

func NewHandler(cfg *config.Config, store storage.StorageI, logger logger.LoggerI) *Handler {
	return &Handler{
		cfg:      cfg,
		logger:   logger,
		storages: store,
	}
}

func (h *Handler) handlerResponse(c *gin.Context, path string, code int, message interface{}) {
	response := Response{
		Status:      code,
		Description: path,
		Data:        message,
	}

	switch {
	case code < 300:
		h.logger.Info(path, logger.Any("info", response.Description))
	case code >= 400:
		h.logger.Error(path, logger.Any("info", response))
	}

	c.JSON(code, response)
}

func (h *Handler) getOffsetQuery(offset string) (int, error) {
	if len(offset) <= 0 {
		return h.cfg.DefaultOffset, nil
	}

	return strconv.Atoi(offset)
}

func (h *Handler) getLimitQuery(limit string) (int, error) {

	if len(limit) <= 0 {
		return h.cfg.DefaultLimit, nil
	}

	return strconv.Atoi(limit)
}

func (h *Handler) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *Handler) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
