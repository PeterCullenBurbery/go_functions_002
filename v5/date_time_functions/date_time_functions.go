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

	"github.com/PeterCullenBurbery/go_functions_002/v5/system_management_functions"
)

// Format_now returns the current time formatted as "2006-01-02 15:04:05"
func Format_now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Date_time_stamp returns a timestamp string formatted via a temporary Java program.
// It takes no arguments. Java will be installed via Chocolatey if needed.
func Date_time_stamp() (string, error) {
	// Ensure Java is installed
	if err := system_management_functions.Install_Java(); err != nil {
		return "", fmt.Errorf("❌ Java installation failed: %w", err)
	}

	// Try to find java and javac from PATH
	java_cmd, err_java := exec.LookPath("java")
	javac_cmd, err_javac := exec.LookPath("javac")

	// If either is missing, fallback to known Adoptium path
	if err_java != nil || err_javac != nil {
		fallback_base := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin`
		java_fallback := filepath.Join(fallback_base, "java.exe")
		javac_fallback := filepath.Join(fallback_base, "javac.exe")

		if system_management_functions.File_exists(java_fallback) && system_management_functions.File_exists(javac_fallback) {
			java_cmd = java_fallback
			javac_cmd = javac_fallback
		} else {
			return "", fmt.Errorf("❌ Could not locate java or javac in PATH or fallback directory")
		}
	}

	// Create temp directory for Java source and class files
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

	// Trim output
	return strings.TrimSpace(output_buffer.String()), nil
}

// Safe_time_stamp optionally replaces "/" with " slash " if mode == 1.
func Safe_time_stamp(timestamp string, mode int) string {
	if mode == 1 {
		return strings.ReplaceAll(timestamp, "/", " slash ")
	}
	return timestamp
}

// Generate_pdb_name_from_timestamp returns a dynamic PDB name in the format:
// pdb_<YYYY>_<MMM>_<DDD>_<HHH>_<MMM>_<SSS>
//
// Example:
//     pdb_2025_007_031_017_020_008
func Generate_pdb_name_from_timestamp() (string, error) {
	// Ensure Java is installed if needed for time zone detection
	if err := system_management_functions.Install_Java(); err != nil {
		return "", fmt.Errorf("❌ Java installation failed: %w", err)
	}

	// Get the current local time
	now := time.Now()

	// Format each component accordingly
	year := now.Year()
	month := fmt.Sprintf("%03d", int(now.Month()))
	day := fmt.Sprintf("%03d", now.Day())
	hour := fmt.Sprintf("%03d", now.Hour())
	minute := fmt.Sprintf("%03d", now.Minute())
	second := fmt.Sprintf("%03d", now.Second())

	// Assemble and return the PDB name
	return fmt.Sprintf("pdb_%d_%s_%s_%s_%s_%s", year, month, day, hour, minute, second), nil
}