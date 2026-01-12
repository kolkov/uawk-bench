# Email pattern matching
# Input: text file
# Measures: character class with special chars [\w.+-]
# Pattern: [\w.+-]+@[\w.-]+\.[\w.-]+
/[a-zA-Z0-9_.+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9.-]+/ { count++ }
END { print count }
