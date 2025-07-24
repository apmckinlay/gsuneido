// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.table = .MakeTable('(id) key(id)')
		for (i = 100; i < 200; ++i)
			QueryOutput(.table, Record(id: String(i)))
		}
	Test_Main()
		{
		Transaction(read:)
			{ |t|
			sq = SeekCursor(.table)
			Assert(sq.Next(t) is: #(id: "100"))
			Assert(sq.Next(t) is: #(id: "101"))
			Assert(sq.Prev(t) is: #(id: "100"))
			Assert(sq.Prev(t) is: false)
			sq.Rewind()
			Assert(sq.Prev(t) is: #(id: "199"))
			sq.Close()
			}
		}
	Test_Seek()
		{
		Transaction(read:)
			{ |t|
			sq = SeekCursor(.table)

			sq.Seek('id', '15')
			Assert(sq.Next(t) is: #(id: '150'))
			sq.Seek('id', '15')
			Assert(sq.Prev(t) is: #(id: '149'))
			Assert(sq.Prev(t) is: #(id: '148'))
			Assert(sq.Next(t) is: #(id: '149'))
			Assert(sq.Next(t) is: #(id: '150'))
			Assert(sq.Next(t) is: #(id: '151'))
			Assert(sq.Prev(t) is: #(id: '150'))
			Assert(sq.Prev(t) is: #(id: '149'))
			Assert(sq.Prev(t) is: #(id: '148'))
			sq.Seek('id', '15')
			Assert(sq.Next(t) is: #(id: '150'))
			Assert(sq.Prev(t) is: #(id: '149'))

			sq.Seek('id', '15')
			sq.Rewind()
			Assert(sq.Next(t) is: #(id: "100"))

			sq.Seek('id', '')
			Assert(sq.Prev(t) is: false)
			sq.Seek('id', '')
			Assert(sq.Next(t) is: #(id: "100"))

			sq.Seek('id', '3')
			Assert(sq.Next(t) is: false)
			sq.Seek('id', '3')
			Assert(sq.Prev(t) is: #(id: "199"))

			sq.Close()
			}
		}
	Test_seek_bad_query()
		{
		rec = [id: 1, fld: 1]
		.table = .MakeTable('(id, fld) key(id)', rec)
		sq = SeekCursor(.table $ ' sort id')
		Assert({ sq.Seek('fld', '') } throws: 'invalid query')
		Transaction(read:)
			{ |t|
			Assert(sq.Next(t) is: rec)
			}
		sq.Close()
		}
	Test_back_and_forth()
		{
		Transaction(read:)
			{ |t|
			sq = SeekCursor(.table)
			sq.Seek('id', '125')
			sq.Next(t)
			for (i = 124; i >= 103; --i)
				Assert(sq.Prev(t) is: Object(id: String(i)))
			for (i = 104; i <= 130; ++i)
				Assert(sq.Next(t) is: Object(id: String(i)))
			for (i = 129; i >= 100; --i)
				Assert(sq.Prev(t) is: Object(id: String(i)))
			sq.Close()
			}
		}
	}