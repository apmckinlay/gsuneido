// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('try {}',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, LIST, RCURLY], CATCHSTMT]]')
		.CheckTdop('try;',
			'[LIST, [TRYSTMT, TRY, SEMICOLON, CATCHSTMT]]')
		.CheckTdop('try; catch;',
			'[LIST, [TRYSTMT, TRY, SEMICOLON, ' $
				'[CATCHSTMT, CATCH, CATCH_COND, SEMICOLON]]]')
		.CheckTdop('try {} catch {}',
			'[LIST, [TRYSTMT, TRY, ' $
				'[STMTS, LCURLY, LIST, RCURLY], ' $
				'[CATCHSTMT, CATCH, CATCH_COND, [STMTS, LCURLY, LIST, RCURLY]]]]')
		.CheckTdop('try a',
			'[LIST, [TRYSTMT, TRY, [STMT, IDENTIFIER(a), SEMICOLON], CATCHSTMT]]')
		.CheckTdop('try {a;b}',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], CATCHSTMT]]')
		.CheckTdop('try {a;b} catch{}',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[CATCHSTMT, CATCH, CATCH_COND, [STMTS, LCURLY, LIST, RCURLY]]]]')
		.CheckTdop('try {a;b} catch c',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[CATCHSTMT, CATCH, CATCH_COND, [STMT, IDENTIFIER(c), SEMICOLON]]]]')
		.CheckTdop('try {a;b} catch (e) d',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[CATCHSTMT, CATCH, ' $
					'[CATCH_COND, LPAREN, IDENTIFIER(e), COMMA, STRING, RPAREN], ' $
					'[STMT, IDENTIFIER(d), SEMICOLON]]]]')
		.CheckTdop('try {a;b} catch (e, "*test") d',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[CATCHSTMT, CATCH, ' $
					'[CATCH_COND, LPAREN, ' $
						'IDENTIFIER(e), COMMA, STRING(*test), RPAREN], ' $
					'[STMT, IDENTIFIER(d), SEMICOLON]]]]')
		.CheckTdop('try {a;b} catch () d',
			'[LIST, [TRYSTMT, TRY, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[CATCHSTMT, CATCH, ' $
					'[CATCH_COND, LPAREN, IDENTIFIER, COMMA, STRING, RPAREN], ' $
					'[STMT, IDENTIFIER(d), SEMICOLON]]]]')

		.CheckTdopCatch('try {} catch ("a") d', 'expected RPAREN, but got STRING(a)')
		.CheckTdopCatch('try {} catch (e "a") d', 'expected RPAREN, but got STRING(a)')
		.CheckTdopCatch('try {} catch (e, 1) d', 'expected STRING, but got NUMBER(1)')
		}
	}
