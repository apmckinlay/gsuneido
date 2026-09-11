// Copyright (C) 2026 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_Select2()
		{
		// test string op with DateOdometerControl
		dd = Field_dateodometer
		op = #(contains, "=~", pre: "(?i)(?q)", suf: "")
		Assert(Select2.Invalid_operator?(op, dd))
		}

	Test_Encode()
		{
		Assert(Field_dateodometer.Encode(#20130401) is: #20130401)
		Assert(Field_dateodometer.Encode("April 1, 2013") is: #20130401)
		Assert(Field_dateodometer.Encode(123) is: 123)
		Assert(Field_dateodometer.Encode(#hello) is: #hello)
		}
	}
