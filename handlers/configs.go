package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/cache"
)

type ConfigHandler struct {
	Cache *cache.ConfigCache
}

func (h *ConfigHandler) RegistrationAllowed(c *gin.Context) {
	isAllowed, err := h.Cache.GetRegistrationAllowed()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_allowed": isAllowed,
	})
}
