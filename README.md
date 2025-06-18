# gopher-tidy

[![Go Report Card](https://goreportcard.com/badge/github.com/CatGrepFind/gopher-tidy)](https://goreportcard.com/report/github.com/CatGrepFind/gopher-tidy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

My personal interactive command-line tool written in Go to find and remove leftover application files on macOS.

`gopher-tidy` helps you reclaim disk space and keep your system clean by hunting down residual files from uninstalled applications that often get left behind in hidden Library folders.

---

## Funcionallity

* **Interactive Interface**: Guides you through the cleanup process step-by-step.
* **Targeted Search**: Scans all common macOS locations for application remnants (`Application Support`, `Caches`, `Preferences`, etc.).
* **Safety First**: **Never deletes anything without your explicit confirmation.**
* **Secure by Design**: Does not require running as `sudo`. If a file needs admin privileges to be deleted, `gopher-tidy` will provide the exact, safe `sudo` command for you to run manually.
* **Pure Go**: Written in 100% standard library Go, making it lightweight and dependency-free.

## Demo

Here's a quick look at the `gopher-tidy` workflow:

```bash
$ ./gopher-tidy
--- macOS Application Cleaner (Go Edition) ---
Enter the name of the application to clean up (e.g., Docker): OldVPN

🔍 Searching for files related to 'OldVPN'...

Found potential leftover files and folders:
  [1] /Users/youruser/Library/Application Support/OldVPN
  [2] /Users/youruser/Library/Caches/com.company.oldvpn.plist
  [3] /Library/LaunchDaemons/com.company.oldvpn.helper.plist

Enter numbers to delete (e.g., 1 3 4), 'all' to delete everything, or 'quit' to exit.
> 1 3

--- DELETION SUMMARY ---
  - /Users/youruser/Library/Application Support/OldVPN
  - /Library/LaunchDaemons/com.company.oldvpn.helper.plist
Proceed with deleting these items? [y/N]: y

🚀 Starting deletion...
Attempting to delete: /Users/youruser/Library/Application Support/OldVPN
✅ DELETED: /Users/youruser/Library/Application Support/OldVPN
Attempting to delete: /Library/LaunchDaemons/com.company.oldvpn.helper.plist
❌ PERMISSION DENIED for: /Library/LaunchDaemons/com.company.oldvpn.helper.plist
   This file requires administrator privileges. To delete it, run this command in your terminal:
   sudo rm -rf "/Library/LaunchDaemons/com.company.oldvpn.helper.plist"


✨ Cleanup process finished.
