# Build the script
go build -o convert-jsonl main.go

# Convert first 100 lines (default)
./convert-jsonl -input occupations.jsonl -output seed_data.sql

# Convert specific number of lines
./convert-jsonl -input occupations.jsonl -output seed_data.sql -lines 50

# Convert all lines
./convert-jsonl -input occupations.jsonl -output seed_data.sql -lines 0

# Custom file names
./convert-jsonl -input data.jsonl -output output.sql -lines 100