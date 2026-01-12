# CSV field sum (comma-separated)
# Input: CSV file
# Measures: non-default FS handling
BEGIN { FS = "," }
{ sum += $3 }
END { print sum }
