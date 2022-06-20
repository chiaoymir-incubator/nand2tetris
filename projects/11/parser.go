package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CompilationEngine struct {
	w              *VMWriter
	tknzr          Tokenizer
	indent         int
	buf            []Token
	msgBuf         []string
	st             *SymbolTable
	className      string
	subroutineType string
	subroutineName string
	returnType     string
	labelIndex     int
}

func NewCompilationEngine(tknzr Tokenizer, w *VMWriter) CompilationEngine {
	return CompilationEngine{
		tknzr:  tknzr,
		w:      w,
		indent: 0,
		buf:    []Token{},
		msgBuf: []string{},
		st:     NewSymbolTable(),
	}
}

/* helper methods */
func (ce *CompilationEngine) appendSymbolTable(st *SymbolTable) {
	if ce.st == nil {
		ce.st = st
	} else {
		st.parent = ce.st
		ce.st = st
	}
}

func (ce *CompilationEngine) popSymbolTable() {
	ce.st = ce.st.parent
}

func (ce *CompilationEngine) addSymbol(name, _type string, _kind SymbolKind) {
	ce.st.Define(name, _type, _kind)
}

func (ce *CompilationEngine) hasSymbol(name string) bool {
	return ce.st.Contains(name)
}

func (ce *CompilationEngine) symbolType(name string) string {
	return ce.st.TypeOf(name)
}

func (ce *CompilationEngine) symbolKind(name string) SymbolKind {
	return ce.st.KindOf(name)
}

func (ce *CompilationEngine) symbolIndex(name string) int {
	return ce.st.IndexOf(name)
}

func (ce *CompilationEngine) printSymbolTable() {
	ce.st.print()
}

func (ce *CompilationEngine) getFuncName() string {
	return ce.className + "." + ce.subroutineName
}

func (ce *CompilationEngine) advance() {
	ce.tknzr.Advance()
}

func (ce *CompilationEngine) hasMoreTokens() bool {
	return ce.hasBufferedToken() || ce.tknzr.HasMoreTokens()
}

func (ce *CompilationEngine) getToken() Token {
	if ce.hasBufferedToken() {
		return ce.peekBufferedToken()
	}
	return ce.tknzr.GetToken()
}
func (ce *CompilationEngine) getTokenContent() string {
	return ce.getToken().content
}

func (ce *CompilationEngine) hasBufferedToken() bool {
	return len(ce.buf) > 0
}

func (ce *CompilationEngine) peekBufferedToken() Token {
	return ce.buf[0]
}

func (ce *CompilationEngine) consumeBufferedToken() Token {
	token := ce.buf[0]
	ce.buf = ce.buf[1:]
	return token
}

func (ce *CompilationEngine) addBufferedToken(t Token) {
	ce.buf = append(ce.buf, t)
}

func (ce *CompilationEngine) consume() (Token, error) {
	if ce.hasBufferedToken() {
		return ce.consumeBufferedToken(), nil
	}
	if !ce.tknzr.HasMoreTokens() {
		return Token{}, ErrNoMoreToken
	}
	token := ce.tknzr.GetToken()
	ce.tknzr.Advance()
	ce.msgBuf = append(ce.msgBuf, token.content)
	return token, nil
}

func (ce *CompilationEngine) consumeCompare(strs []string) error {
	token, err := ce.consume()
	if err != nil {
		return err
	}
	if strs != nil && !cc.containsTarget(strs, token.content) {
		return ErrInvalidToken
	}
	return nil
}

func (ce *CompilationEngine) consumeCompareIdentifier(strs []string) error {
	token, err := ce.consume()
	if err != nil {
		return err
	}
	// we assume tokenizer has already mark the compared string corresponding with the correct type
	if !(cc.isIdentifier(token) || cc.containsTarget(strs, token.content)) {
		return ErrInvalidToken
	}
	return nil
}

func (ce *CompilationEngine) clearDebugMsg() {
	ce.msgBuf = []string{}
}

func (ce *CompilationEngine) printDebugMsg() {
	ce.w.WriteString(fmt.Sprintf("  // %s\n", strings.Join(ce.msgBuf, " ")))
}

/* Main compilation logic. */
func (ce *CompilationEngine) Compile() error {
	ce.advance()
	// empty tokenizer
	if !ce.hasMoreTokens() {
		return ErrEmptyTokenizer
	}

	token := ce.getToken()
	if token.content != "class" {
		return ErrInvalidToken
	}

	return ce.compileClass()
}

func (ce *CompilationEngine) compileClass() error {
	/* 'class' className '{' classVarDec* subroutineDec* '}' */
	if err := ce.consumeCompare([]string{"class"}); err != nil {
		return err
	}
	// remember class name for subroutine naming
	ce.className = ce.getTokenContent()
	if err := ce.compileClassName(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"{"}); err != nil {
		return err
	}

LOOP1:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch {
		case cc.isClassVarDec(token):
			if err := ce.compileClassVarDec(); err != nil {
				return err
			}
		default:
			break LOOP1
		}
	}

LOOP2:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch {
		case cc.isSubroutineDec(token):
			if err := ce.compileSubroutineDec(); err != nil {
				return err
			}
		default:
			break LOOP2
		}
	}

	if err := ce.consumeCompare([]string{"}"}); err != nil {
		return err
	}

	if ce.hasMoreTokens() {
		return ErrTooManyToken
	}

	return nil
}

func (ce *CompilationEngine) compileClassVarDec() error {
	/* ('static'|'field') type varName (',' varName)* ';' */
	varDec := ce.getTokenContent()
	if err := ce.consumeCompare(cc.classVarDecSet); err != nil {
		return err
	}
	var varKind SymbolKind
	if varDec == "static" {
		varKind = SymbolKindStatic
	} else {
		varKind = SymbolKindField
	}

	varType := ce.getTokenContent()
	if err := ce.compileType(cc.typeDecSet); err != nil {
		return err
	}

	varName := ce.getTokenContent()
	if err := ce.compileVarName(); err != nil {
		return err
	}

	ce.addSymbol(varName, varType, varKind)

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch token.content {
		case ",":
			if _, err := ce.consume(); err != nil {
				return err
			}
			varName = ce.getToken().content
			if err := ce.compileVarName(); err != nil {
				return err
			}

			ce.addSymbol(varName, varType, varKind)
		default:
			break LOOP
		}
	}

	if err := ce.consumeCompare([]string{";"}); err != nil {
		return err
	}

	return nil
}

func (ce *CompilationEngine) compileSubroutineDec() error {
	/* ('constructor'|'method'|'function') ('void'|type) subroutineName '(' parameterList ')' subroutineBody */
	newTable := NewSymbolTable()
	ce.appendSymbolTable(newTable)
	ce.labelIndex = 0

	ce.subroutineType = ce.getTokenContent()
	if err := ce.consumeCompare(cc.subroutineDecSet); err != nil {
		return err
	}

	ce.returnType = ce.getTokenContent()
	if err := ce.compileType(cc.typeReturnSet); err != nil {
		return err
	}

	ce.subroutineName = ce.getTokenContent()
	if err := ce.compileSubroutineName(); err != nil {
		return err
	}

	if err := ce.consumeCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileParameterList(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{")"}); err != nil {
		return err
	}
	if err := ce.compileSubroutineBody(); err != nil {
		return err
	}

	ce.popSymbolTable()
	ce.w.WriteString("\n")

	return nil
}

func (ce *CompilationEngine) compileParameterList() error {
	/* ( (type varName) (',' type varName)* )? */
	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	if ce.subroutineType == "method" {
		ce.addSymbol("this", ce.subroutineType, SymbolKindArgument)
	}

	token := ce.getToken()
	if cc.isType(token) {
		varType := ce.getTokenContent()
		if err := ce.compileType(cc.typeDecSet); err != nil {
			return err
		}
		varName := ce.getTokenContent()
		if err := ce.compileVarName(); err != nil {
			return err
		}
		ce.addSymbol(varName, varType, SymbolKindArgument)

	LOOP:
		for ce.hasMoreTokens() {
			token := ce.getToken()
			switch token.content {
			case ",":
				if _, err := ce.consume(); err != nil {
					return err
				}
				varType := ce.getTokenContent()
				if err := ce.compileType(cc.typeDecSet); err != nil {
					return err
				}
				varName := ce.getTokenContent()
				if err := ce.compileVarName(); err != nil {
					return err
				}
				ce.addSymbol(varName, varType, SymbolKindArgument)

			default:
				break LOOP
			}
		}
	}

	return nil
}

func (ce *CompilationEngine) compileSubroutineBody() error {
	/* '{' varDec* statements '}' */
	if err := ce.consumeCompare([]string{"{"}); err != nil {
		return err
	}

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch {
		case cc.isVarDec(token):
			if err := ce.compileVarDec(); err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	fmt.Printf("subroutine %s:\n", ce.subroutineName)
	ce.printSymbolTable()

	ce.w.WriteFunction(ce.getFuncName(), ce.st.getVarCount(SymbolKindLocal))
	if ce.subroutineType == "constructor" {
		ce.w.WritePush(SegmentConst, ce.st.getAllVarCount(SymbolKindField))
		ce.w.WriteCall("Memory.alloc", 1)
		ce.w.WritePop(SegmentPointer, 0)
	} else if ce.subroutineType == "method" {
		// anchor this
		ce.w.WritePush(SegmentArg, 0)
		ce.w.WritePop(SegmentPointer, 0)
	}

	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"}"}); err != nil {
		return err
	}

	return nil
}

func (ce *CompilationEngine) compileVarDec() error {
	/* 'var' type varName (',' varName)* ';' */
	if err := ce.consumeCompare(cc.varDecSet); err != nil {
		return err
	}

	varType := ce.getTokenContent()
	if err := ce.compileType(cc.typeDecSet); err != nil {
		return err
	}
	varName := ce.getTokenContent()
	if err := ce.compileVarName(); err != nil {
		return err
	}
	ce.addSymbol(varName, varType, SymbolKindLocal)

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch token.content {
		case ",":
			if _, err := ce.consume(); err != nil {
				return err
			}
			varName := ce.getTokenContent()
			if err := ce.compileVarName(); err != nil {
				return err
			}
			ce.addSymbol(varName, varType, SymbolKindLocal)

		default:
			break LOOP
		}
	}

	if err := ce.consumeCompare([]string{";"}); err != nil {
		return err
	}

	return nil
}

func (ce *CompilationEngine) compileStatements() error {
	/* statements: ((let|if|while|do|return)Statement)* */
LOOP:
	for ce.hasMoreTokens() {
		ce.clearDebugMsg()
		token := ce.getToken()
		switch {
		case cc.isLetStatement(token):
			if err := ce.compileLetStatement(); err != nil {
				return err
			}
		case cc.isIfStatement(token):
			if err := ce.compileIfStatement(); err != nil {
				return err
			}
		case cc.isWhileStatement(token):
			if err := ce.compileWhileStatement(); err != nil {
				return err
			}
		case cc.isDoStatement(token):
			if err := ce.compileDoStatement(); err != nil {
				return err
			}
		case cc.isReturnStatement(token):
			if err := ce.compileReturnStatement(); err != nil {
				return err
			}
		default:
			break LOOP
		}
		ce.printDebugMsg()
	}

	return nil
}

func (ce *CompilationEngine) compileLetStatement() error {
	/* 'let' varName ('[' expression ']')? '=' expression ';' */
	if err := ce.consumeCompare(cc.letStatement); err != nil {
		return err
	}

	varName := ce.getTokenContent()
	if err := ce.compileVarName(); err != nil {
		return err
	}
	varNameKind := KindToSegment(ce.symbolKind(varName))
	varNameIndex := ce.symbolIndex(varName)

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	var isArrayIndexing bool
	token := ce.getToken()
	if cc.isArrayIndex(token) {
		// push arr
		ce.w.WritePush(varNameKind, varNameIndex)
		if err := ce.consumeCompare([]string{"["}); err != nil {
			return err
		}
		if err := ce.compileExpression(); err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{"]"}); err != nil {
			return err
		}
		ce.w.WriteArithmetic(CommandAdd)
		isArrayIndexing = true
	}

	if err := ce.consumeCompare([]string{"="}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{";"}); err != nil {
		return err
	}

	if isArrayIndexing {
		// pop temp 0
		// pop pointer 1
		// push temp 0
		// pop that 0
		ce.w.WritePop(SegmentTemp, 0)
		ce.w.WritePop(SegmentPointer, 1)
		ce.w.WritePush(SegmentTemp, 0)
		ce.w.WritePop(SegmentThat, 0)
	} else {
		// pop varName
		ce.w.WritePop(varNameKind, varNameIndex)
	}

	return nil
}

func (ce *CompilationEngine) compileIfStatement() error {
	/* 'if' '(' expression ')' '{' statements '}' ( 'else' '{' statements '}' )? */
	label1 := fmt.Sprintf("%s.%d", ce.getFuncName(), ce.labelIndex)
	ce.labelIndex++
	label2 := fmt.Sprintf("%s.%d", ce.getFuncName(), ce.labelIndex)
	ce.labelIndex++

	if err := ce.consumeCompare(cc.ifStatement); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{")"}); err != nil {
		return err
	}

	ce.w.WriteArithmetic(CommandNot)
	ce.w.WriteIfGoto(label1)

	if err := ce.consumeCompare([]string{"{"}); err != nil {
		return err
	}
	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"}"}); err != nil {
		return err
	}

	ce.w.WriteGoto(label2)

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	ce.w.WriteLabel(label1)

	token := ce.getToken()
	if cc.isElseBranch(token) {
		if _, err := ce.consume(); err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{"{"}); err != nil {
			return err
		}
		if err := ce.compileStatements(); err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{"}"}); err != nil {
			return err
		}
	}

	ce.w.WriteLabel(label2)
	return nil
}

func (ce *CompilationEngine) compileWhileStatement() error {
	/* 'while' '(' expression ')' '{' statements '}' */
	label1 := fmt.Sprintf("%s.%d", ce.getFuncName(), ce.labelIndex)
	ce.labelIndex++
	label2 := fmt.Sprintf("%s.%d", ce.getFuncName(), ce.labelIndex)
	ce.labelIndex++
	ce.w.WriteLabel(label1)

	if err := ce.consumeCompare(cc.whileStatement); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{")"}); err != nil {
		return err
	}

	ce.w.WriteArithmetic(CommandNot)
	ce.w.WriteIfGoto(label2)

	if err := ce.consumeCompare([]string{"{"}); err != nil {
		return err
	}
	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{"}"}); err != nil {
		return err
	}

	ce.w.WriteGoto(label1)
	ce.w.WriteLabel(label2)
	return nil
}

func (ce *CompilationEngine) compileDoStatement() error {
	/* 'do' subroutineCall ';' */
	if err := ce.consumeCompare([]string{"do"}); err != nil {
		return err
	}
	if err := ce.compileSubroutineCall(); err != nil {
		return err
	}
	if err := ce.consumeCompare([]string{";"}); err != nil {
		return err
	}

	// pop dummy return value
	ce.w.WritePop(SegmentTemp, 0)

	return nil
}

func (ce *CompilationEngine) compileReturnStatement() error {
	/* 'return' expression? ';' */
	if err := ce.consumeCompare(cc.returnStatement); err != nil {
		return err
	}

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	if !cc.isEndOfStatement(token) {
		if err := ce.compileExpression(); err != nil {
			return err
		}
	}

	if err := ce.consumeCompare([]string{";"}); err != nil {
		return err
	}

	if ce.returnType == "void" {
		// push dummy value to stack as return value
		ce.w.WritePush(SegmentConst, 0)
	}
	ce.w.WriteReturn()

	return nil
}

func (ce *CompilationEngine) compileExpression() error {
	/* term (op term)* */
	if err := ce.compileTerm(); err != nil {
		return err
	}

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch {
		case cc.isOp(token):
			op := ce.getTokenContent()
			if err := ce.compileOP(); err != nil {
				return err
			}
			if err := ce.compileTerm(); err != nil {
				return err
			}

			switch op {
			case "*":
				ce.w.WriteCall("Math.multiply", 2)
			case "/":
				ce.w.WriteCall("Math.divide", 2)
			default:
				ce.w.WriteArithmetic(opCommand(op))
			}

		default:
			break LOOP
		}
	}

	return nil
}

func (ce *CompilationEngine) compileTerm() error {
	/* intergerConstant | stringConstant | keywordConstant | varName | varName'[' expression ']' | subroutineCall | '(' expression ')' | unaryOP term */
	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	switch {
	case cc.isIntConst(token):
		if _, err := ce.consume(); err != nil {
			return err
		}
		i, _ := strconv.Atoi(token.content)
		ce.w.WritePush(SegmentConst, i)

	case cc.isStringConst(token):
		if _, err := ce.consume(); err != nil {
			return err
		}
		ce.w.WritePush(SegmentConst, len(token.content))
		ce.w.WriteCall("String.new", 1)
		for _, v := range token.content {
			// convert to ascii
			ce.w.WritePush(SegmentConst, int(v))
			ce.w.WriteCall("String.appendChar", 2)
		}

	case cc.isKeywordConst(token):
		if _, err := ce.consume(); err != nil {
			return err
		}
		switch token.content {
		case "null", "false":
			ce.w.WritePush(SegmentConst, 0)
		case "true":
			ce.w.WritePush(SegmentConst, 1)
			ce.w.WriteArithmetic(CommandNeg)
		case "this":
			ce.w.WritePush(SegmentPointer, 0)
		}

	case cc.isIdentifier(token):
		firstToken, err := ce.consume()
		if err != nil {
			return ErrInvalidToken
		}
		// add token to buffer to be rewinded
		ce.addBufferedToken(firstToken)

		// intentionally check next token
		if !ce.tknzr.HasMoreTokens() {
			return ErrNoMoreToken
		}
		// intentionally peek next token
		secondToken := ce.tknzr.GetToken()

		switch {
		case cc.isArrayIndex(secondToken):
			// push arr
			varName := ce.getTokenContent()
			varNameKind, varNameIndex := KindToSegment(ce.symbolKind(varName)), ce.symbolIndex(varName)
			ce.w.WritePush(varNameKind, varNameIndex)
			if err := ce.compileVarName(); err != nil {
				return err
			}
			if err := ce.consumeCompare([]string{"["}); err != nil {
				return err
			}
			if err := ce.compileExpression(); err != nil {
				return err
			}
			if err := ce.consumeCompare([]string{"]"}); err != nil {
				return err
			}
			ce.w.WriteArithmetic(CommandAdd)
			ce.w.WritePop(SegmentPointer, 1)
			ce.w.WritePush(SegmentThat, 0)

		case cc.isSubroutineCall(secondToken):
			if err := ce.compileSubroutineCall(); err != nil {
				return err
			}

		default:
			varName := ce.getTokenContent()
			if err := ce.compileVarName(); err != nil {
				return err
			}
			ce.w.WritePush(KindToSegment(ce.symbolKind(varName)), ce.symbolIndex(varName))
		}

	case token.content == "(":
		if err := ce.consumeCompare([]string{"("}); err != nil {
			return err
		}
		if err := ce.compileExpression(); err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{")"}); err != nil {
			return err
		}

	case cc.isUnaryOp(token):
		unaryOP := ce.getTokenContent()
		if err := ce.compileUnaryOP(); err != nil {
			return err
		}
		if err := ce.compileTerm(); err != nil {
			return err
		}
		ce.w.WriteArithmetic(unaryOPCommand(unaryOP))
	}

	return nil
}

func (ce *CompilationEngine) compileSubroutineCall() error {
	/* subroutineName '(' expressionList ')' | (className|varName) '.' subroutineName '(' expressionList ')' */
	// this method is called after 'do' statement
	name := ce.getTokenContent()
	if err := ce.consumeCompareIdentifier(nil); err != nil {
		return err
	}

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	switch token.content {
	case "(":
		// the subroutine is in the current class (a method)
		ce.w.WritePush(SegmentPointer, 0)

		if err := ce.consumeCompare([]string{"("}); err != nil {
			return err
		}
		nArgs, err := ce.compileExpressionList()
		if err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{")"}); err != nil {
			return err
		}

		ce.w.WriteCall(ce.className+"."+name, nArgs+1)

	case ".":
		if _, err := ce.consume(); err != nil {
			return err
		}

		secondName := ce.getTokenContent()
		if err := ce.compileSubroutineName(); err != nil {
			return err
		}

		var isMethod bool
		var typeName string
		// varName (a method)
		if ce.hasSymbol(name) {
			typeName = ce.symbolType(name)
			ce.w.WritePush(KindToSegment(ce.symbolKind(name)), ce.symbolIndex(name))
			isMethod = true
		}
		// else it is a className and we do nothing (a function)

		if err := ce.consumeCompare([]string{"("}); err != nil {
			return err
		}
		nArgs, err := ce.compileExpressionList()
		if err != nil {
			return err
		}
		if err := ce.consumeCompare([]string{")"}); err != nil {
			return err
		}

		if isMethod {
			ce.w.WriteCall(typeName+"."+secondName, nArgs+1)
		} else {
			ce.w.WriteCall(name+"."+secondName, nArgs)
		}

	default:
		return ErrInvalidToken
	}

	return nil
}

func (ce *CompilationEngine) compileExpressionList() (int, error) {
	/* (expression (',' expression)* )? */
	if !ce.hasMoreTokens() {
		return 0, ErrNoMoreToken
	}

	nArgs := 0
	token := ce.getToken()
	if cc.isExpression(token) {
		if err := ce.compileExpression(); err != nil {
			return 0, err
		}
		nArgs++

	LOOP:
		for ce.hasMoreTokens() {
			token := ce.getToken()
			switch token.content {
			case ",":
				if _, err := ce.consume(); err != nil {
					return 0, err
				}
				if err := ce.compileExpression(); err != nil {
					return 0, err
				}
				nArgs++

			default:
				break LOOP
			}
		}
	}

	return nArgs, nil
}

func (ce *CompilationEngine) compileOP() error {
	return ce.consumeCompare(cc.opSet)
}

func (ce *CompilationEngine) compileUnaryOP() error {
	return ce.consumeCompare(cc.unaryOpSet)
}

func (ce *CompilationEngine) compileType(strs []string) error {
	return ce.consumeCompareIdentifier(strs)
}

func (ce *CompilationEngine) compileClassName() error {
	token, err := ce.consume()
	if err != nil || !cc.isIdentifier(token) {
		return err
	}
	return nil
}

func (ce *CompilationEngine) compileSubroutineName() error {
	token, err := ce.consume()
	if err != nil || cc.isIdentifier(token) {
		return err
	}
	return nil
}

func (ce *CompilationEngine) compileVarName() error {
	token, err := ce.consume()
	if err != nil || cc.isIdentifier(token) {
		return err
	}
	if !ce.hasSymbol(token.content) {
		return ErrInvalidToken
	}
	return nil
}
