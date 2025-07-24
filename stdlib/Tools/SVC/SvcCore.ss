// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/*
Suneido Version Control

SvcCore is run on the server (or local in standalone mode) and updates/returns info from
the master tables and server.

Note: AllMasterChanges() returns all master changes since the last get.

In addition to writing a send to the master table, Put() also creates the lib_after_hash
for the change.
*/
class
	{
	Exists?(table)
		{ return TableExists?(table $ '_master') }

	AllMasterChanges(table, since, to = '')
		{
		if not .Exists?(table)
			return Object()
		return QueryAll(table $ '_master where lib_committed > ' $ Display(since) $
			(to is '' ? '' : ' where lib_committed < ' $ Display(to)) $ ' remove text')
		}

	ListAllChanges(since, dir = false)
		{
		list = Object().Set_default(Object())
		if dir isnt false
			.ensureDir(dir)
		for master_table in .masterTableList()
			{
			tablename = master_table.RemoveSuffix('_master')
			reclist = .tableChanges(since, master_table)

			for recname in reclist.Members()
				{
				if dir isnt false
					.exportRecord(tablename, recname, reclist[recname], dir)
				if reclist[recname].type is '-'
					recname = '-' $ recname
				list[tablename].Add(recname)
				}
			}
		return list
		}

	ensureDir(dir)
		{
		EnsureDir(dir)
		}

	// refactored so we can test
	masterTableList()
		{
		return QueryList('tables where table.Suffix?("_master")', 'table')
		}

	tableChanges(since, table, queryextra = '')
		{
		master_changes = QueryAll(table $ ' where lib_committed > ' $ Display(since) $
			queryextra $ ' extend modified = lib_committed sort lib_committed')
		return Svc.MostRecentChanges(master_changes)
		}

	// refactored so we can test
	exportRecord(filename, recname, rec, dir)
		{
		if rec.type is '-'
			{
			.addToDeletes(dir, filename, recname)
			return
			}

		if false is type = .MasterType(filename)
			return

		// need to tweek this if a book
		if type is "book"
			rec = .bookInfo(rec)

		.export(filename, rec, dir $ '/' $ filename, type is "book" ? rec.path : false)
		}

	addToDeletes(dir, filename, recname)
		{
		AddFile(dir $ '/deleted.txt', filename $ ', ' $ recname $ '\r\n')
		}

	MasterType(table)
		{
		type = false
		columns = QueryColumns(table $ '_master')
		if columns.Has?('svc_library')
			type = 'lib'
		else if columns.Has?('svc_book')
			type = 'book'
		return type
		}

	bookInfo(rec)
		{
		rec.name = rec.name.AfterFirst(rec.path $ '/')
		if rec.text.Prefix?('Order:')
			{
			ob = rec.text.Lines()
			rec.order = ob[0].AfterFirst('Order: ')
			rec.text = rec.text.AfterFirst(ob[0] $ '\n\n')
			return rec
			}
		return rec
		}

	export(filename, rec, exportFile, path)
		{
		SVCLibIO.Export(filename, rec, exportFile, :path)
		}

	Get(table, name)
		{
		if not .Exists?(table)
			return false
		x = QueryLast(table $ '_master where name = ' $ Display(name) $
			' sort lib_committed')
		return x is false or x.type is '-' ? false : x
		}

	GetOld(table, name, committed)
		{
		if not .Exists?(table)
			return false
		return QueryLast(table $ '_master
				where name = ' $ Display(name) $ ' and
				lib_committed <= ' $ Display(committed) $ ' and
				type isnt "-"
				sort lib_committed')
		}

	GetDel(table, name)
		{
		if not .Exists?(table)
			return false
		x = QueryLast(table $ '_master where name = ' $ Display(name) $
			' sort lib_committed')
		return x is false or x.type isnt '-' ? false : x
		}

	GetDelByDate(table, name, committed)
		{
		if not .Exists?(table)
			return false
		return QueryLast(table $ '_master
			where name = ' $ Display(name) $ ' and
			lib_committed <= ' $ Display(committed) $ ' and
			type is "-"
			sort lib_committed')
		}

	Put(table, type, id, asof, rec) // used for add & update
		{
		if .lastCommit(table) > asof
			return false // someone else has sent changes

		rec.id = id
		rec.type = .Get(table, rec.name) is false ? '+' : ' '
		rec.lib_committed = Timestamp()
		rec.lib_after_hash = Svc.Hash(rec.text)
		if true isnt result = .EnsureMaster(master = table $ '_master', type)
			return result

		QueryOutput(master, rec)

		.svcHooks(master, 'Put')
		SvcSyncServer.ResetCache()
		return rec.lib_committed
		}

	lastCommit(table)
		{
		return .Exists?(table)
			? QueryMax(table $ '_master', 'lib_committed', Date.Begin())
			: Date.Begin()
		}

	EnsureMaster(master, type)
		{
		if not #(lib, book).Has?(type)
			return 'ERR ensuring master table: ' $ master $ ', invalid type: ' $ type
		typeColumn = type is 'lib' ? 'svc_library' : 'svc_book'
		Database('ensure ' $ master $
			' (name, path, text, lib_committed, id, comment, type, lib_before_hash,
				lib_after_hash,'  $ typeColumn $')
			key(lib_committed, name) key(name, lib_committed)
			index(lib_after_hash, lib_committed)')
		return true
		}

	svcHooks(master, type)
		{
		for f in Contributions('SvcCommitHooks')
			f(master, type)
		}

	Remove(table, name, id, asof, comment)
		{
		if not .Exists?(table)
			return 'ERR table: ' $ table $ ' does not exist in Svc, remove ignored'
		if .lastCommit(table) > asof
			return false // Someone else has sent changes, stop SvcClient send

		rec = [:name, lib_committed: Timestamp(), :id, :comment, type: '-']
		QueryOutput(master = table $ '_master', rec)
		.svcHooks(master, 'Remove')
		SvcSyncServer.ResetCache()

		return rec.lib_committed
		}

	GetBefore(table, name, when)
		{
		return .Exists?(table)
			? QueryFirst(.getbefore_query(table, name, when))
			: false
		}

	Get10Before(table, name, when)
		{
		list = Object()
		if .Exists?(table)
			QueryApply(.getbefore_query(table, name, when))
				{ |x|
				list.Add(x)
				if list.Size() >= 10 /*= requested list size */
					break
				}
		return list
		}

	getbefore_query(table, name, when)
		{
		return table $ '_master
			where name = ' $ Display(name) $
			' and lib_committed < ' $ Display(when) $
			' project name, id, comment, lib_committed, type
			sort reverse lib_committed'
		}

	GetChecksums(table, from, to)
		{
		return .Exists?(table)
			? SvcSyncServer(table $ '_master', from, to)
			: Object()
		}

	// this should only happen when library doesn't have folder,
	// and last changes were only deletes
	OnlyDeletedChangesBetween?(table, localMaxLibCommitted, savedMaxLibCommitted)
		{
		if not .Exists?(table)
			return false
		master = table $ '_master'
		return QueryEmpty?('(' $ master $
			' where lib_committed > ' $ Display(localMaxLibCommitted) $
			' and lib_committed <= ' $ Display(savedMaxLibCommitted) $
			' summarize name, max lib_committed
				rename max_lib_committed to lib_committed)
			join by(name, lib_committed) ' $ master $
			' where type isnt "-"')
		}

	SearchForRename(table, name)
		{
		if not .Exists?(table)
			return Object()
		master = table $ '_master'
		masterItem = QueryFirst(master $ ' where name is ' $ Display(name) $
			' sort lib_committed')
		if masterItem is false or "" is beforeHash = masterItem.lib_before_hash
			return Object()
		lib_committed = masterItem.lib_committed
		return QueryList(master $ ' where name isnt ' $ Display(name) $
			' and lib_committed < ' $ Display(lib_committed) $
			' and lib_after_hash is ' $ Display(beforeHash), 'name')
		}

	SvcTime()
		{
		return Timestamp()
		}

	CheckSvcStatus()
		{
		return ""
		}

	CurrentLibraryRecordsQuery(masterTable)
		{
		return masterTable $ " summarize name, max lib_committed
			rename max_lib_committed to lib_committed
			join by(name, lib_committed) " $ masterTable $ " where type isnt '-'
			sort name"
		}

	LibraryChecksumFromMaster(table)
		{
		masterTable = table $ "_master"
		query = .CurrentLibraryRecordsQuery(masterTable)
		ChecksumLibraries.Calc_cksums(masterTable, query, false, cksumOb = Object())
		if cksumOb.Empty?()
			return cksumOb

		cksum = cksumOb[0]
		cksum.lib = table
		return cksum
		}
	}
