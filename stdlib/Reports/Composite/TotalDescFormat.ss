// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
TotalFormat
	{
	New(item, skip = .16, w = false)
		{
		super(Object('Horz', 'Hfill',
			ob = Object?(item) ? item.Copy().Add(w, at: #w) : item), skip, noline:)
		.item = _report.Construct(ob)
		}
	ExportCSV(data = '')
		{
		if Record?(data) and .item.Member?("Field")
			data = data[.item.Field]
		return .item.ExportCSV(data)
		}
	}