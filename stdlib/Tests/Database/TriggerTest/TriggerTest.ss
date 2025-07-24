// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		tt1 = .MakeTable("(a, b, c) key(a)")
		tt2 = .MakeTable("(when, oldrec, newrec) key(when)")
		.MakeLibraryRecord(Record(name: 'Trigger_' $ tt1,
			text: 'function (transaction, oldrec, newrec)
				{
				Assert(oldrec is false or oldrec.Transaction() is transaction)
				Assert(newrec is false or newrec.Transaction() is transaction)
				Transaction(update:)
					{|t|
					t.QueryOutput(' $ Display(tt2) $ ',	Record(when: Timestamp(),
						:oldrec, :newrec))
					}
				}'))

		Transaction(update:)
			{ |t|
			q = t.Query(tt1)

			q.Output(Record(a: 1, b: 2, c: 3))
			Assert(QueryLast(tt2 $ ' sort when').Project(#(oldrec, newrec)) is:
				#(oldrec: false, newrec: (a: 1, b: 2, c: 3)))

			x = q.Next()
			Assert(x is: #(a: 1, b: 2, c: 3))
			x.c = 33
			x.Update()
			Assert(QueryLast(tt2 $ ' sort when').Project(#(oldrec, newrec)) is:
				#(oldrec: (a: 1, b: 2, c: 3), newrec: (a: 1, b: 2, c: 33)))

			q.Rewind()
			x = q.Next()
			Assert(x is: #(a: 1, b: 2, c: 33))
			x.Delete()
			Assert(QueryLast(tt2 $ ' sort when').Project(#(oldrec, newrec)) is:
				#(oldrec: (a: 1, b: 2, c: 33), newrec: false))
			}
		}
	}