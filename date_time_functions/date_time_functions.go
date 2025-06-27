// date_time_functions.go
package date_time_functions

import "time"

// FormatNow returns the current time formatted as "2006-01-02 15:04:05"
func FormatNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}