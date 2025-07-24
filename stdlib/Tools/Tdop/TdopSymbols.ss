// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
MemoizeSingle
	{
	Func()
		{
		map = Object()
		for m, v in .SymbolsMap
			{
			if .ExcludeSymbols.Has?(m)
				continue

			fn = Global(v[0])
			args = v.Copy().Delete(0)
			args[0] = Global(args[0])  // token
			map[m] = fn(@args)
			}
		return map
		}

	ExcludeSymbols: ()
	SymbolsMap: (
		IDENTIFIER: (Curry, TdopSymbolIdentifier)
		NUMBER: (Curry, TdopSymbolNumber)
		STRING: (Curry, TdopSymbolString)

		'++': (TdopSymbolIncDec, "TDOPTOKEN.INC", lbp: 90, rbp: 80)
		'--': (TdopSymbolIncDec, "TDOPTOKEN.DEC", lbp: 90, rbp: 80)
		'.': (TdopSymbolDot, "TDOPTOKEN.DOT", lbp: 90)
		'*': (TdopSymbolInfix, "TDOPTOKEN.MUL", lbp: 80)
		'/': (TdopSymbolInfix, "TDOPTOKEN.DIV", lbp: 80)
		'%': (TdopSymbolInfix, "TDOPTOKEN.MOD", lbp: 80)

		'$': (TdopSymbolInfix, "TDOPTOKEN.CAT", lbp: 70)
		'+': (TdopSymbolPreInfix, "TDOPTOKEN.ADD", lbp: 70, rbp: 80)
		'-': (TdopSymbolPreInfix, "TDOPTOKEN.SUB", lbp: 70, rbp: 80)

		'<<': (TdopSymbolInfix, "TDOPTOKEN.LSHIFT", lbp: 65)
		'>>': (TdopSymbolInfix, "TDOPTOKEN.RSHIFT", lbp: 65)

		'<': (TdopSymbolInfix, "TDOPTOKEN.LT", lbp: 60)
		'<=': (TdopSymbolInfix, "TDOPTOKEN.LTE", lbp: 60)
		'>': (TdopSymbolInfix, "TDOPTOKEN.GT", lbp: 60)
		'>=': (TdopSymbolInfix, "TDOPTOKEN.GTE", lbp: 60)

		'is': (TdopSymbolIdentifierInfix, "TDOPTOKEN.IS", lbp: 50, value: 'is')
		'==': (TdopSymbolInfix, "TDOPTOKEN.IS", lbp: 50)
		'isnt': (TdopSymbolIdentifierInfix, "TDOPTOKEN.ISNT", lbp: 50, value: 'isnt')
		'!=': (TdopSymbolInfix, "TDOPTOKEN.ISNT", lbp: 50)
		'=~': (TdopSymbolInfix, "TDOPTOKEN.MATCH", lbp: 50)
		'!~': (TdopSymbolInfix, "TDOPTOKEN.MATCHNOT", lbp: 50)

		'&': (TdopSymbolInfix, "TDOPTOKEN.BITAND", lbp: 48)
		'^': (TdopSymbolInfix, "TDOPTOKEN.BITXOR", lbp: 46)
		'|': (TdopSymbolInfix, "TDOPTOKEN.BITOR", lbp: 44)
		'~': (TdopSymbolPrefix, "TDOPTOKEN.BITNOT", lbp: 0, rbp: 80)

		'in': (TdopSymbolIn, "TDOPTOKEN.IN", lbp: 42, value: 'in')

		'and': (TdopSymbolIdentifierInfixR, "TDOPTOKEN.AND", lbp: 40, value: 'and')
		'&&': (TdopSymbolInfixR, "TDOPTOKEN.AND", lbp: 40)
		'or': (TdopSymbolIdentifierInfixR, "TDOPTOKEN.OR", lbp: 30, value: 'or')
		'||': (TdopSymbolInfixR, "TDOPTOKEN.OR", lbp: 40)
		'not': (TdopSymbolNot, "TDOPTOKEN.NOT", lbp: 42, rbp: 80, value: 'not')
		'!': (TdopSymbolPrefix, "TDOPTOKEN.NOT", lbp: 0, rbp: 80)
		'#': (TdopSymbolHash, "TDOPTOKEN.HASH")
		'@': (TdopSymbolAt, "TDOPTOKEN.AT")
		'?': (TdopSymbolTrinary, "TDOPTOKEN.Q_MARK", lbp: 20)
		':': (TdopSymbol, "TDOPTOKEN.COLON")

		'=': (TdopSymbolInfixR, "TDOPTOKEN.EQ", lbp: 90, rbp: 10)
		'+=': (TdopSymbolInfixR, "TDOPTOKEN.ADDEQ", lbp: 90, rbp: 10)
		'-=': (TdopSymbolInfixR, "TDOPTOKEN.SUBEQ", lbp: 90, rbp: 10)
		'$=': (TdopSymbolInfixR, "TDOPTOKEN.CATEQ", lbp: 90, rbp: 10)
		'*=': (TdopSymbolInfixR, "TDOPTOKEN.MULEQ", lbp: 90, rbp: 10)
		'/=': (TdopSymbolInfixR, "TDOPTOKEN.DIVEQ", lbp: 90, rbp: 10)
		'%=': (TdopSymbolInfixR, "TDOPTOKEN.MODEQ", lbp: 90, rbp: 10)
		'<<=': (TdopSymbolInfixR, "TDOPTOKEN.LSHIFTEQ", lbp: 90, rbp: 10)
		'>>=': (TdopSymbolInfixR, "TDOPTOKEN.RSHIFTEQ", lbp: 90, rbp: 10)
		'|=': (TdopSymbolInfixR, "TDOPTOKEN.BITOREQ", lbp: 90, rbp: 10)
		'&=': (TdopSymbolInfixR, "TDOPTOKEN.BITANDEQ", lbp: 90, rbp: 10)
		'^=': (TdopSymbolInfixR, "TDOPTOKEN.BITXOREQ", lbp: 90, rbp: 10)

		'..': (TdopSymbol, "TDOPTOKEN.RANGETO")
		'::': (TdopSymbol, "TDOPTOKEN.RANGELEN")

		'(': (TdopSymbolParens, "TDOPTOKEN.LPAREN", lbp: 85)
		')': (TdopSymbol, "TDOPTOKEN.RPAREN")
		'[': (TdopSymbolBracket, "TDOPTOKEN.LBRACKET", lbp: 90)
		']': (TdopSymbol, "TDOPTOKEN.RBRACKET")
		';': (TdopSymbol, "TDOPTOKEN.SEMICOLON")
		'{': (TdopSymbolCurly, "TDOPTOKEN.LCURLY", lbp: 90)
		'}': (TdopSymbol, "TDOPTOKEN.RCURLY")
		',': (TdopSymbol, "TDOPTOKEN.COMMA")

		'true': (TdopSymbolReserved, "TDOPTOKEN.TRUE", 'true')
		'false': (TdopSymbolReserved, "TDOPTOKEN.FALSE", 'false')
		'function': (TdopSymbolFunction, "TDOPTOKEN.FUNCTION", 'function')
		'class': (TdopSymbolClass, "TDOPTOKEN.CLASS", 'class')
		'dll': (TdopSymbolDll, "TDOPTOKEN.DLL", 'dll')
		'struct': (TdopSymbolStruct, "TDOPTOKEN.STRUCT", 'struct')
		'callback': (TdopSymbolCallback, "TDOPTOKEN.CALLBACK", 'callback')
		'super': (TdopSymbolSuper, "TDOPTOKEN.SUPER", 'super')
		'new': (TdopSymbolNew, "TDOPTOKEN.NEW", 'new', rbp: 87)
		'if': (TdopSymbolIfElse, "TDOPTOKEN.IF", 'if')
		'else': (TdopSymbolStmtIdentifier, "TDOPTOKEN.ELSE", 'else')
		'while': (TdopSymbolWhile, "TDOPTOKEN.WHILE", 'while')
		'do': (TdopSymbolDo, "TDOPTOKEN.DO", 'do')
		'for': (TdopSymbolFor, "TDOPTOKEN.FOR", 'for')
		'forever': (TdopSymbolForever, "TDOPTOKEN.FOREVER", 'forever')
		'switch': (TdopSymbolSwitch, "TDOPTOKEN.SWITCH", 'switch')
		'case': (TdopSymbolStmtIdentifier, "TDOPTOKEN.CASE", 'case')
		'default': (TdopSymbolStmtIdentifier, "TDOPTOKEN.DEFAULT", 'default')
		'break': (TdopSymbolBreakContinue, "TDOPTOKEN.BREAK", 'break')
		'continue': (TdopSymbolBreakContinue, "TDOPTOKEN.CONTINUE", 'continue')
		'try': (TdopSymbolTry, "TDOPTOKEN.TRY", 'try')
		'catch': (TdopSymbolStmtIdentifier, "TDOPTOKEN.CATCH", 'catch')
		'throw': (TdopSymbolThrow, "TDOPTOKEN.THROW", 'throw')
		'return': (TdopSymbolReturn, "TDOPTOKEN.RETURN", 'return')
		)
	}
