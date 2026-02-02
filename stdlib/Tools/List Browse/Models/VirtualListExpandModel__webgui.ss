// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{
		.expandedRows = Object()
		.recycledExpands = Object()
		}

	CustomizableLayout()
		{
		return Object(ctrl: Object('Record',
			Object('Customizable', tabName: CustomizeExpandControl.LayoutName,
				user: Suneido.User)))
		}

	ConstructAt(layoutOb, rowIndex, grid, model, rowHeight/*unused*/)
		{
		layoutOb.ctrl = layoutOb.ctrl.Copy()
		recordExpand? = layoutOb.ctrl[0] is 'Record'
		if recordExpand?
			layoutOb.ctrl.custom = model.ColModel.GetCustomFields()

		grid.ActWith()
			{
			reuse? = true
			if false is ctrl = .reuseRecycledExpand(recordExpand?)
				{
				reuse? = false
				ctrl = grid.Construct(
					Object('WndPane', layoutOb.ctrl, windowClass: "SuBtnfaceArrow"))
				}

			rows = layoutOb.GetDefault('rows', false)
			layoutOb.flex? = not layoutOb.Member?('rows')
			layoutOb.rows = 0
			DoStartup(ctrl)
			ctrl.SetVisible(true) // for recycled expands
			layoutOb.ctrl = ctrl

			Object('VirtualListExpand_ContructAt', rowIndex + model.Offset,
				reuse? is true ? ctrl.UniqueId : ctrl.GetLayout(), :rows)
			}

		}

	Expand(rec, layout, model, readOnly? = false)
		{
		child = layout.ctrl.GetControl()
		if child.Base?(RecordControl)
			{
			child.SetProtectField(model.EditModel.ProtectField)
			child.Set(rec)
			child.SetReadOnly(readOnly?, layout.ctrl.Hwnd)
			if not readOnly?
				child.Valid(forceCheck:)
			}

		Assert(layout hasMember: 'ctrl')
		Assert(layout hasMember: 'rows')
		.expandedRows.Add([:rec, :layout])
		}

	Collapse(rec, grid)
		{
		i = .expandedRows.FindIf({ it.rec is rec })
		.destroy(.expandedRows[i].layout, grid)
		.expandedRows.Delete(i)
		}

	CollapseAll()
		{
		for row in .expandedRows
			.destroy(row.layout)
		.expandedRows = Object()
		}

	ClearAllSelections(except = false)
		{
		for row in .expandedRows
			{
			layout = row.layout
			if except isnt false and except is layout.ctrl
				continue
			if layout.ctrl.Member?('ClearSelect')
				layout.ctrl.ClearSelect()
			}
		}

	SetReadOnly(readOnly)
		{
		for row in .expandedRows
			row.layout.ctrl.SetReadOnly(readOnly)
		}

	GetCurrentFocusedRecord(focusHwnd)
		{
		for row in .expandedRows
			{
			ctrl = row.layout.ctrl
			if .isChild(ctrl, focusHwnd)
				{
				recordCtrl = ctrl.GetControl()
				if recordCtrl.Base?(RecordControl)
					return recordCtrl.Get()
				}
			}
		return false
		}

	isChild(ctrl, targetHwnd)
		{
		if ctrl.Hwnd is targetHwnd
			return true
		for c in ctrl.GetChildren()
			if .isChild(c, targetHwnd)
				return true
		return false
		}

	GetControl(source)
		{
		if source.Base?(VirtualListDummyEditButton)
			{
			rec = source.GetRecord()
			return .expandedRows.FindOne({ it.rec is rec }).layout.ctrl
			}

		for row in .expandedRows
			{
			layout = row.layout
			if .findRecordControl(source, layout.ctrl.GetControl())
				return layout.ctrl
			c = layout.ctrl.FindControl(source.Name)
			if Same?(c, source)
				return layout.ctrl
			}
		return false
		}

	findRecordControl(source, recordCtrl)
		{
		if not source.Member?('Controller') or source.Name is "VirtualListView" or
			source.Base?(WindowBase)
			return false
		if Same?(source.Controller, recordCtrl)
			return true
		return .findRecordControl(source.Controller, recordCtrl)
		}

	GetExpandedControl(rec)
		{
		row = .expandedRows.FindOne({ it.rec is rec })
		return row is false ? false : row.layout
		}

	SetExpandReadOnly(rec, readonly = false)
		{
		if false is ctrl = .GetRecordControl(rec)
			return

		ctrl.SetReadOnly(readonly)
		if false isnt editBtn = ctrl.FindControl('Edit')
			editBtn.Pushed?(not readonly)
		}

	SetExpandRecord(rec, oldrec)
		{
		if false is row = .expandedRows.FindOne({ it.rec is oldrec })
			return
		row.rec = rec
		if false isnt ctrl = .GetRecordControl(rec)
			ctrl.Set(rec)
		}

	GetRecordControl(rec)
		{
		if false is ctrlOb = .GetExpandedControl(rec)
			return false
		ctrl = ctrlOb.ctrl.GetControl()
		if ctrl.Base?(RecordControl)
			return ctrl
		return false
		}

	GetControls()
		{
		ob = Object()
		for row in .expandedRows
			ob.Add(row.layout.ctrl)
		return ob
		}

	GetExpanded()
		{
		return .expandedRows.Map({ it.rec })
		}

	UpdateExpand(@unused)
		{
		}

	destroy(layoutOb, grid = false)
		{
		ctrl = layoutOb.ctrl
		if grid is false
			grid = ctrl.Parent
		child = ctrl.GetControl()
		if child.Base?(RecordControl)
			{
			// .SetVisible(false) kills the focus automatically in Win32
			// do here manually to be consistent
			if child.HasFocus?()
				SetFocus(NULL)
			child.SetVisible(false)
			child.Set(Record())
			.recycledExpands.Add(ctrl)
			grid.Act('VirtualListExpand_Recycle', ctrl.UniqueId)
			}
		else
			{
			grid.Act('VirtualListExpand_Destroy', ctrl.UniqueId)
			ctrl.Destroy()
			}
		}

	reuseRecycledExpand(recordExpand?)
		{
		if not recordExpand?
			return false
		return .recycledExpands.Extract(0, false)
		}

	RecycleExpands()
		{
		.CollapseAll()
		return .recycledExpands
		}

	SetRecycledExpands(.recycledExpands)
		{
		}

	Customize(query, fields, defaultExpandLayout, customKey)
		{
		return CustomizeExpandControl(query, fields, defaultExpandLayout, customKey)
		}

	CustomizableExpand?(layoutOb)
		{
		return .hasCustomizableExpand?(layoutOb)
		}

	hasCustomizableExpand?(layoutOb)
		{
		for idx in layoutOb.Members()
			{
			// NOTE: this relies on the fact that .hasCustomizable? does not get called
			// IF .isCustomizableControl? returns true
			item = layoutOb[idx]
			if Object?(item) and item.Size() > 0
				if .isCustomizableExpandControl?(item) or
					.hasCustomizableExpand?(item) is true
					return true
			}
		return false
		}

	isCustomizableExpandControl?(item)
		{
		item[0] is 'Customizable' and item.Has?(CustomizeExpandControl.LayoutName)
		}

	Map_GetAddress(source)
		{
		addressFieldControl =  source.Controller
		recCtrl = .GetControl(addressFieldControl)
		record = recCtrl.GetControl().Get()
		address1 = addressFieldControl.Name
		prefix = address1.BeforeLast('address1')
		return Object(address1: record[address1],
			address2: record[prefix $ 'address2'],
			city: record[prefix $ 'city'],
			state_prov: record[prefix $ 'state_prov'],
			zip_postal: record[prefix $ 'zip_postal'],
			country: record[prefix $ 'country'])
		}

	GetExpandRecord(source = false, _expandRec = false)
		{
		if expandRec isnt false
			return expandRec
		else if false isnt ctrl = .GetControl(source)
			return ctrl.GetControl().Get()
		else
			return false
		}

	DestroyAll()
		{
		for row in .expandedRows
			row.layout.ctrl.Destroy()
		.expandedRows = Object()

		while false isnt expand = .recycledExpands.Extract(0, false)
			expand.Destroy()
		}
	}
