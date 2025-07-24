// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
CustomizableFieldDialogPropertiesEditor_ChooseListControl
	{
	GetValidMsg(list)
		{
		if list.HasIf?( { it.Has?(',') } )
			return 'Items values can not contain commas'
		super.GetValidMsg(list)
		}
	}