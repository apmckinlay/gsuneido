// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		model = VirtualListExpandModel()
		Assert(model.GetExpandedControl([]) is: false)
		model.DestroyAll()

		instance = class{}()
		layoutOb = [ctrl: mock = Mock(), rows: 3]
		layoutOb.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '1'], layoutOb, model: false)
		Assert(model.GetExpandedControl([value: '1']) is: layoutOb)

		layoutOb2 = [ctrl: mock2 = Mock(), rows: 4]
		mock2.Index = 2
		layoutOb2.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '2'], layoutOb2, model: false)
		Assert(model.GetExpandedControl([value: '2']) is: layoutOb2)

		model.Collapse([value: '1'])
		Assert(model.GetExpandedControl([value: '1']) is: false)
		mock.Verify.Destroy()
		mock2.Verify.Never().Destroy()

		Assert(model.GetExpandedControl([value: '2']) is: layoutOb2)
		model.DestroyAll()
		mock.Verify.Destroy()
		mock2.Verify.Destroy()
		Assert(model.GetExpandedControl([value: '2']) is: false)
		}

	Test_ClearAllSelections()
		{
		instance = class{}()
		model = VirtualListExpandModel()
		layoutOb2 = [ctrl: mock2 = Mock(), rows: 3]
		mock2.Index = 2
		layoutOb2.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '2'], layoutOb2, model: false)

		layoutOb3 = [ctrl: mock3 = Mock(), rows: 4]
		layoutOb3.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '3'], layoutOb3, model: false)
		mock3.Index = 3

		mock3.ClearSelect = true
		model.ClearAllSelections(mock2)
		mock2.Verify.Never().ClearSelect()
		mock3.Verify.ClearSelect()

		mock2.ClearSelect = true
		model.ClearAllSelections()
		mock2.Verify.ClearSelect()
		mock3.Verify.Times(2).ClearSelect()
		}

	Test_SetExpandReadOnly()
		{
		model = VirtualListExpandModel()
		model.SetExpandReadOnly([], readonly:)

		instance = class{}()
		layoutOb = [ctrl: ctrl = Mock(), rows: 3]
		layoutOb.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '2'], layoutOb, model: false)

		ctrl.When.GetControl().Return(instance)
		model.SetExpandReadOnly([value: '2'], readonly:)

		_recCtrl = Object()
		ctrl.When.GetControl().Return(RecordControl
			{
			SetReadOnly(readonly)
				{
				_recCtrl.readonly = readonly
				}
			FindControl(unused)
				{
				return _recCtrl.edit = Mock()
				}
			})
		model.SetExpandReadOnly([value: '2'], readonly:)
		Assert(_recCtrl.readonly is: true)
		_recCtrl.edit.When.Pushed?(false)

		model.SetExpandReadOnly([value: '2'], readonly: false)
		Assert(_recCtrl.readonly is: false)
		_recCtrl.edit.When.Pushed?(true)
		}

	recCtrl: RecordControl
		{
		SetProtectField(@unused) { }
		Set(@unused) { }
		SetReadOnly(@unused) { }
		}
	Test_RecycledExpands()
		{
		expandModel = VirtualListExpandModel
			{
			VirtualListExpandModel_updateTabIndex(@unused)
				{
				}
			}()
		grid = Mock()
		grid.When.Construct([anyArgs:]).Return(constructedCtrl = Mock())
		constructedCtrl.Ymin = 100
		constructedCtrl.Hwnd = 0
		constructedCtrl.When.GetControl().Return(
			VirtualListControl,
			VirtualListControl,
			VirtualListControl, // for destroy
			VirtualListControl, // for destroy
			.recCtrl
			)
		model = Object(ColModel: Mock())
		model.ColModel.Offset = 2
		model.ColModel.When.GetCustomFields().Return(#())
		model.EditModel = Object(ProtectField: 'protect')
		expandModel.ConstructAt(
			layoutOb = Object(ctrl: Object('VirtuaList')),
			1, grid, model, rowHeight: 20)
		grid.Verify.Construct([anyArgs:])
		expandModel.Expand([value: '1'], layoutOb, model, readOnly?:)

		expandModel.ConstructAt(
			layoutOb = Object(ctrl: Object('VirtuaList')),
			1, grid, model, rowHeight: 20)
		grid.Verify.Times(2).Construct([anyArgs:])
		expandModel.Expand([value: '2'], layoutOb, model, readOnly?:)

		constructedCtrl.Verify.Times(2).GetControl()

		expandModel.CollapseAll()
		constructedCtrl.Verify.Times(2).Destroy() // 2 virtual lists

		expandModel.ConstructAt(
			layoutOb = Object(ctrl: Object('Record')),
			1, grid, model, rowHeight: 20)
		grid.Verify.Times(3).Construct([anyArgs:])
		expandModel.Expand([value: '3'], layoutOb, model, readOnly?:)

		expandModel.ConstructAt(
			layoutOb = Object(ctrl: Object('Record')),
			1, grid, model, rowHeight: 20)
		grid.Verify.Times(4).Construct([anyArgs:])
		expandModel.Expand([value: '4'], layoutOb, model, readOnly?:)

		// should collapse all
		Assert(expandModel.RecycleExpands() isSize: 2)
		constructedCtrl.Verify.Times(2).Destroy() // still 2 virtual lists

		expandModel.ConstructAt(
			layoutOb = Object(ctrl: Object('Record')),
			1, grid, model, rowHeight: 20)
		expandModel.Expand([value: '5'], layoutOb, model, readOnly?:)
		// Construct should be called and still 4 times
		grid.Verify.Times(4).Construct([anyArgs:])
		Assert(expandModel.VirtualListExpandModel_recycledExpands isSize: 1)
		Assert(expandModel.VirtualListExpandModel_expandedRows isSize: 1)
		Assert(expandModel.GetControls() isSize: 1)

		// should collapse all and destroy all controls
		expandModel.DestroyAll()
		constructedCtrl.Verify.Times(4).Destroy() // 4 controls in totoal
		Assert(expandModel.VirtualListExpandModel_recycledExpands isSize: 0)
		Assert(expandModel.VirtualListExpandModel_expandedRows isSize: 0)
		}

	Test_SetExpandRecord()
		{
		expandModel = VirtualListExpandModel()
		expandModel.SetExpandRecord([status: 'new'], [status: 'old'])

		instance = class{}()
		layoutOb = [ctrl: Mock(), rows: 3]
		layoutOb.ctrl.When.GetControl().Return(instance)
		expandModel.Expand(old = [status: 'old'], layoutOb, model: false)
		expandModel.SetExpandRecord(newRec = [status: 'new'], old)
		Assert(expandModel.VirtualListExpandModel_expandedRows[0].rec is: newRec)

		expandModel.Collapse(newRec)
		Assert(expandModel.VirtualListExpandModel_expandedRows isSize: 0)
		expandModel.SetExpandRecord(newRec, newRec)
		}

	model: class
		{
		GetRecord(index)
			{
			switch(index)
				{
			case 0: return [rec: 0]
			case 1: return [rec: 1, vl_expanded_rows: 1]
			case 2: return []
			case 3: return [rec: 2]
			case 4: return [rec: 3, vl_expanded_rows: 1]
			case 5: return []
			default: return false
				}
			}
		}
	Test_find_hwnd()
		{
		expandModel = VirtualListExpandModel()
		expandModel.VirtualListExpandModel_expandedRows = Object(
			[rec: [rec: 1, vl_expanded_rows: 1], layout: [ctrl: [Hwnd: 1]]],
			[rec: [rec: 3, vl_expanded_rows: 1], layout: [ctrl: [Hwnd: 3]]],
			)
		expandModel.VirtualListExpandModel_expandedRows
		find = expandModel.VirtualListExpandModel_findPreviousHwnd
		Assert(find(.model, 0, #(Hwnd: 'grid')) is: 'grid')
		Assert(find(.model, 2, #()) is: 1)

		find = expandModel.VirtualListExpandModel_findAfterHwnd
		Assert(find(.model, 2) is: 3)
		}

	Test_FindIfOverLimit()
		{
		expandModel = VirtualListExpandModel
			{
			New()
				{
				.VirtualListExpandModel_expandedRows = #(#(num: 1), #(num: 3), #(num: 5),
				#(num: 7), #(num: 9))
				}
			VirtualListExpandModel_expandLimit: 4
			GetExpanded()
				{
				return .VirtualListExpandModel_expandedRows
				}
			}
		model = expandModel()
		getRowNumFn = function (rec) { return rec.num }

		// over limit
		Assert(model.FindIfOverLimit([num: 7], getRowNumFn) is: 1)
		Assert(model.FindIfOverLimit([num: 1], getRowNumFn) is: 9)
		Assert(model.FindIfOverLimit([num: 9], getRowNumFn) is: 1)

		// on the limit
		model.VirtualListExpandModel_expandLimit = 5
		Assert(model.FindIfOverLimit([num: 9], getRowNumFn) is: false)

		// under limit
		model.VirtualListExpandModel_expandLimit = 6
		Assert(model.FindIfOverLimit([num: 9], getRowNumFn) is: false)
		}

	Test_CustomizableExpand?()
		{
		fn = VirtualListExpandModel.CustomizableExpand?
		layoutOb = #(ctrl: ())
		Assert(fn(layoutOb) is: false)

		layoutOb = #(ctrl: (Customizable))
		Assert(fn(layoutOb) is: false)

		layoutOb = #(ctrl: (Customizable, tabName: 'test_tab'))
		Assert(fn(layoutOb) is: false)

		layoutOb = Object(ctrl: Object('Customizable',
			tabName: CustomizeExpandControl.LayoutName))
		Assert(fn(layoutOb) is: true)

		layoutOb = #(ctrl: (Customizable, tabName: 'CustomizableExpand'))
		Assert(fn(layoutOb) is: true)

		layoutOb = #(ctrl: (Vert, Form,
			(Customizable, tabName: CustomizableExpand)))
		Assert(fn(layoutOb) is: true)

		layoutOb = #(ctrl: (Customizable, test_table, test_tab, CustomizableExpand))
		Assert(fn(layoutOb) is: true)

		layoutOb = #(ctrl: (Form,
			Customizable, nl,
			(Customizable, tabName: 'CustomizableExpand')))
		Assert(fn(layoutOb) is: true)
		}
	}
