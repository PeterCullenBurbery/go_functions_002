# go_functions_002

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)  
[![Go Reference](https://pkg.go.dev/badge/github.com/PeterCullenBurbery/go_functions_002/v5.svg)](https://pkg.go.dev/github.com/PeterCullenBurbery/go_functions_002/v5)

**Author:** Peter Cullen Burbery  
**Language:** Go 1.24+  
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
├── CHANGELOG.md
├── LICENSE
├── README.md
└── v5
    ├── go.mod
    ├── go.sum
    ├── date_time_functions/
    │   └── date_time_functions.go
    ├── math_functions/
    │   └── math_functions.go
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