// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
OptionalNumberFormat
	{
	Convert(data, mask)
		{
		try
			data = Number(data)
		catch
			data = String(data)
		if (mask isnt false and Number?(data))
			data = data.EuroFormat(mask)
		return data
		}
	}