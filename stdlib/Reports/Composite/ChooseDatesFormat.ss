// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	ConvertToStr(data)
		{
		return data.Split(',').Map!({ Date(it).ShortDate() }).Join(',')
		}
	}
