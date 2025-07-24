// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_create()
		{
		master = .MakeTable('(id) key(id)')
		tran = .MakeTable('(id, date) key(id,date)',
			[id: 1, date: 991111])

		Assert({ Database("alter " $ tran $ " create index(id) in " $ master) }
			throws: "blocked by foreign key")
		QueryOutput(master, [id: 1])
		Database("alter " $ tran $ " create index(id) in " $ master)
		}
	Test_block()
		{
		table1 = .MakeTable('(id, name) key(id)')
		table2 = .MakeTable('(hid, date) key(hid,date)
			index(hid) in ' $ table1 $ ' (id)')
		Transaction(update:)
			{|t|
			cq = t.Query(table1)
			cq.Output(#(id: 1, name: 'one'))
			cq.Output(#(id: 2, name: 'two'))
			hq = t.Query(table2)
			// successful source add
			hq.Output(#(hid: 2, date: 991126))
			// blocked source add
			Assert({ hq.Output(#(hid: 3, date: 991126)) }
				throws: "blocked by foreign key")
			cq.Output(#(id: 3, name: 'three'))
			// successful source add
			hq.Output(#(hid: 3, date: 991126))
			Assert(x = cq.Next() isnt: false, msg: "get cus 1")
			x.Delete()
			// blocked source add
			Assert({ hq.Output(#(hid: 1, date: 991126)) }
				throws: "blocked by foreign key")
			// blocked source update
			Assert(x = hq.Next() isnt: false, msg: "get his 2")
			x.hid = 4; x.date = 991127
			Assert({ x.Update() }, throws: "blocked by foreign key")
			// successful source update
			Assert(x = hq.Next() isnt: false, msg: "get his 3")
			x.hid = 2; x.date = 991127
			x.Update()
			// blocked target delete
			Assert(x = cq.Next() isnt: false, msg: "get cus 2")
			Assert({ x.Delete() } throws: "blocked by foreign key")
			// blocked target update
			hq.Output(#(hid: 3, date: 991127))
			Assert(x = cq.Next() isnt: false, msg: "get cus 3")
			x.id = 4
			Assert({ x.Update() } throws: "blocked by foreign key")
			// successful target update
			Assert(t.QueryDo("delete " $ table2 $ " where hid = 3") is: 1)
			Assert(t.QueryDo("update " $ table1 $ " where id = 3 set id = 4") is: 1)
			}
		}
	// the trigger is added to catch a gSuneido bug where in cascade delete/update the
	// line table trigger can't find its corresponding header record within the same
	// transaction (t.Query1)
	makeTrigger(headerTable, lineTable)
		{
		.MakeLibraryRecord(
			Object(
				name: 'Trigger_' $ lineTable,
				text: "function (t, oldrec, newrec)
					{
					if oldrec is false and newrec isnt false
						return

					headerRec = t.Query1(" $ Display(headerTable) $ ", num: oldrec.num)
					Assert(headerRec isnt false)
					}"))
		}
	Test_cascade_delete()
		{
		table1 = .MakeTable('(num) key(num)')
		table2 = .MakeTable('(num,part) key(num,part)
			index(num) in ' $ table1 $ ' cascade')
		.makeTrigger(table1, table2)
		Transaction(update:)
			{|t|
			q1 = t.Query(table1)
			q2 = t.Query(table2)
			q1.Output(#(num: 1))
				q2.Output(#(num: 1, part: 1))
			q1.Output(#(num: 2))
				q2.Output(#(num: 2, part: 1))
				q2.Output(#(num: 2, part: 2))
				q2.Output(#(num: 2, part: 3))
			q1.Output(#(num: 3))

			Assert(t.QueryDo("delete " $ table1 $ " where num = 1") is: 1)
			Assert(t.Query(table2 $ " where num = 1").Next() is: false, msg: "get lin 1")

			Assert(t.QueryDo("delete " $ table1 $ " where num = 2") is: 1)
			Assert(t.Query(table2 $ " where num = 2").Next() is: false, msg: "get lin 2")

			Assert(t.QueryDo("delete " $ table1 $ " where num = 3") is: 1)
			Assert(t.Query(table2 $ " where num = 3").Next() is: false, msg: "get lin 3")
			}
		}
	Test_cascade_update()
		{
		table1 = .MakeTable('(num) key(num)')
		table2 = .MakeTable('(num,part) key(num,part)
			index(num) in ' $ table1 $ ' cascade')
		.makeTrigger(table1, table2)
		Transaction(update:)
			{|t|
			q1 = t.Query(table1)
			q2 = t.Query(table2)
			q1.Output(#(num: 1))
				q2.Output(#(num: 1, part: 1))
			q1.Output(#(num: 2))
				q2.Output(#(num: 2, part: 1))
				q2.Output(#(num: 2, part: 2))
				q2.Output(#(num: 2, part: 3))
			q1.Output(#(num: 3))

			Assert(t.QueryDo("update " $ table1 $ " where num = 2 set num = 4") is: 1)
			Assert(t.Query(table2 $ " where num = 2").Next() is: false, msg: "get lin 2")
			q = t.Query(table2 $ " where num = 4")
			for (i = 1; i <= 3; ++i)
				{
				Assert(x = q.Next() isnt: false)
				Assert(x.part is: i)
				}

			Assert(t.QueryDo("update " $ table1 $ " where num = 3 set num = 6") is: 1)
			}
		}
	Test_recursive()
		{
		try Database("drop recursive_foreign_key_test")
		Database("create recursive_foreign_key_test (employee, manager)
			key(employee)
			index(manager) in recursive_foreign_key_test(employee)")
		Database("drop recursive_foreign_key_test")
		}

	Test_conflict()
		{
		for act1 in #(head_update, head_delete)
			for act2 in #(line_output, line_update)
				{
				.test(act1, act2)
				.test(act2, act1)
				}
		}
	test(act1, act2)
		{
		head = .MakeTable("(head) key(head)")
		line = .MakeTable("(line, head) key(line) index(head) in " $ head)
		QueryOutput(head, [head: 1])
		QueryOutput(line, [head: '', line: 3])
		Assert({
			Transaction(update:)
				{|t|
				.perform(head, line, t, act1)
				Transaction(update:)
					{|t2|
					.perform(head, line, t2, act2)
					}
				}
			}, throws: 'conflict')
		QueryApply(line)
			{|y|
			if y.head isnt "" and Query1(head, head: y.head) is false
				throw "line " $ y.line $ " has no head"
			}
		}
	perform(head, line, t, action)
		{
		switch action
			{
		case 'head_delete':
			x = t.Query1(head, head: 1)
			x.Delete()
		case 'head_update':
			x = t.Query1(head, head: 1)
			x.head = 9
			x.Update()
		case 'line_update':
			x = t.Query1(line, line: 3)
			x.head = 1
			x.Update()
		case 'line_output':
			t.QueryOutput(line, [line: 1.1, head: 1])
			}
		}

	Test_empty()
		{
		table1 = .MakeTable('(k) key(k)')
		table2 = .MakeTable('(a, k) key(a) index(k) in ' $ table1)
		QueryOutput(table2, Record(a: 1, k: ''))
		QueryOutput(table1, Record(k: ''))
		QueryDo('delete ' $ table1)
		}

	Test_foreign_key()
		{
		cus = .MakeTable('(c) key(c)', [c: 1])
		for mode in #('cascade', 'cascade update')
			{
			loc = .MakeTable('(k, c) key(k) index (c) in ' $ cus $ ' ' $ mode,
				#(k: a, c: 1))
			Assert({ QueryOutput(loc, #(k: b, c: 2)) }
				throws: "blocked by foreign key")
			Assert({ QueryDo('update ' $ loc $ ' set c = 2') }
				throws: "blocked by foreign key")
			}
		}
	Test_gsuneido_bug()
		{
		k = "<\x00\x00>"
		main = .MakeTable("(k) key(k)", [:k])
		tbls = [
			.MakeTable("(a,k) key(a) index(k) in " $ main),
			.MakeTable("(k) key(k) in " $ main),
			.MakeTable("(a,k) key(a) index(k) in " $ main),
			.MakeTable("(k) key(k) in " $ main),
			.MakeTable("(a,k) key(a) index(k) in " $ main),
			.MakeTable("(k) key(k) in " $ main)
			]
		for tbl in tbls
			{
			Transaction(update:)
				{|t|
				t.QueryOutput(tbl, [:k])
				x = t.Query1(main)
				Assert({ x.Delete() } throws: "blocked by foreign key")
				t.Rollback()
				}
			}
		}
	}
