// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	validate: false
	SetValidate(.validate?)
		{
		}
	Valid?()
		{
		if .validate
			try
				'' =~ .Get()
			catch (unused, '*regex')
				return false
		return super.Valid?()
		}
	}