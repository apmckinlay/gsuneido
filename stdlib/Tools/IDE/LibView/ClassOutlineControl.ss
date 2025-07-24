// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// contributed by Santiago Ottonello
Controller
	{
	Name: 'ClassOutline'

	New()
		{
		.tree = .FindControl('classoutlinetree')
		.sort = .FindControl('sort')
		.public = .FindControl('public')
		.private = .FindControl('private')
		.meth = .FindControl('meth')
		.data = .FindControl('data')
		.tree.ResetTheme()
		}

	Controls: (Vert,
		(Horz,
			Fill,
			(ButtonToggle, 'a-z', set: true, tip: 'sort members', name: sort),
			(Fill, fill: .15),
			(AtLeastOne, (Horz,
				(ButtonToggle, 'pub', set: true, tip: 'show public members',
					name: public),
				(Fill fill: .05),
				(ButtonToggle, 'pri', set: true, tip: 'show private members',
					name: private),
				)),
			(Fill, fill: .15),
			(AtLeastOne, (Horz,
				(ButtonToggle, 'fun', set: true, tip: 'show methods', name: meth),
				(Fill, fill: .05),
				(ButtonToggle, 'dat', set: true, tip: 'show data members', name: data)
				)),
			name: 'mainHorz'
			)
		(Skip, 2),
		(TreeView, name: 'classoutlinetree', xmin: 0 readonly:)
		)

	Reset()
		{
		.items = false
		.rootLabel = ''
		}

	items: false
	rootLabel: ''
	initialSet?: false
	Set(data)
		{
		.lastdata = data
		try
			.set(data)
		catch (e)
			{
			.items = []
			.rootLabel = ''
			.treeReset()
			Alert(e)
			}
		}

	set(data)
		{
		sourceCode = .getSourceCode(data)
		switch type = LibRecordType(sourceCode)
			{
		case 'class':
			.buttons(true)
			rootLabel = .buildRootLabel()
			items = .filterItems(ClassHelp.ClassMembers(sourceCode))
		case 'function':
			.buttons(false)
			rootLabel = .buildRootLabel() $ ' Referenced:'
			items = .functionItems(sourceCode)
			items.Sort!().Unique!()
		case 'object' :
			.buttons(true)
			rootLabel = .buildRootLabel()
			items = .filterItems(ClassHelp.AllObjectMembers(sourceCode))
		default :
			.buttons(false)
			rootLabel = ''
			items = []
			}
		if items is .items and rootLabel is .rootLabel
			return
		.items = items
		.rootLabel = rootLabel
		.loadTree(type, sourceCode)
		.initialSet? = true
		}

	treeReset()
		{
		if .initialSet?
			.tree.Reset()
		}

	filterItems(items)
		{
		if .sort.Get()
			items.Sort!()
		return items.Filter()
			{
			not .data.Get() and it.Suffix?(':')
				? false
				: not .meth.Get() and not it.Suffix?(':')
					? false
					: .caseFilter(it)
			}
		}

	caseFilter(item)
		{
		return not .public.Get() and item[0].Upper?()
			? false
			: .private.Get() or not item[0].Lower?()
		}

	getSourceCode(data)
		{
		if String?(data)
			sourceCode = data
		else if Object?(data)
			sourceCode = data.GetDefault(#text, '')
		else
			sourceCode = ''
		return sourceCode
		}

	loadTree(type, sourceCode)
		{
		.treeReset()
		if .items.Empty?()
			return
		root = .tree.AddItem(TVI.ROOT, .rootLabel, 0, true)
		for item in .items
			.tree.AddItem(root, item)
		if type is 'class' and .Send('ClassOutline_SkipHierarchy?') isnt true
			.loadHierarchy(.Send('CurrentTable'), ClassHelp.SuperClass(sourceCode))
		.tree.ExpandItem(root)
		// gives tree an inital select, so it won't trigger selection when click on '+'
		.tree.SelectItem(root)
		}

	UpdateRootLabel()
		{
		.rootLabel = .buildRootLabel()
		item = .tree.GetSelectedItem()
		parent = .tree.GetParent(item)

		.tree.SetName(parent is 0 ? item : parent, .rootLabel)
		}

	buildRootLabel()
		{
		lib = .Send('CurrentTable')
		if not Libraries().Has?(lib)
			lib = '(' $ lib $ ')'
		name = .Send('CurrentName')
		return lib $ ':' $ name
		}

	buttons(visible)
		{
		.sort.SetVisible(visible)
		.public.SetVisible(visible)
		.private.SetVisible(visible)
		.data.SetVisible(visible)
		.meth.SetVisible(visible)
		}

	functionItems(sourceCode)
		{
		list = []
		scanner = ScannerWithContext(sourceCode)
		while scanner isnt token = scanner.Next()
			if scanner.Type() is #IDENTIFIER and
				token.GlobalName?() and
				scanner.Prev() isnt '.' and scanner.Prev() isnt '#' and
				scanner.Next() isnt ':'
				list.Add(token)
		return list
		}

	maxInheritances: 50
	loadHierarchy(lib, parent)
		{
		libs = Libraries().Reverse!()
		if not Libraries().Has?(lib)
			return true
		loop = 0
		while parent isnt false
			{
			loop++
			if loop > .maxInheritances
				{
				Print('ClassOutlineControl: too many levels of inheritance')
				return false // for test
				}
			libs = libs[libs.Find(lib) ..]
			if parent.Prefix?('_')
				libs.Remove(lib)
			found = false
			for l in libs
				if false isnt x = Query1(l, name: parent.LeftTrim('_'), group: -1)
					{
					lib = l
					.addTreeNode(lib, parent, x.text)
					parent = ClassHelp.SuperClass(x.text)
					found = true
					break
					}
			if found is false
				break
			}
		return true
		}

	addTreeNode(lib, parent, text)
		{
		root = .tree.AddItem(TVI.ROOT, lib $ ':' $ parent.LeftTrim('_'), 0, true)
		items = .filterItems(ClassHelp.ClassMembers(text))
		for method in items
			.tree.AddItem(root, method)
		}

	Children?(item)
		{
		return .tree.HasChildren?(item)
		}

	TreeView_ItemClicked(item)
		{
		itemName = .tree.GetName(item)
		curRecName = .Send('CurrentTable') $ ':' $ .Send('CurrentName')
		if not .tree.HasChildren?(item) // members
			{
			member = '.' $ itemName.Tr(':')
			root = .tree.GetName(.tree.GetParent(item))
			if not root.Has?('(') and root isnt curRecName
				.gotoDifferentRecord(root, member)
			else
				.Send('ClassOutline_SelectItem', member)
			}
		else if itemName.Tr('()') isnt curRecName.Tr('()') // class
			.gotoDifferentRecord(itemName)
		}

	gotoDifferentRecord(root, member = false)
		{
		rootLib = root.BeforeFirst(':')
		name = root.AfterFirst(':')
		if member isnt false
			name $= member
		libView = .Send('CurrentLibView')
		GotoLibView(name, libs: Object(rootLib), libview: libView is 0 ? false : libView)
		}

	Outline_Highlight(parent, child)
		{
		.tree.Highlight(parent, child)
		}

	NewValue(value/*unused*/)
		{
		.Refresh()
		}

	lastdata: ''
	Refresh()
		{
		.Set(.lastdata)
		}
	}
