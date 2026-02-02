// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	ConvertToStr(data)
		{
		suffix = ''
		if data.Suffix?('...')
			{
			data = data.BeforeLast(',')
			suffix = ',...'
			}
		return data.Split(',').Map!({ Date(it).ShortDate() }).Join(',') $ suffix
		}
	}
