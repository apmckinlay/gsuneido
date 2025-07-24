// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('do {} while a',
			'[LIST, [DOSTMT, DO, ' $
				'[STMTS, LCURLY, LIST, RCURLY], ' $
				'WHILE, LPAREN, IDENTIFIER(a), RPAREN]]')
		.CheckTdop('do a while fn()\r\n{b}',
			'[LIST, [DOSTMT, DO, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'WHILE, LPAREN, ' $
					'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK], RPAREN], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]]')
		.CheckTdop('do a; while fn()\r\n{b}',
			'[LIST, [DOSTMT, DO, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'WHILE, LPAREN, ' $
					'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK], RPAREN], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]]')
		.CheckTdop('do
			{
			a = 10
			do
				a--
			while a > 0
			} while true',
			'[LIST, [DOSTMT, DO, [STMTS, LCURLY, [LIST, ' $
				'[STMT, [BINARYOP, IDENTIFIER(a), EQ, NUMBER(10)], SEMICOLON], ' $
				'[DOSTMT, DO, ' $
					'[STMT, [POSTINCDEC, IDENTIFIER(a), DEC], SEMICOLON], ' $
					'WHILE, LPAREN, ' $
						'[BINARYOP, IDENTIFIER(a), GT, NUMBER(0)], ' $
					'RPAREN]], RCURLY], WHILE, LPAREN, TRUE, RPAREN]]',
			[1, [1, 1, [8, 8, [14,
				[14, [14, 14, 16, 18], -1],
				[25, 25,
					[33, [33, 33, 34], -1],
					41, -1,
						[47, 47, 49, 51],
					-1]], 57], 59, -1, 65, -1]])
		}
	}