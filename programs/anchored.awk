# Anchored pattern matching
# Input: log file (HTTP logs)
# Measures: start anchor optimization
# Pattern: ^HTTP/[12]\.[01]
/^HTTP\/[12]\.[01]/ { count++ }
END { print count }
