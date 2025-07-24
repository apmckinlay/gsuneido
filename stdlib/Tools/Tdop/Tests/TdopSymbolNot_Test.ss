// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a and b not in ()',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, ' $
				'[NOTINOP, IDENTIFIER(b), NOT, IN, LPAREN, LIST, RPAREN]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [7, 7, 9, 13, 16, -1, 17]], -1]])
		.CheckTdop('a | b not in ()',
			'[LIST, [STMT, [NOTINOP, ' $
				'[BINARYOP, IDENTIFIER(a), BITOR, IDENTIFIER(b)], ' $
				'NOT, IN, ' $
				'LPAREN, ' $
				'LIST, ' $
				'RPAREN], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 5], 7, 11 14, -1, 15], -1]])
		.CheckTdop('fn(a not: b not in (c, d))',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(not), COLON, ' $
					'[NOTINOP, ' $
						'IDENTIFIER(b), ' $
						'NOT, IN, ' $
						'LPAREN, ' $
						'[LIST, ' $
							'[EXPR_ELEM, IDENTIFIER(c), COMMA], ' $
							'[EXPR_ELEM, IDENTIFIER(d), COMMA]], ' $
						'RPAREN]], COMMA]], RPAREN, BLOCK], SEMICOLON]]')
		.CheckTdop('1 \r\nnot in (1, 2)',
			'[LIST, ' $
				'[STMT, NUMBER(1), SEMICOLON], ' $
				'[STMT, [UNARYOP, NOT, [CALL, IDENTIFIER(in), LPAREN, [LIST, ' $
					'[ARG_ELEM, [ARG, NUMBER(1)], COMMA], ' $
					'[ARG_ELEM, [ARG, NUMBER(2)], COMMA]], RPAREN, BLOCK]], SEMICOLON]]')
		}
	}