// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_leaves()
		{
		Assert(AstFmtDoc.Text('a'), is: [t: #text, s: 'a', hb: false])
		Assert(AstFmtDoc.Tok('a'), is: [t: #text, s: 'a', hb: false])
		Assert(AstFmtDoc.Tok('a\nb'), is: [t: #verb, s: 'a\nb', hb:])
		Assert(AstFmtDoc.Tokc("// c"), is: [t: #text, s: "// c", hb:])
		Assert(AstFmtDoc.Tokc("/* c */"), is: [t: #text, s: "/* c */", hb: false])
		Assert(AstFmtDoc.Str('"s"'), is: [t: #str, s: '"s"', hb: false])
		}

	Test_breakConstants()
		{
		Assert(AstFmtDoc.Line.s, is: ' ')
		Assert(AstFmtDoc.Soft.s, is: "")
		Assert(AstFmtDoc.Semi.s, is: "; ")
		Assert(AstFmtDoc.Hard.hb)
		Assert(AstFmtDoc.Blank.t, is: #blank)
		Assert(AstFmtDoc.Blank.hb)
		}

	Test_cat()
		{
		d = AstFmtDoc.Cat('a', false, AstFmtDoc.Text('b'))
		Assert(d.t, is: #cat)
		Assert(d.a.Size(), is: 2) // false dropped, string coerced to text
		Assert(d.a[0], is: [t: #text, s: 'a', hb: false])
		Assert(d.hb, is: false)
		}

	Test_hbPropagation()
		{
		Assert(AstFmtDoc.Cat('a', AstFmtDoc.Hard).hb)
		Assert(AstFmtDoc.Group(AstFmtDoc.Text('x')).hb, is: false)
		Assert(AstFmtDoc.Group(AstFmtDoc.Hard).hb)
		Assert(AstFmtDoc.Nest(AstFmtDoc.Hard).hb)
		Assert(AstFmtDoc.Fill([AstFmtDoc.Text('x'), AstFmtDoc.Hard]).hb)
		Assert(AstFmtDoc.Root(AstFmtDoc.Text('x')).hb) // root never flattens
		}

	Test_combinators()
		{
		a = AstFmtDoc.Interleave(['x', 'y', 'z'], #sep)
		Assert(a, is: #(x, sep, y, sep, z))
		s = AstFmtDoc.Seq([AstFmtDoc.Text('x'), AstFmtDoc.Text('y')], AstFmtDoc.Line)
		Assert(s.t, is: #cat)
		Assert(s.a.Size(), is: 3)
		Assert(s.a[1].t, is: #line)
		f = AstFmtDoc.Fillsep([AstFmtDoc.Text('x'), AstFmtDoc.Text('y')], AstFmtDoc.Line)
		Assert(f.t, is: #fill)
		Assert(f.a.Size(), is: 3)
		}
	}
