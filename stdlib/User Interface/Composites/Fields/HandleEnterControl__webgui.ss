// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
FieldControl
	{
	ENTER()
		{
		dirty? = .Dirty?()
		.KillFocus()
		if dirty?
			.Send("NewValue", .Get())
		}
	}