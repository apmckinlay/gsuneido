// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('a', '[LIST, [STMT, IDENTIFIER(a), SEMICOLON]]')
		.CheckTdop('STRING', '[LIST, [STMT, IDENTIFIER(STRING), SEMICOLON]]')
		.CheckTdop('NUMBER', '[LIST, [STMT, IDENTIFIER(NUMBER), SEMICOLON]]')
		.CheckTdop('a\r\n {b}',
			'[LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]]')
		.CheckTdop('A\r\n {b: 1}',
			'[LIST, [STMT, [CLASSDEF, CLASS, COLON, STRING(A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(b), COLON, NUMBER(1), SEMICOLON]], RCURLY], ' $
				'SEMICOLON]]')
		.CheckTdop('A {b: 1}',
			'[LIST, [STMT, [CLASSDEF, CLASS, COLON, STRING(A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(b), COLON, NUMBER(1), SEMICOLON]], RCURLY], ' $
				'SEMICOLON]]')
		}
	}