// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Field: false
	CallClass(data)
		{
		if .Field is false
			throw 'Must Inheret'
		return FormatValue.FormatDataToString(Datadict(.Field), data)
		}
	}