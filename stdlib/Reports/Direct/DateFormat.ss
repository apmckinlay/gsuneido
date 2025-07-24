// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// abstract base class for ShortDateFormat, LongDateFormat, and DateTimeFormat
// derived classes must define Method: (e.g. ShortDate)
TextFormat
	{
	New(data = false, justify = "left", font = false)
		{
		super(data, width: .Width(), :justify, :font)
		}
	WidthChar: '9' // not quite right if format has e.g. MMM (month name)
	Width()
		{
		// date-time is chosen to be max width
		// + 1 is to allow for non-digit characters e.g. May or AM
		// can't just use length of format because e.g. "M" may expand to "12"
		// need two digit month for digits, September for longest text
		return Max(#19991111.191919[.Method]().Size(),
			#19990911.191919[.Method]().Size()) + 1
		}

	GetDefaultWidth()
		{
		return .Width()
		}

	ConvertToStr(data)
		{
		return Date?(data) ? data[.Method]() : data
		}
	}
