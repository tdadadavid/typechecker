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
