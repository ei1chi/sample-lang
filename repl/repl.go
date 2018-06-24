package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ei1chi/sample-lang/eval"
	"github.com/ei1chi/sample-lang/lexer"
	"github.com/ei1chi/sample-lang/parser"
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
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaled := eval.Eval(program)
		if evaled != nil {
			io.WriteString(out, evaled.Inspect())
			io.WriteString(out, "\n")
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
