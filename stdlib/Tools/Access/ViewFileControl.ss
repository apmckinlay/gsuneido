// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 550
	Ymin: 150
	CallClass(hwnd, args)
		{
		OkCancel(Object(this, args[0], args.rootDir, args.fileType, file: args.file,
			filter: args.filter, folders: args.folders),
			args.fileType, hwnd)
		}

	New(.ctrl, .rootDir, type, file = '', .filter = '', folders = #())
		{
		super(.layout(type))
		.list = .FindControl('fileList')
		.folderList = .FindControl('folderList')
		.warning = .FindControl('overLimitWarning')
		folders = .getFolders(folders)
		.folderList.SetList(folders)
		first = .getFirstDir(file, folders)
		.folderList.Set(first)
		.load_list_records(folders.Empty?() ? #() : .getDirs(first, details:))
		.list.SetReadOnly(true, grayOut: false)
		.openFile = .FindControl('openFile')
		}

	Startup()
		{
		if false is row = .list.FindRowIdx('name', .stickyFile)
			return

		.list.SelectRow(row)
		.list.ScrollRowToView(row)
		.openFile.Set(Paths.Combine(.folderList.Get(), .list.GetRow(row)['name']))
		}

	folderLimit: 6
	getFolders(folders)
		{
		folders = folders.Empty?()
			? .getDirs(filter: .filter)
			: folders
		return folders.Map({ it.Trim(`\/`) }).Sort!().Reverse!().Take(.folderLimit)
		}

	stickyFile: false
	getFirstDir(file, folders)
		{
		first = file.Blank?()
			? folders.Empty?()
				? ''
				: folders.First()
			: file =~ .filter
				? file
				: folders.First()

		.stickyFile = first.AfterFirst(`/`)
		return first.BeforeFirst(`/`)
		}

	columns: #(name, date, file_size)
	layout(type)
		{
		return Object('Record'
			Object('Vert'
			Object('Horz'
				Object('StaticText' Paths.Basename(.rootDir) $ `/`)
				'Skip'
				#('ChooseList' name: 'folderList') 'Skip'
				#('StaticText' name: 'overLimitWarning'))
			'Skip'
			Object('ListStretch',
				columns: .columns,
				alwaysHighlightSelected:,
				columnsSaveName: .Title $ ' - ' $ type,
				name: "fileList"
				)
			'Skip'
			#(Horz #(OpenFile name: 'openFile'))
			))
		}

	getDirs(subDirName = '', details = false, filter = '')
		{
		return .GetDirList(Paths.Combine(.rootDir, subDirName), details, filter)
		}

	GetDirList(dir, details = false, filter = '')
		{
		dirs = ServerEval('Dir', Paths.Combine(dir, '*'), :details, files: details)
		return details
			? dirs
			: dirs.RemoveIf({ not it.Suffix?(`/`) or it !~ filter })
		}

	Record_NewValue(field, value)
		{
		if field is 'folderList'
			.load_list_records(.getDirs(subDirName: value, details:))
		}

	listLimit: 500
	load_list_records(data)
		{
		.list.DeleteAll()
		.warning.Set('')
		if data.Empty?()
			return

		counter = 0
		data.Each({ it.file_size = .formatSize(it.size) })
		for rec in data.Sort!({ |x,y| x.date > y.date })
			{
			if ++counter > .listLimit
				{
				.warning.Set('Max file list reached (' $ .listLimit $ ')')
				return
				}

			if Object?(rec)
				.list.AddRow(rec)
			}
		}

	formatSize(size)
		{
		if size is 0 /*= empty file*/
			return '0 kb'

		if size < 1024 /*= 1 kb*/
			return '1 kb'

		return ReadableSize(size)
		}

	List_DoubleClick(row, unused)
		{
		if row isnt false
			.Window.Result(Paths.Combine(.folderList.Get(), .list.GetRow(row)['name']))

		return false
		}

	List_SingleClick(row, unused)
		{
		if row is false
			return 1

		.openFile.Set(Paths.Combine(.folderList.Get(), .list.GetRow(row)['name']))
		return 0
		}

	OK()
		{
		file = .openFile.Get()
		if file.Blank?()
			{
			.AlertWarn('View File', 'No File Selected')
			return false
			}

		return file
		}
	}