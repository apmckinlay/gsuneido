// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Document_Builder
	{
	New(.filename, fileCl = false)
		{
		.file = fileCl is false ? .createFile() : fileCl
		.lineOb = Object()
		}

	createFile()
		{
		return File(.filename, mode: 'w')
		}

	PlainText?() { return true }

	Process(fmt)
		{
		if not Instance?(fmt)
			return
		str = fmt.ExportCSV().Trim()
		if str isnt ''
			.file.Writeline(str)
		}

	Finish(status)
		{
		.file.Close()
		if status is ReportStatus.NODATA
			DeleteFile(.filename)
		return status
		}
	}
