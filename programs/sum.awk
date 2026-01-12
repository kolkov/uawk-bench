# Sum numeric columns
# Input: numeric data with whitespace-separated fields
# Measures: field parsing + numeric operations
{ sum1 += $1; sum2 += $2 }
END { print sum1, sum2 }
