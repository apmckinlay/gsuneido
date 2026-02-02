// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
ButtonControl
	{
	ComponentName: 'EnhancedButton'
	New(.text = false, command = false,
		tabover = false, .defaultButton = false, style/*unused*/ = 0,
		tip = false, pad = false, font = "", size = "", weight = "", textColor = false,
		.width = false, .buttonWidth = false, .buttonHeight = false,
		italic = false, underline = false, strikeout = false,
		.image = false, .mouseOverImage = false, .mouseDownImage = false,
		.imageColor = false, .mouseOverImageColor = false,
		.book = 'imagebook', .mouseEffect = false, .imagePadding = 0,
		.buttonStyle = false, classic/*unused*/ = false, .noBgnd = false,
		.enlargeOnHover = false, hidden = false)
		{
		super(.text is false ? 'EnhancedButton' : .text, command, :defaultButton, :tip,
			:pad, :font, :size, :weight, color: .textColor = TranslateColor(textColor),
			:italic, :underline, :strikeout, :hidden)
		.buttonName = command

		.ComponentArgs = Object(text, tabover, defaultButton, tip, pad,
			font, size, weight, .textColor, width, buttonWidth, buttonHeight,
			italic, underline, strikeout,
			.toCharCode(image), .toCharCode(.mouseOverImage),
			.toCharCode(.mouseDownImage), imageColor, mouseOverImageColor, book,
			mouseEffect, imagePadding, buttonStyle, enlargeOnHover)
		}

	toCharCode(img)
		{
		if String?(img) and false isnt char = IconFont().MapToCharCode(img)
			return char
		if Object?(img)
			return Object(IconFont().MapToCharCode(img[0]),
				IconFont().MapToCharCode(img[1]),
				img.GetDefault(#highlighted, 0), img.GetDefault(#gap, 1))
		if String?(img) and img.Suffix?('.emf')
			{
			name = img.RemoveSuffix('emf') $ 'svg'
			if false isnt rec = QueryFirst(.book $ ' where name is ' $ Display(name) $
				' sort num')
				img = rec.text[rec.text.Find('<svg')..]
			}
		return img
		}

	GetImage()
		{
		return .image
		}
	SetImage(image, mouseOverImage = false, mouseDownImage = false)
		{
		.image = image
		if mouseOverImage isnt false
			.mouseOverImage = mouseOverImage
		if mouseDownImage isnt false
			.mouseDownImage = mouseDownImage
		.Act('SetImage', .toCharCode(image),
			.toCharCode(mouseOverImage),
			.toCharCode(mouseDownImage))
		}

	SetImageColor(.imageColor = false, mouseOverImageColor = false)
		{
		if mouseOverImageColor isnt false
			.mouseOverImageColor = mouseOverImageColor
		.Act('SetImageColor', imageColor, mouseOverImageColor)
		}

	GetImageColor()
		{
		return .imageColor
		}

	Set(text)
		{
		.SetText(text)
		}
	SetText(.text)
		{
		.Act('SetText', .text)
		}

	pushed: false
	Pushed?(state = -1)
		{
		if state isnt -1 and state isnt .pushed
			{
			.pushed = state
			.Act('Pushed?', state)
			}
		return .pushed
		}

	SetTextColor(color)
		{
		.Act(#SetTextColor, .textColor = TranslateColor(color))
		}

	GetTextColor()
		{
		return .textColor
		}

	gray: 0xaaaaaa
	Grayed(state = -1)
		{
		if state isnt -1
			.SetTextColor(state is true ? .gray : false)
		return .textColor is .gray
		}

	SetMouseEffect(.mouseEffect)
		{
		.Act(#SetMouseEffect, mouseEffect)
		}

	prevFocusHwnd: false
	SyncPrevFocus(.prevFocusHwnd)
		{
		}
	BeforeCallCommand()
		{
		if .prevFocusHwnd isnt false and .buttonStyle is false
			SetFocus(.prevFocusHwnd)
		}

	ContextMenu(x, y)
		{
		if 0 isnt .Send('EnhancedButton_ContextMenu', :x, :y)
			return 0
		return super.ContextMenu(x, y)
		}

	RBUTTONDOWN(r)
		{
		.Send('EnhancedRButtonDown', r)
		}

	GetButtonName()
		{
		return .buttonName
		}
	}
