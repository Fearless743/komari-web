package database

import (
	"github.com/Fearless743/komari/database/dbcore"
	"github.com/Fearless743/komari/database/models"
)

func SaveDdnsSyncHistory(history models.DdnsSyncHistory) error {
	db := dbcore.GetDBInstance()
	return db.Create(&history).Error
}

func GetDdnsSyncHistory(clientUUID string, limit int) ([]models.DdnsSyncHistory, error) {
	db := dbcore.GetDBInstance()
	var histories []models.DdnsSyncHistory
	if limit <= 0 {
		limit = 50
	}
	if clientUUID != "" {
		if err := db.Where("client_uuid = ?", clientUUID).
			Order("id DESC").
			Limit(limit).
			Find(&histories).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Order("id DESC").
			Limit(limit).
			Find(&histories).Error; err != nil {
			return nil, err
		}
	}
	for i := range histories {
		histories[i].SyncedAt = histories[i].SyncedAt.Local()
	}
	return histories, nil
}

func DeleteDdnsSyncHistoryBefore(beforeTime models.LocalTime) error {
	db := dbcore.GetDBInstance()
	return db.Where("synced_at < ?", beforeTime).Delete(&models.DdnsSyncHistory{}).Error
}
