# Inner literal matching
# Input: log file
# Measures: inner literal optimization (bidirectional search)
# Pattern: .*error.*
/.*error.*/ { count++ }
END { print count }
