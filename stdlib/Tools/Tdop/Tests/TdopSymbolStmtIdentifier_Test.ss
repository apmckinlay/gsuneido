// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('fn(break+1 :in continue:)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BINARYOP, IDENTIFIER(break), ADD, NUMBER(1)]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(in), COLON, IDENTIFIER(in)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(continue), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4, 9, 10]], -1],
				[12, [12, -1, 12, 13], -1],
				[16, [16, 16, 24, -1], -1]],
				25, -1], -1]])
		.CheckTdop('fn(is: true)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [KEYARG, STRING(is), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, 4, 6, 8], -1]], 12, -1], -1]])
		.CheckTdop('fn(a is true is: true)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(a), IS, TRUE]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(is), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6, 9]], -1],
				[14, [14, 14, 16, 18], -1]], 22, -1], -1]])
		.CheckTdop('fn(a and b and c and: true)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, [BINARYOP, IDENTIFIER(a), AND, ' $
					'[BINARYOP, IDENTIFIER(b), AND, IDENTIFIER(c)]]], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(and), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, [4, 4, 6,
					[10, 10, 12, 16]]], -1],
				[18, [18, 18, 21, 23], -1]], 27, -1], -1]])
		.CheckTdop('fn(a and: true)',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(and), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4, 4], -1],
				[6, [6, 6, 9, 11], -1]], 15, -1], -1]])
		.CheckTdop('break;continue',
			'[LIST, ' $
				'[BREAKCONTINUESTMT, BREAK, SEMICOLON], ' $
				'[BREAKCONTINUESTMT, CONTINUE, SEMICOLON]]')
		.CheckTdop('break;continue;;;',
			'[LIST, ' $
				'[BREAKCONTINUESTMT, BREAK, SEMICOLON], ' $
				'[BREAKCONTINUESTMT, CONTINUE, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON]]')
		}
	}
