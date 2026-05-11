package services





type StoreHealthRecord struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	TotalConnections int32 `json:"total_connections,omitempty"`
	IdleConnections int32 `json:"idle_connections,omitempty"`
	AcquiredConnections int32 `json:"acquired_connections,omitempty"`
	Error string `json:"error,omitempty"`
}

