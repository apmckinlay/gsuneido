// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('++a',
			'[LIST, [STMT, [PREINCDEC, INC, IDENTIFIER(a)], SEMICOLON]]',
			[1, [1, [1, 1, 3], -1]])
		.CheckTdop('--a',
			'[LIST, [STMT, [PREINCDEC, DEC, IDENTIFIER(a)], SEMICOLON]]',
			[1, [1, [1, 1, 3], -1]])
		.CheckTdop('--a.b',
			'[LIST, [STMT, ' $
				'[PREINCDEC, DEC, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)]], ' $
				'SEMICOLON]]')
		.CheckTdop('--a*b',
			'[LIST, [STMT, ' $
				'[BINARYOP, [PREINCDEC, DEC, IDENTIFIER(a)], MUL, IDENTIFIER(b)], ' $
				'SEMICOLON]]',
			[1, [1, [1, [1, 1, 3], 4, 5], -1]])

		.CheckTdop('a++',
			'[LIST, [STMT, [POSTINCDEC, IDENTIFIER(a), INC], SEMICOLON]]',
			[1, [1, [1, 1, 2], -1]])
		.CheckTdop('a--',
			'[LIST, [STMT, [POSTINCDEC, IDENTIFIER(a), DEC], SEMICOLON]]',
			[1, [1, [1, 1, 2], -1]])
		.CheckTdop('a.b--',
			'[LIST, [STMT, ' $
				'[POSTINCDEC, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], DEC], ' $
				'SEMICOLON]]')
		.CheckTdop('a*b--',
			'[LIST, [STMT, ' $
				'[BINARYOP, IDENTIFIER(a), MUL, [POSTINCDEC, IDENTIFIER(b), DEC]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, 2, [3, 3, 4]], -1]])
		}
	}