// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	allLibsView: svc_all_changes
	New(table, svcEnsure = false)
		{
		table = table.Tr('()')
		.svcTable = .tableClass(table)
		if table is .allLibsView
			return
		// This is required in order to handle older DBs / Tables
		Database('ensure ' $ table $ '(lib_modified)')
		if svcEnsure and not .SvcEnabled?()
			.ensureSvcColumns(table)
		}

	tableClass(table)
		{
		return QueryColumns(table).Has?(#group) or table is .allLibsView
			? SvcLibrary(table)
			: SvcBook(table)
		}

	SvcColumns: #(lib_committed, lib_before_hash, lib_before_text, lib_before_path)
	ensureSvcColumns(table)
		{
		Database('ensure ' $ table $ ' (' $ .SvcColumns.Join(', ') $ ')')
		.ResetSvcDisabledCache()
		}

	Getter_Type()
		{
		return .svcTable.Type
		}

	Default(@args)
		{
		return (.svcTable[args.Extract(0)])(@args)
		}

	optionalCall(@args)
		{
		method = args.Extract(0)
		if .svcTable.Method?(method)
			(.svcTable[method])(@args)
		}

	Output(rec, t = false, deleted = false, committed = false)
		{
		DoWithTran(t, update:)
			{|t|
			.FormatData(rec, t, :deleted)
			outputTable = .Table()
			if .needNewNum?(rec, t)
				rec.num = NextTableNum(outputTable, t)
			if not committed
				rec.lib_modified = Date()
			t.QueryOutput(outputTable, rec)
			}
		.Publish(#TreeChange, name: rec.name)
		}

	// Book records are output during "Move", as a result, we need to keep the previous
	// num in order to avoid issues with the ExplorerMultiTreeControl
	// Additionally, this allows us to "recycle" nums during renames
	needNewNum?(rec, t)
		{
		return rec.Member?(#num)
			? false isnt t.Query1(.Table(), num: rec.num)
			: true
		}

	Get(name, t = false, deleted = false)
		{
		DoWithTran(t)
			{|t|
			rec = t.Query1(.NameQuery(name, :deleted))
			if rec isnt false and not deleted
				.GetData(rec, t)
			}
		return rec
		}

	// If record is new: it is deleted immediately
	// Otherwise: Stage for deletion, (not removed until committed to Svc)
	// newNum: true, is used for deletes which occur with Moves / Renames.
	// The live record keeps its number while the deleted version is given a new number
	StageDelete(name, t = false, newNum = false, skipPublish = false)
		{
		committed? = false
		DoWithTran(t, update:)
			{|t|
			if false is rec = .Get(name, :t)
				return
			if not committed? = rec.lib_committed isnt ''
				rec.Delete()
			else
				{
				.Remove(name, :t, deleted:)
				.versionData(rec, t, newNum)
				.optionalCall(#StageDelete, rec, :t, :newNum)
				.Update(rec, :t, deleted:)
				}
			}
		if not skipPublish
			.Publish(#TreeChange, :name)
		.optionalCall(#PostStageDelete, :name, :skipPublish, :committed?)
		}

	versionData(rec, t, newNum)
		{
		if newNum
			rec.num = NextTableNum(.Table(), t)
		rec.lib_modified = Date()
		rec.lib_before_path = rec.GetDefault(#lib_before_path, rec.path)
		rec.lib_before_text = rec.GetDefault(#lib_before_text, rec.text)
		}

	Restore(name, t = false)
		{
		restored? = false
		DoWithTran(t, update:)
			{|t|
			if false isnt rec = .Get(name, :t)
				restored? = .restore(rec, t)
			else if false isnt rec = .Get(name, :t, deleted:)
				restored? = .restore(rec, t, deleted:)
			}
		.optionalCall(#PostRestore, :name, :restored?)
		}

	restore(rec, t, deleted = false)
		{
		if rec.lib_committed is '' // New record, simply delete it
			{
			rec.Delete()
			.Publish(#TreeChange, name: rec.name, force:)
			return true
			}
		if rec.lib_modified is '' // Nothing to restore as record is not modified
			return false
		.restoreBeforeFields(rec)
		rec.name = .MakeName(rec)
		rec.lib_modified = ''
		.optionalCall(#Restore, rec)
		.FormatData(rec, :t, deleted: false)
		rec.Update()
		.Publish(deleted ? #TreeChange : #RecordChange, name: rec.name, force:)
		return true
		}

	restoreBeforeFields(rec)
		{
		rec.text = rec.Extract(#lib_before_text, rec.text)
		rec.path = rec.Extract(#lib_before_path, rec.path)
		}

	Rename(rec, newName, t = false)
		{
		DoWithTran(t, update:)
			{|t|
			oldName = .MakeName(rec)
			rec.name = newName
			rec.name = .MakeName(rec)
			.optionalCall(#Rename, rec, oldName, t)
			if false isnt deletedRec = .Get(rec.name, t, deleted:)
				{
				rec.lib_before_text = deletedRec.lib_before_text
				rec.lib_before_path = deletedRec.lib_before_path
				rec.lib_committed = deletedRec.lib_committed
				deletedRec.Delete()
				}
			else
				rec.Delete(#lib_before_text, #lib_before_path, #lib_committed)
			.StageDelete(oldName, t, newNum:)
			.VerifyModified(rec)
			.Output(rec, t, committed: rec.lib_modified is '')
			}
		.Publish(#RecordChange, name: newName)
		.Publish(#TreeChange)
		}

	VerifyModified(rec)
		{
		.pathRestored(rec)
		.textRestored(rec)
		.modified(rec)
		}

	pathRestored(modifiedRec)
		{
		pathRestored = modifiedRec.lib_before_path is '' or
			modifiedRec.path is modifiedRec.lib_before_path
		if pathRestored
			modifiedRec.Delete(#lib_before_path)
		}

	textRestored(modifiedRec)
		{
		beforeText = modifiedRec.lib_before_text
		if beforeText is '' or modifiedRec.lib_current_text is beforeText
			{
			// Overwrite text with lib_before_text in order to undo FormatText
			if beforeText isnt ''
				modifiedRec.text = beforeText
			modifiedRec.Delete(#lib_before_text)
			}
		}

	modified(rec)
		{
		modifiedState = Date?(rec.lib_modified)
		rec.lib_modified = rec.Member?(#lib_before_text) or
			rec.Member?(#lib_before_path) or rec.GetDefault(#lib_committed, '') is ''
			? Date()
			: ''
		restored? = modifiedState and rec.lib_modified is ''
		.optionalCall(#PostRestore, name: .MakeName(rec), :restored?)
		}

	Move(rec, newParent, t = false)
		{
		DoWithTran(t, update:)
			{|t|
			if rec.lib_before_path is ''
				rec.lib_before_path = .GetPath(rec)
			.optionalCall(#Move, rec, newParent, t)
			}
		}

	// NOTE: Update needs to be called within a pre-existing update query transaction
	// 	IE: QueryApply1, QueryApplyMulti(..., update:), etc.
	Update(rec, t = false, newText = false, deleted = false, force = false)
		{
		if rec.lib_before_text is ''
			rec.lib_before_text = rec.text
		name = .MakeName(rec)
		if '' isnt invalidText = rec.GetDefault('lib_invalid_text', '')
			rec.lib_invalid_text = .FormatText(invalidText, :name)
		else if newText isnt false
			rec.text = .FormatText(newText, :name)
		if not deleted
			.VerifyModified(rec)
		.FormatData(rec, :t, :deleted)
		rec.Update()
		if not deleted
			.Publish(#RecordChange, name: rec.name, :force)
		.optionalCall(#PostUpdate, rec, t)
		}

	Getter_ExcludedTables()
		{
		return #(configlib, Test_lib)
		}

	// Completely removes the record from the table, no staging for delete occurs
	Remove(name, t = false, deleted = false)
		{
		removed? = false
		DoWithTran(t, update:)
			{|t|
			removed? = t.QueryDo('delete ' $ .NameQuery(name, :deleted)) isnt 0
			}
		// Do not publish TreeChange if removing a deleted record
		if not deleted and removed?
			.Publish(#TreeChange, :name)
		return removed?
		}

	GetMaxCommitted()
		{
		if not .SvcEnabled?()
			return Date.Begin()
		return false is (maxCommit = .maxCommit())
			? .outputMaxCommit()
			: maxCommit.lib_committed
		}

	maxCommit()
		{
		return Query1(.MaxCommitQuery())
		}

	MaxCommitName: #libCommitted
	outputMaxCommit()
		{
		Transaction(update:)
			{|t|
			lib_committed = t.QueryMax(.Table(), #lib_committed, Date.Begin())
			.FormatMaxCommitRecord(rec = [
				name: .MaxCommitName,
				num: NextTableNum(.Table(), t),
				:lib_committed
				])
			t.QueryOutput(.Table(), rec)
			}
		return lib_committed
		}

	SetMaxCommitted(date, force = false)
		{
		if false is .maxCommit()
			.outputMaxCommit()
		QueryApply1(.MaxCommitQuery())
			{
			it.lib_committed = not force
				? Max(it.lib_committed, date)
				: date
			it.Update()
			}
		}

	ModifiedQuery()
		{
		extend = not .SvcEnabled?()
			? ' extend lib_committed'
			: ''
		modifiedWhere = ' where lib_committed is "" or lib_modified isnt ""'
		return .Table() $ extend $ .Where() $ modifiedWhere
		}

	// This method can be called without initializing SvcTable. This is to keep the
	// arguments consistent when calling PubSub.PublishConsolidate
	// IF Being called with an instance, it will always use the instance values
	Publish(event, name = '', force = false, type = false, table = false)
		{
		if TestRunner.RunningTests?()
			return
		if Instance?(this)
			{
			table = .Table()
			type = .Type
			}
		event = (type is #lib ? #Library : #Book) $ event
		PubSub.PublishConsolidate(event, :table, :name, :force)
		}

	Compare(local, master)
		{
		discrepancies = Object()
		if local.lib_committed isnt master.lib_committed
			discrepancies.Add(#lib_committed)
		if local.path isnt master.path
			discrepancies.Add(#path)
		.optionalCall(#Compare, local, master, discrepancies)
		return discrepancies
		}

	Check(name, type = '')
		{
		if false is rec = .Get(name)
			return 'Get failed'
		return .CheckRecord(rec, type)
		}
	}