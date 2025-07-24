// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('struct { int32 left; pointer* right }',
			'[STRUCTDEF, STRUCT, LCURLY, [LIST, ' $
				'[STRUCT_MEMBER, [DLL_NORMAL, IDENTIFIER(int32)], ' $
					'IDENTIFIER(left), SEMICOLON], ' $
				'[STRUCT_MEMBER, [DLL_POINTER, IDENTIFIER(pointer), MUL], ' $
					'IDENTIFIER(right), SEMICOLON]], RCURLY]',
			type: 'constant')
		.CheckTdop('struct { int32 left\r\n pointer* right\r\n }',
			'[STRUCTDEF, STRUCT, LCURLY, [LIST, ' $
				'[STRUCT_MEMBER, [DLL_NORMAL, IDENTIFIER(int32)], ' $
					'IDENTIFIER(left), SEMICOLON], ' $
				'[STRUCT_MEMBER, [DLL_POINTER, IDENTIFIER(pointer), MUL], ' $
					'IDENTIFIER(right), SEMICOLON]], RCURLY]',
			type: 'constant')

		.CheckTdopCatch('struct { int32 left pointer* right }')
		}
	}