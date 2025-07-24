// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_basename()
		{
		basename = FindReferences.FindReferences_basename
		Assert(basename("Field_testing") is: "testing")
		Assert(basename("Rule_testing") is: "testing")
		Assert(basename("Trigger_testing") is: "testing")
		Assert(basename("testingControl") is: "testing")
		Assert(basename("testingFormat") is: "testing")
		Assert(basename("testingComponent") is: "testing")
		Assert(basename("Control") is: "Control")
		Assert(basename("Format") is: "Format")
		Assert(basename("Component") is: "Component")
		}

	Test_book_references()
		{
		book = .MakeBook()
		recName = .TempName()
		.MakeBookRecord(book, 'Text ' $ recName $ ' More')
		FindReferences.FindReferences_book_references(recName, list = Object())
		Assert(list isSize: 1)

		recName = .TempName()
		rec = .MakeBookRecord(book, 'Text ' $ recName $ ' More')
		QueryDo('update ' $ book $
			' where num is ' $ Display(rec.num) $ ' set name = "testing.jpg"')
		FindReferences.FindReferences_book_references(recName, list = Object())
		Assert(list isSize: 0)
		}

	Test_AllOccurrences()
		{
		.MakeLibraryRecord([name: "Record_For_FindReferences_Test",
			text: `function() {
				Abc(123)
				Bcd
				Abc(efg)
				Efg
				AbcControl(321)
				}`])
		LibraryTables.ResetCache()

		all = FindReferences.AllOccurrences(.TestLibName(),
			'Record_For_FindReferences_Test', 'Abc', 'AbcControl')
		Assert(all.Tr(' ', '')
			is: '2:Abc(123)\r\n' $ '4:Abc(efg)\r\n' $ '6:AbcControl(321)')

		all = FindReferences.AllOccurrences(.TestLibName(),
			'Record_For_FindReferences_Test', 'Efg', 'Efg')
		Assert(all.Tr(' ', '') is: '5:Efg')

		all = FindReferences.AllOccurrences(.TestLibName(),
			'Record_For_FindReferences_Test', 'AAA', 'AAA')
		Assert(all is: '')

		all = FindReferences.AllOccurrences('test',
			'Record_For_FindReferences_Test', 'AAA', 'AAA')
		Assert(all is: '')

		all = FindReferences.AllOccurrences('stdlib',
			'Record_For_FindReferences_Test', 'AAA', 'AAA')
		Assert(all is: '')
		}

	Test_references()
		{
		fn = FindReferences.FindReferences_references

		notFoundClass = '//Suneido ...
			class
				{
				}'
		foundClass = '//Suneido ...
			class
				{
				Fn() { Fake1() }
				}'
		unusedLib = .MakeLibrary([name: 'RefTest1', text: '#(UN USED library)'],
			[name: 'RefTest4_Test', text: foundClass],
			[name: 'RefTest3', text: foundClass])
		otherLib = .MakeLibrary([name: 'RefTest2', text: '#(not a real class)'])

		recsToBuild = Object()
		recsToBuild.Add([name: 'Fake1', text: notFoundClass])
		recsToBuild.Add([name: 'NotARealRecord', text: notFoundClass])
		recsToBuild.Add([name: 'AFakeRecord', text: foundClass])
		recsToBuild.Add([name: 'OneMoreFakeRecord', text: foundClass])
		recsToBuild.Add([name: 'Fake1_Test', text: foundClass])
		.MakeLibraryRecord(@recsToBuild)

		basename = FindReferences.FindReferences_basename('Fake1')
		list = Object()
		fn(Object('Test_lib', unusedLib), #('Test_lib'), 'Fake1', basename, list, false)
		Assert(list.Size() is: 5)
		Assert(list[0].Location is: 'AFakeRecord')
		Assert(list[1].Location is: 'Fake1_Test')
		Assert(list[2].Location is: 'OneMoreFakeRecord')
		Assert(list[3].Location is: 'RefTest3')
		Assert(list[4].Location is: 'RefTest4_Test')

		Assert(list[0].Table is: 'Test_lib')
		Assert(list[3].Table is: "(" $ unusedLib $ ")")

		list = Object()
		fn(Object('Test_lib', unusedLib), #('Test_lib'), 'Fake1', basename, list, true)
		Assert(list.Size() is: 3)
		Assert(list[0].Location is: 'AFakeRecord')
		Assert(list[1].Location is: 'OneMoreFakeRecord')
		Assert(list[2].Location is: 'RefTest3')

		list = Object()
		fn(Object(otherLib), Object(otherLib), 'Fake1', basename, list, false)
		Assert(list.Size() is: 0)
		}
	}
