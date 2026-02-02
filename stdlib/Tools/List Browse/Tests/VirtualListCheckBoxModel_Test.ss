// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Setup()
		{
		super.Setup()
		.checkBoxColModel = VirtualListCheckBoxModel('test_checkbox_column', 'num', false)
		}

	Test_selectItems()
		{
		model = VirtualListModel(.VirtualList_Table,
			checkBoxColumn: 'test_checkbox_column', keyField: 'num')
		.AddTeardownModel(model)
		model.CheckRecordByKey(1)
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecordByKey(1, col: 'other_field')
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecordByKey(100, forceCheck:)
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecord([num: 100], forceCheck:)
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecordByKey(1, forceCheck:)
		Assert(model.GetCheckedRecords().list
			is:	#([test_checkbox_column: true, num: 1]))

		model.CheckRecordByKey(2, 'test_checkbox_column')
		list = model.GetCheckedRecords().list
		expected = #(1: [test_checkbox_column: true, num: 1],
			2: [test_checkbox_column: true, num: 2])
		Assert(list is: expected.Values())

		model.CheckRecordByKey(1, 'test_checkbox_column')
		Assert(model.GetCheckedRecords().list
			is: #([test_checkbox_column: true, num: 2]))

		model.CheckAll()
		Assert(model.GetCheckedRecords().state is: 'all')

		model.CheckRecordByKey(1, forceCheck:)
		Assert(model.GetCheckedRecords().state is: 'allbut')
		Assert(model.GetCheckedRecords().list
			is: #([test_checkbox_column: false, num: 1]))

		model.UncheckAll()
		Assert(model.GetCheckedRecords().state is: 'selected')

		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)

		model.CheckRecordByKey(1, col: 'other_field')
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecordByKey(100, forceCheck:)
		Assert(model.GetCheckedRecords().list is: #())

		model.CheckRecordByKey(1, forceCheck:)
		rec = model.GetCheckedRecords().list[0]
		Assert(rec.test_checkbox_column)
		Assert(rec.num is: 1)
		}

	Test_select()
		{
		// select one item
		.selectItem('Selected Item Test')
		Assert(.isSelected('Selected Item Test'),
			msg: 'Select One Item IsSelected Failed')
		selected = .checkBoxColModel.GetSelectedInfo()
		Assert(selected.state is: 'selected', msg: 'Select Item state Failed')
		Assert(selected.list is: .selectedList(#("Selected Item Test")),
			msg: 'Select One Item GetSelected Failed')
		Assert(.isSelected('Not Selected Item') is: false,
			msg: 'Select One Item Failed')

		// select 'false' item
		.checkBoxColModel.SelectItem(false)
		Assert(.checkBoxColModel.GetSelectedInfo().list is:
			.selectedList(#('Selected Item Test')), msg: 'Select False Item Failed')

		.selectItem('Second Selected Item Test')
		// unselect and item
		.unselectItem('Selected Item Test')
		Assert(.isSelected('Selected Item Test') is: false,
			msg: 'Unselect Item IsSelected Failed')
		Assert(.checkBoxColModel.GetSelectedInfo().list is:
			.selectedList(#('Second Selected Item Test')),
				msg: 'Unselect Item GetSelected Failed')

		// unselect false items
		.checkBoxColModel.UnselectItem(false)
		Assert(.checkBoxColModel.GetSelectedInfo().list is:
			.selectedList(#('Second Selected Item Test')),
				msg: 'Unselect False Item Failed')

		// select all items
		.checkBoxColModel.SelectAll()
		selected = .checkBoxColModel.GetSelectedInfo()
		Assert(selected.state is: 'all', msg: 'Select All Failed')
		Assert(.isSelected('Test Item'), msg: 'Select All IsSelected Failed')

		// unselect one
		.unselectItem('Unselect This Item Test')
		selected = .checkBoxColModel.GetSelectedInfo()
		Assert(.isSelected('Unselect This Item Test') is: false,
			msg: 'Unselect One IsSelected Failed')
		Assert(selected.state is: 'allbut', msg: 'Select AllBut Failed')
		Assert(selected.list is:
			#([test_checkbox_column: false, num: 'Unselect This Item Test']),
				msg: 'Unselect One GetSelected Failed')

		// unselect all
		.checkBoxColModel.UnselectAll()
		Assert(.checkBoxColModel.GetSelectedInfo().list is: #(),
			msg: 'Unselect All failed')
		Assert(.isSelected('Test Item') is: false, msg: 'Unselect All IsSelected Failed')
		}

	selectItem(item)
		{
		.checkBoxColModel.SelectItem([num: item])
		}

	isSelected(item)
		{
		return .checkBoxColModel.IsSelected([num: item])
		}

	unselectItem(item)
		{
		.checkBoxColModel.UnselectItem([num: item])
		}

	selectedList(items)
		{
		return items.Map({ [test_checkbox_column: true, num: it] })
		}

	Test_totalSelected()
		{
		model = VirtualListCheckBoxModel('test_checkbox_col', 'num', 'test_amount_col')
		Assert(model.GetSelectedTotal() is: 0)

		model.SelectItem([test_checkbox_col: true, num: 1, test_amount_col: 1])
		model.SelectItem([test_checkbox_col: true, num: 2, test_amount_col: 2])
		model.SelectItem([test_checkbox_col: true, num: 3, test_amount_col: 3])
		model.SelectItem([test_checkbox_col: true, num: 4, test_amount_col: 4])

		Assert(model.GetSelectedTotal(recalc:) is: 10)
		Assert(model.GetSelectedTotal() is: 10)
		}

	Test_autoSelectByAmount()
		{
		model = VirtualListCheckBoxModel('test_checkbox_col', 'num', 'test_amount_col')
		model.SelectItem([test_checkbox_col: true, num: 1, test_amount_col: 1])
		model.SelectItem([test_checkbox_col: true, num: 2, test_amount_col: 2])
		model.SelectItem([test_checkbox_col: true, num: 3, test_amount_col: 3])
		model.SelectItem([test_checkbox_col: true, num: 4, test_amount_col: 4])
		Assert(model.GetSelectedInfo().list.Size() is: 4)

		model.AutoSelectByAmount('test_amount_col', 0,
			[test_checkbox_col: true, num: 3, test_amount_col: 0])
		Assert(model.GetSelectedInfo().list.Size() is: 3)

		model.AutoSelectByAmount('test_checkbox_col', 0,
			[test_checkbox_col: false, num: 3, test_amount_col: ])
		Assert(model.GetSelectedInfo().list.Size() is: 3)

		model.AutoSelectByAmount('test_amount_col', 15,
			[test_checkbox_col: false, num: 3, test_amount_col: 15])
		Assert(model.GetSelectedInfo().list.Size() is: 4)


		}
	}
