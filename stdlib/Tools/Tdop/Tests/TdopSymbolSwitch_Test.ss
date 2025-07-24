// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('switch {}',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, TRUE, RPAREN, LCURLY, LIST, RCURLY]]')
		.CheckTdop('switch fn()\r\n{default:}',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK], RPAREN, LCURLY, ' $
				'[LIST, [CASE_ELEM, DEFAULT, LIST, COLON, LIST]], RCURLY]]')
		.CheckTdop('switch (a) { case 1: \r\n case 2: \r\n default: }',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'LCURLY, [LIST, ' $
					'[CASE_ELEM, CASE, [LIST, ' $
						'[EXPR_ELEM, NUMBER(1), COMMA]], COLON, LIST], ' $
					'[CASE_ELEM, CASE, [LIST, ' $
						'[EXPR_ELEM, NUMBER(2), COMMA]], COLON, LIST], ' $
					'[CASE_ELEM, DEFAULT, LIST, COLON, LIST]], RCURLY]]')
		.CheckTdop('switch (a) { case 1: \r\na;b\r\nc\r\n case 2: \r\n default: }',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'LCURLY, [LIST, ' $
					'[CASE_ELEM, CASE, [LIST, ' $
						'[EXPR_ELEM, NUMBER(1), COMMA]], COLON, [LIST, ' $
							'[STMT, IDENTIFIER(a), SEMICOLON], ' $
							'[STMT, IDENTIFIER(b), SEMICOLON], ' $
							'[STMT, IDENTIFIER(c), SEMICOLON]]], ' $
					'[CASE_ELEM, CASE, [LIST, ' $
						'[EXPR_ELEM, NUMBER(2), COMMA]], COLON, LIST], ' $
					'[CASE_ELEM, DEFAULT, LIST, COLON, LIST]], RCURLY]]')
		.CheckTdop('switch a.b {case 1, i>1, b[0]: c \r\ncase 2: d \r\ndefault: e}',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, ' $
				'[MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(b)], RPAREN, ' $
				'LCURLY, [LIST, ' $
					'[CASE_ELEM, CASE, ' $
						'[LIST, ' $
							'[EXPR_ELEM, NUMBER(1), COMMA], ' $
							'[EXPR_ELEM, ' $
								'[BINARYOP, IDENTIFIER(i), GT, NUMBER(1)], COMMA], ' $
							'[EXPR_ELEM, [SUBSCRIPT, IDENTIFIER(b), ' $
								'LBRACKET, NUMBER(0), RBRACKET], COMMA]], ' $
						'COLON, [LIST, [STMT, IDENTIFIER(c), SEMICOLON]]], ' $
					'[CASE_ELEM, CASE, ' $
						'[LIST, ' $
							'[EXPR_ELEM, NUMBER(2), COMMA]], ' $
						'COLON, [LIST, [STMT, IDENTIFIER(d), SEMICOLON]]], ' $
					'[CASE_ELEM, DEFAULT, LIST, COLON, ' $
						'[LIST, [STMT, IDENTIFIER(e), SEMICOLON]]]], RCURLY]]',
			[1, [1, 1, -1,
				[8, 8, 9, 10], -1,
				12, [13,
					[13, 13,
						[18,
							[18, 18, 19],
							[21,
								[21, 21, 22, 23], 24],
							[26, [26, 26,
								27, 28, 29], -1]],
						30, [32, [32, 32, -1]]],
					[36, 36,
						[41,
							[41, 41, -1]],
						42, [44, [44, 44, -1]]],
					[48, 48, -1, 55,
						[57, [57, 57, -1]]]], 58]])
		.CheckTdop('switch a {case 1: if b>0 b else c \r\ndefault:}',
			'[LIST, [SWITCHSTMT, SWITCH, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'LCURLY, [LIST, ' $
					'[CASE_ELEM, CASE, ' $
						'[LIST, ' $
							'[EXPR_ELEM, NUMBER(1), COMMA]], ' $
						'COLON, [LIST, ' $
							'[IFSTMT, IF, LPAREN, ' $
								'[BINARYOP, IDENTIFIER(b), GT, NUMBER(0)], RPAREN, ' $
								'[STMT, IDENTIFIER(b), SEMICOLON], ' $
								'ELSE, ' $
								'[STMT, IDENTIFIER(c), SEMICOLON]]]], ' $
					'[CASE_ELEM, DEFAULT, LIST, COLON, LIST]], RCURLY]]',
			[1, [1, 1, -1, 8, -1,
				10, [11,
					[11, 11,
						[16,
							[16, 16, -1]],
						17, [19,
							[19, 19, -1,
								[22, 22, 23, 24], -1,
								[26, 26, -1],
								28,
								[33, 33, -1]]]],
					[37, 37, -1, 44, -1]], 45]])

		.CheckTdopCatch('switch a { default: 1 \r\n case 1: 2}',
			'Invalid switch: un-reachable case after default')
		.CheckTdopCatch('switch a { if a b}', 'Unexpected: IDENTIFIER(if)')
		}
	}