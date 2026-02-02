// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	IdleAfterChange()
		{
		display = .Send('GetDisplay')
		if display in (false, 0)
			return

		md = .Get()
		html = MarkdownToHtml(md, addons: [Md_Addon_Table, Md_Addon_Definition])
		display.Set(html)
		}
	}