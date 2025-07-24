// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('Object(forever:)',
			'[LIST, [STMT, [CALL, IDENTIFIER(Object), LPAREN, ' $
				'[LIST, [ARG_ELEM, [KEYARG, STRING(forever), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]')

		.CheckTdop('forever {}',
			'[LIST, [FOREVERSTMT, FOREVER, [STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('forever;',
			'[LIST, [FOREVERSTMT, FOREVER, SEMICOLON]]')
		.CheckTdop('forever a',
			'[LIST, [FOREVERSTMT, FOREVER, [STMT, IDENTIFIER(a), SEMICOLON]]]')
		.CheckTdop('forever a \r\nb',
			'[LIST, ' $
				'[FOREVERSTMT, FOREVER, [STMT, IDENTIFIER(a), SEMICOLON]], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1,
				[1, 1, [9, 9, -1]],
				[13, 13, -1]])
		}
	}