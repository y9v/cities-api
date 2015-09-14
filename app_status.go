package main

import (
	"github.com/lebedev-yury/cities/ds"
)

type AppStatus struct {
	Status     string         `json:"status"`
	Statistics *ds.Statistics `json:"statistics"`
}

func (appStatus *AppStatus) IsIndexed() bool {
	return appStatus.Status == "ok"
}

func GetAppStatus() *AppStatus {
	var appStatus AppStatus

	appStatus.Statistics = ds.GetStatistics(db)
	if appStatus.Statistics.CitiesCount == 0 {
		appStatus.Status = "indexing"
	} else {
		appStatus.Status = "ok"
	}

	return &appStatus
}
