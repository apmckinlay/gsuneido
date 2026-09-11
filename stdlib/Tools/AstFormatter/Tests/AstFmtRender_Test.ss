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
		Assert(.r(flat, width: 4), is: "aa\nbb")
		Assert(.r(d.Group(d.Cat('a', d.Soft, 'b'))), is: #ab)
		Assert(.r(d.Group(d.Cat('a', d.Semi, 'b'))), is: "a; b")
		// hard never flattens, even when it would fit
		Assert(.r(d.Group(d.Cat('a', d.Hard, 'b'))), is: "a\nb")
		}

	Test_hardIsIdempotent() // only Blank makes an empty line
		{
		d = AstFmtDoc
		Assert(.r(d.Cat('a', d.Hard, d.Hard, 'b')), is: "a\nb")
		Assert(.r(d.Cat('a', d.Hard, d.Blank, 'b')), is: "a\n\nb")
		// a // comment already ends its line; the hard after it adds nothing
		Assert(.r(d.Cat(d.Tokc("// c"), d.Hard, d.Hard, 'x')), is: "// c\nx")
		// the second hard still sets the indent
		Assert(.r(d.Cat('a', d.Hard, d.Nest(d.Cat(d.Hard, 'b')))), is: "a\n\tb")
		}

	Test_nest()
		{
		d = AstFmtDoc
		Assert(.r(d.Cat("f(", d.Nest(d.Cat(d.Hard, 'x')), d.Hard, ')')), is: "f(\n\tx\n)")
		}

	Test_trailingWhitespaceTrimmed()
		{
		d = AstFmtDoc
		Assert(.r(d.Cat("a ", d.Hard, 'b')), is: "a\nb")
		}

	Test_eolComment() // nothing may join a // comment's line
		{
		d = AstFmtDoc
		Assert(.r(d.Cat(d.Tokc("// c"), d.Text('x'))), is: "// c\nx")
		}

	Test_root() // absolute column 0, out of any nesting
		{
		d = AstFmtDoc
		Assert(.r(d.Nest(d.Cat(d.Hard, d.Root(d.Text("P()"))))), is: "\nP()")
		}

	// callee(<lead>arg) — the shape AstFmtExpr builds for a sole argument
	lead(callee, arg)
		{
		return AstFmtDoc.Cat(AstFmtDoc.Text(callee $ '('), AstFmtDoc.Lead(arg),
			AstFmtDoc.Text(')'))
		}

	// a lead break is taken when it is what makes the arg fit, or when staying
	// on the callee's line would overflow however the arg wraps
	Test_lead()
		{
		d = AstFmtDoc
		// fits whole: no break at all
		Assert(.r(.lead('f', d.Text(#ab)), width: 12), is: "f(ab)")
		// cannot sit beside the callee, but fits flat one indent deeper: break
		Assert(.r(.lead(#ffff, d.Text(#aaaaaaa)), width: 12), is: "ffff(\n\taaaaaaa)")
		// must break either way, but its first chunk still fits beside the
		// callee: stay on the callee's line, wrap in place instead
		wraps = d.Group(d.Cat(d.Text(#aaa), d.Nest(d.Cat(d.Line, d.Text(#bbbbb)))))
		Assert(.r(.lead(#ffff, wraps), width: 12), is: "ffff(aaa\n\tbbbbb)")
		// must break either way AND its first chunk does not fit beside the
		// callee: lead, since staying put overflows however the arg wraps
		mustLead = d.Group(d.Cat(d.Text(#aaaaaa), d.Nest(d.Cat(d.Line, d.Text(#bbbbbb)))))
		Assert(.r(.lead(#ffffff, mustLead), width: 12),
			is: "ffffff(\n\taaaaaa\n\t\tbbbbbb)")
		// a hard break can never fit flat, so it never earns the lead break
		Assert(.r(.lead('f', d.Cat(d.Text('a'), d.Hard, d.Text('b')))), is: "f(a\nb)")
		}

	Test_fill()
		{
		d = AstFmtDoc
		items = d.Interleave([d.Text(#aa), d.Text(#bb), d.Text(#cc)], d.Line)
		Assert(.r(d.Fill(items)), is: "aa bb cc")
		Assert(.r(d.Fill(items), width: 7), is: "aa bb\ncc")
		Assert(.r(d.Fill(items), width: 2), is: "aa\nbb\ncc")
		}

	Test_strSplit()
		{
		d = AstFmtDoc
		Assert(.r(d.Str('"aaa bbb ccc"')), is: '"aaa bbb ccc"')
		Assert(.r(d.Str('"aaa bbb ccc"'), width: 12), is: '"aaa bbb " $\n\t"ccc"')
		}

	Test_strSplitTail()
		{
		d = AstFmtDoc
		s = d.Str('"aaa bbb ccc"') // 13 wide
		Assert(.r(d.Cat(s, d.Text("))")), width: 15), is: '"aaa bbb ccc"))')
		Assert(.r(d.Cat(s, d.Text("))")), width: 13), is: '"aaa bbb " $\n\t"ccc"))')
		// a trailing comment is part of the tail
		Assert(.r(d.Cat(s, d.Text(' '), d.Tokc("// c")), width: 15),
			is: '"aaa bbb " $\n\t"ccc" // c')
		// the tail stops at the next break: the next line does not count
		Assert(.r(d.Cat(s, d.Hard, d.Text(#zzzzzzzzzz)), width: 13),
			is: '"aaa bbb ccc"\nzzzzzzzzzz')
		}

	// AstFmtRender_Test: a separator pushed to a new line by a // comment is dropped
	Test_noSeparatorAtLineStart()
		{
		d = AstFmtDoc
		Assert(.r(d.Cat(d.Text(#else), d.Text(' '), d.Tokc("// c"), d.Text(' '),
					d.Text("if x"))),
			is: "else // c\nif x")
		}
	}
