# Group by key and aggregate
# Input: key-value data (col1=key, col2=value)
# Measures: associative arrays
{ count[$1]++; sum[$1] += $2 }
END { for (k in count) print k, sum[k]/count[k] }
