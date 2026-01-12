# Changelog

All notable changes to the uawk benchmark suite will be documented in this file.

## [2026-01-12] uawk v0.2.1

### Benchmark Results (Linux CI, 10MB, 10 runs)

| Benchmark | uawk | uawk-j4 | goawk | gawk | mawk |
|-----------|------|---------|-------|------|------|
| alternation | **23ms** | 17ms | 715ms | 33ms | 28ms |
| anchored | **16ms** | 15ms | 21ms | 32ms | 7ms |
| charclass | **18ms** | 15ms | 41ms | 29ms | 9ms |
| count | **36ms** | 24ms | 61ms | 61ms | 42ms |
| csv | **67ms** | 44ms | 96ms | 118ms | 90ms |
| email | **24ms** | 23ms | 339ms | 49ms | 629ms |
| filter | **86ms** | 53ms | 108ms | 114ms | 89ms |
| groupby | **195ms** | 107ms | 269ms | 312ms | 145ms |
| inner | **22ms** | 22ms | 299ms | 40ms | 12ms |
| ipaddr | 45ms | 32ms | 136ms | **38ms** | 103ms |
| regex | 76ms | 45ms | 247ms | **44ms** | 452ms |
| select | **70ms** | 44ms | 92ms | 128ms | 66ms |
| suffix | **24ms** | 20ms | 50ms | 32ms | 21ms |
| sum | **77ms** | 49ms | 99ms | 117ms | 80ms |
| version | 44ms | 30ms | 128ms | **37ms** | 96ms |
| wordcount | **210ms** | 110ms | 236ms | 298ms | 166ms |

### Summary

| Comparison | Score | Notable |
|------------|-------|---------|
| uawk vs GoAWK | **16/16** | alternation 31x, email 14x, inner 14x |
| uawk vs gawk | **13/16** | loses: ipaddr, regex, version |
| uawk vs mawk | **10/16** | loses: anchored, charclass, inner, groupby, wordcount, suffix |

### Parallel Mode (-j4) Speedups

| Benchmark | Sequential | Parallel | Improvement |
|-----------|------------|----------|-------------|
| wordcount | 210ms | 110ms | **-48%** |
| groupby | 195ms | 107ms | **-45%** |
| filter | 86ms | 53ms | **-38%** |
| csv | 67ms | 44ms | **-34%** |
| count | 36ms | 24ms | **-33%** |

### Changes in v0.2.1

- **Lazy ENVIRON loading**: -56% VM creation time, -43% memory
- No regressions from v0.2.0

---

## [2026-01-12] uawk v0.2.0

### Changes

- **Parallel execution** (`-j N` flag) for multi-file processing
- **Global array opcodes** for 7% improvement on array-heavy benchmarks
- **CompositeSearcher fix** for overlapping patterns (email benchmark was broken)

### Notable Improvements vs v0.1.6

| Benchmark | v0.1.6 | v0.2.0 | Change |
|-----------|--------|--------|--------|
| email | broken | 24ms | **fixed** |
| wordcount | 236ms | 210ms | **-11%** |

---

## [2026-01-12] Initial Release

- 16 benchmark programs covering compute, I/O, and regex patterns
- Support for uawk, goawk, gawk, mawk
- GitHub Actions CI with weekly runs
- Multiple output formats: Markdown, JSON, CSV
