when reading large files, run `wc -l` first to check the line count.
if the files is over 2000 lines, use the 'offset` and `limit` parameters on the read toold to read
in chunks rather than attemting to read the entire file at once


