// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_InsertNewRecord_DeleteRecord()
		{
		model = .model('num to test_timestamp', 'test_timestamp')
		model.UpdateVisibleRows(20)

		rec = model.InsertNewRecord(record: [hello: 'world'], row_num: 0)
		Assert(rec.New?())
		Assert(rec.test_timestamp isDate:)
		Assert(rec.hello is: 'world')
		Assert(model.GetRecord(0) is: rec)
		Assert(model.GetRecord(1).test_timestamp is: 0)

		rec2 = model.InsertNewRecord(row_num: 9)
		Assert(rec2.New?())
		Assert(rec2.test_timestamp isDate:)
		Assert(model.GetRecord(9) is: rec2)
		Assert(model.GetRecord(10).test_timestamp is: 8)

		model.DeleteRecord(rec)
		Assert(model.GetRecord(0).test_timestamp is: 0)
		Assert(model.GetRecord(9).test_timestamp is: 8)

		model.DeleteRecord(rec2)
		Assert(model.GetRecord(0).test_timestamp is: 0)
		Assert(model.GetRecord(8).test_timestamp is: 8)

		model.UpdateOffset(25)
		rec3 = model.InsertNewRecord()
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(0).test_timestamp is: 11)
		Assert(model.GetRecord(19) is: rec3)
		Assert(model.GetRecord(20) is: false) // one extra empty row at end
		}

	model(rename = '', sort = 'num', startLast = false, classOverride? = false,
		customKey = false)
		{
		cl = classOverride?
			? VirtualListModel
				{
				VirtualListModel_limit: 10
				VirtualListModel_segment: 3
				VirtualListModel_closeCursorsIfAllRead() {}
				}
			: VirtualListModel
		query = .VirtualList_Table $ Opt(' rename ', rename) $ Opt(' sort ', sort)
		model = cl(query, :startLast, :customKey)

		mock = Mock()
		mock.When.GetModel().Return(model)
		mock.When.GetColumns().Return(QueryColumns(query))
		model.VirtualListModel_observerList = mock
		.AddTeardownModel(model)
		return model
		}

	Test_InsertNewRecord_DeleteRecord_reverse()
		{
		model = .model('num to test_timestamp', 'test_timestamp', startLast:)
		model.UpdateVisibleRows(20)

		rec = model.InsertNewRecord(row_num: 0)
		Assert(rec.New?())
		Assert(rec.test_timestamp isDate:)
		Assert(model.GetRecord(0) is: rec)
		Assert(model.GetRecord(1).test_timestamp is: 11)

		rec2 = model.InsertNewRecord(row_num: 9)
		Assert(rec2.New?())
		Assert(rec2.test_timestamp isDate:)
		Assert(model.GetRecord(9) is: rec2)
		Assert(model.GetRecord(10).test_timestamp is: 19)

		model.DeleteRecord(rec)
		Assert(model.GetRecord(0).test_timestamp is: 11)
		Assert(model.GetRecord(1).test_timestamp is: 12)

		model.DeleteRecord(rec2)
		Assert(model.GetRecord(9).test_timestamp is: 20)
		Assert(model.GetRecord(10).test_timestamp is: 21)

		model.UpdateOffset(25)
		rec3 = model.InsertNewRecord()
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(0).test_timestamp is: 11)
		Assert(model.GetRecord(19) is: rec3)
		Assert(model.GetRecord(20) is: false)
		}

	Test_InsertNewRecord_after_recyled_at_end()
		{
		model = .model(classOverride?:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(25, .FakeSaveAndCollapse)

		rec3 = model.InsertNewRecord()
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(5) is: rec3)

		model = .model(startLast:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(-25, .FakeSaveAndCollapse)

		rec3 = model.InsertNewRecord(row_num: 0)
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(0) is: rec3)
		}

	Test_InsertNewRecord_after_recyled()
		{
		model = .model(classOverride?:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(25, .FakeSaveAndCollapse)

		rec3 = model.InsertNewRecord(row_num: 2)
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(0).num is: 25)
		Assert(model.GetRecord(1).num is: 26)
		Assert(model.GetRecord(2) is: rec3)
		Assert(model.GetRecord(3).num is: 27)
		}

	Test_InsertNewRecord_after_recyled_revese()
		{
		model = .model(startLast:, classOverride?:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(-26, .FakeSaveAndCollapse)
		Assert(model.GetRecord(0).num is: 0)

		rec3 = model.InsertNewRecord(row_num: 2)
		Assert(rec3.New?())
		Assert(rec3.test_timestamp isDate:)
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.GetRecord(1).num is: 1)
		Assert(model.GetRecord(2) is: rec3)
		Assert(model.GetRecord(3).num is: 2)
		}

	Test_DeleteRecord_afterRecyled()
		{
		model = .model(classOverride?:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(25, .FakeSaveAndCollapse)
		rec = model.GetRecord(4)
		Assert(rec.num is: 29)
		model.DeleteRecord(rec)
		Assert(model.GetRecord(4) is: false)

		rec = model.GetRecord(0)
		Assert(rec.num is: 25)
		model.DeleteRecord(rec)
		Assert(model.GetRecord(0).num is: 26)
		}

	Test_DeleteRecord_afterRecyled_reverse()
		{
		model = .model(startLast:, classOverride?:)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(-26, .FakeSaveAndCollapse)
		rec = model.GetRecord(0)
		Assert(rec.num is: 0)
		model.DeleteRecord(rec)
		Assert(model.GetRecord(0).num is: 1)

		rec = model.GetRecord(4)
		Assert(rec.num is: 5)
		model.DeleteRecord(rec)
		Assert(model.GetRecord(4).num is: 6)
		}

	Test_insert_at_end()
		{
		// testing when visible rows not being set
		model = .model()
		rec1 = model.InsertNewRecord(record: [hello: 'world 1'], row_num: false)
		rec2 = model.InsertNewRecord(record: [hello: 'world 2'], row_num: false)
		Assert(model.GetRecord(0) is: rec1)
		Assert(model.GetRecord(1) is: rec2)

		// load all
		model = .model()
		model.UpdateVisibleRows(10)

		rec = model.InsertNewRecord(record: [hello: 'world'], row_num: false)
		Assert(model.GetRecord(10) is: rec)

		// load all and read reversely
		model.SetStartLast(true)
		rec = model.InsertNewRecord(record: [hello: 'world'], row_num: false)
		Assert(model.GetRecord(9) is: rec)

		_stopLoadAll = true
		model = .model(startLast:)
		model.UpdateVisibleRows(10)

		rec = model.InsertNewRecord(record: [hello: 'world'], row_num: false)
		Assert(model.GetRecord(9) is: rec)

		model.SetStartLast(false)
		rec = model.InsertNewRecord(record: [hello: 'world'], row_num: false)
		Assert(model.GetRecord(10) is: rec)
		}

	Test_new_record_highlightInvalidFields()
		{
		.SpyOn(VirtualListColModel.VirtualListColModel_loadSavedCols).Return('')
		mandatoryField = .MakeCustomField(
			.VirtualList_Table, 'Text, single line', prompt: 'Mandatory Test')
		.MakeDatadict(fieldName: 'field1', Prompt: 'Field 1')

		record = []
		model = .model(customKey: .VirtualList_Table)
		model.InsertNewRecord(:record, row_num: false)
		Assert(model.EditModel.HasInvalidCols?(model.GetRecord(0)) is: false)

		.MakeCustomizeField(
			.VirtualList_Table, mandatoryField.field, extrafields: [custfield_mandatory:],
			key: .VirtualList_Table)
		record = []
		model = .model(customKey: .VirtualList_Table)
		model.InsertNewRecord(:record, row_num: false)
		Assert(model.EditModel.HasInvalidCols?(model.GetRecord(0)))

		record = []
		model = .model(customKey: .VirtualList_Table)
		record[mandatoryField.field] = 'value'
		model.InsertNewRecord(:record, row_num: false)
		Assert(model.EditModel.HasInvalidCols?(record = model.GetRecord(0)) is: false)
		}
	}
