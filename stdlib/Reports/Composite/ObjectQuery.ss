// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(data)
		{
		.data = data
		.i = 0
		}
	Columns()
		{
		return .data.GetDefault(#columns, .data.Member?(0) ? .data[0].Members() : #())
		}
	Order()
		{
		return .data.GetDefault(#order, #())
		}
	Next()
		{
		return .data.Member?(.i) ? .data[.i++] : false
		}
	}