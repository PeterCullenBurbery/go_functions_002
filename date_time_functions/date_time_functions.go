// date_time_functions.go

package date_time_functions

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Format_now returns the current time formatted as "2006-01-02 15:04:05"
func Format_now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Date_time_stamp returns a timestamp string formatted via a temporary Java program.
// It supports optional overrides for javac/java paths.
func Date_time_stamp(args ...string) (string, error) {
	var javac_cmd, java_cmd string

	switch len(args) {
	case 0:
		// Default: look in PATH
		var err error
		javac_cmd, err = exec.LookPath("javac")
		if err != nil {
			return "", fmt.Errorf("❌ 'javac' not found in PATH. Please ensure JDK is installed")
		}
		java_cmd, err = exec.LookPath("java")
		if err != nil {
			return "", fmt.Errorf("❌ 'java' not found in PATH. Please ensure JRE is installed")
		}
	case 2:
		javac_cmd = args[0]
		java_cmd = args[1]
	default:
		return "", fmt.Errorf("❌ Date_time_stamp() expects 0 or 2 arguments (javac_path, java_path)")
	}

	temp_dir, err := os.MkdirTemp("", "date_time_stamp")
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(temp_dir)

	const java_file_name = "date_time_stamp.java"
	const class_name = "date_time_stamp"
	java_file_path := filepath.Join(temp_dir, java_file_name)

	java_code := `import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.temporal.WeekFields;

public class date_time_stamp {
    public static void main(String[] args) {
        ZonedDateTime now = ZonedDateTime.now();
        ZoneId tz = now.getZone();
        String date_part = now.format(DateTimeFormatter.ofPattern("yyyy-0MM-0dd"));
        String time_part = now.format(DateTimeFormatter.ofPattern("0HH.0mm.0ss.nnnnnnnnn"));
        WeekFields wf = WeekFields.ISO;
        int week = now.get(wf.weekOfWeekBasedYear());
        int weekday = now.get(wf.dayOfWeek());
        int iso_year = now.get(wf.weekBasedYear());
        int day_of_year = now.getDayOfYear();
        String output = String.format(
            "%s %s %04d-W%03d-%03d %04d-%03d",
            date_part, time_part, iso_year, week, weekday, now.getYear(), day_of_year
        );
        output = output.replace(time_part, time_part + " " + tz);
        System.out.println(output);
    }
}`

	if err := os.WriteFile(java_file_path, []byte(java_code), 0644); err != nil {
		return "", fmt.Errorf("❌ Failed to write Java file: %w", err)
	}

	// Compile
	cmd_compile := exec.Command(javac_cmd, java_file_name)
	cmd_compile.Dir = temp_dir
	if err := cmd_compile.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to compile Java file: %w", err)
	}

	// Run
	cmd_run := exec.Command(java_cmd, class_name)
	cmd_run.Dir = temp_dir
	var output_buffer bytes.Buffer
	cmd_run.Stdout = &output_buffer
	cmd_run.Stderr = &output_buffer
	if err := cmd_run.Run(); err != nil {
		return "", fmt.Errorf("❌ Failed to run Java class: %w\nOutput:\n%s", err, output_buffer.String())
	}

	// Trim any trailing newline or carriage return
	return strings.TrimSpace(output_buffer.String()), nil
}

// Safe_time_stamp optionally replaces "/" with " slash " if mode == 1.
func Safe_time_stamp(timestamp string, mode int) string {
	if mode == 1 {
		return strings.ReplaceAll(timestamp, "/", " slash ")
	}
	return timestamp
}