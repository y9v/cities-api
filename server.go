package main

import (
	"net/http"
	"time"
)

func Server() *http.Server {
	return &http.Server{
		Addr:           ":" + configuration.Port,
		Handler:        newRouter(),
		ReadTimeout:    time.Duration(configuration.Timeout) * time.Second,
		WriteTimeout:   time.Duration(configuration.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
