// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// iterate thru a sequence of formats which may include generators
Generator
	{
	New(@args)
		{
		.args = args
		.i = 0
		.generators = Stack()
		}
	pushback: false
	Next()
		{
		while (false isnt (.current_item = item = .input()) and
			Instance?(item) and item.Generator?())
			{
			.generators.Push(item)
			if (false isnt (hdr = item.Header()))
				{
				hdr = _report.Construct(hdr)
				hdr.Header? = true // used by Report
				return hdr
				}
			}
		return item
		}
	input()
		{
		if (.pushback isnt false)
			{ item = .pushback; .pushback = false; return item }
		while (.generators.Count() > 0)
			if (false isnt (item = .generators.Top().Next()))
				return item // NOTE: generator does construct
			else
				.generators.Pop()
		if (.args.Member?(.i))
			return _report.Construct(.args[.i++])
		return false
		}
	Pushback(item)
		{
		Assert(.pushback is false)
		.pushback = item
		}
	PageHeader() // called by Report
		{
		hdrs = Object()
		for (g in .generators.List())
			if (false isnt (hdr = g.PageHeader()))
				hdrs.Add(hdr)
		if (false isnt (hdr = .Header()))
			hdrs.Add(hdr)
		if (hdrs.Size() is 0)
			return false
		else if (hdrs.Size() is 1)
			return hdrs[0]
		else // (hdrs.Size() > 1)
			{
			hdrs.xstretch = 1
			return hdrs.Add('Vert', at: 0)
			}
		}
	Header()
		{
		return .args.Member?('header') ? .args.header : false
		}
	PageFooter() // called by Report
		{
		ftrs = Object()
		for (g in .generators.List())
			if (false isnt (ftr = g.PageFooter()))
				ftrs.Add(ftr)
		if (ftrs.Size() is 0)
			return false
		else if (ftrs.Size() is 1)
			return ftrs[0]
		else // (ftrs.Size() > 1)
			{
			ftrs.xstretch = 1
			return ftrs.Add('Vert', at: 0)
			}
		}
	}
