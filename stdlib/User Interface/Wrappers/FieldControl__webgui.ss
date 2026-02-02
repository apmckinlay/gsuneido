// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
EditControl
	{
	Name: 'Field'
	ComponentName: 'Field'
	MaxCharacters: 512
	New(width = 20, status = "", readonly = false,
		font = "", size = "", weight = "", underline = false,
		password = false, justify = "LEFT", style /*unused*/ = 0,
		set = false, mandatory = false, .trim = true,
		bgndcolor = "", textcolor = "", hidden = false, tabover = false,
		cue = false, readOnlyBgndColor = false, acceptDrop /*unused*/= false,
		.upper = false, .lower = false, maxCharacters = false)
		{
		super(mandatory, readonly, :bgndcolor, :textcolor, :hidden, :tabover, :width,
			height: 1, :cue, :font, :size, :weight, :underline, :readOnlyBgndColor,
			:status)
		.MaxCharacters = maxCharacters is false
			? .MaxCharacters
			: maxCharacters
		if set isnt false
			{
			.Set(set)
			.Send("NewValue", .Get())
			}

		.ComponentArgs = Object(width, readonly, font, size, weight,
			underline, password, justify, bgndcolor, textcolor, tabover,
			readOnlyBgndColor, upper, lower, .MaxCharacters)
		}

	EN_KILLFOCUS()
		{
		dirty? = .Dirty?()
		// has to call super here
		super.EN_KILLFOCUS()
		if dirty?
			.Send("NewValue", .Get())
		return 0
		}

	Get()
		{
		text = super.Get()
		if .upper isnt false
			text = text.Upper()
		if .lower isnt false
			text = text.Lower()
		return .trim ? text.Trim() : text
		}
	Set(value)
		{
		if not String?(value)
			value = Display(value)
		.Dirty?(false)
		super.Set(value)
		}

	SetBackColor(color)
		{
		.Act('SetBackColor', color)
		}

	changed?: false // killfocus sets dirty back to false
	Valid?()
		{
		if super.Valid?() is false
			return false

		return .ValidLength?(.Get())
		}

	ValidLength?(data)
		{
		if ((.Dirty?() or .changed?) and
			.ValidTextLength?(data, .MaxCharacters) is false)
			{
			if not .changed? and // only show alert once
				not .Window.Base?(ListEditWindow) 	// show alert in ListEditWindow causes
													// infinite alert and call overflow
				.AlertInfo("Invalid Text", 'data exceeds the ' $ .MaxCharacters $
					' character limit for this field.')
			.changed? = true
			return false
			}
		return true
		}

	ValidTextLength?(data, maxCharacters)
		{
		if not String?(data)
			return true
		return data.Size() <= maxCharacters
		}

	ValidData?(@args)
		{
		maxCharacters = args.GetDefault('maxCharacters', .MaxCharacters)
		return super.ValidData?(@args) and .ValidTextLength?(args[0], maxCharacters)
		}
	}