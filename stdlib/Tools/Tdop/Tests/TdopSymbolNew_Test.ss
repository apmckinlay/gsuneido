// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('new A',
			'[LIST, [STMT, ' $
				'[NEWOP, NEW, IDENTIFIER(A), LPAREN, LIST, RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 5, -1, -1, -1, -1], -1]])
		.CheckTdop('new A()',
			'[LIST, [STMT, ' $
				'[NEWOP, NEW, IDENTIFIER(A), LPAREN, LIST, RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 5, 6, -1, 7, -1], -1]])
		.CheckTdop('new A[1]()',
			'[LIST, [STMT, [NEWOP, NEW, ' $
				'[SUBSCRIPT, IDENTIFIER(A), LBRACKET, NUMBER(1), RBRACKET], ' $
					'LPAREN, LIST, RPAREN, BLOCK], SEMICOLON]]')
		.CheckTdop('new A.B()',
			'[LIST, [STMT, [NEWOP, NEW, ' $
				'[MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(B)], ' $
				'LPAREN, LIST, RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1,
				[5, 5, 6, 7],
				8, -1, 9, -1], -1]])
		.CheckTdop('new A = B()',
			'[LIST, [STMT, [NEWOP, NEW, ' $
				'[BINARYOP, IDENTIFIER(A), EQ, ' $
					'[CALL, IDENTIFIER(B), LPAREN, LIST, RPAREN, BLOCK]], ' $
				'LPAREN, LIST, RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1,
				[5, 5, 7,
					[9, 9, 10, -1, 11, -1]],
				-1, -1, -1, -1], -1]])
		.CheckTdop('a*new b.c(){d}*e',
			'[LIST, [STMT, [BINARYOP, ' $
				'[BINARYOP, ' $
					'IDENTIFIER(a), MUL, ' $
					'[NEWOP, NEW, [MEMBEROP, IDENTIFIER(b), DOT, IDENTIFIER(c)], ' $
						'LPAREN, LIST, RPAREN, ' $
						'[BLOCK, LCURLY, BITOR, LIST, BITOR, ' $
							'[LIST, [STMT, IDENTIFIER(d), SEMICOLON]], RCURLY]]], ' $
				'MUL, IDENTIFIER(e)], SEMICOLON]]')
		}
	}