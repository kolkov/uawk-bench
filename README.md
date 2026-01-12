# uawk-bench

Benchmark suite for comparing AWK implementations.

## Latest Results (uawk v0.2.1)

**uawk wins 16/16 benchmarks** vs GoAWK on Linux (10MB dataset, 10 runs).

[View CI Results](https://github.com/kolkov/uawk-bench/actions/runs/20934114950)

### Summary

| vs | Wins | Notable |
|----|------|---------|
| GoAWK | **16/16** | alternation 29x, email 14x, inner 13x |
| gawk | **13/16** | loses on ipaddr, regex, version |
| mawk | **10/16** | loses on anchored, charclass, inner, groupby, wordcount, suffix |

### Detailed Results (Linux CI)

| Benchmark | uawk | goawk | gawk | mawk |
|-----------|------|-------|------|------|
| alternation | **25ms** | 723ms | 34ms | 28ms |
| anchored | **17ms** | 22ms | 33ms | 8ms |
| charclass | **19ms** | 43ms | 30ms | 10ms |
| count | **39ms** | 65ms | 61ms | 43ms |
| csv | **69ms** | 102ms | 117ms | 92ms |
| email | **25ms** | 343ms | 49ms | 630ms |
| filter | **96ms** | 112ms | 125ms | 98ms |
| groupby | **196ms** | 277ms | 308ms | 146ms |
| inner | **23ms** | 301ms | 40ms | 12ms |
| ipaddr | 47ms | 138ms | **39ms** | 104ms |
| regex | 78ms | 250ms | **44ms** | 513ms |
| select | **77ms** | 97ms | 138ms | 73ms |
| suffix | **24ms** | 51ms | 33ms | 22ms |
| sum | **79ms** | 99ms | 118ms | 81ms |
| version | 45ms | 130ms | **37ms** | 97ms |
| wordcount | **212ms** | 241ms | 303ms | 165ms |

*Test environment: Ubuntu 24.04, GitHub Actions runner, Go 1.25.5*

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
