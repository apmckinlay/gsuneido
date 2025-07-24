// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(t /*unused*/, oldrec, newrec)
		{
		key = oldrec isnt false
			? oldrec.custfield_name
			: newrec isnt false
				? newrec.custfield_name
				: ''

		if key isnt ''
			CustomizableOnServer.NotifyCustomizableFieldsChanges(key)
		}
	}
