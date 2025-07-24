// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Create Plugin Wizard"
	CallClass()
		{
		ToolDialog(0, this)
		}
	images: false
	New()
		{
		sz = 16
		.images = CreateImageList(sz, sz)
		ImageList_AddVectorImage(.images, 'folder.emf', CLR.black, sz, sz)
		ImageList_AddVectorImage(.images, 'open_folder.emf', CLR.black, sz, sz)
		.tree = .Data.Vert.TreeView
		.tree.SetImageList(.images)
		.model = Construct(LibViewNewItemModel)
		.addchildren(TVI.ROOT, 0)
		.Data.Get().Observer(.Record_changed)
		}
	Controls()
		{
		return Object('Record'
			Object('Vert'
				#(Pair (Static Name) (Field name: name))
				#(Pair (Static Path) (Field readonly:, name: path))
				'Skip'
				Object('TreeView', readonly: true)
				#(Skip 5)
				#(Horz Fill (Button, Create, xmin: 50) Skip (Button Cancel xmin: 50))))
		}
	DefaultButton: Create
	On_Create()
		{
		data = .Data.Get()
		if data.name is ""
			{
			Alert('Please enter a Name')
			return
			}
		lib = data.path.Split('/')[0]
		if false isnt Query1(lib, name: data.item, group: -1)
			{
			Alert(data.name $ ' already exist in ' $ lib)
			return
			}
		data.num = .getnum(.selected) % 100000
		data.text = '#(\nExtensionPoints:\n\t(\n\t(/*<extention point name>*/)\n' $
			'\t)\nContributions:\n\t(\n\t(/*<plugin name>, <extention point>*/)\n\t)\n)'
		.output_new_item(data, lib)
		.Window.Result(true)
		}
	output_new_item(new_item, lib)
		{
		KeyException.TryCatch()
			{
			treeModel = new TreeModel(lib)
			rec = Record(
				parent: new_item.num
				name: new_item.name
				text: new_item.text
				group: false)
			treeModel.NewItem(rec)
			}
		SvcTable(lib).Publish('TreeChange', force:)
		}
	Record_changed(member)
		{
		data = .Data.Get()
		if member isnt 'name' or data.name.Prefix?('Plugin_')
			return
		data.name = 'Plugin_' $ data.name
		}
	selected: false
	SelectTreeItem(olditem /*unused*/, newitem)
		{
		.selected = newitem
		path = .getpath(newitem)
		.Data.SetField('path', path)
		}
	folder1: 0
	folder2: 1
	Expanding(item)
		{
		if .expanded?(item) or not .addchildren(item, .getnum(item))
			return 0
		.tree.SetImage(item, .folder2)
		return 0
		}
	expanded?(item)
		{ return .tree.HasChildren?(item) }
	addchildren(item, num)
		{
		children = .model.Children(num)
		for child in children	// folders
			if (child.group)
				.tree.AddItem(item, child.name, .folder1, true, child.num)
		return children.Size() > 0
		}
	Collapsed(item)
		{
		.tree.SetImage(item, .folder1)
		for item in .tree.GetChildren(item) 	// remove the children
			.delitem(item)
		}
	delitem(item)
		{
		if (item is TVI.ROOT)
			{
			.tree.DeleteItem(item)
			return
			}
		parent = .tree.GetParent(item)
		.tree.DeleteItem(item)
		if (not .expanded?(parent))
			.tree.SetImage(item, .folder1)
		}
	Children?(item)
		{ return .model.Children?(.getnum(item)) }
	getname(item)
		{ return .tree.GetName(item) }
	getnum(item)
		{ return .tree.GetParam(item) }
	getpath(item)
		{
		s = .getname(item)
		while (0 isnt (item = .tree.GetParent(item)))
			s = .getname(item) $ "/" $ s
		return s
		}
	get_curitem()
		{
		x = .tree.Selection
		return (x.Size() > 0) ? x[0] : false
		}
	Destroy()
		{
		if .images isnt false
			ImageList_Destroy(.images)
		super.Destroy()
		}
	}
