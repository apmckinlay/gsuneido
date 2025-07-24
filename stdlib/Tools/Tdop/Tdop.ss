// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// Top Down Operator Precedence parser
// based on: Top Down Operator Precedence by Douglas Crockford in Beautiful Code
// which is in turn based on the paper with the same name
// by Vaughan R. Pratt

// stmts:
//	- SEMICOLON
//	- STMTS [LCURLY, LIST{stmt}, RCURLY]
//	- stmt

// stmt:
//	- STMT [NIL, SEMICOLON]
//	- STMT [expr, SEMICOLON]
//	- IFSTMT [IF, LPAREN, expr, RPAREN, stmts, ELSE, stmts]
//	- WHILESTMT [WHILE, LPAREN, expr, RPAREN, stmts]
//	- DOSTMT [DO, stmts, WHILE, LPAREN, expr, RPAREN]
//	- SWITCHSTMT [SWITCH, LPAREN, expr, RPAREN, LCURLY, LIST{CASE_ELEM}, RCURLY]
//		- CASE_ELEM [CASE|DEFAULT, LIST{EXPR_ELEM}, COLON, LIST{stmt}]
//			- EXPR_ELEM [expr, COMMA]
//	- FORSTMT [FOR, LPAREN, LIST{EXPR_ELEM}, SEMICOLON, expr, SEMICOLON, LIST{EXPR_ELEM}, RPAREN, stmts]
//	- FORINSTMT [FOR, LPAREN, IDENTIFIER, IN, expr|RANGE, RPAREN, stmts]
//	- FOREVERSTMT [FOREVER, stmts]
//	- RETURNSTMT [RETURN, THROW, LIST{EXPR_ELEM}, SEMICOLON]
//	- BREAKCONTINUESTMT [BREAK|CONTINUE, SEMICOLON]
//	- THROWSTMT [THROW, expr, SEMICOLON]
//	- TRYSTMT [TRY, stmts, CATCHSTMT]
//		- CATCHSTMT [CATCH, CATCH_COND, stmts]
//			- CATCH_COND [LPAREN, IDENTIFIER, COMMA, STRING, RPAREN]
//  - MULTIASSIGNSTMT [LIST{EXPR_ELEM}, EQ, expr]

// expr:
//	- constant
//	- IDENTIFIER
//	- MEMBEROP [expr|SELFREF, DOT, IDENTIFIER]
//	- TRINARYOP [expr, Q_MARK, expr, COLON, expr]
//	- BINARYOP [expr, bop, expr]
//	- UNARYOP [uop, expr]
//	- INOP [expr, IN, LPAREN, LIST{EXPR_ELEM}, PRAREN]
// 	- NOTINOP [expr, NOT, IN, LPAREN, LIST{EXPR_ELEM}, PRAREN]
//	- RVALUE [LPAREN, expr, RPAREN]
//	- CALL [expr, LPAREN, ATOP|LIST{ARG_ELEM}, RPAREN, BLOCK]
//		- ATOP [AT, ADD, NUMBER, expr]
// 		- ARG_ELEM [KEYARG|ARG, COMMA]
//			- KEYARG [STRING|NUMBER, COLON, expr]
//			- ARG [expr]
//	- PREINCDEC [INC|DEC, expr]
//	- POSTINCDEC [expr, INC|DEC]
//	- NEWOP [NEW, expr, LPAREN, ATOP|LIST{ARG_ELEM}, RPAREN, BLOCK]
//	- BLOCK [LCURLY, BITOR, BPAREM_AT|LIST{BPAREM}, BITOR, LIST{stmt}, RCURLY]
//		- BPAREM_AT [AT, IDENTIFIER]
//		- BPAREM [IDENTIFIER, COMMA]
//	- SUBSCRIPT [expr, LBRACKET, RANGE|expr, RBRACKET]
//		- RANGE [expr, RANGETO|RANGELEN, expr]

// constant:
//	- STRING
//	- BINARYOP [STRING, CAT, STRING]
//	- NUMBER
//	- UNARYOP [ADD|SUB, NUMBER]
//	- DATE
//	- SYMBOL
//	- TRUE
//	- FALSE
// 	- OBJECT [HASH, LPAREN, LIST(CONST_MEMBER|CONST_KEYMEMBER), RPAREN]
//		- CONST_KEYMEMBER [STRING|NUMBER, COLON, constant, COMMA|SEMICOLON]
//		- CONST_MEMBER [constant, COMMA|SEMICOLON]
// 	- RECORD [HASH, LCURLY|LBRACKET, LIST(CONST_MEMBER|CONST_KEYMEMBER), RCURLY|RBRACKET]
//	- CLASSDEF [CLASS, COLON, STRING, LCURLY, LIST(CONST_MEMBER), RCURLY]
//	- FUNCTIONDEF [FUNCTION, LPAREN, PAREM_AT|LIST{PAREM|PAREM_DEFAULT}, RPAREN, LCURLY, LIST{stmt}, RCURLY]
//		- PAREM_AT [AT, IDENTIFIER]
//		- PAREM [DOT, IDENTIFIER, COMMA]
//		- PAREM_DEFAULT [DOT, IDENTIFIER, EQ, const, COMMA]
//	- DLLDEF [DLL, IDENTIFIER, IDENTIFIER, COLON, STRING, LPAREN, LIST{DLL_PAREM}, RPAREN]
//		- DLL_PAREM [DLL_IN, DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, COMMA]
//			- DLL_IN [LBRACKET, IN, RBRACKET]
//			- DLL_NORMAL [IDENTIFIER]
//			- DLL_POINTER [IDENTIFIER, MUL]
//			- DLL_ARRAY [IDENTIFIER, LBRACKET, NUMBER, RBRACKET]
//	- STRUCTDEF [STRUCT, LCURLY, LIST{STRUCT_MEMBER}, RCURLY]
//		- STRUCT_MEMBER [DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, SEMICOLON]
//	- CALLBACKDEF [CALLBACK, LPAREN, LIST{CALLBACK_PAREM}, RPAREN]
//		- CALLBACK_PAREM [DLL_POINTER|DLL_ARRAY|DLL_NORMAL, IDENTIFIER, COMMA]

// bop:
//	- EQ
//	- ADDEQ
//	- SUBEQ
//	- CATEQ
//	- MULEQ
//	- DIVEQ
//	- MODEQ
//	- LSHIFTEQ
//	- RSHIFTEQ
//	- BITOREQ
//	- BITANDEQ
//	- BITXOREQ
//	- AND
//	- OR
//	- ADD
//	- SUB
//	- MUL
//	- DIV
//	- MOD
//	- CAT
//	- LT
//	- LTE
//	- GT
//	- GTE
//	- IS
//	- ISNT
//	- MATCH
//	- MATCHNOT
//	- BITOR
//	- BITAND
//	- BITXOR
//	- LSHIFT
//	- RSHIFT

// uop:
//	- NOT
//	- ADD
//	- SUB
//	- BITNOT
class
	{
	CallClass(src, type = 'constant', nodes = false, symbols = false)
		{
		parser = new this(src, symbols is false ? TdopSymbols() : symbols)
		_nodes = nodes is false ? Object() : nodes
		_expr = parser.Expression
		_stmt = parser.Statement
		_stmts = parser.Statements
		_token = parser.Scan.Token
		_advance = parser.Advance
		_ahead = parser.Scan.Ahead
		_isNewline = parser.Scan.IsNewline
		_getStmtnest = parser.GetStmtnest
		_setStmtnest = parser.SetStmtnest
		_end = parser.Scan.End
		_expectingCompound = false

		switch (type)
			{
		case 'constant':
			res = TdopConstant()
		case 'expression':
			res = parser.Expression()
		case 'statement':
			res = parser.Statement()
		case 'statements':
			res = parser.Statements()
		default:
			throw 'Unhandled Tdop Type'
			}
		Assert(parser.Scan.Token() is: parser.Scan.End)
		return res
		}

	New(src, .Symbols)
		{
		.Scan = TdopScanner(src, .Symbols)
		}

	stmtnest: 0
	GetStmtnest()
		{
		return .stmtnest
		}

	SetStmtnest(newStmtnest)
		{
		oldStmtnest = .stmtnest
		.stmtnest = newStmtnest
		return oldStmtnest
		}

	Advance(match = false)
		{
		newline = .Scan.IsNewline()
		token = .Scan.Token()
		end = .Scan.End
		if match is TDOPTOKEN.SEMICOLON and
			(token is end or token.Token is TDOPTOKEN.RCURLY or
				(token.Token isnt TDOPTOKEN.SEMICOLON and newline))
			return // implicit semicolon insertion
		if match isnt false and not token.Match(match)
			throw "expected " $ match $ " but got " $ Display(token)
		.updateStmtnest(token)
		.Scan.Advance()
		}

	updateStmtnest(token)
		{
		if token is false
			return

		if token.Token in (TDOPTOKEN.LBRACKET, TDOPTOKEN.LCURLY, TDOPTOKEN.LPAREN)
			.stmtnest++
		if token.Token in (TDOPTOKEN.RBRACKET, TDOPTOKEN.RCURLY, TDOPTOKEN.RPAREN)
			.stmtnest--
		}

	Expression(rbp = 0, expectingCompound = false)
		{
		if expectingCompound is true
			_expectingCompound = true
		t = .Scan.Token()
		.Advance()
		left = t.Nud()
		while rbp < .Scan.Token().Lbp
			{
			t = .Scan.Token()
			if t.Break(left)
				break
			.Advance()
			left = t.Led(left)
			}
		return left
		}

	// STMT [stmt, SEMICOLON]
	Statement(alt_end = false)
		{
		t = .Scan.Token()
		if t.Token is TDOPTOKEN.SEMICOLON
			{
			.Advance(TDOPTOKEN.SEMICOLON)
			return TdopCreateNode(TDOPTOKEN.STMT,
				children: Object(TdopCreateNode(TDOPTOKEN.NIL), t))
			}
		if t.Method?(#Std)
			{
			.Advance()
			return t.Std()
			}
		expr = TdopStmtExpr()
		if expr.Match(TDOPTOKEN.IDENTIFIER) and .Scan.Token().Match(TDOPTOKEN.COMMA)
			return TdopMultiAssignStmt(expr)

		children = Object()
		TdopAddChild(children, token: expr)
		TdopAddChild(children, match: TDOPTOKEN.SEMICOLON,
			mustMatch: not .Scan.Token().Match(alt_end), implicit:)
		return TdopCreateNode(TDOPTOKEN.STMT, :children)
		}

	Statements()
		{
		return TdopCreateList()
			{ |list|
			while .Scan.Token() isnt .Scan.End and
				not .Scan.Token().Match(TDOPTOKEN.RCURLY)
				list.Add(.Statement())
			}
		}
	}
