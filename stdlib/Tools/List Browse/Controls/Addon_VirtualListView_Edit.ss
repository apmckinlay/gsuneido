// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	On_New()
		{
		if .SaveFirst() is false
			return
		.Grid.KEYDOWN(VK.END)
		.Grid.InsertRow(pos: 'end')
		}

	On_Context_New()
		{
		c = VirtualListViewExtra
		if c.Member?('On_Context_New')
			c.On_Context_New(addon: this)
		else
			.ContextNew()
		}

	ContextNew()
		{
		if .GetContextMenu().ContextRec is false
			{
			.Grid.KEYDOWN(VK.END)
			.Grid.InsertRow(pos: 'end')
			}
		else
			.Grid.InsertRow(pos: 'current')
		}

	Record_NewValue(field, value, source)
		{
		if source.Name is 'VirtualListParams'
			{
			if source.Valid(forceCheck:) isnt true
				return

			.Parent.Addons.Send('ResetWhere')
			.ResetWhere()
			return
			}
		// from expand record
		if not source.Base?(RecordControl)
			return
		if false is fieldCtrl = source.GetControl(field)
			return
		valid? = fieldCtrl.Valid?()
		VirtualListEdit.CommitChange(.Grid, source.Get(), field, value, valid?)
		}

	On_Edit(source, force = false)
		{
		curFocus = GetFocus()
		// set focus first, to commit onging change from expand
		.Grid.SetFocus()
		retVal = false
		if false is result = .GetExpandCtrlAndRecord(source, curFocus)
			{
			rec = .GetSelectedRecord()
			if rec isnt false and rec.vl_expanded_rows is ''
				.Grid.EditField(rec, .Model.ColModel.Get(0))
			retVal = false
			}
		else
			{
			if .Model.EditModel.HasOtherLockedRecord?(result.rec)
				.SaveOutstandingChanges()
			rec = .Model.ReloadRecord(result.rec)
			if String?(rec)
				{
				.Grid.SelectRecord(result.rec)
				.AlertInfo(.GetAlertTitle(), rec)
				retVal = false
				}
			else
				{
				.Grid.SelectRecord(rec)
				if not force and .Model.EditModel.RecordLocked?(rec)
					{
					.VirtualListGrid_SaveRecord(rec)
					retVal = true
					}
				else
					retVal = .edit(rec, result.ctrl)
				}
			}
		.GetViewControls().expandBar.RefreshEditState()
		return retVal
		}

	edit(rec, ctrl)
		{
		if .GetReadOnly() is true
			{
			.AlertInfo('List Protected',
				'Edit is not allowed because the list is protected')
			return false
			}
		if true isnt msg = .Model.LockRecord(rec)
			{
			.AlertInfo('Reason Protected', msg)
			return false
			}
		validField = .Model.EditModel.ValidField
		if validField isnt false
			rec[validField]
		if ctrl isnt false
			.FocusFirst(ctrl.Hwnd)
		return true
		}

	ForceEditMode(rec)
		{
		rec = .Model.ReloadRecord(rec)
		if String?(rec)
			{
			.AlertInfo(.GetAlertTitle(), rec)
			return false
			}
		.Grid.SelectRecord(rec)
		return .edit(rec, false) ? rec : false
		}

	SaveOutstandingChanges()
		{
		if .Model is false
			return true

		if false is .Send('VirtualList_SaveOutstandingChanges?')
			return true
		// set focus first, to commit onging change from expand
		if .Model.ExpandModel isnt false and
			false isnt .Model.ExpandModel.GetCurrentFocusedRecord(GetFocus())
			.Grid.SetFocus()
		records = .Model.EditModel.GetOutstandingChanges()
		for key in .Model.EditModel.LockedKeys()
			{
			rec = .Model.GetRecordByKeyPair(key, .Model.EditModel.LockKeyField, str?:)
			Assert(rec isnt false)
			records.Add(rec)
			}
		for rec in records.Copy()
			if not .Grid.CommitRecord(rec, highlighErr?:)
				return false
		return true
		}

	SaveFirst()
		{
		if .Editable?()
			{
			.SaveOutstandingChanges()
			if .Model.EditModel.HasChanges?()
				{
				msg = .IsLinked?()
					? 'Please save outstanding changes.'
					: 'Please correct the highlighted row.'
				.AlertInfo(.GetAlertTitle(), msg)
				return false
				}
			}
		return true
		}

	On_Context_Edit_Field()
		{
		c = VirtualListViewExtra
		if c.Member?('On_Context_Edit_Field')
			c.On_Context_Edit_Field(addon: this)
		else
			.ContextEditField()
		}

	ContextEditField()
		{
		if .GetContextMenu().ContextCol isnt false
			.Grid.EditField(.GetContextMenu().ContextRec, .GetContextMenu().ContextCol)
		}

	On_Context_Restore()
		{
		.RestoreAttachmentFiles()
		rec = .GetContextMenu().ContextRec
		.Model.EditModel.ClearChanges(rec)
		.Model.NextNum.CheckPutBackNextNum(rec)
		prevRec = rec
		if rec.New?()
			.Send('VirtualList_NewRowAdded', rec = .Model.RestoreNewRecord(rec))
		else
			{
			.Model.UnlockRecord(rec)
			rec = .Model.ReloadRecord(rec)
			}
		if String?(rec)
			{
			.AlertInfo(.GetAlertTitle(), rec)
			return
			}
		.Grid.ClearHighlightRecord(prevRec)
		.Grid.SelectRecord(rec)
		if .Model.ExpandModel isnt false and
			false isnt ctrlOb = .Model.ExpandModel.GetExpandedControl(rec)
			{
			recCtrl = ctrlOb.ctrl.GetControl() // RecordControl
			recCtrl.Set(rec)
			recCtrl.SetAllValid()
			}
		.RefreshValid(rec)
		.Send('VirtualList_AfterChanged', record: rec, saved: false)
		}

	VirtualListGrid_SaveRecord(rec, highlighErr? = false)
		{
		if not .Model.EditModel.RecordChanged?(rec)
			return VirtualListSave.ResetRecord(.Model, rec)

		.Send('VirtualList_BeforeSave_PreValid', rec)
		if .Model.ExpandModel isnt false and
			false isnt ctrlOb = .Model.ExpandModel.GetExpandedControl(rec)
			{
			recCtrl = ctrlOb.ctrl.GetControl() // RecordControl
			if recCtrl.Base?(RecordControl) and true isnt msg = recCtrl.Valid(forceCheck:)
				{
				.SetStatusBar(msg)
				return false
				}
			}

		if false is .Send('VirtualList_BeforeSave_PreTran', :rec)
			return false
		save = VirtualListSave(.Grid, .Model, .GetContextMenu().UpdateHistory)
		return save.SaveRecord(rec, highlighErr?)
		}

	Keyboard_Delete(rec)
		{
		if .Grid.GetSelectedRecords().Size() isnt 1
			{
			.AlertInfo('Delete Record', 'Please Select a single record to delete')
			return
			}

		.On_Context_Delete(rec)
			{
			OkCancel("Delete immediately and permanently. " $
				"This cannot be undone.", title: "Delete Record")
			}
		}

	On_Context_Delete(item /*unused*/= false, block = false)
		{
		if not .Editable?()
			return

		save = VirtualListSave(.Grid, .Model, .GetContextMenu().UpdateHistory)
		selects = .Grid.GetSelectedRecords()
		for m in selects.Members()
			{
			rec = selects[m]
			if false is .allowDelete?(rec, save, block)
				return

			.Send("VirtualList_DeleteRecord", rec)
			.Send('VirtualList_AfterChanged', record: rec, saved: not rec.New?())
			}

		.AfterDelete()
		}

	allowDelete?(rec, save, block)
		{
		if rec.vl_deleted isnt true and false is .Send("VirtualList_AllowDelete", rec)
			return false

		preOwnLock? = .Model.OwnLock?(rec)
		if not preOwnLock? and false is rec = .ForceEditMode(rec)
			return false

		if '' isnt msg = RecordAllowDelete(.Model.GetQuery(), rec)
			return .cannotDeleteRecord(preOwnLock?, rec, msg)

		return .deleteRec(rec, save, preOwnLock?, block)
		}

	cannotDeleteRecord(preOwnLock?, rec, msg)
		{
		if not preOwnLock? // there could be changes not saved
			.Model.UnlockRecord(rec)
		.Grid.Send('SetStatusBar', 'This record can not be deleted. ' $ msg)
		return false
		}

	deleteRec(rec, save, preOwnLock?, block)
		{
		if true is result = save.DeleteRecord(rec, block)
			{
			.DeleteRow(rec)
			return true
			}
		else if not preOwnLock?
			.Model.UnlockRecord(rec)
		return String?(result)
			? .cannotDeleteRecord(true, rec, result)
			: result
		}

	AfterDelete()
		{
		.Model.Selection.ClearSelect(false)
		.RecordDirty?(true)
		.Grid.SelectFocusedRow()
		.Grid.Repaint()
		expandBar = .GetViewControls().expandBar
		expandBar.HideButton()
		expandBar.Repaint()
		.SetStatusBar('')
		}

	DeleteRow(rec)
		{
		.Model.UnlockRecord(rec)
		.removeRow(rec)
		.Model.NextNum.CheckPutBackNextNum(rec)
		.RepaintExpandBar()
		}

	RemoveRowByKeyPair(keyVal, keyField)
		{
		if false is rec = .Model.GetRecordByKeyPair(keyVal, keyField)
			return
		.removeRow(rec)
		.RepaintExpandBar()

		.Grid.SelectFocusedRow()
		.Grid.Repaint()
		}

	removeRow(rec)
		{
		rowNum = .Model.GetRecordRowNum(rec)
		.Grid.ToggleExpand(rowNum, expand: false)
		.Grid.DeleteRow(rec, rowNum)
		}

	On_Context_Save()
		{
		rec = .GetSelectedRecord()
		if .Model.EditModel.RecordLocked?(rec)
			{
			if false isnt .VirtualListGrid_SaveRecord(rec)
				return false
			}
		else
			.Model.ReloadRecord(rec)
		.Grid.RepaintSelectedRows()
		return true
		}
	}
