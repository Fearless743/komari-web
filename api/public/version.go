package public

import (
	"github.com/gin-gonic/gin"
	"github.com/Fearless743/komari/api"
	"github.com/Fearless743/komari/utils"
)

func GetVersion(c *gin.Context) {
	api.RespondSuccess(c, gin.H{
		"version": utils.CurrentVersion,
		"hash":    utils.VersionHash,
	})
}
