package taskexeclogs

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Fearless743/komari/api"
	"github.com/Fearless743/komari/database/taskexeclogs"
)

func GetTaskExecLogs(c *gin.Context) {
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
	logs, total, err := taskexeclogs.GetTaskExecLogs(limitInt, offset)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve task exec logs: "+err.Error())
		return
	}

	api.RespondSuccess(c, gin.H{"logs": logs, "total": total, "page": pageInt, "limit": limitInt})
}

func GetTaskExecLogDetails(c *gin.Context) {
	taskId := c.Param("task_id")

	execLog, err := taskexeclogs.GetTaskExecLogByTaskID(taskId)
	if err != nil {
		api.RespondError(c, 404, "Task exec log not found: "+err.Error())
		return
	}

	detailLogs, err := taskexeclogs.GetDetailLogsByTaskID(taskId)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve detail logs: "+err.Error())
		return
	}

	api.RespondSuccess(c, gin.H{"log": execLog, "details": detailLogs})
}

func ExportTaskExecLogsCSV(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "1000"
	}
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 || limitInt > 10000 {
		limitInt = 1000
	}

	logs, _, err := taskexeclogs.GetTaskExecLogs(limitInt, 0)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve task exec logs: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=task_exec_logs_"+strconv.Itoa(int(time.Now().Unix()))+".csv")
	c.Writer.WriteString("ID,Task ID,Command,Admin UUID,Admin IP,Start Time,End Time,Status,Total Clients,Completed Clients,Failed Clients\n")
	for _, log := range logs {
		endTime := ""
		if log.EndTime != nil {
			endTime = log.EndTime.ToTime().Format(time.RFC3339)
		}
		c.Writer.WriteString(strconv.FormatUint(uint64(log.ID), 10) + ",")
		c.Writer.WriteString(log.TaskID + ",")
		c.Writer.WriteString(`"` + escapeCSV(log.Command) + `",`)
		c.Writer.WriteString(log.AdminUUID + ",")
		c.Writer.WriteString(log.AdminIP + ",")
		c.Writer.WriteString(log.StartTime.ToTime().Format(time.RFC3339) + ",")
		c.Writer.WriteString(endTime + ",")
		c.Writer.WriteString(log.Status + ",")
		c.Writer.WriteString(strconv.Itoa(log.TotalClients) + ",")
		c.Writer.WriteString(strconv.Itoa(log.CompletedClients) + ",")
		c.Writer.WriteString(strconv.Itoa(log.FailedClients) + "\n")
	}
}

func ExportTaskExecDetailsCSV(c *gin.Context) {
	taskId := c.Param("task_id")

	details, err := taskexeclogs.GetDetailLogsByTaskID(taskId)
	if err != nil {
		api.RespondError(c, 500, "Failed to retrieve detail logs: "+err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=task_exec_details_"+taskId+".csv")
	c.Writer.WriteString("ID,Task ID,Client UUID,Client Name,Start Time,End Time,Exit Code,Output,Status,Error\n")
	for _, log := range details {
		endTime := ""
		if log.EndTime != nil {
			endTime = log.EndTime.ToTime().Format(time.RFC3339)
		}
		exitCode := ""
		if log.ExitCode != nil {
			exitCode = strconv.Itoa(*log.ExitCode)
		}
		c.Writer.WriteString(strconv.FormatUint(uint64(log.ID), 10) + ",")
		c.Writer.WriteString(log.TaskID + ",")
		c.Writer.WriteString(log.ClientUUID + ",")
		c.Writer.WriteString(log.ClientName + ",")
		c.Writer.WriteString(log.StartTime.ToTime().Format(time.RFC3339) + ",")
		c.Writer.WriteString(endTime + ",")
		c.Writer.WriteString(exitCode + ",")
		c.Writer.WriteString(`"` + escapeCSV(log.Output) + `",`)
		c.Writer.WriteString(log.Status + ",")
		c.Writer.WriteString(`"` + escapeCSV(log.Error) + `"` + "\n")
	}
}

func escapeCSV(s string) string {
	s = strconv.Quote(s)
	return s[1 : len(s)-1]
}
