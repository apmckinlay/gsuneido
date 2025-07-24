// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (query)
	{
	ob = QueryColumns(QueryStripSort(query))
	removeOb = Customizable.GetNonPermissableFields(query)
	return ob.Difference(removeOb).RemoveIf(Customizable.DeletedField?)
	}
