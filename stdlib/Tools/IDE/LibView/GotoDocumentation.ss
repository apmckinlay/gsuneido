// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (name)
	{
	if name is "" or not TableExists?('suneidoc')
		return

	find_items = Object()
	find_items.Add(name)

	// if name contains 'Control' or 'Format' then remove it
	// because Gotofind will add it back on
	if name.Suffix?('Control')
		find_items.Add(name.Replace("Control$", ""))
	if name.Suffix?('Format')
		find_items.Add(name.Replace("Format$", ""))

	address = ''
	for item in find_items
		if (false isnt x = QueryFirst('suneidoc where name is ' $ Display(item) $
			' sort path'))
			{
			address = x.path $ '/' $ x.name
			break
			}

	if address is ''
		return

	GotoUserManual(address)
	}
