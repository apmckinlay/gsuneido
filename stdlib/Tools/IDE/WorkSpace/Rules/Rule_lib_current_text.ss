// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function()
	{
	return this.GetDefault(#lib_invalid_text, '') is ''
		? .text
		: .lib_invalid_text
	}
