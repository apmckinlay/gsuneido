// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_FontControl
	{
	On_Customize_Font()
		{
		if false isnt newFont = SuFontControl(.Font.Copy())
			.ChangeFont(newFont)
		}
	}