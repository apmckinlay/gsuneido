// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Import"
	CallClass(query)
		{
		OkCancelModal(Object(this, query), title: .Title)
		}

	New(.query)
		{}

	Controls:
		(Record
			(Vert
				(Pair
					(Static "File name:")
					(OpenFile title: 'Import' name: file))
				Skip
				(Pair
					(Static 'Format:')
					(ChooseList listField: import_formats, selectFirst:,
						name: format))
				)
			)
	OK()
		{
		fn = false
		data = .Data.Get()
		for c in Plugins().Contributions("ImportExport", "import_formats")
			if c.name is data.format
				fn = Global(c.impl)
		try
			{
			if fn is false
				throw 'unrecognized format: ' $ data.format
			fn(data.file, .query, header:)
			return true
			}
		catch (x)
			{
			Alert("Error during import: " $ x, 'Global Import',
				flags: MB.ICONERROR)
			return false
			}
		}
	}
