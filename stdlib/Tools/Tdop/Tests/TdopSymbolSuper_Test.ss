// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('super.super(super: super())',
			'[LIST, [STMT, [CALL, ' $
				'[MEMBEROP, SUPER, DOT, IDENTIFIER(super)], LPAREN, [LIST, ' $
					'[ARG_ELEM, [KEYARG, ' $
						'STRING(super), ' $
						'COLON, ' $
						'[CALL, [MEMBEROP, SUPER, DOT, IDENTIFIER(New)], ' $
							'LPAREN, LIST, RPAREN, BLOCK]], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1,
				[1, 1, 6, 7], 12, [13,
					[13, [13,
						13,
						18,
						[20, [20, 20, -1, -1], 25, -1, 26, -1]], -1]]
				27, -1], -1]])

		.CheckTdopCatch('super')
		}
	}
