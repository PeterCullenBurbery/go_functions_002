package system_management_functions

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Choco_install installs the given Chocolatey package and checks if it was installed successfully.
func Choco_install(package_name string) error {
	log.Printf("🚀 Starting installation of %s via Chocolatey...", package_name)

	// Try to resolve choco.exe
	choco_path, err := exec.LookPath("choco")
	if err != nil {
		// Fallback to default path
		choco_path = `C:\ProgramData\chocolatey\bin\choco.exe`
		check_cmd := exec.Command("cmd", "/c", "if exist \""+choco_path+"\" (exit 0) else (exit 1)")
		if err := check_cmd.Run(); err != nil {
			return fmt.Errorf("❌ Chocolatey not found. Please install Chocolatey first")
		}
	}

	// Install the package
	args := []string{"install", package_name, "--yes"}
	install_cmd := exec.Command(choco_path, args...)
	install_cmd.Stdout = log.Writer()
	install_cmd.Stderr = log.Writer()

	if err := install_cmd.Run(); err != nil {
		log.Printf("⚠️ Install command failed or exited with warning: %v", err)
		// Continue to verification anyway
	}

	// Verify installation (specific package)
	verify_cmd := exec.Command(choco_path, "list", "--local-only", package_name)
	output, err := verify_cmd.CombinedOutput()
	output_str := string(output)

	if err != nil && !strings.Contains(output_str, package_name) {
		return fmt.Errorf("⚠️ Could not verify installation of %s: %w", package_name, err)
	}

	if strings.Contains(output_str, package_name) {
		log.Printf("✅ %s installed successfully or already present.", package_name)
	} else {
		log.Printf("⚠️ Install command ran, but %s may not be fully installed.", package_name)
	}

	return nil
}

// Winget_install installs the specified package using winget with standard flags.
// Example: Winget_install("Visual Studio Code", "Microsoft.VisualStudioCode")
func Winget_install(package_name string, package_id string) error {
	log.Printf("🚀 Starting installation of %s via winget...", package_name)

	args := []string{
		"install",
		"-e",
		"--id", package_id,
		"--scope", "machine",
		"--silent",
		"--accept-package-agreements",
		"--accept-source-agreements",
	}

	cmd := exec.Command("winget", args...)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("❌ Failed to install %s via winget: %w", package_name, err)
	}

	log.Printf("✅ %s installed successfully via winget.", package_name)
	return nil
}