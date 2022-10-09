package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
)

type ConfigsHandler struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
}

func (h *ConfigsHandler) RegistrationEnabled(c *gin.Context) {
	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_enabled": configs.RegistrationEnabled,
	})
}

func (h *ConfigsHandler) TellerAppID(c *gin.Context) {
	if _, _, err := auth.AuthorizeUser(c, h.Db); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teller_app_id": configs.TellerApplicationID,
	})
}

func (h *ConfigsHandler) GetConfigs(c *gin.Context) {
	if _, userType, err := auth.AuthorizeUser(c, h.Db); err != nil || *userType != models.AdminUser {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ConfigsHandler) UpdateConfigs(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	if _, userType, err := auth.AuthorizeUser(c, h.Db); err != nil || *userType != models.AdminUser {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var input *cache.ConfigsInput
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.ConfigsCache.UpdateConfigsCache(ctx, h.Db, input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, input)
}
