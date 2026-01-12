# Large alternation matching
# Input: log file with various log levels and keywords
# Measures: UseAhoCorasick optimization (coregex v0.9.0, >8 alternations)
# Pattern: 10+ alternations triggers Aho-Corasick multi-pattern matching
/ERROR|WARN|INFO|DEBUG|TRACE|FATAL|CRITICAL|NOTICE|ALERT|EMERGENCY/ { count++ }
END { print count }
