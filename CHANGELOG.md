# Changelog

All notable changes to this project will be documented in this file.

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