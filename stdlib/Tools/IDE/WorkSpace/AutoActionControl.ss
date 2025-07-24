// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
AutoChooseControl
	{
	IgnoreAfterPick?: false
	New(@args)
		{
		super(@args)
		.Send('SetFieldProcOverrideCheck', .allowOverride?)
		}

	allowOverride?()
		{
		return not .ListOpen?()
		}

	InsertChoice(s)
		{
		.Send('SelectAction', action: s)
		}

	Escape()
		{
		if .ListOpen?()
			.Send('CancelList')
		return super.Escape()
		}

	EN_CHANGE()
		{
		.Send("Edit_Change")
		return super.EN_CHANGE()
		}
	
	Choices(prefix, list)
		{
		if prefix isnt ''
			return super.Choices(prefix, list)
		
		listForEmpty = .Send('GetListForEmpty')
		return listForEmpty isnt 0 ? listForEmpty : #()
		}

	// override to avoid auto pick a value if it is only choice
	PickIfOneChoice() {}
	}