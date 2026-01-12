# Regex matching
# Input: text file
# Measures: regex engine performance
/[a-zA-Z]+[0-9]+/ { count++ }
END { print count }
