// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('true; false;',
			'[LIST, [STMT, TRUE, SEMICOLON], [STMT, FALSE, SEMICOLON]]',
			[1, [1, 1, 5], [7, 7, 12]])
		.CheckTdop('Object(true, :true false:)',
			'[LIST, [STMT, [CALL, IDENTIFIER(Object), ' $
				'LPAREN, [LIST, ' $
					'[ARG_ELEM, [ARG, TRUE], COMMA], ' $
					'[ARG_ELEM, ' $
						'[KEYARG, STRING(true), COLON, IDENTIFIER(true)], COMMA], ' $
					'[ARG_ELEM, [KEYARG, STRING(false), COLON, TRUE], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1,
				7, [8,
					[8, [8, 8], 12],
					[14, [14, -1, 14, 15], -1],
					[20, [20, 20, 25, -1], -1]]
				26, -1], -1]])
		.CheckTdop('#(true:, true)',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(true), COLON, TRUE, COMMA], ' $
				'[CONST_MEMBER, TRUE, COMMA]], RPAREN]',
			[1, 1, 2, [3,
				[3, 3, 7, -1, 8],
				[10, 10, -1]], 14],
			type: 'constant')
		}
	}
