// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// REFERENCE: BrowseControl
class
	{
	New(.grid = false, .model = false)
		{
		}
	DoSave(t, data)
		{
		saveResult = true
		keyExceptionResult = KeyException.TryCatch(
			block:
				{
				DoWithTran(t, update:)
					{|t|
					saveResult = .saveList(data, t)
					}
				if t.Ended?() // when transaction rolled back
					{
					AlertDelayed('Unable to save')
					return .selectErrorLine()
					}
				}
			catch_block:
				{|e|
				if not t.Ended?()
					t.Rollback()
				KeyException(e)
				return .selectErrorLine()
				}
			)
		.model.CleanupAttachments()
		return keyExceptionResult and saveResult
		}

	selectErrorLine()
		{
		// if error occurred after list deletions in the save method,
		// save_rec may not be in list, must check
		if .save_rec isnt false
			.grid.SelectRecord(.save_rec)
		return false
		}

	saveList(data, tran)
		{
		.save_rec = false

		if data.Size() is 0
			return true

		query = QueryStripSort(.model.GetSaveQuery())
		outputs = Object()
		oldrecs = Object()

		// do ALL query deletes first to prevent duplicate keys when a key is
		// deleted, then inserted again.
		if false is .handleQueryDeletes(data, tran, oldrecs)
			return false

		// process the rest of the data
		.handleListUpdateDelete(data, outputs, oldrecs, tran)

		// do outputs last so we don't get conflicts with updates
		.outputListUpdates(outputs, tran, query)

		return true
		}

	conflictMsg: "Another user has modified records, restore or discard changes"
	handleQueryDeletes(data, tran, oldrecs)
		{
		valid = true
		for row in data.Members()
			{
			.save_rec = record = data[row]
			if not .existingRecordChanged?(record)
				continue

			Assert(record.vl_origin isObject:)
			cur = tran.Query1(.model.GetKeyQuery(record.vl_origin, save?:))
			if .RecordConflict?(record, cur, quiet?: valid is false)
				{
				if not .HighlightConflictRecord?()
					return false

				.highlightInvalidLine(valid, record)
				valid = false
				continue
				}

			if record.vl_deleted is true
				{
				if false is .HandleDelete(cur, :record, :tran)
					return false
				}
			else // record dirty
				oldrecs[row] = cur
			}
		.setStatusBar(valid)
		return valid
		}

	existingRecordChanged?(record)
		{
		return ((.model.EditModel.RecordChanged?(record) and not record.New?()) or
			record.vl_deleted is true)
		}

	highlightInvalidLine(currentValid, record)
		{
		// if valid is currenttly true, then this is the first record conflict
		if currentValid
			.grid.SelectRecord(record)
		.grid.HighlightRecords(Object(record), CLR.ErrorColor)
		}

	HandleDelete(cur)
		{
		cur.Delete()
		return true
		}

	HighlightConflictRecord?()
		{
		return true
		}

	setStatusBar(valid)
		{
		if not valid
			.grid.Send('SetStatusBar', .conflictMsg)
		}

	handleListUpdateDelete(data, outputs, oldrecs, t)
		{
		for (row in data.Members())
			{
			record = data[row]
			.save_rec = record
			if (record.listrow_deleted is true)
				continue
			// have to queue up outputs and do last so that
			// they don't conflict with updates
			if record.New?()
				{
				if record.vl_excluded isnt true and
					.model.EditModel.RecordChanged?(record)
					outputs.Add(record)
				}
			else if oldrecs.Member?(row)
				{
				oldrecs[row].Update(record)
				.grid.Controller.Send("VirtualList_AfterSave", data: record, :t)
				}
			}
		}

	outputListUpdates(outputs, tran, query)
		{
		if outputs.Empty?()
			return

		tran.Query(query)
			{|q|
			for (record in outputs)
				{
				.save_rec = record
				q.Output(record)
				.grid.Controller.Send("VirtualList_AfterSave", data: record, t: tran)
				}
			}
		}

	Check_valid_field(record, evalRule? = false)
		{
		.grid.Send('SetStatusBar', '')
		if record.vl_deleted is true
			return true

		if .model.EditModel.HasInvalidCols?(record)
			{
			.grid.Send('SetStatusBar', .model.EditModel.GetInvalidMsg(record))
			return false
			}

		if .model.EditModel.ValidField isnt false
			{
			msg = .validationMsg(evalRule?, record)
			.grid.Send('SetStatusBar', msg)
			return msg is ""
			}
		return true
		}

	validationMsg(evalRule?, record)
		{
		try
			msg = evalRule? is true
				? record.Eval(Global('Rule_' $ .model.EditModel.ValidField))
				: record[.model.EditModel.ValidField]
		catch (err)
			{
			msg = "There was a problem validating the record"
			SuneidoLog("ERROR: (CAUGHT) " $ err, caughtMsg: 'user alerted: ' $ msg)
			}
		return msg
		}

	RecordConflict?(record, cur, quiet? = false)
		{
		if cur is false
			{
			.alertOriginalRecordNotFound()
			SuneidoLog("INFO: List can't get the original record to update",
				params: record)
			return true
			}

		return RecordConflict?(record.vl_origin, cur,
			.model.AllAvailableColumns(), .grid.Window.Hwnd, :quiet?)
		}

	alertOriginalRecordNotFound()
		{
		msg = "List: can't get record to update.\n" $
			"Another user may have deleted the line.\n"
		msg $= .OriginalRecordNotFound_ExtraMsg()
		AlertDelayed(msg, title: 'Error', hwnd: .grid.Window.Hwnd, flags: MB.ICONERROR)
		}

	OriginalRecordNotFound_ExtraMsg()
		{
		return ''
		}

	Valid?(evalRule? = false)
		{
		.model.FirstInvalidRecord = false
		data = .model.GetLoadedData()
		for row in data.Members()
			{
			rec = data[row]
			if rec.vl_deleted is true
				continue
			dirty? = .model.EditModel.RecordChanged?(rec)
			if not dirty?
				continue
			if not .Check_valid_field(rec, evalRule?)
				{
				.model.FirstInvalidRecord = rec
				return false
				}
			}
		return true
		}
	}