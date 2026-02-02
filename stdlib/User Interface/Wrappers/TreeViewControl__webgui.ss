// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: 'TreeView'
	ComponentName: 'TreeView'
	New(readonly = false, style = 0)
		{
		.ComponentArgs = Object(readonly, style)
		.root = Object(id: 0, name: '', image: 0, children: Object(), expanded?:)
		.trees = Object(.root)
		.next = 1
		}

	AddItem(parent, name, image = 0, container? = false, param = 0)
		{
		parent = .convert(parent)
		if .trees.Member?(parent) is false
			throw 'TreeViewControl.AddItem: Cannot not find parent ' $ Display(parent)

		id = .next++
		image = image is 0
			? false
			: .images.GetDefault(image - 1, false)
		newItem = Object(:parent, :id, :name, :image, :param, expanded?: false)
		if container? is true
			newItem.children = Object()
		.trees[parent].children.Add(newItem)
		.trees[id] = newItem
		.Act('AddItem', parent, id, name, image, container?)
		return id
		}

	convert(item)
		{
		return item is TVI.ROOT ? 0 : item
		}

	TREEVIEW_TOGGLE(collapsed?, hItem)
		{
		.trees[hItem].expanded? = not collapsed?
		if collapsed?
			{
			.AttemptSend("Collapsing", hItem)
			.AttemptSend("Collapsed", hItem)
			}
		else
			{
			.AttemptSend("Expanding", hItem)
			.AttemptSend("Expanded", hItem)
			}
		}

	TVN_SELCHANGED(oldSelect, newSelect)
		{
		Assert(oldSelect is: .selected)
		.selected = newSelect
		.Tree_SelChanged(oldSelect, newSelect)
		}

	Tree_SelChanged(olditem, newitem)
		{
		.Send("SelectTreeItem", olditem, newitem)
		}

	GetSelectedItem()
		{
		return .selected
		}

	GetChildren(item = false)
		{
		item = .convert(item)
		if not .trees.Member?(item) or not .trees[item].Member?(#children)
			return #()
		return .trees[item].children.Map({ it.id })
		}

	ExpandItem(item, collapse = false)
		{
		item = .convert(item)
		.Act("ExpandItem", item, collapse)
		.TREEVIEW_TOGGLE(collapse, item)
		}

	Expanded?(item)
		{
		item = .convert(item)
		return .trees[item].expanded?
		}

	DeleteItem(item)
		{
		item = .convert(item)
		if not .ItemExists?(item)
			return false

		parent = .trees[item].parent
		.trees[parent].children.RemoveIf({ it.id is item })
		toDelete = Object()
		.deleteItem(item, toDelete)
		.Act('DeleteItem', item, toDelete)
		if toDelete.Has?(.selected)
			{
			next = toDelete.Max() + 1
			if not .ItemExists?(next)
				next = toDelete.Min() - 1
			if .ItemExists?(next)
				.SelectItem(next)
			}
		return true
		}

	deleteItem(item, toDelete)
		{
		ob = .trees[item]
		.trees.Erase(item)
		toDelete.Add(item)
		for child in ob.GetDefault(#children, #())
			.deleteItem(child.id, toDelete)
		}

	Children?(item)
		{
		return Boolean?(x = .Send("Children?", item)) ? x : false
		}

	GetName(item)
		{
		item = .convert(item)
		if not .trees.Member?(item)
			return ''
		return .trees[item].name
		}

	GetParent(item)
		{
		if not .trees.Member?(item)
			return 0
		return .trees[item].parent
		}

	GetParam(item)
		{
		if not .trees.Member?(item)
			return false
		return .trees[item].param
		}

	EditLabel(item)
		{
		item = .convert(item)
		if not .trees.Member?(item)
			return ''
		.SelectItem(item)
		.Act('EditLabel', item)
		}

	EndLabelEdit(id, text)
		{
		.Send("Rename", id, text)
		}

	SetName(item, name)
		{
		.trees[item].name = name
		.Act('SetName', item, name)
		}

	ItemExists?(item)
		{
		return .trees.Member?(item)
		}

	HasChildren?(item)
		{
		if not .ItemExists?(item)
			return false
		return .trees[item].GetDefault(#children, #()).NotEmpty?()
		}

	Container?(item)
		{
		if not .ItemExists?(item)
			return false
		return .trees[item].Member?(#children)
		}

	selected: 0
	SelectItem(item)
		{
		item = .convert(item)
		if not .trees.Member?(item)
			return
		oldSelect = .selected
		.Act('SelectItem', .selected = item)
		if .selected isnt 0
			.Send("SelectTreeItem", oldSelect, .selected)
		}

	images: #()
	SetImageList(.images) {}

	SetImage(item, image)
		{
		item = .convert(item)
		if not .ItemExists?(item)
			return
		image = image is 0
			? false
			: .images.GetDefault(image - 1, false)
		.Act('SetImage', item, image)
		}

	// lpfnCompare(lParam1, lParam2, lParamSort) returns:
	// 	-# if lParam1 comes first
	//	+# if lParam2 comes first
	// 	NOTE: lParamSort, corresponds to the optional lParam member in TVSORTCB.
	//		  Currently unused
	SortChildren(hParent, lpfnCompare)
		{
		hParent = .convert(hParent)
		if not .ItemExists?(hParent)
			return

		if false is children = .trees[hParent].GetDefault(#children, false)
			return

		for i in ..children.Size()
			children[i].tree_sort = i
		children.Sort!({ |a, b| lpfnCompare(a.param, b.param, false) < 0 })

		newOrders = Object().AddMany!(0, children.Size())
		for m, v in children
			newOrders[v.tree_sort] = m
		.Act(#ReorderChildren, hParent, newOrders)
		}

	ForEachChild(parent, block)
		{
		item = .convert(parent)
		if not .ItemExists?(item)
			return
		.trees[item].GetDefault(#children, #()).Each({ block(it.id) })
		}

	AttemptSend(@args)
		{
		if .Send(@args) is 0 and .Method?(args[0])
			(this[args[0]])(@+1args)
		return 0
		}

	Default(@args)
		{
		SuServerPrint(@args)
		}
	}
