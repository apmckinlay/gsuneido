// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	r(doc, width = 90)
		{
		return AstFmtRender(width).Render(doc)
		}

	Test_breaks()
		{
		d = AstFmtDoc
		flat = d.Group(d.Cat(#aa, d.Line, #bb))
		Assert(.r(flat), is: "aa bb")
		Assert(.r(flat, width: 4), is: 'aa\nbb')
		Assert(.r(d.Group(d.Cat('a', d.Soft, 'b'))), is: #ab)
		Assert(.r(d.Group(d.Cat('a', d.Semi, 'b'))), is: "a; b")
		// hard never flattens, even when it would fit
		Assert(.r(d.Group(d.Cat('a', d.Hard, 'b'))), is: 'a\nb')
		}

	Test_nest()
		{
		d = AstFmtDoc
		Assert(.r(d.Cat("f(", d.Nest(d.Cat(d.Hard, 'x')), d.Hard, ')')), is: 'f(\n\tx\n)')
		}

	Test_trailingWhitespaceTrimmed()
		{
		d = AstFmtDoc
		Assert(.r(d.Cat("a ", d.Hard, 'b')), is: 'a\nb')
		}

	Test_eolComment() // nothing may join a // comment's line
		{
		d = AstFmtDoc
		Assert(.r(d.Cat(d.Tokc("// c"), d.Text('x'))), is: '// c\nx')
		}

	Test_root() // absolute column 0, out of any nesting
		{
		d = AstFmtDoc
		Assert(.r(d.Nest(d.Cat(d.Hard, d.Root(d.Text("P()"))))), is: '\nP()')
		}

	Test_fill()
		{
		d = AstFmtDoc
		items = d.Interleave([d.Text(#aa), d.Text(#bb), d.Text(#cc)], d.Line)
		Assert(.r(d.Fill(items)), is: "aa bb cc")
		Assert(.r(d.Fill(items), width: 7), is: 'aa bb\ncc')
		Assert(.r(d.Fill(items), width: 2), is: 'aa\nbb\ncc')
		}

	Test_strSplit()
		{
		d = AstFmtDoc
		Assert(.r(d.Str('"aaa bbb ccc"')), is: '"aaa bbb ccc"')
		Assert(.r(d.Str('"aaa bbb ccc"'), width: 12), is: '"aaa bbb " $\n\t"ccc"')
		}
	}
