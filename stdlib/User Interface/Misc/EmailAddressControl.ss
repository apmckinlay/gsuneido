// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: "EmailAddress"
	Valid?()
		{
		addr = .Get()
		return super.Valid?() and (addr is "" or ValidEmailAddress?(addr))
		}
	}