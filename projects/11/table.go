package main

import "fmt"

type SymbolKind int

const (
	SymbolKindInvalid SymbolKind = iota
	SymbolKindStatic
	SymbolKindField
	SymbolKindArgument
	SymbolKindLocal
)

func (sk SymbolKind) String() string {
	var str string
	switch sk {
	case SymbolKindStatic:
		str = "static"
	case SymbolKindField:
		str = "field"
	case SymbolKindArgument:
		str = "arg"
	case SymbolKindLocal:
		str = "local"
	}
	return str
}

type SymbolTableEntry struct {
	Type  string
	Kind  SymbolKind
	Index int
}

type SymbolTable struct {
	table     map[string]SymbolTableEntry
	kindTable map[SymbolKind]int
	parent    *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		table: map[string]SymbolTableEntry{},
		kindTable: map[SymbolKind]int{
			SymbolKindStatic:   0,
			SymbolKindField:    0,
			SymbolKindArgument: 0,
			SymbolKindLocal:    0,
		},
	}
}

func (st *SymbolTable) Define(name, _type string, _kind SymbolKind) {
	index := st.getVarCount(_kind)
	st.table[name] = SymbolTableEntry{
		Type:  _type,
		Kind:  _kind,
		Index: index,
	}
	st.addVarCount(_kind)
}

func (st *SymbolTable) getVarCount(kind SymbolKind) int {
	return st.kindTable[kind]
}

func (st *SymbolTable) addVarCount(kind SymbolKind) {
	st.kindTable[kind] += 1
}

func (st *SymbolTable) getAllVarCount(kind SymbolKind) int {
	if val, ok := st.kindTable[kind]; !ok {
		return 0
	} else {
		if st.parent != nil {
			return val + st.parent.getAllVarCount(kind)
		}
		return val
	}
}

func (st *SymbolTable) Contains(name string) bool {
	if _, ok := st.table[name]; ok {
		return true
	}
	if st.parent == nil {
		return false
	}
	return st.parent.Contains(name)
}

func (st *SymbolTable) KindOf(name string) SymbolKind {
	if _, ok := st.table[name]; !ok {
		if st.parent != nil {
			return st.parent.KindOf(name)
		}
		return SymbolKindInvalid
	}
	return st.table[name].Kind
}

func (st *SymbolTable) TypeOf(name string) string {
	if _, ok := st.table[name]; !ok {
		if st.parent != nil {
			return st.parent.TypeOf(name)
		}
		return "invalid"
	}
	return st.table[name].Type
}

func (st *SymbolTable) IndexOf(name string) int {
	if _, ok := st.table[name]; !ok {
		if st.parent != nil {
			return st.parent.IndexOf(name)
		}
		return -1
	}
	return st.table[name].Index
}

func (st *SymbolTable) print() {
	fmt.Printf("    Name\tType\tKind\tIndex\n")
	for name, entry := range st.table {
		fmt.Printf("    %s\t%s\t%s\t%d\n", name, entry.Type, entry.Kind, entry.Index)
	}
	fmt.Println()
	if st.parent != nil {
		st.parent.print()
	}
}
