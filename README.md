# go_functions_002

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)  
[![Go Reference](https://pkg.go.dev/badge/github.com/PeterCullenBurbery/go_functions_002/v5.svg)](https://pkg.go.dev/github.com/PeterCullenBurbery/go_functions_002/v5)

**Author:** Peter Cullen Burbery  
**Language:** Go 1.21+  
**Platform:** Primarily Windows, with some cross-platform functions

---

## 📖 Overview

`go_functions_002` is a utility library written in Go, offering a wide range of functions for:

- **System management** (PATH updates, registry tweaks, Windows Explorer settings)
- **Date and time handling**
- **YAML parsing**
- **Mathematical algorithms**
- **Package management automation** (Chocolatey, Winget)
- **File and network utilities**

This package is designed for automation, system configuration, and utility scripting — especially for Windows system administration in Go.

---

## 📦 Installation

```bash
go get github.com/PeterCullenBurbery/go_functions_002/v5
```

---

## ✨ Features

### 📅 Date/Time Functions
- **`Date_time_stamp()`** – Returns a precise timestamp using a temporary Java helper (auto-installs Java if needed).
- **`Format_now()`** – Returns current time in `"2006-01-02 15:04:05"` format.
- **`Safe_time_stamp()`** – Produces a safe filename timestamp (replaces `/` with ` slash ` when needed).

---

### 🧮 Math Functions
- **`Topological_sort()`** – Deterministic Kahn’s algorithm, sorts nodes alphabetically when precedence is equal.
- **`Reverse_topological_sort()`** – Returns reversed topological order (useful for teardown sequences).

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
- **Registry-Based System Config**
  - `Set_24_hour_format`, `Do_not_use_24_hour_format`
  - `Set_short_date_pattern`, `Reset_short_date_pattern`
  - `Set_long_date_pattern`, `Reset_long_date_pattern`
  - `Set_time_pattern`, `Reset_time_pattern`
  - `Seconds_in_taskbar`, `Take_seconds_out_of_taskbar`
  - `Enable_long_file_paths`, `Are_long_file_paths_enabled`
- **Security & Installation Helpers**
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

## 🖥 Example Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/PeterCullenBurbery/go_functions_002/v5/system_management_functions"
    "github.com/PeterCullenBurbery/go_functions_002/v5/date_time_functions"
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
├── date_time_functions/
├── math_functions/
├── system_management_functions/
├── yaml_functions/
├── CHANGELOG.md
└── LICENSE
```

---

## 📜 License

MIT License — see [LICENSE](LICENSE) for details.

---

## 📘 Citation

If you use this module in your work, please cite:

> Peter Cullen Burbery. (2025). go_functions_002 [Software]. GitHub. https://github.com/PeterCullenBurbery/go_functions_002

---

Do you want me to also make you a **condensed features table** in the README so people can quickly scan all 60+ functions without reading full descriptions? That could make it look more like your Python and PowerShell README styles.

>> USER: Yes

>> ASSISTANT: Got it ✅  
Here’s your **updated README.md** for `go_functions_002` with a **condensed features table** section added so people can quickly scan all functions.

---

# go_functions_002

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)  
[![Go Reference](https://pkg.go.dev/badge/github.com/PeterCullenBurbery/go_functions_002.svg)](https://pkg.go.dev/github.com/PeterCullenBurbery/go_functions_002)

**Author:** Peter Cullen Burbery  
**Language:** Go 1.21+  
**Platform:** Primarily Windows, with some cross-platform functions

---

## 📖 Overview

`go_functions_002` is a utility library written in Go, offering a wide range of functions for:

- **System management** (PATH updates, registry tweaks, Windows Explorer settings)
- **Date and time handling**
- **YAML parsing**
- **Mathematical algorithms**
- **Package management automation** (Chocolatey, Winget)
- **File and network utilities**

This package is designed for automation, system configuration, and utility scripting — especially for Windows system administration in Go.

---

## 📦 Installation

```bash
go get github.com/PeterCullenBurbery/go_functions_002/v5
```

---

## ✨ Features

### 📅 Date/Time Functions
- **`Date_time_stamp()`** – Precise timestamp with Java helper (auto-installs Java if needed).
- **`Format_now()`** – Current time in `"2006-01-02 15:04:05"` format.
- **`Safe_time_stamp()`** – Safe filename timestamp (`/` → ` slash `).

---

### 🧮 Math Functions
- **`Topological_sort()`** – Deterministic Kahn’s algorithm with alphabetical tie-breaks.
- **`Reverse_topological_sort()`** – Reverse of topological order (leaf-to-root).

---

### ⚙️ System Management Functions
Includes over **50 Windows utilities** such as:
- **PATH & Environment:** `Add_to_path`, `Remove_from_path`, `Clean_path`, `Add_to_ps_module_path`, `Remove_from_ps_module_path`, `Set_system_environment_variable`
- **Explorer Tweaks:** `Show_hidden_files`, `Do_not_show_hidden_files`, `Show_file_extensions`, `Do_not_show_file_extensions`, `Bring_back_the_right_click_menu`, `Use_Windows_11_right_click_menu`, `Hide_search_box`, `Do_not_hide_search_box`, `Set_dark_mode`, `Set_light_mode`
- **Registry Config:** `Set_24_hour_format`, `Do_not_use_24_hour_format`, `Set_short_date_pattern`, `Reset_short_date_pattern`, `Set_long_date_pattern`, `Reset_long_date_pattern`, `Set_time_pattern`, `Reset_time_pattern`, `Seconds_in_taskbar`, `Take_seconds_out_of_taskbar`, `Enable_long_file_paths`, `Are_long_file_paths_enabled`
- **Security & Installation:** `Exclude_from_Microsoft_Windows_Defender`, `Choco_install`, `Winget_install`, `Install_choco`, `Enable_SSH`, `Enable_SSH_through_firewall`
- **File Utilities:** `Download_file`, `Extract_zip`, `Extract_password_protected_zip`, `File_exists`, `Get_file_size`, `Get_file_size_human_readable`, `Create_desktop_shortcut`
- **Networking:** `Get_primary_ipv4_address`

---

### 📄 YAML Functions
- `GetCaseInsensitiveString`, `GetCaseInsensitiveMap`, `GetCaseInsensitiveList`, `GetNestedString`, `GetNestedMap`

---

## 📊 Condensed Features Table

| Category | Function | Description |
|----------|----------|-------------|
| **Date/Time** | `Date_time_stamp()` | Precise timestamp with nanoseconds, week & ordinal date |
|  | `Format_now()` | Current time in `YYYY-MM-DD HH:MM:SS` |
|  | `Safe_time_stamp()` | Replaces `/` with ` slash ` for safe filenames |
| **Math** | `Topological_sort()` | Deterministic DAG topological sort |
|  | `Reverse_topological_sort()` | Reverse order of topological sort |
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

    "github.com/PeterCullenBurbery/go_functions_002/v5/system_management_functions"
    "github.com/PeterCullenBurbery/go_functions_002/v5/date_time_functions"
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
├── date_time_functions/
├── math_functions/
├── system_management_functions/
├── yaml_functions/
├── CHANGELOG.md
└── LICENSE
```

---

## 📜 License

MIT License — see [LICENSE](LICENSE) for details.

---

## 📘 Citation

If you use this module in your work, please cite:

> Peter Cullen Burbery. (2025). go_functions_002 [Software]. GitHub. https://github.com/PeterCullenBurbery/go_functions_002