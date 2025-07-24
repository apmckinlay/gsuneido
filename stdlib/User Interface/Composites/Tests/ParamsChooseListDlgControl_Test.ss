// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getFieldSaveName()
		{
		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: '' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: '')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'barney' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'barney')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'test_number' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'test_number')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'test_identifier' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'test_identifier')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'field_date' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'field_date')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'bizpartner_num' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'bizpartner_num')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'bizpartner_num_supplier' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'bizpartner_num')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'gldept_id' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'gldept_id')

		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'gldept_id_param' }
		Assert(c.ParamsChooseListDlgControl_getFieldSaveName() is: 'gldept_id')
		}

	Test_buildEditorDisplayValue()
		{
		c = ParamsChooseListDlgControl
			{ ParamsChooseListDlgControl_field: 'string' }


		m = c.ParamsChooseListDlgControl_buildEditorDisplayValue
		Assert(m(#()) is: "")
		Assert(m(#("Test")) is: "Test")
		Assert(m(#("Test" "Test2" "Test3" "Test4")) is: "Test,Test2,Test3,Test4")
		inputOb = Object()
		i = 0
		for .. 250
			inputOb.Add(Display(++i))
		expected = inputOb[..200].Join(',') $ '...'
		Assert(m(inputOb) is: expected)
		}

	Test_getFieldPrompt()
		{
		cl = ParamsChooseListDlgControl
			{
			ParamsChooseListDlgControl_getSelectPrompt(field)
				{
				promptMap = [
					equipment1_num: 'Equipment1',
					equipment1_name: 'Equipment1 Name',
					equipment1_abbrev: 'Equipment1 Abbrev',
					equipment2_num: 'Equipmenttt2',
					equipment2_name: 'Equipment2 Namer',
					equipment2_abbrev: 'Equipment2 Abbreviation',
					equipment3_num: 'Equipment3'
					]

				return promptMap[field] isnt '' ? promptMap[field] : field
				}
			}

		fn = cl.ParamsChooseListDlgControl_getFieldPrompt
		Assert(fn('') is: '')

		Assert(fn('equipment_noprompt') is: 'equipment_noprompt')
		Assert(fn('equipment_noprompt_num') is: 'equipment_noprompt_num')
		Assert(fn('equipment_num_noprompt') is: 'equipment_num_noprompt')

		Assert(fn('equipment1_num') is: 'Equipment1')
		Assert(fn('equipment1_name') is: 'Equipment1 Name')
		Assert(fn('equipment1_abbrev') is: 'Equipment1 Abbrev')

		Assert(fn('equipment2_num') is: 'Equipmenttt2')
		Assert(fn('equipment2_name') is: 'Equipment2 Namer')
		Assert(fn('equipment2_abbrev') is: 'Equipment2 Abbreviation')

		Assert(fn('equipment3_num') is: 'Equipment3')
		Assert(fn('equipment3_name') is: 'Equipment3 Name')
		Assert(fn('equipment3_abbrev') is: 'Equipment3 Abbrev')
		}
	}