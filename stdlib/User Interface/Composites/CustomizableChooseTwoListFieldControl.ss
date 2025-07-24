// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// used by CustomizableChooseFieldsControl
ChooseTwoListFieldControl
	{
	Name: 'CustomizableChooseTwoListField'
	KillFocus()
		{
		dirty? = .Dirty?()
		.Parent.ProcessFieldValues()
		if dirty?
			.Send("NewValue", .Get())
		}

	Valid?()
		{
		return .Parent.Valid?()
		}
	}