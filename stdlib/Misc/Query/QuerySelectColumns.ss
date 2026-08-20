// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query)
		{
		ob = QueryColumns(QueryStripSort(query))
		removeOb = .nonPermissableFields(query)
		return ob.Difference(removeOb).RemoveIf(Customizable.DeletedField?)
		}

	// extracted for test
	nonPermissableFields(query)
		{
		return Customizable.GetNonPermissableFields(query)
		}
	}