// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_sort()
		{
		s = QueryStrategy('stdlib
			sort parent, name, lib_committed')
		Assert(s has: 'tempindex')

		s = QueryStrategy('stdlib
			where lib_committed = #20030902.094223671
			sort parent, name, lib_committed')
		Assert(s hasnt: 'tempindex')
		Assert(s has: 'stdlib^(parent,name)')
		}
	Test_project()
		{
		s = QueryStrategy('stdlib
			project parent, name, lib_committed
			/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */')
		Assert(s hasnt: 'project-seq')

		s = QueryStrategy('stdlib
			where lib_committed = #20030902.094223671
			project parent, name, lib_committed
			/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */')
		Assert(s has: 'project-seq')
		Assert(s has: 'stdlib^(parent,name)')
		}
	}
