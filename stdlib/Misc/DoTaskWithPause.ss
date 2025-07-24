// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (msg/*unused*/, block)
	{
	while (block()) {}
	return true
	}