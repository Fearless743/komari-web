package terminallogs

import (
	"time"

	"github.com/Fearless743/komari/database/dbcore"
	"github.com/Fearless743/komari/database/models"
)

func CreateSessionLog(sessionID, clientUUID, clientName, adminUUID, adminIP string) (*models.TerminalLog, error) {
	db := dbcore.GetDBInstance()
	log := &models.TerminalLog{
		SessionID:  sessionID,
		ClientUUID: clientUUID,
		ClientName: clientName,
		AdminUUID:  adminUUID,
		AdminIP:    adminIP,
		StartTime:  models.FromTime(time.Now()),
		Status:     "connected",
	}
	if err := db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateSessionEnd(sessionID string, endTime time.Time, status string, duration string, bytesSent, bytesRecv int64) error {
	db := dbcore.GetDBInstance()
	lt := models.FromTime(endTime)
	return db.Model(&models.TerminalLog{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"end_time":   lt,
			"status":     status,
			"duration":   duration,
			"bytes_sent": bytesSent,
			"bytes_recv": bytesRecv,
		}).Error
}

func GetTerminalLogs(limit, offset int) ([]models.TerminalLog, int64, error) {
	db := dbcore.GetDBInstance()
	var logs []models.TerminalLog
	var total int64
	if err := db.Model(&models.TerminalLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("start_time desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func GetTerminalLogByID(id uint) (*models.TerminalLog, error) {
	db := dbcore.GetDBInstance()
	var log models.TerminalLog
	if err := db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func GetTerminalLogsBySessionID(sessionID string) ([]models.TerminalLog, error) {
	db := dbcore.GetDBInstance()
	var logs []models.TerminalLog
	if err := db.Where("session_id = ?", sessionID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetTerminalDataLogs(sessionID string, limit int) ([]models.TerminalDataLog, error) {
	db := dbcore.GetDBInstance()
	var logs []models.TerminalDataLog
	query := db.Where("session_id = ?", sessionID)
	if limit > 0 {
		query = query.Order("timestamp desc").Limit(limit)
	}
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func CreateDataLog(sessionID, direction, dataSample string, dataSize int, isBinary bool) error {
	db := dbcore.GetDBInstance()
	log := &models.TerminalDataLog{
		SessionID:  sessionID,
		Direction:  direction,
		Timestamp:  models.FromTime(time.Now()),
		DataSize:   dataSize,
		DataSample: dataSample,
		IsBinary:   isBinary,
	}
	return db.Create(log).Error
}

func ClearOldTerminalLogs(before time.Time) error {
	db := dbcore.GetDBInstance()
	if err := db.Where("start_time < ?", before).Delete(&models.TerminalLog{}).Error; err != nil {
		return err
	}
	return db.Where("timestamp < ?", before).Delete(&models.TerminalDataLog{}).Error
}
