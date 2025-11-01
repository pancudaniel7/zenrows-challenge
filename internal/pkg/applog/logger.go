package applog

// AppLogger describes the logging contract used throughout the application.
type AppLogger interface {
	// Info logs informational events.
	Info(msg string, args ...any)
	// Warn logs events that might need attention but are not errors.
	Warn(msg string, args ...any)
	// Error logs failures that prevented successful processing.
	Error(msg string, args ...any)
	// Debug logs verbose diagnostic details for developers.
	Debug(msg string, args ...any)
	// Trace logs highly granular events, typically for deep troubleshooting.
	Trace(msg string, args ...any)
	// Fatal logs a critical error and terminates the process.
	Fatal(msg string, args ...any)
}
