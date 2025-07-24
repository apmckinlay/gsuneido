// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.table = .MakeTable('(k, x) key(k)',
			[k: 0], [k: 1], [k: 2], [k: 3], [k: 4], [k: 5])
		}

	Test_switch_direction()
		{
		x = false
		Cursor(.table)
			{|c|
			check =
				{|t, dir, i|
				x = c[dir](t)
				Assert(x is: [k: i])
				}
			Transaction(read:)
				{|t|
				check(t, #Prev, 5)
				check(t, #Prev, 4)
				}
			Transaction(read:)
				{|t|
				check(t, #Prev, 3)
				check(t, #Next, 4)
				check(t, #Next, 5)
				}
			Transaction(read:)
				{|t|
				check(t, #Prev, 4)
				check(t, #Prev, 3)
				check(t, #Prev, 2)
				}

			c.Rewind()
			Transaction(update:)
				{|t|
				check(t, #Next, 0)
				check(t, #Next, 1)
				check(t, #Prev, 0)
				x.x++
				x.Update()
				check(t, #Next, 1)
				}
			}
		}

	Test_update_and_switch_directions()
		{
		Cursor(.table)
			{|c|
			next = { Transaction(read:) {|t| c.Next(t).k } }
			Assert(next() is: 0)
			Assert(next() is: 1)
			Assert(next() is: 2)
			Transaction(update:)
				{|t|
				x = c.Prev(t)
				Assert(x.k is: 1)
				x.x++
				x.Update()
				}
			Assert(next() is: 2)
			}
		}

	Test_delete_and_switch_directions()
		{
		table = .MakeTable('(k) key(k)',
			[k: 0], [k: 1], [k: 2], [k: 3], [k: 4], [k: 5])
		Cursor(table)
			{|c|
			next = { Transaction(read:) {|t| c.Next(t).k } }
			Assert(next() is: 0)
			Assert(next() is: 1)
			Assert(next() is: 2)
			Transaction(update:)
				{|t|
				x = c.Prev(t)
				Assert(x.k is: 1)
				x.Delete()
				}
			Assert(next() is: 2)
			}
		}

	Test_record_transaction()
		{
		Cursor(.table)
			{|c|
			Transaction(read:)
				{|t|
				x = c.Next(t)
				Assert(x.Transaction() isnt: false)
				}
			}
		}

	Test_modification_during_iteration()
		{
		for d in #((Next, 0, 1, 2), (Prev, 5, 4, 3))
			for op in #(Update, Delete)
				for f in [.test1, .test2, .test3, .test4]
					{
					f(d[0], op, d[1], d[2], d[3])
					if op is 'Delete'
						.table = .MakeTable('(k, x) key(k)',
							[k: 0], [k: 1], [k: 2], [k: 3], [k: 4], [k: 5])
					}
		}

	test1(dir, op, x, y, z)
		{
		Cursor(.table)
			{|c|
			Transaction(update:)
				{|t|
				get = { c[dir](t).k }
				Assert(get() is: x)
				x = c[dir](t)
				Assert(x.k is: y)
				x.x++
				x[op]()
				Assert(get() is: z)
				}
			}
		}

	test2(dir, op, x, y, z)
		{
		Cursor(.table)
			{|c|
			get = { Transaction(read:) {|t| c[dir](t).k } }
			Assert(get() is: x)
			Transaction(update:)
				{|t|
				x = c[dir](t)
				Assert(x.k is: y)
				x.x++
				x[op]()
				}
			Assert(get() is: z)
			}
		}

	test3(dir, op, x, y, z)
		{
		Cursor(.table)
			{|c|
			get = { Transaction(read:) {|t| c[dir](t).k } }
			Assert(get() is: x)
			Transaction(update:)
				{|t|
				x = c[dir](t)
				Assert(x.k is: y)
				x.x++
				x[op]()
				Assert(c[dir](t).k is: z)
				}
			}
		}

	test4(dir, op, x, y, z)
		{
		Cursor(.table)
			{|c|
			get = { Transaction(read:) {|t| c[dir](t).k } }
			Assert(get() is: x)
			Assert(get() is: y)
			Transaction(update:)
				{|t|
				x = t.Query1(.table, k: 1)
				x.x++
				x[op]()
				}
			Assert(get() is: z)
			}
		}
	Test_bug()
		{
		tbl = .MakeTable('(k) key(k)',
			[k: 0], [k: 1], [k: 2], [k: 3], [k: 4], [k: 5])
		Cursor(tbl)
			{|c|
			Transaction(update:)
				{|t|
				Assert(c.Next(t) is: #(k: 0))
				Assert(c.Next(t) is: #(k: 1))
				Assert(t.QueryDo('delete ' $ tbl $ ' where k is 0') is: 1)
				Assert(c.Next(t) is: #(k: 2))
				Assert(c.Prev(t) is: #(k: 1))
				Assert(c.Prev(t) is: false)
				}
			}
		}

	}