# Word frequency count
# Input: text file
# Measures: split + associative arrays + sorting
{
    for (i = 1; i <= NF; i++)
        words[tolower($i)]++
}
END {
    for (w in words)
        print words[w], w
}
