//go:build ignore

// Script to generate test data for AWK benchmarks.
// Run with: go run scripts/generate-data.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const (
	outputDir = "testdata"
	seed      = 42
)

// Dataset sizes
var sizes = map[string]int64{
	"1MB":   1 << 20,
	"10MB":  10 << 20,
	"100MB": 100 << 20,
}

var (
	words = []string{
		"the", "quick", "brown", "fox", "jumps", "over", "lazy", "dog",
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing",
		"data", "processing", "benchmark", "performance", "test123", "value42",
		"alpha", "beta", "gamma", "delta", "epsilon", "server", "client",
	}

	names      = []string{"alice", "bob", "charlie", "david", "eve", "frank"}
	categories = []string{"A", "B", "C", "D"}
)

func main() {
	rng := rand.New(rand.NewSource(seed))

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	for sizeName, sizeBytes := range sizes {
		fmt.Printf("Generating %s datasets...\n", sizeName)

		generateNumeric(rng, sizeName, sizeBytes)
		generateText(rng, sizeName, sizeBytes)
		generateCSV(rng, sizeName, sizeBytes)
		generateKeyValue(rng, sizeName, sizeBytes)
	}

	fmt.Println("Done!")
}

func generateNumeric(rng *rand.Rand, sizeName string, size int64) {
	filename := filepath.Join(outputDir, fmt.Sprintf("numeric_%s.txt", sizeName))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	for written < size {
		line := fmt.Sprintf("%d %.6f %d %.6f %d\n",
			rng.Intn(1000),
			rng.Float64()*1000,
			rng.Intn(1000),
			rng.Float64()*1000,
			rng.Intn(1000),
		)
		n, _ := w.WriteString(line)
		written += int64(n)
	}
	w.Flush()
	fmt.Printf("  %s: %.1f MB\n", filename, float64(written)/(1<<20))
}

func generateText(rng *rand.Rand, sizeName string, size int64) {
	filename := filepath.Join(outputDir, fmt.Sprintf("text_%s.txt", sizeName))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	for written < size {
		numWords := 5 + rng.Intn(11)
		lineWords := make([]string, numWords)
		for i := range lineWords {
			lineWords[i] = words[rng.Intn(len(words))]
		}
		line := strings.Join(lineWords, " ") + "\n"
		n, _ := w.WriteString(line)
		written += int64(n)
	}
	w.Flush()
	fmt.Printf("  %s: %.1f MB\n", filename, float64(written)/(1<<20))
}

func generateCSV(rng *rand.Rand, sizeName string, size int64) {
	filename := filepath.Join(outputDir, fmt.Sprintf("data_%s.csv", sizeName))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	header := "id,name,value,category,score\n"
	n, _ := w.WriteString(header)
	written += int64(n)

	id := 1
	for written < size {
		line := fmt.Sprintf("%d,%s,%.2f,%s,%d\n",
			id,
			names[rng.Intn(len(names))],
			rng.Float64()*1000,
			categories[rng.Intn(len(categories))],
			rng.Intn(100),
		)
		n, _ := w.WriteString(line)
		written += int64(n)
		id++
	}
	w.Flush()
	fmt.Printf("  %s: %.1f MB\n", filename, float64(written)/(1<<20))
}

func generateKeyValue(rng *rand.Rand, sizeName string, size int64) {
	filename := filepath.Join(outputDir, fmt.Sprintf("keyvalue_%s.txt", sizeName))
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("key%03d", i)
	}

	for written < size {
		line := fmt.Sprintf("%s %d\n",
			keys[rng.Intn(len(keys))],
			rng.Intn(1000),
		)
		n, _ := w.WriteString(line)
		written += int64(n)
	}
	w.Flush()
	fmt.Printf("  %s: %.1f MB\n", filename, float64(written)/(1<<20))
}
