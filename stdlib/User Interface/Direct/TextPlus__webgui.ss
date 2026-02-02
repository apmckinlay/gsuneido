// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	valid: true
	New(.text = "", lefttext = false, font = "", size = "",
		weight = "", .readonly = false, tip = "", hidden = false,
		tabover = false, .align = 'center')
		{
		.SuSetHidden(hidden)
		.ComponentArgs = Object(text, lefttext, font, size, weight, readonly, tip, hidden,
			tabover, align)
		}
	value: ""
	Toggle()
		{
		if not .GetReadOnly()
			{
			.value = .value isnt true
			.dirty? = true
			}
		}

	readonly: false
	set_readonly: false
	GetReadOnly()
		{
		return .readonly or .set_readonly
		}
	SetReadOnly(ro)
		{
		.set_readonly = ro
		.Act('SetReadOnly', ro)
		}
	SetEnabled(enabled)
		{
		.Act('SetEnabled', enabled)
		}

	Get()
		{
		return .value is true
		}
	Set(value)
		{
		if value is .value
			return

		.value = value
		.dirty? = false
		.Act('Set', value)
		}

	Valid?()
		{
		.valid = BooleanOrEmpty?(.value)
		return .valid
		}

	GetUnvalidated()
		{
		return .value
		}

	ValidData?(value)
		{
		return BooleanOrEmpty?(value)
		}

	dirty?: false
	Dirty?(dirty = "")
		{
		Assert(dirty is true or dirty is false or dirty is "")
		if (dirty isnt "")
			.dirty? = dirty
		return .dirty?
		}

	SetColor(color)
		{
		.Act('SetColor', TranslateColor(color))
		}

	GetText()
		{
		return .text
		}

	// for HighlightControl
	FindControl(name)
		{
		return name is 'Static' ? this : false
		}
	}
