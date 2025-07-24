// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: 'SSNSIN'

	New(width = 11, .mandatory = false)
		{
		super(:width)
		}

	Get()
		{
		return super.Get().Xor(EncryptControlKey())
		}

	Set(value)
		{
		super.Set(value.Xor(EncryptControlKey()))
		}

	KillFocus()
		{
		s = .Get().Xor(EncryptControlKey())

		if s is '' or s =~ .Pattern or s =~ '[a-z]|[A-Z]'
			return

		dirty? = .Dirty?()
		s = s.Tr('^0-9')
		.Set(.Match(s, .Mask).Xor(EncryptControlKey()))
		.Dirty?(dirty?)
		return
		}

	Match(s, pattern)
		{
		t = ''
		si = 0
		for pc in pattern
			{
			sc = s[si++]
			if pc is '#' or pc is sc
				t $= sc
			else // fill in literal char
				{
				t $= pc
				--si // undo increment
				}
			}
		return t
		}

	Valid?()
		{
		value = .Get().Xor(EncryptControlKey())
		if .mandatory is false and value is ''
			return true
		return super.Valid?()
		}

	// used by ParamsSelect in list
	DisplayValues(control /*unused*/, vals)
		{
		return vals.Map({ it.Xor(EncryptControlKey()) })
		}
	}