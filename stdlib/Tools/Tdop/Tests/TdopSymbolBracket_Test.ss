// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_Subscript()
		{
		.CheckTdop('a[b]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(a), LBRACKET, IDENTIFIER(b), ' $
				'RBRACKET], SEMICOLON]]',
			[1, [1, [1, 1, 2, 3,
				4], -1]])
		.CheckTdop('c*a[b]',
			'[LIST, [STMT, ' $
				'[BINARYOP, IDENTIFIER(c), MUL, ' $
					'[SUBSCRIPT, IDENTIFIER(a), LBRACKET, IDENTIFIER(b), RBRACKET]], ' $
				'SEMICOLON]]',
			[1, [1,
				[1, 1, 2,
					[3, 3, 4, 5, 6]],
				-1]])
		.CheckTdop('a[b].c',
			'[LIST, [STMT, [MEMBEROP, ' $
				'[SUBSCRIPT, IDENTIFIER(a), LBRACKET, IDENTIFIER(b), RBRACKET], ' $
				'DOT, IDENTIFIER(c)], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, 2, 3, 4],
				5, 6], -1]])
		.CheckTdop('a.c[b]',
			'[LIST, [STMT, [SUBSCRIPT, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(c)], ' $
				'LBRACKET, IDENTIFIER(b), RBRACKET], SEMICOLON]]',
			[1, [1, [1, [1, 1, 2, 3],
				4, 5, 6], -1]])

		.CheckTdop('a[c][b]',
			'[LIST, [STMT, [SUBSCRIPT, ' $
				'[SUBSCRIPT, IDENTIFIER(a), LBRACKET, IDENTIFIER(c), RBRACKET], ' $
				'LBRACKET, IDENTIFIER(b), RBRACKET], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, 2, 3, 4],
				5, 6, 7], -1]])

		.CheckTdop('ob[a..b]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), ' $
				'LBRACKET, [RANGE, IDENTIFIER(a), RANGETO, IDENTIFIER(b)], RBRACKET], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1,
				3, [4, 4, 5, 7], 8],
				-1]])
		.CheckTdop('ob[1+2::2+1]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), LBRACKET, ' $
				'[RANGE, ' $
					'[BINARYOP, NUMBER(1), ADD, NUMBER(2)], ' $
					'RANGELEN, ' $
					'[BINARYOP, NUMBER(2), ADD, NUMBER(1)]], RBRACKET], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4,
					[4, 4, 5, 6],
					7,
					[9, 9, 10, 11]], 12], -1]])
		.CheckTdop('ob[..]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), LBRACKET, ' $
				'[RANGE, NUMBER(0), RANGETO, NUMBER(2147483647)], RBRACKET], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4, -1, 4, -1], 6], -1]])
		.CheckTdop('ob[1..]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), LBRACKET, ' $
				'[RANGE, NUMBER(1), RANGETO, NUMBER(2147483647)], RBRACKET], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4, 4, 5, -1], 7], -1]])
		.CheckTdop('ob[..1]',
			'[LIST, [STMT, [SUBSCRIPT, IDENTIFIER(ob), LBRACKET, ' $
				'[RANGE, NUMBER(0), RANGETO, NUMBER(1)], RBRACKET], SEMICOLON]]',
			[1, [1, [1, 1, 3,
				[4, -1, 4, 6], 7], -1]])
		}

	Test_Record()
		{
		.CheckTdop('[a b c]',
			'[LIST, [STMT, [CALL, IDENTIFIER(Record), LBRACKET, [LIST, ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA], ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(b)], COMMA], ' $
				'[ARG_ELEM, [ARG, IDENTIFIER(c)], COMMA]], RBRACKET, BLOCK], SEMICOLON]]',
			[1, [1, [1, -1, 1, [2,
				[2, [2, 2], -1],
				[4, [4, 4], -1],
				[6, [6, 6], -1]], 7, -1], -1]])
		.CheckTdopCatch('[@ob]')
		}
	}
