package handler

import (
    // "bytes"
    // "encoding/json"
    // "net/http"
    // "net/http/httptest"
    "testing"

    "github.com/Bladerunner2014/snap-task/internal/controller"
    "github.com/Bladerunner2014/snap-task/pkg/model"
)

func TestNew(t *testing.T) {
    ctrl := &controller.SQLiteController{}
    scheduler := New(10, ctrl)

    if scheduler.MaxJobs != 10 {
        t.Errorf("Expected MaxJobs to be 10, got %d", scheduler.MaxJobs)
    }

    if scheduler.Ctrl != ctrl {
        t.Error("Expected Ctrl to be set correctly")
    }

    if len(scheduler.Jobs) != 0 {
        t.Errorf("Expected Jobs to be empty, got %d jobs", len(scheduler.Jobs))
    }
}

func TestAddJob(t *testing.T) {
    scheduler := New(2, &controller.SQLiteController{})

    job1 := model.Job{URL: "http://example.com"}
    err := scheduler.AddJob(job1)
    if err != nil {
        t.Errorf("Unexpected error adding job: %v", err)
    }

    job2 := model.Job{URL: "http://example.org"}
    err = scheduler.AddJob(job2)
    if err != nil {
        t.Errorf("Unexpected error adding job: %v", err)
    }

    job3 := model.Job{URL: "http://example.net"}
    err = scheduler.AddJob(job3)
    if err == nil {
        t.Error("Expected error when adding job above MaxJobs, got nil")
    }
}


