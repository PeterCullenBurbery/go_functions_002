# Changelog

All notable changes to this project will be documented in this file.

## [3.0.0] - 2025-007-005@006.036 PM

### Added
- Added folder for version 3.
- Removed folder for v2.

## [2.9.0] - 2025-007-005@006.019 PM

- Added Enable_SSH. Enable_SSH ensures the "sshd" service is set to Automatic and Running.

## [2.8.0] - 2025-007-005@012.028 PM

- Add_to_ps_module_path modified. I'm not what the difference is.
- Remove_from_ps_module_path modified. I'm not sure what the difference is.
- I'm not sure what the difference is.

## [2.7.0] - 2025-007-005@012.001 PM

- Added Add_to_ps_module_path. Add_to_ps_module_path adds the given directory to the system-wide PSModulePath environment variable.
- Added Remove_from_ps_module_path. Remove_from_ps_module_path removes the given directory from the system-wide PSModulePath environment variable.

## [2.6.0] - 2025-007-005@011.038 AM

- Added Convert_blob_to_raw_github_url. Convert_blob_to_raw_github_url transforms a GitHub "blob" URL into a "raw" content URL.

## [2.5.0] - 2025-007-004@010.005 PM

### Added
- Added Set_first_day_of_week_Monday. Set_first_day_of_week_Monday sets Monday as the first day of the week in Windows regional settings.
- Added Set_first_day_of_week_Sunday. Set_first_day_of_week_Sunday sets Sunday as the first day of the week in Windows regional settings.
- Dang! That's close!.

## [2.4.0] - 2025-007-004@009.059 PM

### Added
- Added Set_24_hour_format. Set_24_hour_format configures Windows to use 24-hour time by setting iTime = 1.
- Added Do_not_use_24_hour_format. Do_not_use_24_hour_format configures Windows to use 12-hour time by setting iTime = 0.

## [2.3.0] - 2025-007-004@009.054 PM

### Added
- Added Set_time_pattern.
  - Set_time_pattern sets custom time patterns and separator:
    - Long time:  "HH.mm.ss"
    - Short time: "HH.mm.ss"
    - Separator:  "."
- Added Reset_time_pattern. Reset_time_pattern resets long/short time format and separator to system defaults.

## [2.2.0] - 2025-007-004@007.023 PM

### Added
- Added folder for v2.

## [2.1.0] - 2025-007-004@006.018 PM

### Added
- Added Set_long_date_pattern. Set_long_date_pattern sets the long date pattern to "yyyy-MM-dd-dddd" and broadcasts the change to the system.
- Added Reset_long_date_pattern. Reset_long_date_pattern resets the long date pattern to the default "dddd, MMMM d, yyyy" and broadcasts the change to the system.

## [2.0.0] - 2025-007-004@006.012 PM

### Added
- Added Set_short_date_pattern. Set_short_date_pattern sets the short date pattern to "yyyy-MM-dd-dddd" and broadcasts the change to the system.
- Added Reset_short_date_pattern. Reset_short_date_pattern resets the short date pattern to "M/d/yyyy" and broadcasts the change to the system.
- Okay back to short timing between repository commits.

## [1.9.0] - 2025-007-004@004.057 PM

### Added
- Added Seconds_in_taskbar. Seconds_in_taskbar enables seconds on the taskbar clock by setting ShowSecondsInSystemClock = 1.
- Added Take_seconds_out_of_taskbar. Take_seconds_out_of_taskbar disables seconds on the taskbar clock by setting ShowSecondsInSystemClock = 0.

## [1.8.0] - 2025-007-004@004.054 PM

### Added
- Added Hide_search_box. Hide_search_box sets SearchboxTaskbarMode = 0 to hide the taskbar search box.
- Added Do_not_hide_search_box. Do_not_hide_search_box sets SearchboxTaskbarMode = 2 to show the full search box on the taskbar.
- Dang! That's close. 2 minutes.

## [1.7.0] - 2025-007-004@004.052 PM

### Added
- Added Show_hidden_files. Show_hidden_files sets Hidden = 1 to show hidden files in File Explorer.
- Added Do_not_show_hidden_files. Do_not_show_hidden_files sets Hidden = 2 to hide hidden files in File Explorer.
- Dang! That's also close.

## [1.6.0] - 2025-007-004@004.046 PM

### Added
- Added Show_file_extensions. Show_file_extensions sets HideFileExt = 0 to make file extensions visible.
- Added Do_not_show_file_extensions. Do_not_show_file_extensions sets HideFileExt = 1 to hide file extensions.
- Dang, that's close on the heels of 2025-007-004@004.044 PM! 2 minutes! 1.5.0.

## [1.5.0] - 2025-007-004@004.044 PM

### Added
- Added Set_start_menu_to_left. Set_start_menu_to_left sets the Windows 11 taskbar alignment to the left by writing TaskbarAl=0 in the registry and restarting Explorer.
- Added Set_start_menu_to_center. Set_start_menu_to_center sets the Windows 11 taskbar alignment to the center by writing TaskbarAl=1 in the registry and restarting Explorer.

## [1.4.0] - 2025-007-004@004.030 PM

### Added
- Added Set_dark_mode. Set_dark_mode sets Windows to dark mode for both system and apps.
- Added Set_light_mode. Set_light_mode sets Windows to light mode for both system and apps.

## [1.3.0] - 2025-007-004@003.010 PM

### Added
- Added Extract_password_protected_zip. Extract_password_protected_zip extracts a password-protected ZIP archive using AES or ZipCrypto.

## [1.2.0] - 2025-007-004@001.057 PM

### Added
- Added Exclude_from_Microsoft_Windows_Defender. Exclude_from_Microsoft_Windows_Defender excludes the given file or folder from Microsoft Defender. If a file is given, its parent folder is excluded instead.

## [1.1.0] - 2025-007-004@001.026 PM

### Added
- Added Remove_from_path. Remove_from_path removes the given path from the system PATH if present. It normalizes the path, modifies HKLM registry, and broadcasts environment changes.
- Added Create_desktop_shortcut. Create_desktop_shortcut creates a .lnk shortcut on the desktop. It accepts the target path, shortcut name (optional), description (optional), window style (3 = maximized), and allUsers flag.
- Added Extract_zip. Extract_zip extracts a ZIP archive specified by src into the destination directory dest. It protects against Zip Slip attacks by ensuring all extracted paths are within dest.

## [1.0.0] - 2025-007-004@012.023 PM

### Added
- Added Add_to_path. Add_to_path adds the given path to the top of the system PATH (HKLM) if not already present. It broadcasts the environment change so apps like Explorer pick it up.

## [0.9.0] - 2025-007-004@011.048 AM

### Added
- Added Download_file. Download_file downloads a file from the given URL and saves it to the specified destination path.

## [0.8.0] - 2025-006-029@007.045 PM

- I don't know what changed.

## [0.7.2] - 2025-006-029@005.021 PM

### Updated
- Updated fileExists to File_exists.

## [0.7.1] - 2025-006-029@004.042 PM

### Updated
- `Date_time_stamp`: Now automatically installs Java using `system_management_functions` if it is not found, simplifying usage by requiring no arguments.

## [0.7.0] - 2025-006-029@004.035 PM

### Added

- Added Is_Java_installed. Is_Java_installed checks if both java.exe and javac.exe are available in PATH, or in the default Eclipse Adoptium installation directory.
- Added fileExists. fileExists checks if a file exists and is not a directory.
- Added Install_Java. Install_Java ensures Java is installed by checking Is_Java_installed(). If not found, it installs the temurin21 JDK via Chocolatey.

## [0.6.0] - 2025-006-029@004.014 PM

### Added
- Added Is_Choco_installed. Is_Choco_installed checks if Chocolatey is installed.

### Updated
- `Choco_install`: Now uses `Is_Choco_installed` to verify and optionally trigger Chocolatey installation if missing.

## [0.5.2] - 2025-006-029@006.048 AM

- Choco_install updated.
  - "choco list --limit-output --exact msys2" instead of "choco list --local-only msys2"

## [0.5.1] - 2025-006-028@008.007 PM

- Choco_install updated.
  - Clearer fallback logic using os.Stat instead of a cmd shell workaround.

## [0.5.0] - 2025-006-028@006.059 PM

- Added Install_choco. Install_choco installs Chocolatey using the official PowerShell script.

## [0.4.3] - 2025-006-028@003.013 PM

- I don't understand what the difference is. Something with yaml functions.

## [0.4.2] - 2025-006-027@008.052 PM

- Added Winget_install. Winget_install installs the specified package using winget with standard flags.

## [0.4.1] - 2025-006-027@008.048 PM

- I don't understand what the difference is.

## [0.4.0] - 2025-006-027@008.040 PM

- Added Choco_install. Choco_install installs the given Chocolatey package and checks if it was installed successfully.

## [0.3.0] - 2025-006-027@008.027 PM

- Added 005 yaml functions
  - GetCaseInsensitiveMap
  - GetCaseInsensitiveList
  - GetCaseInsensitiveString
  - GetNestedString
  - GetNestedMap

## [0.2.1] - 2025-006-027@007.006 PM

- Added license. Added MIT license.

## [0.2.0] - 2025-006-027@006.055 PM

- FormatNow added. FormatNow returns the current time formatted as "2006-01-02 15:04:05"

## [0.1.0] - 2025-006-027@006.046 PM

- SayHello added. This function("Peter") says "Hello, Peter!".