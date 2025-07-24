// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('#(if:)',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(if), COLON, TRUE, COMMA]], RPAREN]',
			[1, 1, 2, [3,
				[3, 3, 5, -1, -1]], 6]
			type: 'constant')
		.CheckTdop('if 1<2;',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[BINARYOP, NUMBER(1), LT, NUMBER(2)], RPAREN, ' $
				'SEMICOLON, ELSE, STMTS]]',
			[1, [1, 1, -1,
				[4, 4, 5, 6], -1,
				7, -1, -1]])
		.CheckTdop('if a\r\n;else b',
			'[LIST, [IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'SEMICOLON, ELSE, [STMT, IDENTIFIER(b), SEMICOLON]]]',
			[1, [1, 1, -1, 4, -1,
				7, 8, [13, 13, -1]]])
		.CheckTdop('if a\r\n;else ;b',
			'[LIST, [IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'SEMICOLON, ELSE, SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]',
			[1, [1, 1, -1, 4, -1,
				7, 8, 13],
				[14, 14, -1]])
		.CheckTdop('if a>1 \r\n ++a',
			'[LIST, [IFSTMT, IF, ' $
				'LPAREN, [BINARYOP, IDENTIFIER(a), GT, NUMBER(1)], RPAREN, ' $
				'[STMT, [PREINCDEC, INC, IDENTIFIER(a)], SEMICOLON], ' $
				'ELSE, STMTS]]')
		.CheckTdop('if a>1 \r\n ++a else --a',
			'[LIST, [IFSTMT, IF, ' $
				'LPAREN, [BINARYOP, IDENTIFIER(a), GT, NUMBER(1)], RPAREN, ' $
				'[STMT, [PREINCDEC, INC, IDENTIFIER(a)], SEMICOLON], ' $
				'ELSE, [STMT, [PREINCDEC, DEC, IDENTIFIER(a)], SEMICOLON]]]')
		.CheckTdop('if a {b} c',
			'[LIST, [IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'ELSE, STMTS], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]]')
		.CheckTdop('if a { } else { }',
			'[LIST, [IFSTMT, IF, LPAREN, IDENTIFIER(a), RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY], ' $
				'ELSE, [STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('if fn()\r\n{a}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, BLOCK], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(a), SEMICOLON]], RCURLY], ' $
				'ELSE, STMTS]]')
		.CheckTdop('if a>1
			{
			++a
			Print(a)
			}
			else if (a > -1)
				a++
			else
				a--',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[BINARYOP, IDENTIFIER(a), GT, NUMBER(1)], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, [PREINCDEC, INC, IDENTIFIER(a)], SEMICOLON], ' $
					'[STMT, [CALL, IDENTIFIER(Print), LPAREN, [LIST, ' $
						'[ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA]], RPAREN, BLOCK], ' $
						'SEMICOLON]], ' $
					'RCURLY], ' $
				'ELSE, [IFSTMT, IF, LPAREN, ' $
					'[BINARYOP, IDENTIFIER(a), GT, [UNARYOP, SUB, NUMBER(1)]], RPAREN, ' $
					'[STMT, [POSTINCDEC, IDENTIFIER(a), INC], SEMICOLON], ' $
					'ELSE, [STMT, [POSTINCDEC, IDENTIFIER(a), DEC], SEMICOLON]]]]',
			[1, [1, 1, -1,
				[4, 4, 5, 6], -1,
				[12, 12, [18,
					[18, [18, 18, 20], -1],
					[26, [26, 26, 31, [32,
						[32, [32, 32], -1]], 33, -1],
						-1]],
					39],
				45, [50, 50, 53,
					[54, 54, 56, [58, 58, 59]], 60,
					[67, [67, 67, 68], -1],
					75, [85, [85, 85, 86], -1]]]])
		}

	}