// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('.b',
			'[LIST, [STMT, [MEMBEROP, SELFREF, DOT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, -1, 1, 2], -1]])
		.CheckTdop('a.b',
			'[LIST, [STMT, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 2, 3], -1]])
		.CheckTdop('a .b',
			'[LIST, [STMT, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 4], -1]])
		.CheckTdop('a . b',
			'[LIST, [STMT, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, 1, 3, 5], -1]])
		.CheckTdop('a.b.c',
			'[LIST, [STMT, [MEMBEROP, ' $
				'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], DOT, IDENTIFIER(c)], ' $
				'SEMICOLON]]',
			[1, [1, [1, [1, 1, 2, 3], 4, 5], -1]])
		.CheckTdop('a.\nb.c',
			'[LIST, [STMT, [MEMBEROP, ' $
				'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], DOT, IDENTIFIER(c)], ' $
				'SEMICOLON]]',
			[1, [1, [1, [1, 1, 2, 4], 5, 6], -1]])
		.CheckTdop('a\n.b.c',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, [MEMBEROP, [MEMBEROP, SELFREF, DOT, IDENTIFIER(b)], ' $
					'DOT, IDENTIFIER(c)], SEMICOLON]]',
			[1, [1, 1, -1], [3, [3, [3, -1, 3, 4], 5, 6], -1]])

		.CheckTdop('.a * b', '[LIST, [STMT, [BINARYOP, ' $
			'[MEMBEROP, SELFREF, DOT, IDENTIFIER(a)], MUL, IDENTIFIER(b)], SEMICOLON]]',
			[1, [1, [1, [1, -1, 1, 2], 4, 6], -1]])
		.CheckTdop('a.a * b', '[LIST, [STMT, [BINARYOP, ' $
			'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(a)], MUL, IDENTIFIER(b)],' $
				' SEMICOLON]]',
			[1, [1, [1, [1, 1, 2, 3], 5, 7], -1]])
		.CheckTdop('not a.a',
			'[LIST, [STMT, [UNARYOP, NOT, [MEMBEROP, IDENTIFIER(a), DOT, ' $
				'IDENTIFIER(a)]], SEMICOLON]]',
			[1, [1, [1, 1, [5, 5, 6, 7]], -1]])
		.CheckTdop('a.new.class.function',
			'[LIST, [STMT, ' $
				'[MEMBEROP, ' $
					'[MEMBEROP, ' $
						'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(new)], ' $
						'DOT, IDENTIFIER(class)], ' $
					'DOT, IDENTIFIER(function)], SEMICOLON]]',
			[1, [1, [1, [1, [1, 1, 2, 3], 6, 7], 12, 13], -1]])

		.CheckTdopCatch('a."1"')
		}
	}