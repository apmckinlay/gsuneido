// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// Overrides settings in LibIO to enable exporting from SvcCore
LibIO
	{
	SetParent(ob, lib_record, lib /*unused*/)
		{
		ob.path = lib_record.path.Split('/').Map({ Object(name: it) })
		ob.parent_name = ob.path.Empty?() ? '' : ob.path.Last().name
		}

	Get_record(lib /*unused*/, record, path /*unused*/)
		{
		if record.GetDefault('type', ' ') is '-'
			record.delete = true
		return record
		}
	}
