// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
TdopSymbols
	{
	SymbolsMap: (
		IDENTIFIER: (Curry, TdopSymbolIdentifier)
		NUMBER: (Curry, TdopSymbolNumber)
		STRING: (Curry, TdopSymbolString)
		'*': (TdopSymbolInfix, "TDOPTOKEN.MUL", lbp: 80)
		'/': (TdopSymbolInfix, "TDOPTOKEN.DIV", lbp: 80)
		'%': (TdopSymbolInfix, "TDOPTOKEN.MOD", lbp: 80)
		'$': (TdopSymbolInfix, "TDOPTOKEN.CAT", lbp: 70)
		'+': (TdopSymbolInfix, "TDOPTOKEN.ADD", lbp: 70, rbp: 80)
		'-': (TdopSymbolPreInfix, "TDOPTOKEN.SUB", lbp: 70, rbp: 80)
		'(': (TdopSymbolParens, "TDOPTOKEN.LPAREN", lbp: 85)
		')': (TdopSymbol, "TDOPTOKEN.RPAREN")
		',': (TdopSymbol, "TDOPTOKEN.COMMA")

		'and': (TdopSymbolIdentifierInfixR, "TDOPTOKEN.AND", lbp: 40, value: 'and')
		'or': (TdopSymbolIdentifierInfixR, "TDOPTOKEN.OR", lbp: 30, value: 'or')
		'not': (TdopSymbolNot, "TDOPTOKEN.NOT", lbp: 42, rbp: 80, value: 'not')

		'<': (TdopSymbolInfix, "TDOPTOKEN.LT", lbp: 60)
		'<=': (TdopSymbolInfix, "TDOPTOKEN.LTE", lbp: 60)
		'>': (TdopSymbolInfix, "TDOPTOKEN.GT", lbp: 60)
		'>=': (TdopSymbolInfix, "TDOPTOKEN.GTE", lbp: 60)
		'is': (TdopSymbolIdentifierInfix, "TDOPTOKEN.IS", lbp: 50, value: 'is')
		'isnt': (TdopSymbolIdentifierInfix, "TDOPTOKEN.ISNT", lbp: 50, value: 'isnt')
		'true': (TdopSymbolReserved, "TDOPTOKEN.TRUE", 'true')
		'false': (TdopSymbolReserved, "TDOPTOKEN.FALSE", 'false'))
	}
