// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TreeViewControl
	{
	New(model)
		{
		super(true, TVS.TRACKSELECT)
		.model = model
		.load()
		}
	TVN_BEGINRDRAG()
		{
		// Disallow right button dragging
		return 0
		}
	TVN_SELCHANGED(lParam)
		{
		tv = NMTREEVIEW(lParam)
		.Send("Goto", .getpath(tv.itemNew.hItem))
		return 0
		}
	RBUTTONUP()
		{
		return 0
		}
	load(path = "", parent = 0)	// Recursive
		{
		// Load root items (chapters)
		for (child in .model.Children(path))
			{
			x = .AddItem(parent, child.name, 0, container?: true)
			.load(child.path $ '/' $ child.name, x)
			}
		}
	getpath(item)
		{
		// Returns complete path to item including name
		path = ""
		for (; item isnt 0; item = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.PARENT, item))
			path = "/" $ .GetName(item) $ path
		return path
		}
	GotoPath(path)
		{
		// Go to a given item; if item is not visible, expand tree to make it visible
		path = path.Split("/")
		list = .GetChildren(TVI.ROOT)
		item = false
		for (pathitem in path)
			{
			if (pathitem is "")
				continue
			item = false
			for (listitem in list)
				if (pathitem is .GetName(listitem))
					{
					item = listitem
					break
					}
			if (item is false)
				return false
			list = .GetChildren(item)
			}
		.SelectItem(item)
		return true
		}
	SelectItem(item)	// Override
		{
		if (item isnt SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.CARET, 0))
			super.SelectItem(item)
		SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.FIRSTVISIBLE, item)
		}
	Children?(item)
		{ return .model.Children?(.getpath(item)) }
	// The following two functions are not necessary
	// because drag functionality does not exist in TreeView
	//Container?(item)
	//	{ return true }
	//Static?(item))
	//	{return true }
	}
