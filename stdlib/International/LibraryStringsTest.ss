// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		table1 = .MakeTable('(trlang_from) key (trlang_from)')
		table2 = .MakeTable('(num, parent, group, name, text) key(num)')
		QueryOutput(table2, Object(num: 666, parent:0, group: -1, name: 'testlib',
			text: 'Controller
				{
				Controls:
				(Record
						(Vert
								(RadioGroups
									(Horz name: "HorzDaily"
										(Field name: "daily")
										label: "DAILY")
										#(Field name: "weekly" label: Weekly)
										#(Horz
											(ComboBox "1" "&Two" "Three..." "Four"
												"Five" name: numbers)
											label: Monthly))
								(Button Inspect)
								(Button Update)))
				On_Inspect()
						{
						.list = .Data.Get()
						Alert(.list.RadioGroups)
						}
				On_Update()
					{
					.Data.Vert.RadioGroups.Set("DAILY")
			//		Alert(.Data.Vert.RadioGroups.RadioGroupsControl_args[0])
					.Data.SetField("daily", "Hello")
					.Data.SetField("weekly", "Goodbye")
					.Data.SetField("numbers", "One")
					}
				}'))
		QueryOutput(table1, #(trlang_from: 'Four'))
		n = LibraryStrings(table2, table1, quiet:)
		Assert(n is: 11)
		}
	}