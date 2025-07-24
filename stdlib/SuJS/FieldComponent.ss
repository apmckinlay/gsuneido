// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
EditComponent
	{
	Name: 'Field'

	New(width = 20, readonly = false,
		font = "", size = "", weight = "", underline = false,
		password = false, justify = "LEFT",
		bgndcolor = "", textcolor = "", tabover = false
		readOnlyBgndColor = false, upper = false, lower = false, maxCharacters = false)
		{
		super(readonly, :bgndcolor, :textcolor, :tabover,
			:font, :size, :weight, :underline,
			:width, height: 1, :readOnlyBgndColor)

		.setStyles(justify, password, upper, lower, maxCharacters)
		}

	setStyles(justify, password, upper, lower, maxCharacters)
		{
		styles = Object()
		styles['text-align'] = justify
		if upper isnt false
			styles['text-transform'] = 'uppercase'
		if lower isnt false
			styles['text-transform'] = 'lowercase'
		.SetStyles(styles)

		if password is true
			.El.type = 'password'

		maxChar = maxCharacters isnt false ? maxCharacters : 512
		.El.SetAttribute(#maxlength, maxChar)
		}

	SetBackColor(color)
		{
		.SetBgndColor(color)
		}

	Set(value)
		{
		super.Set(value)
		// select all text, this is useful for Browse and other controls
		// which create fields on the fly
		.SelectAll()
		}
	}
