package system_management_functions

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"net/http"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// Install_choco installs Chocolatey using the official PowerShell script.
// It takes no arguments and logs output to the standard logger.
// You could have Install_Choco check if choco is installed before installing. Then you could just call Install_Choco, and it would handle the details of whether Choco was installed or not.
// Install_choco installs Chocolatey if it is not already installed.
// It logs all steps to the standard logger.
func Install_choco() error {
	if Is_Choco_installed() {
		log.Println("✅ Chocolatey is already installed. Skipping installation.")
		return nil
	}

	log.Println("📦 Chocolatey not found. Starting installation...")

	powershellCommand := `Set-ExecutionPolicy Bypass -Scope Process -Force; ` +
		`[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; ` +
		`iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", powershellCommand)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("❌ Chocolatey installation failed: %w", err)
	}

	// Recheck to confirm installation succeeded
	if !Is_Choco_installed() {
		return fmt.Errorf("❌ Chocolatey installation script ran, but choco.exe was not found afterward")
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

// Choco_install installs the given Chocolatey package and checks if it was installed successfully.
func Choco_install(package_name string) error {
	log.Printf("🚀 Starting installation of %s via Chocolatey...", package_name)

	// Ensure Chocolatey is installed
	if err := Install_choco(); err != nil {
		return fmt.Errorf("❌ Failed to install or locate Chocolatey: %w", err)
	}

	// Resolve choco.exe path
	choco_path, err := exec.LookPath("choco")
	if err != nil {
		choco_path = `C:\ProgramData\chocolatey\bin\choco.exe`
		if _, statErr := os.Stat(choco_path); os.IsNotExist(statErr) {
			return fmt.Errorf("❌ Chocolatey not found at %s even after attempted installation", choco_path)
		}
	}

	// Run installation
	args := []string{"install", package_name, "--yes"}
	install_cmd := exec.Command(choco_path, args...)
	install_cmd.Stdout = log.Writer()
	install_cmd.Stderr = log.Writer()

	if err := install_cmd.Run(); err != nil {
		log.Printf("⚠️ Install command failed or exited with warning: %v", err)
		// Continue to verification anyway
	}

	// Verify installation (via choco list)
	verify_cmd := exec.Command(choco_path, "list", "--limit-output", "--exact", package_name)
	output, _ := verify_cmd.CombinedOutput()
	output_str := strings.TrimSpace(string(output))

	if strings.HasPrefix(strings.ToLower(output_str), strings.ToLower(package_name)+"|") {
		log.Printf("✅ %s installed successfully or already present.", package_name)
		return nil
	}

	return fmt.Errorf("⚠️ Could not verify installation of %s. Raw output:\n%s", package_name, output_str)
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

	java_exists := File_exists(java_fallback)
	javac_exists := File_exists(javac_fallback)

	return java_exists && javac_exists
}

// File_exists checks if a file exists and is not a directory.
func File_exists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Install_Java ensures Java is installed by checking Is_Java_installed().
// If not found, it installs the temurin21 JDK via Chocolatey.
// You could use Is java installed after choco_install(java).
// Install_Java ensures Java is installed by checking Is_Java_installed().
// If not found, it installs the temurin21 JDK via Chocolatey and verifies the result.
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

	// Re-check after installation
	if !Is_Java_installed() {
		return fmt.Errorf("❌ temurin21 JDK was installed, but Java was still not detected")
	}

	log.Println("✅ temurin21 JDK installation complete and verified.")
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

// Download_file downloads a file from the given URL and saves it to the specified destination path.
//
// Parameters:
//   - destination_path: The full file path (including filename) where the downloaded content will be saved.
//   - url: The HTTP or HTTPS URL from which to download the file.
//
// Returns:
//   - An error if the download, file creation, or writing fails; otherwise, nil.
//
// Example:
//   err := Download_file("C:\\downloads\\example.exe", "https://example.com/file.exe")
//   if err != nil {
//       log.Fatalf("Download failed: %v", err)
//   }
func Download_file(destination_path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(destination_path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
// Add_to_path adds the given path to the top of the system PATH (HKLM) if not already present.
// It broadcasts the environment change so apps like Explorer pick it up.
func Add_to_path(path_to_add string) error {
	// Resolve full absolute path
	absPath, err := filepath.Abs(path_to_add)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve absolute path: %w", err)
	}

	// If it's a file, get parent folder
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to stat path: %w", err)
	}
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}
	normalizedPath := strings.TrimRight(absPath, `\`)

	// Open registry key for system environment variables
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open system environment registry key: %w", err)
	}
	defer key.Close()

	// Get current PATH
	currentPath, _, err := key.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("❌ Failed to read PATH: %w", err)
	}

	// Normalize and check if already in PATH
	entries := strings.Split(currentPath, ";")
	for i := range entries {
		entries[i] = strings.TrimRight(entries[i], `\`)
	}
	for _, entry := range entries {
		if strings.EqualFold(entry, normalizedPath) {
			fmt.Println("✅ Path already present in system PATH.")
			return nil
		}
	}

	// Prepend and set new PATH
	newPath := normalizedPath + ";" + currentPath
	if err := key.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("❌ Failed to update PATH in registry: %w", err)
	}
	fmt.Println("✅ Path added to the top of system PATH.")

	// Broadcast environment change
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	procSendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ret, _, _ := procSendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Environment"))),
		uintptr(SMTO_ABORTIFHUNG),
		5000,
		uintptr(0),
	)

	if ret == 0 {
		fmt.Println("⚠️ Environment change broadcast may have failed.")
	} else {
		fmt.Println("📢 Environment update broadcast sent.")
	}

	return nil
}

// Remove_from_path removes the given path from the system PATH if present.
// It normalizes the path, modifies HKLM registry, and broadcasts environment changes.
func Remove_from_path(path_to_remove string) error {
	// Resolve to absolute path
	absPath, err := filepath.Abs(path_to_remove)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve absolute path: %w", err)
	}

	// If it's a file, get parent directory
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to stat path: %w", err)
	}
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}
	normalizedPath := strings.TrimRight(absPath, `\`)

	// Open system environment key
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open system environment registry key: %w", err)
	}
	defer key.Close()

	// Get current PATH
	currentPath, _, err := key.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("❌ Failed to read PATH: %w", err)
	}

	entries := strings.Split(currentPath, ";")
	normalizedEntries := make([]string, 0, len(entries))
	found := false

	for _, entry := range entries {
		trimmed := strings.TrimRight(entry, `\`)
		if strings.EqualFold(trimmed, normalizedPath) {
			found = true
			continue
		}
		normalizedEntries = append(normalizedEntries, entry)
	}

	if !found {
		fmt.Println("ℹ️ Path not found in system PATH.")
		return nil
	}

	newPath := strings.Join(normalizedEntries, ";")
	if err := key.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("❌ Failed to update PATH in registry: %w", err)
	}
	fmt.Printf("✅ Path '%s' removed from system PATH.\n", normalizedPath)

	// Broadcast environment change
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	procSendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ret, _, _ := procSendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Environment"))),
		uintptr(SMTO_ABORTIFHUNG),
		5000,
		uintptr(0),
	)

	if ret == 0 {
		fmt.Println("⚠️ Environment change broadcast may have failed.")
	} else {
		fmt.Println("📢 Environment update broadcast sent.")
	}

	return nil
}