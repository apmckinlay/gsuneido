// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(GetControlListInfo('') is: #())

		cl = GetControlListInfo
			{
			GetControlFromField(field /*unused*/)
				{
				return #(ChooseList #('test1', 'test2'))
				}
			}
		Assert(cl('testField') is: #('test1', 'test2'))
		}

	Test_specialChooseListWithHardcodedList()
		{
		.MakeLibraryRecord(Record(name: 'XYZDoNotUseControl',
			text: 'class { GetList() { return Object("one", "two", "three") } }'))
		ctrl = Object('XYZDoNotUse')
		m = GetControlListInfo.GetControlListInfo_getControlListInfo
		list = m(ctrl)
		Assert(list is: #(one two three))
		}

	Test_customKeyGetListInfo()
		{
		fn = GetControlListInfo.GetControlListInfo_customKeyGetListInfo
		ctrl = Object(customField: 'test_customKey_test')
		Database('create test_customKey_test_table (name, desc) key (name)')
		QueryOutput('test_customKey_test_table', [name: 'test', desc: "I'm Helping"])

		Assert(fn(ctrl) is: #(field: "name", query: "test_customKey_test_table",
			displayField: "name", columns: #(#name, #desc), customizeQueryCols: true,
			keys: false))
		}

	Test_keyGetListInfo()
		{
		getCtrlLstInfo = GetControlListInfo
			{
			GetControlListInfo_getRestrictions(ctrl /*unused*/, query /*unused*/)
				{
				return false
				}
			}
		fn = getCtrlLstInfo.GetControlListInfo_keyGetListInfo
		ctrlClass = class
			{
			Key_BuildQuery(@unused)
				{
				return 'eta_rates'
				}
			}
		ctrl = #(query: eta_rates, columns: #(etarate_id, etarate_desc),
			numField: etarate_num, nameField: etarate_id, excludeSelect: #(etarate_uom))
		Assert(fn(ctrl, ctrlClass)
			is: #(query: eta_rates, columns: #(etarate_id,etarate_desc),
				field: etarate_num, displayField: etarate_id, customizeQueryCols: false,
				keys: false, optionalRestrictions: #(), excludeSelect: #(etarate_uom)))

		ctrl = #(query: eta_rates, columns: #(etarate_id, etarate_desc),
			numField: etarate_num, nameField: etarate_id,
			excludeSelect: 'GetControlListInfo_Test.ExcludeSelectText')
		Assert(fn(ctrl, ctrlClass)
			is: #(query: eta_rates, columns: #(etarate_id,etarate_desc),
				field: etarate_num, displayField: etarate_id, customizeQueryCols: false,
				keys: false, optionalRestrictions: #(), excludeSelect: #(test)))


		ctrl = #(Key, eta_rates, etarate_num,  #(etarate_id, etarate_desc))
		Assert(fn(ctrl, ctrlClass)
			is: #(query: eta_rates, columns: #(etarate_id, etarate_desc),
				field: etarate_num, displayField: etarate_name, customizeQueryCols: false,
				keys: false, optionalRestrictions: #(), excludeSelect: #()))

		ctrl = #(query: eta_rates, columns: #(etarate_id, etarate_desc),
			numField: etarate_num, nameField: etarate_id, keys: etarate_id)
		Assert(fn(ctrl, ctrlClass)
			is: #(query: eta_rates, columns: #(etarate_id, etarate_desc),
				field: etarate_num, displayField: etarate_id, customizeQueryCols: false,
				keys: etarate_id, optionalRestrictions: #(), excludeSelect: #()))
		}

	ExcludeSelectText()
		{
		return #(test)
		}

	Teardown()
		{
		super.Teardown()
		if TableExists?("test_customKey_test_table")
			Database("destroy test_customKey_test_table")
		}
	}
