// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Test behaviour of built-in numbers. The goal is to establish a predictable
// set of rules that the built-in number type must meet on GSuneido.
Test
	{
	int_min:         0x80000000   // min. value storeable in a 32-bit integer
	Test_RightShift()
		{
		Assert(.int_min >> 1, is: 0x40000000)
		}
	}