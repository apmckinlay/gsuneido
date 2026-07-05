// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
ZoomControl
	{
	New(.text, readonly = false, font = '', size = '', textLimit = false)
		{
		// Need to Encrypt the text in preparation for: Encrypted Control > Set
		super(.text.Xor(EncryptControlKey()), readonly, font, size, textLimit)
		}

	OK()
		{
		// Need to Encrypt the text in preparation for: Encrypted Control > Set
		return super.OK().Xor(EncryptControlKey())
		}

	Cancel()
		{
		return .text // Return the original text as .Set is not called on cancel
		}
	}
