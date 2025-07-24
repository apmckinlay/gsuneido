// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
HandleEnterComponent
	{
	Name: 'Number'
	DefaultWidth: 15
	New(.mask = '-###,###,###', readonly = false, width = false, justify = "RIGHT",
		font = "", size = "", weight = "", underline = false, tabover = false)
		{
		super(:mask, :width, :readonly, :justify,
			:font, :size, :weight, :underline, :tabover)
		}

	SetFontAndSize(font, size, weight, underline, width, height/*unused*/)
		{ // overrides EditControl
		// NOTE: using '9' assumes digits are all the same width (normal)
		// NOTE: ignores width if mask is set
		if .mask isnt false
			super.SetFontAndSize(font, size, weight, underline, 1, 1,
				text: .mask.Tr('#', '9'))
		else
			super.SetFontAndSize(font, size, weight, underline, width, 1, text: '9')
		}

	SetVScroll(.increase, .rollover, .rangefrom, .rangeto)
		{
		.El.AddEventListener('wheel', .scroll)
		}

	scroll(event)
		{
		if .GetReadOnly() or event.deltaY is 0
			return
		// Integer? will throw when .El.value is not a number string and Number() returns NaN
		try
			{
			if ((false is val = Number(.El.value)) or not Integer?(val))
				return
			}
		catch
			return

		val += event.deltaY < 0 ? .increase : -.increase
		val = .handleBoundaries(val)
		.SelectAll()
		.Set(val)
		.EN_CHANGE()
		}

	handleBoundaries(val)
		{
		if .rollover
			{
			if val > .rangeto
				val = .rangefrom
			if val < .rangefrom
				val = .rangeto
			}
		else
			{
			if val > .rangeto
				val = .rangeto
			if val < .rangefrom
				val = .rangefrom
			}
		return val
		}
	}
