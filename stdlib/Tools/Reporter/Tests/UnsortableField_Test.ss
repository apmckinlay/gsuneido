// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(UnsortableField?("attachment")) // image
		Assert(UnsortableField?("text")) // editor
		Assert(UnsortableField?("name") is: false) // field
		Assert(UnsortableField?("name", Datadict("name")) is: false)
		Assert(UnsortableField?("scintilla_rich")) //scintilla rich
		Assert(UnsortableField?("scintilla_rich", Datadict("scintilla_rich")))
		}

	Test_ControlSpecIsClass()
		{
		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (class { Name: "TestName" }) }`])
		Assert(UnsortableField?(name) is: false)

		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (class { Name: "TestName"; Unsortable: false }) }`])
		Assert(UnsortableField?(name) is: false)

		name = .TempName().Lower()
		.MakeLibraryRecord([name: "Field_" $ name,
			text: `class { Control: (class { Name: "TestName"; Unsortable: true }) }`])
		Assert(UnsortableField?(name))
		}
	}