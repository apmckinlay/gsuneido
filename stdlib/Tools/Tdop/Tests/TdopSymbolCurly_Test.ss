// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_Std()
		{
		.CheckTdop('{}', '[LIST, [STMTS, LCURLY, LIST, RCURLY]]',
			[1, [1, 1, -1, 2]])
		.CheckTdop('{}{}', '[LIST, ' $
			'[STMTS, LCURLY, LIST, RCURLY], [STMTS, LCURLY, LIST, RCURLY]]',
			[1, [1, 1, -1, 2], [3, 3, -1, 4]])
		.CheckTdop('{;;;}',
			'[LIST, [STMTS, LCURLY, [LIST, ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON], ' $
				'[STMT, NIL, SEMICOLON]], RCURLY]]',
			[1, [1, 1, [2,
				[2, -1, 2],
				[3, -1, 3],
				[4, -1, 4]], 5]])
		.CheckTdop('{};', '[LIST, ' $
			'[STMTS, LCURLY, LIST, RCURLY], [STMT, NIL, SEMICOLON]]',
			[1, [1, 1, -1, 2], [3, -1, 3]])
		.CheckTdop('{a; {b} c}',
			'[LIST, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]], RCURLY]]',
			[1, [1, 1, [2,
				[2, 2, 3],
				[5, 5, [6, [6, 6, -1]], 7],
				[9, 9, -1]], 10]])
		.CheckTdop('{a;b\r\nc}d',
			'[LIST, [STMTS, LCURLY, [LIST, ' $
				'[STMT, IDENTIFIER(a), SEMICOLON], ' $
				'[STMT, IDENTIFIER(b), SEMICOLON], ' $
				'[STMT, IDENTIFIER(c), SEMICOLON]], RCURLY], ' $
				'[STMT, IDENTIFIER(d), SEMICOLON]]',
			[1, [1, 1, [2,
				[2, 2, 3],
				[4, 4, -1],
				[7, 7, -1]], 8],
				[9, 9, -1]])
		}

	Test_Block()
		{
		.CheckTdop('fn({}){}',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, [ARG, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], COMMA]], ' $
				'RPAREN, [BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4, [4,
					[4, 4, -1, -1, -1, -1, 5]], -1]],
				6, [7, 7, -1, -1, -1, -1, 8]], -1]])
		.CheckTdop('fn({|a|})',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BLOCK, LCURLY, ' $
						'BITOR, [LIST, [BPAREM, IDENTIFIER(a), COMMA]], BITOR, ' $
						'LIST, RCURLY]], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4,
						5, [6, [6, 6, -1]], 7,
						-1, 8]], -1]],
				9, -1], -1]])
		.CheckTdop('fn({|@a|})',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BLOCK, LCURLY, ' $
						'BITOR, [BPAREM_AT, AT, IDENTIFIER(a)], BITOR, ' $
						'LIST, RCURLY]], COMMA]], ' $
				'RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4,
						5, [6, 6, 7], 8,
						-1, 9]], -1]]
				10, -1], -1]])
		.CheckTdop('fn({|a| b})',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BLOCK, LCURLY, ' $
						'BITOR, [LIST, [BPAREM, IDENTIFIER(a), COMMA]], BITOR, ' $
						'[LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]], ' $
					'COMMA]], RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4,
						5, [6, [6, 6, -1]], 7,
						[9, [9, 9, -1]], 10]],
					-1]], 11, -1], -1]])
		.CheckTdop('fn({|a, c| b})',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BLOCK, LCURLY, ' $
						'BITOR, [LIST, ' $
							'[BPAREM, IDENTIFIER(a), COMMA], ' $
							'[BPAREM, IDENTIFIER(c), COMMA]], BITOR, ' $
						'[LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]], ' $
					'COMMA]], RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4,
						5, [6,
							[6, 6, 7],
							[9, 9, -1]], 10,
						[12, [12, 12, -1]], 13]],
					-1]], 14, -1], -1]])
		.CheckTdop('fn({|a| b; {c}})',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, [LIST, ' $
				'[ARG_ELEM, ' $
					'[ARG, [BLOCK, LCURLY, ' $
						'BITOR, [LIST, [BPAREM, IDENTIFIER(a), COMMA]], BITOR, ' $
						'[LIST, ' $
							'[STMT, IDENTIFIER(b), SEMICOLON], ' $
							'[STMTS, LCURLY, [LIST, ' $
								'[STMT, IDENTIFIER(c), SEMICOLON]], RCURLY]], ' $
						'RCURLY]], ' $
					'COMMA]], RPAREN, BLOCK], SEMICOLON]]',
			[1, [1, [1, 1, 3, [4,
				[4,
					[4, [4, 4,
						5, [6, [6, 6, -1]], 7,
						[9,
							[9, 9, 10],
							[12, 12, [13,
								[13, 13, -1]], 14]],
						15]],
					-1]], 16, -1], -1]])

		.CheckTdopCatch('fn({|a, @b|})')
		.CheckTdopCatch('fn({|@a, b|})')
		.CheckTdopCatch('fn({|@+1a|})')
		}

	Test_CallBlock()
		{
		.CheckTdop('fn {}',
			'[LIST, [STMT, [CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]',
			[1, [1, [1, 1, -1, -1, -1,
				[4, 4, -1, -1, -1, -1, 5]], -1]])
		.CheckTdop('fn \r\n{}',
			'[LIST, [STMT, IDENTIFIER(fn), SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]')
		.CheckTdop('a * fn {}',
			'[LIST, [STMT, [BINARYOP, ' $
				'IDENTIFIER(a), MUL, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]]], SEMICOLON]]')
		.CheckTdop('a.fn {}',
			'[LIST, [STMT, [CALL, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(fn)], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('a.Fn {}',
			'[LIST, [STMT, [CALL, [MEMBEROP, IDENTIFIER(a), DOT, IDENTIFIER(Fn)], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')

		.CheckTdop('"a"{}',
			'[LIST, [STMT, [CALL, STRING(a), LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('fn(){}{}',
			'[LIST, [STMT, [CALL, ' $
				'[CALL, IDENTIFIER(fn), LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')

		.CheckTdop('A.B{}{}',
			'[LIST, [STMT, [CALL, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(B)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('A.B\r\n{}{}',
			'[LIST, [STMT, [CALL, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(B)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('A.B{}\r\n{}',
			'[LIST, [STMT, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(B)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]')
		.CheckTdop('A.B\r\n{}\r\n{}',
			'[LIST, [STMT, ' $
				'[CALL, [MEMBEROP, IDENTIFIER(A), DOT, IDENTIFIER(B)], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]')

		.CheckTdop('A[B]{}{}',
			'[LIST, [STMT, [CALL, ' $
				'[CALL, [SUBSCRIPT, IDENTIFIER(A), LBRACKET, IDENTIFIER(B), RBRACKET], ' $
					'LPAREN, LIST, RPAREN, ' $
					'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], ' $
				'LPAREN, LIST, RPAREN, ' $
				'[BLOCK, LCURLY, BITOR, LIST, BITOR, LIST, RCURLY]], SEMICOLON]]')
		.CheckTdop('A[B]\r\n{}{}',
			'[LIST, ' $
				'[STMT, [SUBSCRIPT, IDENTIFIER(A), LBRACKET, IDENTIFIER(B), RBRACKET], ' $
					'SEMICOLON], ' $
				'[STMTS, LCURLY, LIST, RCURLY], ' $
				'[STMTS, LCURLY, LIST, RCURLY]]')
		}

	Test_ConstantRecord()
		{
		.CheckTdop('{abc, fn: function(){}, cl: Test{}}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_MEMBER, STRING(abc), COMMA], ' $
				'[CONST_KEYMEMBER, STRING(fn), COLON, ' $
					'[FUNCTIONDEF, FUNCTION, ' $
						'LPAREN, LIST, RPAREN, LCURLY, LIST, RCURLY], COMMA], ' $
				'[CONST_KEYMEMBER, STRING(cl), COLON, ' $
					'[CLASSDEF, CLASS, COLON, STRING(Test), LCURLY, LIST, RCURLY], ' $
					'COMMA]], ' $
				'RCURLY]'
			type: 'constant')
		}

	Test_Class()
		{
		.CheckTdop('A{}',
			'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]',
			type: 'constant')
		.CheckTdop('A\r\n{}',
			'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]',
			type: 'constant')
		.CheckTdop('if a is A{} {b}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[BINARYOP, IDENTIFIER(a), IS, ' $
					'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]], ' $
				'RPAREN, ' $
				'[STMTS, LCURLY, [LIST, ' $
					'[STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'ELSE, STMTS]]')
		.CheckTdop('if a is A\r\n{} {b}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[BINARYOP, IDENTIFIER(a), IS, IDENTIFIER(A)], RPAREN, ' $
				'[STMTS, LCURLY, LIST, RCURLY], ELSE, STMTS], ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY]]')
		.CheckTdop('if (a is A\r\n{}) {b}',
			'[LIST, [IFSTMT, IF, LPAREN, ' $
				'[BINARYOP, IDENTIFIER(a), IS, ' $
					'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]], ' $
				'RPAREN, ' $
				'[STMTS, LCURLY, [LIST, [STMT, IDENTIFIER(b), SEMICOLON]], RCURLY], ' $
				'ELSE, STMTS]]')
		}
	}