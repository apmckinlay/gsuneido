// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('return',
			'[LIST, [RETURNSTMT, RETURN, THROW, LIST, SEMICOLON]]')
		.CheckTdop('return a',
			'[LIST, [RETURNSTMT, RETURN, THROW, ' $
				'[LIST, [EXPR_ELEM, IDENTIFIER(a), COMMA]], SEMICOLON]]')
		.CheckTdop('return \r\na',
			'[LIST, [RETURNSTMT, RETURN, THROW, LIST, SEMICOLON], ' $
				'[STMT, IDENTIFIER(a), SEMICOLON]]')
		.CheckTdop('return; a',
			'[LIST, [RETURNSTMT, RETURN, THROW, LIST, SEMICOLON], ' $
				'[STMT, IDENTIFIER(a), SEMICOLON]]',
			[1, [1, 1, -1, -1, 7],
				[9, 9, -1]])
		.CheckTdop('return {}',
			'[LIST, [RETURNSTMT, RETURN, THROW, ' $
				'[LIST, [EXPR_ELEM, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY], COMMA]], ' $
				'SEMICOLON]]')
		.CheckTdop('function(){return}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, ' $
				'LCURLY, [LIST, [RETURNSTMT, RETURN, THROW, LIST, SEMICOLON]], RCURLY]',
			type: 'constant')
		.CheckTdop('return fn()\r\n{a}',
			'[LIST, [RETURNSTMT, RETURN, THROW, ' $
				'[LIST, [EXPR_ELEM, [CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, [LIST, ' $
						'[STMT, IDENTIFIER(a), SEMICOLON]], RCURLY]], COMMA]], ' $
				'SEMICOLON]]',
			[1, [1, 1, -1,
				[8, [8, [8, 8, 10, -1, 11,
					[14, 14, -1, -1, -1, [15,
						[15, 15, -1]], 16]], -1]],
				-1]])
		.CheckTdop('return throw a',
			'[LIST, [RETURNSTMT, RETURN, THROW, ' $
				'[LIST, [EXPR_ELEM, IDENTIFIER(a), COMMA]], ' $
				'SEMICOLON]]',
			[1, [1, 1, 8,
				[14, [14, 14, -1]],
				-1]])
		.CheckTdop('return a, b+1, c()',
			'[LIST, [RETURNSTMT, RETURN, THROW, ' $
				'[LIST, ' $
					'[EXPR_ELEM, IDENTIFIER(a), COMMA], ' $
					'[EXPR_ELEM, [BINARYOP, IDENTIFIER(b), ADD, NUMBER(1)], COMMA], ' $
					'[EXPR_ELEM, [CALL, IDENTIFIER(c), LPAREN, LIST, RPAREN, BLOCK], ' $
						'COMMA]], ' $
				'SEMICOLON]]',
			[1, [1, 1, -1,
				[8,
					[8, 8, 9],
					[11, [11, 11, 12, 13], 14],
					[16, [16, 16, 17, -1, 18, -1],
						-1]],
				-1]])
		}
	}