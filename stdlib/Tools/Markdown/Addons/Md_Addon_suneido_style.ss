// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Addon_Base
	{
	ConvertToHtml(writer, item)
		{
		if not item.Base?(Md_Code) or item.Info isnt 'suneido'
			return false

		code = Addon_suneido_style.BuildWithStyles(item.Codes.Join('\n'),
			Addon_suneido_style.DefaultStyles())

		writer.Add('code', code)
		return true
		}
	}