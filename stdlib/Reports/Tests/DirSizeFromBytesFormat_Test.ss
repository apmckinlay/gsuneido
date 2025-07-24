// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(DirSizeFromBytesFormat.Convert('str') is: 'str')
		Assert(DirSizeFromBytesFormat.Convert('5368709120') is: '5 gb')
		Assert(DirSizeFromBytesFormat.Convert(5368709120) is: '5 gb')
		Assert(DirSizeFromBytesFormat.Convert('') is: '0')
		Assert(DirSizeFromBytesFormat.Convert('5 gb') is: '5 gb')
		}
	}