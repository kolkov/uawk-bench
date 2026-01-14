# uawk-bench

Benchmark suite for comparing AWK implementations.

## Latest Results (uawk v0.2.2)

**uawk wins 16/16 benchmarks** vs GoAWK on Linux (10MB dataset, 10 runs).

[View CI Results](https://github.com/kolkov/uawk-bench/actions/runs/20991974572)

### Summary

| vs | Wins | Notable |
|----|------|---------|
| GoAWK | **16/16** | alternation 31x, email 14x, inner 14x |
| gawk | **13/16** | loses on ipaddr, regex, version |
| mawk | **10/16** | loses on anchored, charclass, inner, groupby, wordcount, suffix |

### Detailed Results (Linux CI)

| Benchmark | uawk | uawk-j4 | goawk | gawk | mawk |
|-----------|------|---------|-------|------|------|
| alternation | **23ms** | 17ms | 719ms | 33ms | 28ms |
| anchored | **15ms** | 15ms | 21ms | 32ms | 7ms |
| charclass | **17ms** | 16ms | 42ms | 29ms | 9ms |
| count | **37ms** | 24ms | 61ms | 61ms | 43ms |
| csv | **67ms** | 42ms | 96ms | 117ms | 90ms |
| email | **23ms** | 22ms | 340ms | 48ms | 631ms |
| filter | **83ms** | 48ms | 108ms | 117ms | 88ms |
| groupby | **200ms** | 108ms | 269ms | 311ms | 145ms |
| inner | **22ms** | 22ms | 298ms | 42ms | 12ms |
| ipaddr | 46ms | 34ms | 136ms | **39ms** | 103ms |
| regex | 80ms | 49ms | 248ms | **44ms** | 455ms |
| select | **68ms** | 44ms | 92ms | 127ms | 66ms |
| suffix | **23ms** | 21ms | 50ms | 32ms | 21ms |
| sum | **74ms** | 46ms | 97ms | 118ms | 80ms |
| version | 44ms | 33ms | 128ms | **37ms** | 96ms |
| wordcount | **211ms** | 108ms | 236ms | 303ms | 165ms |

### Parallel Mode Speedups (-j4)

| Benchmark | Sequential | Parallel | Improvement |
|-----------|------------|----------|-------------|
| wordcount | 211ms | 108ms | **-49%** |
| groupby | 200ms | 108ms | **-46%** |
| filter | 83ms | 48ms | **-42%** |
| regex | 80ms | 49ms | **-39%** |

*Test environment: Ubuntu 24.04, GitHub Actions runner, Go 1.25.5*

See [CHANGELOG.md](CHANGELOG.md) for version history.

## Supported AWKs

| AWK | Description |
|-----|-------------|
| uawk | High-performance AWK in Go with coregex regex engine |
| goawk | Reference Go AWK by Ben Hoyt |
| gawk | GNU AWK |
| mawk | Fast C AWK (Linux only) |

## Benchmarks

| Program | Description | Data | Pattern Type |
|---------|-------------|------|--------------|
| sum.awk | Sum numeric columns | numeric | - |
| count.awk | Count lines and fields | text | - |
| filter.awk | Filter by condition | numeric | - |
| select.awk | Extract specific fields | numeric | - |
| groupby.awk | Group and aggregate | keyvalue | - |
| wordcount.awk | Word frequency | text | - |
| regex.awk | Pattern `[a-zA-Z]+[0-9]+` | text | Composite |
| csv.awk | CSV field sum | csv | - |
| ipaddr.awk | IP address matching | log | DigitPrefilter |
| alternation.awk | Log level matching | log | Aho-Corasick |
| email.awk | Email pattern matching | text | CharClass |
| suffix.awk | Suffix pattern matching | log | Reverse search |
| version.awk | Version number matching | log | Digit sequences |
| charclass.awk | Character class patterns | text | CharClass |
| inner.awk | Inner literal patterns | log | Inner literal |
| anchored.awk | Anchored patterns | log | Start anchor |

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

Note: Parallel mode (`-j N`) requires multiple input files to show benefit. Single-file benchmarks run sequentially.

## Output

Results are written to `results/` directory:
- `results.md` — Markdown with detailed statistics
- `results.json` — JSON for programmatic analysis
- `results.csv` — CSV for spreadsheets

## CI

Benchmarks run on push to main and weekly. Results are in GitHub Actions artifacts.

```bash
gh workflow run benchmark.yml -f size=100MB
```

## Requirements

- Go 1.25+
- uawk, goawk (`go install`)
- gawk, mawk (system packages)

## License

MIT
