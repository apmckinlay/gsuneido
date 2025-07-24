// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Setup()
		{
		super.Setup()

		.protectField = .TempTableName()
		.MakeLibraryRecord([name: 'Rule_' $ .protectField,
			text: `function () { return '' }`])

		.grid = Mock()
		.grid.Window = Object(Hwnd: 0)
		.grid.Controller = Mock()
		.grid.Controller.When.Send([anyArgs:]).Return(true)
		.grid.Controller.Addons = Mock()
		.grid.Controller.Addons.When.Collect([anyArgs:]).Return(Object(''))
		.updateHistoryFn = function(@unused) {}
		}

	Test_main()
		{
		// delete
		model = VirtualListModel(.VirtualList_Table, protectField: .protectField)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		rec0 = model.GetLoadedData()[0]  // record 0
		Assert(Query1(.VirtualList_Table, num: 0) isnt: false)
		VirtualListSave(.grid, model).DeleteRecord(rec0)
		Assert(Query1(.VirtualList_Table, num: 0) is: false)

		list = Mock()
		list.When.GetModel().Return(model)
		list.When.GetColumns().Return(#('num', 'field1'))
		list.Parent = [Window: [Parent: 0]]

		// modify
		rec9 = model.GetLoadedData()[9]  // record 9
		rec9.PreSet('vl_list', list)
		rec9.field1 = 'rec 9'
		Assert(Query1(.VirtualList_Table, num: 9).field1 is: '')
		VirtualListSave(.grid, model, .updateHistoryFn).SaveRecord(rec9)
		Assert(Query1(.VirtualList_Table, num: 9).field1 is: 'rec 9')

		// new
		recNew = Record()
		recNew.vl_list = list
		recNew.num = 60
		recNew.Observer(VirtualListObserverOnChange)
		recNew.field1 = 'new record'
		model.GetLoadedData().Add(recNew)
		Assert(Query1(.VirtualList_Table, num: 60) is: false)
		VirtualListSave(.grid, model, .updateHistoryFn).SaveRecord(recNew)
		Assert(Query1(.VirtualList_Table, num: 60) isnt: false)
		}

	Test_saveQuery()
		{
		// testing saveQuery option
		// using TS field to determine if original record is changed or not.
		Database('ensure ' $ .VirtualList_Table $ ' (num_TS)')
		viewQuery = .VirtualList_Table $ ' summarize num, field1, num_TS, count'
		model = VirtualListModel(viewQuery, saveQuery: .VirtualList_Table,
			protectField: .protectField)
		.AddTeardownModel(model)

		list = Mock()
		list.When.GetModel().Return(model)
		list.When.GetColumns().Return(#('num', 'field1'))
		list.Parent = [Window: [Parent: 0]]

		model.UpdateVisibleRows(10)
		rec1 = model.GetLoadedData()[0]
		Assert(Query1(.VirtualList_Table, num: rec1.num) isnt: false)
		VirtualListSave(.grid, model).DeleteRecord(rec1)
		Assert(Query1(.VirtualList_Table, num: rec1.num) is: false)

		// modify
		rec5 = model.GetLoadedData()[4]
		rec5.PreSet('vl_list', list)
		rec5.field1 = 'changed'
		Assert(Query1(.VirtualList_Table, num: rec5.num).field1 is: '')
		VirtualListSave(.grid, model, .updateHistoryFn).SaveRecord(rec5)
		Assert(Query1(.VirtualList_Table, num: rec5.num).field1 is: 'changed')

		// new
		recNew = Record()
		recNew.vl_list = list
		recNew.num = 50
		recNew.Observer(VirtualListObserverOnChange)
		recNew.field1 = 'new record'
		model.GetLoadedData().Add(recNew)
		Assert(Query1(.VirtualList_Table, num: 50) is: false)
		VirtualListSave(.grid, model, .updateHistoryFn).SaveRecord(recNew)
		Assert(Query1(.VirtualList_Table, num: 50) isnt: false)
		}

	Test_delete_keyException()
		{
		newTable = .MakeTable('(test_num, num) index (num) in ' $ .VirtualList_Table $
			' key(test_num)')
		keyExceptSpy = .SpyOn(KeyException).Return('')

		model = VirtualListModel(.VirtualList_Table, protectField: .protectField)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)

		recLast = model.GetLoadedData().Last()
		QueryOutput(newTable, [test_num: Timestamp(), num: recLast.num])

		Assert(Query1(.VirtualList_Table, num: recLast.num) isnt: false)
		VirtualListSave(.grid, model).DeleteRecord(recLast)
		Assert(keyExceptSpy.CallLogs() isSize: 1)
		Assert(keyExceptSpy.CallLogs()[0].e has: 'blocked by foreign key')
		Assert(Query1(.VirtualList_Table, num: recLast.num) isnt: false)
		}
	}