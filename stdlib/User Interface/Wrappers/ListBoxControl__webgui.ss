// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name:		"ListBox"
	ComponentName: "ListBox"
	Xmin: 		100
	Xstretch: 	1
	Ymin: 		100
	Ystretch: 	1

	New(@args)
		{
		.sort = args.GetDefault("sort", false)

		if args.Size(list:) is 1 and Object?(args[0])
			args = args[0]

		.items = Object()
		for s in args.Values(list:)
			.AddItem(s)

		.ComponentArgs = Object(
			multicolumn: args.GetDefault(#multicolumn, false),
			font: args.GetDefault(#font, ''),
			size: args.GetDefault(#size, ''),
			weight: args.GetDefault(#weight, ''))
		}

	n: 0
	maxCharacters: 250
	AddItem(s, n = false)
		{
		s = String(s)[.. .maxCharacters]
		i = .insertString(s)
		if (n is false)
			.SetData(i, .n++)
		else
			.SetData(i, n)
		}

	insertString(s, i = false)
		{
		item = Object(:s)
		at = i isnt false
			? i
			: .sort is false
				? .items.Size()
				: .items.BinarySearch(item, By(#s))
		.items.Add(item, :at)
		.Act('InsertItem', s, at)
		return at
		}

	InsertItem(s, i)
		{
		s = String(s)[.. .maxCharacters]
		i = .insertString(s, i)
		}

	curSel: false
	LBN_DBLCLK(.curSel)
		{
		.Send("ListBoxDoubleClick", .curSel)
		return 0
		}

	SELCHANGE(.curSel)
		{
		.Send("ListBoxSelect", .curSel)
		return 0
		}

	GetCount()
		{
		return .items.Size()
		}

	Count()
		{
		return .GetCount()
		}

	DeleteItem(i)
		{
		if .curSel is i
			.curSel = false
		if .items.Member?(i)
			{
			// WARNING: this will mess up .n
			.items.Delete(i)
			.Act('DeleteItem', i)
			}
		}

	DeleteAll()
		{
		for (i = .GetCount() - 1; i >= 0; --i)
			.DeleteItem(i)
		.n = 0
		}

	SetCurSel(i)
		{
		if i isnt -1 and not .items.Member?(i)
			return
		if i is -1
			.curSel = false
		else
			.curSel = i
		.Act('SetCurSel', i)
		}

	GetCurSel()
		{
		return .curSel is false ? LB.ERR : .curSel
		}

	GetSelected()
		{
		return .GetData(.GetCurSel())
		}

	Get()
		{
		return .GetText(.GetCurSel())
		}

	SetData(i, n)
		{
		if i is -1
			.items.Each({ it.n = n })
		else if .items.Member?(i)
			.items[i].n = n
		}

	GetData(i)
		{
		if i is LB.ERR or not .items.Member?(i)
			return LB.ERR
		return .items[i].GetDefault(#n, LB.ERR)
		}

	GetText(i)
		{
		if i is LB.ERR or not .items.Member?(i)
			return ''
		return .items[i].s
		}

	SetColumnWidth(@unused) { }

	FindString(text)
		{
		if false is i = .items.FindIf({ it.s.Lower().Prefix?(text.Lower()) })
			return -1
		return i
		}

	CONTEXTMENU(x, y, i)
		{
		.SetCurSel(i)
		.Send("ListBox_ContextMenu", x, y)
		}
	}