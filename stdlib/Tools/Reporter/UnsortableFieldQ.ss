// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (field, dd = false)
	{
	if dd is false
		dd = Datadict(field)
	if dd.GetDefault('Unsortable', false) is true
		return true

	ctrlName = dd.Control.GetDefault(0, '')
	if Class?(ctrlName)
		ctrlClass = ctrlName
	else
		{
		if not ctrlName.Suffix?('Control')
			ctrlName $= 'Control'
		try
			ctrlClass = Global(ctrlName)
		catch (err /*unused*/, "can't find")
			return false
		}

	// can't use GetDefault on a class
	return ctrlClass.Member?('Unsortable') and ctrlClass.Unsortable
	}
