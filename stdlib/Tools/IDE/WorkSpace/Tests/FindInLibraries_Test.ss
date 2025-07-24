// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_query()
		{
		q = FindInLibraries.FindInLibraries_query(#(
			libs: 'stdlib', exclude: false,
			nameRepeat: ([text: 'Test']),
			textRepeat: ([text: 'foo'], [text: 'bar', exclude: true]),
			bytoken: 'Map',
			byexpression: 'e is true',
			findUid: #30000101))
		Assert(q like: '((stdlib extend lib="stdlib", lib_current_text) where group = -1)
	where name =~ "(?i)(?q)Test(?-q)"
	where lib_current_text !~ "(?i)(?q)bar(?-q)"
	extend lines = FindCombine(' $
		'Object(#RegexMatchLines, lib_current_text, "(?i)(?q)foo(?-q)"),
		Object(#FindByToken, lib_current_text, #("Map"), #(`\<Map\>`)),
		Object(#FindByExpression, name, lib_current_text, "e is true", false), #30000101)
	where not lines.Empty?()
	sort name')
		}
	Test_libsQuery()
		{
		f = FindInLibraries.FindInLibraries_libsQuery
		data = #(libs: 'stdlib', exclude: true)
		Assert(f(data) hasnt: 'stdlib')

		data = #(libs: 'stdlib', exclude: false)
		Assert(f(data)
			like: '((stdlib extend lib="stdlib", lib_current_text) where group = -1)')
		}
	Test_nameWhere()
		{
		f = FindInLibraries.FindInLibraries_nameWhere
		data = Object(nameRepeat: Object())
		Assert(f(data) is: '')

		data.nameRepeat.Add([])
		Assert(f(data) is: '')

		data.nameRepeat.Add([text: 'foo'])
		Assert(f(data) is: ' where name =~ "(?i)(?q)foo(?-q)"\n')
		data.nameRepeat.Add([text: 'bar', exclude:])
		Assert(f(data)
			is: ' where name =~ "(?i)(?q)foo(?-q)" and name !~ "(?i)(?q)bar(?-q)"\n')
		}
	Test_textExcludes()
		{
		f = FindInLibraries.FindInLibraries_textExcludes
		data = Object(textRepeat: Object())
		Assert(f(data) is: '')

		data.textRepeat.Add([])
		Assert(f(data) is: '')

		data.textRepeat.Add([text: 'foo'])
		Assert(f(data) is: '')

		data.textRepeat.Add([text: 'bar', exclude:])
		Assert(f(data) is: ' where lib_current_text !~ "(?i)(?q)bar(?-q)"\n')

		data.textRepeat.Add([text: 'baz', exclude:])
		Assert(f(data)
			is: ' where lib_current_text !~ "(?i)(?q)bar(?-q)" and ' $
				'lib_current_text !~ "(?i)(?q)baz(?-q)"\n')
		}
	Test_linesExpr()
		{
		f = FindInLibraries.FindInLibraries_linesExpr
		data = [textRepeat: Object()]
		Assert(f(data) is: '')

		data.textRepeat.Add([])
		Assert(f(data) is: '')

		data.bytoken = '++'
		Assert(f(data) is: 'Object(#FindByToken, lib_current_text, #("++"), #())')

		data.bytoken = "is foo"
		Assert(f(data)
			is: 'Object(#FindByToken, lib_current_text, #("is", "foo"), #(`\<foo\>`))')

		data.textRepeat.Add([text: 'foo'])
		Assert(f(data)
			like: 'Object(#RegexMatchLines, lib_current_text, "(?i)(?q)foo(?-q)"),
			Object(#FindByToken, lib_current_text, #("is", "foo"), #(`\<foo\>`))')
		}
	Test_quickCheck()
		{
		f = FindInLibraries.FindInLibraries_quickCheck
		Assert(f(#('if', '123', '++', 'x')) is: #())
		Assert(f(#('if', 'foo', '++', 'bars', '123')) is: #('\<bars\>', '\<foo\>'))
		}
	Test_keyword()
		{
		k = FindInLibraries.FindInLibraries_keyword?
		Assert(k('if'))
		Assert(k('foo') is: false)
		Assert(k('') is: false)
		}
	Test_limitReach()
		{
		f = FindInLibraries.FindInLibraries_limitReached?
		hitFn = { |unused| throw 'should not run' }

		Assert(f(0, 0, hitFn) is: false)
		Assert(f(1000, 1000, hitFn) is: false)
		Assert(f(1001, 0, { |unused| }))
		Assert(f(0, 1001, { |unused| }))
		Assert(f(1001, 1001, { |unused| }))
		}
	Test_printlines()
		{
		cl = FindInLibraries
			{
			Limit: 10
			}
		fn = cl.FindInLibraries_printLines
		toPrint = Object()
		printFn = { |s| toPrint.Add(s) }

		x = [lib: 'testLib', name: 'testName',
			lib_current_text: 'A\r\nB\r\n\r\n\tC\r\n\t\tD']
		x.lines = #((1))

		Assert(fn(printFn, x, 0) is: 1)
		Assert(toPrint is: #('testLib:testName:2: B'))

		x.lines = #((1, 2, 3, 4))
		toPrint.Delete(all:)
		Assert(fn(printFn, x, 0) is: 4)
		Assert(toPrint is: #(
			'testLib:testName:2: B',
			'                    ',
			'                    \tC',
			'                    \t\tD'))

		x.lines = #((2, 3, 4))
		toPrint.Delete(all:)
		Assert(fn(printFn, x, 0) is: 3)
		Assert(toPrint is: #(
			'testLib:testName:3: ',
			'                    C',
			'                    \tD'))

		x.lines = #((2, 3, 4))
		toPrint.Delete(all:)
		Assert(fn(printFn, x, 9) is: 11)
		Assert(toPrint is: #(
			'testLib:testName:3: ',
			'                    C'))

		x.lines = #((1), (3, 4))
		toPrint.Delete(all:)
		Assert(fn(printFn, x, 0) is: 3)
		Assert(toPrint is: #(
			'testLib:testName:2: B',
			'testLib:testName:4: C',
			'                    \tD'))

		x.lines = #((2))
		toPrint.Delete(all:)
		Assert(fn(printFn, x, 0) is: 1)
		Assert(toPrint is: #('testLib:testName:3: '))
		}
	}
