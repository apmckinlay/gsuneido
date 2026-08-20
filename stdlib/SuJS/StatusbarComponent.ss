// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
StatusComponent
	{
	Name: 'Statusbar'

	New()
		{
		super()
		.El.SetStyle('display', 'flex')
		}

	Set(text)
		{
		parts = text.Split('\t').Set_default('')
		.El.innerHTML =
			'<span>' $ XmlEntityEncode(parts[0]) $ '</span>' $
			'<span style="flex:1;text-align:center">' $
				XmlEntityEncode(parts[1]) $ '</span>' $
			'<span style="text-align:right; margin-right: .5em;">' $
				XmlEntityEncode(parts[2]) $ '</span>'
		}
	}