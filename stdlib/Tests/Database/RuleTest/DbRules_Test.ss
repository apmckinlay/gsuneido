// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		try Database("destroy testrule_tmp")
		Database("create testrule_tmp
			(name,
			Test_simple_pull,
			test_amount,
			Test_running_total,
			Test_simple_default,
			Test_local_default)
			key(name)")
		Transaction(update:)
			{|t|
			q = t.Query("testrule_tmp")
			x = Record()
			Assert(x.name is: '')
			Assert(x.test_amount is: '')
			Assert(x.test_simple_pull is: " and ")
			Assert(x.test_running_total is: "")
			Assert(x.test_simple_default is: 'simple default')
			Assert(x.test_local_default is: '')

			x.name = 'joe'
			x.test_amount = 99
			Assert(x.test_simple_pull is: "joe and 99")
			Assert(x.test_amount__valid is: false)
			x.test_amount = 3
			Assert(x.test_amount__valid)
			Assert(x.name is: 'joe')
			Assert(x.test_amount is: 3)
			Assert(x.test_simple_pull is: "joe and 3")
			Assert(x.test_running_total is: "")
			Assert(x.test_simple_default is: 'simple default')
			Assert(x.test_local_default is: 'j')
			q.Output(x)

			x = q.Next()
			Assert(x.name is: 'joe')
			Assert(x.test_amount is: 3)
			Assert(x.test_simple_pull is: "joe and 3", msg: "stored")
			Assert(x.test_running_total is: "joe ")
			Assert(x.test_simple_default is: 'simple default')
			Assert(x.test_local_default is: 'j')

			x.name = 'sue'
			Assert(x.test_simple_pull is: "sue and 3")
			Assert(x.test_local_default is: 's')
			Assert(x.test_running_total is: "joe ")
			x.Update()

			q.Output(#(name: 'tom', test_amount: 456))
			x.Invalidate("test_running_total")
			Assert(x.test_running_total is: "sue tom ")

			// test indirect invalidate
			Assert(x.test_indirect is: 8)
			y = Object()
			x.Observer({|member| y[member] = true; })
			x.Invalidate("test_running_total")
			Assert(y is: #(test_running_total:, test_indirect:))
			}
		}
	Test_copy_deps()
		{
		r = Record(name: 'fred', test_amount: 123)
		Assert(r.test_simple_pull is: "fred and 123")
		r2 = r.Copy()
		r2.name = 'joe'
		Assert(r2.test_simple_pull is: "joe and 123")
		}
	Test_getdeps()
		{
		r = Record(name: 'fred', test_amount: 123)
		r.test_simple_pull
		Assert(r.GetDeps("test_simple_pull").Split(',') equalsSet: #(name,test_amount))
		}
	Test_setdeps()
		{
		r = Record(name: 'fred', test_amount: 123, test_simple_pull: "fred and 123")
		r.name = "joe"
		Assert(r.test_simple_pull is: "fred and 123")
		r.SetDeps("test_simple_pull", 'name,test_amount')
		r.name = "sue"
		Assert(r.test_simple_pull is: "sue and 123")
		}

	Test_where()
		{
		table = .MakeTable('(name, test_amount, Test_simple_pull) key(name)'
			#(name: 'joe', test_amount: 100),
			#(name: 'sue', test_amount: 200),
			#(name: 'ann', test_amount: 300),
			#(name: 'bob', test_amount: 400))

		WithQuery(table $ " where test_simple_pull.Has?('o')")
			{ |q|
			Assert(q.Next().name is: 'bob')
			Assert(q.Next().name is: 'joe')
			Assert(q.Next() is: false)
			}
		}
	Test_saved_deps() // test that dependencies persist
		{
		table = .MakeTable('(name, test_amount, test_simple_pull, test_simple_pull_deps)
			key(name)')
		QueryOutput(table, r = Record(name: 'fred', test_amount: 100))
		Assert(.getDeps(r, #test_simple_pull)
			equalsSet: #(name, test_amount), msg: "output:")

		r = Query1(table $ ' rename test_simple_pull_deps to deps')
		Assert(r.deps.Split(',') equalsSet: #(name, test_amount), msg: "stored:")

		QueryApply1(table)
			{ |x|
			Assert(.getDeps(x, #test_simple_pull)
				equalsSet: #(name, test_amount), msg: "input:")
			x.name = 'tom'
			x.Update()
			}

		QueryDo("update " $ table $ " set test_amount = 200")
		Assert(Query1(table).test_simple_pull is: 'tom and 200')
		}
	getDeps(record, field)
		{
		return record.GetDeps(field).Split(',')
		}
	Test_timestamp()
		{
		table = .MakeTable('(test_timestamp, other) key(test_timestamp)')
		QueryOutput(table, Record(other: 'fred'))

		timestamp = Query1(table).test_timestamp
		Assert(Date?(timestamp))

		QueryApply1(table)
			{ |x|
			x.other = 'tom'
			x.Update()
			}

		Assert(Query1(table).test_timestamp is: timestamp)
		}
	Test_extend_rules()
		{
		table = .MakeTable('(name, test_amount) key(name)'
			#(name: 'joe', test_amount: 100),
			#(name: 'sue', test_amount: 200),
			#(name: 'ann', test_amount: 300),
			#(name: 'bob', test_amount: 400))
		QueryApply(table $ ' extend test_simple_pull')
			{ |x|
			Assert(x.test_simple_pull is: x.name $ " and " $ x.test_amount)
			}
		Assert(QueryCount(table $ ' extend test_simple_pull
				where test_simple_pull.Has?("e")') is: 2)
		}

	Teardown()
		{
		Database("destroy testrule_tmp")
		super.Teardown()
		}
	}