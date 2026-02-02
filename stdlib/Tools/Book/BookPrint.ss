// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// export html pages of book to single large file
// Contributions by Oliver Ackermann
//TODO: factor out duplicate code with BookExportOne
//TODO: clean up res temp files - question is when?
class
	{
	New(.book, start, file, cmd)
		{
		_table = .book
		try File(file, "w")
			{|f|
			.f = f
			f.Writeline(.prefix(cmd))
			if start isnt ""
				{
				name = start.Extract("[^/]*$")
				path = start[.. -name.Size() - 1]
				if false isnt x = Query1(.book, :path, :name)
					.writeHtml(x)
				}
			.process(start)
			f.Writeline("</body></html>")
			}
		catch (err /*unused*/, "File: can't open")
			Alert("Unable to create temp file for printing. Check Permissions.",
				"Book Print", 0, MB.ICONWARNING)
		}
	prefix(cmd)
		{
		s = '<object id="printWB" width=0 height=0
			classid="clsid:8856F961-340A-11D0-A96B-00C04FD705A2"></object>
			<body onLoad="print && print() || printWB.ExecWB(' $ cmd $ ', 1)" '
		return HtmlWrap("", .book).BeforeFirst("<body") $ s
		}
	process(path) // recursive
		{
		QueryApply(.book $ " where path = " $ Display(path) $
			" sort order, name")
			{|x|
			name = x.path $ "/" $ x.name
			if .getBookAuthorize(name) is "hidden"
				continue
			.writeHtml(x)
			if (name isnt "/res")
				.process(name) // do children (if any)
			}
		}
	getBookAuthorize(name)
		{
		cont = Contributions('BookAuthorize')
		if cont.Empty?()
			return true
		func = Global(cont.Last().f)
		return func(.book, name)
		}
	writeHtml(x)
		{
		_name = x.name
		_path = x.path
		text = BookContent.Match(.book, x.text)
			? BookContent.ToHtml(.book, x.text)
			: x.text.Eval() // needs Eval
		text = Asup(text)
		.f.Writeline(HtmlWrap.Embed(text))
		}
	}
