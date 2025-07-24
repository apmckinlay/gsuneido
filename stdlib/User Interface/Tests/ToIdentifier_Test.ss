// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(ToIdentifier('') is: '')
		Assert(ToIdentifier('Abc123') is: 'Abc123')
		Assert(ToIdentifier('Abc_123') is: 'Abc_123')
		Assert(ToIdentifier('\tA\x00 \tb \tc 123 ') is: 'A_b_c_123')
		}
	}