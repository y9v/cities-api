package main

import (
	"github.com/boltdb/bolt"
	"github.com/lebedev-yury/cities/cache"
	"github.com/lebedev-yury/cities/config"
	"net/http"
	"time"
)

func Server(db *bolt.DB, options *config.Options, c *cache.Cache) *http.Server {
	return &http.Server{
		Addr:           ":" + options.Port,
		Handler:        newRouter(db, options, c),
		ReadTimeout:    time.Duration(options.Timeout) * time.Second,
		WriteTimeout:   time.Duration(options.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
