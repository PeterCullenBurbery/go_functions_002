// system_management_functions.go

package system_management_functions

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"net/http"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
	yekazip "github.com/yeka/zip"
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
// If not found, it installs the temurin21 JDK via Chocolatey,
// verifies the result, and sets JAVA_HOME to the default installation path using PowerShell.
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

	// Set JAVA_HOME to standard Temurin path
	java_home := `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot`
	log.Printf("🔧 Setting JAVA_HOME to: %s", java_home)

	if err := Set_system_environment_variable("JAVA_HOME", java_home); err != nil {
		return fmt.Errorf("❌ Failed to set JAVA_HOME: %w", err)
	}

	log.Println("✅ JAVA_HOME environment variable set successfully.")
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

// Expand_windows_env expands environment variables using the Windows API.
// For example, %SystemRoot% becomes C:\Windows.
func Expand_windows_env(input string) string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procExpand := kernel32.NewProc("ExpandEnvironmentStringsW")

	inputPtr, _ := syscall.UTF16PtrFromString(input)
	buf := make([]uint16, 32767) // MAX_PATH

	ret, _, _ := procExpand.Call(
		uintptr(unsafe.Pointer(inputPtr)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)

	if ret == 0 {
		return input // fallback if expansion fails
	}

	return syscall.UTF16ToString(buf[:ret])
}

// Add_to_path adds the given path to the top of the system PATH (HKLM) if not already present.
// It expands environment variables, removes redundant entries (like %SystemRoot%), avoids duplicates,
// and broadcasts the environment change to Explorer. It also prints PowerShell instructions to refresh the session.
func Add_to_path(path_to_add string) error {
	fmt.Printf("🔧 Input path: %s\n", path_to_add)

	// Step 1: Resolve absolute path
	abs_path, err := filepath.Abs(path_to_add)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve absolute path: %w", err)
	}
	fmt.Printf("📁 Absolute path: %s\n", abs_path)

	info, err := os.Stat(abs_path)
	if err != nil {
		return fmt.Errorf("❌ Failed to stat path: %w", err)
	}
	if !info.IsDir() {
		abs_path = filepath.Dir(abs_path)
	}
	normalized := strings.TrimRight(abs_path, `\`)
	fmt.Printf("🧹 Normalized path: %s\n", normalized)

	// Step 2: Open system PATH from registry
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()
	fmt.Println("🔑 Opened HKLM system environment registry key.")

	raw_path, _, err := key.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("❌ Failed to read PATH: %w", err)
	}
	fmt.Println("📍 Current PATH (raw):")
	fmt.Println(raw_path)

	// Step 3: Process PATH entries
	entries := strings.Split(raw_path, ";")
	fmt.Println("🔍 Checking each existing PATH entry against target:")

	normalized_lower := strings.ToLower(normalized)
	already_exists := false
	seen := make(map[string]bool)
	rebuilt := []string{normalized} // New path goes first
	seen[normalized_lower] = true

	for _, entry := range entries {
		entry_trimmed := strings.TrimSpace(strings.TrimRight(entry, `\`))
		if entry_trimmed == "" {
			continue
		}

		expanded := strings.TrimRight(Expand_windows_env(entry_trimmed), `\`)
		lower_expanded := strings.ToLower(expanded)

		if !strings.EqualFold(entry_trimmed, expanded) {
			fmt.Printf("   - Original: %-70s →  Expanded: %s\n", entry_trimmed, expanded)
		}

		if lower_expanded == normalized_lower {
			already_exists = true
		}

		if !seen[lower_expanded] {
			rebuilt = append(rebuilt, expanded)
			seen[lower_expanded] = true
		}
	}

	if already_exists {
		fmt.Println("✅ Path already present in system PATH (via expanded match).")
		return nil
	}

	new_path := strings.Join(rebuilt, ";")
	fmt.Println("🧩 New PATH to set in registry:")
	fmt.Println(new_path)

	// Step 4: Write back to registry
	if err := key.SetStringValue("Path", new_path); err != nil {
		return fmt.Errorf("❌ Failed to update PATH in registry: %w", err)
	}
	fmt.Println("✅ Path added to the top of system PATH.")

	// Step 5: Broadcast change
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

	// Step 6: Check for refreshenv and print accordingly
	if _, err := exec.LookPath("refreshenv"); err == nil {
		fmt.Println("♻️  'refreshenv' is available. To update this session, run:")
		fmt.Println("    refreshenv")
	} else {
		fmt.Println("ℹ️  'refreshenv' not available in this session.")
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

// Create_desktop_shortcut creates a .lnk shortcut on the desktop.
// It accepts the target path, shortcut name (optional), description (optional),
// window style (3 = maximized), and allUsers flag.
func Create_desktop_shortcut(target_path, shortcut_name, description string, window_style int, all_users bool) error {
	// Ensure target exists
	if _, err := os.Stat(target_path); os.IsNotExist(err) {
		return fmt.Errorf("❌ Target path does not exist: %s", target_path)
	}

	// Determine desktop path
	var desktopPath string
	if all_users {
		public := os.Getenv("PUBLIC")
		desktopPath = filepath.Join(public, "Desktop")
	} else {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("❌ Could not determine current user: %w", err)
		}
		desktopPath = filepath.Join(usr.HomeDir, "Desktop")
	}

	// Determine shortcut name
	if shortcut_name == "" {
		base := filepath.Base(target_path)
		shortcut_name = strings.TrimSuffix(base, ".exe") + ".lnk"
	}

	shortcutPath := filepath.Join(desktopPath, shortcut_name)

	// Initialize COM
	if err := ole.CoInitialize(0); err != nil {
		return fmt.Errorf("❌ Failed to initialize COM: %w", err)
	}
	defer ole.CoUninitialize()

	// Create Shell COM object
	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("❌ Failed to create WScript.Shell COM object: %w", err)
	}
	defer shell.Release()

	dispatch, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("❌ Failed to get IDispatch: %w", err)
	}
	defer dispatch.Release()

	// Create the shortcut
	shortcutRaw, err := oleutil.CallMethod(dispatch, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to create shortcut: %w", err)
	}
	shortcut := shortcutRaw.ToIDispatch()
	defer shortcut.Release()

	// Set properties
	_, _ = oleutil.PutProperty(shortcut, "TargetPath", target_path)
	_, _ = oleutil.PutProperty(shortcut, "WorkingDirectory", filepath.Dir(target_path))
	_, _ = oleutil.PutProperty(shortcut, "WindowStyle", window_style)
	_, _ = oleutil.PutProperty(shortcut, "Description", description)
	_, _ = oleutil.PutProperty(shortcut, "IconLocation", fmt.Sprintf("%s, 0", target_path))

	// Save
	_, err = oleutil.CallMethod(shortcut, "Save")
	if err != nil {
		return fmt.Errorf("❌ Failed to save shortcut: %w", err)
	}

	fmt.Printf("✅ Shortcut created at: %s\n", shortcutPath)
	return nil
}

// Extract_zip extracts a ZIP archive specified by src into the destination directory dest.
//
// It protects against Zip Slip attacks by ensuring all extracted paths are within dest.
//
// Parameters:
//   - src:  Full path to the ZIP archive.
//   - dest: Destination directory where the contents will be extracted.
//
// Returns:
//   - An error if extraction fails, or nil on success.
func Extract_zip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("❌ failed to open zip archive: %w", err)
	}
	defer r.Close()

	for _, file := range r.File {
		fpath := filepath.Join(dest, file.Name)

		// Zip Slip protection
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("❌ illegal file path in archive: %s", fpath)
		}

		// Directory entry
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, file.Mode()); err != nil {
				return fmt.Errorf("❌ failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("❌ failed to create parent directory for %s: %w", fpath, err)
		}

		// Open destination file
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("❌ failed to create file %s: %w", fpath, err)
		}

		// Open zip file entry
		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("❌ failed to open zip entry %s: %w", file.Name, err)
		}

		// Copy contents
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("❌ failed to copy data to file %s: %w", fpath, err)
		}
	}

	return nil
}

// Exclude_from_Microsoft_Windows_Defender attempts to exclude the specified file or folder
// from Microsoft Defender's real-time protection.
//
// This function first checks whether the Microsoft Defender Antivirus service (WinDefend)
// is currently running. If it is not running (e.g., disabled by policy or replaced by another
// antivirus solution), the function skips the exclusion step without error.
//
// If the provided path refers to a file, its parent directory is excluded instead.
// If the path refers to a directory, it is used directly.
//
// This operation requires administrative privileges if Microsoft Defender is enabled.
//
// Parameters:
//   - path_to_exclude: The absolute or relative path to a file or folder to exclude.
//
// Returns:
//   - nil if the exclusion was successful, unnecessary (because Defender is not running),
//     or if the path was already excluded.
//   - An error if any part of the exclusion process fails (e.g., bad path, PowerShell failure).
//
// Example:
//
//      err := Exclude_from_Microsoft_Windows_Defender("C:\\downloads\\nirsoft")
//      if err != nil {
//          log.Fatalf("Failed to exclude: %v", err)
//      }
func Exclude_from_Microsoft_Windows_Defender(path_to_exclude string) error {
        // Step 0: Check if Microsoft Defender is running
        check_cmd := exec.Command("powershell", "-NoProfile", "-Command",
                `(Get-Service WinDefend).Status`)
        output_bytes, err := check_cmd.Output()
        if err != nil {
                log.Println("ℹ️ Unable to query WinDefend service; skipping exclusion step.")
                return nil
        }
        output := string(output_bytes)
        if output != "Running\r\n" && output != "Running\n" {
                log.Println("ℹ️ Microsoft Defender is not running; skipping exclusion step.")
                return nil
        }

        // Resolve absolute path
        absolute_path, err := filepath.Abs(path_to_exclude)
        if err != nil {
                return fmt.Errorf("❌ Failed to resolve absolute path: %w", err)
        }

        // Stat to determine if it's a file or folder
        file_info, err := os.Stat(absolute_path)
        if err != nil {
                return fmt.Errorf("❌ Failed to stat path: %w", err)
        }

        // If it's a file, get parent directory
        if !file_info.IsDir() {
                absolute_path = filepath.Dir(absolute_path)
        }

        // Normalize (trim trailing slash)
        normalized_path := filepath.Clean(absolute_path)

        // Build PowerShell command to exclude from Defender
        exclude_cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
                fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, normalized_path))

        exclude_output_bytes, err := exclude_cmd.CombinedOutput()
        if err != nil {
                return fmt.Errorf("❌ Failed to exclude from Defender: %w\nOutput: %s", err, string(exclude_output_bytes))
        }

        fmt.Printf("✅ Excluded from Microsoft Defender: %s\n", normalized_path)
        return nil
}

// Extract_password_protected_zip extracts a password-protected ZIP archive using AES or ZipCrypto.
//
// Parameters:
//   - src: full path to the .zip archive
//   - dest: directory where the files should be extracted
//   - password: the password used to decrypt the archive
//
// Returns:
//   - error if extraction fails, otherwise nil
func Extract_password_protected_zip(src, dest, password string) error {
	reader, err := yekazip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("❌ failed to open zip archive: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		fpath := filepath.Join(dest, file.Name)

		// Zip Slip protection
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("❌ illegal file path in archive: %s", fpath)
		}

		// Directory
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, file.Mode()); err != nil {
				return fmt.Errorf("❌ failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("❌ failed to create parent directory: %w", err)
		}

		// Set password and open the file
		file.SetPassword(password)
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("❌ failed to open encrypted file %s: %w", file.Name, err)
		}
		defer rc.Close()

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("❌ failed to create file %s: %w", fpath, err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return fmt.Errorf("❌ failed to write file %s: %w", fpath, err)
		}
	}

	return nil
}

// Set_dark_mode sets Windows to dark mode for both system and apps.
// If restartExplorer is true, it restarts Explorer to apply the change.
func Set_dark_mode(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("AppsUseLightTheme", 0); err != nil {
		return fmt.Errorf("❌ Failed to set AppsUseLightTheme: %w", err)
	}
	if err := key.SetDWordValue("SystemUsesLightTheme", 0); err != nil {
		return fmt.Errorf("❌ Failed to set SystemUsesLightTheme: %w", err)
	}

	fmt.Println("✅ Dark mode set: AppsUseLightTheme & SystemUsesLightTheme = 0")

	if restartExplorer {
		cmd := exec.Command("taskkill", "/f", "/im", "explorer.exe")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		cmd = exec.Command("explorer.exe")
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("❌ Failed to launch Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply Dark Mode.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Set_light_mode sets Windows to light mode for both system and apps.
// If restartExplorer is true, it restarts Explorer to apply the change.
func Set_light_mode(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("AppsUseLightTheme", 1); err != nil {
		return fmt.Errorf("❌ Failed to set AppsUseLightTheme: %w", err)
	}
	if err := key.SetDWordValue("SystemUsesLightTheme", 1); err != nil {
		return fmt.Errorf("❌ Failed to set SystemUsesLightTheme: %w", err)
	}

	fmt.Println("✅ Light mode set: AppsUseLightTheme & SystemUsesLightTheme = 1")

	if restartExplorer {
		cmd := exec.Command("taskkill", "/f", "/im", "explorer.exe")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		cmd = exec.Command("explorer.exe")
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("❌ Failed to launch Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply Light Mode.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Set_start_menu_to_left sets the Windows 11 taskbar alignment to the left
// by writing TaskbarAl=0 in the registry and restarting Explorer.
func Set_start_menu_to_left() error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("TaskbarAl", 0); err != nil {
		return fmt.Errorf("❌ Failed to set TaskbarAl to 0: %w", err)
	}

	fmt.Println("✅ Registry updated: TaskbarAl = 0 (left)")

	// Restart Explorer
	if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
		return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
	}
	if err := exec.Command("explorer.exe").Start(); err != nil {
		return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
	}

	fmt.Println("🔁 Explorer restarted to apply Start menu alignment (left)")
	return nil
}

// Set_start_menu_to_center sets the Windows 11 taskbar alignment to the center
// by writing TaskbarAl=1 in the registry and restarting Explorer.
func Set_start_menu_to_center() error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("TaskbarAl", 1); err != nil {
		return fmt.Errorf("❌ Failed to set TaskbarAl to 1: %w", err)
	}

	fmt.Println("✅ Registry updated: TaskbarAl = 1 (center)")

	// Restart Explorer
	if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
		return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
	}
	if err := exec.Command("explorer.exe").Start(); err != nil {
		return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
	}

	fmt.Println("🔁 Explorer restarted to apply Start menu alignment (center)")
	return nil
}

// Show_file_extensions sets HideFileExt = 0 to make file extensions visible.
// If restartExplorer is true, Explorer is restarted to apply the change immediately.
func Show_file_extensions(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("HideFileExt", 0); err != nil {
		return fmt.Errorf("❌ Failed to set HideFileExt = 0: %w", err)
	}

	fmt.Println("✅ File extensions will be visible (HideFileExt = 0)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply visibility of file extensions.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Do_not_show_file_extensions sets HideFileExt = 1 to hide file extensions.
// If restartExplorer is true, Explorer is restarted to apply the change immediately.
func Do_not_show_file_extensions(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("HideFileExt", 1); err != nil {
		return fmt.Errorf("❌ Failed to set HideFileExt = 1: %w", err)
	}

	fmt.Println("✅ File extensions will be hidden (HideFileExt = 1)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply hiding of file extensions.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Show_hidden_files sets Hidden = 1 to show hidden files in File Explorer.
// If restartExplorer is true, Explorer will be restarted to apply the change immediately.
func Show_hidden_files(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("Hidden", 1); err != nil {
		return fmt.Errorf("❌ Failed to set Hidden = 1: %w", err)
	}

	fmt.Println("✅ Hidden files will be shown (Hidden = 1)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply hidden file visibility.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Do_not_show_hidden_files sets Hidden = 2 to hide hidden files in File Explorer.
// If restartExplorer is true, Explorer will be restarted to apply the change immediately.
func Do_not_show_hidden_files(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("Hidden", 2); err != nil {
		return fmt.Errorf("❌ Failed to set Hidden = 2: %w", err)
	}

	fmt.Println("✅ Hidden files will be hidden (Hidden = 2)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply hiding of hidden files.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Hide_search_box sets SearchboxTaskbarMode = 0 to hide the taskbar search box.
// If restartExplorer is true, Explorer will be restarted to apply the change.
func Hide_search_box(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Search`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("SearchboxTaskbarMode", 0); err != nil {
		return fmt.Errorf("❌ Failed to set SearchboxTaskbarMode = 0: %w", err)
	}

	fmt.Println("✅ Search box will be hidden (SearchboxTaskbarMode = 0)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply hiding of search box.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Do_not_hide_search_box sets SearchboxTaskbarMode = 2 to show the full search box on the taskbar.
// If restartExplorer is true, Explorer will be restarted to apply the change.
func Do_not_hide_search_box(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Search`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("SearchboxTaskbarMode", 2); err != nil {
		return fmt.Errorf("❌ Failed to set SearchboxTaskbarMode = 2: %w", err)
	}

	fmt.Println("✅ Search box will be shown (SearchboxTaskbarMode = 2)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply showing of search box.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Seconds_in_taskbar enables seconds on the taskbar clock by setting ShowSecondsInSystemClock = 1.
// If restartExplorer is true, Explorer will be restarted to apply the change.
func Seconds_in_taskbar(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("ShowSecondsInSystemClock", 1); err != nil {
		return fmt.Errorf("❌ Failed to set ShowSecondsInSystemClock = 1: %w", err)
	}

	fmt.Println("✅ Taskbar clock will display seconds (ShowSecondsInSystemClock = 1)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply the seconds display.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Take_seconds_out_of_taskbar disables seconds on the taskbar clock by setting ShowSecondsInSystemClock = 0.
// If restartExplorer is true, Explorer will be restarted to apply the change.
func Take_seconds_out_of_taskbar(restartExplorer bool) error {
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open/create registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue("ShowSecondsInSystemClock", 0); err != nil {
		return fmt.Errorf("❌ Failed to set ShowSecondsInSystemClock = 0: %w", err)
	}

	fmt.Println("✅ Taskbar clock will hide seconds (ShowSecondsInSystemClock = 0)")

	if restartExplorer {
		if err := exec.Command("taskkill", "/f", "/im", "explorer.exe").Run(); err != nil {
			return fmt.Errorf("❌ Failed to stop Explorer: %w", err)
		}
		if err := exec.Command("explorer.exe").Start(); err != nil {
			return fmt.Errorf("❌ Failed to restart Explorer: %w", err)
		}
		fmt.Println("🔁 Explorer restarted to apply the change.")
	} else {
		fmt.Println("ℹ️ Explorer restart skipped.")
	}

	return nil
}

// Set_short_date_pattern sets the short date pattern to "yyyy-MM-dd-dddd"
// and broadcasts the change to the system.
func Set_short_date_pattern() error {
	const (
		keyPath      = `Control Panel\International`
		valueName    = "sShortDate"
		newPattern   = "yyyy-MM-dd-dddd"
		broadcastMsg = "Intl"
	)

	// Write to registry
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(valueName, newPattern); err != nil {
		return fmt.Errorf("❌ Failed to set short date pattern: %w", err)
	}

	fmt.Printf("✅ Short date pattern set to '%s'.\n", newPattern)

	// Inline SendMessageTimeoutW
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	categoryPtr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(categoryPtr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("📢 System broadcast sent to apply the setting.")
	return nil
}

// Reset_short_date_pattern resets the short date pattern to "M/d/yyyy"
// and broadcasts the change to the system.
func Reset_short_date_pattern() error {
	const (
		keyPath       = `Control Panel\International`
		valueName     = "sShortDate"
		defaultFormat = "M/d/yyyy"
		broadcastMsg  = "Intl"
	)

	// Write to registry
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(valueName, defaultFormat); err != nil {
		return fmt.Errorf("❌ Failed to reset short date pattern: %w", err)
	}

	fmt.Printf("✅ Short date pattern reset to '%s'.\n", defaultFormat)

	// Inline SendMessageTimeoutW
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	categoryPtr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(categoryPtr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("📢 System broadcast sent to apply the setting.")
	return nil
}

// Set_long_date_pattern sets the long date pattern to "yyyy-MM-dd-dddd"
// and broadcasts the change to the system.
func Set_long_date_pattern() error {
	const (
		keyPath      = `Control Panel\International`
		valueName    = "sLongDate"
		newPattern   = "yyyy-MM-dd-dddd"
		broadcastMsg = "Intl"
	)

	// Write to registry
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(valueName, newPattern); err != nil {
		return fmt.Errorf("❌ Failed to set long date pattern: %w", err)
	}

	fmt.Printf("✅ Long date pattern set to '%s'.\n", newPattern)

	// Broadcast setting change (inline SendMessageTimeout)
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	categoryPtr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(categoryPtr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("📢 System broadcast sent to apply the setting.")
	return nil
}

// Reset_long_date_pattern resets the long date pattern to the default "dddd, MMMM d, yyyy"
// and broadcasts the change to the system.
func Reset_long_date_pattern() error {
	const (
		keyPath       = `Control Panel\International`
		valueName     = "sLongDate"
		defaultFormat = "dddd, MMMM d, yyyy"
		broadcastMsg  = "Intl"
	)

	// Write to registry
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(valueName, defaultFormat); err != nil {
		return fmt.Errorf("❌ Failed to reset long date pattern: %w", err)
	}

	fmt.Printf("✅ Long date pattern reset to '%s'.\n", defaultFormat)

	// Broadcast setting change (inline)
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	categoryPtr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(categoryPtr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("📢 System broadcast sent to apply the setting.")
	return nil
}

// Set_time_pattern sets custom time patterns and separator:
// - Long time:  "HH.mm.ss"
// - Short time: "HH.mm.ss"
// - Separator:  "."
func Set_time_pattern() error {
	const (
		keyPath       = `Control Panel\International`
		sTimeFormat   = "HH.mm.ss"
		sShortTime    = "HH.mm.ss"
		sTime         = "."
		broadcastMsg  = "Intl"
	)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("sTimeFormat", sTimeFormat); err != nil {
		return fmt.Errorf("❌ Failed to set sTimeFormat: %w", err)
	}
	if err := key.SetStringValue("sShortTime", sShortTime); err != nil {
		return fmt.Errorf("❌ Failed to set sShortTime: %w", err)
	}
	if err := key.SetStringValue("sTime", sTime); err != nil {
		return fmt.Errorf("❌ Failed to set sTime (separator): %w", err)
	}

	fmt.Println("✅ Time format set:")
	fmt.Printf("   Long time  (sTimeFormat): %s\n", sTimeFormat)
	fmt.Printf("   Short time (sShortTime) : %s\n", sShortTime)
	fmt.Printf("   Time separator (sTime)  : %s\n", sTime)

	// Broadcast setting change
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ptr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("🔄 System broadcast completed to apply time settings.")
	return nil
}

// Reset_time_pattern resets long/short time format and separator to system defaults.
func Reset_time_pattern() error {
	const (
		keyPath           = `Control Panel\International`
		defaultTimeFormat = "HH:mm:ss"   // Long time
		defaultShortTime  = "h:mm tt"    // Short time
		defaultSeparator  = ":"
		broadcastMsg      = "Intl"
	)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("sTimeFormat", defaultTimeFormat); err != nil {
		return fmt.Errorf("❌ Failed to reset sTimeFormat: %w", err)
	}
	if err := key.SetStringValue("sShortTime", defaultShortTime); err != nil {
		return fmt.Errorf("❌ Failed to reset sShortTime: %w", err)
	}
	if err := key.SetStringValue("sTime", defaultSeparator); err != nil {
		return fmt.Errorf("❌ Failed to reset sTime (separator): %w", err)
	}

	fmt.Println("✅ Time settings reset to system defaults:")
	fmt.Printf("   Long time  (sTimeFormat): %s\n", defaultTimeFormat)
	fmt.Printf("   Short time (sShortTime) : %s\n", defaultShortTime)
	fmt.Printf("   Time separator (sTime)  : %s\n", defaultSeparator)

	// Broadcast setting change
	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ptr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("🔄 System broadcast completed to apply default time settings.")
	return nil
}

// Set_24_hour_format configures Windows to use 24-hour time by setting iTime = 1.
func Set_24_hour_format() error {
	const (
		keyPath = `Control Panel\International`
		broadcastMsg = "Intl"
	)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("iTime", "1"); err != nil {
		return fmt.Errorf("❌ Failed to set iTime = 1: %w", err)
	}

	fmt.Println("✅ Windows is now configured to use 24-hour time (iTime = 1).")

	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ptr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("🔄 System broadcast completed to apply the setting.")
	return nil
}

// Do_not_use_24_hour_format configures Windows to use 12-hour time by setting iTime = 0.
func Do_not_use_24_hour_format() error {
	const (
		keyPath = `Control Panel\International`
		broadcastMsg = "Intl"
	)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("iTime", "0"); err != nil {
		return fmt.Errorf("❌ Failed to set iTime = 0: %w", err)
	}

	fmt.Println("✅ Windows is now configured to use 12-hour time (iTime = 0).")

	const (
		HWND_BROADCAST   = 0xffff
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	ptr := syscall.StringToUTF16Ptr(broadcastMsg)
	var result uintptr

	_, _, _ = sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(SMTO_ABORTIFHUNG),
		100,
		uintptr(unsafe.Pointer(&result)),
	)

	fmt.Println("🔄 System broadcast completed to apply the setting.")
	return nil
}

// Set_first_day_of_week_Monday sets Monday as the first day of the week in Windows regional settings.
func Set_first_day_of_week_Monday() error {
	const keyPath = `Control Panel\International`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("iFirstDayOfWeek", "0"); err != nil {
		return fmt.Errorf("❌ Failed to set iFirstDayOfWeek = 0: %w", err)
	}

	fmt.Println("✅ First day of the week set to Monday (iFirstDayOfWeek = 0).")
	return nil
}

// Set_first_day_of_week_Sunday sets Sunday as the first day of the week in Windows regional settings.
func Set_first_day_of_week_Sunday() error {
	const keyPath = `Control Panel\International`

	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue("iFirstDayOfWeek", "6"); err != nil {
		return fmt.Errorf("❌ Failed to set iFirstDayOfWeek = 6: %w", err)
	}

	fmt.Println("✅ First day of the week set to Sunday (iFirstDayOfWeek = 6).")
	return nil
}

// Convert_blob_to_raw_github_url transforms a GitHub "blob" URL into a "raw" content URL.
//
// GitHub's web interface often shows files using a URL like:
//   https://github.com/{user}/{repo}/blob/{branch}/{path/to/file}
//
// But to access the raw file directly through GitHub (still using the github.com domain),
// the equivalent raw content URL is:
//   https://github.com/{user}/{repo}/raw/{branch}/{path/to/file}
//
// This function performs the necessary transformation by replacing the "/blob/" segment
// in the URL with "/raw/".
//
// For example:
//   Input:
//     https://github.com/user/repo/blob/main/script.ps1
//   Output:
//     https://github.com/user/repo/raw/main/script.ps1
//
// Parameters:
//   - blob_url: the GitHub "blob" URL to convert
//
// Returns:
//   - The corresponding "raw" content URL
//   - An error if the input does not contain the expected "/blob/" segment
func Convert_blob_to_raw_github_url(blob_url string) (string, error) {
	const blob_segment = "/blob/"
	const raw_segment = "/raw/"

	if !strings.Contains(blob_url, blob_segment) {
		return "", fmt.Errorf("❌ input does not contain /blob/: %s", blob_url)
	}

	raw_url := strings.Replace(blob_url, blob_segment, raw_segment, 1)
	return raw_url, nil
}

// Add_to_ps_module_path resolves the appropriate parent directory
// based on a given file or folder path and adds it to the system-wide PSModulePath
//
// Supported input:
// - .psm1 or .psd1 file => adds the grandparent directory
// - folder with only .psm1 or .psd1 files => adds parent of that folder
// - otherwise adds folder itself
func Add_to_ps_module_path(input_path string) error {
	resolved_path, err := filepath.Abs(input_path)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve path: %w", err)
	}

	info, err := os.Stat(resolved_path)
	if err != nil {
		return fmt.Errorf("❌ Path error: %w", err)
	}

	var directory_to_add string

	if !info.IsDir() {
		ext := strings.ToLower(filepath.Ext(resolved_path))
		if ext == ".psm1" || ext == ".psd1" {
			module_folder := filepath.Dir(resolved_path)
			directory_to_add = filepath.Dir(module_folder)
		} else {
			return errors.New("❌ Input is a file but not .psm1 or .psd1")
		}
	} else {
		dir_entries, err := os.ReadDir(resolved_path)
		if err != nil {
			return fmt.Errorf("❌ Failed to read folder: %w", err)
		}

		has_psm1 := false
		invalid := false
		for _, entry := range dir_entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			switch ext {
			case ".psm1":
				has_psm1 = true
			case ".psd1":
				// optional
			default:
				invalid = true
			}
		}

		if has_psm1 && !invalid {
			directory_to_add = filepath.Dir(resolved_path)
		} else {
			directory_to_add = resolved_path
		}
	}

	// Trim trailing slashes for comparison
	directory_to_add = strings.TrimRight(directory_to_add, `\`)

	// Modify the registry
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open environment registry key: %w", err)
	}
	defer key.Close()

	currentPath, _, err := key.GetStringValue("PSModulePath")
	if err != nil {
		return fmt.Errorf("❌ Failed to read PSModulePath: %w", err)
	}

	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		if strings.EqualFold(strings.TrimRight(p, `\`), directory_to_add) {
			fmt.Println("⚠️ Already exists in PSModulePath:", directory_to_add)
			return nil
		}
	}

	newPath := currentPath + ";" + directory_to_add
	if err := key.SetStringValue("PSModulePath", newPath); err != nil {
		return fmt.Errorf("❌ Failed to update PSModulePath: %w", err)
	}

	fmt.Println("✅ Added to PSModulePath:", directory_to_add)
	return nil
}

// Remove_from_ps_module_path resolves the appropriate parent directory
// based on a given file or folder path and removes it from the system-wide PSModulePath.
func Remove_from_ps_module_path(input_path string) error {
	resolved_path, err := filepath.Abs(input_path)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve path: %w", err)
	}

	info, err := os.Stat(resolved_path)
	if err != nil {
		return fmt.Errorf("❌ Path error: %w", err)
	}

	var directory_to_remove string

	if !info.IsDir() {
		ext := strings.ToLower(filepath.Ext(resolved_path))
		if ext == ".psm1" || ext == ".psd1" {
			module_folder := filepath.Dir(resolved_path)
			directory_to_remove = filepath.Dir(module_folder)
		} else {
			return errors.New("❌ Input is a file but not .psm1 or .psd1")
		}
	} else {
		dir_entries, err := os.ReadDir(resolved_path)
		if err != nil {
			return fmt.Errorf("❌ Failed to read folder: %w", err)
		}

		has_psm1 := false
		invalid := false
		for _, entry := range dir_entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			switch ext {
			case ".psm1":
				has_psm1 = true
			case ".psd1":
				// optional
			default:
				invalid = true
			}
		}

		if has_psm1 && !invalid {
			directory_to_remove = filepath.Dir(resolved_path)
		} else {
			directory_to_remove = resolved_path
		}
	}

	directory_to_remove = strings.TrimRight(directory_to_remove, `\`)

	// Modify the registry
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open environment registry key: %w", err)
	}
	defer key.Close()

	currentPath, _, err := key.GetStringValue("PSModulePath")
	if err != nil {
		return fmt.Errorf("❌ Failed to read PSModulePath: %w", err)
	}

	paths := strings.Split(currentPath, ";")
	newPaths := make([]string, 0, len(paths))
	found := false
	for _, p := range paths {
		if strings.EqualFold(strings.TrimRight(p, `\`), directory_to_remove) {
			found = true
			continue // skip this entry
		}
		newPaths = append(newPaths, p)
	}

	if !found {
		fmt.Println("⚠️ Path not found in PSModulePath:", directory_to_remove)
		return nil
	}

	newPath := strings.Join(newPaths, ";")
	if err := key.SetStringValue("PSModulePath", newPath); err != nil {
		return fmt.Errorf("❌ Failed to update PSModulePath: %w", err)
	}

	fmt.Println("✅ Removed from PSModulePath:", directory_to_remove)
	return nil
}

// Enable_SSH ensures the "sshd" service is set to Automatic and Running.
func Enable_SSH() error {
	serviceName := "sshd"

	for {
		// Check if service exists
		checkCmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Get-Service -Name '%s'`, serviceName))
		if err := checkCmd.Run(); err != nil {
			return fmt.Errorf("❌ Service '%s' not found", serviceName)
		}

		// Get current status and start type
		var queryOutput bytes.Buffer
		queryCmd := exec.Command("powershell", "-Command", fmt.Sprintf(`$s = Get-Service -Name '%s'; $mode = (Get-CimInstance -ClassName Win32_Service -Filter "Name='%s'").StartMode; "$($s.Status)|$mode"`, serviceName, serviceName))
		queryCmd.Stdout = &queryOutput
		if err := queryCmd.Run(); err != nil {
			return fmt.Errorf("❌ Failed to query service state: %w", err)
		}

		output := strings.TrimSpace(queryOutput.String())
		parts := strings.Split(output, "|")
		if len(parts) != 2 {
			return fmt.Errorf("❌ Unexpected service query output: %s", output)
		}

		status := parts[0]
		startMode := parts[1]

		fmt.Printf("🔎 Current State — Name: %s | Status: %s | StartType: %s\n", serviceName, status, startMode)

		changed := false

		if !strings.EqualFold(startMode, "Auto") {
			fmt.Println("⚙️ Setting StartType to 'Automatic'...")
			setCmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Set-Service -Name '%s' -StartupType Automatic`, serviceName))
			if err := setCmd.Run(); err != nil {
				return fmt.Errorf("❌ Failed to set start mode: %w", err)
			}
			changed = true
		}

		if !strings.EqualFold(status, "Running") {
			fmt.Println("🚀 Starting SSHD service...")
			startCmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Start-Service -Name '%s'`, serviceName))
			if err := startCmd.Run(); err != nil {
				return fmt.Errorf("❌ Failed to start sshd: %w", err)
			}
			changed = true
		}

		if !changed {
			fmt.Println("✅ SSHD is Running and set to Automatic. Done.")
			break
		}

		time.Sleep(2 * time.Second)
	}

	// Final confirmation
	fmt.Println("📋 Final State:")
	var finalOutput bytes.Buffer
	statusCmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Get-Service -Name '%s' | Select-Object Name, Status, StartType`, serviceName))
	statusCmd.Stdout = &finalOutput
	if err := statusCmd.Run(); err != nil {
		return fmt.Errorf("❌ Failed to fetch final state: %w", err)
	}
	fmt.Print(finalOutput.String())

	return nil
}

// run_powershell returns trimmed stdout of a PowerShell command.
func run_powershell(command string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("❌ PowerShell error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Enable_SSH_through_firewall ensures that TCP port 22 is allowed in the firewall for all profiles.
func Enable_SSH_through_firewall() error {
	// Step 1: Get active network profile (Public/Private/Domain)
	active_profile, err := run_powershell(`(Get-NetConnectionProfile).NetworkCategory`)
	if err != nil {
		return fmt.Errorf("❌ Failed to get active network profile: %v", err)
	}
	fmt.Printf("🌐 Active profile: %s\n", active_profile)

	// Step 2: Check if any enabled inbound allow rule permits TCP port 22 AND matches profile
	check_rule := `
$port = 22
$match = Get-NetFirewallRule -Enabled True -Direction Inbound -Action Allow |
Where-Object {
    ($_ | Get-NetFirewallPortFilter).Protocol -eq "TCP" -and
    ($_ | Get-NetFirewallPortFilter).LocalPort -eq "$port" -and
    ($_ | Get-NetFirewallProfile).Profile -match "` + active_profile + `"
}
if ($match) { "exists" } else { "missing" }`

	rule_status, err := run_powershell(check_rule)
	if err != nil {
		return fmt.Errorf("❌ Failed to check existing firewall rule: %v", err)
	}

	if rule_status == "exists" {
		fmt.Println("✅ SSH firewall rule already exists for profile:", active_profile)
		return nil
	}

	fmt.Printf("🔐 No rule found for SSH on port 22. Creating rule for all profiles...\n")

	// Step 3: Create rule for all profiles
	create_rule := `
New-NetFirewallRule -Name "Allow-SSH" -DisplayName "Allow SSH on Port 22" `
	create_rule += `-Enabled True -Direction Inbound -Protocol TCP -Action Allow `
	create_rule += `-LocalPort 22 -Profile Domain,Private,Public -ErrorAction Stop`

	_, err = run_powershell(create_rule)
	if err != nil {
		return fmt.Errorf("❌ Failed to create firewall rule: %v", err)
	}

	// Step 4: Verify the rule was created correctly
	verify_rule := `
$r = Get-NetFirewallRule -Name "Allow-SSH" -ErrorAction Stop
if ($r.Enabled -eq 'True' -and $r.Direction -eq 'Inbound' -and $r.Action -eq 'Allow') {
    "verified"
} else {
    "mismatch"
}`

	verify_status, err := run_powershell(verify_rule)
	if err != nil {
		return fmt.Errorf("❌ Failed to verify created firewall rule: %v", err)
	}

	if verify_status == "verified" {
		fmt.Println("✅ Rule 'Allow-SSH' successfully created and verified.")
	} else {
		fmt.Println("⚠️ Rule 'Allow-SSH' created but verification failed.")
	}

	return nil
}

// Set_system_environment_variable sets a system-wide environment variable in the registry under HKLM.
// It also broadcasts the environment change so that Explorer and other processes recognize the update.
func Set_system_environment_variable(variable_name string, variable_value string) error {
	if variable_name == "" {
		return fmt.Errorf("❌ Variable name cannot be empty")
	}

	fmt.Printf("🧾 Setting system environment variable: %s = %s\n", variable_name, variable_value)

	// Step 1: Open the registry key
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	// Step 2: Set the value
	if err := key.SetStringValue(variable_name, variable_value); err != nil {
		return fmt.Errorf("❌ Failed to set environment variable: %w", err)
	}

	fmt.Println("✅ Environment variable written to registry.")

	// Step 3: Broadcast the environment change
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

// Get_primary_ipv4_address returns the most appropriate local IPv4 address
// from the available network interfaces.
//
// It prioritizes interfaces whose names contain preferred keywords such as
// "Wi-Fi", "Ethernet", or "Tailscale", and excludes interfaces that are
// likely virtual, loopback, or otherwise irrelevant, such as those containing
// "VMware", "Virtual", "Bluetooth", "Loopback", "OpenVPN", or "Disconnected".
//
// The function performs the following steps:
//   1. Lists all active, non-loopback interfaces.
//   2. Filters out interfaces matching any excluded keywords.
//   3. Searches for an interface whose name contains a preferred keyword.
//   4. Falls back to any remaining valid interface if no preferred one is found.
//   5. Returns the first usable IPv4 address found.
//
// Returns the IPv4 address as a string, or an empty string if none are found.
// If an error occurs while listing interfaces, it is returned.
func Get_primary_ipv4_address() (string, error) {
	preferred_keywords := []string{"Wi-Fi", "Ethernet", "Tailscale"}
	excluded_keywords := []string{"VMware", "Virtual", "Bluetooth", "Loopback", "OpenVPN", "Disconnected"}

	all_interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("❌ failed to get network interfaces: %w", err)
	}

	var candidates []net.Interface

	// Step 1: Filter interfaces that are up and not excluded
	for _, iface := range all_interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if is_excluded_interface(iface.Name, excluded_keywords) {
			continue
		}
		candidates = append(candidates, iface)
	}

	// Step 2: Try preferred keywords
	for _, keyword := range preferred_keywords {
		for _, iface := range candidates {
			if strings.Contains(strings.ToLower(iface.Name), strings.ToLower(keyword)) {
				ip, err := get_ipv4_from_interface(iface)
				if err == nil && ip != "" {
					return ip, nil
				}
			}
		}
	}

	// Step 3: Fallback to any remaining candidate
	for _, iface := range candidates {
		ip, err := get_ipv4_from_interface(iface)
		if err == nil && ip != "" {
			return ip, nil
		}
	}

	return "", nil
}

// get_ipv4_from_interface extracts the first usable IPv4 address from a network interface.
func get_ipv4_from_interface(iface net.Interface) (string, error) {
	addresses, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addresses {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil && ip.To4() != nil && !ip.IsLoopback() && !ip.IsUnspecified() && !ip.IsLinkLocalUnicast() {
			return ip.String(), nil
		}
	}
	return "", nil
}

// is_excluded_interface checks if an interface name matches any excluded keywords.
func is_excluded_interface(interface_name string, excluded_keywords []string) bool {
	lower_name := strings.ToLower(interface_name)
	for _, keyword := range excluded_keywords {
		if strings.Contains(lower_name, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// Restart_file_explorer uses PowerShell to stop and restart Windows File Explorer,
// and waits until explorer.exe is running again.
func Restart_file_explorer() error {
	log.Println("🔄 Stopping Explorer via PowerShell...")

	cmdKill := exec.Command("powershell.exe", "-Command", `Stop-Process -Name explorer -Force`)
	if err := cmdKill.Run(); err != nil {
		log.Printf("❌ Failed to stop Explorer via PowerShell: %v", err)
		return err
	}

	time.Sleep(1 * time.Second)

	log.Println("🚀 Starting Explorer via PowerShell...")

	cmdStart := exec.Command("powershell.exe", "-Command", `Start-Process explorer.exe`)
	if err := cmdStart.Run(); err != nil {
		log.Printf("❌ Failed to start Explorer via PowerShell: %v", err)
		return err
	}

	log.Println("⏳ Waiting for Explorer to relaunch...")

	// Poll for explorer.exe to appear again
	for i := 0; i < 10; i++ {
		cmdCheck := exec.Command("powershell.exe", "-Command", `Get-Process explorer -ErrorAction SilentlyContinue`)
		if err := cmdCheck.Run(); err == nil {
			log.Println("✅ Explorer process is running.")
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("⚠️ Timeout: Explorer process did not appear.")
	return nil
}

// Get_file_size returns the size in bytes of the specified path.
// If the path is a regular file, its size is returned directly.
// If the path is a directory, the function walks through all files
// and returns the cumulative size of all non-directory files within it.
//
// Parameters:
//   - path: The path to the file or directory.
//
// Returns:
//   - int64: Total size in bytes.
//   - error: Any error encountered while accessing the file system.
//
// Example:
//   size, err := Get_file_size("C:\\Users\\Administrator\\Desktop")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println("Total size:", size)
func Get_file_size(path string) (int64, error) {
	var totalSize int64

	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	err = filepath.Walk(path, func(_ string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			totalSize += fileInfo.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

// Get_file_size_human_readable returns the size of a file or directory at the given path
// in a human-readable string format with three decimal places of precision.
//
// The function supports size units in ascending order:
// - bytes
// - KB (Kilobytes)
// - MB (Megabytes)
// - GB (Gigabytes)
// - TB (Terabytes)
//
// For files, the function directly returns the file size in the most appropriate unit.
// For directories, it recursively calculates the total size of all non-directory files inside.
//
// Parameters:
//   - path: The file or directory path as a string.
//
// Returns:
//   - string: A formatted string representing the human-readable size (e.g., "123.456 MB").
//   - error: An error if the path does not exist or cannot be read.
//
// Example:
//   sizeStr, err := Get_file_size_human_readable("C:\\Users\\Administrator\\Desktop")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println("Size:", sizeStr)
func Get_file_size_human_readable(path string) (string, error) {
	size, err := Get_file_size(path)
	if err != nil {
		return "", err
	}

	const (
		_          = iota
		kilobyte   = 1 << (10 * iota)
		megabyte
		gigabyte
		terabyte
	)

	switch {
	case size >= terabyte:
		return fmt.Sprintf("%.3f TB", float64(size)/float64(terabyte)), nil
	case size >= gigabyte:
		return fmt.Sprintf("%.3f GB", float64(size)/float64(gigabyte)), nil
	case size >= megabyte:
		return fmt.Sprintf("%.3f MB", float64(size)/float64(megabyte)), nil
	case size >= kilobyte:
		return fmt.Sprintf("%.3f KB", float64(size)/float64(kilobyte)), nil
	default:
		return fmt.Sprintf("%d bytes", size), nil
	}
}

// Bring_back_the_right_click_menu enables the classic Windows 10-style
// context menu in Windows 11 by setting a specific registry key.
//
// It creates the following key in the current user registry hive:
//   HKEY_CURRENT_USER\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32
// and sets its default value to an empty string.
//
// After applying the registry modification, the function restarts
// Windows File Explorer using the Restart_file_explorer function to apply the change.
//
// Returns an error if the registry key cannot be written or Explorer fails to restart.
func Bring_back_the_right_click_menu() error {
	const keyPath = `Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32`

	// Open or create the registry key
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create/open registry key: %w", err)
	}
	defer key.Close()

	// Set the default value to an empty string
	if err := key.SetStringValue("", ""); err != nil {
		return fmt.Errorf("failed to set registry value: %w", err)
	}

	log.Println("✅ Right-click menu registry tweak applied.")

	// Restart File Explorer to apply changes
	if err := Restart_file_explorer(); err != nil {
		return fmt.Errorf("failed to restart Explorer: %w", err)
	}

	return nil
}

// Use_Windows_11_right_click_menu restores the default Windows 11-style
// right-click context menu by removing a specific registry override.
//
// It deletes the following keys from the current user registry hive:
//   HKEY_CURRENT_USER\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}
//   HKEY_CURRENT_USER\Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}\InprocServer32
//
// These keys are used to force the classic Windows 10-style context menu in Windows 11.
// Removing them reverts Explorer to its default behavior.
//
// After deleting the registry keys, the function restarts Windows File Explorer
// via Restart_file_explorer to apply the change.
//
// Returns an error if any key deletion fails (unless the key doesn't exist),
// or if restarting Explorer fails.
func Use_Windows_11_right_click_menu() error {
	const baseKeyPath = `Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}`

	// Delete the entire CLSID key to revert to the Windows 11 context menu
	err := registry.DeleteKey(registry.CURRENT_USER, baseKeyPath+`\InprocServer32`)
	if err != nil && err != syscall.ERROR_FILE_NOT_FOUND {
		return fmt.Errorf("failed to delete InprocServer32 subkey: %w", err)
	}

	err = registry.DeleteKey(registry.CURRENT_USER, baseKeyPath)
	if err != nil && err != syscall.ERROR_FILE_NOT_FOUND {
		return fmt.Errorf("failed to delete CLSID key: %w", err)
	}

	log.Println("🔄 Restored Windows 11 right-click menu by removing registry override.")

	// Restart Explorer to apply the change
	if err := Restart_file_explorer(); err != nil {
		return fmt.Errorf("failed to restart Explorer: %w", err)
	}

	return nil
}

// Enable_long_file_paths enables long file path support in Windows by setting
// LongPathsEnabled=1 under HKLM\SYSTEM\CurrentControlSet\Control\FileSystem.
// It first checks the current value using Are_long_file_paths_enabled() and only
// modifies the registry if needed. Administrator privileges are required.
func Enable_long_file_paths() error {
	const registryPath = `SYSTEM\CurrentControlSet\Control\FileSystem`
	const valueName = "LongPathsEnabled"

	enabled, err := Are_long_file_paths_enabled()
	if err != nil {
		return fmt.Errorf("❌ Could not check if long file paths are enabled: %w", err)
	}

	if enabled {
		fmt.Println("ℹ️ Long file paths are already enabled (LongPathsEnabled = 1).")
		return nil
	}

	// Open the registry key with write access to change the value
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.READ|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("❌ Failed to open registry key (requires Admin): %w", err)
	}
	defer key.Close()

	if err := key.SetDWordValue(valueName, 1); err != nil {
		return fmt.Errorf("❌ Failed to set %s = 1: %w", valueName, err)
	}

	fmt.Println("✅ Long file paths have been enabled (LongPathsEnabled = 1).")

	return nil
}

// Are_long_file_paths_enabled checks if long file path support is currently enabled.
// It returns true if LongPathsEnabled == 1, false otherwise.
// Returns an error if the registry cannot be read (e.g., due to insufficient privileges).
func Are_long_file_paths_enabled() (bool, error) {
	const registryPath = `SYSTEM\CurrentControlSet\Control\FileSystem`
	const valueName = "LongPathsEnabled"

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.READ)
	if err != nil {
		return false, fmt.Errorf("❌ Failed to open registry key: %w", err)
	}
	defer key.Close()

	currentVal, _, err := key.GetIntegerValue(valueName)
	if err != nil {
		return false, fmt.Errorf("❌ Failed to read registry value: %w", err)
	}

	return currentVal == 1, nil
}