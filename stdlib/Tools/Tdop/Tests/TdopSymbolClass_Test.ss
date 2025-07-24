// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('class{}',
			'[CLASSDEF, CLASS, COLON, STRING, LCURLY, LIST, RCURLY]',
			[1, 1, -1, -1, 6, -1, 7],
			type: 'constant')
		.CheckTdop('class: A{}',
			'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]',
			[1, 1, 6, 8, 9, -1, 10],
			type: 'constant')
		.CheckTdop('class A{}',
			'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, LIST, RCURLY]',
			[1, 1, -1, 7, 8, -1, 9],
			type: 'constant')
		.CheckTdop('class: _A{}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, LIST, RCURLY]',
			[1, 1, 6, 8, 10, -1, 11],
			type: 'constant')
		.CheckTdop('class: _A{a: (1)}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[OBJECT, HASH, LPAREN, [LIST, [CONST_MEMBER, NUMBER(1), COMMA]], ' $
						'RPAREN], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, -1, 14, [15, [15, 15, -1]],
						16], -1]], 17],
			type: 'constant')
		.CheckTdop('class: _A{a: {1}}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[RECORD, HASH, LCURLY, [LIST, [CONST_MEMBER, NUMBER(1), COMMA]], ' $
						'RCURLY], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, -1, 14, [15, [15, 15, -1]],
						16], -1]], 17],
			type: 'constant')
		.CheckTdop('class: _A{a: [1]}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[RECORD, HASH, LBRACKET, ' $
						'[LIST, [CONST_MEMBER, NUMBER(1), COMMA]], ' $
						'RBRACKET], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, -1, 14,
						[15, [15, 15, -1]],
						16], -1]], 17],
			type: 'constant')
		.CheckTdop('class: _A{a: #(1)}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[OBJECT, HASH, LPAREN, [LIST, [CONST_MEMBER, NUMBER(1), COMMA]], ' $
						'RPAREN], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, 14, 15, [16, [16, 16, -1]],
						17], -1]], 18],
			type: 'constant')
		.CheckTdop('class: _A{a: #{1}}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[RECORD, HASH, LCURLY, [LIST, [CONST_MEMBER, NUMBER(1), COMMA]], ' $
						'RCURLY], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, 14, 15, [16, [16, 16, -1]],
						17], -1]], 18],
			type: 'constant')
		.CheckTdop('class: _A{A: true B: test}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON, TRUE, SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(B), COLON, STRING(test), SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12, 14, -1],
				[19, 19, 20, 22, -1]], 26],
			type: 'constant')
		.CheckTdop('class: _A{A: function(){} B: test}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON, ' $
					'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, ' $
						'RPAREN, LCURLY, LIST, RCURLY], SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(B), COLON, STRING(test), SEMICOLON]], RCURLY]',
			type: 'constant')
		.CheckTdop('class: _A{A: class{b: 1} B: test}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON,  ' $
					'[CLASSDEF, CLASS, COLON, STRING, LCURLY, [LIST, ' $
						'[CONST_KEYMEMBER, STRING(b), COLON, NUMBER(1), SEMICOLON]], ' $
						'RCURLY], SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(B), COLON, STRING(test), SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, 14, -1, -1, 19, [20,
						[20, 20, 21, 23, -1]],
						24], -1],
				[26, 26, 27, 29, -1]], 33],
			type: 'constant')
		.CheckTdop('class: _A{A: B{b: 1} c: test}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON, ' $
					'[CLASSDEF, CLASS, COLON, STRING(B), LCURLY, [LIST, ' $
						'[CONST_KEYMEMBER, STRING(b), COLON, NUMBER(1), SEMICOLON]], ' $
						'RCURLY], SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(c), COLON, STRING(test), SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, -1, -1, 14, 15, [16,
						[16, 16, 17, 19, -1]],
						20], -1],
				[22, 22, 23, 25, -1]], 29],
			type: 'constant')
		.CheckTdop('class: _A{A: +1.1 B: -1.2 C: 1.3 -1: "a" $ "b"}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON, ' $
					'[UNARYOP, ADD, NUMBER(1.1)], SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(B), COLON, ' $
					'[UNARYOP, SUB, NUMBER(1.2)], SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(C), COLON, NUMBER(1.3), SEMICOLON], ' $
				'[CONST_KEYMEMBER, [UNARYOP, SUB, NUMBER(1)], COLON, ' $
					'[BINARYOP, STRING(a), CAT, STRING(b)], SEMICOLON]], ' $
				'RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 12,
					[14, 14, 15], -1],
				[19, 19, 20,
					[22, 22, 23], -1],
				[27, 27, 28, 30, -1],
				[34, [34, 34, 35], 36, [38, 38, 42, 44], -1]], 47]
			type: 'constant')
		.CheckTdop('class: _A{"A": "TEST" "B": "A" $ "B" $ "C"}',
			'[CLASSDEF, CLASS, COLON, STRING(_A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(A), COLON, STRING(TEST), SEMICOLON], ' $
				'[CONST_KEYMEMBER, STRING(B), COLON, ' $
					'[BINARYOP, [BINARYOP, STRING(A), CAT, STRING(B)], ' $
						'CAT, STRING(C)], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 10, [11,
				[11, 11, 14, 16, -1],
				[23, 23, 26,
					[28, [28, 28, 32 34],
						38 40], -1]], 43],
			type: 'constant')
		.CheckTdop('class: A{a(){}}',
			'[CLASSDEF, CLASS, COLON, STRING(A), LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, ' $
						'LCURLY, LIST, RCURLY], SEMICOLON]], RCURLY]',
			[1, 1, 6, 8, 9, [10,
				[10, 10, -1,
					[11, -1, 11, -1, 12,
						13, -1, 14], -1]], 15]
			type: 'constant')

		.CheckTdopCatch('class: a{}', 'base class must be global defined in library',
			type: 'constant')
		.CheckTdopCatch('class: _a{}', 'base class must be global defined in library',
			type: 'constant')
		.CheckTdopCatch('class: __A{}', 'base class must be global defined in library',
			type: 'constant')
		.CheckTdopCatch('class: A{1}', 'class members must be named',
			type: 'constant')
		.CheckTdopCatch('class: A{A: 1+2}', 'class members must be named',
			type: 'constant')
		}
	}
