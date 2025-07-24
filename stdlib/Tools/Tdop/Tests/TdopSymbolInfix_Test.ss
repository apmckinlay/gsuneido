// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a + b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), ADD, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 5], -1]])
		.CheckTdop('a + b + c',
			'[LIST, [STMT, [BINARYOP, ' $
				'[BINARYOP, IDENTIFIER(a), ADD, IDENTIFIER(b)], ADD, IDENTIFIER(c)], ' $
				'SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 5], 7, 9], -1]])
		.CheckTdop('a + b * c', '[LIST, [STMT, [BINARYOP, IDENTIFIER(a), ADD, ' $
			'[BINARYOP, IDENTIFIER(b), MUL, IDENTIFIER(c)]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7 9]], -1]])
		.CheckTdop('(a + b) * c',
			'[LIST, [STMT, [BINARYOP, ' $
				'[RVALUE, LPAREN, ' $
					'[BINARYOP, IDENTIFIER(a), ADD, IDENTIFIER(b)], RPAREN], ' $
				'MUL, IDENTIFIER(c)], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, [2, 2, 4, 6], 7],
				9, 11], -1]])

		.CheckTdop('a + \n b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), ADD, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 7], -1]])
		.CheckTdop('a \n + b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, [UNARYOP, ADD, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, 1, -1], [5, [5, 5, 7], -1]])

		.CheckTdop('a $ b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), CAT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 5], -1]])
		.CheckTdop('a $ \n b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), CAT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 7], -1]])
		.CheckTdop('a \n $ b',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), CAT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 5, 7], -1]])

		.CheckTdop('1<=2',
			'[LIST, [STMT, [BINARYOP, NUMBER(1), LTE, NUMBER(2)], SEMICOLON]]',
			[1, [1, [1, 1, 2, 4], -1]])
		.CheckTdop('1\n<=2',
			'[LIST, [STMT, [BINARYOP, NUMBER(1), LTE, NUMBER(2)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 5], -1]])
		.CheckTdop('1<=\n2',
			'[LIST, [STMT, [BINARYOP, NUMBER(1), LTE, NUMBER(2)], SEMICOLON]]',
			[1, [1, [1, 1, 2, 5], -1]])

		.CheckTdop('a | b ^ c & ~d',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), BITOR, ' $
				'[BINARYOP, IDENTIFIER(b), BITXOR, ' $
					'[BINARYOP, IDENTIFIER(c), BITAND, ' $
						'[UNARYOP, BITNOT, IDENTIFIER(d)]]]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [5, 5, 7, [9, 9, 11, [13, 13, 14]]]], -1]])
		.CheckTdop('a | b is c < d << e + f * g.h[i]',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), BITOR, ' $
				'[BINARYOP, IDENTIFIER(b), IS, ' $
					'[BINARYOP, IDENTIFIER(c), LT, ' $
						'[BINARYOP, IDENTIFIER(d), LSHIFT, ' $
							'[BINARYOP, IDENTIFIER(e), ADD, ' $
								'[BINARYOP, IDENTIFIER(f), MUL, ' $
									'[SUBSCRIPT, ' $
										'[MEMBEROP, ' $
											'IDENTIFIER(g), DOT, IDENTIFIER(h)], ' $
										'LBRACKET, IDENTIFIER(i), RBRACKET]]]]]]], ' $
				'SEMICOLON]]')
		}
	}