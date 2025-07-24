// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		// Can't use .MakeTable because we need the table name to end with _table
		.custFieldName = .TempName()
		.custTableName = .custFieldName $ '_table'
		Database("ensure " $ .custTableName $ " (name, desc) key(name)")
		QueryOutput(.custTableName, Record(name: 'My Custom Field',
			desc: 'My Custom Field'))
		}
	Test_ValidData()
		{
		args = Object("My Custom Field", allowOther: false, mandatory: false,
			customField: .custFieldName)

		Assert(CustomKeyControl.ValidData?(@args))

		args = Object("My Invalid Field", allowOther: false, mandatory: false,
			customField: .custFieldName)
		Assert(CustomKeyControl.ValidData?(@args) is: false)
		}

	Teardown()
		{
		Database("destroy " $ .custTableName)
		super.Teardown()
		}
	}