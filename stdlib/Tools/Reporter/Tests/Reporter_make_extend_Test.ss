// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FormulasWhichExtend()
		{
		w = Reporter_make_extend.Reporter_make_extend_which_extend

		deps = #(none: (), before: (x), after: (count),
			indirect_before: (before), indirect_after: (after),
			before2: (x, force_before), force_before: (id)
			after2: (total_amount, force_after), force_after: (desc))
		b = #(id, desc, x, y, amount)
		sf = #(count, total_amount)
		sb = #(id, desc)
		which = w(deps, b, sb, sf)
		Assert(which is: #(before: 0, after: 1,
			indirect_before: 0, indirect_after: 1,
			before2: 0, force_before: 0, after2: 1, force_after: 1,
			desc: 0, id: 0))

		b = #(date, amount)
		sb = #(month)
		sf = #(count, total_amount, total_days)
		deps = #(month: (), days: ())
		which = w(deps, b, sb, sf)
		Assert(which is: #(month: 0, days: 0))
		}

	Test_compilable?()
		{
		buildFn = Reporter_make_extend.Reporter_make_extend_buildFormulaFunction
		compilable? = Reporter_make_extend.Reporter_make_extend_compilable?

		fn = buildFn(#(a, b), 'a + b', 'Test Formula')
		Assert(compilable?(fn))

		fn = buildFn(#(), '2 + 5', 'Test Formula')
		Assert(compilable?(fn))

		fn = buildFn(#(a, b), 'a + ', 'Test Formula')
		Assert(compilable?(fn) is: false)

		fn = buildFn(#(), 'if', 'Test Formula')
		Assert(compilable?(fn) is: false)
		}

	Test_buildFormulaFunction()
		{
		buildFn = Reporter_make_extend.Reporter_make_extend_buildFormulaFunction
		fn = buildFn(#(a, b), 'a + b', 'Test Formula')
		Assert(fn.Compile()(1, 2) is: 3)

		fn = buildFn(#(), 'Max', 'Test Formula')
		fn.Compile()() // should not throw error
		}
	}