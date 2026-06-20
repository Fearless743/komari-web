package taskexeclogs

import (
	"time"

	"github.com/Fearless743/komari/database/dbcore"
	"github.com/Fearless743/komari/database/models"
)

func CreateTaskExecLog(taskID, command, adminUUID, adminIP string, totalClients int) (*models.TaskExecLog, error) {
	db := dbcore.GetDBInstance()
	log := &models.TaskExecLog{
		TaskID:       taskID,
		Command:      command,
		AdminUUID:    adminUUID,
		AdminIP:      adminIP,
		StartTime:    models.FromTime(time.Now()),
		Status:       "running",
		TotalClients: totalClients,
	}
	if err := db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateTaskExecLogEnd(taskID string, endTime time.Time, status string, completed, failed int) error {
	db := dbcore.GetDBInstance()
	lt := models.FromTime(endTime)
	return db.Model(&models.TaskExecLog{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"end_time":           lt,
			"status":             status,
			"completed_clients":  completed,
			"failed_clients":     failed,
		}).Error
}

func CreateDetailLog(logID uint, taskID, clientUUID, clientName, status string) (*models.TaskExecDetailLog, error) {
	db := dbcore.GetDBInstance()
	detail := &models.TaskExecDetailLog{
		LogID:      logID,
		TaskID:     taskID,
		ClientUUID: clientUUID,
		ClientName: clientName,
		StartTime:  models.FromTime(time.Now()),
		Status:     status,
	}
	if err := db.Create(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func UpdateDetailLogByTaskAndClient(taskID, clientUUID, output string, exitCode int, status string, endTime time.Time) error {
	db := dbcore.GetDBInstance()
	updates := map[string]interface{}{
		"output":    output,
		"exit_code": exitCode,
		"status":    status,
	}
	if !endTime.IsZero() {
		updates["end_time"] = models.FromTime(endTime)
	}
	return db.Model(&models.TaskExecDetailLog{}).
		Where("task_id = ? AND client_uuid = ?", taskID, clientUUID).
		Updates(updates).Error
}

func GetTaskExecLogs(limit, offset int) ([]models.TaskExecLog, int64, error) {
	db := dbcore.GetDBInstance()
	var logs []models.TaskExecLog
	var total int64
	if err := db.Model(&models.TaskExecLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("start_time desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func GetTaskExecLogByTaskID(taskID string) (*models.TaskExecLog, error) {
	db := dbcore.GetDBInstance()
	var log models.TaskExecLog
	if err := db.Where("task_id = ?", taskID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func GetDetailLogsByTaskID(taskID string) ([]models.TaskExecDetailLog, error) {
	db := dbcore.GetDBInstance()
	var logs []models.TaskExecDetailLog
	if err := db.Where("task_id = ?", taskID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetDetailLogByTaskAndClient(taskID, clientUUID string) (*models.TaskExecDetailLog, error) {
	db := dbcore.GetDBInstance()
	var log models.TaskExecDetailLog
	if err := db.Where("task_id = ? AND client_uuid = ?", taskID, clientUUID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func ClearOldTaskExecLogs(before time.Time) error {
	db := dbcore.GetDBInstance()
	if err := db.Where("start_time < ?", before).Delete(&models.TaskExecLog{}).Error; err != nil {
		return err
	}
	return db.Where("start_time < ?", before).Delete(&models.TaskExecDetailLog{}).Error
}
