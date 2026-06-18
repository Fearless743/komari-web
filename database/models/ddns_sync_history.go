package models

import "time"

type DdnsSyncHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ClientUUID    string    `json:"client_uuid" gorm:"type:varchar(36);index"`
	ClientName    string    `json:"client_name" gorm:"type:varchar(100)"`
	Hostname      string    `json:"hostname" gorm:"type:varchar(255)"`
	RecordType    string    `json:"record_type" gorm:"type:varchar(10)"`
	IPV4          string    `json:"ipv4" gorm:"type:varchar(100)"`
	IPV6          string    `json:"ipv6" gorm:"type:varchar(100)"`
	RecordID      string    `json:"record_id" gorm:"type:varchar(255)"`
	Status        string    `json:"status" gorm:"type:varchar(20)"`
	Error         string    `json:"error" gorm:"type:text"`
	TriggeredBy   string    `json:"triggered_by" gorm:"type:varchar(50)"`
	SyncedAt      time.Time `json:"synced_at" gorm:"autoCreateTime"`
}

func (DdnsSyncHistory) TableName() string {
	return "ddns_sync_history"
}
