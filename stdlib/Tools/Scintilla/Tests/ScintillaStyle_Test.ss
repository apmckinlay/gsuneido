// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// marks returns one character per source character: '^' where the
	// character is coloured as a type annotation, '.' everywhere else
	marks(src)
		{
		styles = ScintillaStyle.StyleString(src)
		result = ''
		for i in ..styles.Size()
			result $= styles[i] is '\x07' /*= ANNOTATION */ ? '^' : '.'
		return result
		}

	check(src, expected)
		{
		Assert(.marks(src) is: expected, msg: src)
		}

	styleOf(src, sub)
		{
		return ScintillaStyle.StyleString(src)[src.Find(sub) :: sub.Size()]
		}

	Test_definitions()
		{
		.check('Add(a: number) : boolean {}',
			'.....^.^^^^^^..^.^^^^^^^...')
		.check('foo(x: boolean|string) :number|object {}',
			'.....^.^^^^^^^^^^^^^^..^^^^^^^^^^^^^^...')
		.check('Origin() : Point {}',
			'.........^.^^^^^...')
		.check('function (x: number) {}',
			'...........^.^^^^^^....')
		.check('Configure(.disabled: object) {}',
			'...................^.^^^^^^....')
		}

	Test_noAnnotations()
		{
		.check('Bar() {}',
			'........')
		.check('Sum(@args) {}',
			'.............')
		.check('if (a ? b : c) {}',
			'.................')
		}

	Test_callsNotColoured()
		{
		.check("Query1('tables', table: 'test')",
			'...............................')
		.check('Query1(table: tablename)',
			'........................')
		.check('x = Foo(a: 1)',
			'.............')
		}

	Test_memberCallsNotColoured()
		{
		.check('.Configure(disabled: false) {}',
			'..............................')
		.check('.Configure(disabled: tablename) {}',
			'..................................')
		}

	Test_valueKeywordsNotColoured()
		{
		.check('Foo(disabled: false) {}',
			'.......................')
		.check('Configure(.disabled: false) {}',
			'..............................')
		}

	Test_keepsUnderlyingColour()
		{
		Assert(.styleOf('Add(a: number) {}', 'number') is: '\x07'.Repeat(6),
			msg: 'type name should be ANNOTATION')
		Assert(.styleOf('Foo(disabled: false) {}', 'false') is: '\x04'.Repeat(5),
			msg: 'false should stay KEYWORD')
		}
	}
