package main

import "fmt"

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
