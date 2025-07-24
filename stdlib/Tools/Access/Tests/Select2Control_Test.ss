// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_layout_printParams()
		{
		ctrl = Select2Control
			{
			Getter_Select2Control_sf()
				{
				return FakeObject(Prompts: function ()
					{ return Object("One", "Two", "Three") })
				}
			Getter_Select2Control_ops_desc()
				{
				return Select2.TranslateOps()
				}
			}
		layout = ctrl.Select2Control_layout2(
			printParams: true, option: '', title: '', menuOptions: false)
		Assert(layout is: .layout_printParams())
		}

	layout_printParams()
		{
		n = Select2Control.Numrows
		layout_printParams_ob = Object(#(("Static", ""),
			("Horz", ("Skip", 4), ("Static", "Field")),
			("Horz", ("Skip", 4), ("Static", "Operator")),
			("Horz", ("Skip", 4), ("Static", "Value")),
			("Horz", ("Static", "Print?"))))
		for i in .. n
			{
			layout_printParams_ob.Add(Object(Object("CheckBox", name: "checkbox" $ i),
				Object("ChooseList",
					#("One", "Three", "Two"), width: 20, name: "fieldlist" $ i),
				Object("ChooseList",
					#("greater than", "greater than or equal to", "less than",
						"less than or equal to", "equals", "not equal to", "empty",
						"not empty", "contains", "does not contain", "starts with",
						"ends with", "matches", "does not match"), name: "oplist" $ i),
				Object("Field", name: "val" $ i), Object("CheckBox", name: "print" $ i)))
			}
		return Object('Record',
			Object('Vert',
				Object('Horz'
					#(Static "Checkmark the rows you want to apply")
					'Fill'
					'Skip')
				'Skip',
				Object("Grid", layout_printParams_ob))
			)
		}

	Test_layout()
		{
		ctrl = Select2Control
			{
			Getter_Select2Control_sf()
				{
				return FakeObject(Prompts: function ()
					{ return Object("One", "Two", "Three") })
				}
			Getter_Select2Control_ops_desc()
				{
				return Select2.TranslateOps()
				}
			}
		layout = ctrl.Select2Control_layout2(printParams: false,
			option: '', title: '', menuOptions: false)
		Assert(layout is: .expectedLayout())
		}

	expectedLayout()
		{
		n = Select2Control.Numrows
		layout_ob = Object(#(#("Static", ""),
			("Horz", ("Skip", 4), ("Static", "Field")),
			("Horz", ("Skip", 4), ("Static", "Operator")),
			("Horz", ("Skip", 4), ("Static", "Value"))))
		for i in .. n
			layout_ob.Add(Object(Object("CheckBox", name: "checkbox" $ i),
				Object("ChooseList",
					#("One", "Three", "Two"), width: 20, name: "fieldlist" $ i),
				Object("ChooseList", #("greater than", "greater than or equal to",
					"less than", "less than or equal to", "equals", "not equal to",
					"empty", "not empty", "contains", "does not contain", "starts with",
					"ends with", "matches", "does not match"), name: "oplist" $ i),
				Object("Field", name: "val" $ i)))
		return Object('Record',
			Object('Vert',
				Object('Horz'
					#(Static "Checkmark the rows you want to apply")
					'Fill'
					'Skip')
				'Skip',
				Object("Grid", layout_ob))
			)
		}

	Test_checkMenuOptions()
		{
		select = Select2Control { Numrows: 8 }
		data = Record()
		Assert(select.Select2Control_checkMenuOptions(data))
		data = Record(fieldlist1: 'test', menu_option1: true)
		Assert(select.Select2Control_checkMenuOptions(data))
		data = Record(fieldlist1: 'test', menu_option1: true,
			fieldlist2: 'test2', menu_option2: true,
			fieldlist3: 'test3', menu_option3: true)
		Assert(select.Select2Control_checkMenuOptions(data))
		data = Record(fieldlist1: 'test', menu_option1: true,
			fieldlist2: 'test2', menu_option2: true,
			fieldlist3: 'test3', menu_option3: true,
			fieldlist4: 'test3', menu_option4: true)
		Assert(select.Select2Control_checkMenuOptions(data)
			is: "Can't have the same field as a Menu Option more than once: test3")
		data = Record(fieldlist1: 'test', menu_option1: true,
			fieldlist2: 'test2', menu_option2: true,
			fieldlist3: 'test3', menu_option3: true,
			fieldlist4: 'test3', menu_option4: true,
			fieldlist5: 'test', menu_option5: true,
			fieldlist6: 'test2', menu_option6: true)
		Assert(select.Select2Control_checkMenuOptions(data)
			is: "Can't have the same field as a Menu Option more than once:" $
				" test3,test,test2")
		data = Record(fieldlist1: 'test', menu_option1: true,
			fieldlist2: 'test2', menu_option2: true,
			fieldlist3: 'test3', menu_option3: true,
			fieldlist4: 'test3', menu_option4: false,
			fieldlist5: 'test', menu_option5: false,
			fieldlist6: 'test2', menu_option6: false)
		Assert(select.Select2Control_checkMenuOptions(data))
		}
	}
