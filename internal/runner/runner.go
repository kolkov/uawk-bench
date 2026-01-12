// Package runner executes AWK implementations and measures performance.
package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AWK represents an AWK implementation to benchmark.
type AWK struct {
	Name    string   // Display name (e.g., "uawk", "goawk")
	Command string   // Executable name or path
	Args    []string // Additional arguments (e.g., ["-b"] for gawk)
}

// Result holds benchmark results for a single run.
type Result struct {
	AWK      string        // AWK implementation name
	Program  string        // AWK program name
	Duration time.Duration // Execution time
	Output   string        // Program output (for verification)
	Error    error         // Error if execution failed
}

// BenchmarkResult holds aggregated results for multiple runs.
type BenchmarkResult struct {
	AWK       string
	Program   string
	Runs      int
	Min       time.Duration
	Max       time.Duration
	Mean      time.Duration
	Median    time.Duration
	StdDev    time.Duration
	Throughput float64 // MB/s based on input size
}

// Runner executes AWK benchmarks.
type Runner struct {
	AWKs    []AWK
	Timeout time.Duration
	Warmup  int // Number of warmup runs
	Runs    int // Number of measured runs
}

// NewRunner creates a runner with default settings.
func NewRunner() *Runner {
	return &Runner{
		Timeout: 5 * time.Minute,
		Warmup:  1,
		Runs:    5,
	}
}

// DefaultAWKs returns the standard set of AWK implementations to test.
func DefaultAWKs() []AWK {
	return []AWK{
		{Name: "uawk", Command: "uawk"},                                    // POSIX mode (default)
		{Name: "uawk-fast", Command: "uawk", Args: []string{"--no-posix"}}, // Fast mode (no Longest)
		{Name: "goawk", Command: "goawk"},
		{Name: "gawk", Command: "gawk", Args: []string{"-b"}}, // -b disables multibyte
		{Name: "mawk", Command: "mawk"},
	}
}

// FindAvailable returns AWKs that are installed on the system.
func FindAvailable(awks []AWK) []AWK {
	var available []AWK
	for _, awk := range awks {
		// Try command as-is first
		if path, err := exec.LookPath(awk.Command); err == nil {
			awk.Command = path
			available = append(available, awk)
			continue
		}
		// Try with .exe extension on Windows
		if path, err := exec.LookPath(awk.Command + ".exe"); err == nil {
			awk.Command = path
			available = append(available, awk)
			continue
		}
		// Try in GOPATH/bin
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			goBin := filepath.Join(gopath, "bin", awk.Command)
			if _, err := os.Stat(goBin); err == nil {
				awk.Command = goBin
				available = append(available, awk)
				continue
			}
			goBin = filepath.Join(gopath, "bin", awk.Command+".exe")
			if _, err := os.Stat(goBin); err == nil {
				awk.Command = goBin
				available = append(available, awk)
				continue
			}
		}
		// Try in HOME/go/bin
		if home := os.Getenv("HOME"); home != "" {
			goBin := filepath.Join(home, "go", "bin", awk.Command)
			if _, err := os.Stat(goBin); err == nil {
				awk.Command = goBin
				available = append(available, awk)
				continue
			}
			goBin = filepath.Join(home, "go", "bin", awk.Command+".exe")
			if _, err := os.Stat(goBin); err == nil {
				awk.Command = goBin
				available = append(available, awk)
				continue
			}
		}
		// Try in USERPROFILE/go/bin (Windows)
		if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
			goBin := filepath.Join(userprofile, "go", "bin", awk.Command+".exe")
			if _, err := os.Stat(goBin); err == nil {
				awk.Command = goBin
				available = append(available, awk)
				continue
			}
		}
	}
	return available
}

// Run executes an AWK program with the given input file.
func (r *Runner) Run(ctx context.Context, awk AWK, programFile, inputFile string) Result {
	args := append([]string{}, awk.Args...)
	args = append(args, "-f", programFile, inputFile)

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, awk.Command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return Result{
			AWK:      awk.Name,
			Program:  programFile,
			Duration: duration,
			Error:    fmt.Errorf("%w: %s", err, stderr.String()),
		}
	}

	return Result{
		AWK:      awk.Name,
		Program:  programFile,
		Duration: duration,
		Output:   stdout.String(),
	}
}

// RunInline executes an inline AWK program.
func (r *Runner) RunInline(ctx context.Context, awk AWK, program, inputFile string) Result {
	args := append([]string{}, awk.Args...)
	args = append(args, program, inputFile)

	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, awk.Command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return Result{
			AWK:      awk.Name,
			Program:  program,
			Duration: duration,
			Error:    fmt.Errorf("%w: %s", err, stderr.String()),
		}
	}

	return Result{
		AWK:      awk.Name,
		Program:  program,
		Duration: duration,
		Output:   stdout.String(),
	}
}

// Benchmark runs multiple iterations and returns aggregated results.
func (r *Runner) Benchmark(ctx context.Context, awk AWK, programFile, inputFile string, inputSize int64) (*BenchmarkResult, error) {
	// Warmup runs
	for i := 0; i < r.Warmup; i++ {
		result := r.Run(ctx, awk, programFile, inputFile)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Measured runs
	durations := make([]time.Duration, r.Runs)
	for i := 0; i < r.Runs; i++ {
		result := r.Run(ctx, awk, programFile, inputFile)
		if result.Error != nil {
			return nil, result.Error
		}
		durations[i] = result.Duration
	}

	// Calculate statistics
	return calculateStats(awk.Name, programFile, durations, inputSize), nil
}

func calculateStats(awkName, program string, durations []time.Duration, inputSize int64) *BenchmarkResult {
	n := len(durations)
	if n == 0 {
		return nil
	}

	// Sort for median
	sorted := make([]time.Duration, n)
	copy(sorted, durations)
	sortDurations(sorted)

	// Calculate min, max, mean
	var sum time.Duration
	min := durations[0]
	max := durations[0]
	for _, d := range durations {
		sum += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	mean := sum / time.Duration(n)

	// Median
	var median time.Duration
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}

	// Standard deviation
	var variance float64
	meanFloat := float64(mean)
	for _, d := range durations {
		diff := float64(d) - meanFloat
		variance += diff * diff
	}
	variance /= float64(n)
	stdDev := time.Duration(sqrt(variance))

	// Throughput (MB/s)
	var throughput float64
	if mean > 0 && inputSize > 0 {
		throughput = float64(inputSize) / (1024 * 1024) / mean.Seconds()
	}

	return &BenchmarkResult{
		AWK:        awkName,
		Program:    program,
		Runs:       n,
		Min:        min,
		Max:        max,
		Mean:       mean,
		Median:     median,
		StdDev:     stdDev,
		Throughput: throughput,
	}
}

func sortDurations(d []time.Duration) {
	// Simple insertion sort (small arrays)
	for i := 1; i < len(d); i++ {
		key := d[i]
		j := i - 1
		for j >= 0 && d[j] > key {
			d[j+1] = d[j]
			j--
		}
		d[j+1] = key
	}
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}
