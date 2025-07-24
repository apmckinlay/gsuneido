// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('callback(long wnd, long msg)',
			'[CALLBACKDEF, CALLBACK, LPAREN, [LIST, ' $
				'[CALLBACK_PAREM, [DLL_NORMAL, IDENTIFIER(long)], ' $
					'IDENTIFIER(wnd), COMMA], ' $
				'[CALLBACK_PAREM, [DLL_NORMAL, IDENTIFIER(long)], ' $
					'IDENTIFIER(msg), COMMA]], RPAREN]',
			type: 'constant')
		}
	}