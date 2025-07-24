// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_with_key()
		{
		table = .MakeTable('(a b c) key(a)')
		.testNone3ways(table)

		QueryOutput(table, #(a: 'a' b: 'b' c: 'c'))
		.testOne(table)
		.testOne2ways(table)

		QueryOutput(table, #(a: 'aa' b: 'b' c: 'c'))
		Assert({ QueryApply1(table){|unused|} }, throws: 'not unique')

		.testOne2ways(table)
		}
	Test_empty_key()
		{
		table = .MakeTable('(a b c) key()')
		.testNone3ways(table)

		QueryOutput(table, #(a: 'a' b: 'b' c: 'c'))
		.testOne2ways(table)
		}

	testNone3ways(table)
		{
		.testNone(table)
		.testNone(table $ ' where a is "a"')
		.testNone(table, a: "a")
		}
	testNone(@args)
		{
		args.block =  {|unused| throw "should not get any records" }
		QueryApply1(@args)
		}
	testOne2ways(table)
		{
		.testOne(table $ ' where a is "a"')
		.testOne(table, a: "a")
		}
	testOne(@args)
		{
		got = false
		args.block =  {|x| got = x }
		QueryApply1(@args)
		Assert(got isnt: false)
		}
	}