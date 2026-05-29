package record

// Category defines the severity level or type of a log entry.
type Category string

const (
	// MESSAGE represents standard informational log entries.
	MESSAGE Category = "MESSAGE"

	// WARNING represents entries that indicate potential issues or important states.
	WARNING Category = "WARNING"

	// ERROR represents entries for failures that do not stop the application.
	ERROR   Category = "ERROR"
)

// Record represents a single log entry containing its metadata and content.
// It includes tags for JSON and BSON serialization, making it compatible 
// with modern storage engines and web services.
type Record struct {
	// Category is the classification of the log (e.g., MESSAGE, ERROR).
	Category  Category `json:"category" bson:"category"`
	// Message is the actual text content of the log entry.
	Message   string   `json:"message" bson:"message"`
	// Timestamp is the creation time in Unix milliseconds.
	Timestamp int64    `json:"timestamp" bson:"timestamp"`
}
