package main

import (
    "net/http"
    "github.com/joho/godotenv"
	"github.com/Bladerunner2014/snap-task/internal/controller"
	"github.com/Bladerunner2014/snap-task/internal/handler/http"
    "github.com/Bladerunner2014/snap-task/pkg/log"
    "fmt"
    "os"
    "strconv"


)

var factory = log.LoggerFactory{}
var logger = factory.NewLogger()

func loadEnv() {
    // Attempt to load the .env file. If it fails, proceed with default environment variables.
    err:=godotenv.Load("../.env")
    if err != nil {
        logger.Fatal().Msg("error loading .env file")
      }  
}

func main() {
    loadEnv()
    db_name := os.Getenv("SQLITE_DB_NAME")

    ctrl, err := controller.New(db_name)
    if err != nil {
        logger.Fatal().Msg("error creating sqlite controller")
    }
    defer ctrl.Close()
    maxGoStr := os.Getenv("MAX_GO")
    maxGo, err := strconv.Atoi(maxGoStr)
	if err != nil {
        logger.Fatal().Msg("error loading .env file")
	}
    job_info:=fmt.Sprintf("max job set to :%d", maxGo)
    logger.Info().Msg(job_info)
    scheduler := handler.New(maxGo, ctrl)
    port := os.Getenv("PORT")
    logger.Info().Msg("Server starting on 0.0.0.0:8080")
    http.Handle("/job", http.HandlerFunc(scheduler.Handle))
    if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)}

}