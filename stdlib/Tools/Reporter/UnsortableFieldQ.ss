// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (field, dd = false)
	{
	if dd is false
		dd = Datadict(field)
	if dd.GetDefault('Unsortable', false) is true
		return true

	if false is ctrlClass = GetControlClass.FromControl(dd.Control)
		return false
	// can't use GetDefault on a class
	return ctrlClass.Member?('Unsortable') and ctrlClass.Unsortable
	}
