// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// marks returns one character per source character: '^' where the
	// character is coloured as a type annotation, '.' everywhere else
	marks(src)
		{
		styles = ScintillaStyle.StyleString(src)
		result = ""
		for i in ..styles.Size()
			result $= styles[i] is '\x07' /*= ANNOTATION */ ? '^' : '.'
		return result
		}

	check(src, expected)
		{
		Assert(.marks(src) is: expected, msg: src)
		}

	checkNone(src)
		{
		Assert(.marks(src).Has?('^') is: false, msg: src)
		}

	styleOf(src, sub)
		{
		return ScintillaStyle.StyleString(src)[src.Find(sub) :: sub.Size()]
		}

	Test_definitions()
		{
		.check("Add(a: number) : boolean {}",
			".....^.^^^^^^..^.^^^^^^^...")
		.check("foo(x: boolean|string) :number|object {}",
			".....^.^^^^^^^^^^^^^^..^^^^^^^^^^^^^^...")
		.check("Origin() : Point {}",
			".........^.^^^^^...")
		.check("function (x: number) {}",
			"...........^.^^^^^^....")
		.check("Configure(.disabled: object) {}",
			"...................^.^^^^^^....")
		.check("Add(a: number)\n\t{\n\treturn a\n\t}",
			".....^.^^^^^^.................")
		}

	Test_blockArgumentsNotColoured()
		{
		.checkNone("QueryApply1('etaweb_users', authHash: oldHash)\n" $
			"\t{|x|\n\tx.Update()\n\t}")
		.checkNone("Foo(a: b)\n\t{ |x| x }")
		.checkNone("Foo(a: b) {|x| x }")
		// a block argument doesn't have to declare parameters
		.checkNone("QueryApply('eta_workers_orders', etaao_num: assocNum, update:)\n" $
			"\t{\n\tit.etapay_readyto_pay? = true\n\tit.Update()\n\t}")
		.checkNone("QueryApply(table, etaao_num: assocNum, update:)\n\t{\n\t}")
		.checkNone("Transaction(update:)\n\t{\n\t}")
		}

	// only what a parameter list can hold, i.e. names with optional type and
	// default, keeps `name(...) { }` a definition rather than a call
	Test_notParameterLists()
		{
		.checkNone("Foo('str', a: b)\n\t{\n\t}")
		.checkNone("Foo(a: b, 'str')\n\t{\n\t}")
		.checkNone("Foo(a: 1)\n\t{\n\t}")
		.checkNone("Foo(num: t.num)\n\t{\n\t}")
		.checkNone("Foo(bar(x): y)\n\t{\n\t}")
		.checkNone("Foo(@args, b: c) {}") // @args has to be the whole list
		.checkNone("Foo(a: number b: string) {}") // missing comma
		}

	Test_defaultParameters()
		{
		.check("Foo(a: number, b = 'x') {}",
			".....^.^^^^^^.............")
		.check("Foo(a: number, b = false) {}",
			".....^.^^^^^^...............")
		.check("Foo(a: number, b = -1) {}",
			".....^.^^^^^^............")
		.check("Foo(_a: number) {}",
			// dynamic parameter
			"......^.^^^^^^....")
		}

	Test_noAnnotations()
		{
		.check("Bar() {}",
			"........")
		.check("Sum(@args) {}",
			".............")
		.check("if (a ? b : c) {}",
			".................")
		}

	Test_callsNotColoured()
		{
		.check("Query1('tables', table: 'test')",
			"...............................")
		.check("Query1(table: tablename)",
			"........................")
		.check("x = Foo(a: 1)",
			".............")
		}

	Test_memberCallsNotColoured()
		{
		.check(".Configure(disabled: false) {}",
			"..............................")
		.check(".Configure(disabled: tablename) {}",
			"..................................")
		}

	Test_valueKeywordsNotColoured()
		{
		.check("Foo(disabled: false) {}",
			".......................")
		.check("Configure(.disabled: false) {}",
			"..............................")
		}

	Test_keepsUnderlyingColour()
		{
		Assert(.styleOf("Add(a: number) {}", #number) is: '\x07'.Repeat(6),
			msg: "type name should be ANNOTATION")
		Assert(.styleOf("Foo(disabled: false) {}", "false") is: '\x04'.Repeat(5),
			msg: "false should stay KEYWORD")
		}

	Test_conditionsNotColoured()
		{
		.checkNone("if QueryEmpty?('gl_transactions', gltran_subsys_num: t.num,\n" $
			"\tgltran_subsys_detail_num: t.num)\n\t{\n\treturn\n\t}")
		.checkNone("if not QueryEmpty?(query, num: t.num)\n\t{\n\t}")
		.checkNone("while Foo(a: 1)\n\t{\n\t}")
		.checkNone("x = Foo(a: 1)\n\t{\n\t}")
		}

	Test_definitionsNotAtStart()
		{
		.check("x: 5\nAdd(a: number) {}",
			"..........^.^^^^^^....")
		.check("Foo(a: number,\n\tb: string) {}",
			".....^.^^^^^^....^.^^^^^^....")
		}

	Test_unionsWithValueKeywords()
		{
		.check("ObjectClash(o :object|false)\n\t{\n\t}",
			"..............^^^^^^^^^^^^^.......")
		.check("Foo(x: object|false) :string|false {}",
			".....^.^^^^^^^^^^^^..^^^^^^^^^^^^^...")
		}
	}
