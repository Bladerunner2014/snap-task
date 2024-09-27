package handler

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "reflect"
    "time"

    "github.com/Bladerunner2014/snap-task/internal/controller"
    "github.com/Bladerunner2014/snap-task/pkg/model"
)

type JobScheduler struct {
    Jobs    []model.Job
    MaxJobs int
    Ctrl    *controller.SQLiteController
}

func New(maxJobs int, ctrl *controller.SQLiteController) *JobScheduler {
    return &JobScheduler{
        Jobs:    make([]model.Job, 0),
        MaxJobs: maxJobs,
        Ctrl:    ctrl,
    }
}

func (s *JobScheduler) AddJob(job model.Job) error {
    if len(s.Jobs) >= s.MaxJobs {
        return fmt.Errorf("maximum number of jobs (%d) reached", s.MaxJobs)
    }
    s.Jobs = append(s.Jobs, job)
    go s.runJob(job)
    return nil
}

func (s *JobScheduler) runJob(job model.Job) {
    ticker := time.NewTicker(job.Interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            resp, err := http.Get(job.URL)
            if err != nil {
                log.Printf("Error fetching %s: %v", job.URL, err)
                continue
            }
            defer resp.Body.Close()

            var responseData map[string]interface{}
            if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
                log.Printf("Error parsing JSON from %s: %v", job.URL, err)
                continue
            }

            if matchesPattern(responseData, job.Pattern) {
                s.Ctrl.StoreResponse(job.URL, fmt.Sprintf("%v", responseData))
            }
        case <-job.StopChan:
            log.Printf("Job for URL %s stopped manually", job.URL)
            return
        }
    }
}

func matchesPattern(data, pattern map[string]interface{}) bool {
    for key, expectedType := range pattern {
        value, exists := data[key]
        if !exists {
            return false
        }

        actualType := reflect.TypeOf(value).String()
        if actualType != expectedType.(string) {
            return false
        }
    }
    return true
}

func (s *JobScheduler) Handle(w http.ResponseWriter, r *http.Request) {
    // Set the content type to JSON
    w.Header().Set("Content-Type", "application/json")

    if r.Method != http.MethodPost {
        jsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
        return
    }

    var jobReq model.JobRequest
    if err := json.NewDecoder(r.Body).Decode(&jobReq); err != nil {
        jsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
        return
    }

    job := model.Job{
        URL:      jobReq.URL,
        Pattern:  jobReq.Pattern,
        Interval: time.Duration(jobReq.Interval) * time.Second,
        Duration: time.Duration(jobReq.Duration) * time.Second,
        StopChan: make(chan struct{}),
    }

    if err := s.AddJob(job); err != nil {
        jsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
        return
    }

    jsonResponse(w, map[string]string{"message": "Job added successfully"}, http.StatusCreated)
}

// Helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}