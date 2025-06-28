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
		return fmt.Errorf("❌ Failed to install %s: %w", package_name, err)
	}

	// Verify installation
	verify_cmd := exec.Command(choco_path, "list", "--local-only")
	output, err := verify_cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("⚠️ Could not verify installation of %s: %w", package_name, err)
	}

	if strings.Contains(string(output), package_name) {
		log.Printf("✅ %s installed successfully.", package_name)
	} else {
		log.Printf("⚠️ Install command ran, but %s may not be fully installed.", package_name)
	}

	return nil
}