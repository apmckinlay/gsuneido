// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('throw true',
			'[LIST, [THROWSTMT, THROW, TRUE, SEMICOLON]]')
		.CheckTdop('throw true;',
			'[LIST, [THROWSTMT, THROW, TRUE, SEMICOLON]]')
		.CheckTdop('throw \r\n"error"',
			'[LIST, [THROWSTMT, THROW, STRING(error), SEMICOLON]]')

		.CheckTdopCatch('throw;')
		}
	}
