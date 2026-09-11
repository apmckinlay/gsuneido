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
			'<span style="white-space:pre;">' $ XmlEntityEncode(parts[0]) $ '</span>' $
			'<span style="flex:1;text-align:center;white-space:pre;">' $
				XmlEntityEncode(parts[1]) $ '</span>' $
			'<span style="text-align:right;white-space:pre;">' $
				XmlEntityEncode(parts[2]) $ '</span>'
		}
	}