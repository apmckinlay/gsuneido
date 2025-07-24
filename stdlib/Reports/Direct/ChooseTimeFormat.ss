// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
TextFormat
	{
	New(.data = false, font = false, justify = 'right')
		{
		super(false, width: .findWidth(), :justify, :font)
		}

	findWidth()
		{
		format = Settings.Get('TimeFormat')
		reqWidth = format.Size()
		if format !~ '(?i)hh'
			++reqWidth
		return reqWidth
		}

	WidthChar: '9'
	Print(x, y, w, h, data = "")
		{
		if .data isnt false
			data = .data
		super.Print(x, y, w, h, .formatData(data))
		}

	formatData(data)
		{
		if data is ''
			return ''
		data = String(data)
		if not data.Numeric?()
			return data
		timeStr = data.LeftFill(4, '0') /*=hhmm*/
		// set second and millisecond to zero for safety
		date = Date(hour: Number(timeStr[..2]), minute: Number(timeStr[-2..]),
			second: 0, millisecond: 0)
		return date.Time()
		}

	DataToString(value, record /*unused*/)
		{
		return .formatData(value)
		}
	}