# go_functions_002

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)  
[![Go Reference](https://pkg.go.dev/badge/github.com/PeterCullenBurbery/go_functions_002/v6.svg)](https://pkg.go.dev/github.com/PeterCullenBurbery/go_functions_002/v6)  
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.16734055.svg)](https://doi.org/10.5281/zenodo.16734055)

**Author:** Peter Cullen Burbery  
**Language:** Go 1.24+  
**Platform:** Primarily Windows, with some cross-platform functions

---

For the most up to date documentation, please visit [go_functions_002](https://pkg.go.dev/github.com/PeterCullenBurbery/go_functions_002/v6).

---

## 📖 Overview

`go_functions_002` is a utility library written in Go, offering a wide range of functions for:

- **System management** (PATH updates, registry tweaks, Windows Explorer settings)
- **Date and time handling**
- **YAML parsing**
- **Mathematical algorithms**
- **Oracle Database CDB/PDB lifecycle management**
- **Package management automation** (Chocolatey, Winget)
- **File and network utilities**

This package is designed for automation, system configuration, database management, and utility scripting — especially for **Windows system administration in Go**.

---

## 📦 Installation

```bash
go get github.com/PeterCullenBurbery/go_functions_002/v6
```

---

## ✨ Features

### 📅 Date/Time Functions
- **`Date_time_stamp()`** – Returns a precise timestamp with nanoseconds, ISO week, and ordinal date.
- **`Format_now()`** – Returns current time in `"2006-01-02 15:04:05"` format.
- **`Safe_time_stamp()`** – Produces a safe filename timestamp (replaces `/` with ` slash `).
- **`Generate_pdb_name_from_timestamp()`** – Generates a unique PDB name from the current timestamp.
- **`Get_timestamp()`** – Returns an underscore-delimited, time zone–aware, nanosecond-precision timestamp like `2025_008_004_014_017_048_822529300_America_slash_New_York_2025_W032_001_2025_216`.
- **`Get_dash_separated_timestamp()`** – Returns a dash-separated timestamp like `2025-008-005-020-058-035-258752600-America-slash-New-York-2025-W032-002-2025-217`.

---

### 🧮 Math Functions
- **`Topological_sort()`** – Deterministic Kahn’s algorithm, sorts nodes alphabetically when precedence is equal.
- **`Reverse_topological_sort()`** – Returns reversed topological order.

---

### 🗄 Oracle Database System Management Functions
Tools for **Oracle CDB/PDB lifecycle management**:

- **Path Helpers**
  - `Get_root_datafile_directory()` – Locate CDB$ROOT SYSTEM01.DBF directory.
  - `Get_pdbseed_datafile_directory()` – Locate PDB$SEED SYSTEM01.DBF directory.
  - `Verify_pdbseed_directory_matches_expected()` – Validate that PDB$SEED path matches expectations.

- **CDB Container Guards**
  - `Ensure_connected_to_cdb_root()` – Verify connected container is CDB$ROOT.

- **PDB Lifecycle**
  - `Create_pluggable_database_from_seed()` – Create a new PDB from PDB$SEED.
  - `Open_pluggable_database_read_write()` – Open PDB in READ WRITE mode.
  - `Save_pluggable_database_state()` – Save PDB auto-open state.
  - `Get_pdb_status()` – Retrieve PDB OPEN_MODE.
  - `Get_saved_state_info()` – Retrieve saved state details from DBA_PDB_SAVED_STATES.
  - `Close_pluggable_database_immediate()` – Close PDB immediately.
  - `Discard_pluggable_database_state()` – Remove PDB saved state.
  - `Drop_pluggable_database_including_datafiles()` – Drop PDB and delete datafiles.
  - `Verify_pluggable_database_dropped()` – Confirm PDB no longer exists.

- **Session Management**
  - `Get_user_sessions()` – List USER sessions in a PDB.
  - `Kill_user_sessions_in_pdb()` – Kill USER sessions (one pass).
  - `Kill_user_sessions_in_pdb_until_gone()` – Retry killing until no sessions remain.

- **Convenience Operations**
  - `Create_open_save_state_pdb_from_seed()` – Create, open, save state in one call.
  - `Teardown_drop_pdb()` – Close, discard state, and drop PDB.

---

### ⚙️ System Management Functions
Includes 50+ Windows utilities such as:
- **PATH & Environment Management**
  - `Add_to_path`, `Remove_from_path`, `Clean_path`
  - `Add_to_ps_module_path`, `Remove_from_ps_module_path`
  - `Set_system_environment_variable`
- **Windows Explorer Tweaks**
  - `Show_hidden_files`, `Do_not_show_hidden_files`
  - `Show_file_extensions`, `Do_not_show_file_extensions`
  - `Bring_back_the_right_click_menu`, `Use_Windows_11_right_click_menu`
  - `Hide_search_box`, `Do_not_hide_search_box`
  - `Set_dark_mode`, `Set_light_mode`
- **Registry Config**
  - `Set_24_hour_format`, `Do_not_use_24_hour_format`
  - `Set_short_date_pattern`, `Reset_short_date_pattern`
  - `Set_long_date_pattern`, `Reset_long_date_pattern`
  - `Set_time_pattern`, `Reset_time_pattern`
  - `Seconds_in_taskbar`, `Take_seconds_out_of_taskbar`
  - `Enable_long_file_paths`, `Are_long_file_paths_enabled`
- **Security & Install**
  - `Exclude_from_Microsoft_Windows_Defender`
  - `Choco_install`, `Winget_install`, `Install_choco`
  - `Enable_SSH`, `Enable_SSH_through_firewall`
- **File Utilities**
  - `Download_file`
  - `Extract_zip`, `Extract_password_protected_zip`
  - `File_exists`
  - `Get_file_size`, `Get_file_size_human_readable`
  - `Create_desktop_shortcut`
- **Networking**
  - `Get_primary_ipv4_address`

---

### 📄 YAML Functions
Convenience helpers for working with `map[string]interface{}` parsed from YAML:
- `GetCaseInsensitiveString`
- `GetCaseInsensitiveMap`
- `GetCaseInsensitiveList`
- `GetNestedString`
- `GetNestedMap`

---

## 📊 Condensed Features Table

| Category | Function | Description |
|----------|----------|-------------|
| **Date/Time** | `Date_time_stamp()` | Precise timestamp with nanoseconds, week & ordinal date |
|  | `Format_now()` | Current time in `YYYY-MM-DD HH:MM:SS` |
|  | `Safe_time_stamp()` | Replaces `/` with ` slash ` for safe filenames |
|  | `Generate_pdb_name_from_timestamp()` | PDB name based on timestamp |
| **Math** | `Topological_sort()` | Deterministic DAG topological sort |
|  | `Reverse_topological_sort()` | Reverse order of topological sort |
| **Oracle DB** | `Get_root_datafile_directory()` | Locate CDB$ROOT SYSTEM01.DBF directory |
|  | `Get_pdbseed_datafile_directory()` | Locate PDB$SEED SYSTEM01.DBF directory |
|  | `Verify_pdbseed_directory_matches_expected()` | Validate PDB$SEED path |
|  | `Ensure_connected_to_cdb_root()` | Verify CDB$ROOT connection |
|  | `Create_pluggable_database_from_seed()` | Create new PDB from PDB$SEED |
|  | `Open_pluggable_database_read_write()` | Open PDB READ WRITE |
|  | `Save_pluggable_database_state()` | Save auto-open state |
|  | `Get_pdb_status()` | Get PDB open mode |
|  | `Get_saved_state_info()` | Retrieve saved state |
|  | `Close_pluggable_database_immediate()` | Close PDB immediately |
|  | `Discard_pluggable_database_state()` | Remove saved state |
|  | `Drop_pluggable_database_including_datafiles()` | Drop PDB with datafiles |
|  | `Verify_pluggable_database_dropped()` | Confirm PDB is dropped |
|  | `Get_user_sessions()` | List USER sessions |
|  | `Kill_user_sessions_in_pdb()` | Kill USER sessions (one pass) |
|  | `Kill_user_sessions_in_pdb_until_gone()` | Keep killing until none remain |
|  | `Create_open_save_state_pdb_from_seed()` | Create, open, save state in one call |
|  | `Teardown_drop_pdb()` | Close, discard, drop PDB |
| **PATH & Env** | `Add_to_path()` | Add folder to system PATH |
|  | `Remove_from_path()` | Remove folder from PATH |
|  | `Clean_path()` | Deduplicate & normalize PATH |
|  | `Add_to_ps_module_path()` | Add directory to PSModulePath |
|  | `Remove_from_ps_module_path()` | Remove directory from PSModulePath |
|  | `Set_system_environment_variable()` | Set system-wide env variable |
| **Explorer Tweaks** | `Show_hidden_files()` | Show hidden files |
|  | `Do_not_show_hidden_files()` | Hide hidden files |
|  | `Show_file_extensions()` | Show file extensions |
|  | `Do_not_show_file_extensions()` | Hide file extensions |
|  | `Bring_back_the_right_click_menu()` | Enable Win10 context menu on Win11 |
|  | `Use_Windows_11_right_click_menu()` | Restore Win11 default context menu |
|  | `Hide_search_box()` | Hide taskbar search box |
|  | `Do_not_hide_search_box()` | Show taskbar search box |
|  | `Set_dark_mode()` | Set system/apps to dark mode |
|  | `Set_light_mode()` | Set system/apps to light mode |
| **Registry Config** | `Set_24_hour_format()` | 24-hour clock |
|  | `Do_not_use_24_hour_format()` | 12-hour clock |
|  | `Set_short_date_pattern()` | Custom short date format |
|  | `Reset_short_date_pattern()` | Restore default short date format |
|  | `Set_long_date_pattern()` | Custom long date format |
|  | `Reset_long_date_pattern()` | Restore default long date format |
|  | `Set_time_pattern()` | Custom time format |
|  | `Reset_time_pattern()` | Restore default time format |
|  | `Seconds_in_taskbar()` | Show seconds in taskbar clock |
|  | `Take_seconds_out_of_taskbar()` | Hide seconds in taskbar clock |
|  | `Enable_long_file_paths()` | Enable >260 char paths |
|  | `Are_long_file_paths_enabled()` | Check if long paths are enabled |
| **Security & Install** | `Exclude_from_Microsoft_Windows_Defender()` | Exclude path from Defender |
|  | `Choco_install()` | Install Chocolatey package |
|  | `Winget_install()` | Install package via Winget |
|  | `Install_choco()` | Install Chocolatey itself |
|  | `Enable_SSH()` | Enable SSH server |
|  | `Enable_SSH_through_firewall()` | Open SSH port in firewall |
| **File Utilities** | `Download_file()` | Download from URL |
|  | `Extract_zip()` | Extract ZIP |
|  | `Extract_password_protected_zip()` | Extract password-protected ZIP |
|  | `File_exists()` | Check if file exists |
|  | `Get_file_size()` | Get file/dir size in bytes |
|  | `Get_file_size_human_readable()` | Get file/dir size in readable format |
|  | `Create_desktop_shortcut()` | Create Windows `.lnk` shortcut |
| **Networking** | `Get_primary_ipv4_address()` | Get best local IPv4 address |
| **YAML** | `GetCaseInsensitiveString()` | Case-insensitive key lookup (string) |
|  | `GetCaseInsensitiveMap()` | Case-insensitive key lookup (map) |
|  | `GetCaseInsensitiveList()` | Case-insensitive key lookup (list) |
|  | `GetNestedString()` | Nested string lookup |
|  | `GetNestedMap()` | Nested map lookup |

---

## 🖥 Example Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/PeterCullenBurbery/go_functions_002/v6/system_management_functions"
    "github.com/PeterCullenBurbery/go_functions_002/v6/date_time_functions"
)

func main() {
    ts, err := date_time_functions.Date_time_stamp()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Timestamp:", ts)

    err = system_management_functions.Add_to_path("C:\\MyTools")
    if err != nil {
        log.Fatal(err)
    }

    size, err := system_management_functions.Get_file_size_human_readable("C:\\Windows")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Windows folder size:", size)
}
```

---

## 📂 Repository Structure

```
go_functions_002/
├── CHANGELOG.md
├── CITATION.cff
├── LICENSE
├── README.md
└── v6
    ├── go.mod
    ├── go.sum
    ├── date_time_functions/
    │   └── date_time_functions.go
    ├── math_functions/
    │   └── math_functions.go
    ├── oracle_database_system_management_functions/
    │   └── oracle_database_system_management_functions.go
    ├── system_management_functions/
    │   └── system_management_functions.go
    └── yaml_functions/
        └── yaml_functions.go
```

---

## 📜 License

MIT License — see [LICENSE](LICENSE) for details.

---

## 📘 Citation

If you use this module in your work, please cite:

> Peter Cullen Burbery. (2025). go_functions_002 [Software]. GitHub. https://github.com/PeterCullenBurbery/go_functions_002