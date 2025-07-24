// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (from)
	{
	field = 'trlang_' $ GetLanguage().name
	try
		return Query1('translatelanguage
			where trlang_from is ' $ Display(from))[field]
	return ""
	}