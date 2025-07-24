// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
/*
SvcModel encapsulates the Svc instance. Through Svc, SvcCore can be accessed (through
SvcClient and SvcServer in client/server mode, or directly in standalone mode), which runs
on the server. This class is called any time information or action is needed from the
server and it simplifies the requests with an interface.

Notable methods:

SetTable() gets all changes from the server. svc_all_changes loops through all
LibraryTables.

If there is more than one master change, GetPrevMasterRec() returns the master change
before the current one. This can represent the previous item in the list, or the last
change got locally.

GetMasterFromLocal() returns the master change at the lib_committed date of the local
record. It represents the state of the record before any of the oustanding local/master
changes were made.

GetMergedRec() returns the text of the record if the local were to be merged with the
master.

MoveConflict() returns the necessary items to insert into the master window from a moved
conflict. It does not modify any code. merge_conflict() does the actual merge.

UpdateLocalModified() updates the currently selected local change modified date. This is
called when any item in the local list is selected. A selection change refreshes the code
window, so the date should also be refreshed, because the user is up-to-date on the
change.
*/
class
	{
	settings: false
	local_changes: #()
	conflicts: #()
	master_changes: #()
	svc: false
	Default(@args)
		{
		return .svc[args[0]](@+1 args)
		}

	SetSettings(settings = false)
		{
		.settings = settings
		}

	SetTable(table)
		{
		.Clear()

		if table is false or table is "" or .settings is false
			return

		.svc = .initSvc()

		changes = .GetChanges(table)
		.local_changes = changes.local_changes
		.conflicts = changes.conflicts
		.master_changes = changes.master_changes
		}

	initSvc()
		{
		return Svc(server: .settings.svc_server, local?: .settings.svc_local?)
		}

	GetChanges(table)
		{
		if table is "svc_all_changes"
			{
			allChanges = Object(local_changes: Object(), master_changes: Object(),
				conflicts: Object())
			libraries = SvcControl.SvcLibraryTables().Remove(@SvcDisabledLibraries())
			for library in libraries
				{
				changes = .svc.GetChanges(library)
				allChanges.local_changes.MergeUnion(changes.local_changes)
				allChanges.master_changes.MergeUnion(changes.master_changes)
				allChanges.conflicts.MergeUnion(changes.conflicts)
				}
			return allChanges
			}
		else
			return .svc.GetChanges(table)
		}

	Clear()
		{
		.local_changes = Object()
		.conflicts = Object()
		.master_changes = Object()
		.svc = false
		}

	LocalChangeNeedsUpdate?(table, log = false)
		{
		loaded = .local_changes.Map(Display)
		pending = .GetChanges(table).local_changes.Map(Display)
		if not loaded.EqualSet?(pending)
			{
			if log is true
				SuneidoLog('Local Change Needs Update',
					params: Object(:table,
						loaded: loaded.Difference(pending),
						pending: pending.Difference(loaded)))
			return true
			}
		return false
		}

	Getter_LocalChanges()
		{
		return .local_changes
		}

	Getter_MasterChanges()
		{
		return .master_changes
		}

	Getter_Conflicts()
		{
		return .conflicts
		}

	GetConflict(name, lib)
		{
		return .conflicts.FindOne({ it.name is name and it.lib is lib })
		}

	GetLocalRec(lib, name, t = false, deleted = false)
		{
		table = SvcTable(lib)
		validRec = table.Get(name, t, :deleted)
		if table.Type is #lib and validRec isnt false
			validRec.text = validRec.lib_current_text
		return validRec
		}

	GetMasterRec(lib, name, lib_committed = false, delete? = false)
		{
		if lib_committed isnt false and delete? is false
			return .svc.GetOld(lib, name, lib_committed)
		else if lib_committed isnt false and delete?
			return .svc.GetDelByDate(lib, name, lib_committed)
		else if lib_committed is false and delete?
			return .svc.GetDel(lib, name)
		else
			return .svc.Get(lib, name)
		}

	GetPrevMasterRec(lib, name, lib_committed)
		{
		if lib_committed is false
			return false
		table = SvcTable(lib)
		list = .master_changes.Filter({
			it.name is name and
			it.lib is lib and
			it.lib is table.Table() and
			it.modified < lib_committed })
		if list.Empty?()
			return false
		rec = list.MaxWith({ it.modified })
		return .GetMasterRec(rec.lib, name, rec.modified)
		}

	GetPrevDefinition(name, curTable)
		{
		if curTable is 'stdlib'
			return false
		svcTable = SvcTable(curTable)
		if svcTable.Type isnt 'lib'
			return false
		webgui? = curTable.Suffix?('webgui')
		for lib in .libraries()
			{
			if webgui? and not lib.Suffix?('webgui')
				continue
			if lib is curTable
				return false
			svcTable = SvcTable(lib)
			if false isnt rec = svcTable.Get(name)
				return rec.Merge([table: lib])
			}
		return false
		}

	libraries()
		{
		return Libraries()
		}

	GetMasterFromLocal(lib, name)
		{
		if false is local = .GetLocalRec(lib, name)
			return .svc.GetOld(lib, name, SvcTable(lib).Get(name, deleted:).lib_committed)
		if false is master = .svc.GetOld(lib, name, local.lib_committed)
			return []
		return master
		}

	GetMergedRec(lib, name, local, master)
		{
		if local is false
			return false

		base = .GetMasterFromLocal(lib, name)
		three = Diff.Three(base.text.Lines(), master.text.Lines(), local.text.Lines())
		merged = Object()
		for line in three
			{
			if not line[1].Prefix?('-')
				merged.Add(line[2])
			}
		merged = merged.Join('\r\n')

		return Object(:merged, :base)
		}

	Getter_Svc()
		{
		return .svc
		}

	MoveConflict(name, lib, merge? = false)
		{
		conflict = .GetConflict(name, lib)
		recAdded = Object()
		for send in conflict.sends
			{
			rec = Object(:name, type: send.masterType, modified: send.masterModified,
				:lib, who: send.who)
			.master_changes.Add(rec)
			recAdded.Add(rec)
			}

		if merge? and conflict.sends.Last().masterType isnt '-'
			{
			type = '#'
			rec = Object(:name, :type, modified: Timestamp(), :lib)
			.master_changes.Add(rec)
			recAdded.Add(rec)
			}
		.conflicts.RemoveIf({ it.name is name and it.lib is lib})
		return recAdded
		}

	Restore(name, lib, type)
		{
		if type not in ('-', ' ', '%', '+')
			return false
		.svc.Restore(lib, name)
		return true
		}

	UpdateLibrary(lib)
		{
		if not .conflicts.Empty?()
			return -1
		if .master_changes.Empty?()
			return 0

		changes = lib is false
			? .master_changes
			: .master_changes.Filter({ it.lib is lib })

		return .svc.UpdateLibrary(changes, .merge_conflict)
		}

	merge_conflict(lib, name)
		{
		Transaction(update:)
			{ |t|
			if false is x = .GetLocalRec(lib, name, t)
				return
			master = .GetMasterRec(lib, name)
			local_lines = x.text.Entab().Lines()
			master_lines = master.text.Entab().Lines()
			// if there is conflict between new local record (never been committed)
			// and same named record on SVC, merging will cause program error
			// if only trying to find next newest (.GetOld) record from SVC
			// instead, use master.lib_committed
			base = not x.Member?('lib_committed')
				? []
				: .GetMasterRec(lib, name, x.lib_committed)
			if base.path isnt master.path
				x.path = master.path
			base_lines = base.text.Entab().Lines()
			text = Diff.Merge(base_lines, master_lines, local_lines).Join('\r\n')
			.mergeUpdate(lib, x, master, text, t)
			}
		}

	mergeUpdate(lib, local, master, text, t)
		{
		table = SvcTable(lib)
		if table.Type is 'lib'
			{
			validMerge? = CodeState.Valid?(table.Table(), [name: local.name, :text])
			local.lib_invalid_text = validMerge? ? '' : text
			text = validMerge? ? text : master.text
			}
		local.text = text
		local.name = table.MakeName(local)
		local.lib_committed = master.lib_committed
		local.lib_modified = Date()
		local.lib_before_hash = Svc.Hash(master.text)
		local.lib_before_text = master.text
		table.Update(local, t)
		}

	SvcExists?(table)
		{
		return .svc.Exists?(table)
		}

	SvcTime()
		{
		return .svc.SvcTime()
		}

	SendLocalChanges(changes, desc, userid, asof = false,
		afterEachSendingFn = function (@unused) {})
		{
		return .svc.SendLocalChanges(changes, desc, userid, asof,
			{ |name, lib|
			.local_changes.RemoveIf({ it.name is name and it.lib is lib })
			afterEachSendingFn(name, lib)
			})
		}

	SendLocalChangesFromAll(allChanges, desc, userid, asof = false,
		afterEachSendingFn = function (@unused) {})
		{
		if false is libsSendingTo = .NoMasterChangesInLibs?(allChanges)
			return false

		libsSendingTo.Each(
			{ |lib|
			changes = allChanges.Filter({ it.lib is lib })
			if not .SendLocalChanges(:changes, :desc, :userid, :asof, :afterEachSendingFn)
				return false
			})
		return true
		}

	NoMasterChangesInLibs?(changes)
		{
		librariesSendingTo = Object()
		changes.Each({ librariesSendingTo.AddUnique(it.lib) })
		return .svc.Outstanding?(librariesSendingTo) ? false : librariesSendingTo
		}

	SvcCompare(table)
		{
		return SvcCompare(.svc, table)
		}

	UpdateLocalModified(name, lib, changeOb, dateModified)
		{
		if false isnt local = .GetLocalRec(lib, name)
			{
			localDate = local.lib_modified
			changeOb.FindOne({ it.name is name and it.lib is lib }).modified = localDate
			return localDate
			}
		return dateModified
		}

	Library?(table)
		{
		return .svc.Library?(table)
		}
	}
