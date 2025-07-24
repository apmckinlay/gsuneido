// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: "Encrypt"
	Get()
		{
		return super.Get().Xor(EncryptControlKey())
		}
	Set(value)
		{
		super.Set(value.Xor(EncryptControlKey()))
		}
	ZoomReadonly(value)
		{
		ZoomControl(0, Decrypt(value), readonly:)
		}
	}