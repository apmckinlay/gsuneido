// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// this control allows both the choosing of what columns will appear and in what
// order, but also the width of the columns. By default it will show the value
// stored in the datadict for the given field, but allows users to change the value
// which will be stored in a new field
Controller
	{
	Xstretch: 1
	Ystretch: 1
	Xmin: 200
	Ymin: 100
	Name: "CustomizableTwoList"
	New(.displayMemberName, list = #(), initial_list = #())
		{
		super(.layout(list))
		.storeOriginalValues(list)

		.availableListCtrl = .FindControl('AvailableOptions')
		.selectedList = .FindControl('SelectedOptions')
		.populateLists(list, initial_list.DeepCopy())
		}

	layout(list)
		{
		.cols = list.Empty?()
			? Object(.displayMemberName)
			: list[0].Members().Sort!().Remove(.displayMemberName).
				Add(.displayMemberName, at: 0)
		return Object('HorzSplit',
			Object('Horz'
				#(Vert
				(Static " Available")
				(Skip 2)
				(ListBox name: "AvailableOptions", sort:) name: 'Vert1')
				#Skip,
				TwoListControl.AddRemoveButtons(),
				xstretch: 1
			),
			Object('Horz',
				Object('Vert',
					#(Horz #(Static " Selected") Skip Fill
						(Static 'Right-click a line to reset'))
					#(Skip 2)
					Object('ListStretch' name: "SelectedOptions", columns: .cols,
						alwaysHighlightSelected:, noHeaderButtons:, noDragDrop:,
						stretchColumn: .displayMemberName)
						name: 'Vert2')
				'Skip',
				TwoListControl.MoveButtons()
				xstretch: .cols.Size()
			))
		}

	storeOriginalValues(list)
		{
		.origValues = list.DeepCopy()
		}

	populateLists(available, selected)
		{
		for item in available
			if not selected.HasIf?({ it[.displayMemberName] is item[.displayMemberName]})
				.availableListCtrl.AddItem(item[.displayMemberName])
		.selectedList.Set(selected)

		invalid = selected.FindAllIf(
			{|inSel| not available.HasIf?(
				{|inAvail| inAvail[.displayMemberName] is inSel[.displayMemberName]})
			})
		invalid.Each({ .selectedList.AddInvalidCell(0, it) })
		}

	Get()
		{
		return .selectedList.Get().Map({ it.Project(.cols) }).Instantiate()
		}

	ListBoxDoubleClick(item)
		{
		if item is -1 or false is .addItem(item)
			return
		.availableListCtrl.DeleteItem(item)
		.selectedList.SetSelection(.selectedList.GetNumRows() - 1)
		}

	addItem(rowIdx)
		{
		text = .availableListCtrl.GetText(rowIdx)
		if not .selectedList.Get().HasIf?({ it[.displayMemberName] is text })
			{
			if false is row = .origValues.FindOne({ it[.displayMemberName] is text })
				return false
			.selectedList.AddRow(row)
			}
		return 0
		}

	List_AllowCellEdit(col, row)
		{
		return col isnt 0 and row isnt false
		}

	List_DoubleClick(row, col)
		{
		return not .List_AllowCellEdit(col, row) ? false : 0
		}

	List_WantNewRow(@unused)
		{
		return false
		}

	List_DeleteKeyDown()
		{
		return false
		}

	List_AfterEdit(col, row, data /*unused*/, valid?)
		{
		if not valid?
			.selectedList.AddInvalidCell(col, row)
		else
			.selectedList.RemoveInvalidCell(col, row)
		}

	Valid()
		{
		for row in .selectedList.Get()
			if not row.GetDefault('list_invalid_cells', #()).Empty?()
				return false
		return true
		}

	Context_Menu: ('&Reset')
	List_ContextMenu(x, y)
		{
		selection = .selectedList.GetSelection()
		if selection.Size() is 0
			return 0
		ContextMenu(.Context_Menu).ShowCall(this, x, y)
		return 0
		}

	On_Context_Reset()
		{
		sel = .selectedList.GetSelection()[0]
		row = .selectedList.GetRow(sel)
		orig = .origValues.FindOne(
			{ it[.displayMemberName] is row[.displayMemberName] })
		if orig isnt false
			{
			.selectedList.DeleteSelection()
			.selectedList.InsertRow(sel, orig.Copy())
			.selectedList.SetSelection(sel)
			}
		}

	On_All()
		{
		for i in .. .availableListCtrl.GetCount()
			.addItem(i)
		.availableListCtrl.DeleteAll()
		.selectedList.SetSelection(.selectedList.GetNumRows() - 1)
		}

	On_Move()
		{
		if ((-1 is rowIdx = .availableListCtrl.GetCurSel()) or
			(false is .addItem(rowIdx)))
			return
		.availableListCtrl.DeleteItem(rowIdx)
		.availableListCtrl.SetCurSel(-1)
		.selectedList.SetSelection(.selectedList.GetNumRows() - 1)
		}

	On_MoveBack()
		{
		sel = .selectedList.GetCurrentRecord()
		if false is sel or false is .removeItem(sel)
			return
		.selectedList.DeleteSelection()
		.selectedList.ClearSelectFocus()
		}

	removeItem(sel)
		{
		text = sel[.displayMemberName]
		if false is .origValues.FindOne({ it[.displayMemberName] is text })	and
			not sel.Member?('list_invalid_cells')
			return false
		.availableListCtrl.AddItem(text)
		return 0
		}

	On_AllBack()
		{
		deleted = Object()
		for i in .. .selectedList.Get().Size()
			if false isnt .removeItem(.selectedList.GetRow(i))
				deleted.Add(i)
		.selectedList.DeleteRows(@deleted)
		}

	// Move Up and Move Down are only done on list2 (ListControl)
	On_Move_Up()
		{
		sel = .selectedList.GetSelection()
		if sel is #() or sel[0] is 0
			return
		sel = sel[0]
		.swapWith(sel, sel - 1)
		}

	On_Move_Down()
		{
		sel = .selectedList.GetSelection()
		if sel is #() or sel[0] is .selectedList.Get().Size() - 1
			return
		sel = sel[0]
		.swapWith(sel, sel + 1)
		}

	swapWith(sel, next)
		{
		nextOb = .selectedList.Get()[sel].Copy()
		.selectedList.DeleteRows(sel)
		.selectedList.InsertRow(next, nextOb)
		.selectedList.SetSelection(next)
		}
	}
