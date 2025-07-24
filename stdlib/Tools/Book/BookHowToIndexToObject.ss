// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (indexTable, lib, newLibRecordName)
	{
	new_ob = QueryAll(indexTable)
	nextnum = Query1(lib $ " summarize max num")
	num = nextnum is false ? 1 : ++nextnum.max_num
	QueryOutput(lib, Record(
		num: ++num, parent: 0, group: -1,
		name: newLibRecordName,
		text: Display(new_ob)))
	}