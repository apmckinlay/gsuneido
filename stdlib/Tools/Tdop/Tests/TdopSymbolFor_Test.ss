// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_For()
		{
		.CheckTdop('for (;;){}',
			'[LIST, [FORSTMT, FOR, ' $
				'LPAREN, LIST, SEMICOLON, TRUE, SEMICOLON, LIST, RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY]]]',
			[1, [1, 1,
				5, -1, 6, -1, 7, -1, 8,
				[9, 9, -1, 10]]])
		.CheckTdop('for (;;);{}',
			'[LIST, [FORSTMT, FOR, ' $
				'LPAREN, LIST, SEMICOLON, TRUE, SEMICOLON, LIST, RPAREN, SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]',
			[1, [1, 1,
				5, -1, 6, -1, 7, -1, 8, 9],
				[10, 10, -1, 11]])
		.CheckTdop('for (;;) a',
			'[LIST, [FORSTMT, FOR, ' $
				'LPAREN, LIST, SEMICOLON, TRUE, SEMICOLON, LIST, RPAREN, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON]]]')
		.CheckTdop('for(a=0, b=c>1\r\n?1\r\n:2;a<10;a++,b++)\r\n{Print(a+b)}',
			'[LIST, [FORSTMT, FOR, LPAREN, ' $
				'[LIST, ' $
					'[EXPR_ELEM, [BINARYOP, IDENTIFIER(a), EQ, NUMBER(0)], COMMA], ' $
					'[EXPR_ELEM, [BINARYOP, IDENTIFIER(b), EQ, [TRINARYOP, ' $
						'[BINARYOP, IDENTIFIER(c), GT, NUMBER(1)], ' $
						'Q_MARK, ' $
						'NUMBER(1), ' $
						'COLON, ' $
						'NUMBER(2)]], COMMA]], SEMICOLON, ' $
				'[BINARYOP, IDENTIFIER(a), LT, NUMBER(10)], SEMICOLON, ' $
				'[LIST, ' $
					'[EXPR_ELEM, [POSTINCDEC, IDENTIFIER(a), INC], COMMA], ' $
					'[EXPR_ELEM, [POSTINCDEC, IDENTIFIER(b), INC], COMMA]], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, [CALL, IDENTIFIER(Print), LPAREN, [LIST, ' $
						'[ARG_ELEM, ' $
							'[ARG, [BINARYOP, IDENTIFIER(a), ADD, IDENTIFIER(b)]], ' $
							'COMMA]], ' $
						'RPAREN, BLOCK], SEMICOLON]], RCURLY]]]',
			[1, [1, 1, 4,
				[5,
					[5, [5, 5, 6, 7], 8],
					[10, [10, 10, 11, [12,
						[12, 12, 13, 14],
						17,
						18,
						21,
						22]], -1]], 23,
				[24, 24, 25, 26], 28,
				[29,
					[29, [29, 29, 30], 32],
					[33, [33, 33, 34], -1]], 36,
				[39, 39, [40,
					[40, [40, 40, 45, [46,
						[46,
							[46, [46, 46, 47, 48]],
							-1]],
						49, -1], -1]], 50]]])

		.CheckTdopCatch('for a=0; a<5; a++ \r\n{}', 'expected RANGETO')
		}

	Test_ForIn()
		{
		.CheckTdop('for i in ob {}',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, IDENTIFIER(ob), RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('for i in ob; {}',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, IDENTIFIER(ob), RPAREN, SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]')
		.CheckTdop('for i in Ob\r\n {}',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, IDENTIFIER(Ob), RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('for i in Ob {} {}',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, ' $
				'[CLASSDEF, CLASS, COLON, STRING(Ob), LCURLY, LIST, RCURLY], RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY]]]')
		.CheckTdop('for i in ob i',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, IDENTIFIER(ob), RPAREN, ' $
				'[STMT, IDENTIFIER(i), SEMICOLON]]]')
		.CheckTdop('for i in a..b i',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, ' $
					'[RANGE, IDENTIFIER(a), RANGETO, IDENTIFIER(b)], RPAREN, ' $
				'[STMT, IDENTIFIER(i), SEMICOLON]]]',
			[1, [1, 1,
				-1, 5, 7,
					[10, 10, 11, 13], -1, [15, 15, -1]]])
		.CheckTdop('for i in ..b i',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(i), IN, ' $
					'[RANGE, NUMBER(0), RANGETO, IDENTIFIER(b)], RPAREN, ' $
				'[STMT, IDENTIFIER(i), SEMICOLON]]]',
			[1, [1, 1,
				-1, 5, 7,
					[10, -1, 10, 12], -1, [14, 14, -1]]])
		.CheckTdop('for ..b i',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER, IN, ' $
					'[RANGE, NUMBER(0), RANGETO, IDENTIFIER(b)], RPAREN, ' $
				'[STMT, IDENTIFIER(i), SEMICOLON]]]',
			[1, [1, 1,
				-1, -1, -1,
					[5, -1, 5, 7], -1,
				[9, 9, -1]]])
		.CheckTdop('for for in in.Split(",")
			{
			a = for
			Print(a)
			}',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(for), IN, ' $
				'[CALL, ' $
					'[MEMBEROP, IDENTIFIER(in), DOT, IDENTIFIER(Split)], ' $
					'LPAREN, [LIST, [ARG_ELEM, [ARG, STRING(,)], COMMA]], ' $
					'RPAREN, BLOCK], RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, ' $
						'[BINARYOP, IDENTIFIER(a), EQ, IDENTIFIER(for)], SEMICOLON], ' $
					'[STMT, [CALL, IDENTIFIER(Print), ' $
						'LPAREN, [LIST, [ARG_ELEM, [ARG, IDENTIFIER(a)], COMMA]], ' $
						'RPAREN, BLOCK], SEMICOLON]], RCURLY]]]',
			[1, [1, 1,
				-1, 5, 9,
				[12,
					[12, 12, 14, 15],
					20, [21, [21, [21, 21], -1]],
					24, -1], -1,
				[30, 30, [36,
					[36,
						[36, 36, 38, 40], -1],
					[48, [48, 48,
						53, [54, [54, [54, 54], -1]],
						55, -1], -1]], 61]]])
		.CheckTdop('for a in ob.Map(){it} b',
			'[LIST, [FORINSTMT, FOR, ' $
				'LPAREN, IDENTIFIER(a), IN, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(ob), DOT, IDENTIFIER(Map)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, ' $
						'[LIST, [STMT, IDENTIFIER(it), SEMICOLON]], RCURLY]], RPAREN, ' $
				'[STMT, IDENTIFIER(b), SEMICOLON]]]')

		.CheckTdopCatch('for 1 in ob {}')
		.CheckTdopCatch('for a,b in ob {}')
		.CheckTdopCatch('for i in a.. {}')
		.CheckTdopCatch('for a {}')
		.CheckTdopCatch('for .. {}')
		.CheckTdopCatch('for (..b) {}')
		}
	}