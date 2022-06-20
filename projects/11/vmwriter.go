package main

import (
	"bufio"
	"fmt"
	"os"
)

type Segment string

const (
	SegmentConst   Segment = "constant"
	SegmentArg     Segment = "argument"
	SegmentLocal   Segment = "local"
	SegmentStatic  Segment = "static"
	SegmentThis    Segment = "this"
	SegmentThat    Segment = "that"
	SegmentPointer Segment = "pointer"
	SegmentTemp    Segment = "temp"
)

func KindToSegment(kind SymbolKind) Segment {
	var seg Segment
	switch kind {
	case SymbolKindField:
		return SegmentThis
	case SymbolKindLocal:
		return SegmentLocal
	case SymbolKindArgument:
		return SegmentArg
	case SymbolKindStatic:
		return SegmentStatic
	}

	return seg
}

type Command string

const (
	CommandPush     Command = "push"
	CommandPop      Command = "pop"
	CommandAdd      Command = "add"
	CommandSub      Command = "sub"
	CommandNeg      Command = "neg"
	CommandEq       Command = "eq"
	CommandGt       Command = "gt"
	CommandLt       Command = "lt"
	CommandAnd      Command = "and"
	CommandOr       Command = "or"
	CommandNot      Command = "not"
	CommandLabel    Command = "label"
	CommandGoto     Command = "goto"
	CommandIfGoto   Command = "if-goto"
	CommandCall     Command = "call"
	CommandFunction Command = "function"
)

func opCommand(op string) Command {
	var command Command
	switch op {
	case "+":
		command = CommandAdd
	case "-":
		command = CommandSub
	case "=":
		command = CommandEq
	case ">":
		command = CommandGt
	case "<":
		command = CommandLt
	case "&":
		command = CommandAnd
	case "|":
		command = CommandOr
	}
	return command
}

func unaryOPCommand(unaryOP string) Command {
	var command Command
	switch unaryOP {
	case "-":
		command = CommandNeg
	case "~":
		command = CommandNot
	}
	return command
}

const LineIndent string = "  "

type VMWriter struct {
	f *os.File
	w *bufio.Writer
}

func NewVMWriter(f *os.File) *VMWriter {
	return &VMWriter{
		f: f,
		w: bufio.NewWriter(f),
	}
}

func (vmw *VMWriter) WritePush(seg Segment, i int) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s %s %d\n", CommandPush, seg, i))
}

func (vmw *VMWriter) WritePop(seg Segment, i int) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s %s %d\n", CommandPop, seg, i))
}

func (vmw *VMWriter) WriteArithmetic(command Command) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s\n", command))
}

func (vmw *VMWriter) WriteLabel(name string) {
	vmw.w.WriteString(fmt.Sprintf("%s %s\n", CommandLabel, name))
}

func (vmw *VMWriter) WriteGoto(name string) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s %s\n", CommandGoto, name))
}

func (vmw *VMWriter) WriteIfGoto(name string) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s %s\n", CommandIfGoto, name))
}

func (vmw *VMWriter) WriteCall(name string, nArgs int) {
	vmw.w.WriteString(LineIndent + fmt.Sprintf("%s %s %d\n", CommandCall, name, nArgs))
}

func (vmw *VMWriter) WriteFunction(name string, nLocals int) {
	vmw.w.WriteString(fmt.Sprintf("%s %s %d\n", CommandFunction, name, nLocals))
}

func (vmw *VMWriter) WriteReturn() {
	vmw.w.WriteString(LineIndent + "return" + "\n")
}

func (vmw *VMWriter) WriteString(s string) {
	vmw.w.WriteString(s)
}

func (vmw *VMWriter) Flush() error {
	return vmw.w.Flush()
}

func (vmw *VMWriter) Close() {
	vmw.f.Close()
}
