// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a and b ? 1 + 2 : 10 / 2',
			'[LIST, [STMT, [TRINARYOP, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], ' $
				'Q_MARK, ' $
				'[BINARYOP, NUMBER(1), ADD, NUMBER(2)], ' $
				'COLON, ' $
				'[BINARYOP, NUMBER(10), DIV, NUMBER(2)]], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 7],
				9,
				[11, 11, 13, 15],
				17,
				[19, 19, 22, 24]], -1]])
		.CheckTdop('a and b \n? 1 + 2 \n: 10 / 2',
			'[LIST, [STMT, [TRINARYOP, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], ' $
				'Q_MARK, ' $
				'[BINARYOP, NUMBER(1), ADD, NUMBER(2)], ' $
				'COLON, ' $
				'[BINARYOP, NUMBER(10), DIV, NUMBER(2)]], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 7],
				10,
				[12, 12, 14, 16],
				19,
				[21, 21, 24, 26]], -1]])
		.CheckTdop('a and b ?\n 1 + 2 :\n 10 / 2',
			'[LIST, [STMT, [TRINARYOP, [BINARYOP, IDENTIFIER(a), AND, IDENTIFIER(b)], ' $
				'Q_MARK, ' $
				'[BINARYOP, NUMBER(1), ADD, NUMBER(2)], ' $
				'COLON, ' $
				'[BINARYOP, NUMBER(10), DIV, NUMBER(2)]], SEMICOLON]]',
			[1, [1, [1, [1, 1, 3, 7],
				9,
				[12, 12, 14, 16],
				18,
				[21, 21, 24 26]], -1]])
		}
	}