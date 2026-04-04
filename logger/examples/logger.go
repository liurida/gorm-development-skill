package examples

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- Custom Logger Implementation ---

// CustomLogger is an example of a custom logger that implements gorm's logger.Interface
type CustomLogger struct {
	// Add any fields your logger might need
	SlowThreshold time.Duration
}

// NewCustomLogger creates a new instance of our custom logger
func NewCustomLogger() logger.Interface {
	return &CustomLogger{
		SlowThreshold: 200 * time.Millisecond, // Default slow threshold
	}
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	// Custom loggers should allow changing log levels
	newLogger := *l
	// Here you would implement logic to handle different levels
	return &newLogger
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[INFO] "+msg, data...)
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[WARN] "+msg, data...)
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[ERROR] "+msg, data...)
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	duration := time.Since(begin)
	sql, rows := fc()

	// Log all queries
	log.Printf("[TRACE] SQL: %s | Rows: %d | Duration: %s", sql, rows, duration)

	// Log slow queries
	if duration > l.SlowThreshold {
		log.Printf("[SLOW] SQL: %s | Duration: %s", sql, duration)
	}

	// Log errors
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FAIL] SQL: %s | Error: %v", sql, err)
	}
}

// --- Using GORM's Default Logger ---

// ConfigureDefaultLogger demonstrates how to configure GORM's built-in logger.
func ConfigureDefaultLogger(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,      // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,             // Don't log ErrRecordNotFound
			ParameterizedQueries:      false,            // Include params in the SQL log
			Colorful:                  true,             // Enable color
		},
	)

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
}

// --- Changing Log Level ---

// ChangeLogLevel shows how to change the log level for a session.
func ChangeLogLevel(db *gorm.DB) {
	// Original session with default level
	db.Find(&[]User{})

	// Change level to Silent for this session
	tx := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
	tx.Find(&[]User{})
}

// --- Debug Mode ---

// DebugSingleOperation demonstrates using Debug() for a single operation.
func DebugSingleOperation(db *gorm.DB) {
	// This operation will be logged at Info level, regardless of global setting
	db.Debug().Where("name = ?", "jinzhu").First(&User{})
}

// --- Using a Custom Logger ---

// UseCustomLogger demonstrates how to initialize GORM with a custom logger.
func UseCustomLogger(dsn string) (*gorm.DB, error) {
	customLogger := NewCustomLogger()

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
}

// --- Contextual Logging ---

// ContextualLogger is an example of a logger that uses context values.
type ContextualLogger struct{
	logger.Interface
}

func (l *ContextualLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if requestID := ctx.Value("request_id"); requestID != nil {
		msg = fmt.Sprintf("[RequestID: %v] %s", requestID, msg)
	}
	log.Printf("[INFO] "+msg, data...)
}

// Trace demonstrates adding context to trace logs.
func (l *ContextualLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	duration := time.Since(begin)

	if requestID := ctx.Value("request_id"); requestID != nil {
		log.Printf("[TRACE][RequestID: %v] SQL: %s | Rows: %d | Duration: %s", requestID, sql, rows, duration)
	} else {
		log.Printf("[TRACE] SQL: %s | Rows: %d | Duration: %s", sql, rows, duration)
	}
}

// RunWithContextualLogger shows how to use a logger that reads from context.
func RunWithContextualLogger(db *gorm.DB) {
	// Create a context with a request ID
	ctx := context.WithValue(context.Background(), "request_id", "user-123-abc")

	// Use the context with the DB operation
	db.WithContext(ctx).First(&User{})
}
