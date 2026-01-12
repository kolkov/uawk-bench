// Package report generates benchmark reports in various formats.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/kolkov/uawk-bench/internal/runner"
)

// Report holds all benchmark results.
type Report struct {
	Generated time.Time
	System    SystemInfo
	Results   []runner.BenchmarkResult
}

// SystemInfo describes the benchmark environment.
type SystemInfo struct {
	OS       string
	Arch     string
	CPUs     int
	GoVersion string
}

// WriteMarkdown writes results as a Markdown table.
func WriteMarkdown(w io.Writer, results []runner.BenchmarkResult) error {
	if len(results) == 0 {
		return nil
	}

	// Group by program
	byProgram := make(map[string][]runner.BenchmarkResult)
	for _, r := range results {
		byProgram[r.Program] = append(byProgram[r.Program], r)
	}

	// Get sorted program names
	programs := make([]string, 0, len(byProgram))
	for p := range byProgram {
		programs = append(programs, p)
	}
	sort.Strings(programs)

	fmt.Fprintf(w, "# AWK Benchmark Results\n\n")
	fmt.Fprintf(w, "Generated: %s\n\n", time.Now().Format(time.RFC3339))

	for _, prog := range programs {
		progResults := byProgram[prog]

		// Sort by mean time (fastest first)
		sort.Slice(progResults, func(i, j int) bool {
			return progResults[i].Mean < progResults[j].Mean
		})

		fmt.Fprintf(w, "## %s\n\n", prog)
		fmt.Fprintf(w, "| AWK | Mean | Min | Max | StdDev | Throughput |\n")
		fmt.Fprintf(w, "|-----|------|-----|-----|--------|------------|\n")

		baseline := progResults[0].Mean
		for _, r := range progResults {
			speedup := ""
			if r.Mean != baseline {
				ratio := float64(r.Mean) / float64(baseline)
				speedup = fmt.Sprintf(" (%.2fx)", ratio)
			}

			fmt.Fprintf(w, "| %s | %s%s | %s | %s | %s | %.1f MB/s |\n",
				r.AWK,
				formatDuration(r.Mean), speedup,
				formatDuration(r.Min),
				formatDuration(r.Max),
				formatDuration(r.StdDev),
				r.Throughput,
			)
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}

// WriteJSON writes results as JSON.
func WriteJSON(w io.Writer, results []runner.BenchmarkResult) error {
	report := struct {
		Generated string                    `json:"generated"`
		Results   []runner.BenchmarkResult `json:"results"`
	}{
		Generated: time.Now().Format(time.RFC3339),
		Results:   results,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

// WriteCSV writes results as CSV.
func WriteCSV(w io.Writer, results []runner.BenchmarkResult) error {
	fmt.Fprintf(w, "awk,program,runs,mean_ns,min_ns,max_ns,stddev_ns,throughput_mbps\n")
	for _, r := range results {
		fmt.Fprintf(w, "%s,%s,%d,%d,%d,%d,%d,%.2f\n",
			r.AWK,
			r.Program,
			r.Runs,
			r.Mean.Nanoseconds(),
			r.Min.Nanoseconds(),
			r.Max.Nanoseconds(),
			r.StdDev.Nanoseconds(),
			r.Throughput,
		)
	}
	return nil
}

// WriteSummary writes a brief summary comparing AWK implementations.
func WriteSummary(w io.Writer, results []runner.BenchmarkResult) error {
	if len(results) == 0 {
		return nil
	}

	// Aggregate by AWK
	byAWK := make(map[string][]time.Duration)
	for _, r := range results {
		byAWK[r.AWK] = append(byAWK[r.AWK], r.Mean)
	}

	// Calculate geometric mean for each AWK
	type awkScore struct {
		name  string
		score float64
	}
	scores := make([]awkScore, 0, len(byAWK))

	for awk, durations := range byAWK {
		// Geometric mean
		product := 1.0
		for _, d := range durations {
			product *= float64(d.Nanoseconds())
		}
		geoMean := pow(product, 1.0/float64(len(durations)))
		scores = append(scores, awkScore{awk, geoMean})
	}

	// Sort by score (fastest first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score < scores[j].score
	})

	fmt.Fprintf(w, "## Summary (Geometric Mean)\n\n")
	fmt.Fprintf(w, "| Rank | AWK | Relative Speed |\n")
	fmt.Fprintf(w, "|------|-----|----------------|\n")

	baseline := scores[0].score
	for i, s := range scores {
		ratio := s.score / baseline
		bar := strings.Repeat("█", int(10/ratio))
		fmt.Fprintf(w, "| %d | %s | %.2fx %s |\n", i+1, s.name, ratio, bar)
	}
	fmt.Fprintf(w, "\n")

	return nil
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%.0fns", float64(d.Nanoseconds()))
	case d < time.Millisecond:
		return fmt.Sprintf("%.1fµs", float64(d.Microseconds()))
	case d < time.Second:
		return fmt.Sprintf("%.1fms", float64(d.Milliseconds()))
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func pow(x, y float64) float64 {
	return math.Pow(x, y)
}
