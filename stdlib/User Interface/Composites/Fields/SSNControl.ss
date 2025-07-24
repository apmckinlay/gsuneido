// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
SSNSINControl
	{
	Name: 'SSN'
	Pattern: '^[0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]$'
	Mask: '###-##-####'

	Valid?()
		{
		value = .Get().Xor(EncryptControlKey())
		result = super.Valid?()
		return result and .validValue?(value)
		}
	validValue?(value)
		{
		return (value =~ .Pattern or value is '')
		}
	ValidData?(@args)
		{
		value = args[0].Xor(EncryptControlKey())
		return .validValue?(value)
		}
	}
