// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: "Encrypt"
	Get()
		{
		return Decrypt(super.Get())
		}
	Set(value)
		{
		super.Set(Decrypt(value))
		}
	ZoomReadonly(value)
		{
		ZoomControl(0, Decrypt(value), readonly:)
		}
	}