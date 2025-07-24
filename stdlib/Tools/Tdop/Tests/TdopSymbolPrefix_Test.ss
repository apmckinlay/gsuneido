// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('! ~a',
			'[LIST, [STMT, [UNARYOP, NOT, [UNARYOP, BITNOT, IDENTIFIER(a)]], SEMICOLON]]',
			[1, [1, [1, 1, [3, 3, 4]], -1]])
		.CheckTdop('a or not b',
			'[LIST, [STMT, ' $
				'[BINARYOP, IDENTIFIER(a), OR, [UNARYOP, NOT, IDENTIFIER(b)]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 3, [6, 6, 10]], -1]])
		.CheckTdop('a or\n not b',
			'[LIST, [STMT, ' $
				'[BINARYOP, IDENTIFIER(a), OR, [UNARYOP, NOT, IDENTIFIER(b)]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 3, [7, 7, 11]], -1]])
		.CheckTdop('a or not (b and c)',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), OR, ' $
				'[UNARYOP, NOT, ' $
					'[RVALUE, LPAREN, ' $
						'[BINARYOP, IDENTIFIER(b), AND, IDENTIFIER(c)], RPAREN]]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[6, 6,
					[10, 10,
						[11, 11, 13, 17], 18]]], -1]])
		.CheckTdop('a or not b and c',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), OR, ' $
				'[BINARYOP, [UNARYOP, NOT, IDENTIFIER(b)], AND, IDENTIFIER(c)]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[6, [6, 6, 10], 12, 16]], -1]])
		.CheckTdop('a and not b or c',
			'[LIST, [STMT, [BINARYOP, ' $
				'[BINARYOP, IDENTIFIER(a), AND, [UNARYOP, NOT, IDENTIFIER(b)]], ' $
				'OR, ' $
				'IDENTIFIER(c)], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, 3, [7, 7, 11]],
				13,
				16], -1]])
		}
	}