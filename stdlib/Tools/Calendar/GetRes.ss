// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ReturnWithDelay()
		{
		Thread.Sleep(1.SecondsInMs()) // slow down brute force attack
		return -1
		}

	Headers(name)
		{
		headers = Object(
			Cache_Control: 'max-age=484200'
			Last_Modified: #20000101
			Expires: Date().Plus(years: 1))

		ext = name.AfterLast('.')
		if MimeTypes.Member?(ext)
			headers.Content_Type = MimeTypes[ext]

		return headers
		}
	}