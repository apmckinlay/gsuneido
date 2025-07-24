// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('+a',
			'[LIST, [STMT, [UNARYOP, ADD, IDENTIFIER(a)], SEMICOLON]]',
			[1, [1, [1, 1, 2], -1]])
		.CheckTdop('+-a',
			'[LIST, [STMT, [UNARYOP, ADD, [UNARYOP, SUB, IDENTIFIER(a)]], SEMICOLON]]',
			[1, [1, [1, 1, [2, 2, 3]], -1]])
		.CheckTdop('+a.b',
			'[LIST, [STMT, ' $
				'[UNARYOP, ADD, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)]], ' $
				'SEMICOLON]]',
			[1, [1, [1, 1, [2, 2, 3, 4]], -1]])
		.CheckTdop('+a * b',
			'[LIST, [STMT, ' $
				'[BINARYOP, [UNARYOP, ADD, IDENTIFIER(a)], MUL, IDENTIFIER(b)], ' $
				'SEMICOLON]]',
			[1, [1, [1, [1, 1, 2], 4, 6], -1]])
		}
	}