// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (query)
	{
	QueryColumns(query).RemoveIf(Customizable.DeletedField?)
	}