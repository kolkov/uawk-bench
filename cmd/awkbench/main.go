// awkbench - AWK implementations benchmark tool
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kolkov/uawk-bench/internal/dataset"
	"github.com/kolkov/uawk-bench/internal/report"
	"github.com/kolkov/uawk-bench/internal/runner"
)

var (
	dataDir     = flag.String("data", "testdata", "Directory for test data")
	programDir  = flag.String("programs", "programs", "Directory with AWK programs")
	outputDir   = flag.String("output", "results", "Directory for results")
	size        = flag.String("size", "10MB", "Dataset size: 1MB, 10MB, 100MB, 500MB")
	runs        = flag.Int("runs", 5, "Number of benchmark runs")
	warmup      = flag.Int("warmup", 1, "Number of warmup runs")
	awkList     = flag.String("awk", "", "Comma-separated list of AWKs to test (default: all available)")
	generateOnly = flag.Bool("generate", false, "Only generate test data, don't benchmark")
	format      = flag.String("format", "markdown", "Output format: markdown, json, csv")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse size
	datasetSize := parseSize(*size)
	if datasetSize == 0 {
		return fmt.Errorf("invalid size: %s (use 1MB, 10MB, 100MB, 500MB)", *size)
	}

	// Generate test data
	fmt.Printf("Generating test data (%s)...\n", *size)
	gen := dataset.NewGenerator(42) // Fixed seed for reproducibility
	files, err := gen.GenerateAll(*dataDir, datasetSize)
	if err != nil {
		return fmt.Errorf("generating data: %w", err)
	}
	fmt.Printf("Generated %d datasets in %s\n", len(files), *dataDir)

	if *generateOnly {
		for name, path := range files {
			info, _ := os.Stat(path)
			fmt.Printf("  %s: %s (%.1f MB)\n", name, path, float64(info.Size())/(1<<20))
		}
		return nil
	}

	// Find available AWKs
	// Note: frawk skipped - Cranelift backend crashes, LLVM requires complex CI setup
	numCPU := runtime.NumCPU()
	allAWKs := []runner.AWK{
		{Name: "uawk", Command: "uawk"},                                                          // POSIX mode (default)
		{Name: "uawk-fast", Command: "uawk", Args: []string{"--no-posix"}},                       // Fast mode
		{Name: "uawk-j4", Command: "uawk", Args: []string{"-j", "4"}},                            // Parallel 4 workers
		{Name: fmt.Sprintf("uawk-j%d", numCPU), Command: "uawk", Args: []string{"-j", fmt.Sprintf("%d", numCPU)}}, // Parallel max workers
		{Name: "goawk", Command: "goawk"},
		{Name: "gawk", Command: "gawk", Args: []string{"-b"}},
		{Name: "mawk", Command: "mawk"},
	}

	var awks []runner.AWK
	if *awkList != "" {
		// Filter to requested AWKs
		requested := strings.Split(*awkList, ",")
		for _, awk := range allAWKs {
			for _, req := range requested {
				if strings.TrimSpace(req) == awk.Name {
					awks = append(awks, awk)
					break
				}
			}
		}
	} else {
		awks = runner.FindAvailable(allAWKs)
	}

	if len(awks) == 0 {
		return fmt.Errorf("no AWK implementations found")
	}

	fmt.Printf("Testing AWK implementations: ")
	for i, awk := range awks {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(awk.Name)
	}
	fmt.Println()

	// Find AWK programs
	programs, err := filepath.Glob(filepath.Join(*programDir, "*.awk"))
	if err != nil {
		return fmt.Errorf("finding programs: %w", err)
	}
	if len(programs) == 0 {
		return fmt.Errorf("no AWK programs found in %s", *programDir)
	}

	fmt.Printf("Running %d programs with %d runs each...\n\n", len(programs), *runs)

	// Setup runner
	r := runner.NewRunner()
	r.Runs = *runs
	r.Warmup = *warmup

	ctx := context.Background()
	var results []runner.BenchmarkResult

	// Map programs to appropriate data files
	programData := map[string]string{
		"sum.awk":         files["numeric"],
		"count.awk":       files["text"],
		"filter.awk":      files["numeric"],
		"select.awk":      files["numeric"],
		"groupby.awk":     files["keyvalue"],
		"wordcount.awk":   files["text"],
		"regex.awk":       files["text"],
		"csv.awk":         files["csv"],
		"ipaddr.awk":      files["log"],      // coregex: DigitPrefilter
		"alternation.awk": files["log"],      // coregex: Aho-Corasick
		"email.awk":       files["text"],     // coregex: char class with special chars
		"suffix.awk":      files["log"],      // coregex: reverse search
		"version.awk":     files["log"],      // coregex: digit sequences
		"charclass.awk":   files["text"],     // uawk: CharClassSearcher fast path
		"inner.awk":       files["log"],      // coregex: inner literal optimization
		"anchored.awk":    files["log"],      // coregex: start anchor
	}

	// Run benchmarks
	for _, prog := range programs {
		progName := filepath.Base(prog)
		dataFile, ok := programData[progName]
		if !ok {
			dataFile = files["numeric"] // Default
		}

		info, _ := os.Stat(dataFile)
		inputSize := info.Size()

		fmt.Printf("%-20s ", progName)

		for _, awk := range awks {
			result, err := r.Benchmark(ctx, awk, prog, dataFile, inputSize)
			if err != nil {
				fmt.Printf("%s:ERR ", awk.Name)
				fmt.Fprintf(os.Stderr, "  [%s error: %v]\n", awk.Name, err)
				continue
			}
			results = append(results, *result)
			fmt.Printf("%s:%.1fms ", awk.Name, float64(result.Mean.Milliseconds()))
		}
		fmt.Println()
	}

	// Write results
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Write all formats
	mdFile := filepath.Join(*outputDir, "results.md")
	jsonFile := filepath.Join(*outputDir, "results.json")
	csvFile := filepath.Join(*outputDir, "results.csv")

	// Markdown
	f, err := os.Create(mdFile)
	if err != nil {
		return err
	}
	report.WriteSummary(f, results)
	report.WriteMarkdown(f, results)
	writeSystemInfo(f)
	f.Close()
	fmt.Printf("\nResults written to %s\n", mdFile)

	// JSON
	f, err = os.Create(jsonFile)
	if err != nil {
		return err
	}
	report.WriteJSON(f, results)
	f.Close()

	// CSV
	f, err = os.Create(csvFile)
	if err != nil {
		return err
	}
	report.WriteCSV(f, results)
	f.Close()

	return nil
}

func parseSize(s string) dataset.Size {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "1MB":
		return dataset.Small
	case "10MB":
		return dataset.Medium
	case "100MB":
		return dataset.Large
	case "500MB":
		return dataset.XLarge
	default:
		return 0
	}
}

func writeSystemInfo(f *os.File) {
	fmt.Fprintf(f, "## System Info\n\n")
	fmt.Fprintf(f, "- OS: %s\n", runtime.GOOS)
	fmt.Fprintf(f, "- Arch: %s\n", runtime.GOARCH)
	fmt.Fprintf(f, "- CPUs: %d\n", runtime.NumCPU())
	fmt.Fprintf(f, "- Go: %s\n", runtime.Version())
}
