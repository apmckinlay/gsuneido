// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_VirtualListExpandModel_Test
	{
	recCtrl: RecordControl
		{
		SetProtectField(@unused) { }
		Set(@unused) { }
		SetReadOnly(@unused) { }
		SetVisible(@unused) { }
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
		grid.When.ActWith([anyArgs:]).Do({ |call| (call.block)() })
		grid.When.Act([anyArgs:])
		constructedCtrl.Parent = grid
		constructedCtrl.UniqueId = 0
		constructedCtrl.Hwnd = 0
		constructedCtrl.When.GetLayout().Return(Object('VirtualList'))
		constructedCtrl.When.GetControl().Return(
			VirtualListControl,
			VirtualListControl,
			VirtualListControl, // for destroy
			VirtualListControl, // for destroy
			.recCtrl
			)
		model = Object(ColModel: Mock())
		model.Offset = 0
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
		grid = Mock()
		grid.When.Act([anyArgs:])

		instance = class{}()
		layoutOb = [ctrl: Mock(), rows: 3]
		layoutOb.ctrl.UniqueId = 0
		layoutOb.ctrl.When.GetControl().Return(instance)
		expandModel.Expand(old = [status: 'old'], layoutOb, model: false)
		expandModel.SetExpandRecord(newRec = [status: 'new'], old)
		Assert(expandModel.VirtualListExpandModel_expandedRows[0].rec is: newRec)

		expandModel.Collapse(newRec, grid)
		Assert(expandModel.VirtualListExpandModel_expandedRows isSize: 0)
		expandModel.SetExpandRecord(newRec, newRec)
		grid.Verify.Act('VirtualListExpand_Destroy', 0)
		}

	Test_main()
		{
		model = VirtualListExpandModel()
		grid = Mock()
		grid.When.Act([anyArgs:])

		Assert(model.GetExpandedControl([]) is: false)
		model.DestroyAll()

		instance = class{}()
		layoutOb = [ctrl: mock = Mock(), rows: 3]
		layoutOb.ctrl.UniqueId = 0
		layoutOb.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '1'], layoutOb, model: false)
		Assert(model.GetExpandedControl([value: '1']) is: layoutOb)

		layoutOb2 = [ctrl: mock2 = Mock(), rows: 4]
		mock2.Index = 2
		layoutOb2.ctrl.When.GetControl().Return(instance)
		model.Expand([value: '2'], layoutOb2, model: false)
		Assert(model.GetExpandedControl([value: '2']) is: layoutOb2)

		model.Collapse([value: '1'], grid)
		Assert(model.GetExpandedControl([value: '1']) is: false)
		mock.Verify.Destroy()
		mock2.Verify.Never().Destroy()

		Assert(model.GetExpandedControl([value: '2']) is: layoutOb2)
		model.DestroyAll()
		mock.Verify.Destroy()
		mock2.Verify.Destroy()
		Assert(model.GetExpandedControl([value: '2']) is: false)
		}

	Test_FindIfOverLimit() { }
	Test_find_hwnd() { }
	}