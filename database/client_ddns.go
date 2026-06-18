package database

import (
	"time"

	"github.com/Fearless743/komari/database/dbcore"
	"github.com/Fearless743/komari/database/models"
)

func UpdateClientDdnsRecordID(uuid string, recordID string) error {
	return dbcore.GetDBInstance().Model(&models.Client{}).Where("uuid = ?", uuid).Updates(map[string]any{
		"ddns_record_id": recordID,
		"updated_at":     time.Now(),
	}).Error
}
