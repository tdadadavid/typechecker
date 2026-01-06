package repl

import (
	"bufio"
	"fmt"
	"io"
	"typc/tc"
)

const PROMPT = ">>"

func Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	printHelp(out)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		input, err := parseInput(line)
		if err != nil {
			if err == errEmptyInput {
				continue
			}
			fmt.Fprintf(out, "Error: %v\n", err)
			continue
		}

		eva := tc.New()
		result, err := eva.Check(input)
		if err != nil {
			fmt.Fprintf(out, "Error: %v\n", err)
			continue
		}

		fmt.Fprintf(out, "%+v\n", result)
	}
}

func printHelp(out io.Writer) {
	fmt.Fprintln(out, "Typechecker REPL")
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, "  1")
	fmt.Fprintln(out, "  'hello'")
	fmt.Fprintln(out, "  \"hello\"")
	fmt.Fprintln(out, "  true")
	fmt.Fprintln(out, "  nil")
	fmt.Fprintln(out, "  + 1 2")
	fmt.Fprintln(out, "  (+ 1 2)")
	fmt.Fprintln(out, "  [\"+\", 1, 2]")
}
