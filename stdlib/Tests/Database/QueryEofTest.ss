// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_create()
		{
		table = .MakeTable('(id, date) key(id)',
			#(id: 1, date: 991111),
			#(id: 2, date: 991111))

		WithQuery(table)
			{|q|
			q.Next()
			q.Next()
			Assert(q.Next() is: false)
			Assert(q.Next() is: false, msg: "next after eof")
			}

		WithQuery(table)
			{|q|
			q.Prev()
			q.Prev()
			Assert(q.Prev() is: false)
			Assert(q.Prev() is: false, msg: "prev after eof")
			}
		}
	Test_eoftwice()
		{
		WithQuery('stdlib')
			{|q|
			.check(q)
			}

		WithQuery('stdlib where name > "M"')
			{|q|
			.check(q)
			}

		WithQuery('stdlib sort lib_committed')
			{|q|
			Assert(q.Strategy() has: 'tempindex')
			.check(q)
			}

		WithQuery('(stdlib union stdlib) sort lib_committed
			/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */')
			{|q|
			Assert(q.Strategy() has: 'tempindex')
			.check(q)
			}
		}

	check(q)
		{
		first = q.Next()
		Assert(first isnt: false)
		Assert(q.Prev() is: false)
		Assert(q.Next() is: first)
		Assert(q.Next() isnt: false)
		Assert(q.Prev() is: first)
		}
	}