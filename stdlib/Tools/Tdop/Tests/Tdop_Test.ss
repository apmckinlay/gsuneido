// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_BasicTypes()
		{
		.CheckTdop('abc', '[LIST, [STMT, IDENTIFIER(abc), SEMICOLON]]')

		.CheckTdop('123', 'NUMBER(123)', type: 'constant')
		.CheckTdop('.001', 'NUMBER(.001)', type: 'constant')
		.CheckTdop('1.123e10', 'NUMBER(1.123e10)', type: 'constant')
		.CheckTdop('0x1aF', 'NUMBER(0x1aF)', type: 'constant') //hex
		.CheckTdop('017', 'NUMBER(017)', type: 'constant')

		.CheckTdop('#20170620', 'DATE(20170620)', type: 'constant')
		.CheckTdop('#20170620.010101', 'DATE(20170620.010101)', type: 'constant')

		.CheckTdop('"abc"', 'STRING(abc)', type: 'constant')
		}

	Test_ComplexExpression()
		{
		.CheckTdop('!(a.b * (c + -12)) $ "d"',
			'[LIST, [STMT, [BINARYOP, ' $
				'[UNARYOP, NOT, [RVALUE, LPAREN, [BINARYOP, ' $
					'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], MUL, ' $
					'[RVALUE, LPAREN, ' $
						'[BINARYOP, IDENTIFIER(c), ADD, ' $
							'[UNARYOP, SUB, NUMBER(12)]], RPAREN]]' $
					', RPAREN]], ' $
				'CAT, STRING(d)], SEMICOLON]]')
		}

	Test_Statements()
		{
		.CheckTdop('if a is 1 Print(a); else Print(2); Print(3)',
			'[LIST, [IFSTMT, IF, ' $
				'LPAREN, [BINARYOP, IDENTIFIER(a), IS, NUMBER(1)], RPAREN, ' $
				'[STMT, [CALL, IDENTIFIER(Print), LPAREN, [LIST, ' $
					'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA]], RPAREN, BLOCK], ' $
					'SEMICOLON], ' $
				'ELSE, [STMT, [CALL, IDENTIFIER(Print), LPAREN, [LIST, ' $
					'[ARG_ELEM, [ARG, NUMBER(2)], COMMA]], RPAREN, BLOCK], ' $
					'SEMICOLON]], ' $
			'[STMT, [CALL, IDENTIFIER(Print), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, NUMBER(3)], COMMA]], RPAREN, BLOCK], SEMICOLON]]')
		.CheckTdop('a\r\nb',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]')
		.CheckTdop('a;b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]')
		.CheckTdop('a;;;;b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]')
		.CheckTdop('if true {a;;;} else b;;;;c',
			'[LIST, [IFSTMT, IF, LPAREN, TRUE, RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, IDENTIFIER(a), SEMICOLON], ' $
					'[STMT, NIL, SEMICOLON], ' $
					'[STMT, NIL, SEMICOLON]], RCURLY], ' $
				'ELSE, [STMT, IDENTIFIER(b), SEMICOLON]], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]]')
		}

	Test_MultiAssignStmt()
		{
		.CheckTdop('a, b, c = fn()',
			'[LIST, ' $
				'[MULTIASSIGNSTMT, ' $
					'[LIST, ' $
						'[EXPR_ELEM, IDENTIFIER(a), COMMA], ' $
						'[EXPR_ELEM, IDENTIFIER(b), COMMA], ' $
						'[EXPR_ELEM, IDENTIFIER(c), COMMA]], ' $
					'EQ, [CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK]]]',
			[1,
				[1,
					[1,
						[1, 1, 2],
						[4, 4, 5],
						[7, 7, -1]],
					9, [11, 11, 13, -1, 14, -1]]])

		.CheckTdopCatch('1, a = fn()')
		.CheckTdopCatch('fn(), a = fn()')
		.CheckTdopCatch('a, b = 12')
		}

	Test_StatementNestNewline()
		{
		.CheckTdop('true ? a\r\n+1 : b\r\n+1',
			'[LIST, [STMT, [TRINARYOP, TRUE, Q_MARK, ' $
				'[BINARYOP, IDENTIFIER(a), ADD, NUMBER(1)], COLON, IDENTIFIER(b)], ' $
				'SEMICOLON], ' $
				'[STMT, [UNARYOP, ADD, NUMBER(1)], SEMICOLON]]')
		.CheckTdop('(a\r\n.b)',
			'[LIST, [STMT, [RVALUE, LPAREN, ' $
				'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], RPAREN], SEMICOLON]]')
		.CheckTdop('[a: b\r\n*c]',
			'[LIST, [STMT, [CALL, IDENTIFIER(Record), LBRACKET, [LIST, ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, ' $
					'[BINARYOP, IDENTIFIER(b), MUL, IDENTIFIER(c)]], COMMA]], ' $
				'RBRACKET, BLOCK], SEMICOLON]]')
		.CheckTdop('ob[1\r\n+2\r\n::\r\na\r\n()]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), LBRACKET, [RANGE, ' $
				'[BINARYOP, NUMBER(1), ADD, NUMBER(2)], RANGELEN, ' $
				'[CALL, IDENTIFIER(a), LPAREN, LIST, RPAREN, BLOCK]], RBRACKET], ' $
				'SEMICOLON]]')
		.CheckTdop('(1\r\n+\r\na\r\n{})',
			'[LIST, [STMT, [RVALUE, LPAREN, ' $
				'[BINARYOP, NUMBER(1), ADD, ' $
				'[CALL, IDENTIFIER(a), LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]]], ' $
				'RPAREN], SEMICOLON]]')
		.CheckTdop('(1\r\n+\r\n(2\r\n+\r\n3)\r\n+\r\n4)',
			'[LIST, [STMT, [RVALUE, LPAREN, [BINARYOP, ' $
				'[BINARYOP, ' $
					'NUMBER(1), ' $
					'ADD, ' $
					'[RVALUE, LPAREN, ' $
						'[BINARYOP, NUMBER(2), ADD, NUMBER(3)], RPAREN]], ' $
				'ADD, NUMBER(4)], RPAREN], SEMICOLON]]')
		.CheckTdop('(a \r\nin (1, b not in (2, 3)))',
			'[LIST, [STMT, [RVALUE, LPAREN, ' $
				'[INOP, ' $
					'IDENTIFIER(a), IN, ' $
					'LPAREN, ' $
					'[LIST, ' $
						'[EXPR_ELEM, NUMBER(1), COMMA], ' $
						'[EXPR_ELEM, [NOTINOP, ' $
							'IDENTIFIER(b), NOT, IN, ' $
							'LPAREN, ' $
							'[LIST, ' $
								'[EXPR_ELEM, NUMBER(2), COMMA], ' $
								'[EXPR_ELEM, NUMBER(3), COMMA]], ' $
							'RPAREN], COMMA]], ' $
					'RPAREN], RPAREN], SEMICOLON]]')
		}
	}
