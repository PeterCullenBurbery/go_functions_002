package system_management_functions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Choco_install installs the given Chocolatey package and checks if it was installed successfully.
func Choco_install(package_name string) error {
	log.Printf("🚀 Starting installation of %s via Chocolatey...", package_name)

	// Check if Chocolatey is installed
	if !Is_Choco_installed() {
		log.Println("🔍 Chocolatey not found. Attempting to install it...")
		if err := Install_choco(); err != nil {
			return fmt.Errorf("❌ Failed to install Chocolatey: %w", err)
		}
		log.Println("✅ Chocolatey installation complete. Proceeding with package installation...")
	}

	// Try to resolve choco.exe
	choco_path, err := exec.LookPath("choco")
	if err != nil {
		choco_path = `C:\ProgramData\chocolatey\bin\choco.exe`
		if _, statErr := os.Stat(choco_path); os.IsNotExist(statErr) {
			return fmt.Errorf("❌ Chocolatey not found at %s even after attempted installation", choco_path)
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

	// Verify installation (new method: --limit-output --exact <name>)
	verify_cmd := exec.Command(choco_path, "list", "--limit-output", "--exact", package_name)
	output, _ := verify_cmd.CombinedOutput()
	output_str := strings.TrimSpace(string(output))

	if strings.HasPrefix(strings.ToLower(output_str), strings.ToLower(package_name)+"|") {
		log.Printf("✅ %s installed successfully or already present.", package_name)
		return nil
	}

	return fmt.Errorf("⚠️ Could not verify installation of %s. Raw output:\n%s", package_name, output_str)
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

// Install_choco installs Chocolatey using the official PowerShell script.
// It takes no arguments and logs output to the standard logger.
func Install_choco() error {
	log.Println("📦 Starting Chocolatey installation...")

	powershellCommand := `Set-ExecutionPolicy Bypass -Scope Process -Force; ` +
		`[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; ` +
		`iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", powershellCommand)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("❌ Chocolatey installation failed: %w", err)
	}

	log.Println("✅ Chocolatey installed successfully.")
	return nil
}

// Is_Choco_installed checks if Chocolatey is installed.
// It returns true if choco.exe is found in PATH or at the default location.
func Is_Choco_installed() bool {
	// First try to resolve choco.exe from PATH
	if _, err := exec.LookPath("choco"); err == nil {
		return true
	}

	// Fallback to default Chocolatey path
	default_choco_path := `C:\ProgramData\chocolatey\bin\choco.exe`
	if _, err := os.Stat(default_choco_path); err == nil {
		return true
	}

	return false
}

// Is_Java_installed checks if both java.exe and javac.exe are available in PATH,
// or in the default Eclipse Adoptium installation directory.
func Is_Java_installed() bool {
	// Check PATH using exec.LookPath
	if _, err := exec.LookPath("java"); err == nil {
		if _, err := exec.LookPath("javac"); err == nil {
			return true
		}
	}

	// Fallback: Check default Eclipse Adoptium JDK location
	base_path := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot\bin`
	java_fallback := filepath.Join(base_path, "java.exe")
	javac_fallback := filepath.Join(base_path, "javac.exe")

	java_exists := fileExists(java_fallback)
	javac_exists := fileExists(javac_fallback)

	return java_exists && javac_exists
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Install_Java ensures Java is installed by checking Is_Java_installed().
// If not found, it installs the temurin21 JDK via Chocolatey.
func Install_Java() error {
	log.Println("📦 Checking if Java is already installed...")

	if Is_Java_installed() {
		log.Println("✅ Java is already installed. Skipping installation.")
		return nil
	}

	log.Println("❌ Java not found. Proceeding with installation via Chocolatey...")

	if err := Choco_install("temurin21"); err != nil {
		return fmt.Errorf("❌ Failed to install temurin21 JDK: %w", err)
	}

	log.Println("✅ temurin21 JDK installation complete.")
	return nil
}