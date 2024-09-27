package log

import (
    "os"
    "strings"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
    "github.com/natefinch/lumberjack"
    "io"
    "time"
)

// LoggerFactory is a struct to hold log settings.
type LoggerFactory struct {
    logFilePath   string
    logLevel      string
    humanReadable bool
}

// loadEnv loads environment variables from the .env file.
func loadEnv() {
    // Attempt to load the .env file. If it fails, proceed with default environment variables.
    godotenv.Load("../.env")
   
}

// NewLogger initializes and returns a zerolog logger based on the settings provided.
// If the caller does not provide configs, default configs will be loaded from environment variables.
func (lf *LoggerFactory) NewLogger() zerolog.Logger {
    // Load the .env file to get environment variables for default config
    loadEnv()

    // If logFilePath is not set, get it from the environment or use a default value.
    if lf.logFilePath == "" {
        lf.logFilePath = os.Getenv("LOG_FILE_PATH")
        if lf.logFilePath == "" {
            lf.logFilePath = "default.log" // Set default log file path if env var is not set
        }
    }

    // If logLevel is not set, get it from the environment or use a default value.
    if lf.logLevel == "" {
        lf.logLevel = os.Getenv("LOG_LEVEL")
        if lf.logLevel == "" {
            lf.logLevel = "info" // Default log level
        }
    }

    // If humanReadable is not set, get it from the environment (true/false) or use a default value.
    if !lf.humanReadable {
        humanReadableEnv := os.Getenv("LOG_HUMAN_READABLE")
        if humanReadableEnv == "" {
            lf.humanReadable = false // Default to false (JSON format)
        } else {
            lf.humanReadable = strings.ToLower(humanReadableEnv) == "true"
        }
    }

    // Open or create the log file with lumberjack for log rotation
    logFile := &lumberjack.Logger{
        Filename:   lf.logFilePath,
        MaxSize:    10, // megabytes
        MaxBackups: 3,  // keep up to 3 old log files
        MaxAge:     28, // days
        Compress:   true, // compress old log files
    }

    // Set up the multi-output writer to log to both file and console
    var writers []io.Writer

    // Add file writer (JSON format by default)
    writers = append(writers, logFile)

    // Add console output (human-readable if configured)
    if lf.humanReadable {
        consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
        writers = append(writers, consoleWriter)
    } else {
        writers = append(writers, os.Stdout) // Add stdout for JSON logs in the console
    }

    // Multi-level writer to write to both file and console
    multi := zerolog.MultiLevelWriter(writers...)

    // Create the logger and set output to multi
    logger := zerolog.New(multi).With().Timestamp().Logger()

    // Set log level from string
    switch strings.ToLower(lf.logLevel) {
    case "debug":
        logger = logger.Level(zerolog.DebugLevel)
    case "info":
        logger = logger.Level(zerolog.InfoLevel)
    case "warn":
        logger = logger.Level(zerolog.WarnLevel)
    case "error":
        logger = logger.Level(zerolog.ErrorLevel)
    default:
        logger = logger.Level(zerolog.InfoLevel) // Default to info if level not recognized
    }

    return logger
}
