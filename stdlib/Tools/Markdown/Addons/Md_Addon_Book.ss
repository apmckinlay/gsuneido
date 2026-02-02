// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Addon_Base
	{
	Match(line, start)
		{
		line = line[start..].Trim()
		if not line.Prefix?('<$') or not line.Suffix?('$>')
			return false

		return new Md_Asup(line)
		}

	ConvertToHtml(writer, item)
		{
		if not item.Base?(Md_Asup)
			return false

		writer.AddPure(item.Content)
		return true
		}

	PreInline(item)
		{
		if item[0] is #link and item.href.Prefix?('http')
			item.href = `suneido:/eval?ShellExecute(0, 'open', '` $ item.href $ `')`
		return false
		}
	}