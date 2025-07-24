// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'Status'
	Xstretch: 1
	New()
		{
		.CreateElement('div')
		.El.SetStyle('min-height', '1.5em')
		.El.SetStyle('line-height', '1.5em')
		.SetValid()
		}

	Set(text)
		{
		.El.innerHTML = XmlEntityEncode(text)
		}

	SetValid(valid = true)
		{
		.SetBkColor(valid ? CLR.ButtonFace : CLR.ErrorColor)
		}

	SetWarning(warn = true)
		{
		.SetBkColor(warn ? CLR.WarnColor : CLR.ButtonFace)
		}

	SetBkColor(color)
		{
		if color is false
			color = CLR.ButtonFace
		.El.SetStyle('background-color', ToCssColor(color))
		}
	}
