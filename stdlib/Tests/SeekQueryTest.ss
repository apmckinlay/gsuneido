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
			sq = t.SeekQuery(.table)
			Assert(sq.Next() is: #(id: "100"))
			Assert(sq.Next() is: #(id: "101"))
			Assert(sq.Prev() is: #(id: "100"))
			Assert(sq.Prev() is: false)
			sq.Rewind()
			Assert(sq.Prev() is: #(id: "199"))
			}
		}
	Test_Seek()
		{
		Transaction(read:)
			{ |t|
			sq = t.SeekQuery(.table)

			Assert(sq.Seek('id', '15') is: 50)
			Assert(sq.Next() is: #(id: '150'))
			sq.Seek('id', '15')
			Assert(sq.Prev() is: #(id: '149'))
			Assert(sq.Prev() is: #(id: '148'))
			Assert(sq.Next() is: #(id: '149'))
			Assert(sq.Next() is: #(id: '150'))
			Assert(sq.Next() is: #(id: '151'))
			Assert(sq.Prev() is: #(id: '150'))
			Assert(sq.Prev() is: #(id: '149'))
			Assert(sq.Prev() is: #(id: '148'))
			sq.Seek('id', '15')
			Assert(sq.Next() is: #(id: '150'))
			Assert(sq.Prev() is: #(id: '149'))

			sq.Seek('id', '15')
			sq.Rewind()
			Assert(sq.Next() is: #(id: "100"))

			Assert(sq.Seek('id', '') is: 0)
			Assert(sq.Prev() is: false)
			sq.Seek('id', '')
			Assert(sq.Next() is: #(id: "100"))

			Assert(sq.Seek('id', '3') is: 100)
			Assert(sq.Next() is: false)
			sq.Seek('id', '3')
			Assert(sq.Prev() is: #(id: "199"))
			}
		}
	Test_back_and_forth()
		{
		Transaction(read:)
			{ |t|
			sq = t.SeekQuery(.table)
			sq.Seek('id', '125')
			sq.Next()
			for (i = 124; i >= 103; --i)
				Assert(sq.Prev() is: Object(id: String(i)))
			for (i = 104; i <= 130; ++i)
				Assert(sq.Next() is: Object(id: String(i)))
			for (i = 129; i >= 100; --i)
				Assert(sq.Prev() is: Object(id: String(i)))
			}
		}
	Test_Close()
		{
		Transaction(read:)
			{ |t|
			sq = t.SeekQuery(.table)
			sq.SeekQuery_qbefore =  false
			sq.SeekQuery_qafter = false
			sq.SeekQuery_q = false
			// Close() should handle any of the queries that are false
			sq.Close()
			}
		}
	}