// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_DisplayValues()
		{
		fn = ParamsChooseListControl.DisplayValues

		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (FirstLastName) }`])

		Assert(fn(#(), name) is: #())
		Assert(fn(#(a, b), name) is: #(a, b))

		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (FirstLastNameControl) }`])
		Assert(fn(#(), name) is: #())
		Assert(fn(#(a, b), name) is: #(a, b))

		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (ChooseDate) }`])
		Assert(fn(#(), name) is: #())
		today = Date().NoTime()
		Assert(fn(Object(today), name) is: Object(DateControl.FormatValue(today)))
		}
	}