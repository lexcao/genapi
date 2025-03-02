package printer

import (
	"fmt"
	"strings"
)

type Item interface {
	Name() string
	Print(p *Printer)
}

type Printer struct {
	ident   int
	builder strings.Builder
}

func (p *Printer) PushIndent() {
	p.ident++
}

func (p *Printer) Indent() {
	p.builder.WriteString(strings.Repeat("\t", p.ident))
}

func (p *Printer) PopIndent() {
	if p.ident > 0 {
		p.ident--
	}
}

func (p *Printer) NewLine() {
	p.builder.WriteRune('\n')
}

func (p *Printer) Quote(s string) {
	p.builder.WriteString(fmt.Sprintf("%q", s))
}

func (p *Printer) Unquoted(s string) {
	p.builder.WriteString(s)
}

func (p *Printer) Item(item Item) {
	p.builder.WriteString(item.Name())
	p.builder.WriteRune('{')
	p.NewLine()
	p.PushIndent()
	item.Print(p)
	p.PopIndent()
	p.Indent()
	p.builder.WriteRune('}')
}

func (p *Printer) KeyValueLine(key func(p *Printer), value func(p *Printer)) {
	p.Indent()
	key(p)
	p.Unquoted(": ")
	value(p)
	p.Unquoted(",")
	p.NewLine()
}

func (p *Printer) String() string {
	return p.builder.String()
}

func Print(item Item) string {
	var p Printer
	p.Item(item)
	return p.String()
}

func PrintWith(name string, printFn func(p *Printer)) string {
	return Print(itemFn{name: name, printFn: printFn})
}

type itemFn struct {
	name    string
	printFn func(p *Printer)
}

func (i itemFn) Name() string {
	return i.name
}

func (i itemFn) Print(p *Printer) {
	i.printFn(p)
}
