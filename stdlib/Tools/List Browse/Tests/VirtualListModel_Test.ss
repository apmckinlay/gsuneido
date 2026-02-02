// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_UpdateVisibleRows()
		{
		mock = Mock(VirtualListModel)
		mock.Offset = 0
		mock.VisibleRows = 0
		mock.VirtualListModel_startLast = false
		mock.VirtualListModel_data = []
		mock.When.loadAllIfSmall([anyArgs:]).Return(false)
		mock.When.updateVisibleRows?([anyArgs:]).CallThrough()
		mock.When.keepAtBottom?([anyArgs:]).CallThrough()
		mock.VirtualListModel_curBottom = Mock()
		mock.VirtualListModel_curBottom.When.ReadDown([anyArgs:]).Return(1)
		curBottom = mock.VirtualListModel_curBottom
		curBottom.Pos = 0
		curBottom.Seeking = false
		mock.When.VirtualListModel_overLimit?().Return(false)
		mock.Eval(VirtualListModel.UpdateVisibleRows, 40)
		Assert(mock.VisibleRows is: 40)
		}

	Test_UpdateVisibleRows_reverse()
		{
		mock = Mock(VirtualListModel)
		mock.Offset = 0
		mock.VisibleRows = 0
		mock.VirtualListModel_startLast = false
		mock.VirtualListModel_data = []
		mock.When.loadAllIfSmall([anyArgs:]).Return(false)
		mock.When.updateVisibleRows?([anyArgs:]).CallThrough()
		mock.When.keepAtBottom?([anyArgs:]).CallThrough()
		mock.VirtualListModel_curBottom = Mock()
		mock.VirtualListModel_curBottom.When.ReadDown([anyArgs:]).Return(1)
		curBottom = mock.VirtualListModel_curBottom
		curBottom.Pos = 1
		curBottom.Seeking = false
		mock.When.VirtualListModel_overLimit?().Return(false)
		mock.Eval(VirtualListModel.UpdateVisibleRows, 40)
		Assert(mock.VisibleRows is: 40)
		}

	Test_UpdateVisibleRows_begin_end()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 35)
		Assert(model.Begin?(), msg: 'begin 35')
		Assert(model.End?(), msg: 'end 35')

		model.UpdateVisibleRows(31)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 31)
		Assert(model.Begin?(), msg: 'begin 31')
		Assert(model.End?(), msg: 'end 31')

		model.UpdateVisibleRows(30)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 30)
		Assert(model.Begin?(), msg: 'begin 30')
		Assert(model.End?(), msg: 'end 30')

		model.UpdateVisibleRows(29)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 29)
		Assert(model.Begin?(), msg: 'begin 29')
		Assert(model.End?() is: false, msg: 'end 29')
		}

	Test_UpdateVisibleRows_begin_end_reverse()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(35)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 35)
		Assert(model.Begin?(), msg: 'begin 35')
		Assert(model.End?(), msg: 'end 35')

		model.SetStartLast(true)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 35)
		Assert(model.Begin?(), msg: 'begin -30')
		Assert(model.End?(), msg: 'end -30')

		model.UpdateVisibleRows(31)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 31)
		Assert(model.Begin?(), msg: 'begin 31')
		Assert(model.End?(), msg: 'end 31')

		model.UpdateVisibleRows(30)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 30)
		Assert(model.Begin?(), msg: 'begin 30')
		Assert(model.End?(), msg: 'end 30')

		model.UpdateVisibleRows(29)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 29)
		Assert(model.Begin?(), msg: 'begin 29')
		Assert(model.End?() is: false, msg: 'end 29')
		}

	Test_UpdateOffset()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 20)

		model.UpdateOffset(5)
		Assert(model.Offset is: 5)

		model.UpdateOffset(-5)
		Assert(model.Offset is: 0)

		model.UpdateOffset(-5, fromRefresh?:)
		Assert(model.Offset is: 0)

		model.UpdateOffset(-5)
		Assert(model.Offset is: 0)

		model.UpdateOffset(5)
		Assert(model.Offset is: 5)

		model.UpdateOffset(20)
		Assert(model.Offset is: 11)

		model.UpdateOffset(-5, fromRefresh?:)
		Assert(model.Offset is: 6)
		}

	Test_UpdateOffset_reverse()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 20)

		model.SetStartLast(true)
		Assert(model.Offset is: -19)
		Assert(model.VisibleRows is: 20)

		model.UpdateOffset(-5)
		Assert(model.Offset is: -24)

		model.UpdateOffset(10)
		Assert(model.Offset is: -19)

		model.UpdateOffset(-5)
		Assert(model.Offset is: -24)

		model.UpdateOffset(5)
		Assert(model.Offset is: -19)

		model.UpdateOffset(-5)
		Assert(model.Offset is: -24)

		model.UpdateOffset(-10)
		Assert(model.Offset is: -30)

		model.UpdateOffset(-20)
		Assert(model.Offset is: -30)
		}

	Test_reverse_option()
		{
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.Offset is: -19)
		Assert(model.VisibleRows is: 20)

		model.UpdateVisibleRows(20)
		Assert(model.Offset is: -19)
		Assert(model.VisibleRows is: 20)

		model.UpdateVisibleRows(25) // resize doesn't read more
		Assert(model.Offset is: -25)
		Assert(model.VisibleRows is: 25)

		model.UpdateVisibleRows(35)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 35)

		model.UpdateVisibleRows(10)
		Assert(model.Offset is: -9)
		Assert(model.VisibleRows is: 10)

		model.UpdateVisibleRows(20)
		Assert(model.Offset is: -20)
		Assert(model.VisibleRows is: 20)

		model.UpdateVisibleRows(35)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 35)
		}

	Test_SetRecordsToTop()
		{
		model = VirtualListModel(.VirtualList_Table $ ' sort num')
		.AddTeardownModel(model)

		model.UpdateVisibleRows(20)

		model.SetRecordsToTop('num', #(25, 15))

		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 20)

		Assert(model.GetRecord(0).num is: 15)
		Assert(model.GetRecord(1).num is: 25)
		Assert(model.GetRecord(2).num is: 0)
		Assert(model.GetRecord(3).num is: 1)

		model.SetSort('reverse num')
		model.SetRecordsToTop('num', #(25, 15))

		Assert(model.GetRecord(0).num is: 25)
		Assert(model.GetRecord(1).num is: 15)
		Assert(model.GetRecord(2).num is: 29)
		Assert(model.GetRecord(3).num is: 28)
		}

	Test_expand_collapse()
		{
		_stopLoadAll = true
		model = VirtualListModel(.VirtualList_Table $ ' sort num')
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.VirtualListModel_curBottom.Pos is: 21)

		model.SetRecordExpanded(0, 2)
		Assert(model.Offset is: 0)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 one')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 one')
		Assert(model.GetRecord(3).num is: 1)
		Assert(model.VirtualListModel_curBottom.Pos is: 23)

		model.SetRecordCollapsed(0)
		Assert(model.Offset is: 0)
		Assert(model.GetRecord(0).vl_expanded_rows is: '')
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.GetRecord(1).vl_expanded_rows is: '')
		Assert(model.GetRecord(1).num is: 1)
		Assert(model.VirtualListModel_curBottom.Pos is: 21)

		model.SetRecordExpanded(19, 3)
		Assert(model.Offset is: 0)
		Assert(model.GetRecord(19).vl_expanded_rows is: 3)
		Assert(model.VirtualListModel_curBottom.Pos is: 24)

		model.SetRecordCollapsed(19)
		Assert(model.Offset is: 0)
		Assert(model.GetRecord(19).vl_expanded_rows is: '')
		Assert(model.VirtualListModel_curBottom.Pos is: 21)

		model.UpdateOffset(5)
		Assert(model.VirtualListModel_curBottom.Pos is: 26)
		model.SetRecordExpanded(10, 2)
		Assert(model.GetRecord(10).vl_expanded_rows is: 2)
		Assert(model.GetRecord(11).vl_expand?, msg: 'expand 11')
		Assert(model.GetRecord(12).vl_expand?, msg: 'expand 12')
		Assert(model.VirtualListModel_curBottom.Pos is: 28)

		model.SetRecordCollapsed(10)
		Assert(model.GetRecord(10).vl_expanded_rows is: '')
		Assert(model.GetRecord(10).num is: 15)
		Assert(model.GetRecord(11).num is: 16)
		Assert(model.VirtualListModel_curBottom.Pos is: 26)

		// test collapse record above the .Offset
		model.UpdateOffset(-5)
		model.SetRecordExpanded(0, 2)
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 two')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 two')
		Assert(model.VirtualListModel_curBottom.Pos is: 29) // UpdateOffset+1, expand+2

		model.UpdateOffset(3)
		Assert(model.Offset is: 3)
		Assert(model.GetRecord(0).num is: 1)

		model.SetRecordCollapsed(-3, keepPosition?:)
		Assert(model.Offset is: 3)
		Assert(model.GetRecord(0).num is: 3)

		model.SetRecordExpanded(-3, 2) // restore
		Assert(model.Offset is: 3)

		model.SetRecordCollapsed(-3)
		Assert(model.Offset is: 1)
		Assert(model.GetRecord(0).num is: 1)
		Assert(model.GetRecord(-1).num is: 0)
		Assert(model.VirtualListModel_curBottom.Pos is: 28) // UpdateOffset+1, collapse-2
		}

	Test_expand_collapse_reverse()
		{
		_stopLoadAll = true
		model = VirtualListModel(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)

		Assert(model.Offset is: -19)
		Assert(model.VirtualListModel_curTop.Pos is: -20)
		Assert(model.GetRecord(0).num is: 11)

		model.SetRecordExpanded(0, 2)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(0).num is: 11)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 one')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 one')
		Assert(model.GetRecord(3).num is: 12)
		Assert(model.Offset is: -21)
		Assert(model.VirtualListModel_curTop.Pos is: -22)

		model.SetRecordCollapsed(0)
		Assert(model.GetRecord(0).vl_expanded_rows is: '')
		Assert(model.GetRecord(0).num is: 11)
		Assert(model.GetRecord(1).num is: 12)
		Assert(model.Offset is: -19)
		Assert(model.VirtualListModel_curTop.Pos is: -20)

		model.SetRecordExpanded(18, 3)
		Assert(model.Offset is: -22)
		Assert(model.GetRecord(18).vl_expanded_rows is: 3)
		Assert(model.VirtualListModel_curTop.Pos is: -23)

		model.SetRecordCollapsed(18)
		Assert(model.Offset is: -19)
		Assert(model.GetRecord(18).vl_expanded_rows is: '')
		Assert(model.VirtualListModel_curTop.Pos is: -20)

		model.UpdateOffset(-5)
		Assert(model.Offset is: -24)
		Assert(model.VirtualListModel_curTop.Pos is: -25)

		model.SetRecordExpanded(10, 2)
		Assert(model.GetRecord(10).vl_expanded_rows is: 2)
		Assert(model.GetRecord(11).vl_expand?, msg: 'expand 11')
		Assert(model.GetRecord(12).vl_expand?, msg: 'expand 12')
		Assert(model.VirtualListModel_curTop.Pos is: -27)
		Assert(model.Offset is: -26)

		model.SetRecordCollapsed(10)
		Assert(model.GetRecord(10).vl_expanded_rows is: '')
		Assert(model.GetRecord(10).num is: 16)
		Assert(model.GetRecord(11).num is: 17)
		Assert(model.Offset is: -24)
		Assert(model.VirtualListModel_curTop.Pos is: -25)

		// test collapse record above the .Offset
		model.SetRecordExpanded(0, 5) // expand at top of the viewport
		Assert(model.GetRecord(0).num is: 6)
		Assert(model.GetRecord(0).vl_expanded_rows is: 5)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 two')
		Assert(model.GetRecord(5).vl_expand?, msg: 'expand 5')
		Assert(model.VirtualListModel_curTop.Pos is: -30)
		Assert(model.Offset is: -29) // -25 - 5

		model.UpdateOffset(3) // scroll down
		Assert(model.Offset is: -26) // -29 + 3
		Assert(model.GetRecord(0).vl_expand?, msg: 'expand 0 one')

		model.SetRecordCollapsed(-3, keepPosition?:) // collapse the expanded record
		Assert(model.Offset is: -21)

		model.SetRecordExpanded(-3, 5) // restore
		Assert(model.Offset is: -26)
		Assert(model.GetRecord(0).vl_expand?, msg: 'expand 0 two')

		model.SetRecordCollapsed(-3) // collapse the expanded record
		Assert(model.Offset is: -24)
		Assert(model.GetRecord(0).num is: 6)
		Assert(model.GetRecord(1).num is: 7)
		}

	Test_closing_cursors_if_all_read()
		{
		_stopLoadAll = true
		// does not add tear down, so we know if test would leave transaction open or not
		model = VirtualListModel(.VirtualList_Table $ ' sort field1')
		model.UpdateVisibleRows(20)
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed is: false)
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed is: false)

		model.UpdateVisibleRows(40) 	// resize
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed,
			msg: 'resize top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed,
			msg: 'resize bottom')

		model = VirtualListModel(.VirtualList_Table $ ' sort field1')
		model.UpdateVisibleRows(20)
		model.UpdateOffset(10) 			// move down 10
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed,
			msg: 'move down 10 top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed,
			msg: 'move down 10 bottom')

		model.UpdateOffset(10) 			// move down 10 again
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed,
			msg: 'move down 10 again top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed,
			msg: 'move down 10 again bottom')

		model = VirtualListModel(.VirtualList_Table $ ' sort field1', startLast:)
		model.UpdateVisibleRows(40)
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed,
			msg: 'start last 40 top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed,
			msg: 'start last 40 bottom')

		model = VirtualListModel(.VirtualList_Table $ ' sort field1', startLast:)
		model.UpdateVisibleRows(20)
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed is: false,
			msg: 'start last 20 top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed is: false,
			msg: 'start last 20 bottom')

		model.UpdateOffset(-20) 			// move up 10
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed,
			msg: 'move up 10 top')
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed,
			msg: 'move up 10 bottom')

		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table $ ' sort field1')
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		model.UpdateOffset(40, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curTop.VirtualListModelCursor_closed is: false)
		Assert(model.VirtualListModel_curBottom.VirtualListModelCursor_closed is: false)
		}

	Test_GetLastVisibleRowIndex()
		{
		model = VirtualListModel(.VirtualList_Table $
			' rename num to test_timestamp
			sort test_timestamp')
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.GetLastVisibleRowIndex() is: 20)

		model.UpdateOffset(25)
		Assert(model.GetLastVisibleRowIndex() is: 29)

		model.UpdateVisibleRows(30) // stay at the end
		Assert(model.GetLastVisibleRowIndex() is: 29)
		}

	Test_GetLastVisibleRowIndex_reverse()
		{
		_stopLoadAll = true
		model = VirtualListModel(.VirtualList_Table $
			' rename num to test_timestamp
			sort test_timestamp', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)

		Assert(model.Offset is: -19)
		Assert(model.GetLastVisibleRowIndex() is: -1)

		model.UpdateVisibleRows(30) // stay at the end
		Assert(model.GetLastVisibleRowIndex() is: -1)
		}

	Test_ReloadRecord()
		{
		model = VirtualListModel(.VirtualList_Table, protectField: 'test',
			observerList: list = Mock())
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)

		rec = model.GetRecord(0)
		Assert(rec.field1 is: '')

		model.AutoSave? = true
		Assert(model.ReloadRecord(rec) isObject:)
		Assert(model.GetRecord(0).field1 is: '')

		rec = model.GetRecord(0)
		QueryDo('update ' $ .VirtualList_Table $ ' where num is 0 set field1 = "f1"')
		Assert(model.ReloadRecord(rec) isObject:)
		Assert(model.GetRecord(0).field1 is: 'f1')

		model = VirtualListModel(.VirtualList_Table, protectField: 'test',
			observerList: list = Mock())
		model.EditModel = edit = Mock()
		edit.When.Editable?().Return(true)
		edit.LockKeyField = 'num'
		edit.ProtectField = 'test'
		edit.When.RecordLocked?([anyArgs:]).Return(true)
		edit.When.GetOutstandingChanges().Return(#())
		list.When.GetColumns([anyArgs:]).Return(#())
		list.Parent = [Window: [Parent: 0]]
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.GetRecord(0).field1 is: 'f1')

		list.When.GetModel().Return(model)
		model.GetRecord(0).field1 = 'new_field1'
		model.AutoSave? = true
		Assert(model.ReloadRecord(rec) isObject:)
		Assert(model.GetRecord(0).field1 is: 'new_field1')

		model.AutoSave = false
		rec = model.GetRecord(0)
		QueryDo('update ' $ .VirtualList_Table $ ' where num is 0 set field1 = "f2"')
		Assert(model.ReloadRecord(rec) isObject:)
		Assert(model.GetRecord(0).field1 is: 'new_field1')

		model = VirtualListModel(.VirtualList_Table, protectField: 'test',
			observerList: list = Mock())
		.AddTeardownModel(model)
		model.AutoSave? = true
		model.UpdateVisibleRows(20)
		rec = model.GetRecord(0)
		QueryDo('update ' $ .VirtualList_Table $ ' where num is 0 set field1 = "f3"')
		Assert(model.ReloadRecord(rec) isObject:)
		Assert(model.GetRecord(0).field1 is: 'f3')

		model.UpdateVisibleRows(20)
		rec = model.GetRecord(0)
		QueryDo('update ' $ .VirtualList_Table $ ' where num is 0
			set num = 100, field1 = "f4"')
		Assert(model.ReloadRecord(rec, force: true, newRec: [num: 100]) isObject:)
		rec2 = model.GetRecord(0)
		Assert(rec2.field1 is: 'f4')

		rec = model.GetRecord(0)
		QueryDo('update ' $ .VirtualList_Table $ ' where num is 100
			set num = 200, field1 = "f5"')
		Assert(model.ReplaceRecord([num: 100], newRec: [num: 200]) isObject:)
		rec2 = model.GetRecord(0)
		Assert(rec2.field1 is: 'f5')
		Assert(rec2.num is: 200)

		// cannot find the new record
		rec = model.GetRecord(0)
		Assert(model.ReplaceRecord([num: 200], newRec: [num: 300]) isString:)
		rec2 = model.GetRecord(0)
		Assert(rec2.field1 is: 'f5')
		Assert(rec2.num is: 200)

		QueryDo('update ' $ .VirtualList_Table $ ' where num is 200 set num = 0')
		}

	Test_expand_collapse_over_limit()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table $ ' sort num')
		.AddTeardownModel(model)
		model.UpdateVisibleRows(5)
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.Offset is: 0)
		Assert(model.VirtualListModel_curBottom.Pos is: 6)

		model.SetRecordExpanded(0, 2)
		Assert(model.Offset is: 0)
		Assert(model.GetRecord(0).num is: 0)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 one')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 one')
		Assert(model.GetRecord(3).num is: 1)
		Assert(model.VirtualListModel_curBottom.Pos is: 8)

		for .. 5
			model.UpdateOffset(2, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_data.Size() lessThanOrEqualTo: 10)
		Assert(model.Offset is: 8) // 2 records collapsed
		Assert(model.GetRecord(0).num is: 8)

		model.SetRecordExpanded(0, 2)
		Assert(model.Offset is: 8)
		Assert(model.GetRecord(0).num is: 8)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 two')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 two')
		Assert(model.GetRecord(3).num is: 9)

		model.SetRecordCollapsed(0)
		Assert(model.VirtualListModel_data.Size() lessThanOrEqualTo: 10)
		Assert(model.Offset is: 8)
		Assert(model.GetRecord(0).num is: 8)
		Assert(model.GetRecord(0).vl_expanded_rows is: '')
		Assert(model.GetRecord(1).num is: 9)
		}

	Test_expand_collapse_over_limit_reverse()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(5)
		model.UpdateOffset(-1, .FakeSaveAndCollapse)
		Assert(model.GetRecord(0).num is: 25)
		Assert(model.Offset is: -5)

		model.SetRecordExpanded(0, 2)
		Assert(model.Offset is: -7)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(0).num is: 25)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 one')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 one')
		Assert(model.GetRecord(3).num is: 26)
		Assert(model.VirtualListModel_curTop.Pos is: -8)

		for .. 5
			model.UpdateOffset(-2, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_data.Size() lessThanOrEqualTo: 10)
		Assert(model.Offset is: -15)
		model.SetRecordExpanded(0, 2)
		Assert(model.Offset is: -17)
		Assert(model.GetRecord(0).vl_expanded_rows is: 2)
		Assert(model.GetRecord(0).num is: 15)
		Assert(model.GetRecord(1).vl_expand?, msg: 'expand 1 two')
		Assert(model.GetRecord(2).vl_expand?, msg: 'expand 2 two')
		Assert(model.GetRecord(3).num is: 16)

		model.SetRecordCollapsed(0)
		Assert(model.Offset is: -15)
		Assert(model.GetRecord(0).vl_expanded_rows is: '')
		Assert(model.GetRecord(0).num is: 15)
		Assert(model.GetRecord(1).num is: 16)
		}

	Test_expand_collapse_over_limit_but_all_read()
		{
		_stopLoadAll = true
		spy = .SpyOn(SuneidoLog)
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 35
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table $ ' sort num')
		.AddTeardownModel(model)

		// all 30 records are read, cursors should be closed
		model.UpdateVisibleRows(30)
		Assert(model.AllRead?, msg: 'all read?')

		// expand a record, so now there are 40 lines which exceeds the limit
		model.SetRecordExpanded(0, 10)
		Assert(model.VirtualListModel_curBottom.Pos is: 40)

		// scrolling down should not trigger recyling
		// because all data is read and cursors are closed
		model.UpdateOffset(15, .FakeSaveAndCollapse)
		Assert(spy.CallLogs() isSize: 0)
		}

	Test_ReadAllData()
		{
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			}
		model = modelClass(.VirtualList_Table $ ' sort num')
		model.ReadAllData()
		Assert(model.GetLoadedData() isSize: 30)
		for i in ..30
			Assert(model.GetRecord(i).num is: i)
		}

	Test_Dependancies_On_Reload()
		{
		.SpyOn(SelectFields.SelectFields_warnIfNoPrompt).Return('')
		table = .MakeTable(
			'(num, test_field_a, test_field_c, test_field_c_deps) key (num)')
		.MakeLibraryRecord([name: 'Rule_test_field_b', text: 'function()
			{
			return .test_field_a + 1
			}'])
		.MakeLibraryRecord([name: 'Rule_test_field_c', text: 'function()
			{
			return .test_field_b + 2
			}'])
		rec = [num: 1, test_field_a: 100]
		rec.test_field_b
		rec.test_field_c
		QueryOutput(table, rec)
		model = VirtualListModel(table $ ' extend test_field_b',
			columns: #(test_field_a, test_field_b)
			protectField: 'editable')
		model.AutoSave? = true
		model.ReadAllData()
		rec = model.GetRecord(0)
		freshRec = model.ReloadRecord(rec)
		m = Mock()
		m.When.GetModel().Return(model)
		m.When.GetColumns([anyArgs:]).Return(model.Columns())
		m.Parent = [Window: [Hwnd: 0]]
		freshRec.vl_list = m
		freshRec.test_field_a = 200
		QueryApply1(table)
			{
			it.Update(freshRec)
			}
		Assert(QueryFirst(table $ ' sort num').test_field_c is: 203)
		}

	Test_Seeking()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
		model = modelClass(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(10)
		model.Seek('num', 999)
		Assert(model.GetRecord(0).num is: 20)
		Assert(model.GetRecord(1).num is: 21)
		Assert(model.GetRecord(9).num is: 29)
		model.UpdateVisibleRows(20)
		Assert(model.GetRecord(0).num is: 10)
		Assert(model.GetRecord(19).num is: 29)
		}

	Test_reading_twice()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			VirtualListModel_closeCursorsIfAllRead() {}
			}
		model = modelClass(.VirtualList_Table $ ' sort num')
		.AddTeardownModel(model)
		model.UpdateVisibleRows(5)
		Assert(model.VirtualListModel_curTop.Pos is: 0)
		Assert(model.VirtualListModel_curBottom.Pos is: 5)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// top cursor move down to release
		model.UpdateOffset(10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curTop.Pos is: 3)
		Assert(model.VirtualListModel_curBottom.Pos is: 15)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// top cursor move up to read
		model.UpdateOffset(-10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curTop.Pos is: 0)
		Assert(model.VirtualListModel_curBottom.Pos is: 12)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// top cursor move down to release
		model.UpdateOffset(10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curTop.Pos is: 3)
		Assert(model.VirtualListModel_curBottom.Pos is: 15)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)
		}

	Test_reading_twice_with_startLast()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			VirtualListModel_closeCursorsIfAllRead() {}
			}
		model = modelClass(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(5)
		Assert(model.VirtualListModel_curBottom.Pos is: 0)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move up to release
		model.UpdateOffset(-10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: -3)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move down to read
		model.UpdateOffset(10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: 0)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move up to release
		model.UpdateOffset(-10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: -3)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)
		}

	Test_reading_twice_with_seeking()
		{
		_stopLoadAll = true
		modelClass = VirtualListModel
			{
			VirtualListModel_limit: 10
			VirtualListModel_segment: 3
			VirtualListModel_closeCursorsIfAllRead() {}
			}
		model = modelClass(.VirtualList_Table $ ' sort num', startLast:)
		.AddTeardownModel(model)

		model.UpdateVisibleRows(5)
		Assert(model.VirtualListModel_curBottom.Pos is: 0)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		model.Seek('num', 10)
		Assert(model.VirtualListModel_curBottom.Pos is: 5)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move up to release
		model.UpdateOffset(-10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: 2)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move down to read
		model.UpdateOffset(10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: 5)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)

		// bottom cursor move up to release
		model.UpdateOffset(-10, .FakeSaveAndCollapse)
		Assert(model.VirtualListModel_curBottom.Pos is: 2)
		Assert(model.GetLoadedData() hasntMember: model.VirtualListModel_curBottom.Pos)
		}

	Test_GetKeyQuery()
		{
		// no key
		table = .MakeTable('(a, b, c) key ()')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table)

		model = VirtualListModel(table $ ' rename a to a2')
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table $ ' rename a to a2')

		// one key
		table = .MakeTable('(a, b, c) key (a)')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table $ ' where a is 1')

		model = VirtualListModel(table $ ' rename a to a2')
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a2: 1, b: 1, c: 1])
			is: table $ ' rename a to a2 where a2 is 1')

		// two keys
		table = .MakeTable('(a, b, c) key (a) key (b, c)')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table $ ' where a is 1')

		table = .MakeTable('(a, b, c) key (a) key (b, c)')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table $ ' where a is 1')

		model = VirtualListModel(table $ ' rename a to a2')
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a2: 1, b: 1, c: 1])
			is: table $ ' rename a to a2 where a2 is 1')

		table = .MakeTable('(a, b, c) key (b, c) key (a)')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1]) is: table $ ' where a is 1')

		model = VirtualListModel(table $ ' rename a to a2')
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a2: 1, b: 1, c: 1])
			is: table $ ' rename a to a2 where a2 is 1')

		// composite keys
		table = .MakeTable('(a, b, c) key (b, c)')
		model = VirtualListModel(table)
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a: 1, b: 1, c: 1])
			is: table $ ' where b is 1 and c is 1')

		model = VirtualListModel(table $ ' rename b to b2')
		.AddTeardownModel(model)
		Assert(model.GetKeyQuery([a2: 1, b2: 1, c: 1])
			is: table $ ' rename b to b2 where b2 is 1 and c is 1')
		}

	Test_ValidateRow()
		{
		mock = Mock(VirtualListModel)
		mock.When.ValidateRow([anyArgs:]).CallThrough()
		mock.VirtualListModel_data = []
		Assert(mock.ValidateRow(0) is: false)

		mock.VirtualListModel_data = [0, 1, 2, 3, 4, 5]
		mock.VirtualListModel_curTop = Object(Pos: 0)
		mock.VirtualListModel_curBottom = Object(Pos: 6)
		Assert(mock.ValidateRow(-1) is: false)
		Assert(mock.ValidateRow(0) is: 0)
		Assert(mock.ValidateRow(3) is: 3)
		Assert(mock.ValidateRow(5) is: 5)
		Assert(mock.ValidateRow(6) is: false)
		Assert(mock.ValidateRow(-1, true) is: 0)
		Assert(mock.ValidateRow(0, true) is: 0)
		Assert(mock.ValidateRow(3, true) is: 3)
		Assert(mock.ValidateRow(5, true) is: 5)
		Assert(mock.ValidateRow(6, true) is: 5)
		}

	Test_recPos()
		{
		mock = Mock(VirtualListModel)
		mock.VirtualListModel_data = data = Object()
		mock.VirtualListModel_baseQuery = 'valid'
		mock.When.keys('valid').Return(Object('a', 'b'))
		mock.When.keys('invalid').Throw('Forced Error')
		mock.When.recPos([anyArgs:]).CallThrough()
		mock.When.logDebugging([anyArgs:]).Do({ })

		keyPairs = Object(a: Timestamp(), b: Timestamp())
		rec = [a: keyPairs.a, b: keyPairs.b, c: 'Other']

		// data is empty, rec is not found
		Assert(mock.recPos(rec) is: false)
		mock.Verify.logDebugging(reason: 'Record no longer exists in .data', :keyPairs)

		// data has values now, rec is still not found
		data.Append([
			[a: 'rec1', b: 'rec1', c: 'Other'],
			[a: 'rec2', b: 'rec2', c: 'Other'],
			[a: 'rec3', b: 'rec3', c: 'Other']])
		Assert(mock.recPos(rec) is: false)
		mock.Verify.Times(2).
			logDebugging(reason: 'Record no longer exists in .data', :keyPairs)

		// data has rec now, pos is returned
		data.Add(dataRec = rec.Copy(), at: 2)
		Assert(mock.recPos(rec) is: 2)

		// rec no longer matches data, selected rec has an extra member
		rec.vl_member = ts = Timestamp()
		Assert(mock.recPos(rec) is: false)
		mock.Verify.logDebugging(reason: 'Record out of sync',
			diffs: Object(vl_member: 'dataVal: "", recVal: ' $ Display(ts)))

		// rec still doesn't match data, data now has a different value
		rec.Delete('vl_member')
		dataRec.c = 'Different'
		Assert(mock.recPos(rec) is: false)
		mock.Verify.logDebugging(reason: 'Record out of sync',
			diffs: Object(c: 'dataVal: "Different", recVal: "Other"'))

		// several differences between rec and data are detected
		dataRec.extra_column = 'Extra Column'
		rec.vl_other_member = 'Other Member'
		Assert(mock.recPos(rec) is: false)
		mock.Verify.logDebugging(reason: 'Record out of sync',
			diffs: Object(
				c: 'dataVal: "Different", recVal: "Other"',
				extra_column: 'dataVal: "Extra Column", recVal: ""',
				vl_other_member: 'dataVal: "", recVal: "Other Member"'))

		// Test error handling
		mock.VirtualListModel_baseQuery = 'invalid'
		Assert(mock.recPos(rec) is: false)
		mock.Verify.logDebugging(reason: 'Error caught', error: 'Forced Error')
		}
	}
