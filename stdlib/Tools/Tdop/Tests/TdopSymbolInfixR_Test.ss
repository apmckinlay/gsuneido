// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a and b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 7], -1]])
		.CheckTdop('a \nand b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 4, 8], -1]])
		.CheckTdop('a and\n b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 8], -1]])
		.CheckTdop('a and b and c',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, ' $
				'[BINARYOP, IDENTIFIER(b), AND, IDENTIFIER(c)]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [7, 7, 9, 13]], -1]])
		.CheckTdop('a or b and c or d', '[LIST, [STMT, [BINARYOP, ' $
			'IDENTIFIER(a), OR, ' $
				'[BINARYOP, [BINARYOP, IDENTIFIER(b), AND, IDENTIFIER(c)], ' $
					'OR, IDENTIFIER(d)]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [6, [6, 6, 8, 12], 14, 17]], -1]])
		.CheckTdop('a and b or c and d',
			'[LIST, [STMT, ' $
				'[BINARYOP, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], OR, ' $
					'[BINARYOP, IDENTIFIER(c), AND, IDENTIFIER(d)]], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 7], 9, [12, 12, 14, 18]], -1]])
		.CheckTdop('1 > 2 and 3 > 4',
			'[LIST, [STMT, [BINARYOP, [BINARYOP, NUMBER(1), GT, NUMBER(2)], AND, ' $
				'[BINARYOP, NUMBER(3), GT, NUMBER(4)]], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 5], 7, [11, 11, 13, 15]], -1]])

		.CheckTdop('a = b = c',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), EQ, ' $
				'[BINARYOP, IDENTIFIER(b), EQ, IDENTIFIER(c)]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7, 9]], -1]])
		.CheckTdop('a < b = c',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), LT, ' $
				'[BINARYOP, IDENTIFIER(b), EQ, IDENTIFIER(c)]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7, 9]], -1]])
		.CheckTdop('a * b += c',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), MUL, ' $
				'[BINARYOP, IDENTIFIER(b), ADDEQ, IDENTIFIER(c)]], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[5, 5, 7, 10]], -1]])
		.CheckTdop('a * b.c = d() + 1',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), MUL, ' $
				'[BINARYOP, [MEMBEROP, IDENTIFIER(b), DOT, IDENTIFIER(c)], EQ, ' $
					'[BINARYOP, [CALL, IDENTIFIER(d), LPAREN, LIST, RPAREN, BLOCK], ' $
						'ADD, NUMBER(1)]]], SEMICOLON]]')
		.CheckTdop('a = b and c > 1',
			'[LIST, [STMT, ' $
				'[BINARYOP, IDENTIFIER(a), EQ, [BINARYOP, IDENTIFIER(b), AND, ' $
					'[BINARYOP, IDENTIFIER(c), GT, NUMBER(1)]]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7, [11, 11, 13 15]]], -1]])

		.CheckTdop('a or b or c and d and e',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), OR, [BINARYOP, IDENTIFIER(b), OR, ' $
				'[BINARYOP, IDENTIFIER(c), AND, ' $
					'[BINARYOP, IDENTIFIER(d), AND, IDENTIFIER(e)]]]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [6, 6, 8, [11, 11, 13, [17, 17, 19, 23]]]], -1]])

		.CheckTdop('a = b += c <<= d = e',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), EQ, ' $
				'[BINARYOP, IDENTIFIER(b), ADDEQ, ' $
					'[BINARYOP, IDENTIFIER(c), LSHIFTEQ, ' $
						'[BINARYOP, IDENTIFIER(d), EQ, IDENTIFIER(e)]]]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7, [10, 10, 12, [16, 16, 18, 20]]]], -1]])
		}
	}
