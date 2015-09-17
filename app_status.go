package main

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/ds"
)

type AppStatus struct {
	Status     string         `json:"status"`
	Statistics *ds.Statistics `json:"statistics"`
}

func (appStatus *AppStatus) IsIndexed() bool {
	return appStatus.Status == "ok"
}

func GetAppStatus(db *bolt.DB) *AppStatus {
	var err error
	var appStatus AppStatus

	appStatus.Statistics, err = ds.GetStatistics(db)
	if err == nil && appStatus.Statistics.CityNamesCount > 0 {
		appStatus.Status = "ok"
	} else {
		appStatus.Status = "indexing"
	}

	return &appStatus
}
