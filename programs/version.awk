# Version number matching
# Input: log file
# Measures: digit sequences with dots
# Pattern: [0-9]+\.[0-9]+\.[0-9]+
/[0-9]+\.[0-9]+\.[0-9]+/ { count++ }
END { print count }
