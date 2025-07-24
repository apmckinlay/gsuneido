// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(md)
		{
		lines = Lines.Iter(md)
		result = ""
		closing = ""
		while lines isnt line = lines.Next()
			{
			if line is "```"
				{
				result = Opt(result, '\n') $ "<pre>"
				while ((line = lines.Next()) not in ("```", lines))
					result $= line $ '\n'
				result $= "</pre>"
				if line is lines
					break
				}
			else if line.Blank?()
				{
				if closing isnt ""
					{
					result $= closing
					closing = ""
					}
				}
			else if closing is ""
				{
				line = "<p>" $ line
				closing = "</p>"
				}
			result = Opt(result, '\n') $ line
			}
		if closing isnt ""
			result $= closing
		if result isnt ""
			result $= '\n'
		return result
		}
	}