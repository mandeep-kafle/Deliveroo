# Cron Expression Parser

A command-line application written in Go that parses cron expressions and displays their execution schedule. This parser handles the standard cron format with five time fields plus a command.

## Requirements

* Linux/macOS environment
* [Go 1.19](https://golang.org/doc/install) or higher (In MacOS: Binary can run without installing go if needed)



## Features Supported

* **Time Fields**: Minutes, Hours, Day of Month, Month, Day of Week
* **Patterns**:
  * Numbers (e.g., "5")
  * Lists (e.g., "1,2,3")
  * Ranges (e.g., "1-5")
  * Steps/Intervals (e.g., "*/15", "1-5/2")
  * Wildcards ("*")

* **Field Boundaries**:
  * Minutes: 0-59
  * Hours: 0-23
  * Day of Month: 1-31
  * Month: 1-12
  * Day of Week: 0-6  


## Development Setup

1. Clone the repository
```bash
git clone 
cd Deliveroo
```
2. Run with your input

```bash
Method 1 (With Go Installed): go run main.go "*/15 0 1,15 * 1-5 /usr/bin/find"
Method 2 (Without Go Run MacOS Binary) : sudo ./cronParserBuild "*/15 0 1,15 * 1-5 /usr/bin/find"
```
3. Run tests (Requires Go)
```bash
go test -v ./cronparser
```

## Current Limitations

* No support for special time strings (e.g., "@yearly")
* No support for special characters ('?' or 'L')
* No validation for varying days in months (e.g., February)

