// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup(fieldNumber? = false)
		{
		.VirtualList_Table = .MakeTable('(num, field1) key(num)')
		for i in ..30 /*= 30 records in the test table*/
			{
			rec = Record(num: i)
			if fieldNumber?
				rec.field1 = i * i
			QueryOutput(.VirtualList_Table, rec)
			}

		Suneido.Delete('ForeignKeyTables')
		.WatchTable('slow_queries')
		.teardownModels = Object()
		}

	AddTeardownModel(model)
		{
		.teardownModels.Add(model)
		}

	FakeSaveAndCollapse(rec, row_num, model)
		{
		if rec.vl_expanded_rows isnt ''
			{
			model.SetRecordCollapsed(row_num)
			}
		return true
		}

	Teardown()
		{
		for m in .teardownModels
			m.Destroy()
		super.Teardown()
		Suneido.Delete('ForeignKeyTables')
		}
	}
