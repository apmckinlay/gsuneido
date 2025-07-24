// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_FunctionCall()
		{
		.CheckTdop('A.Fn().B',
			'[LIST, [STMT, [MEMBEROP, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(Fn)], ' $
					'LPAREN, LIST, RPAREN, BLOCK], ' $
				'DOT, IDENTIFIER(B)], SEMICOLON]]',
			[1, [1, [1,
				[1, [1, 1, 2, 3],
					5, -1, 6, -1]
				7, 8], -1]])

		.CheckTdop('A*Fn()[1]',
			'[LIST, [STMT, [BINARYOP, ' $
				'IDENTIFIER(A), MUL, ' $
				'[SUBSCRIPT, ' $
					'[CALL, IDENTIFIER(Fn), LPAREN, LIST, RPAREN, BLOCK], ' $
					'LBRACKET, NUMBER(1), RBRACKET]], SEMICOLON]]')
		.CheckTdop('fn(1, :b, a: 3)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, NUMBER(1)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, NUMBER(3)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, 4], 5],
				[7, [7, -1, 7, 8], 9],
				[11, [11, 11, 12, 14], -1]]
				15, -1], -1]])
		.CheckTdop('fn(1, c:, :b, a: 3)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, NUMBER(1)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(c), COLON, TRUE], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, NUMBER(3)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, 4], 5],
				[7, [7, 7, 8, -1], 9],
				[11, [11, -1, 11, 12], 13],
				[15, [15, 15, 16, 18], -1]],
				19, -1], -1]])
		.CheckTdop('fn(d = 1, c:, :b, a: 3)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(d), EQ, NUMBER(1)]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(c), COLON, TRUE], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, NUMBER(3)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6, 8]], 9],
				[11, [11, 11, 12, -1], 13],
				[15, [15, -1, 15, 16], 17],
				[19, [19, 19, 20, 22], -1]],
				23, -1], -1]])
		.CheckTdop('fn(d = 1 :b a: 3)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(d), EQ, NUMBER(1)]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, NUMBER(3)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6, 8]], -1],
				[10, [10, -1, 10, 11], -1],
				[13, [13, 13, 14, 16], -1]],
				17, -1], -1]])
		.CheckTdop('fn(d = 1 :b :a)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(d), EQ, NUMBER(1)]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, IDENTIFIER(a)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6, 8]], -1],
				[10, [10, -1, 10, 11], -1],
				[13, [13, -1, 13, 14], -1]],
				15, -1], -1]])
		.CheckTdop('fn(d = 1 c: a: 3)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(d), EQ, NUMBER(1)]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(c), COLON, TRUE], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(a), COLON, NUMBER(3)], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6, 8]], -1],
				[10, [10, 10, 11, -1], -1],
				[13, [13, 13, 14, 16], -1]]
				17, -1], -1]])
		.CheckTdop('fn(true)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, 4], -1]],
				8, -1], -1]])
		.CheckTdop('fn(@ob)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, ' $
				'[ATOP, AT, ADD, NUMBER, IDENTIFIER(ob)], RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4, 4, -1, -1, 5], 7, -1], -1]])
		.CheckTdop('fn(@+1ob)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, ' $
				'[ATOP, AT, ADD, NUMBER(1), IDENTIFIER(ob)], RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4, 4, 5, 6, 7], 9, -1], -1]])
		.CheckTdop('fn(a, b:c){}',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(b), COLON, IDENTIFIER(c)], COMMA]], ' $
				'RPAREN, [BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('fn(){}',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('fn()\r\n{}',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('if fn(){} {b}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'ELSE, STMTS]]')
		.CheckTdop('if fn()\r\n{} {b}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK], RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY], ELSE, STMTS], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]]')

		.CheckTdop('if a\r\n(b)\r\nc',
			'[LIST, [IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'[STMT, [RVALUE, LPAREN, IDENTIFIER(b), RPAREN], SEMICOLON], ' $
				'ELSE, STMTS], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]]')
		.CheckTdop('if a\r\n[b]\r\nc',
			'[LIST, ' $
				'[IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
					'[STMT, [CALL, IDENTIFIER(Record), LBRACKET, [LIST, ' $
						'[ARG_ELEM, [ARG, IDENTIFIER(b)], COMMA]], RBRACKET, BLOCK], ' $
						'SEMICOLON], ' $
					'ELSE, STMTS], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]]')

		.CheckTdopCatch('fn(1, @ob)')
		.CheckTdopCatch('fn(@ob 1)')
		.CheckTdopCatch('fn(@ob :a)')
		.CheckTdopCatch('fn(1, ,2)')
		.CheckTdopCatch('fn(a: :b)')
		.CheckTdopCatch('fn(:1)', 'Invalid argument list')
		.CheckTdopCatch('fn(:"1")', 'Invalid argument list')
		.CheckTdopCatch('fn(a: 1 4)',
			"un-named arguments must come before named arguments")
		}
	}
