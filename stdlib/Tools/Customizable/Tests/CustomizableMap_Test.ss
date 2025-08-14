// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		ServerSuneido.DeleteMember('Memoize_CustomizableMap.CustomizableMap_cached')
		}

	Test_CopyCustomFields()
		{
		table1 = .MakeTable('(key, field) key(key)')
		fld1 = .MakeCustomField(table1, 'Text, single line', 'TEST FIELD1').field
		fld2 = .MakeCustomField(table1, 'Text, single line', 'TEST FIELD2').field

		table2 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table2,
			text: `class { Name: ` $ table2 $ ` }`])
		fld3 = .MakeCustomField(table2, 'Text, single line', 'TEST FIELD1').field
		.MakeCustomField(table2, 'Text, single line', 'TEST FIELDNOTMAPPED').field

		map = CustomizableMap(table1, table2)
		Assert(map.CustomizableMap_map
			is: [[from_field: fld1, to_field: fld3, trim: false]])

		fromRec = [a: 1, b: 2]
		fromRec[fld1] = 1000
		fromRec[fld2] = 2000

		toRec = [a: 1, c: 3]
		map.CopyCustomFields(toRec, fromRec)
		Assert(toRec[fld3] is: 1000)

		map.CopyCustomFields(fromRec, fromRec)
		Assert(fromRec[fld3] is: 1000)

		// test query
		map = CustomizableMap(table1 $ ' rename key to key2', table2 $ ' sort field')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld1, to_field: fld3, trim: false]])

		fromRec = [a: 1, b: 2]
		fromRec[fld1] = 1000
		fromRec[fld2] = 2000

		toRec = [a: 1, c: 3]
		map.CopyCustomFields(toRec, fromRec)
		Assert(toRec[fld3] is: 1000)

		map.CopyCustomFields(fromRec, fromRec)
		Assert(fromRec[fld3] is: 1000)

		// testing reset cache after deleting custom field
		c1 = Customizable(table1)
		c1.DeleteField(fld1)
		fromRec = []
		fromRec[fld1] = 1000
		fromRec[fld2] = 2000
		CustomizableMap(table1, table2).CopyCustomFields(toRec = [], fromRec)
		Assert(toRec[fld3] is: '')
		}

	Test_different_types()
		{
		table1 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table1,
			text: `class { Name: ` $ table1 $ ` }`])
		fld5 = .MakeCustomField(table1, 'Text, single line', 'TEST FIELD5').field
		.MakeCustomField(table1, 'Text, single line', 'TEST FIELDUNMATCHEDTYPE').field
		.MakeCustomField(table1, 'Text, single line', 'TEST FIELDUNMATCHEDTYPE2').field

		table2 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table2,
			text: `class { Name: ` $ table2 $ ` }`])
		fld7 = .MakeCustomField(table2, 'Text, single line', 'TEST FIELD5').field
		.MakeCustomField(table2, 'Dollar', 'TEST FIELDUNMATCHEDTYPE')
		.MakeCustomField(table2, 'Info', 'TEST FIELDUNMATCHEDTYPE2')

		map = CustomizableMap(table1, table2)
		// shouldn't have TEST FIELD6, because the type doesn't match
		Assert(map.CustomizableMap_map
			is: [[from_field: fld5, to_field: fld7, trim: false]])
		}

	Test_allowToFill?()
		{
		table1 = .MakeTable('(key, field) key(key)')
		fld9 = .MakeCustomField(table1, 'Text, single line', 'TEST FIELD5').field
		.MakeCustomField(table1, 'Attachment', 'TEST FIELDTYPENOTALLOWED').field

		table2 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table2,
			text: `class { Name: ` $ table2 $ ` }`])
		fld10 = .MakeCustomField(table2, 'Text, single line', 'TEST FIELD5').field
		.MakeCustomField(table2, 'Attachment', 'TEST FIELDTYPENOTALLOWED').field

		map = CustomizableMap(table1, table2, customFields: false,
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld9, to_field: fld10, trim: false]])

		map = CustomizableMap(table1, table2, customFields: #(),
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld9, to_field: fld10, trim: false]])

		customFields = Object()
		customFields[fld10] = #(only_fillin_from: 'ScreenField2')
		map = CustomizableMap(table1, table2, :customFields,
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map is: [])

		customFields = Object()
		customFields[fld10] = #(only_fillin_from: 'ScreenField1')
		map = CustomizableMap(table1, table2, :customFields,
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld9, to_field: fld10, trim: false]])

		customFields = Object()
		customFields[fld10] = #(only_fillin_from: '')
		map = CustomizableMap(table1, table2, :customFields,
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld9, to_field: fld10, trim: false]])

		customFields = Object()
		customFields['myTestField'] = #(only_fillin_from: 'ScreenField2')
		map = CustomizableMap(table1, table2, :customFields,
			current_field: 'ScreenField1')
		Assert(map.CustomizableMap_map
			is: [[from_field: fld9, to_field: fld10, trim: false]])
		}

	Test_differentStringLengths()
		{
		table1 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table1,
			text: `class { Name: ` $ table1 $ ` }`])
		fld1 = .MakeCustomField(table1, 'Text, multi line', 'TEST FIELD1').field

		table2 = .MakeTable('(key, field) key(key)')
		.MakeLibraryRecord([name: "Table_" $ table2,
			text: `class { Name: ` $ table2 $ ` }`])
		fld2 = .MakeCustomField(table2, 'Text, single line', 'TEST FIELD1').field

		map = CustomizableMap(table1, table2)
		Assert(map.CustomizableMap_map is: [[from_field: fld1, to_field: fld2, trim:]])

		fromRec = [a: 1, b: 2]
		fromRec[fld1] = 'Simple Test'

		toRec = [a: 1, c: 3]
		map.CopyCustomFields(toRec, fromRec)
		Assert(toRec[fld2] is: 'Simple Test')

		fromRec[fld1] = 'A'.Repeat(600)
		map.CopyCustomFields(toRec, fromRec)
		// 512 - 3 for elipsis
		Assert(toRec[fld2] is: 'A'.Repeat(509) $ '...')
		}
	}