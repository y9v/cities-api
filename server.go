package main

import (
	"net/http"
	"time"
)

func Server() *http.Server {
	return &http.Server{
		Addr:           ":" + options.Port,
		Handler:        newRouter(),
		ReadTimeout:    time.Duration(options.Timeout) * time.Second,
		WriteTimeout:   time.Duration(options.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
