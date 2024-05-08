# go-tablewriter

This module implements a writer for table data, which can be output as CSV or Text.

Example:

```go
package main

import (
    tablewriter "github.com/djthorpe/go-tablewriter"

)

func main() {
	table := []TableData{
		{A: "hello", B: "world"},
	}
    writer := tablewriter.New(os.Stdout)
    writer.Write(table, tablewriter.OptHeader())
}
```

The `Write` function expects a single struct or a slice of structs as the first argument. Each struct represents a row in the table. 
The struct fields (including any which are embedded) are used as columns in the table.

## Table Options

The following options can be used to customize the output:

  - `tablewriter.OptHeader()`: Output the header row.
  - `tablewriter.OptTag("json")`: Set the struct field tag which is used to set options, default is "json".
  - `tablewriter.OptFieldDelim('|')`: Set the field delimiter, default is ',' for CSV and '|' for Text.
  - `tablewriter.OptOutputCSV()`: Output as CSV.
  - `tablewriter.OptOutputText()`: Output as Text.
  - `tablewriter.OptNull("<nil>")`: Set how the nil value is represented in the output, defaults to `<nil>`.

## Struct Tags

Tags on struct fields can determine how the field is output. The `json` tag is used by default.

  - `json:"-"`: Skip the field.
  - `json:"Name"`: Set the column header to "Name".
  - `json:",omitdefault"`: If all values in the table are zero-valued, skip output of the column.

## Customize Field Output

You can implement the following interface on any field to customize how it is output:

```go
type Marshaller interface {
	Marshal() ([]byte, error)
}
```

By default, strings are output as-is and other values are marshalled using the `encoding/json` package.

## Contribution and License

See the [LICENSE](LICENSE) file for license rights and limitations, currently Apache. Pull requests and [issues](https://github.com/djthorpe/go-tablewriter/issues) are welcome.

## Changelog

  * v0.0.1 (May 2024) Initial version
  * v0.0.2 (May 2024) Documentation updates

Future versions will include more options for customizing the output:

tablewriter:

  * Omitting columns based on zero-value
  * Adding JSON and SQL output

pkg/text:

  * Sizing the width of fields for the text package
  * Setting the width of the table based on terminal width
  * Wrapping text in columns
  * Outputing fields with ANSI color codes
