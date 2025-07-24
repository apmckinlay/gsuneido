// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
function (value, key = false)
	{
	if key is false
		key = EncryptControlKey()
	return String?(value) ? value.Xor(key) : value
	}