// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
RecordFormat
	{
	New(rec, field_ob)
		{
		super(rec, .fields(field_ob), font: Object(size: '+2'))
		}

	fields(field_ob)
		{
		list = Object()
		for field in _report.Params.choosefields.Split(",")
			if '' isnt field = field.Trim()
				list.Add(field_ob.GetDefault(field, field))
		return list
		}
	}