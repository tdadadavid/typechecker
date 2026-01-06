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
		eva := tc.New()

		fmt.Printf("Checking line: %T\n", line)
		result, err := eva.Check(line)
		if err != nil {
			fmt.Fprintf(out, "Error: %v\n", err)
			continue
		}

		fmt.Fprintf(out, "%+v\n", result)
	}
}
