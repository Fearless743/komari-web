package models

type TaskExecLog struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID    string    `json:"task_id" gorm:"type:varchar(64);index:idx_task_exec_logs_task"`
	Command   string    `json:"command" gorm:"type:text"`
	AdminUUID string    `json:"admin_uuid" gorm:"type:varchar(36);index"`
	AdminIP   string    `json:"admin_ip" gorm:"type:varchar(45)"`
	StartTime LocalTime `json:"start_time" gorm:"type:timestamp;index"`
	EndTime   *LocalTime `json:"end_time" gorm:"type:timestamp"`
	Status    string    `json:"status" gorm:"type:varchar(20)"` // running, completed, failed, timeout
	TotalClients   int       `json:"total_clients" gorm:"default:0"`
	CompletedClients int   `json:"completed_clients" gorm:"default:0"`
	FailedClients    int   `json:"failed_clients" gorm:"default:0"`
}

type TaskExecDetailLog struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LogID     uint      `json:"log_id" gorm:"index"`
	TaskID    string    `json:"task_id" gorm:"type:varchar(64);index"`
	ClientUUID string   `json:"client_uuid" gorm:"type:varchar(36);index"`
	ClientName string   `json:"client_name" gorm:"type:varchar(100)"`
	StartTime LocalTime `json:"start_time" gorm:"type:timestamp"`
	EndTime   *LocalTime `json:"end_time" gorm:"type:timestamp"`
	ExitCode  *int      `json:"exit_code" gorm:"type:int"`
	Output    string    `json:"output" gorm:"type:longtext"`
	Status    string    `json:"status" gorm:"type:varchar(20)"` // running, success, failed, timeout, offline
	Error     string    `json:"error" gorm:"type:text"`
}
