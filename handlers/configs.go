package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/cache"
)

type ConfigsHandler struct {
	Cache *cache.Configs
}

func (h *ConfigsHandler) RegistrationEnabled(c *gin.Context) {
	configs, err := h.Cache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_allowed": configs.RegistrationEnabled,
	})
}
