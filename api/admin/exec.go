package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/Fearless743/komari/api"
	"github.com/Fearless743/komari/database/clients"
	"github.com/Fearless743/komari/database/auditlog"
	"github.com/Fearless743/komari/database/models"
	"github.com/Fearless743/komari/database/tasks"
	"github.com/Fearless743/komari/database/taskexeclogs"
	"github.com/Fearless743/komari/utils"
	"github.com/Fearless743/komari/ws"
)

// 接受数据类型：
// - command: string
// - clients: []string (客户端 UUID 列表)
func Exec(c *gin.Context) {
	var req struct {
		Command string   `json:"command" binding:"required"`
		Clients []string `json:"clients" binding:"required"`
	}
	var onlineClients []string
	var offlineClients []string
	if err := c.ShouldBindJSON(&req); err != nil {
		api.RespondError(c, 400, "Invalid or missing request body: "+err.Error())
		return
	}
	adminUUID, _ := c.Get("uuid")
	for _, clientUuid := range req.Clients {
		if client := ws.GetConnectedClients()[clientUuid]; client != nil {
			onlineClients = append(onlineClients, clientUuid)
		} else {
			offlineClients = append(offlineClients, clientUuid)
		}
	}
	if len(onlineClients) == 0 {
		api.RespondError(c, 400, "No clients connected")
		return
	}
	taskId := utils.GenerateRandomString(16)
	totalClients := len(append(onlineClients, offlineClients...))
	
	// Create task execution log
	execLog, logErr := taskexeclogs.CreateTaskExecLog(taskId, req.Command, adminUUID.(string), c.ClientIP(), totalClients)
	if logErr != nil {
		fmt.Printf("Failed to create task exec log: %v", logErr)
	}

	// Create detail logs for each client
	for _, clientUUID := range append(onlineClients, offlineClients...) {
		clientName := clientUUID
		if client, err := clients.GetClientByUUID(clientUUID); err == nil {
			clientName = client.Name
		}
		status := "running"
		if !contain(onlineClients, clientUUID) {
			status = "offline"
		}
		_, detailErr := taskexeclogs.CreateDetailLog(execLog.ID, taskId, clientUUID, clientName, status)
		if detailErr != nil {
			fmt.Printf("Failed to create task exec detail log: %v", detailErr)
		}
	}

	if err := tasks.CreateTask(taskId, append(onlineClients, offlineClients...), req.Command); err != nil {
		api.RespondError(c, 500, "Failed to create task: "+err.Error())
		return
	}
	for _, clientUuid := range onlineClients {
		var send struct {
			Message string `json:"message"`
			Command string `json:"command"`
			TaskId  string `json:"task_id"`
		}
		send.Message = "exec"
		send.Command = req.Command
		send.TaskId = taskId

		payload, _ := json.Marshal(send)
		client := ws.GetConnectedClients()[clientUuid]
		if client != nil {
			if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
				api.RespondError(c, 400, "Client connection is broke: "+clientUuid)
				return
			}
		} else {
			api.RespondError(c, 400, "Client connection is null: "+clientUuid)
			return
		}
	}
	auditlog.Log(c.ClientIP(), adminUUID.(string), "REC, task id: "+taskId, "warn")
	api.RespondSuccess(c, gin.H{
		"task_id": taskId,
		"clients": onlineClients,
	})
	if len(offlineClients) > 0 {
		for _, clientUUID := range offlineClients {
			tasks.SaveTaskResult(taskId, clientUUID, "Client offline!", -1, models.FromTime(time.Now()))
			taskexeclogs.UpdateDetailLogByTaskAndClient(taskId, clientUUID, "Client offline!", -1, "offline", time.Now())
		}
	}
}

func contain(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// func contain(clients []string, uuid string) bool {
// 	for _, client := range clients {
// 		if client == uuid {
// 			return true
// 		}
// 	}
// 	return false
// }
