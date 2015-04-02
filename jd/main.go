package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/printer"
)

func main() {
	app := cli.NewApp()
	app.Name = "jd"
	app.Usage = "JSON Diff"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 2 {
			fmt.Println("Not enough arguments.\n")
			fmt.Printf("Usage: %s json_file another_json_file\n", app.Name)
			os.Exit(1)
		}

		aFilePath := c.Args()[0]
		bFilePath := c.Args()[1]

		// Prepare your JSON string as `[]byte`, not `string`
		aString, err := ioutil.ReadFile(aFilePath)
		if err != nil {
			fmt.Printf("Failed to load file '%s': %s\n", aFilePath, err.Error())
			os.Exit(2)
		}

		// Another JSON string
		bString, err := ioutil.ReadFile(bFilePath)
		if err != nil {
			fmt.Printf("Failed to load file '%s': %s\n", bFilePath, err.Error())
			os.Exit(2)
		}

		// Then, compare them
		d, err := diff.Compare(aString, bString)
		if err != nil {
			fmt.Printf("Failed to unmarshal file: %s\n", err.Error())
			os.Exit(3)
		}

		// You can access the diff result with `d.Structure()`, however,
		// using `diff.Itterator` is better way to walk through the result.
		// `AsciiPrinter` implements the `diff.Itterator` interface.
		// You can create your own itterator for your purposes.
		printer := printer.NewAsciiPrinter()

		// Walk through
		d.Iterate(printer)

		// Output the result
		fmt.Print(printer.Result())
	}

	app.Run(os.Args)
}
