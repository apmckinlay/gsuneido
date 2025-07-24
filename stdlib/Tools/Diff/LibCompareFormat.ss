// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Generator
	{
	New(.t1, .t2)
		{
		.book = BookTable?(t1) and BookTable?(t2)
		.list = LibCompare(.t1, .t2)
		.i = 0
		}
	Next()
		{
		if .i >= .list.Size()
			return false
		item = .list[.i++]
		record_num = item[1]
		record_name = item[2]
		switch (item[0])
			{
		case '+' :
			// Print the library record from the second library
			lib_record = Query1(.t2, num: record_num)
			Assert(lib_record isnt: false)
			return _report.Construct(Object('Library', lib_record.name, lib_record.text))
		case '-' :
			// print nothing
			return _report.Construct(Object('Text', ''))
		case '#' :
			// Use Diff on the two text fields and print results
			lib_record1 = Query1(.t1, group: -1, name: record_name)
			lib_record2 = Query1(.t2, group: -1, name: record_name)
			Assert(lib_record1 isnt false and lib_record2 isnt false)
			return _report.Construct(
				Object('Diff', .text(lib_record1), .text(lib_record2), lib_record1.name))
		default:
			throw "shouldn't get here"
			}
		}

	text(rec)
		{
		return .book
			? "path:" $ rec.path $ ", order:" $	rec.order $ '\r\n' $ rec.text
			: rec.lib_current_text
		}
	}
