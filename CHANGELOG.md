# Changelog

All notable changes to this project will be documented in this file.

Note:

A lot of differences are uncertain.

After 3.3.4, versions become more certain.

I wasn't using ChatGPT/AI to write commits until 3.4.0.

## [5.4.1] - 2025-008-004 008.014.058.855321700 America/New_York 2025-W032-001 2025-216

- update README.md. The links at the top were not spaced out properly.

## [5.4.0] - 2025-008-004 008.007.016.794655300 America/New_York 2025-W032-001 2025-216

### Added
- **`oracle_database_system_management_functions`** package for Oracle CDB/PDB lifecycle management.
  - **Path Helpers**
    - `Get_root_datafile_directory` — Locate CDB$ROOT SYSTEM01.DBF directory.
    - `Get_pdbseed_datafile_directory` — Locate PDB$SEED SYSTEM01.DBF directory.
    - `Verify_pdbseed_directory_matches_expected` — Validate that PDB$SEED path matches expectations.
  - **CDB Container Guard**
    - `Ensure_connected_to_cdb_root` — Verify connected container is CDB$ROOT.
  - **PDB Lifecycle**
    - `Create_pluggable_database_from_seed` — Create a new PDB from PDB$SEED.
    - `Open_pluggable_database_read_write` — Open PDB in READ WRITE mode.
    - `Save_pluggable_database_state` — Save PDB auto-open state.
    - `Get_pdb_status` — Retrieve PDB OPEN_MODE.
    - `Get_saved_state_info` — Retrieve saved state details from DBA_PDB_SAVED_STATES.
    - `Close_pluggable_database_immediate` — Close PDB immediately.
    - `Discard_pluggable_database_state` — Remove PDB saved state.
    - `Drop_pluggable_database_including_datafiles` — Drop PDB and delete datafiles.
    - `Verify_pluggable_database_dropped` — Confirm PDB no longer exists.
  - **Session Management**
    - `Get_user_sessions` — List USER sessions in a PDB.
    - `Kill_user_sessions_in_pdb` — Kill USER sessions (one pass).
    - `Kill_user_sessions_in_pdb_until_gone` — Retry killing until no sessions remain.
  - **Convenience Operations**
    - `Create_open_save_state_pdb_from_seed` — Create, open, save state in one call.
    - `Teardown_drop_pdb` — Close, discard state, and drop PDB.

## [5.3.0] - 2025-008-003 018.021.016.455497600 America/New_York 2025-W031-007 2025-215

### Added
- `date_time_functions.Generate_pdb_name_from_timestamp`
  - Returns a PDB-style name using zero‑padded components:
    - Format: `pdb_YYYY_MMM_DDD_HHH_MMM_SSS`
    - Example: `pdb_2025_007_031_017_020_008`

## [5.2.0] - 2025-008-003 017.057.022.251943100 America/New_York 2025-W031-007 2025-215

- Zenodo [meta]data has been added.
- Update README.md.
  - DOI's added.

## [5.1.1] - 2025-008-003 017.052.023.246101600 America/New_York 2025-W031-007 2025-215

- CITATION.cff incorrectly was "Peter Cullen Burbery Python functions". CITATION.cff has been updated to "Go functions-002".

## [5.1.0] - 2025-008-003 017.010.043.238504300 America/New_York 2025-W031-007 2025-215

- released for Zenodo. released for Zenodo release.

## [5.0.0] - 2025-008-003 016.006.054.065664300 America/New_York 2025-W031-007 2025-215

### Added
- Added folder for version 5.
- Removed folder for v4.

## [4.9.0] - 2025-007-24@003.004 PM

### Added
- Add retry logic to Create_desktop_shortcut.
  - The Create_desktop_shortcut function now retries up to 100 times with a delay if shortcut creation fails. This improves robustness against transient errors when creating desktop shortcuts. I haven't really tested it. I haven't really tested it, really. I should test it at some point.

## [4.8.0] - 2025-007-24@001.054 PM

### Added
- Added Clean_path. Clean_path cleans the system PATH (HKLM) by removing duplicates and expanding environment variables. It updates the registry and broadcasts the environment change to Explorer.

## [4.7.0] - 2025-007-22@006.023 PM

- Improve desktop shortcut creation robustness and logging
  - Added ensure desktop shortcut exists to Create_desktop_shortcut.

## [4.6.0] - 2025-007-22@009.053 AM

- Improve Remove_from_path with better matching and user advice
  - Enhanced the Remove_from_path function to normalize and compare PATH entries in a case-insensitive manner, handle environment variable expansion, and provide clearer output. Added a check for 'refreshenv' availability and advise the user to run it if present, improving the user experience after modifying the system PATH.

## [4.5.0] - 2025-007-22@009.050 AM

- Add refreshenv message to Add_to_path.

## [4.4.1] - 2025-007-22@008.007 AM

- Fixed Exclude_from_Microsoft_Windows_Defender exclude from Defender only if Microsoft Windows Defender is running. Before Exclude_Defender would error out when asked. The function now works as intended by not adding exclusion if Microsoft Windows Defender is not running.

## [4.4.0] - 2025-007-22@007.002 AM

- Updated Exclude_from_Microsoft_Windows_Defender.
  - Exclude_from_Microsoft_Windows_Defender attempts to find if Defender is running and adds exclusion only if Microsoft Windows Defender is running.

## [4.3.0] - 2025-007-019@001.015 PM
- Added Are_long_file_paths_enabled. Are_long_file_paths_enabled checks if long file path support is currently enabled. It returns true if LongPathsEnabled == 1, false otherwise.
- Updated Enable_long_file_paths to use Are_long_file_paths_enabled() for checking current state before modifying the registry. This improves code clarity and avoids redundant registry access.

## [4.2.0] - 2025-007-017@005.010 PM

- Add function to enable long file paths in Windows.
  - Introduces Enable_long_file_paths, which sets LongPathsEnabled=1 in the Windows registry to allow long file path support. The function checks the current value before updating and requires administrator privileges.

## [4.1.0] - 2025-007-013@009.028 PM

### Removed
- Removed folder for v3.

### Added
- Add function to restore classic right-click menu.
  - Introduces Bring_back_the_right_click_menu, which applies a registry tweak to restore the classic Windows right-click context menu and restarts File Explorer to apply the change.
- Add functions to toggle Windows 11 right-click menu.
  - Introduces Bring_back_the_right_click_menu to enable the classic Windows 10-style context menu and Use_Windows_11_right_click_menu to restore the default Windows 11 context menu. Both functions modify specific registry keys and restart Explorer to apply changes.
- Added documentation to Bring_back_the_right_click_menu.
- Added Use_Windows_11_right_click_menu.

## [4.0.0] - 2025-007-012@008.026 PM

### Added
- Added folder for version 4.
- Add human-readable file size function
  - Introduced Get_file_size_human_readable, which returns the size of a file or directory in a human-readable string format (bytes, KB, MB, GB, TB) with three decimal places. This complements the existing Get_file_size function by providing a more user-friendly output.

## [3.9.0] - 2025-007-012@007.019 PM

- Add Get_file_size function to calculate file or directory size
  - Introduces Get_file_size, which returns the size in bytes of a file or the cumulative size of all files within a directory. This utility aids in determining storage usage for a given path.

## [3.8.1] - 2025-007-011@009.055 PM

- Remove PowerShell version detection functions
  - Deleted the PowershellVersionDetails struct and the Get_powershell_version function, including the embedded PowerShell script and related JSON parsing. This streamlines the code by removing PowerShell version detection logic that is no longer needed.

## [3.8.0] - 2025-007-011@009.046 PM

- Add PowerShell version detection utility
  - Introduces Get_powershell_version, which runs a PowerShell script to detect version details and feature support, returning a structured PowershellVersionDetails object. This helps determine PowerShell capabilities at runtime for improved compatibility handling.

## [3.7.0] - 2025-007-011@010.024 AM

- Add function to restart Windows File Explorer
  - Introduces Restart_file_explorer, which uses PowerShell commands to stop and restart explorer.exe, and waits for the process to relaunch. This utility aids in programmatically refreshing the Windows desktop environment.

## [3.6.0] - 2025-007-010@003.041 PM

- Add function to get primary IPv4 address
  - Introduces Get_primary_ipv4_address, which selects the most appropriate local IPv4 address by prioritizing preferred network interfaces and filtering out virtual, loopback, and irrelevant interfaces. Helper functions for extracting IPv4 addresses and checking excluded interface names are also included.

## [3.5.1] - 2025-007-010@003.035 PM

- Refactor JAVA_HOME setting in Install_Java
  - Replaces direct PowerShell command with a call to Set_system_environment_variable for setting JAVA_HOME. This improves code maintainability and centralizes environment variable management.

## [3.5.0] - 2025-007-010@003.031 PM

- Add function to set system environment variable
  - Introduces Set_system_environment_variable, which sets a system-wide environment variable in the Windows registry and broadcasts the change so that other processes recognize the update. Includes error handling and status messages for each step.

## [3.4.1] - 2025-007-009@007.049 PM

- Updated Topological_sort to sort nodes with the same precedence alphabetically, ensuring deterministic output. This affects both Topological_sort and Reverse_topological_sort, improving consistency for repeated runs with the same input graph.

## [3.4.0] - 2025-007-009@007.041 PM

- Added math_functions.
  - Added Topological_sort. Topological_sort performs a topological sort on a DAG using Kahn's algorithm. Returns the ordered list of tasks, or an error if a cycle is detected.
  - Reverse_topological_sort performs a topological sort and returns the reversed order. Useful for teardown operations or viewing leaf-to-root dependencies.

## [3.3.4] - 2025-007-007@006.020 PM

- Add_to_path's code updated.

## [3.3.3] - 2025-007-007@003.046 PM

- Added Expand_windows_env. Expand_windows_env expands environment variables using the Windows API. For example, %SystemRoot% becomes C:\Windows. Add_to_path now uses Expand_windows_env, instead of os.ExpandEnv.

## [3.3.2] - 2025-007-007@011.032 AM

- Add_to_path now prints fmt.Printf("📝 New PATH to be written:\n%s\n", path to be written.

## [3.3.1] - 2025-007-007@010.049 AM

### Improved
- Updated Add_to_path:
  - Now rewrites the entire PATH with fully expanded, deduplicated, and normalized entries.
  - Removes duplicates even if they differ by case, trailing slashes, or use of environment variables (e.g., `%SystemRoot%`).
  - Stores only literal absolute paths, eliminating any `%VAR%` references in PATH.

## [3.3.0] - 2025-007-007@008.056 AM

- Updated Add_to_path/Remove_from_path to expand environment variables before comparing.

## [3.2.1] - 2025-007-006@003.039 PM

- Update date_time_functions to use "github.com/PeterCullenBurbery/go_functions_002/v3/system_management_functions" instead of "github.com/PeterCullenBurbery/go_functions_002/v2/system_management_functions".

## [3.2.0] - 2025-007-005@007.048 PM

- Updated Install_Java. Install_Java now sets JAVA_HOME to `C:\Program Files\Eclipse Adoptium\jdk-21.0.6.7-hotspot` using ```(`[Environment]::SetEnvironmentVariable("JAVA_HOME", "%s", "Machine")`, java_home)```.

## [3.1.1] - 2025-007-005@007.029 PM

- Renamed Enable_ssh_through_firewall to Enable_SSH_through_firewall.

## [3.1.0] - 2025-007-005@007.026 PM

### Improved
- Improved Enable_SSH_through_firewall:
  - Refactored to use `run_powershell` helper for clarity and maintainability.
  - Simplified logic with cleaner PowerShell invocation.

- I don't understand what the difference is. I didn't understand what the difference is. I asked ChatGPT to write a summary (above).

## [3.0.0] - 2025-007-005@006.036 PM

### Added
- Added folder for version 3.
- Removed folder for v2.
- Added Enable_SSH_through_firewall. Enable_SSH_through_firewall ensures that TCP port 22 is allowed in the firewall for all profiles.

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