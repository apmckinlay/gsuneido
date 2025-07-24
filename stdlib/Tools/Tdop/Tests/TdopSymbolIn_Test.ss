// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a in (1, b, "c")',
			'[LIST, [STMT, [INOP, ' $
				'IDENTIFIER(a), ' $
				'IN, ' $
				'LPAREN, ' $
				'[LIST, ' $
					'[EXPR_ELEM, NUMBER(1), COMMA], ' $
					'[EXPR_ELEM, IDENTIFIER(b), COMMA], ' $
					'[EXPR_ELEM, STRING(c), COMMA]], ' $
				'RPAREN], SEMICOLON]]',
			[1, [1, [1,
				1,
				3,
				6,
				[7,
					[7, 7, 8],
					[10, 10, 11],
					[13, 13, -1]],
				16], -1]])
		.CheckTdop('a and b in ()',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(a), AND, ' $
				'[INOP, IDENTIFIER(b), IN, LPAREN, LIST, RPAREN]], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[7, 7, 9, 12, -1, 13]], -1]])
		.CheckTdop('a | b in ()',
			'[LIST, [STMT, [INOP, ' $
				'[BINARYOP, IDENTIFIER(a), BITOR, IDENTIFIER(b)], ' $
				'IN, LPAREN, LIST, RPAREN], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, 3 5],
				7, 10, -1, 11], -1]])
		.CheckTdop('fn(a, in: b in (c, d))',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA], ' $
				'[ARG_ELEM, [KEYARG, STRING(in), COLON, ' $
					'[INOP, ' $
						'IDENTIFIER(b), ' $
						'IN, ' $
						'LPAREN, ' $
						'[LIST, ' $
							'[EXPR_ELEM, IDENTIFIER(c), COMMA], ' $
							'[EXPR_ELEM, IDENTIFIER(d), COMMA]], ' $
						'RPAREN]], COMMA]], RPAREN, BLOCK], SEMICOLON]]')
		.CheckTdop('1 \r\nin (1, 2)',
			'[LIST, ' $
				'[STMT, NUMBER(1), SEMICOLON], ' $
				'[STMT, [CALL, IDENTIFIER(in), LPAREN, [LIST, ' $
					'[ARG_ELEM, [ARG, NUMBER(1)], COMMA], ' $
					'[ARG_ELEM, [ARG, NUMBER(2)], COMMA]], RPAREN, BLOCK], SEMICOLON]]')
		}
	}