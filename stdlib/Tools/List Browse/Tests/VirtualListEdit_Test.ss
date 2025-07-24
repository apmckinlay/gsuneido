// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_nextColEditable()
		{
		mock = Mock(VirtualListEdit)
		mock.When.nextColEditable([anyArgs:]).CallThrough()
		mock.VirtualListEdit_model = Object(ColModel: colModel = Mock())
		mock.VirtualListEdit_minWidth = VirtualListEdit.VirtualListEdit_minWidth
		colModel.When.GetCustomFields().Return(#(f3: #(tabover:)))

		mock.VirtualListEdit_model.EditModel = VirtualListEditModel()
		Assert(mock.nextColEditable('f0', 0, [vl_deleted:]) is: false)

		// column width < 5
		colModel.When.GetColWidth([anyArgs:]).Return(5)
		Assert(mock.nextColEditable('f1', 1, []) is: false)

		colModel.When.GetColWidth([anyArgs:]).Return(100)
		// no protect field
		mock.VirtualListEdit_model.EditModel = VirtualListEditModel()
		Assert(mock.nextColEditable('f2', 2, []) is: false)

		// tab over
		mock.VirtualListEdit_model.EditModel.ProtectField = 'something'
		Assert(mock.nextColEditable('f3', 3, []) is: false)

		// not allow editing
		mock.When.allowEdit([anyArgs:]).Return(false)
		Assert(mock.nextColEditable('f4', 4, []) is: false)

		mock.When.allowEdit([anyArgs:]).Return('not editable')
		Assert(mock.nextColEditable('f5', 5, []) is: false)

		// allow editing
		mock.When.allowEdit([anyArgs:]).Return(true)
		Assert(mock.nextColEditable('f5', 5, []))

		mock.When.allowEdit([anyArgs:]).Return('')
		Assert(mock.nextColEditable('f5', 5, []))
		}

	Test_editNextCell()
		{
		mock = Mock(VirtualListEdit)
		mock.When.editNextCell([anyArgs:]).CallThrough()
		mock.When.editCell([anyArgs:]).Return('')

		mock.VirtualListEdit_model = Mock(VirtualListModel)
		mock.VirtualListEdit_rec = []
		colModel = mock.VirtualListEdit_model.ColModel = Mock(VirtualListColModel)
		grid = mock.VirtualListEdit_parent = Mock()
		grid.When.Send("VirtualListGrid_SaveRecord", []).Return(true)
		grid.When.InsertRow([anyArgs:]).Return(true)

		mock.When.nextColEditable('b', 1, []).Return(true)

		grid.When.GetCurrentRowCellRect([anyArgs:]).Return('rect')
		grid.When.GetSelectedRecord().Return([])

		colModel.When.GetColumns().Return(#(a, b, c))
		colModel.When.FindCol('a').Return(0)
		colModel.When.Get(1).Return('b')

		mock.editNextCell('a', 1)
		mock.Verify.editCell('b', 'rect')

		colModel.When.FindCol('b').Return(1)
		colModel.When.Get(2).Return('c')
		mock.When.nextColEditable('c', 2, []).Return(true)
		mock.editNextCell('b', 1)
		mock.Verify.editCell('c', 'rect')

		colModel.When.FindCol('c').Return(2)
		colModel.When.Get(0).Return('a')
		grid.When.SelectNextRow(1).Return(true)
		mock.When.nextColEditable('a', 0, []).Return(true)
		mock.When.rowChanged([anyArgs:]).CallThrough()
		mock.editNextCell('c', 1)
		mock.Verify.editCell('a', 'rect')

		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('b', 'rect')

		mock.editNextCell('b', -1)
		mock.Verify.Times(2).editCell('a', 'rect')

		grid.When.SelectNextRow(-1).Return(true)
		mock.When.nextColEditable('a', 0, []).Return(false)
		mock.When.nextColEditable('b', 1, []).Return(false)
		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('c', 'rect')

		grid.When.SelectNextRow(-1).Return(false)
		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('c', 'rect') // no change

		grid.When.SelectNextRow(1).Return(false)
		mock.editNextCell('c', 1)
		mock.Verify.Times(2).editCell('c', 'rect') // no change
		}

	Test_valueChanged?()
		{
		fn = VirtualListEdit.VirtualListEdit_valueChanged?
		mock = Mock()
		mock.VirtualListEdit_rec = rec = []
		mock.VirtualListEdit_col = 'field'
		Assert(mock.Eval(fn, valid?: true, unvalidated_val: '', data: 'f'))

		rec.field = 'f'
		Assert(mock.Eval(fn, valid?: true, unvalidated_val: '', data: 'f') is: false)

		rec.field = 10
		Assert(mock.Eval(fn, valid?: true, unvalidated_val: '', data: 'f'))

		rec.field = ''
		Assert(mock.Eval(fn, valid?: false, unvalidated_val: 'ff', data: ''))

		rec.field = ''
		rec.List_InvalidData = Object(field: 'ff')
		Assert(mock.Eval(fn, valid?: false, unvalidated_val: 'ff', data: '') is: false)

		rec.field = ''
		rec.List_InvalidData = Object(field: 'ff')
		Assert(mock.Eval(fn, valid?:, unvalidated_val: '', data: 20))
		}

	Test_ListEditWindow_Commit()
		{
		mock = Mock(VirtualListEdit)
		mock.VirtualListEdit_col = 'field1'
		mock.VirtualListEdit_rec = ['field1': 'hello']
		mock.When.valueChanged?([anyArgs:]).CallThrough()
		mock.When.commitChange([anyArgs:]).CallThrough()

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 0,
			data: 'hello world', valid?:, unvalidated_val: '', readonly:, dirty?: false)
		mock.Verify.Never().commitChange([anyArgs:])
		mock.Verify.Never().editNextCell([anyArgs:])

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: [], valid?:, unvalidated_val: '', readonly:)
		mock.Verify.Never().commitChange([anyArgs:])
		mock.Verify.editNextCell('field1', 1)

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello', valid?:, unvalidated_val: '', readonly: false, dirty?: false)
		mock.Verify.Never().commitChange([anyArgs:])
		mock.Verify.Times(2).editNextCell('field1', 1)

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello', valid?:, unvalidated_val: '', readonly: false, dirty?:)
		mock.Verify.Never().commitChange([anyArgs:])
		mock.Verify.Times(3).editNextCell('field1', 1)

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello', valid?:, unvalidated_val: '', readonly:, dirty?:)
		mock.Verify.Never().commitChange([anyArgs:])
		mock.Verify.Times(4).editNextCell('field1', 1)
		}

	Test_changes_and_save()
		{
		.SpyOn(SelectFields.SelectFields_warnIfNoPrompt).Return('')
		mock = Mock(VirtualListEdit)
		mock.VirtualListEdit_col = 'field1'
		mock.When.valueChanged?([anyArgs:]).CallThrough()
		mock.When.commitChange([anyArgs:]).CallThrough()

		list = Mock()
		list.When.GetColumns().Return(#(field1))
		list.When.Send([anyArgs:]).Return(0)
		list.Parent = [Window: [Hwnd: 0]]

		model = VirtualListModel(.VirtualList_Table, columns: #(field1),
			protectField: 'editable', observerList: list)
		model.UpdateVisibleRows(10)
		mock.VirtualListEdit_rec = rec = model.GetRecord(0)
		.AddTeardownModel(model)

		list.When.GetModel().Return(model)
		mock.VirtualListEdit_parent = grid = Mock()
		grid.Controller = list
		grid.Hwnd = 0

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello world', valid?:, unvalidated_val: '', readonly: false, dirty?:)
		mock.Verify.commitChange([anyArgs:])
		mock.Verify.Times(1).editNextCell('field1', 1)
		grid.Verify.RepaintSelectedRows()
		Assert(rec.field1 is: 'hello world')
		Assert(model.EditModel.RecordChanged?(rec))
		Assert(model.EditModel.HasChanges?())
		Assert(model.EditModel.AllowLeaving?() is: false)
		Assert(model.EditModel.RecordDirty?(rec))
		Assert(model.EditModel.GetOutstandingChanges() is: Object(rec))
		Assert(model.EditModel.ColumnInvalid?(rec, 'field1') is: false)

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello wrong world', valid?: false,
			unvalidated_val: '', readonly: false, dirty?:)
		mock.Verify.Times(2).commitChange([anyArgs:])
		mock.Verify.Times(2).editNextCell('field1', 1)
		grid.Verify.Times(2).RepaintSelectedRows()
		Assert(rec.field1 is: 'hello wrong world')
		Assert(model.EditModel.RecordChanged?(rec))
		Assert(model.EditModel.HasChanges?())
		Assert(model.EditModel.AllowLeaving?() is: false)
		Assert(model.EditModel.RecordDirty?(rec))
		Assert(model.EditModel.GetOutstandingChanges() is: Object(rec))
		Assert(model.EditModel.ColumnInvalid?(rec, 'field1'))

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: '', valid?: true, unvalidated_val: '', readonly: false, dirty?:)
		mock.Verify.Times(3).commitChange([anyArgs:])
		mock.Verify.Times(3).editNextCell('field1', 1)
		grid.Verify.Times(3).RepaintSelectedRows()
		Assert(rec.field1 is: '')
		Assert(model.EditModel.RecordChanged?(rec) is: false)
		Assert(model.EditModel.HasChanges?() is: false)
		Assert(model.EditModel.AllowLeaving?())
		Assert(model.EditModel.RecordDirty?(rec) is: false)
		Assert(model.EditModel.GetOutstandingChanges() is: Object())
		Assert(model.EditModel.ColumnInvalid?(rec, 'field1') is: false)

		rec.field2 = 'touched' // not table column
		Assert(model.EditModel.RecordChanged?(rec) is: false)
		Assert(model.EditModel.HasChanges?() is: false)
		Assert(model.EditModel.AllowLeaving?())
		Assert(model.EditModel.RecordDirty?(rec) is: false)
		Assert(model.EditModel.GetOutstandingChanges() is: Object())
		Assert(model.EditModel.ColumnInvalid?(rec, 'field1') is: false)

		view = Mock(Addon_VirtualListView_Edit)
		view.Model = model
		Assert(view.Eval(Addon_VirtualListView_Edit.VirtualListGrid_SaveRecord, rec)
			isnt: false)

		Assert(rec.field1 is: '')
		Assert(model.EditModel.RecordChanged?(rec) is: false)
		Assert(model.EditModel.HasChanges?() is: false)
		Assert(model.EditModel.AllowLeaving?())
		Assert(model.EditModel.RecordDirty?(rec) is: false)
		Assert(model.EditModel.GetOutstandingChanges() is: Object())
		Assert(model.EditModel.ColumnInvalid?(rec, 'field1') is: false)
		newrec = model.GetRecord(0)
		Assert(newrec.field1 is: '')
		Assert(newrec.field2 is: '')

		view.When.Send([anyArgs:]).Return(true)
		view.Grid = grid = Mock()
		view.Addons = Mock()
		view.Addons.When.Collect([anyArgs:]).Return(Object())
		grid.Controller = view
		grid.Window = Object(Hwnd: 123)
		view.When.GetContextMenu().Return(
			Object(	UpdateHistory: function(@unused) { } ))
		mock.VirtualListEdit_rec = rec = model.GetRecord(0)

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello wrong world', valid?: false,
			unvalidated_val: '', readonly: false, dirty?:)
		Assert(view.Eval(Addon_VirtualListView_Edit.VirtualListGrid_SaveRecord, rec)
			is: false)
		Assert(Query1(.VirtualList_Table, num: 0).field1 is: '')

		mock.Eval(VirtualListEdit.ListEditWindow_Commit, 'field1', 'row', dir: 1,
			data: 'hello world', valid?:, unvalidated_val: '', readonly: false, dirty?:)
		Assert(view.Eval(Addon_VirtualListView_Edit.VirtualListGrid_SaveRecord, rec)
			isnt: false)
		Assert(model.GetRecord(0).field1 is: 'hello world')
		Assert(Query1(.VirtualList_Table, num: 0).field1 is: 'hello world')
		}
	}