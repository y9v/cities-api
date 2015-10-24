package ds

import (
	"github.com/boltdb/bolt"
)

type AppStatus struct {
	Statistics *Statistics `json:"meta"`
}

func (appStatus *AppStatus) IsIndexed() bool {
	return appStatus.Statistics != nil && appStatus.Statistics.Status == "ok"
}

func GetAppStatus(db *bolt.DB) *AppStatus {
	var appStatus AppStatus
	appStatus.Statistics, _ = GetStatistics(db)

	return &appStatus
}
