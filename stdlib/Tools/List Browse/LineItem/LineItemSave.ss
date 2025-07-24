// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
VirtualListSaveBase
	{
	CallClass(@args)
		{
		instance = new this(@args)
		return instance.Result
		}

	New(.action, t = false, .grid = false, .model = false, .list = false)
		{
		super(.grid, .model)
		// don't need to handle 'restore'
		// since it is already handled by AccessControl calling Set on the record,
		// which reset the LineItemControl
		switch (.action)
			{
		case 'delete':			.Result = .delete()
		case 'valid':			.Result = .Valid?()
		case 'save':			.Result = .save(t)
		case 'after_save': 		.Result = .accessAfterSave()
		case 'accessInvalid': 	.Result = .accessInvalid()
		default:				.Result = true
			}
		}

	// checks if OK to delete, empties the list and deletions list
	// NOTE: 	There should be a foreign key from lines to header
	//			with cascading deletes enabled.
	delete()
		{
		for rec in .model.GetLoadedData()
			if false is ProtectRuleAllowsDelete?(
				rec, .model.EditModel.ProtectField, rec.New?(), header_delete:)
				return false
		for record in .model.GetLoadedData()
			record.RemoveObserver(VirtualListObserverOnChange)
		return true
		}

	save(t)
		{
		.grid.FinishEdit()
		if not .Valid?(evalRule?:, selectFirstInvalid?:)
			return false

		data = .model.GetLoadedData()
		// does not clear changes since AccessControl save transaction could still fail
		// and we want the line item changes to still be processed on the next save call.
		// The accessAfterSave method is responsible for clearing changes
		return .DoSave(t, data)
		}

	HandleDelete(cur, record, tran)
		{
		result = .list.Send('LineItem_AllowDelete', record, tranFromSave: tran)
		if result is 0 or result is true
			{
			cur.Delete()
			.model.DeleteRecordAttachments(record)
			return true
			}
		if String?(result)
			AlertDelayed(result, 'Delete Record')
		return false
		}

	HighlightConflictRecord?()
		{
		return false
		}

	OriginalRecordNotFound_ExtraMsg()
		{
		return "Please use Current > Restore and re-do your changes if necessary."
		}

	accessAfterSave()
		{
		.list.Set(.list.Get())
		return true
		}

	accessInvalid()
		{
		// when the list is flipped to form and user entered invalid data on the form
		// list.Valid?() won't be called before this
		if .model.Member?('FirstInvalidRecord')
			.list.SelectRecord(.model.FirstInvalidRecord)
		return true
		}
	}
