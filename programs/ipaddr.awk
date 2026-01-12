# IP address matching
# Input: log file with IP addresses
# Measures: DigitPrefilter optimization (coregex v0.9.0)
# Pattern: \d+\.\d+\.\d+\.\d+
/[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/ { count++ }
END { print count }
