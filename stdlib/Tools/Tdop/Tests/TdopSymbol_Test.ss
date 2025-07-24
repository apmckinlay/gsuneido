// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_Semicolons()
		{
		.CheckTdop('a; b;',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, 2],
				[4, 4, 5]])
		.CheckTdop('a; b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, 2],
				[4, 4, -1]])
		.CheckTdop('a; \n b;',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, 2],
				[6, 6, 7]])
		.CheckTdop('a; \n b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, 2],
				[6, 6, -1]])
		.CheckTdop('a \n b;',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, -1],
				[5, 5, 6]])
		.CheckTdop('a \n b',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, -1],
				[5, 5, -1]])
		}
	}