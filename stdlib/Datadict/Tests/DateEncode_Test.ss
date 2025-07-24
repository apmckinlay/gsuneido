// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	old_fmt: false
	Setup()
		{
		.old_fmt = Settings.Get('ShortDateFormat')
		}
	Test_main()
		{
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")
		Assert(Field_date.Encode(#20050930) is: #20050930)
		Assert(Field_date.Encode('#20050930') is: #20050930)
		Assert(Field_date.Encode(false) is: false) // can't encode
		Assert(Field_date.Encode('June 31, 2004') is: 'June 31, 2004') // can't encode
		Assert(Field_date.Encode('June 30, 2004') is: #20040630)
		Assert(Field_date.Encode('20050930') is: #20050930)
		}
	Teardown()
		{
		super.Teardown()
		if .old_fmt isnt false
			Settings.Set('ShortDateFormat', .old_fmt)
		}
	}