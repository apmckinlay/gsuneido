// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	bookDeleteOrder: -999
	// SvcCore passes a record as name
	Export(lib, name, filename, path = false, delete = false, interactive = false)
		{
		exportHandler = path isnt false ? .exportBookRec : .exportLibRec

		.FileBlock(filename, mode: 'a')
			{|f|
			if false is rec = .Get_record(lib, name, path)
				{
				if not delete
					throw "LibIO Export can't get " $ name $ " from " $ lib
				rec = [:lib, :name, :path, delete:]
				}
			s = exportHandler(rec, lib, path)
			s $= .sep $ '\r\n'
			f.Write(s)
			if interactive
				Print("Exported", (path is false ? lib $ ':' : path $ "/") $ name $
					" to " $ filename)
			}
		}

	exportBookRec(rec, lib, path)
		{
		header = rec.name $ '\r\n'

		if rec.GetDefault(#delete, false) is true
			rec.order = .bookDeleteOrder

		if path !~ "^/res\>" and rec.text[-1] isnt '\n'
			rec.text $= '\r\n'

		header $= rec.order $ '\r\n' $ path $ '\r\n'
		if path =~ "^/res\>"
			rec.text = Base64.EncodeLines(rec.text) $ '\r\n'

		info = .add_book_recinfo(rec, lib)
		return header $ rec.text $ info $ '\r\n'
		}

	add_book_recinfo(rec, lib)
		{
		if rec.lib_modified is '' and rec.lib_committed is ''
			return ''
		info =  ' lib_modified: ' $ String(rec.lib_modified) $
			', lib_committed: ' $ String(rec.lib_committed)

		return BookContent.Match(lib, rec.text)
			? '<!--' $ info $ '-->'
			: '//' $ info
		}

	exportLibRec(rec, lib, path/*unused*/)
		{
		if rec.lib_invalid_text isnt ''
			throw 'LibIO did not export: ' $ lib $ ':' $ rec.name $
				' as it contains errors. Please correct and re-attempt export'

		header = rec.name $ '\r\n'
		if rec.text[-1] isnt '\n'
			rec.text $= '\r\n'

		ob = Object(lib_committed: rec.lib_committed, orig_lib: lib)

		if rec.GetDefault(#delete, false) is true
			ob.delete = true
		else
			.SetParent(ob, rec, lib)
		header $= 'librec_info: ' $ Display(ob) $ '\r\n'

		return header $ rec.text
		}

	SetParent(ob, lib_record, lib)
		{
		if false isnt parent = Query1(lib, num: lib_record.parent)
			ob.Merge(Object(parent_name: parent.name,
				parent_group: parent.group))

		// Get Path
		parent = lib_record.parent
		path = Object()
		while parent isnt 0
			{
			parentRec = Query1(lib, num: parent)
			path.Add(Object(name: parentRec.name,
				group: parentRec.group) at: 0)
			parent = parentRec.parent
			}
		ob.Merge(Object(:path))
		}

	Import(filename, lib, interactive = false, useSVCDates = false)
		{
		book? = lib isnt false and BookTable?(lib)
		importHandler = book? ? .importBookRec : .importLibRec
		_imported = Object()
		.FileBlock(filename, mode: 'r')
			{|f|
			for (n = 0; false isnt name = f.Readline(); ++n)
				{
				// ignore blank lines in file between records or after records
				if name.Blank?()
					{
					--n
					continue
					}
				importHandler(f, lib, name, interactive, useSVCDates, :filename)
				}
			}
		return n
		}
	// so we can override to use FakeFile in Tests and make use of RetryBool
	FileBlock(filename, block, mode)
		{
		File(filename, mode)
			{ |f|
			block(f)
			}
		return true
		}

	importBookRec(f, lib, name, interactive, useSVCDates)
		{
		bookInfo = .readBookInfo(f, name)
		if interactive
			{
			if not .import?(lib, name, bookInfo, _imported, bookInfo.path)
				return
			Print("Imported", bookInfo.path $ "/" $ name)
			}

		if bookInfo.order is .bookDeleteOrder
			{
			.DeleteRec(lib, SvcBook.MakeName([:name, path: bookInfo.path]),
				useSVCDates, :interactive)
			return
			}

		if false is bookRec = .Get_record(lib, name, bookInfo.path)
			{
			bookRec = Object(:name,
				text: bookInfo.text,
				order: bookInfo.order,
				path: bookInfo.path)
			.output_lib_record(lib, bookRec, bookInfo, useSVCDates, interactive)
			}
		else
			{
			if bookRec.lib_before_text is ''
				bookRec.lib_before_text = 'Order: ' $ bookRec.order $ '\r\n\r\n' $
					bookRec.text
			bookRec.text = bookInfo.text
			bookRec.order = bookInfo.order
			bookRec.path = bookInfo.path
			.update_lib_record(lib, bookRec, bookInfo, useSVCDates, interactive)
			}
		}

	readBookInfo(f, name)
		{
		order = f.Readline()
		if order isnt ""
			order = Number(order)
		path = f.Readline()

		text = ''
		while false isnt line = f.Readline()
			{
			if line is .sep
				break
			text $= line $ '\r\n'
			}

		bookInfo = Object(:order, :path, :text, svcName: path $ "/" $ name)
		.handleImportedBookRecord(bookInfo)
		return bookInfo
		}

	import?(lib, name, info, imported, path = false)
		{
		// Record has already been imported once, continue to import record as specified
		if imported.Member?(key = lib $ ':' $ name)
			return imported[key]

		import? = false
		if false is rec = .Get_record(lib, name, path)
			import? = true
		else if .skipImport?(info, rec)
			{
			.skipImport(lib, name, '(import file matches current text)')
			return false
			}
		else if rec.lib_modified is ''
			import? = true
		else
			import? = .overwrite?(lib, name)
		return imported[key] = import?
		}

	skipImport?(info, rec)
		{
		return info.text.Trim() is rec.lib_current_text.Trim() and
			info.GetDefault('order', '') is rec.GetDefault('order', '')
		}

	skipImport(lib, name, suffix = '')
		{ Print('Skipped importing: ' $ lib $ ':' $ name $ Opt(', ', suffix)) }

	overwrite?(lib, name)
		{
		overwrite? = .askOverwrite(lib, name)
		if not overwrite?
			.skipImport(lib, name)
		return overwrite?
		}

	askOverwrite(lib, name)
		{
		return YesNo('Record ' $ lib $ ':' $ name $ ' has changes.' $
			'\r\n\r\nOverwrite with import file?', 'Import')
		}

	handleImportedBookRecord(bookInfo)
		{
		.extractBookRecInfo(bookInfo)
		if bookInfo.path =~ "^/res\>"
			bookInfo.text = Base64.Decode(bookInfo.text.Tr('\r\n'))
		else
			bookInfo.text = bookInfo.text.Trim('\r\n')
		}

	extractBookRecInfo(bookInfo)
		{
		text = bookInfo.text.Trim()
		pos = text.FindLast('\r\n')
		recinfo_line = pos is false
			? text					// when content is empty
			: text[pos + 2 ..]  // 2: size of '\r\n'

		if not recinfo_line.Has?('lib_modified') or not recinfo_line.Has?('lib_committed')
			return

		modified = Date(recinfo_line.AfterFirst('lib_modified: ').BeforeFirst(','))
		committed = Date(recinfo_line.AfterFirst('lib_committed: ').BeforeFirst('-->'))
		bookInfo.lib_modified = modified is false ? "" : modified
		bookInfo.lib_committed = committed is false ? "" : committed
		bookInfo.text = pos is false ? '' : text[.. pos]
		}

	importLibRec(f, lib, name, interactive, useSVCDates, filename = false)
		{
		result = .buildImportedText(f)
		libInfo = result.librec_info.Copy()
		libInfo.svcName = name
		lib = lib isnt false ? lib : libInfo.GetDefault(#orig_lib, false)
		Assert(lib isnt: false, msg: "Import: library is not specified")
		if not TableExists?(lib)
			{
			SuneidoLog('INFO: table does not exist, skipping ' $ lib $ ':' $ name)
			return
			}

		if interactive
			{
			if not .import?(lib, name, result, _imported)
				return
			LibViewImportRestoreControl.OutputImportHistory(filename, lib, name,
				result.text)
			Print("Imported", lib $ ':' $ name)
			}

		if libInfo.GetDefault(#delete, false) is true
			{
			.DeleteRec(lib, name, useSVCDates, libInfo.lib_committed, :interactive)
			return
			}

		text = result.text

		if false is libRec = .Get_record(lib, name, false)
			{
			libRec = Object(:name, :text, parent: 0, group: -1)
			.setLibRecordParent(lib, libRec, libInfo)
			.output_lib_record(lib, libRec, libInfo, useSVCDates, interactive)
			}
		else
			{
			if libRec.lib_before_text is ''
				libRec.lib_before_text = libRec.text
			libRec.text = text
			.setLibRecordParent(lib, libRec, libInfo)
			.update_lib_record(lib, libRec, libInfo, useSVCDates, interactive)
			}
		LibUnload(name)
		}

	buildImportedText(f)
		{
		text = ''
		count = 0
		librec_info = false
		while ((false isnt line = f.Readline()) and line isnt .sep)
			{
			if count is 0 and line.Prefix?('librec_info: ')
				{
				librec_info = line.Replace('librec_info: ', '').Trim()
				++count
				continue
				}
			text $= line $ '\r\n'
			}

		try
			librec_info = librec_info is false ? Object() : librec_info.SafeEval()
		catch(err)
			{
			SuneidoLog('ERROR: LibIO.buildImportedText cannot read record info',
				params: Object(:librec_info, :err))
			librec_info = Object()
			}
		return Object(:librec_info, :text)
		}

	DeleteRec(lib, name, useSVCDates, maxCommitted = false, interactive = false)
		{
		svcTable = SvcTable(lib)
		// Always attempt to stage record for deletion
		svcTable.StageDelete(name, skipPublish: not interactive)
		if useSVCDates // If import is from svc, delete the potential staged deletion
			{
			svcTable.Remove(name, deleted:)
			if false isnt maxCommitted
				svcTable.SetMaxCommitted(maxCommitted)
			}
		}

	setLibRecordParent(lib, lib_record, librec_info)
		{
		// set parent and group
		if librec_info.Member?('path')
			{
			// check if path exists and create it if it doesn't
			num = QueryMax(lib, 'num', 0) + 1
			prev_parent = 0
			for lvl in librec_info.path
				{
				if false is rec = Query1(lib, group: prev_parent, name: lvl.name)
					{
					QueryOutput(lib, Record(num: ++num, parent: prev_parent,
						group: prev_parent, name: lvl.name))
					prev_parent = num
					}
				else
					prev_parent = rec.num
				}
			lib_record.parent = prev_parent
			}
		else // record was exported 'old way'
			{
			// don't override parent if parent folder does not exist
			if librec_info.Member?('parent_group') and
				librec_info.parent_group isnt '' and
				false isnt rec = Query1(lib, group: librec_info.parent_group,
					name: librec_info.parent_name)
				lib_record.parent = rec.num
			}
		}

	output_lib_record(lib, outputRec, infoRec, useSVCDates, interactive)
		{
		outputRec.num = QueryMax(lib, "num", 0) + 1
		.setModifiedAndCommitted(outputRec, infoRec, useSVCDates)

		svcTable = SvcTable(lib)
		if false isnt delRec = svcTable.Get(infoRec.svcName, deleted:)
			{
			// revert lib committed date if it has a pending delete in IDE mode
			if useSVCDates is false
				outputRec.lib_committed = delRec.lib_committed
			svcTable.Remove(infoRec.svcName, deleted:)
			}

		QueryOutput(lib, outputRec)
		if interactive
			svcTable.Publish(#TreeChange, force:)
		}

	update_lib_record(lib, outputRec, infoRec, useSVCDates, interactive)
		{
		.setModifiedAndCommitted(outputRec, infoRec, useSVCDates)

		QueryApply1(lib, num: outputRec.num)
			{ |rec|
			for mem in outputRec.Members()
				rec[mem] = outputRec[mem]
			rec.Update()
			}

		if interactive
			SvcTable(lib).Publish(#RecordChange, name: outputRec.name, force:)
		}

	setModifiedAndCommitted(outputRec, infoRec, useSVCDates)
		{
		if useSVCDates is true
			{
			committed = infoRec.GetDefault('lib_committed', '')
			if Date?(committed)
				outputRec.lib_committed = committed
			outputRec.lib_modified = ''
			outputRec.lib_before_text = ''
			}
		else
			{
			// do not need to change lib_committed
			outputRec.lib_modified = Date()
			}
		}

	Seperator()
		{ return .sep }
	sep: "==========================================================="
	Get_record(lib, name, path)
		{
		q = lib $ " where name is " $ Display(name)
		q $= path is false
			? " and group is -1"
			: " and path is " $ Display(path)
		return QueryFirst(q $ ' sort num')
		}
	}