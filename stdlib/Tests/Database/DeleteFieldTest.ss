// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table = .MakeTable('(a, R, b, x, bb, c) key(c)')
		Database("alter " $ table $ " drop (b, bb)")
		Transaction(update:)
			{ |t|
			q = t.Query(table)

			x = Record(a: 1, c: 3)
			q.Output(x)
			Assert(q.Next() is: x)

			x = Record(a: 1, c: 4)
			q.Output(x)
			q.Rewind()
			Assert(q.Prev() is: x)
			}
		}
	Test_two()
		{
		table = .MakeTable('
			(test_id, test_customers, test_ref, test_desc, test_created_date,
			test_created_by, test_updated_date, test_updated_by, test_option,
			test_release_date, test_completed_date, test_minor?, test_origin,
			test_selected, test_bumped, test_type, test_dead, test_release_notes)
			key(test_id)',
			[test_id: 123])

		Database('alter ' $ table $ ' drop(test_minor?)')

		Database('alter ' $ table $ ' drop(test_customers)')
		Transaction(update:)
			{ |t|
			x = t.QueryFirst(table $ ' extend test_status = 1 sort test_id')
			x.test_type = 'Change'
			x.Update()
			}
		Assert(Query1(table, test_id: 123) is: [test_id: 123, test_type: 'Change'])
		}
	}
