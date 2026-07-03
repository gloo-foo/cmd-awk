# AWK Test Data Files

This directory contains sample data files for testing and demonstrating AWK functionality.

## Files

### simple_fields.txt

Simple space-separated fields for basic field extraction and manipulation.

- Fields: 3 columns of text
- Use case: Print field examples, field swapping, basic operations

### numbers.txt

One number per line for mathematical operations.

- Content: Single column of integers
- Use case: Sum, average, counting, mathematical operations

### people.csv

Comma-separated values with person information.

- Format: name,age,city
- Use case: CSV processing, field separator examples, data transformation

### fruits.txt

List of fruit names for pattern matching and filtering.

- Content: One fruit name per line
- Use case: Pattern matching, conditional processing, grep-like filtering

### duplicates.txt

Lines with repeated values for deduplication examples.

- Content: Repeated fruit names
- Use case: Unique line processing, state management with variables

### scores.txt

Student names with multiple test scores.

- Format: name score1 score2 score3
- Use case: Multi-field processing, calculating averages, row operations

### log_entries.txt

Simulated log file with timestamps and levels.

- Format: date time level message
- Use case: Log parsing, filtering by level, extracting specific fields

### prices.txt

Product names with prices.

- Format: product price
- Use case: Filtering by value, threshold comparisons, arithmetic

### tab_separated.tsv

Tab-separated employee data with headers.

- Format: Name\tAge\tDepartment\tSalary
- Use case: TSV processing, custom field separators, skipping headers

### mixed_text.txt

Various text lines for pattern matching and text operations.

- Content: Mixed phrases and sentences
- Use case: Text transformation, case-insensitive matching, line counting

## Usage in Tests

Instead of using `strings.NewReader()`, you can now use these files:

```go
// Old approach
gloo.MustRun(
    Awk(
        myProgram{},
        strings.NewReader("data\nmore data"),
    ),
)

// New approach with testdata
gloo.MustRun(
    Awk(
        myProgram{},
        gloo.File("testdata/numbers.txt"),
    ),
)
```
