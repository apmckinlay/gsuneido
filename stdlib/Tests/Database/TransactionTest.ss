// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_complete_within_block()
		{
		Transaction(read:)
			{ |t|
			t.Complete()
			}
		}
	Test_rollback_within_block()
		{
		Transaction(read:)
			{ |t|
			t.Rollback()
			}
		}

	Test_ended()
		{
		t = Transaction(read:)
		Assert(not t.Ended?())
		t.Complete()
		Assert(t.Ended?())

		t = Transaction(update:)
		Assert(not t.Ended?())
		t.Rollback()
		Assert(t.Ended?())
		}

	Test_update()
		{
		Transaction(update:)
			{ |t| Assert(t.Update?()) }
		Transaction(read:)
			{ |t| Assert(not t.Update?()) }
		}

	Test_conflict()
		{
		table = .MakeTable("(a,b) key(a)", [a: 1, b: 2])

		Assert({
			Transaction(update:)
				{|t|
				x = t.QueryFirst(table $ ' sort a')
				Transaction(update:)
					{|t2|
					t2.QueryDo("update " $ table $ " set b = 22")
					}
				// need an update or tran can't fail
				x.b = 3
				x.Update()
				}
			} throws: "conflict", msg: 1)

		Assert({
			Transaction(update:)
				{|t|
				t.QueryFirst(table $ ' sort a')
				Transaction(update:)
					{|t2|
					t2.QueryDo("update " $ table $ " set b = 222")
					}
				// need an update or tran can't fail
				t.QueryOutput(table, [a: 3, b: 4])
				}
			} throws: "conflict", msg: 2)

		Assert({
			Transaction(update:)
				{|t|
				t.QueryFirst(table $ ' sort a')
				Transaction(update:)
					{|t2|
					t2.QueryOutput(table, [a: 0, b: 0])
					}
				// need an update or tran can't fail
				t.QueryOutput(table, [a: 3, b: 4])
				}
			} throws: "conflict", msg: 3)
		}
	Test_null_pointer_exception()
		{
		Assert(Catch({
			Transaction(update:)
				{|t|
				t.QueryApply("stdlib")
					{|unused|
					t.Rollback()
					}
				} })
			matches: "cannot use a completed Transaction|" $
				"can't Rollback ended transaction|" $
				"can't use ended transaction")
		}
	Test_status()
		{
		t = Transaction(update:)
		t.Complete()
		t.Complete() // ok

		t = Transaction(update:)
		t.Rollback()
		t.Rollback() // ok

		t = Transaction(update:)
		t.Complete()
		Assert({ t.Rollback() } throws: "already completed")

		t = Transaction(update:)
		t.Rollback()
		Assert({ t.Complete() } throws: "already aborted")

		t = Transaction(read:)
		q = t.Query("tables")
		q.Close()
		Assert({ q.Close() } throws: "can't use closed query")
		t.Complete()

		t = Transaction(read:)
		c = Cursor("tables")
		c.Close()
		Assert({ c.Close() } throws: "can't use closed cursor")
		t.Complete()

		t = Transaction(read:)
		q = t.Query("tables")
		t.Complete()
		Assert({ q.Next() } throws: "can't use ended transaction")

		t = Transaction(read:)
		c = Cursor("tables")
		t.Complete()
		Assert({ c.Next(t) } throws: "can't use ended transaction")
		c.Close()
		}
	}
