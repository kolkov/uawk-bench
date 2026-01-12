# Count lines and fields
# Input: any text file
# Measures: basic I/O throughput
{ fields += NF }
END { print NR, fields }
