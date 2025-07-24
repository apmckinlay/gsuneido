// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Memoize
	{
	Func(table)
		{
		if false is x = Query1Cached(table $
			" where path = '/res' and name = 'HtmlPrefix'")
			return false
		prefix = Asup(x.text)
		for contrib in GetContributions('HtmlWrap_' $ table $ 'Prefix')
			{
			condition = contrib.GetDefault('condition', function() { return true })
			attribs = contrib.GetDefault('attribs', '')
			if condition()
				prefix = prefix.Replace('(?q)' $ contrib.cssClass $ ' { display: none; }',
					'\=' $ contrib.cssClass $ ' { ' $ attribs $ ' }')
			}
		return prefix
		}
	}