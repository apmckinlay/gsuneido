// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/*
SVC - Suneido Version Control

Svc is the main interface to the version control. It talks to an SvcCore either directly
or via SvcClient & SvcServer and is always run on the same database as the library/book.
SvcControl calls this class through SvcModel to perform local actions (through SvcLibrary
or SvcBook) or access the server (through SvcClient in client/server mode or SvcCore
directly in standalone mode).

Notable methods:

GetChanges() gets a masterlist from the server, and filters it in masterChanges().
localChanges() queries the local tables and gets the local changes, and adds them to
conflicts if they are also in the master changes.

Put(), Remove(), and Restore() accesses local and/or server to respectively send a change,
send a deletion, and restore a local record.

UpdateLibrary() is called when getting master changes. When multiple master changes are
being got, all changes are ignored except the last one, unless it's an add. The ignored
changes are still printed. If multiple changes of type " " were sent without any ending
deletes or merges, then the full change is got on the first one, and the rest are just
printed. UpdateLibrary() returns the number of records that have been updated.
*/
class
	{
	New(server = false, local? = false)
		{
		.svc = local? is true or server is false or server is ''
			? SvcCore()
			: SvcClient()
		}

	master_changes: false
	changes(table)
		{
		changes = .GetChanges(table)
		.Max_new_committed = changes.max_new_committed
		.master_changes = changes.master_changes
		.masterlist = changes.masterlist
		.conflicts = changes.conflicts
		.local_changes = changes.local_changes
		}

	GetChanges(table, asof = '')
		{
		masters = Object()
		master_changes = Object()
		svcTable = SvcTable(table)
		masterlist = .allMasterChanges(table, svcTable.GetMaxCommitted(), asof)
		max_new_committed = .masterChanges(svcTable, master_changes, masters, masterlist)

		conflicts = Object()
		local_changes = Object()
		.localChanges(svcTable, conflicts, local_changes, masters, master_changes)
		return Object(:max_new_committed, :masterlist, :master_changes, :conflicts,
			:local_changes)
		}

	allMasterChanges(table, since, to)
		{
		return .svc.AllMasterChanges(table, since, to)
		}

	masterChanges(svcTable, masterChanges, masters, masterlist)
		{
		maxLibCommitted = ''
		for master in masterlist
			{
			local = svcTable.Get(master.name)
			if not .addMasterChange?(svcTable, master, local, masterlist)
				continue
			masterChanges.Add(Object(type: master.type, name: master.name, who: master.id,
				modified: master.lib_committed, lib: svcTable.Table(),
				committed: local isnt false ? local.lib_committed : ''))
			masters[master.name] = true
			if maxLibCommitted is "" or master.lib_committed > maxLibCommitted
				maxLibCommitted = master.lib_committed
			}
		return maxLibCommitted
		}

	addMasterChange?(svcTable, change, local, masterlist)
		{
		if change.type is '-'
			{
			if local is false and
				masterlist.FindAllIf({ it.name is change.name }).Size() <= 1
				{
				svcTable.Remove(change.name, deleted:)
				return false // already deleted locally
				}
			}
		else // + or change
			{
			if local isnt false and local.lib_committed is change.lib_committed
				return false // already have it
			}
		deleted = svcTable.Get(change.name, deleted:)
		if false isnt deleted and deleted.lib_committed >= change.lib_committed
			return false // had it and deleted it
		return true
		}

	localChanges(svcTable, conflicts, local_changes, masters, master_changes)
		{
		.modifiedLocalChanges(svcTable, masters, master_changes, conflicts, local_changes)
		local_changes.Sort!({|x, y| x.type $ x.name < y.type $ y.name })
		}

	modifiedLocalChanges(svcTable, masters, master_changes, conflicts, local_changes)
		{
		lib = svcTable.Table()
		QueryApply(svcTable.ModifiedQuery() $ ' remove text')
			{ |x|
			x.type = svcTable.Deleted?(x)
				? '-'
				: x.lib_committed is ''
					? '+'
					: ' '
			x.lib = lib
			x.name = svcTable.MakeName(x, clean:)
			.addLocal(x, masters, master_changes, conflicts, local_changes)
			}
		}

	addLocal(change, masters, master_changes, conflicts, local_changes)
		{
		lib = change.lib
		name = change.name
		type = change.type
		modified = change.lib_modified
		committed = change.lib_committed
		if masters.Member?(name)
			{
			masterIdxs = master_changes.FindAllIf(
				{ it.name is name and it.lib is lib })

			masterRec = master_changes[masterIdxs.Last()]
			sends = Object()

			masterIdxs.Reverse!() // delete master_list items from last to first
			for row in masterIdxs // to avoid automatic shuffling during loop
				{
				sends.Add(Object(:name, localModified: modified,
					who: master_changes[row].who,
					masterModified: master_changes[row].modified, localType: type,
					masterType: master_changes[row].type, :lib))
				master_changes.Delete(row)
				}

			// Sort because merging needs to check if a delete
			// is the last modification
			sends.Sort!({ |x,y| x.masterModified < y.masterModified })

			conflicts.Add(Object(:name, localModified: modified, who: masterRec.who,
				masterModified: masterRec.modified, localType: type,
				:modified, masterType: masterRec.type, :lib, :sends, :committed))
			}
		else
			local_changes.Add(Object(:type, :name, :modified, :lib, :committed))
		}

	Conflicts(table)
		{
		.changes(table)
		return .conflicts
		}

	Local_changes(table)
		{
		.changes(table)
		return .local_changes
		}

	Master_list(table)
		{
		.changes(table)
		list = .masterlist.Values()
		list.Sort!({ |x,y| x.lib_committed < y.lib_committed })
		return list
		}

	Master_changes(table)
		{
		.changes(table)
		return .master_changes
		}

	Get(lib, name)
		{
		return .svc.Get(lib, name)
		}

	GetOld(lib, name, committed)
		{
		return .svc.GetOld(lib, name, committed)
		}

	GetDel(lib, name)
		{
		return .svc.GetDel(lib, name)
		}

	GetDelByDate(lib, name, committed)
		{
		return .svc.GetDelByDate(lib, name, committed)
		}

	Put(svcTable, name, id, comment, asof = false)
		{
		if asof is false
			asof = .SvcTime()
		if false is x = svcTable.Get(name)
			return false
		x.name = svcTable.MakeName(x)
		x.comment = comment
		x.path = svcTable.GetPath(x)
		if x.lib_before_hash is '' and Date?(x.lib_committed)
			x.lib_before_hash = .Hash(x.lib_before_text is ''
				? x.text
				: x.lib_before_text)

		// Can't do it all in one transaction because Put creates master table
		// and suneido doesn't allow schema changes while outstanding transactions
		if false is result = .svc.Put(svcTable.Table(), svcTable.Type, id, asof, x)
			return false
		Transaction(update:)
			{|t|
			x = svcTable.Get(name, t)
			x.lib_modified = ''
			x.lib_committed = result
			x.lib_before_hash = .Hash(x.text)
			x.lib_before_text = ''
			x.lib_before_path = ''
			svcTable.Update(x, t)
			}
		svcTable.SetMaxCommitted(result)
		return result
		}

	Remove(svcTable, name, id, comment, asof = false)
		{
		table = svcTable.Table()
		if asof is false
			asof = .SvcTime()
		if not Date?(result = .svc.Remove(table, name, id, asof, comment))
			return false
		svcTable.Remove(name, deleted:)
		svcTable.SetMaxCommitted(result)
		return result
		}

	Restore(lib, name)
		{
		svcTable = SvcTable(lib)
		svcTable.Restore(name)
		.CheckCommitted(svcTable, name, svcTable.Get(name), .Get(lib, name))
		}

	CheckCommitted(svcTable, name, local, master)
		{
		if master is false or local is false
			return

		table = svcTable.Table()
		if .Outstanding?([table], [lib: table, :name])
			return

		discrepancies = svcTable.Compare(local, master)
		if discrepancies.Has?(#text) and local.lib_modified isnt ''
			discrepancies.Remove(#text)
		if not discrepancies.Empty?()
			.recordDiscrepency(table, name, discrepancies)
		}

	Outstanding?(libraries, change = false) // change = [lib: '', name: '']
		{
		return libraries.Any?({ .hasChanges?(it, change) })
		}

	hasChanges?(lib, change)
		{
		changes = .GetChanges(lib)
		if changes.master_changes.Empty?() and changes.conflicts.Empty?()
			return false
		if change is false or lib isnt change.lib
			return true
		return changes.master_changes.Any?({ it.name is change.name }) or
			changes.conflicts.Any?({ it.name is change.name })
		}

	recordDiscrepency(table, name, discrepancies)
		{
		msgOb = Object('Unexpected discrepancy detected:')
		textMatches? = not discrepancies.Has?(#text)
		discrepancies.Map!({ '    ' $ it })
		msgOb.Add(@discrepancies)
		msgOb.Add('\nPlease verify the record\'s integrity via\nSvc : ' $
			table $ ' > Compare')
		flags = textMatches? ? MB.ICONWARNING : MB.ICONERROR
		Alert(msgOb.Join('\n'), title: 'Restore ' $ table $ ':' $ name, :flags)
		}

	GetBefore(table, name, when)
		{
		return .svc.GetBefore(table, name, when)
		}

	Get10Before(table, name, when)
		{
		return .svc.Get10Before(table, name, when)
		}

	Exists?(table)
		{
		return .svc.Exists?(table)
		}

	Compare(svcTable)
		{
		return SvcSyncClient(svcTable, .svc).Check()
		}

	MissingTest?(lib, name)
		{
		name = name.RightTrim('?')
		return .Get(lib, name $ '_Test') is false and .Get(lib, name $ 'Test') is false
		}

	UpdateLibrary(master_changes, mergeConflictFn = false)
		{
		changes = .MostRecentChanges(master_changes).Values().
			Sort!({ |x,y| x.name < y.name })
		feedbackob = Object()
		.processChanges(changes, mergeConflictFn, feedbackob)
		.ProcessFeedbackOb(feedbackob)
		return changes.Size()
		}

	processChanges(changes, mergeConflictFn, feedbackob)
		{
		svcTable = false
		for rec in changes
			{
			svcTable = .changeTable(svcTable, rec.lib)
			changeType = .getPrefix(rec.type, svcTable.Get(rec.name))
			printVal = .processChange(svcTable, rec, changeType, mergeConflictFn)
			desc = rec.name $ Opt(' ', printVal.lib_committed) $ Opt(' ', printVal.id) $
				Opt(' - ', printVal.comment)
			feedbackob.Add(Object(:changeType, lib: svcTable.Table(), name: desc,
				prefix: '<<<'))
			svcTable.SetMaxCommitted(printVal.lib_committed)
			LibUnload(rec.name)
			}
		// Ensure that the tests have to be run after getting changes
		if not changes.Empty?()
			TestRunner.RequireRun()
		}

	changeTable(svcTable, lib)
		{
		return svcTable is false or lib isnt svcTable.Table()
			? SvcTable(lib, svcEnsure:)
			: svcTable
		}

	processChange(svcTable, rec, prefix, mergeConflictFn)
		{
		printVal = []
		lib = svcTable.Table()
		name = rec.name
		modified = rec.modified
		if prefix is '+'
			printVal = .addMasterChange(svcTable, name, modified)
		else if prefix is '-'
			{
			.deleteRecord(svcTable, name)
			printVal = .GetDelByDate(lib, name, modified)
			}
		else if prefix is '#' and mergeConflictFn isnt false
			{
			mergeConflictFn(lib, name)
			printVal = .GetOld(lib, name, Date.End())
			printVal.comment = 'MERGED LOCAL WITH: ' $ printVal.comment
			}
		else
			{
			.getChange(svcTable, lib, name)
			printVal = .GetOld(lib, name, modified)
			}
		return printVal
		}

	Overwrite(records)
		{
		localOnlyRecords = Object()
		masterRecords = Object()
		for rec in records
			if false is master = .Get(rec.table, rec.name)
				localOnlyRecords.Add(rec)
			else
				{
				master.lib = rec.table
				master.type = ' '
				master.modified = Date.End()
				masterRecords.Add(master)
				}
		feedbackob = Object()
		.removeLocalOnlyChanges(localOnlyRecords, feedbackob)
		.processChanges(masterRecords, false, feedbackob)
		.ProcessFeedbackOb(feedbackob)
		}

	removeLocalOnlyChanges(localOnlyRecords, feedbackob)
		{
		svcTable = false
		localOnlyRecords.Sort!().Each()
			{
			svcTable = .changeTable(svcTable, it.table)
			svcTable.Remove(it.name, deleted: svcTable.Get(it.name, deleted:) isnt false)
			feedbackob.Add(Object(changeType: ' ', lib: it.table,
				name: it.name $ ' (only existed in the local table)',
				prefix: 'Deleted:'))
			}
		}

	MostRecentChanges(master_changes)
		{
		list = Object()
		for change in master_changes
			{
			mem = Opt(change.lib, ':') $ change.name
			if not list.Member?(mem) or list[mem].modified < change.modified
				list[mem] = change
			}
		return list
		}

	getPrefix(type, local)
		{
		local? = local isnt false
		if not local? and type is ' '
			type = '+' 	// Never got original record, treat this as a new record
		else if local? and type is '+'
			type = ' ' 	// Record exists locally, treat this as a change
		return type 	// Deletions ('-') are always treated as deletions
		}

	addMasterChange(svcTable, name, modified)
		{
		x = .GetOld(lib = svcTable.Table(), name, modified)
		svcTable.Output(x, committed:)
		svcTable.Remove(name, deleted:)
		return .GetOld(lib, name, modified)
		}

	getChange(svcTable, lib, name)
		{
		Transaction(update:)
			{|t|
			if false is old = svcTable.Get(name, t) // delete on local
				.Restore(lib, name)
			}
		Transaction(update:)
			{|t|
			old = svcTable.Get(name, t)
			x = .Get(lib, name)
			old.lib_committed = x.lib_committed
			old.lib_modified = ''
			old.text = x.text
			old.name = x.name
			old.path = x.path
			old.lib_before_hash = .Hash(x.text)
			old.lib_before_text = ''
			old.lib_before_path = ''
			old.lib_invalid_text = ''
			svcTable.Update(old, t)
			return x
			}
		}

	deleteRecord(svcTable, name)
		{
		.deleteRecordContrib(svcTable, name)
		svcTable.Remove(name)
		}

	// To stop tests from kicking in contributions
	deleteRecordContrib(table, name)
		{
		(OptContribution('SvcDeleteRecord', function (@unused) { }))(table, name)
		}

	Library?(table)
		{
		return SvcTable(table).Type is 'lib'
		}

	CheckSvcStatus()
		{
		return .svc.CheckSvcStatus()
		}

	SendLocalChanges(changes, desc, userid, asof = false,
		afterEachSendingFn = function (@unused) {})
		{
		// As of 29331, if called from SvcControl this situation is not possible.
		// However, adding handling encase this is ever called from a different location
		if changes.Empty?()
			return true
		svcTable = lib = false
		feedbackob = Object()
		for change in changes
			{
			svcTable = .changeTable(svcTable, lib = change.lib)
			type = change.type
			name = change.name
			method = type is '-' ? .Remove : .Put
			if false is asof = (method)(svcTable, name, userid, desc, asof)
				{
				feedbackob.Add(Object(changeType: type, :lib, :name, failed?:))
				return false
				}
			feedbackob.Add(Object(changeType: type, :lib, :name, prefix: '>>>'))
			afterEachSendingFn(name, lib)
			}
		.ProcessFeedbackOb(feedbackob)
		return true
		}
	ProcessFeedbackOb(ob)
		{
		for item in ob
			Print(.buildMsg(item.GetDefault('changeType', ''),
				item.GetDefault('lib', ''), item.GetDefault('name', ''),
				item.GetDefault('prefix', ''), item.GetDefault('failed?', false)))
		}
	buildMsg(changeType, lib, name, prefix = '', failed? = false)
		{
		msg = prefix $ changeType $ lib $ ':' $ name
		if failed?
			msg = '!!! ' $ msg $ ' FAILED !!!'
		return msg
		}

	Hash(text)
		{
		return Sha1(text).ToHex()
		}

	SearchForRename(table, name)
		{
		return .svc.SearchForRename(table, name)
		}

	SvcTime()
		{
		return .svc.SvcTime()
		}
	}
