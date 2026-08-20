// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	style(src = "")
		{
		return AstFmtStyle(src, AstFmtComments(src))
		}

	fn(body)
		{
		return "function ()\n\t{\n\t" $ body $ "\n\t}"
		}

	stmt(body)
		{
		return Suneido.Parse(.fn(body))[0]
		}

	Test_quote()
		{
		st = .style()
		Assert(st.Quote(#foo, '"foo"'), is: "#foo")
		Assert(st.Quote(#foo, '"foo"', sym: false), is: '"foo"')
		Assert(st.Quote('a', '"a"'), is: "'a'")
		Assert(st.Quote("class", '"class"'), is: '"class"')
		Assert(st.Quote("_x", '"_x"'), is: '"_x"')
		Assert(st.Quote("a b", "'a b'"), is: '"a b"')
		Assert(st.Quote('say "hi"', '\'say "hi"\''), is: '\'say "hi"\'')
		Assert(st.Quote('x', "`x`"), is: "`x`")
		Assert(st.Quote("a\tb", '"a\\tb"'), is: '"a\\tb"')
		Assert(st.Quote("a +\nb", "'a +\nb'"), is: '"a +\nb"')
		Assert(st.Quote("a +\nb", '"a +\nb"'), is: '"a +\nb"')
		Assert(st.Quote("a\nb", "`a\nb`"), is: "`a\nb`") // backquote stays raw
		Assert(st.Quote('say "hi"\nbye', '\'say "hi"\nbye\'')
			// a " inside pins it
			is: '\'say "hi"\nbye\'')
		}

	Test_words()
		{
		st = .style()
		Assert(st.Plain?("abc du"))
		Assert(st.Plain?("a\tb"), is: false)
		Assert(st.Plain?("a\\b"), is: false)
		Assert(st.BareWord?(#foo))
		Assert(st.BareWord?("class"), is: false)
		Assert(st.BareWord?("2x"), is: false)
		Assert(st.BareWord?(123), is: false)
		Assert(st.ValueSrc("  42 ,"), is: "42")
		Assert(st.ValueSrc(" 42 "), is: "42")
		}

	Test_nodes()
		{
		st = .style()
		e = .stmt("y = true").expr
		Assert(st.BoolConst?(e.rhs, true))
		Assert(st.BoolConst?(e.rhs, false), is: false)
		r = .stmt("y = x[0 .. n]").expr.rhs
		Assert(st.ZeroIdx?(r.from))
		Assert(st.ZeroIdx?(r.to), is: false)
		// a debug call at the margin keeps the margin, an indented one does not
		dbg = .fn("x = 1\nPrint(1)\n\tPrint(2)")
		stmts = Suneido.Parse(dbg)
		Assert(.style(dbg).Debug?(stmts[1].expr))
		Assert(.style(dbg).Debug?(stmts[2].expr), is: false)
		Assert(.style(dbg).Debug?(.stmt("f(1)").expr), is: false)
		Assert(st.Negatable?(.stmt("y = a is b").expr.rhs))
		Assert(st.Negatable?(.stmt("y = a + b").expr.rhs), is: false)
		Assert(st.SwitchExpr(.stmt("y = (x)").expr.rhs).type, is: #Ident)
		Assert(st.SwitchExpr(.stmt("y = (a $ b)").expr.rhs).type, is: #Unary)
		Assert(st.Fnlen(.stmt("Foobar(1)").expr.func), is: 6)
		Assert(st.Fnlen(.stmt("this.Go(1)").expr.func), is: 3)
		Assert(st.Fnlen(.stmt("f().g(1)").expr.func), is: false)
		Assert(st.Shorthand?(.stmt("f(:a)").expr[0]))
		Assert(st.Shorthand?(.stmt("f(a: b)").expr[0]), is: false)
		Assert(st.Simple(.stmt("y = x").expr.rhs))
		Assert(st.Simple(.stmt("f()").expr), is: false)
//		Assert(st.NewThis(.stmt("y = new this()").expr.rhs))
//		Assert(st.NewThis(.stmt("y = new Foo()").expr.rhs), is: false)
		// scalar runs align; a method (index 2) breaks the run
		cls = Suneido.Parse(
			"class\n\t{\n\tx: 1\n\tabcd: 2\n\tGo()" $ "\n\t\t{\n\t\t}\n\ty: 3\n\t}")
		Assert(st.ClassAlign(cls), is: #(4, 4, false, 1))
		}

	Test_tight()
		{
		st = .style()
		add = .stmt("y = i + 1").expr.rhs
		Assert(st.Tight?(add))
		Assert(st.Tight?(add, min: 13), is: false)
		mod = .stmt("y = i % 2").expr.rhs
		Assert(st.Tight?(mod))
		Assert(st.Tight?(mod, min: 14), is: false)
		Assert(st.Tight?(.stmt("y = f() + 1").expr.rhs), is: false)
		}

	Test_unbrace()
		{
		src = .fn("if a\n\t\t{\n\t\tb()\n\t\t}")
		st = AstFmtStyle(src, AstFmtComments(src))
		Assert(st.Unbrace(Suneido.Parse(src)[0].t, false), isnt: false)
		src2 = .fn("if a\n\t\t{\n\t\tb()\n\t\t// note\n\t\t}")
		// a comment before } would be dropped by the skip: keep the braces
		st2 = AstFmtStyle(src2, AstFmtComments(src2))
		Assert(st2.Unbrace(Suneido.Parse(src2)[0].t, false), is: false)
		// unbraced body ending in an open if would capture the else
		src3 = .fn("if a\n\t\t{\n\t\tif b\n\t\t\tc()\n\t\t}\n\telse\n\t\td()")
		st3 = AstFmtStyle(src3, AstFmtComments(src3))
		node3 = Suneido.Parse(src3)[0]
		Assert(st3.Unbrace(node3.t, #Else), is: false)
		Assert(st3.Unbrace(node3.t, false), isnt: false)
		}

	Test_vertical()
		{
		src = .fn("x = #(\n\t\ta: 1,\n\t\tbb: 2)")
		ob = Suneido.Parse(src)[0].expr.rhs.value
		st = AstFmtStyle(src, AstFmtComments(src))
		Assert(st.Vertical?(ob))
		Assert(st.AlignWidth(ob), is: 2)
		src2 = .fn("x = #(a: 1, bb: 2)")
		ob2 = Suneido.Parse(src2)[0].expr.rhs.value
		Assert(AstFmtStyle(src2, AstFmtComments(src2)).Vertical?(ob2), is: false)
		}

	Test_argBreaks()
		{
		src = .fn("f(a, b,\n\t\tc)")
		call = Suneido.Parse(src)[0][0]
		st = .style(src)
		Assert(st.ArgBrokeInSrc?(call, 1), is: false)
		Assert(st.ArgBrokeInSrc?(call, 2))
		Assert(st.VerticalArgs?(call, 2), is: false)
		src2 = .fn("f(a,\n\t\tb,\n\t\tc)")
		Assert(.style(src2).VerticalArgs?(Suneido.Parse(src2)[0][0], 2))
		}

	Test_tableColumns()
		{
		src = "Base\n\t{\n\tc: (\n\t\t(a: 1, b: 200),\n\t\t(a: 22, b: 3),\n\t\t" $
			"(a: 333, b: 44))\n\t}"
		ob = Suneido.Parse(src)[0].value
		Assert(.style(src).TableColumns(ob, 90, 4), is: #(6, 6))
		// too narrow for the padded row
		Assert(.style(src).TableColumns(ob, 20, 4), is: false)
		// mixed keys
		src2 = "Base\n\t{\n\tc: (\n\t\t(a: 1, b: 2),\n\t\t(x: 1, y: 2),\n\t\t" $
			"(a: 1, b: 2))\n\t}"
		Assert(.style(src2).TableColumns(Suneido.Parse(src2)[0].value, 90, 4), is: false)
		// fewer than 3 rows
		src3 = "Base\n\t{\n\tc: (\n\t\t(a: 1, b: 2),\n\t\t(a: 1, b: 2))\n\t}"
		Assert(.style(src3).TableColumns(Suneido.Parse(src3)[0].value, 90, 4), is: false)
		// a non-scalar cell value
		src4 = "Base\n\t{\n\tc: (\n\t\t(a: (1), b: 2),\n\t\t(a: 1, b: 2),\n\t\t" $
			"(a: 1, b: 2))\n\t}"
		Assert(.style(src4).TableColumns(Suneido.Parse(src4)[0].value, 90, 4), is: false)
		}

	Test_verticalNested()
		{
		// a nested table node has no pos of its own
		src = "Base\n\t{\n\tc: (\n\t\t(a: 1),\n\t\t(b: 2))\n\t}"
		ob = Suneido.Parse(src)[0].value
		Assert(ob.pos, is: 0)
		Assert(.style(src).Vertical?(ob))
		// the open delimiter must break, else it stays filled
		src2 = "Base\n\t{\n\tc: ((a: 1),\n\t\t(b: 2))\n\t}"
		ob2 = Suneido.Parse(src2)[0].value
		Assert(ob2.pos, is: 0)
		Assert(.style(src2).Vertical?(ob2), is: false)
		}
	}
