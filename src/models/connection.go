package models

import "time"

type Connection struct {
	Username string
	Password string
	Host     string
	Port     int
	Timeout  time.Duration
}
