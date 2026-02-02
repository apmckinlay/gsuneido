// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Type(table)
		{
		return GetContributions('BookContentType').GetDefault(table, #html)
		}

	Match(table, text)
		{
		switch .Type(table)
			{
		case #html:
			return text.Prefix?("<")
		case #md:
			return text.Prefix?("<") or text[::1] is '#' and text[1::1] isnt '('
			}
		}

	ToHtml(table, text)
		{
		switch .Type(table)
			{
		case #html:
			return text
		case #md:
			return MarkdownToHtml(text,
				addons: [Md_Addon_Table, Md_Addon_Definition, Md_Addon_Book,
					Md_Addon_suneido_style])
			}
		}
	}