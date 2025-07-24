// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('while false;',
			'[LIST, [WHILESTMT, WHILE, LPAREN, FALSE, RPAREN, SEMICOLON]]')
		.CheckTdop('while false {}',
			'[LIST, [WHILESTMT, WHILE, LPAREN, FALSE, RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('while ob.FindIf(){it is 1} a',
			'[LIST, [WHILESTMT, WHILE, LPAREN, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(ob), DOT, IDENTIFIER(FindIf)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, [LIST, ' $
						'[STMT, [BINARYOP, IDENTIFIER(it), IS, NUMBER(1)], ' $
							'SEMICOLON]], RCURLY]], RPAREN, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON]]]')
		.CheckTdop('while test()\r\n{a}',
			'[LIST, [WHILESTMT, WHILE, LPAREN, ' $
				'[CALL, IDENTIFIER(test), LPAREN, LIST, RPAREN, BLOCK], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(a), SEMICOLON]], RCURLY]]]')
		.CheckTdop('while a > 1
			{
			a--
			b = 10
			while(b>0)
				--b
			}',
			'[LIST, [WHILESTMT, WHILE, LPAREN, ' $
				'[BINARYOP, IDENTIFIER(a), GT, NUMBER(1)], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, [POSTINCDEC, IDENTIFIER(a), DEC], SEMICOLON], ' $
					'[STMT, [BINARYOP, IDENTIFIER(b), EQ, NUMBER(10)], SEMICOLON], ' $
					'[WHILESTMT, WHILE, LPAREN, ' $
						'[BINARYOP, IDENTIFIER(b), GT, NUMBER(0)], RPAREN, ' $
						'[STMT, [PREINCDEC, DEC, IDENTIFIER(b)], SEMICOLON]]], ' $
					'RCURLY]]]',
			[1, [1, 1, -1,
				[7, 7, 9, 11], -1,
				[17, 17, [23,
					[23, [23, 23, 24], -1],
					[31, [31, 31, 33, 35], -1],
					[42, 42, 47,
						[48, 48, 49, 50], 51,
						[58, [58, 58, 60], -1]]],
					66]]])
		}
	}