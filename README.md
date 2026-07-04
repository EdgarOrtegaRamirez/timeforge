# TimeForge

A comprehensive Go CLI toolkit for working with timestamps, timezones, durations, and date ranges.

## Features

- **Parse** — Parse timestamps in 20+ formats (RFC3339, ISO8601, Unix, custom layouts)
- **Convert** — Convert between timezones with full timezone database support
- **Duration** — Duration arithmetic, formatting (compact/human), and parsing with extended units (days, weeks, months)
- **Range** — Generate, split, and analyze time ranges with breakdown calculations
- **Format** — Format timestamps using 25+ named presets or custom Go layouts
- **Analyze** — Analyze timestamp distributions: percentiles, gaps, bursts, histograms
- **Business** — Business day calculations: add/subtract, count between dates, list sequences

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/timeforge/cmd/timeforge@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/timeforge
cd timeforge
go build -o timeforge ./cmd/timeforge/
```

## Quick Start

```bash
# Parse a timestamp
timeforge parse "2024-01-15T10:30:00Z"

# Convert between timezones
timeforge convert "2024-01-15T10:30:00Z" --to US/Eastern

# Compute duration between two timestamps
timeforge duration compute "2024-01-15T10:00:00Z" "2024-01-16T14:30:00Z"

# Parse extended durations
timeforge duration parse "2d12h30m"

# Add business days
timeforge business add "2024-01-15" 5

# Count business days between dates
timeforge business between "2024-01-15" "2024-02-15"

# Generate time ranges
timeforge range generate "2024-01-15T00:00:00Z" 1h 24

# Split a range into intervals
timeforge range split "2024-01-01" "2024-12-31" 1M

# Format with a preset
timeforge format "2024-01-15T10:30:00Z" datetime

# List available format presets
timeforge format list

# Analyze timestamps from stdin
echo -e "2024-01-15T10:00:00Z\n2024-01-15T10:01:00Z\n2024-01-15T10:02:00Z" | timeforge analyze pipe
```

## Commands

### `timeforge parse [timestamp]`

Parse a timestamp in any supported format. Detects format automatically.

```
timeforge parse "2024-01-15T10:30:00Z"
Parsed:   2024-01-15T10:30:00Z
Format:   2006-01-02T15:04:05Z07:00
Unix:     1705312200
Unix MS:  1705312200000
UTC:      2024-01-15T10:30:00Z
Local:    2024-01-15T10:30:00Z
```

### `timeforge convert [timestamp]`

Convert timestamps between timezones.

```
timeforge convert "2024-01-15T10:30:00Z" --from UTC --to US/Eastern
2024-01-15T05:30:00-05:00
```

### `timeforge duration`

Duration arithmetic and formatting.

```bash
# Compute duration between timestamps
timeforge duration compute "2024-01-15T10:00:00Z" "2024-01-16T14:30:00Z"

# Parse extended duration strings
timeforge duration parse "2d12h30m"

# Add multiple durations
timeforge duration add 1h30m 45m 30s
```

### `timeforge range`

Generate and analyze time ranges.

```bash
# Generate hourly ranges for 24 hours
timeforge range generate "2024-01-15T00:00:00Z" 1h 24

# Split a year into monthly intervals
timeforge range split "2024-01-01" "2024-12-31" 1M
```

### `timeforge format [timestamp] [layout]`

Format timestamps using named presets or custom layouts.

```bash
# Use a preset
timeforge format "2024-01-15T10:30:00Z" datetime

# List all presets
timeforge format list

# Custom Go layout
timeforge format "2024-01-15T10:30:00Z" "Mon, 02 Jan 2006"
```

### `timeforge analyze`

Analyze timestamp distributions.

```bash
# Analyze timestamps from arguments
timeforge analyze timestamps 2024-01-15T10:00:00Z 2024-01-15T10:01:00Z 2024-01-15T10:05:00Z

# Analyze from stdin
cat timestamps.txt | timeforge analyze pipe
```

### `timeforge business`

Business day calculations.

```bash
# Add business days
timeforge business add "2024-01-15" 10

# Count between dates
timeforge business between "2024-01-15" "2024-03-15"

# Check if a date is a business day
timeforge business check "2024-01-15"

# List next N business days
timeforge business list "2024-01-15" 5
```

## Architecture

```
timeforge/
├── cmd/timeforge/     CLI entry point (cobra)
├── pkg/
│   ├── parse/         Timestamp parsing (20+ formats, Unix, relative)
│   ├── convert/       Timezone conversion and search
│   ├── duration/      Duration arithmetic and formatting
│   ├── range_/        Time range generation and analysis
│   ├── format/        Timestamp formatting with 25+ presets
│   ├── analyze/       Timestamp distribution analysis
│   └── business/      Business day calendar and calculations
└── tests/             Integration tests
```

## Supported Timestamp Formats

- RFC3339 / RFC3339Nano
- ISO 8601 variants
- Unix timestamps (seconds, milliseconds, nanoseconds)
- Date-only: `2004-01-15`, `01/15/2024`, `Jan 15, 2024`
- Date-time: `2024-01-15 10:30:00`, `15 Jan 2024 10:30:00`
- Log formats: `2024/01/15 10:30:00`
- Relative: `-2h30m`, `+1d`

## Duration Format

Supports extended duration syntax beyond Go's standard:

| Unit | Aliases |
|------|---------|
| Years | `y`, `yr`, `year`, `years` |
| Months | `M`, `mo`, `month`, `months` |
| Weeks | `w`, `wk`, `week`, `weeks` |
| Days | `d`, `day`, `days` |
| Hours | `h`, `hr`, `hour`, `hours` |
| Minutes | `m`, `min`, `minute`, `minutes` |
| Seconds | `s`, `sec`, `second`, `seconds` |
| Milliseconds | `ms`, `milli`, `millis` |

## License

MIT
