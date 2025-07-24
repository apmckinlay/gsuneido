// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		cl = ChooseColumns
			{
			ChooseColumns_getDefaultWidth(col /*unused*/)
				{
				return 15
				}
			}
		field1 = .MakeDatadict(
			fieldName: 'choose_column_test_1',
			Prompt: 'ChooseColumns Test1')
		field2 = .MakeDatadict(
			fieldName: 'choose_column_test_2',
			Prompt: 'ChooseColumns Test2')
		duplicateCustomPrompt = 'Duplicate Custom Prompt'
		uniqueCustomSelectPrompt = 'Unique Select Prompt'
		customField1 = .MakeDatadict(
			fieldName: 'custom_99999998'
			Prompt: duplicateCustomPrompt)
		customField2 = .MakeDatadict(
			fieldName: 'custom_99999999',
			Prompt: duplicateCustomPrompt,
			SelectPrompt: uniqueCustomSelectPrompt)

		.chooseColumn = cl(Object(field1, field2, customField1, customField2))
		.origSettings = Suneido.GetDefault('Settings', []).Copy()
		}

	Test_AvailableList()
		{
		Assert(.chooseColumn.AvailableList() equalsSet:
			#((column: 'ChooseColumns Test1', col_width: 15),
			  (column: 'ChooseColumns Test2', col_width: 15),
			  (column: 'Duplicate Custom Prompt', col_width: 15),
			  (column: 'Unique Select Prompt', col_width: 15)))
		}

	Test_GetSaveData()
		{

		Assert(.chooseColumn.GetSaveData(
			#((column: 'ChooseColumns Test1', col_width: 15),
			  (column: 'ChooseColumns Test2', col_width: 15),
			  (column: 'Duplicate Custom Prompt', col_width: 15)))
			 is: #(choose_column_test_1, choose_column_test_2, custom_99999998))

		Assert(.chooseColumn.GetSaveData(
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 15),
			  (column: 'Unique Select Prompt', col_width: 15)))
			 is:
			 #(#(choose_column_test_1, width: 20), choose_column_test_2, custom_99999999))

		Assert(.chooseColumn.GetSaveData(
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 16)))
			 is:
			 #((choose_column_test_1, width: 20),
			   (choose_column_test_2, width: 16)))
		}

	Test_SetSaveList()
		{
		Assert(.chooseColumn.SetSaveList(
			#('ChooseColumns Test1', 'ChooseColumns Test2'))
			is:
			#((column: 'ChooseColumns Test1', col_width: 15)
			  (column: 'ChooseColumns Test2', col_width: 15)))

		Assert(.chooseColumn.SetSaveList(
			#(#('ChooseColumns Test1', width: 20), 'ChooseColumns Test2'))
			 is:
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 15)))

		Assert(.chooseColumn.SetSaveList(
			#(("ChooseColumns Test1", width: 20),
			  ('ChooseColumns Test2', width: 16)))
			is:
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 16)))
		}

	Test_ValidList()
		{
		Assert(ChooseColumns.ValidList(
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 16)))
			is: '')

		Assert(ChooseColumns.ValidList(
			#((column: 'ChooseColumns Test1', col_width: -20),
			  (column: 'ChooseColumns Test2', col_width: 16)))
			is: 'ChooseColumns Test1 has an invalid Width')

		Assert(ChooseColumns.ValidList(
			#((column: 'ChooseColumns Test1', col_width: 20),
			  (column: 'ChooseColumns Test2', col_width: 'abc')))
			is: 'ChooseColumns Test2 has an invalid Width')
		}

	Test_FindColumn()
		{
		Assert(ChooseColumns.FindColumn(
			#(#("ChooseColumns Test1", width: 20), 'ChooseColumns Test2'),
			'ChooseColumns Test1') is: 0)

		Assert(ChooseColumns.FindColumn(
			#(#("ChooseColumns Test1", width: 20), 'ChooseColumns Test2'),
			'ChooseColumns Test2') is: 1)

		Assert(ChooseColumns.FindColumn(
			#(#("ChooseColumns Test1", width: 20), 'ChooseColumns Test2'),
			'choose_column_test_3') is: false)
		}

	Test_GetFieldName()
		{
		Assert(ChooseColumns.GetFieldName(
			#('ChooseColumns Test1', width: 20))
			is: 'ChooseColumns Test1')

		Assert(ChooseColumns.GetFieldName('ChooseColumns Test2')
			is: 'ChooseColumns Test2')

		Assert(ChooseColumns.GetFieldName(
			#("Wrap", width: 24, field: 'choose_column_test_3'))
			is: 'choose_column_test_3')
		}

	Test_getDefaultWidth()
		{
		fn = ChooseColumns.ChooseColumns_getDefaultWidth
		Assert(fn('garbageField-WillFallbackToFieldDefault') is: 11)
		Assert(fn('text') is: 21)
		Assert(fn('address1') is: 16)
		Assert(fn('address2') is: 16) // inherits

		// will use default width
		fieldText_NoWidth = .TempName()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldText_NoWidth,
				text: `class
					{
					Format: (Text)
					}`])
		Assert(fn(fieldText_NoWidth) is: 10)
		}

	Test_getDefaultWidth_uom()
		{
		fn = ChooseColumns.ChooseColumns_getDefaultWidth
		fieldUOM_WithWidths = .TempName()
		fieldUOM_NoWidths = .TempName()
		fieldUOM_UnNamedWidths = .TempName()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldUOM_WithWidths,
				text: `class
					{
					Format: (UOM numwidth: 10, uomwidth: 5, div: '  ' /*2 char div*/)
					}`],
			[name: "Field_" $ fieldUOM_NoWidths,
				text: `class
					{
					Format: (UOM)
					}`],
			[name: "Field_" $ fieldUOM_UnNamedWidths,
				text: `class
					{
					Format: (UOM false 22 11)
					}`]
			)
		Assert(fn(fieldUOM_WithWidths) is: 17)

		// comes from default uom/num/div widths in UOMFormat. Changing those will make this fail
		Assert(fn(fieldUOM_NoWidths) is: 21)
		Assert(fn(fieldUOM_UnNamedWidths) is: 34)
		}

	Test_getDefaultWidth_number()
		{
		fn = ChooseColumns.ChooseColumns_getDefaultWidth
		fieldNumber_NoWidth = .TempName()
		fieldNumber_WithMask = .TempName()
		fieldNumber_WithWidth = .TempName()
		fieldNumber_WidthAndMask = .TempName()
		.MakeLibraryRecord(
			[name: "Field_" $ fieldNumber_NoWidth,
				text: `class
					{
					Format: (Number)
					}`],
			[name: "Field_" $ fieldNumber_WithMask,
				text: `class
					{
					Format: (OptionalNumber mask: '-###.##')
					}`],
			[name: "Field_" $ fieldNumber_WithWidth,
				text: `class
					{
					Format: (Number width: 8)
					}`],
			[name: "Field_" $ fieldNumber_WidthAndMask,
				text: `class
					{
					Format: (DollarFormat width: 11 mask: '###,###,###.##')
					}`]
			)
		Assert(fn(fieldNumber_NoWidth) is: 10)
		Assert(fn(fieldNumber_WithMask) is: 7)
		Assert(fn(fieldNumber_WithWidth) is: 8)
		Assert(fn(fieldNumber_WidthAndMask) is: 11)
		}

	origSettings: false
	Test_getDefaultWidth_dates()
		{
		fn = ChooseColumns.ChooseColumns_getDefaultWidth
		field_DateTimeSec = .TempName()
		field_ShortDate = .TempName()
		field_DateTime = .TempName()
		field_DateTime_WithWitdth = .TempName()
		field_LongDate = .TempName()
		.MakeLibraryRecord(
			[name: "Field_" $ field_DateTimeSec,
				text: `class
					{
					Format: (DateTimeSec)
					}`],
			[name: "Field_" $ field_ShortDate,
				text: `class
					{
					Format: (ShortDate)
					}`],
			[name: "Field_" $ field_DateTime,
				text: `class
					{
					Format: (DateTime)
					}`],
			[name: "Field_" $ field_DateTime_WithWitdth,
				text: `class
					{
					Format: (DateTime, width: 15)
					}`],
			[name: "Field_" $ field_LongDate,
				text: `class
					{
					Format: (LongDate)
					}`]
			)

		// Not affected by Settings. "1999-09-11 19:19:19" + 1
		Assert(fn(field_DateTimeSec) is: 20)

		Settings.Set('ShortDateFormat', "yyyy-MMM-dd")
		Assert(fn(field_ShortDate) is: 12) // "1999-Sep-11" + 1
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")
		Assert(fn(field_ShortDate) is: 11) // "1999-09-11" + 1

		Settings.Set('TimeFormat', "hh:mm tt")
		Assert(fn(field_DateTime) is: 20) // "1999-09-11 07:19 PM" + 1
		Settings.Set('TimeFormat', "h:mm")
		Assert(fn(field_DateTime) is: 16) // "1999-09-11 7:19" + 1
		Assert(fn(field_DateTime_WithWitdth) is: 15)
		Settings.Set('LongDateFormat', "dddd, dd MMMM, yyyy")
		Assert(fn(field_LongDate) is: 29) // "Saturday, 11 September, 1999" + 1
		Settings.Set('LongDateFormat', "ddd, dd MMM, yyyy")
		Assert(fn(field_LongDate) is: 18) // "Sat, 11 Sep, 1999" + 1
		}

	Test_getDefaultWidth_checkMark()
		{
		fn = ChooseColumns.ChooseColumns_getDefaultWidth
		field_CheckBox = .TempName()
		field_CheckMark_WithWidth = .TempName()
		field_CheckMark = .TempName()
		.MakeLibraryRecord(
			[name: "Field_" $ field_CheckBox,
				text: `class
					{
					Format: (CheckBox)
					}`],
			[name: "Field_" $ field_CheckMark_WithWidth,
				text: `class
					{
					Format: (CheckMarkFormat width: 5)
					}`],
			[name: "Field_" $ field_CheckMark,
				text: `class
					{
					Format: (CheckMarkFormat)
					}`]
			)
		Assert(fn(field_CheckBox) is: 4)
		Assert(fn(field_CheckMark_WithWidth) is: 5)
		Assert(fn(field_CheckMark) is: 4)
		}

	Teardown()
		{
		if .origSettings isnt false
			Suneido.Settings = .origSettings
		super.Teardown()
		}
	}