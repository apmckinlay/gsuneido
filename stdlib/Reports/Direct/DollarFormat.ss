// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
OptionalNumberFormat
	{
	Convert(data, mask = false)
		{
		try
			data = Number(data)
		catch
			data = String(data)
		if IsInf?(data)
			data = ''
		if mask isnt false and Number?(data)
			data = data.DollarFormat(mask)
		return data
		}
	}

