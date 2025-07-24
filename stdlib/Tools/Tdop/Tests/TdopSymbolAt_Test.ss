// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('@a',
			'[LIST, [STMT, [ATOP, AT, ADD, NUMBER, IDENTIFIER(a)], SEMICOLON]]',
			[1, [1, [1, 1, -1, -1, 2], -1]])
		.CheckTdop('@+1a',
			'[LIST, [STMT, [ATOP, AT, ADD, NUMBER(1), IDENTIFIER(a)], SEMICOLON]]',
			#(1, (1, (1, 1, 2, 3, 4), -1)))
		.CheckTdop('@1',
			'[LIST, [STMT, [ATOP, AT, ADD, NUMBER, NUMBER(1)], SEMICOLON]]',
			#(1, (1, (1, 1, -1, -1 2), -1)))
		.CheckTdop('@+1#(1, a: 2)',
			'[LIST, [STMT, [ATOP, AT, ADD, NUMBER(1), [OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_MEMBER, NUMBER(1), COMMA], ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, NUMBER(2), COMMA]], RPAREN]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 2, 3, [4, 4, 5, [6,
				[6, 6, 7],
				[9, 9, 10, 12, -1]], 13]],
				-1]])

		.CheckTdopCatch('@+a', 'expected NUMBER, but got IDENTIFIER(a)')
		}
	}