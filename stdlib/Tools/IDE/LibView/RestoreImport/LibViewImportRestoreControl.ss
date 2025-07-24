// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Restore Imported Records"
	Table: 'libview_import_history'
	CallClass(hwnd)
		{
		ToolDialog(hwnd, Object(this), .Title)
		}

	New()
		{
		.Ensure()
		.cleanupRestoredRecords()
		.setData(.FindControl('List'))
		}

	Controls: #(Vert
		('ListStretch',
			columns: #(lvimport_filename, lvimport_num, lvimport_lib, lvimport_name),
			noHeaderButtons:, stretchColumn: lvimport_name,
			columnsSaveName: 'Restore Imported Records')
		'Skip'
		(Static 'Right-click a filename to restore; ' $
			'right-click a record to go to definition')
		)

	cleanupRestoredRecords()
		{
		QueryApplyMulti(.Table, update:)
			{
			if not TableExists?(it.lvimport_lib) or QueryEmpty?(it.lvimport_lib $
				' where name is ' $ Display(it.lvimport_name) $ ' and lib_modified > ""')
				it.Delete()
			}
		}

	setData(list)
		{
		curfilename = ''
		QueryApply(.Table $ ' extend lvimport_date
			sort lvimport_date, lvimport_filename, lvimport_lib, lvimport_name')
			{
			if it.lvimport_filename isnt curfilename
				{
				list.AddRow(Record(lvimport_filename: it.lvimport_filename))
				curfilename = it.lvimport_filename
				}
			list.AddRow(Record(lvimport_num: it.lvimport_num,
				lvimport_lib: it.lvimport_lib, lvimport_name: it.lvimport_name))
			}
		}
	List_ContextMenu(x, y, source)
		{
		sel = source.GetSelection()
		if sel.Empty?()
			return

		rec = source.GetRow(sel[0])
		if rec.lvimport_filename isnt ''
			ContextMenu(#('Restore Import')).ShowCall(this, x, y)
		else
			ContextMenu(#('Go To Definition')).ShowCall(this, x, y)
		}
	On_Context_Restore_Import()
		{
		if false is selrec = .getSelectedRec()
			return
		list = QueryAll(.Table $
			" where lvimport_filename is " $ Display(selrec.lvimport_filename))
		changed = Object()
		for rec in list
			{
			cur = Query1(rec.lvimport_lib, name: rec.lvimport_name, group: -1)
			curchksum = FormatChecksum(Adler32(cur.text))
			if curchksum isnt rec.lvimport_chksum
				changed.Add(rec.lvimport_lib $ ':' $ rec.lvimport_name)
			}
		if not changed.Empty?()
			{
			.AlertError('Records Changed',
				'The following records have been changed since importing:\n' $
				changed.Join('\n'))
			return
			}
		for rec in list
			{
			SvcTable(rec.lvimport_lib).Restore(rec.lvimport_name)
			QueryDo('delete ' $ .Table $
				' where lvimport_num is ' $ Display(rec.lvimport_num))
			}
		listctrl = .FindControl('List')
		listctrl.Clear()
		.setData(listctrl)
		}
	On_Context_Go_To_Definition()
		{
		if false is rec = .getSelectedRec()
			return
		GotoLibView(rec.lvimport_lib $ ':' $ rec.lvimport_name)
		}
	getSelectedRec()
		{
		list = .FindControl('List')
		sel = list.GetSelection()
		if sel.Empty?()
			return false

		return list.GetRow(sel[0])
		}

	List_DoubleClick(@unused)
		{ return false }

	Ensure()
		{
		Database("ensure " $ .Table $
			" (lvimport_num, lvimport_filename, lvimport_lib, lvimport_name,
				lvimport_chksum)
			key(lvimport_num)")
		}

	OutputImportHistory(filename, lib, name, text)
		{
		.Ensure()
		QueryDo('delete ' $ .Table $
			' where lvimport_filename is ' $ Display(filename) $
				' and lvimport_name is ' $ Display(name))
		rec = Record(lvimport_num: Timestamp(), lvimport_filename: filename,
			lvimport_lib: lib, lvimport_name: name)
		rec.lvimport_chksum = FormatChecksum(Adler32(text))
		QueryOutput(.Table, rec)
		}
	}