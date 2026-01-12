# Character class matching
# Input: text file
# Measures: CharClassSearcher fast path
# Pattern: [a-zA-Z]+
/[a-zA-Z]+/ { count++ }
END { print count }
