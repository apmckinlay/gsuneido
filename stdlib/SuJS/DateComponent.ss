// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
FieldComponent
	{
	Name: "Date"
	New(.readonly = false, .showTime = false, tabover = false, hidden = false)
		{
		super(:readonly, :tabover, :hidden)
		}

	SetFontAndSize(font, size, weight, underline, width/*unused*/, height/*unused*/ = 0)
		{ // overrides EditControl
		// widen year to match what FormatValue does
		fmt = Settings.Get('ShortDateFormat').Replace('\<yy?y?\>', 'yyyy')
		if .showTime is true
			fmt $= ' ' $ Settings.Get('TimeFormat')
		// this date/time must have 2 digit year, month, day, hour, minute, second
		// can't just use format because e.g. "M" may expand to "12"
		s = #19991231.235959.Format(fmt)
		super.SetFontAndSize(font, size, weight, underline, 1, 1, text: s)
		}
	}
