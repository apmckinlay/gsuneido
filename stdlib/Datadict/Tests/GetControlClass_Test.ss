// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Control value is an Object with string
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (Number)
				Format: (Text)
				}`])
		Assert(GetControlClass.FromField(fieldName) is: NumberControl)

		// Control value is an Object with class
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (NumberControl {})
				Format: (Text)
				}`])
		Assert(GetControlClass.FromField(fieldName).Base() is: NumberControl)

		// Control value is a class - invalid
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: NumberControl {}
				Format: (Text)
				}`])
		Assert(GetControlClass.FromField(fieldName) is: false)

		// Control value has an invalid control class name
		fieldName = .TempName()
		.MakeLibraryRecord([name: "Field_" $ fieldName,
			text: `class
				{
				Control: (ShouldNotExist)
				Format: (Text)
				}`])
		Assert(GetControlClass.FromField(fieldName) is: false)

		// try get control from a non-existent Datadict
		Assert(GetControlClass.FromField('should_not_exist')
			is: FieldControl) // default from Datadict
		}
	}