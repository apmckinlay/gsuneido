// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Copy To"
	CallClass(checked)
		{
		return OkCancel(Object(this, checked), .Title, 0)
		}
	New(.srcs)
		{
		.table = .FindControl('table')
		.folder = .FindControl('folder')
		.overwrite = .FindControl('overwrite')
		}
	Controls()
		{
		return Object('Vert'
			Object('Pair' #(Static 'Destination Library')
				Object('AutoFillField', candidates: .tables(), name: 'table' width: 14))
			Object('Pair'
				#(Static 'Folder')
				Object('AutoFillField' getCandidateFn: .getFolderCandidates
					name: 'folder' width: 15))
			#(Pair (Static 'Overwrite?') (CheckBox name: 'overwrite')))
		}

	getFolderCandidates(path)
		{
		if path.Suffix?('/')
			return false

		folders = path.LeftTrim('/').Split('/')
		if folders.Empty?()
			return false

		if '' is dstLib = .table.Get()
			return false

		parent = 0
		for folder in folders[..-1]
			{
			if false is rec = Query1(dstLib, group: parent, name: folder)
				return false
			parent = rec.num
			}
		prefix = path.BeforeLast('/')
		return QueryList(dstLib $ " where group is " $ Display(parent), "name").
			Map({ Opt(prefix, '/') $ it })
		}

	tables()
		{
		return Libraries().Add('').MergeUnion(LibraryTables()).
			Difference(SvcControl.SvcExcludeLibraries)
		}

	OK()
		{
		if '' is dstLib = .table.Get()
			{
			.AlertError('Copy records', 'Please select a target library')
			return false
			}
		if .srcs.Any?({ it.lib is dstLib })
			{
			.AlertError('Copy records', 'Cannot copy records to same library')
			return false
			}
		dstFolder = .folder.Get()
		overwrite? = .overwrite.Get()
		count = 0
		for src in .srcs
			{
			res = CopyRecordsTo(src.lib, Object(src.name), dstLib, dstFolder,
				.print, overwrite?)
			if Number?(res)
				count += res
			else
				.AlertError('Copy records', res)
			}
		return count.ToWords() $
			' record(s) copied to ' $ dstLib $ Opt('/', dstFolder)
		}
	print(@args) // overridden by test
		{
		Print(@args)
		}
	}