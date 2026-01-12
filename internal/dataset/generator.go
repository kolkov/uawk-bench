// Package dataset generates test data for AWK benchmarks.
package dataset

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Size represents dataset size presets.
type Size int

const (
	Small  Size = 1 << 20  // 1 MB
	Medium Size = 10 << 20 // 10 MB
	Large  Size = 100 << 20 // 100 MB
	XLarge Size = 500 << 20 // 500 MB
)

func (s Size) String() string {
	switch {
	case s < 1<<20:
		return fmt.Sprintf("%dKB", s/1024)
	case s < 1<<30:
		return fmt.Sprintf("%dMB", s/(1<<20))
	default:
		return fmt.Sprintf("%dGB", s/(1<<30))
	}
}

// Generator creates test datasets.
type Generator struct {
	Seed int64
	rng  *rand.Rand
}

// NewGenerator creates a generator with the given seed.
func NewGenerator(seed int64) *Generator {
	return &Generator{
		Seed: seed,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// GenerateNumeric creates a file with numeric data.
// Format: "int float int float int" per line
func (g *Generator) GenerateNumeric(dir string, size Size) (string, error) {
	filename := filepath.Join(dir, fmt.Sprintf("numeric_%s.txt", size))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	for written < int64(size) {
		line := fmt.Sprintf("%d %.6f %d %.6f %d\n",
			g.rng.Intn(1000),
			g.rng.Float64()*1000,
			g.rng.Intn(1000),
			g.rng.Float64()*1000,
			g.rng.Intn(1000),
		)
		n, err := w.WriteString(line)
		if err != nil {
			return "", err
		}
		written += int64(n)
	}

	if err := w.Flush(); err != nil {
		return "", err
	}
	return filename, nil
}

// GenerateText creates a file with text data (words).
func (g *Generator) GenerateText(dir string, size Size) (string, error) {
	filename := filepath.Join(dir, fmt.Sprintf("text_%s.txt", size))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	words := []string{
		"the", "quick", "brown", "fox", "jumps", "over", "lazy", "dog",
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing",
		"data", "processing", "benchmark", "performance", "test123", "value42",
		"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta",
	}

	for written < int64(size) {
		// Generate line with 5-15 words
		numWords := 5 + g.rng.Intn(11)
		lineWords := make([]string, numWords)
		for i := range lineWords {
			lineWords[i] = words[g.rng.Intn(len(words))]
		}
		line := strings.Join(lineWords, " ") + "\n"

		n, err := w.WriteString(line)
		if err != nil {
			return "", err
		}
		written += int64(n)
	}

	if err := w.Flush(); err != nil {
		return "", err
	}
	return filename, nil
}

// GenerateCSV creates a CSV file.
func (g *Generator) GenerateCSV(dir string, size Size) (string, error) {
	filename := filepath.Join(dir, fmt.Sprintf("data_%s.csv", size))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	// Header
	header := "id,name,value,category,score\n"
	n, _ := w.WriteString(header)
	written += int64(n)

	names := []string{"alice", "bob", "charlie", "david", "eve", "frank", "grace", "henry"}
	categories := []string{"A", "B", "C", "D"}

	id := 1
	for written < int64(size) {
		line := fmt.Sprintf("%d,%s,%.2f,%s,%d\n",
			id,
			names[g.rng.Intn(len(names))],
			g.rng.Float64()*1000,
			categories[g.rng.Intn(len(categories))],
			g.rng.Intn(100),
		)
		n, err := w.WriteString(line)
		if err != nil {
			return "", err
		}
		written += int64(n)
		id++
	}

	if err := w.Flush(); err != nil {
		return "", err
	}
	return filename, nil
}

// GenerateKeyValue creates a file for group-by benchmarks.
// Format: "key value" where key has limited cardinality
func (g *Generator) GenerateKeyValue(dir string, size Size) (string, error) {
	filename := filepath.Join(dir, fmt.Sprintf("keyvalue_%s.txt", size))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	// 100 unique keys for group-by
	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("key%03d", i)
	}

	for written < int64(size) {
		line := fmt.Sprintf("%s %d\n",
			keys[g.rng.Intn(len(keys))],
			g.rng.Intn(1000),
		)
		n, err := w.WriteString(line)
		if err != nil {
			return "", err
		}
		written += int64(n)
	}

	if err := w.Flush(); err != nil {
		return "", err
	}
	return filename, nil
}

// GenerateLog creates a log file with IP addresses and log levels.
// Format: "2024-01-05 10:30:45 192.168.1.100 INFO Processing request..."
// Used for ipaddr.awk and alternation.awk benchmarks.
func (g *Generator) GenerateLog(dir string, size Size) (string, error) {
	filename := filepath.Join(dir, fmt.Sprintf("log_%s.txt", size))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	var written int64

	levels := []string{"ERROR", "WARN", "INFO", "DEBUG", "TRACE", "FATAL", "CRITICAL", "NOTICE", "ALERT", "EMERGENCY"}
	messages := []string{
		"Processing request from client",
		"Connection established successfully",
		"Database query completed",
		"Cache miss for key",
		"Authentication failed for user",
		"File uploaded successfully",
		"Memory usage threshold exceeded",
		"Service health check passed",
		"Rate limit exceeded",
		"Session expired",
	}

	for written < int64(size) {
		// Generate random IP address
		ip := fmt.Sprintf("%d.%d.%d.%d",
			g.rng.Intn(256),
			g.rng.Intn(256),
			g.rng.Intn(256),
			g.rng.Intn(256),
		)

		// Random timestamp
		hour := g.rng.Intn(24)
		min := g.rng.Intn(60)
		sec := g.rng.Intn(60)

		line := fmt.Sprintf("2024-01-05 %02d:%02d:%02d %s %s %s\n",
			hour, min, sec,
			ip,
			levels[g.rng.Intn(len(levels))],
			messages[g.rng.Intn(len(messages))],
		)

		n, err := w.WriteString(line)
		if err != nil {
			return "", err
		}
		written += int64(n)
	}

	if err := w.Flush(); err != nil {
		return "", err
	}
	return filename, nil
}

// GenerateAll creates all dataset types for the given size.
func (g *Generator) GenerateAll(dir string, size Size) (map[string]string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	files := make(map[string]string)

	numeric, err := g.GenerateNumeric(dir, size)
	if err != nil {
		return nil, fmt.Errorf("numeric: %w", err)
	}
	files["numeric"] = numeric

	text, err := g.GenerateText(dir, size)
	if err != nil {
		return nil, fmt.Errorf("text: %w", err)
	}
	files["text"] = text

	csv, err := g.GenerateCSV(dir, size)
	if err != nil {
		return nil, fmt.Errorf("csv: %w", err)
	}
	files["csv"] = csv

	kv, err := g.GenerateKeyValue(dir, size)
	if err != nil {
		return nil, fmt.Errorf("keyvalue: %w", err)
	}
	files["keyvalue"] = kv

	log, err := g.GenerateLog(dir, size)
	if err != nil {
		return nil, fmt.Errorf("log: %w", err)
	}
	files["log"] = log

	return files, nil
}
