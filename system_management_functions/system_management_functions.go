package system_management_functions

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"net/http"
	"strings"
	"syscall"
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

// Exclude_from_Microsoft_Windows_Defender excludes the given file or folder from Microsoft Defender.
//
// If a file is given, its parent folder is excluded instead.
// This requires administrator privileges.
//
// Parameters:
//   - path_to_exclude: Absolute path to a file or folder to exclude.
//
// Returns:
//   - An error if exclusion fails; nil otherwise.
func Exclude_from_Microsoft_Windows_Defender(path_to_exclude string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path_to_exclude)
	if err != nil {
		return fmt.Errorf("❌ Failed to resolve absolute path: %w", err)
	}

	// Stat to determine if it's a file or folder
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("❌ Failed to stat path: %w", err)
	}

	// If it's a file, get parent directory
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}

	// Normalize (trim trailing slash)
	normalizedPath := filepath.Clean(absPath)

	// Build PowerShell command
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
		fmt.Sprintf(`Add-MpPreference -ExclusionPath "%s"`, normalizedPath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("❌ Failed to exclude from Defender: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Excluded from Microsoft Defender: %s\n", normalizedPath)
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