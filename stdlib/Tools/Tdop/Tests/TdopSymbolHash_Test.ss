// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('#20170101.1234', 'DATE(20170101.1234)', 1, type: 'constant')
		.CheckTdop('#abc', 'STRING(abc)', 1, type: 'constant')
		.CheckTdop('#(1, 2, 3)',
			'[OBJECT, HASH, ' $
				'LPAREN, ' $
				'[LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA], ' $
					'[CONST_MEMBER, NUMBER(2), COMMA], ' $
					'[CONST_MEMBER, NUMBER(3), COMMA]], ' $
				'RPAREN]',
			[1, 1, 2, [3, [3, 3, 4], [6, 6, 7], [9, 9, -1]], 10],
			type: 'constant')
		.CheckTdop('#(1, a: 3)',
			'[OBJECT, HASH, LPAREN, [LIST, [CONST_MEMBER, NUMBER(1), COMMA], ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, NUMBER(3), COMMA]], RPAREN]',
			[1, 1, 2, [3, [3, 3, 4], [6, 6, 7, 9, -1]], 10],
			type: 'constant')
		.CheckTdop('#(#(d c: a) a: (1 2 3) b:b)',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_MEMBER, [OBJECT, HASH, LPAREN, [LIST, ' $
					'[CONST_MEMBER, STRING(d), COMMA], ' $
					'[CONST_KEYMEMBER, STRING(c), COLON, STRING(a), COMMA]], ' $
					'RPAREN], COMMA], ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, [OBJECT, HASH, LPAREN, [LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA], ' $
					'[CONST_MEMBER, NUMBER(2), COMMA], ' $
					'[CONST_MEMBER, NUMBER(3), COMMA]], RPAREN], COMMA], ' $
				'[CONST_KEYMEMBER, STRING(b), COLON, STRING(b), COMMA]], RPAREN]',
			[1, 1, 2, [3,
				[3, [3, 3, 4, [5,
					[5, 5, -1],
					[7, 7, 8, 10, -1]],
					11], -1],
				[13, 13, 14, [16, -1, 16, [17,
					[17, 17, -1],
					[19, 19, -1],
					[21, 21, -1]], 22], -1],
				[24, 24, 25, 26, -1]], 27],
			type: 'constant')
		.CheckTdop('#{1, "2", {3, a: 4}}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_MEMBER, NUMBER(1), COMMA], ' $
				'[CONST_MEMBER, STRING(2), COMMA], ' $
				'[CONST_MEMBER, [RECORD, HASH, LCURLY, [LIST, ' $
					'[CONST_MEMBER, NUMBER(3), COMMA], ' $
					'[CONST_KEYMEMBER, STRING(a), COLON, NUMBER(4), COMMA]], ' $
				'RCURLY], COMMA]], RCURLY]',
			[1, 1, 2, [3,
				[3, 3, 4],
				[6, 6, 9],
				[11, [11, -1, 11, [12,
					[12, 12, 13],
					[15, 15, 16, 18, -1]]
				19], -1]], 20],
			type: 'constant')
		.CheckTdop('#{1, "ab" $ "cd", a: "ab" $ "cd"}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_MEMBER, NUMBER(1), COMMA], ' $
				'[CONST_MEMBER, [BINARYOP, STRING(ab), CAT, STRING(cd)], COMMA], ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, ' $
					'[BINARYOP, STRING(ab), CAT, STRING(cd)], COMMA]], ' $
				'RCURLY]',
			[1, 1, 2, [3,
				[3, 3, 4],
				[6, [6, 6, 11, 13], 17],
				[19, 19, 20, [22, 22, 27, 29], -1]], 33],
			type: 'constant')
		.CheckTdop('#(function: function(){})',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(function), COLON, ' $
					'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, LCURLY, LIST, ' $
						'RCURLY], COMMA]],  ' $
				'RPAREN]',
			type: 'constant')
		.CheckTdop('#(class: class{a: 1 b(){c}} d:)',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(class), COLON, ' $
					'[CLASSDEF, CLASS, COLON, STRING, LCURLY, [LIST, ' $
						'[CONST_KEYMEMBER, STRING(a), COLON, NUMBER(1), SEMICOLON], ' $
						'[CONST_KEYMEMBER, STRING(b), COLON, ' $
							'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, LCURLY, ' $
								'[LIST, [STMT, IDENTIFIER(c), SEMICOLON]], ' $
								'RCURLY], SEMICOLON]], RCURLY], ' $
					'COMMA], ' $
				'[CONST_KEYMEMBER, STRING(d), COLON, TRUE, COMMA]], RPAREN]',
			[1, 1, 2, [3,
				[3, 3, 8,
					[10, 10, -1, -1, 15, [16,
						[16, 16, 17, 19, -1],
						[21, 21, -1,
							[22, -1, 22, -1, 23, 24,
								[25, [25, 25, -1]], 26], -1]], 27],
					-1],
				[29, 29, 30, -1, -1]], 31],
			type: 'constant')
////	NOTE: literal Ojbect/Record allow named arguments before un-named arguments
		.CheckTdop('#(a: 1, "b")',
			'[OBJECT, HASH, LPAREN, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, NUMBER(1), COMMA], ' $
				'[CONST_MEMBER, STRING(b), COMMA]], RPAREN]',
			[1, 1, 2, [3,
				[3, 3, 4, 6, 7],
				[9, 9, -1]], 12],
			type: 'constant')
		.CheckTdop('#{Object(1, 2)}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_MEMBER, STRING(Object), COMMA], ' $
				'[CONST_MEMBER, [OBJECT, HASH, LPAREN, [LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA], ' $
					'[CONST_MEMBER, NUMBER(2), COMMA]], RPAREN], COMMA]], RCURLY]',
			[1, 1, 2, [3,
				[3, 3, -1],
				[9, [9, -1, 9, [10,
					[10, 10, 11],
					[13, 13, -1]], 14], -1]], 15],
			type: 'constant')
		.CheckTdop('#{#20170101}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_MEMBER, DATE(20170101), COMMA]], RCURLY]',
			[1, 1, 2, [3,
				[3, 3, -1]], 12],
			type: 'constant')
		.CheckTdop('#{a: Object(1, 2){ 1 }}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, STRING(Object), COMMA], ' $
				'[CONST_MEMBER, [OBJECT, HASH, LPAREN, [LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA], ' $
					'[CONST_MEMBER, NUMBER(2), COMMA]], RPAREN], COMMA], ' $
				'[CONST_MEMBER, [RECORD, HASH, LCURLY, [LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA]], RCURLY], COMMA]], RCURLY]',
			[1, 1, 2, [3,
				[3, 3, 4, 6, -1],
				[12, [12, -1, 12, [13,
					[13, 13, 14],
					[16, 16, -1]], 17], -1],
				[18, [18, -1, 18, [20,
					[20, 20, -1]], 22], -1]], 23],
			type: 'constant')
		.CheckTdop('#{a: Object(1, block: {1, 2})}',
			'[RECORD, HASH, LCURLY, [LIST, ' $
				'[CONST_KEYMEMBER, STRING(a), COLON, STRING(Object), COMMA], ' $
				'[CONST_MEMBER, [OBJECT, HASH, LPAREN, [LIST, ' $
					'[CONST_MEMBER, NUMBER(1), COMMA], ' $
					'[CONST_KEYMEMBER, STRING(block), COLON, ' $
						'[RECORD, HASH, LCURLY, [LIST, ' $
							'[CONST_MEMBER, NUMBER(1), COMMA], ' $
							'[CONST_MEMBER, NUMBER(2), COMMA]], RCURLY], COMMA]], ' $
					'RPAREN], COMMA]], RCURLY]',
			[1, 1, 2, [3,
				[3, 3, 4, 6, -1],
				[12, [12, -1, 12, [13,
					[13, 13, 14],
					[16, 16, 21,
						[23, -1, 23, [24,
							[24, 24, 25],
							[27, 27, -1]], 28], -1]],
					29], -1]], 30],
			type: 'constant')

		.CheckTdopCatch('#{:a}')
		.CheckTdopCatch('#{a: c + b}')
		.CheckTdopCatch('#{a{}}', 'base class must be global defined in library')
		.CheckTdopCatch('#{a\r\n{}}', 'base class must be global defined in library')
		}
	}