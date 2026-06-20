package terminal

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/Fearless743/komari/database/auditlog"
	"github.com/Fearless743/komari/database/terminallogs"
)

func ForwardTerminal(id string) {
	session, exists := TerminalSessions[id]

	if !exists || session == nil || session.Agent == nil || session.Browser == nil {
		return
	}
	auditlog.Log(session.RequesterIp, session.UserUUID, "established, terminal id:"+id, "terminal")
	established_time := time.Now()
	errChan := make(chan error, 1)
	var browserToAgentCount int
	var agentToBrowserCount int

	go func() {
		for {
			messageType, data, err := session.Browser.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			// Track bytes and log data samples
			isBinary := messageType == websocket.BinaryMessage
			session.BytesSent += int64(len(data))
			browserToAgentCount++
			if browserToAgentCount%100 == 0 {
				sample := string(data)
				if len(sample) > 200 {
					sample = sample[:200]
				}
				direction := "browser_to_agent"
				terminallogs.CreateDataLog(id, direction, sample, len(data), isBinary)
			}

			if messageType == websocket.TextMessage {
				if session.Agent != nil && string(data[0:1]) == "{" {
					err = session.Agent.WriteMessage(websocket.TextMessage, data)
				} else if session.Agent != nil {
					err = session.Agent.WriteMessage(websocket.BinaryMessage, data)
				}
			} else if session.Agent != nil {
				// 二进制消息，原样传递
				err = session.Agent.WriteMessage(websocket.BinaryMessage, data)
			}

			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		for {
			_, data, err := session.Agent.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			session.BytesRecv += int64(len(data))
			agentToBrowserCount++
			if agentToBrowserCount%100 == 0 {
				sample := string(data)
				if len(sample) > 200 {
					sample = sample[:200]
				}
				direction := "agent_to_browser"
				terminallogs.CreateDataLog(id, direction, sample, len(data), true)
			}
			if session.Browser != nil {
				err = session.Browser.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// 等待错误或主动关闭
	<-errChan
	// 关闭连接
	if session.Agent != nil {
		session.Agent.Close()
	}
	if session.Browser != nil {
		session.Browser.Close()
	}
	disconnect_time := time.Now()
	duration := disconnect_time.Sub(established_time).String()
	terminallogs.UpdateSessionEnd(id, disconnect_time, "disconnected", duration, session.BytesSent, session.BytesRecv)
	auditlog.Log(session.RequesterIp, session.UserUUID, "disconnected, terminal id:"+id+", duration:"+duration, "terminal")
	TerminalSessionsMutex.Lock()
	delete(TerminalSessions, id)
	TerminalSessionsMutex.Unlock()
}
