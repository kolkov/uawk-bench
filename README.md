# uawk-bench

Benchmark suite for comparing AWK implementations.

## Supported AWKs

| AWK | Description |
|-----|-------------|
| uawk | AWK interpreter in Go, uses coregex |
| goawk | Reference Go AWK by Ben Hoyt |
| gawk | GNU AWK |
| mawk | Fast C AWK |

## Benchmarks

| Program | Description | Data |
|---------|-------------|------|
| sum.awk | Sum numeric columns | numeric |
| count.awk | Count lines and fields | text |
| filter.awk | Filter by condition | numeric |
| select.awk | Extract specific fields | numeric |
| groupby.awk | Group and aggregate | keyvalue |
| wordcount.awk | Word frequency | text |
| regex.awk | Pattern `[a-zA-Z]+[0-9]+` | text |
| csv.awk | CSV field sum | csv |
| ipaddr.awk | IP address matching | log |
| alternation.awk | Log level matching | log |
| email.awk | Email pattern matching | text |
| suffix.awk | Suffix pattern matching | log |
| version.awk | Version number matching | log |
| charclass.awk | Character class patterns | text |
| inner.awk | Inner literal patterns | log |
| anchored.awk | Anchored patterns | log |

## Usage

```bash
# Build
go build -o bin/awkbench ./cmd/awkbench

# Run benchmarks
./bin/awkbench -size 10MB -runs 10

# Generate data only
./bin/awkbench -generate -size 100MB

# Test specific AWKs
./bin/awkbench -awk uawk,goawk -runs 5

# Custom directories
./bin/awkbench -data ./testdata -output ./results
```

## uawk Modes

```bash
# Default (POSIX regex)
./bin/awkbench -awk uawk

# Fast mode (non-POSIX regex)
./bin/awkbench -awk uawk-fast

# Parallel execution (v0.2.0+)
./bin/awkbench -awk uawk-j4    # 4 workers
./bin/awkbench -awk uawk-j8    # 8 workers
```

## Output

Results are written to `results/` directory:
- `results.md` — Markdown
- `results.json` — JSON
- `results.csv` — CSV

## CI

Benchmarks run on push to main and weekly. Results are in GitHub Actions artifacts.

```bash
gh workflow run benchmark.yml -f size=100MB
```

## Requirements

- Go 1.25+
- uawk, goawk (go install)
- gawk, mawk (system packages)

## License

MIT
