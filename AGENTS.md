# AGENTS.md

## Project Overview

TimeForge is a comprehensive Go CLI toolkit for working with timestamps, timezones, durations, and date ranges.

## Architecture

- `cmd/timeforge/` — CLI entry point using cobra
- `pkg/parse/` — Timestamp parsing (20+ formats, Unix timestamps, relative time)
- `pkg/convert/` — Timezone conversion, search, and info
- `pkg/duration/` — Duration arithmetic, formatting, and parsing
- `pkg/range_/` — Time range generation, splitting, and analysis
- `pkg/format/` — Timestamp formatting with 25+ named presets
- `pkg/analyze/` — Timestamp distribution analysis (percentiles, gaps, bursts)
- `pkg/business/` — Business day calendar and calculations

## Building

```bash
go build -o timeforge ./cmd/timeforge/
```

## Testing

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/parse/ -v

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Key Design Decisions

1. **Multi-format parsing** — Uses cascading format detection with 20+ layouts, plus Unix timestamp and relative time support
2. **Extended duration syntax** — Supports days, weeks, months, years beyond Go's standard duration
3. **Business calendar** — Configurable weekend/holiday system with US holiday defaults
4. **Analysis engine** — Computes percentiles, detects gaps, and identifies activity bursts
5. **Preset system** — 25+ named format presets for common use cases

## Common Tasks

- **Add a new timestamp format**: Add to `formats` slice in `pkg/parse/parse.go`
- **Add a new format preset**: Add to `presets` map in `pkg/format/format.go`
- **Add a new business holiday**: Use `cal.AddHoliday()` in `pkg/business/business.go`
- **Add a new CLI command**: Add a `cobra.Command` in `cmd/timeforge/main.go`
- **Add a new analysis metric**: Add to `AnalysisResult` struct and `Analyze()` function

## Dependencies

- Go 1.24+ standard library
- `github.com/spf13/cobra` — CLI framework

## Testing Conventions

- Table-driven tests for multiple scenarios
- Test both happy paths and error cases
- Use meaningful test names describing the scenario
