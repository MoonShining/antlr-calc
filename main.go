package main

import (
	"fmt"
	"github.com/MoonShining/antlr-calc/parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"strconv"
)

type visitor struct {
	*parser.BaseCalcVisitor
	stack []int
}

func (l *visitor) push(i int) {
	l.stack = append(l.stack, i)
}

func (l *visitor) pop() int {
	if len(l.stack) < 1 {
		panic("stack is empty unable to pop")
	}

	// Get the last value from the stack.
	result := l.stack[len(l.stack)-1]

	// Remove the last element from the stack.
	l.stack = l.stack[:len(l.stack)-1]

	return result
}

func (v *visitor) VisitStart(ctx *parser.StartContext) interface{} {
	ctx.Expression().Accept(v)
	return nil
}

func (v *visitor) VisitNumber(ctx *parser.NumberContext) interface{} {
	i, err := strconv.Atoi(ctx.NUMBER().GetText())
	if err != nil {
		panic(err.Error())
	}
	v.stack = append(v.stack, i)
	return nil
}
func (v *visitor) VisitMulDiv(ctx *parser.MulDivContext) interface{} {
	ctx.Expression(0).Accept(v)
	ctx.Expression(1).Accept(v)

	right, left := v.pop(), v.pop()

	switch ctx.GetOp().GetTokenType() {
	case parser.CalcParserMUL:
		v.push(left * right)
	case parser.CalcParserDIV:
		v.push(left / right)
	}
	return nil
}

func (v *visitor) VisitAddSub(ctx *parser.AddSubContext) interface{} {
	ctx.Expression(0).Accept(v)
	ctx.Expression(1).Accept(v)

	right, left := v.pop(), v.pop()

	switch ctx.GetOp().GetTokenType() {
	case parser.CalcParserADD:
		v.push(left + right)
	case parser.CalcParserSUB:
		v.push(left - right)
	}
	return nil
}

func main() {
	is := antlr.NewInputStream("1+2*3+6/2-10")
	lexer := parser.NewCalcLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewCalcParser(stream)

	v := &visitor{}
	p.Start().Accept(v)
	fmt.Print(v.pop())
}
