// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
NumberFormat
	{
	Convert(data, mask /*unused*/ = false )
		{
		try
			data = (data is false) ? 0 : Number(data)
		catch
			return String(data)
		if IsInf?(data)
			data = ''
		return ReadableSize.FromInt(data)
		}
	}
