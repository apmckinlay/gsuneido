// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
VirtualListSaveBase
	{
	New(.grid = false, .model = false, .updateHistoryFn = false)
		{
		super(grid, model)
		}
	DeleteRecord(rec, block = false)
		{
		if false is ProtectRuleAllowsDelete?(
			rec, .model.EditModel.ProtectField, rec.New?(), header_delete: false)
			return 'Reason Protected may provide additional info'

		if block isnt false and false is block()
			return false

		if rec.New?()
			return true
		rec.RemoveObserver(VirtualListObserverOnChange)
		result = false
		KeyException.TryCatch(
			block:
				{
				RetryTransaction()
					{ |tran|
					cur = tran.Query1(.model.GetKeyQuery(rec.vl_origin, save?:))
					if not .RecordConflict?(rec, cur)
						{
						.grid.Controller.Send('VirtualList_BeforeDelete', rec, tran)
						cur.Delete()
						result = true
						}
					}
				}
			catch_block:
				{|e|
				result = false
				rec.Observer(VirtualListObserverOnChange)
				KeyException(e)
				})
		if result is true
			{
			.grid.Controller.Send('VirtualList_AfterDelete', rec)
			.model.DeleteRecordAttachments(rec)
			.model.CleanupAttachments()
			}
		return result
		}

	SaveRecord(record, highlighErr? = false)
		{
		.grid.Send('SetStatusBar', '')

		if false is .Check_valid_field(record)
			{
			if highlighErr?
				.grid.SelectRecord(record)
			return false
			}

		data = Object(record)
		save = false
		try
			KeyExceptionTransaction(update:)
				{|t|
				if false is .grid.Controller.Send(
					"VirtualList_BeforeSave", data: record, :t)
					return false

				if '' isnt msg = .grid.Controller.Addons.Collect(
					'SaveValid', record, record.vl_origin, t).Remove("").Join('\n\n')
					{
					.grid.AlertError('Save', msg)
					return false
					}

				(.updateHistoryFn)(record)
				save = .DoSave(:t, :data, clearChanges?:)
				}
		catch (unused, "interrupt: KeyException")
			return false
		result = save ? .ResetRecord(.model, record, clear?: .model.AutoSave?) : false
		.grid.Controller.Send('VirtualList_AfterSaving', :record)
		.grid.Controller.Send('VirtualList_AfterChanged', :record, saved:)
		return result
		}

	ResetRecord(model, record, clear? = false)  // also called after no change
		{
		prev_outstandings = model.EditModel.GetOutstandingChanges(all?:).Copy()
		model.EditModel.ClearChanges(record)
		if clear? is true and
			(outstandings = model.EditModel.GetOutstandingChanges(all?:)).NotEmpty?()
			{
			SuneidoLog('INFO: outstanding changes should be empty, but it is not',
				calls: Display([:record, :outstandings, :prev_outstandings]))
			}
		model.NextNum.ConfirmNextNum(record)
		model.UnlockRecord(record)
		try
			return model.ReloadRecord(record, force:)
		catch (err, 'expected the value to not be false but it was')
			{
			SuneidoLog('ERROR: cannot find original record to edit, rethrow error',
				calls:,
				params: Object(:record, loaded: .model.GetLoadedData().Size(),
					offset: .model.Offset, visibleRows: .model.VisibleRows,
					modelCreated: .model.Created))
			throw err
			}
		}
	}