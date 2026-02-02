// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
TreeViewControl
	{
	Name: 'TreeView'
	ComponentName: 'BookTree'
	ComponentArgs: #()
	New(.model)
		{
		super(true)
		.load()
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

	TVN_SELCHANGED(oldSelect, newSelect)
		{
		super.TVN_SELCHANGED(oldSelect, newSelect)
		.Send("Goto", .getpath(newSelect))
		return 0
		}

	getpath(item)
		{
		// Returns complete path to item including name
		path = ""
		for (; item isnt 0; item = .GetParent(item))
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
				{
				if (pathitem is .GetName(listitem))
					{
					item = listitem
					break
					}
				}
			if (item is false)
				return false
			list = .GetChildren(item)
			}
		.SelectItem(item)
		return true
		}
	}