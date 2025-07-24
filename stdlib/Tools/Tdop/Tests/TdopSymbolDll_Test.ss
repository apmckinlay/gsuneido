// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('dll a b:c@2(d e, [in] f[2] g, h* i)'
			'[DLLDEF, DLL, IDENTIFIER(a), IDENTIFIER(b), COLON, STRING(c@2), ' $
				'LPAREN, [LIST, ' $
					'[DLL_PAREM, DLL_IN, ' $
						'[DLL_NORMAL, IDENTIFIER(d)], ' $
						'IDENTIFIER(e), COMMA], ' $
					'[DLL_PAREM, [DLL_IN, LBRACKET, IN, RBRACKET], ' $
						'[DLL_ARRAY, IDENTIFIER(f), LBRACKET, NUMBER(2), RBRACKET], ' $
						'IDENTIFIER(g), COMMA], ' $
					'[DLL_PAREM, DLL_IN, ' $
						'[DLL_POINTER, IDENTIFIER(h), MUL], ' $
						'IDENTIFIER(i), COMMA]], RPAREN]',
			[1, 1, 5, 7, 8, 9,
				12, [13,
					[13, -1,
						[13, 13],
						15, 16],
					[18, [18, 18, 19, 21],
						[23, 23, 24, 25, 26],
						28, 29],
					[31, -1,
						[31, 31, 32],
						34, -1]], 35]
			type: 'constant')
		.CheckTdopCatch('dll a b:c(d e f)')
		.CheckTdopCatch('dll a b:c(d e ,, f g)')
		}
	}