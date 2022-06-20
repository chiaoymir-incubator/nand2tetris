package main

import (
	"bufio"
	"fmt"
	"strings"
)

var (
	ErrInvalidToken      = fmt.Errorf("invalid token")
	ErrEmptyTokenizer    = fmt.Errorf("the tokenizer is empty")
	ErrNoMoreToken       = fmt.Errorf("no more token")
	ErrTooManyToken      = fmt.Errorf("too many token")
	ErrCompilationFailed = fmt.Errorf("compilation failed")
)

type CompilationChecker struct {
	letStatement     []string
	ifStatement      []string
	whileStatement   []string
	doStatement      []string
	returnStatement  []string
	elseBranch       []string
	endOfStatement   []string
	classVarDecSet   []string
	varDecSet        []string
	typeDecSet       []string
	typeReturnSet    []string
	subroutineDecSet []string
	opSet            []string
	unaryOpSet       []string
	keywordConstSet  []string
}

func NewCompilationChecker() CompilationChecker {
	return CompilationChecker{
		letStatement:     []string{"let"},
		ifStatement:      []string{"if"},
		whileStatement:   []string{"while"},
		doStatement:      []string{"do"},
		returnStatement:  []string{"return"},
		elseBranch:       []string{"else"},
		endOfStatement:   []string{";"},
		classVarDecSet:   []string{"field", "static"},
		varDecSet:        []string{"var"},
		typeDecSet:       []string{"int", "char", "boolean"},
		typeReturnSet:    []string{"int", "char", "boolean", "void"},
		subroutineDecSet: []string{"constructor", "method", "function"},
		opSet:            []string{"+", "-", "*", "/", "&", "|", "<", ">", "="},
		unaryOpSet:       []string{"-", "~"},
		keywordConstSet:  []string{"true", "false", "null", "this"},
	}
}

func (cc *CompilationChecker) containsTarget(strs []string, target string) bool {
	result := false
	for _, s := range strs {
		if s == target {
			result = true
			break
		}
	}
	return result
}

func (cc *CompilationChecker) isLetStatement(t Token) bool {
	return cc.containsTarget(cc.letStatement, t.content)
}

func (cc *CompilationChecker) isIfStatement(t Token) bool {
	return cc.containsTarget(cc.ifStatement, t.content)
}

func (cc *CompilationChecker) isWhileStatement(t Token) bool {
	return cc.containsTarget(cc.whileStatement, t.content)
}

func (cc *CompilationChecker) isDoStatement(t Token) bool {
	return cc.containsTarget(cc.doStatement, t.content)
}

func (cc *CompilationChecker) isReturnStatement(t Token) bool {
	return cc.containsTarget(cc.returnStatement, t.content)
}

func (cc *CompilationChecker) isElseBranch(t Token) bool {
	return cc.containsTarget(cc.elseBranch, t.content)
}

func (cc *CompilationChecker) isEndOfStatement(t Token) bool {
	return cc.containsTarget(cc.endOfStatement, t.content)
}

func (cc *CompilationChecker) isClassVarDec(t Token) bool {
	return cc.containsTarget(cc.classVarDecSet, t.content)
}

func (cc *CompilationChecker) isSubroutineDec(t Token) bool {
	return cc.containsTarget(cc.subroutineDecSet, t.content)
}

func (cc *CompilationChecker) isVarDec(t Token) bool {
	return cc.containsTarget(cc.varDecSet, t.content)
}

func (cc *CompilationChecker) isType(t Token) bool {
	return cc.isIdentifier(t) || cc.containsTarget(cc.typeDecSet, t.content)
}

func (cc *CompilationChecker) isExpression(t Token) bool {
	return cc.isTerm(t)
}

func (cc *CompilationChecker) isTerm(t Token) bool {
	switch {
	case cc.isIntConst(t) || cc.isStringConst(t) || cc.isKeywordConst(t) || cc.isIdentifier(t) || cc.isUnaryOp(t) || t.content == "(":
		return true
	default:
		return false
	}
}

func (cc *CompilationChecker) isOp(t Token) bool {
	return cc.containsTarget(cc.opSet, t.content)
}

func (cc *CompilationChecker) isUnaryOp(t Token) bool {
	return cc.containsTarget(cc.unaryOpSet, t.content)
}

func (cc *CompilationChecker) isIntConst(t Token) bool {
	return t.tokenType == TokenIntConst
}

func (cc *CompilationChecker) isStringConst(t Token) bool {
	return t.tokenType == TokenStringConst
}

func (cc *CompilationChecker) isKeywordConst(t Token) bool {
	return t.tokenType == TokenKeyword && cc.containsTarget(cc.keywordConstSet, t.content)
}

func (cc *CompilationChecker) isIdentifier(t Token) bool {
	return t.tokenType == TokenIdentifier
}

// This method is called after retrieving the first identifier
func (cc *CompilationChecker) isSubroutineCall(t Token) bool {
	return t.content == "(" || t.content == "."
}

// This method is called after retrieving the first identifier
func (cc *CompilationChecker) isArrayIndex(t Token) bool {
	return t.content == "["
}

var (
	cc = NewCompilationChecker()
)

type CompilationEngine struct {
	w      *bufio.Writer
	tknzr  Tokenizer
	indent int
	buf    []Token
}

func NewCompilationEngine(tknzr Tokenizer, w *bufio.Writer) CompilationEngine {
	return CompilationEngine{
		tknzr:  tknzr,
		w:      w,
		indent: 0,
		buf:    []Token{},
	}
}

/* helper methods */
func (ce *CompilationEngine) writeToken(t Token) {
	/* for output to xml */
	switch t.content {
	case "<":
		t.content = "&lt;"
	case ">":
		t.content = "&gt;"
	case "\"":
		t.content = "&quot;"
	case "&":
		t.content = "&amp;"
	default:
	}
	ce.w.WriteString(fmt.Sprintf("%s<%s> %s </%s>\n", ce.getIndent(), t.tokenType, t.content, t.tokenType))
}

func (ce *CompilationEngine) writeStartSymbol(symbol string) {
	ce.w.WriteString(fmt.Sprintf("%s<%s>\n", ce.getIndent(), symbol))
}

func (ce *CompilationEngine) writeEndSymbol(symbol string) {
	ce.w.WriteString(fmt.Sprintf("%s</%s>\n", ce.getIndent(), symbol))
}

func (ce *CompilationEngine) incrIndent() {
	ce.indent += 2
}

func (ce *CompilationEngine) decrIndent() {
	ce.indent -= 2
}

func (ce *CompilationEngine) getIndent() string {
	return strings.Repeat(" ", ce.indent)
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
	return token, nil
}

// consume a token, compare it with the given strings, and write a line
func (ce *CompilationEngine) consumeWriteCompare(strs []string) error {
	token, err := ce.consume()
	if err != nil {
		return err
	}
	if strs != nil && !cc.containsTarget(strs, token.content) {
		return ErrInvalidToken
	}
	ce.writeToken(token)
	return nil
}

func (ce *CompilationEngine) consumeWriteIdentifier() error {
	token, err := ce.consume()
	if err != nil {
		return err
	}
	if !cc.isIdentifier(token) {
		return ErrInvalidToken
	}
	ce.writeToken(token)
	return nil
}

func (ce *CompilationEngine) consumeWriteCompareIdentifier(strs []string) error {
	token, err := ce.consume()
	if err != nil {
		return err
	}
	// we assume tokenizer has already mark the compared string corresponding with the correct type
	if !(cc.isIdentifier(token) || cc.containsTarget(strs, token.content)) {
		return ErrInvalidToken
	}
	ce.writeToken(token)
	return nil
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
	symbol := "class"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare([]string{"class"}); err != nil {
		return err
	}
	if err := ce.compileClassName(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"{"}); err != nil {
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

	if err := ce.consumeWriteCompare([]string{"}"}); err != nil {
		return err
	}

	if ce.hasMoreTokens() {
		return ErrTooManyToken
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileClassVarDec() error {
	/* ('static'|'field') type varName (',' varName)* ';' */
	symbol := "classVarDec"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.classVarDecSet); err != nil {
		return err
	}
	if err := ce.compileType(cc.typeDecSet); err != nil {
		return err
	}
	if err := ce.compileVarName(); err != nil {
		return err
	}

LOOP:
	for ce.hasMoreTokens() {
		// peek token
		token := ce.getToken()
		switch token.content {
		case ",":
			if err := ce.consumeWriteCompare([]string{","}); err != nil {
				return err
			}
			if err := ce.compileVarName(); err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	if err := ce.consumeWriteCompare([]string{";"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileSubroutineDec() error {
	/* ('constructor'|'method'|'function') ('void'|type) subroutineName '(' parameterList ')' subroutineBody */
	symbol := "subroutineDec"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.subroutineDecSet); err != nil {
		return err
	}
	if err := ce.compileType(cc.typeReturnSet); err != nil {
		return err
	}
	if err := ce.compileSubroutineName(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileParameterList(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{")"}); err != nil {
		return err
	}
	if err := ce.compileSubroutineBody(); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileParameterList() error {
	/* ( (type varName) (',' type varName)* )? */
	symbol := "parameterList"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	if cc.isType(token) {
		if err := ce.compileType(cc.typeDecSet); err != nil {
			return err
		}
		if err := ce.compileVarName(); err != nil {
			return err
		}

	LOOP:
		for ce.hasMoreTokens() {
			// peek next token
			token := ce.getToken()
			switch token.content {
			case ",":
				if err := ce.consumeWriteCompare([]string{","}); err != nil {
					return err
				}
				if err := ce.compileType(cc.typeDecSet); err != nil {
					return err
				}
				if err := ce.compileVarName(); err != nil {
					return err
				}
			default:
				break LOOP
			}
		}
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileSubroutineBody() error {
	/* '{' varDec* statements '}' */
	symbol := "subroutineBody"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare([]string{"{"}); err != nil {
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

	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"}"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileVarDec() error {
	/* 'var' type varName (',' varName)* ';' */
	symbol := "varDec"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.varDecSet); err != nil {
		return err
	}
	if err := ce.compileType(cc.typeDecSet); err != nil {
		return err
	}
	if err := ce.compileVarName(); err != nil {
		return err
	}

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch token.content {
		case ",":
			if err := ce.consumeWriteCompare([]string{","}); err != nil {
				return err
			}
			if err := ce.compileVarName(); err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	if err := ce.consumeWriteCompare([]string{";"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileStatements() error {
	/* statements: ((let|if|while|do|return)Statement)* */
	symbol := "statements"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

LOOP:
	for ce.hasMoreTokens() {
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
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileLetStatement() error {
	symbol := "letStatement"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.letStatement); err != nil {
		return err
	}
	if err := ce.compileVarName(); err != nil {
		return err
	}

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	if cc.isArrayIndex(token) {
		if err := ce.consumeWriteCompare([]string{"["}); err != nil {
			return err
		}
		if err := ce.compileExpression(); err != nil {
			return err
		}
		if err := ce.consumeWriteCompare([]string{"]"}); err != nil {
			return err
		}
	}

	if err := ce.consumeWriteCompare([]string{"="}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{";"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileIfStatement() error {
	symbol := "ifStatement"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.ifStatement); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{")"}); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"{"}); err != nil {
		return err
	}
	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"}"}); err != nil {
		return err
	}

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	if cc.isElseBranch(token) {
		if err := ce.consumeWriteCompare(cc.elseBranch); err != nil {
			return err
		}
		if err := ce.consumeWriteCompare([]string{"{"}); err != nil {
			return err
		}
		if err := ce.compileStatements(); err != nil {
			return err
		}
		if err := ce.consumeWriteCompare([]string{"}"}); err != nil {
			return err
		}
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileWhileStatement() error {
	symbol := "whileStatement"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.whileStatement); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"("}); err != nil {
		return err
	}
	if err := ce.compileExpression(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{")"}); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"{"}); err != nil {
		return err
	}
	if err := ce.compileStatements(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{"}"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileDoStatement() error {
	symbol := "doStatement"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare([]string{"do"}); err != nil {
		return err
	}
	if err := ce.compileSubroutineCall(); err != nil {
		return err
	}
	if err := ce.consumeWriteCompare([]string{";"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileReturnStatement() error {
	symbol := "returnStatement"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.consumeWriteCompare(cc.returnStatement); err != nil {
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

	if err := ce.consumeWriteCompare([]string{";"}); err != nil {
		return err
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileExpression() error {
	symbol := "expression"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if err := ce.compileTerm(); err != nil {
		return err
	}

LOOP:
	for ce.hasMoreTokens() {
		token := ce.getToken()
		switch {
		case cc.isOp(token):
			if err := ce.compileOP(); err != nil {
				return err
			}
			if err := ce.compileTerm(); err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileTerm() error {
	symbol := "term"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	switch {
	case cc.isIntConst(token) || cc.isStringConst(token) || cc.isKeywordConst(token):
		if err := ce.consumeWriteCompare([]string{token.content}); err != nil {
			return err
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
			if err := ce.compileVarName(); err != nil {
				return err
			}
			if err := ce.consumeWriteCompare([]string{"["}); err != nil {
				return err
			}
			if err := ce.compileExpression(); err != nil {
				return err
			}
			if err := ce.consumeWriteCompare([]string{"]"}); err != nil {
				return err
			}
		case cc.isSubroutineCall(secondToken):
			if err := ce.compileSubroutineCall(); err != nil {
				return err
			}
		default:
			if err := ce.compileVarName(); err != nil {
				return err
			}
		}
	case token.content == "(":
		if err := ce.consumeWriteCompare([]string{"("}); err != nil {
			return err
		}
		if err := ce.compileExpression(); err != nil {
			return err
		}
		if err := ce.consumeWriteCompare([]string{")"}); err != nil {
			return err
		}
	case cc.isUnaryOp(token):
		if err := ce.compileUnaryOP(); err != nil {
			return err
		}
		if err := ce.compileTerm(); err != nil {
			return err
		}
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileSubroutineCall() error {
	if err := ce.compileIdentifier(); err != nil {
		return err
	}

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	switch token.content {
	case ".":
		if err := ce.consumeWriteCompare([]string{"."}); err != nil {
			return err
		}
		if err := ce.compileSubroutineName(); err != nil {
			return err
		}
		fallthrough
	case "(":
		if err := ce.consumeWriteCompare([]string{"("}); err != nil {
			return err
		}
		if err := ce.compileExpressionList(); err != nil {
			return err
		}
		if err := ce.consumeWriteCompare([]string{")"}); err != nil {
			return err
		}
	default:
		return ErrInvalidToken
	}

	return nil
}

func (ce *CompilationEngine) compileExpressionList() error {
	symbol := "expressionList"
	ce.writeStartSymbol(symbol)
	ce.incrIndent()

	if !ce.hasMoreTokens() {
		return ErrNoMoreToken
	}

	token := ce.getToken()
	if cc.isExpression(token) {
		if err := ce.compileExpression(); err != nil {
			return err
		}

	LOOP:
		for ce.hasMoreTokens() {
			token := ce.getToken()
			switch token.content {
			case ",":
				if err := ce.consumeWriteCompare([]string{","}); err != nil {
					return err
				}
				if err := ce.compileExpression(); err != nil {
					return err
				}
			default:
				break LOOP
			}
		}
	}

	ce.decrIndent()
	ce.writeEndSymbol(symbol)
	return nil
}

func (ce *CompilationEngine) compileOP() error {
	return ce.consumeWriteCompare(cc.opSet)
}

func (ce *CompilationEngine) compileUnaryOP() error {
	return ce.consumeWriteCompare(cc.unaryOpSet)
}

func (ce *CompilationEngine) compileType(strs []string) error {
	return ce.consumeWriteCompareIdentifier(strs)
}

func (ce *CompilationEngine) compileClassName() error {
	return ce.compileIdentifier()
}

func (ce *CompilationEngine) compileSubroutineName() error {
	return ce.compileIdentifier()
}

func (ce *CompilationEngine) compileVarName() error {
	return ce.compileIdentifier()
}

func (ce *CompilationEngine) compileIdentifier() error {
	return ce.consumeWriteIdentifier()
}
