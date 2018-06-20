package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ei1chi/sample-lang/lexer"
	"github.com/ei1chi/sample-lang/token"
)

const PROMPT = "(´・ω・`)っ "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		l := lexer.NewLexer(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
