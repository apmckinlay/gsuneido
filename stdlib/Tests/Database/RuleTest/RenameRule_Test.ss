// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(name, test_amount, pull, pull_deps) key(name)')
		QueryOutput(table $ ' rename pull to test_simple_pull', [name: 'fred'])
		Assert(Query1(table $ ' rename pull_deps to d').d.Split(',')
			equalsSet: #(name,test_amount), msg: "deps in db")
		Assert(Query1(table).GetDeps('pull').Split(',')
			equalsSet: #(name,test_amount), msg: "deps on retrieved record")
		}
	}