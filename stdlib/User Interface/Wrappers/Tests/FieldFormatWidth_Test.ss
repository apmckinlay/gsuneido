// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		// test image
		Assert(FieldFormatWidth('image', 10) is: 100)

		// test field width member
		Assert(FieldFormatWidth('city', 10) is: 80)
		Assert(FieldFormatWidth('name', 10) is: 160)

		// test fmt.Width()
		dateStrSize = #20201119.ShortDate().Size() + 1 // use double digit month
		Assert(FieldFormatWidth('date', 10) is: dateStrSize * 10)
		}
	}