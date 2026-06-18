package update

import (
	"github.com/gin-gonic/gin"
	"github.com/Fearless743/komari/api"
	"github.com/Fearless743/komari/database/auditlog"
	"github.com/Fearless743/komari/utils/geoip"
)

func UpdateMmdbGeoIP(c *gin.Context) {
	if err := geoip.UpdateDatabase(); err != nil {
		api.RespondError(c, 500, "Failed to update GeoIP database "+err.Error())
		return
	}
	uuid, _ := c.Get("uuid")
	auditlog.Log(c.ClientIP(), uuid.(string), "GeoIP database updated", "info")
	api.RespondSuccess(c, nil)
}
