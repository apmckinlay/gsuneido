// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	bookOption = Suneido.GetDefault('CurrentBookOption', '').Extract('(/.*?/)')
	return bookOption is false
		? ''
		: bookOption.Tr('/')
	}