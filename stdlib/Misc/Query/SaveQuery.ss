// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (query, table)
	// post: result of query saved as table
	{
	WithQuery(query)
		{|q|
		cols = q.Columns()
		keys = q.Keys()
		}
	create = "create " $ table $ " (" $ cols.Join(',') $ ")"
	for key in keys
		create $= " key(" $ key $ ")"
	Database(create)
	QueryDo("insert (" $ QueryStripSort(query) $ ") into " $ table)
	}