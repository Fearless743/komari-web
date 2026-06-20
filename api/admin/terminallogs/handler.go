package terminallogs

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Fearless743/komari/api"
	"github.com/Fearless743/komari/database/terminallogs"
)

func GetTerminalLogs(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "50"
	}
	page := c.Query("page")
	if page == "" {
		page = "1"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		api.RespondError(c, 400, "Invalid limit: "+limit)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		api.RespondError(c, 400, "Invalid page: "+page)
		return
	}

	offset := (pageInt - 1) * limitInt
	logs, total, err := terminallogs.GetTerminalLogs(limitInt, offset)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve terminal logs: "+err.Error())
		return
	}

	api.RespondSuccess(c, gin.H{"logs": logs, "total": total, "page": pageInt, "limit": limitInt})
}

func GetTerminalLogDetails(c *gin.Context) {
	sessionID := c.Param("session_id")
	
	logs, err := terminallogs.GetTerminalLogsBySessionID(sessionID)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve terminal logs: "+err.Error())
		return
	}

	dataLogs, err := terminallogs.GetTerminalDataLogs(sessionID, 1000)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve terminal data logs: "+err.Error())
		return
	}

	api.RespondSuccess(c, gin.H{"logs": logs, "data_logs": dataLogs})
}

func ExportTerminalLogsCSV(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "1000"
	}
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 10000 {
		limitInt = 1000
	}

	logs, _, err := terminallogs.GetTerminalLogs(limitInt, 0)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve terminal logs: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=terminal_logs_"+strconv.Itoa(int(time.Now().Unix()))+".csv")
	c.Writer.WriteString("ID,Session ID,Client UUID,Client Name,Admin UUID,Admin IP,Start Time,End Time,Duration,Status,Bytes Sent,Bytes Received\n")
	for _, log := range logs {
		endTime := ""
		if log.EndTime != nil {
			endTime = log.EndTime.ToTime().Format(time.RFC3339)
		}
		c.Writer.WriteString(strconv.FormatUint(uint64(log.ID), 10) + ",")
		c.Writer.WriteString(log.SessionID + ",")
		c.Writer.WriteString(log.ClientUUID + ",")
		c.Writer.WriteString(log.ClientName + ",")
		c.Writer.WriteString(log.AdminUUID + ",")
		c.Writer.WriteString(log.AdminIP + ",")
		c.Writer.WriteString(log.StartTime.ToTime().Format(time.RFC3339) + ",")
		c.Writer.WriteString(endTime + ",")
		c.Writer.WriteString(log.Duration + ",")
		c.Writer.WriteString(log.Status + ",")
		c.Writer.WriteString(strconv.FormatInt(log.BytesSent, 10) + ",")
		c.Writer.WriteString(strconv.FormatInt(log.BytesRecv, 10) + "\n")
	}
}
