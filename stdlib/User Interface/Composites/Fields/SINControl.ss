// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
SSNSINControl
	{
	Name: 'SIN'
	Pattern: '^[0-9][0-9][0-9]\s[0-9][0-9][0-9]\s[0-9][0-9][0-9]$'
	Mask: '### ### ###'

	Valid?()
		{
		value = .Get().Xor(EncryptControlKey())
		result = super.Valid?()
		return result and .validValue?(value)
		}
	validValue?(value)
		{
		return SINValid?(value, false)
		}
	ValidData?(@args)
		{
		value = args[0].Xor(EncryptControlKey())
		return .validValue?(value)
		}
	}
