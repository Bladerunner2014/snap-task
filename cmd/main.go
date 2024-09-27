package main

import (
    "log"
    "net/http"
	"github.com/Bladerunner2014/snap-task/internal/controller"
	"github.com/Bladerunner2014/snap-task/internal/handler/http"
)


func main() {
    ctrl, err := controller.New("responses.db")
    if err != nil {
        log.Fatalf("Error creating SQLite controller: %v", err)
    }
    defer ctrl.Close()

    scheduler := handler.New(10, ctrl)
    
    http.Handle("/job", http.HandlerFunc(scheduler.Handle))

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}