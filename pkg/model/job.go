package model

import (
    "time"
)

type Job struct {
    URL      string
    Pattern  map[string]interface{}
    Interval time.Duration
    Duration time.Duration
    StopChan chan struct{}
}

type JobRequest struct {
    URL      string                 `json:"url"`
    Pattern  map[string]interface{} `json:"pattern"`
    Interval int                    `json:"interval"`
    Duration int                    `json:"duration"`
}