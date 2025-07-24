// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
EncryptFormat
	{
	Convert(data)
		{
		if not String?(data)
			data = String(data)

		if data.Prefix?('SIN:')
			data = data.AfterFirst('SIN: ')
		else if data.Prefix?('SSN:')
			data = data.AfterFirst('SSN: ')

		return super.Convert(data)
		}
	}
