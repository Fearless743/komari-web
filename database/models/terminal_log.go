package models

type TerminalLog struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID   string    `json:"session_id" gorm:"type:varchar(64);index:idx_terminal_logs_session"`
	ClientUUID  string    `json:"client_uuid" gorm:"type:varchar(36);index"`
	ClientName  string    `json:"client_name" gorm:"type:varchar(100)"`
	AdminUUID   string    `json:"admin_uuid" gorm:"type:varchar(36);index"`
	AdminIP     string    `json:"admin_ip" gorm:"type:varchar(45)"`
	StartTime   LocalTime `json:"start_time" gorm:"type:timestamp;index"`
	EndTime     *LocalTime `json:"end_time" gorm:"type:timestamp"`
	Duration    string    `json:"duration" gorm:"type:varchar(50)"`
	Status      string    `json:"status" gorm:"type:varchar(20)"` // connected, disconnected, timeout, error
	BytesSent   int64     `json:"bytes_sent" gorm:"default:0"`
	BytesRecv   int64     `json:"bytes_recv" gorm:"default:0"`
}

type TerminalDataLog struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID  string    `json:"session_id" gorm:"type:varchar(64);index:idx_terminal_data_session"`
	Direction  string    `json:"direction" gorm:"type:varchar(10)"` // browser_to_agent, agent_to_browser
	Timestamp  LocalTime `json:"timestamp" gorm:"type:timestamp;index"`
	DataSize   int       `json:"data_size" gorm:"default:0"`
	DataSample string    `json:"data_sample" gorm:"type:text"`
	IsBinary   bool      `json:"is_binary" gorm:"default:false"`
}
