// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
TwoListDlgControl
	{
	Title: 'Customize Columns'

	delimiter: '\x00'
	CallClass(hwnd, listCtrl, columns, title, mandatoryFields, headerSelectPrompt)
		{
		if columns is false
			return
		if listCtrl.GetColumns().Any?(Customizable.DeletedField?)
			{
			.AlertInfo(.Title, 'Your list contains a deleted custom field.\r\n' $
				'Please reload the screen before customizing columns.')
			return
			}

		listCtrl.SetHeaderChanged(true)
		UserColumns.Save(title, listCtrl, columns)

		available_cols = .buildAvailableColumns(columns, headerSelectPrompt, listCtrl)
		cur_cols = .buildCurrentColumns(listCtrl, available_cols)

		dlg = Object(this, available_cols.Values().Sort!(), cur_cols.list,
			mandatory_list: .mandatoryPrompts(
				mandatoryFields, headerSelectPrompt, listCtrl),
			saved_title: title, :available_cols, delimiter: .delimiter)
		if AccessPermissions(Customizable.PermissionOption()) is true
			dlg.extra_buttons = #(#(Button 'Save As Default')
				#(Static ' (for all users)', name: save_result))

		if false is newList = OkCancel(dlg, 'Customize Columns', hwnd)
			return

		newList = newList.Split(.delimiter)
		if newList is cur_cols.list
			return

		.handleReArrangedColumns(listCtrl, available_cols, newList)
		.handleColumnResize(listCtrl, available_cols, cur_cols.sizes, newList)
		}

	unavailableColumns: #('listrow_deleted', 'params_itemselected')
	getProtectedColumns(listCtrl)
		{
		protectedColumns = .unavailableColumns.Copy()
		if listCtrl.Method?('GetCheckBoxField') and
			false isnt listCheckBoxColumn = listCtrl.GetCheckBoxField()
				protectedColumns.AddUnique(listCheckBoxColumn)
		return protectedColumns
		}

	buildAvailableColumns(columns, headerSelectPrompt, listCtrl)
		{
		available_cols = Object()
		for col in columns.Copy().Remove(@.getProtectedColumns(listCtrl))
			if not Customizable.DeletedField?(col)
				available_cols[col] = .getPrompt(col, headerSelectPrompt, available_cols)
		return available_cols
		}
	buildCurrentColumns(listCtrl, available_cols)
		{
		cur_cols = Object()
		cols_size = Object()

		cols = listCtrl.GetColumns()
		protectedColumns = .getProtectedColumns(listCtrl)

		for i in .. cols.Size()
			{
			size = listCtrl.GetColWidth(i)
			cols_size[cols[i]] = size
			// assume col with width 0 is unselected col, so don't add to cur_cols list
			if size isnt 0 and not protectedColumns.Has?(cols[i])
				cur_cols.Add(available_cols[cols[i]])
			}

		return Object(list: cur_cols, sizes: cols_size)
		}
	getPrompt(col, headerSelectPrompt, existingPrompts = #())
		{
		return (Object?(col)
			? col[1]
			: headerSelectPrompt is false
				? Datadict.PromptOrHeading(col)
				: Datadict.GetFieldPrompt(col, existingPrompts)).Trim()
		}
	mandatoryPrompts(mandatoryFields, headerSelectPrompt, listCtrl)
		{
		prompts = Object()
		for field in mandatoryFields.Remove(@.getProtectedColumns(listCtrl))
			{
			prompt = headerSelectPrompt is false
				? Datadict.PromptOrHeading(field)
				: Datadict.GetFieldPrompt(field)
			prompts.Add(prompt.Trim())
			}
		return prompts
		}
	handleReArrangedColumns(listCtrl, available_cols, newList)
		{
		new_cur_cols = Object()
		for (i = 0; i < newList.Size(); ++i)
			new_cur_cols.Add(available_cols.Find(newList[i]))
		if not new_cur_cols.Empty?()
			{
			new_cur_cols = new_cur_cols.MergeUnion(
				listCtrl.GetColumns().Difference(new_cur_cols))
			.handleDeleteColSelectCol(new_cur_cols, listCtrl)
			listCtrl.SetColumns(new_cur_cols)
			}
		}

	handleDeleteColSelectCol(new_cur_cols, listCtrl)
		{
		if new_cur_cols.Has?('listrow_deleted')
			.positionColumnFirst(new_cur_cols, 'listrow_deleted')

		if listCtrl.Method?('GetCheckBoxField') and
			false isnt (checkBoxField = listCtrl.GetCheckBoxField()) and
			new_cur_cols.Has?(checkBoxField)
			.positionColumnFirst(new_cur_cols, checkBoxField)
		}

	positionColumnFirst(columns, field)
		{
		columns.Remove(field)
		columns.Add(field, at: 0)
		}

	handleColumnResize(listCtrl, available_cols, cols_size, newList)
		{
		diff_cols = available_cols.Values().Sort!().Difference(newList.Sort!())
		cols_list = listCtrl.GetColumns()
		protectedColumns = .getProtectedColumns(listCtrl)
		for (i = 0; i < cols_list.Size(); ++i)
			{
			width = not protectedColumns.Has?(cols_list[i]) and
				(not available_cols.Member?(cols_list[i]) or
					diff_cols.Has?(available_cols[cols_list[i]]))
				? 0
				: not cols_size.Member?(cols_list[i]) or cols_size[cols_list[i]] is 0
					? false
					: cols_size[cols_list[i]]
			listCtrl.SetColWidth(i, width)
			}
		listCtrl.Repaint()
		}
	result_ctrl: false
	New(@args)
		{
		super(@args)
		.saved_title = args.GetDefault('saved_title', '')
		.available_cols = args.GetDefault('available_cols', #())
		.result_ctrl = .FindControl('save_result')
		}
	On_Save_As_Default()
		{
		if not .valid()
			return
		newList = .TwoList.GetNewList()
		new_cur_cols = Object()
		for (i = 0; i < newList.Size(); ++i)
			new_cur_cols.Add(.available_cols.Find(newList[i]))
		UserColumns.SaveDefaultColumns(.saved_title, new_cur_cols,
			.available_cols.Members(), deletecol:)

		.result_ctrl.Set(' Default Columns Saved!')
		.result_ctrl.SetColor(0x007f00)
		.result_ctrl.SetFont(weight: FW.BOLD)
		}
	NewValue(data /*unused*/)
		{
		if .result_ctrl isnt false
			{
			.result_ctrl.Set(' (for all users)')
			.result_ctrl.SetColor(0x000000)
			.result_ctrl.SetFont(weight: FW.NORMAL)
			}
		}
	OK()
		{
		if not .valid()
			return false
		super.OK()
		}
	valid()
		{
		if .TwoList.GetNewList().Empty?()
			{
			.AlertWarn("Invalid Columns", "At least one column is required")
			return false
			}
		return true
		}
	}
