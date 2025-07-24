// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
TdopTests
	{
	Test_main()
		{
		.CheckTdop('function(){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, LCURLY, LIST, RCURLY]',
			[1, 1, 9, -1, 10, 11, -1, 12]
			type: 'constant')
		.CheckTdop('function(@a){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, ' $
				'[PAREM_AT, AT, IDENTIFIER(a)], RPAREN, LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10, 10, 11], 12, 13, -1, 14]
			type: 'constant')
		.CheckTdop('function(a, b=1){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM, DOT, IDENTIFIER(a), COMMA], ' $
				'[PAREM_DEFAULT, DOT, IDENTIFIER(b), EQ, NUMBER(1), COMMA]], RPAREN, ' $
				'LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10,
				[10, -1, 10, 11],
				[13, -1, 13, 14, 15, -1]], 16,
				17, -1, 18]
			type: 'constant')
		.CheckTdop('function(a b=#(1)){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM, DOT, IDENTIFIER(a), COMMA], ' $
				'[PAREM_DEFAULT, DOT, IDENTIFIER(b), EQ, ' $
					'[OBJECT, HASH, LPAREN, [LIST, ' $
						'[CONST_MEMBER, NUMBER(1), COMMA]], RPAREN], COMMA]], RPAREN, ' $
				'LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10,
				[10, -1, 10, -1],
				[12, -1, 12, 13,
					[14, 14, 15, [16,
						[16, 16, -1]], 17], -1]], 18,
				19, -1, 20]
			type: 'constant')
		.CheckTdop('function(a = function(){}){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM_DEFAULT, DOT, IDENTIFIER(a), EQ, ' $
					'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, ' $
						'LCURLY, LIST, RCURLY], ' $
					'COMMA]], RPAREN, ' $
				'LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10,
				[10, -1, 10, 12,
					[14, 14, 22, -1, 23,
						24, -1, 25],
					-1]], 26,
				27, -1, 28]
			type: 'constant')
		.CheckTdop('function(.a, .b = #()){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM, DOT, IDENTIFIER(a), COMMA], ' $
				'[PAREM_DEFAULT, DOT, IDENTIFIER(b), EQ, ' $
					'[OBJECT, HASH, LPAREN, LIST, RPAREN], COMMA]], RPAREN, ' $
				'LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10,
				[10, 10, 11, 12],
				[14, 14, 15, 17,
					[19, 19, 20, -1, 21], -1]], 22,
				23, -1, 24]
			type: 'constant')
		.CheckTdop('function(a b){c;;d}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM, DOT, IDENTIFIER(a), COMMA], ' $
				'[PAREM, DOT, IDENTIFIER(b), COMMA]], RPAREN, ' $
				'LCURLY, [LIST, ' $
					'[STMT, IDENTIFIER(c), SEMICOLON], ' $
					'[STMT, NIL, SEMICOLON], ' $
					'[STMT, IDENTIFIER(d), SEMICOLON]], RCURLY]',
			type: 'constant')
		.CheckTdop('fn = function(){}',
			'[LIST, [STMT, [BINARYOP, IDENTIFIER(fn), EQ, ' $
				'[FUNCTIONDEF, FUNCTION, LPAREN, LIST, RPAREN, LCURLY, LIST, RCURLY]], ' $
				'SEMICOLON]]')
		.CheckTdop('function(a=#(function: 1)){}',
			'[FUNCTIONDEF, FUNCTION, LPAREN, [LIST, ' $
				'[PAREM_DEFAULT, DOT, IDENTIFIER(a), EQ, ' $
					'[OBJECT, HASH, LPAREN, [LIST, ' $
						'[CONST_KEYMEMBER, STRING(function), ' $
							'COLON, NUMBER(1), COMMA]], RPAREN], COMMA]], RPAREN, ' $
				'LCURLY, LIST, RCURLY]',
			[1, 1, 9, [10,
				[10, -1, 10, 11,
					[12, 12, 13, [14,
						[14, 14,
							22, 24, -1]], 25], -1]], 26,
				27, -1, 28]
			type: 'constant')

		.CheckTdopCatch('function(a, b=1, c){}', 'Default parameters must come last')
		.CheckTdopCatch('function(a, b=c){}', 'parameter defaults must be constants')
		.CheckTdopCatch('function(a, b=1+2){}')
		}
	}