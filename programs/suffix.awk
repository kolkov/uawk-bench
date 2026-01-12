# Suffix pattern matching
# Input: log file with filenames
# Measures: reverse search optimization
# Pattern: .*\.(txt|log|md)
/\.(txt|log|md)$/ { count++ }
END { print count }
