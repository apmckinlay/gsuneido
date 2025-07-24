// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
/* READ ME:
Model:
	This tree control is designed to be used with a database model of some form.
	As a result, the Controller will need to provide a model (IE: .Controller.Model).
	The model provides the criteria and queries required for populating the tree as
	well as fulfilling other tree requests.
	Some examples of applicable models are:
	- LibTreeModel
	- BookEditModel
	- LibViewNewItemModel

ExplorerMultiImageHandler:
	If the Controller is defining / using ExplorerMultiImageHandler, it will also
	need to be made available to this control so that the images are synced.
	If the controller is not defining / using ExplorerMultiImageHandler, this control
	will create and manage its own.
*/
MultiTreeViewControl
	{
	New(.inorder = false, multi? = true, .readonly = false, style = 0)
		{
		super(multi?, readonly, style)
		.SetImageList(.imageHandler.ImageList)
		.addChildren(TVI.ROOT, 0)
		}

	// Purposely not setting .model = .Controller.Model as
	// .Model can be reconstructed multiple times in ExplorerMultiControl.
	// As a result, we need to ALWAYS point to the ExplorerMultiControl.Model directrly
	getter_model()
		{ return .Controller.Model }

	getter_imageHandler()
		{
		imageHandler = .Controller.GetDefault(#ImageHandler, false)
		if .instanceImageHandler = imageHandler is false
			imageHandler = ExplorerMultiImageHandler()
		return .imageHandler = imageHandler
		}

	getter_contextMenu()
		{ return .contextMenu = ExplorerMultiTreeContextMenu(.Controller, .readonly) }

	DragMove(dragging, target, last? = false)
		{
		.Send(#SaveCode_AfterChange)
		.Move(dragging, target, last?)
		}

	Move(item, target, last? = false)
		{
		if target is oldParent = .GetParent(item)
			return
		if .ChildOfParent?(target, item)
			return .AlertError('Invalid Move',
				'Cannot move parent folder into sub folder')
		if last?
			.SelectItem(oldParent)
		x = .model.Get(.GetParam(item))
		.DeleteItem(item)
		.model.Move(x, .GetParam(target))
		if last?
			.postMove(x.name, target, oldParent)
		}

	postMove(movedName, newParent, oldParent)
		{
		.ExpandItem(newParent)
		.SortChildren(newParent)
		.SelectItem(newParent)
		for child in .GetChildren(newParent)
			if movedName is .GetName(child)
				{
				.SelectItem(child)
				break
				}
		if not .model.Children?(.GetParam(oldParent))
			.ExpandItem(oldParent, collapse:)
		}

	Expanding(item)
		{
		.addChildren(item, .GetParam(item), .GetChildren(item).Map(.GetParam))
		return 0
		}

	addChildren(item, num, prexistingChildren = #())
		{
		children = .childrenToAdd(num, prexistingChildren)
		if .inorder
			for child in children
				.addChildItem(item, child)
		else
			{
			for child in children	// folders
				if child.group
					.addItem(item, child.name, child.num, container?:)
			for child in children   // item
				if not child.group
					.addChildItem(item, child)
			}
		return children.Size() > 0
		}

	childrenToAdd(num, prexistingChildren)
		{ return .model.Children(num).RemoveIf({ prexistingChildren.Has?(it.num) }) }

	addChildItem(item, child)
		{
		modified = child.GetDefault(#lib_modified, '')
		committed = child.GetDefault(#lib_committed, '')
		container? = .model.Container?(child.num)
		return .addItem(item, child.name, child.num, container?, :modified, :committed)
		}

	Expanded(item)
		{
		if .HasChildren?(item)
			.setFolderImage(item, true)
		return 0
		}

	setFolderImage(item, open)
		{
		folderRec = .model.Get(.GetParam(item)).Set_default('')
		.SetImage(item, .folderImage(folderRec, open))
		}

	Collapsed(item)
		{
		if .Container?(item)
			.setFolderImage(item, false)
		return 0
		}

	ChildOfParent?(item, parent)
		{
		while 0 isnt item = .GetParent(item)
			if item is parent
				return true
		return false
		}

	Path(item)
		{
		s = .GetName(item)
		while 0 isnt item = .GetParent(item)
			s = .GetName(item) $ "/" $ s
		return s
		}

	RootName(item)
		{
		while 0 isnt parent = .GetParent(item)
			item = parent
		return .GetName(item)
		}

	DeleteItem(item)
		{
		if item is TVI.ROOT
			{
			super.DeleteItem(item)
			return
			}
		parent = .GetParent(item)
		super.DeleteItem(item)
		if not .Children?(parent)
			.ExpandItem(parent, collapse:)
		}

	DragCopy(dragging, target, fromtree = true, last? = true)		// Recursive
		{
		if fromtree
			{
			.Expanding(target)
			.Send(#SaveCode_AfterChange)
			htarget = target
			dragging = .GetParam(dragging)
			target = .GetParam(target)
			}
		else
			htarget = false
		x = .model.Get(dragging, origText?:)
		oldparent = x.num
		x.parent = target
		.model.EnsureUnique(x)
		if .model.NewItem(x)
			{
			if x.group > -1
				for child in .model.Children(oldparent)
					.DragCopy(child.num, x.num, false)
			if fromtree
				.dragadd(htarget, x, last?)
			}
		else
			.Send(#Tree_UpdateTabData, dragging, [field: #num, value: x.num])
		.SortChildren(fromtree ? htarget : target)
		}

	dragadd(target, x, last?, oldNum = false)
		{
		item = .addChildItem(target, x)
		if oldNum is false
			oldNum = x.num
		.Send(#Tree_UpdateTabData, item, [field: #num, value: oldNum])
		if last?
			.SelectItem(item)
		}

	Rename(item, name)
		{
		x = .model.Get(.GetParam(item))
		if name is x.name
			return 0
		result = true
		try
			.model.Rename(x, name)
		catch (e)
			result = KeyException.Translate(e, 'Rename')
		if result isnt true
			{
			msg = "can't rename " $ x.name $ " to " $ name $
				Opt(' (', result is false ? '' : result, ')')
			Alert(msg, title: 'Rename', flags: MB.ICONERROR)
			return 0
			}
		.SetName(item, name)
		.SortChildren(.GetParent(item))
		.Send(#Tree_UpdateTabData, item, [field: #item, value: item])
		return 1
		}

	addItem(parent, name, num, container?, modified = '', committed = '')
		{
		// Add item to tree
		image = .image(container?, num, modified, committed)
		return .AddItem(parent, name, image, container?, num)
		}

	image(container?, num, lib_modified, lib_committed)
		{
		// Add item to tree
		return container?
			? .folderImage([:num, :lib_modified, :lib_committed, group: true])
			: .documentImage([:num, :lib_modified, :lib_committed, group: false])
		}

	folderImage(rec, open = false)
		{
		imageBase = open ? #FolderOpenTheme : #FolderClosedTheme
		imageModified = open ? #FolderOpenModifiedTheme : #FolderClosedModifiedTheme
		return .recordModified?(rec)
			? .imageHandler[imageModified]
			: .imageHandler[imageBase]
		}

	documentImage(rec)
		{
		return .recordModified?(rec)
			? .imageHandler.ModifiedTheme
			: .imageHandler.DocumentTheme
		}

	recordModified?(data)
		{ return .model.Method?(#Modified?) ? .model.Modified?(data) : false }

	AddNewItem(parent, container?, name, text)
		{
		x = Record(parent: .GetParam(parent), :name, group: container?, :text)
		if not .model.NewItem(.model.EnsureUnique(x))
			return false
		.ExpandItem(parent)
		if NULL is item = .findnum(parent, x.num)
			item = .addItem(parent, x.name, x.num, container?)
		.SortChildren(parent)
		.SetFocus()
		.EditLabel(item)
		return item
		}

	findnum(parent, num)
		{
		.ForEachChild(parent)
			{ |item|
			if .GetParam(item) is num
				{
				return item
				}
			}
		return NULL
		}

	SortChildren(parent)
		{ super.SortChildren(parent, .compareFunc) }

	compareFunc(lParam1, lParam2, lParamSort /*unused*/)
		{ return .model.TreeSort(lParam1, lParam2) }

	Insert(target, itemOb, table)
		{
		select = .insertItem(target, itemOb, table)
		if false is .HasChildren?(target)
			return
		// .insertItem failed to output the record, select the parent instead
		if select is false
			.SelectItem(target)
		.ExpandItem(target)
		.SortChildren(target)
		if select isnt false
			.SelectItem(select)
		}

	insertItem(hItem, itemOb, table)
		{
		rec = .buildInsertRecord(hItem, itemOb, table)
		restore? = SvcCheckRestorability(svcTable = SvcTable(table), rec)
		if not .model.NewItem(rec)
			return false
		if restore?
			svcTable.Remove(svcTable.MakeName(rec), deleted:)
		hItem = .addChildItem(hItem, rec)
		for child in itemOb.children
			.insertItem(hItem, child, table)
		return hItem
		}

	buildInsertRecord(hItem, itemOb, table)
		{
		rec = itemOb.Set_default('').Copy()
		if itemOb.table isnt table
			rec.lib_committed = ''
		rec.parent = .GetParam(hItem)
		rec.path = .Path(hItem).AfterFirst(table $ '/')
		rec.lib_modified = Date()
		.model.EnsureUnique(rec)
		return rec
		}

	Reset(.model)
		{
		.Refresh()
		.SetImageList(.imageHandler.ImageList)
		.addChildren(TVI.ROOT, 0)
		}

	ShowContextMenu(item, x, y)
		{ return .contextMenu.ShowTreeContextMenu(item, x, y) }

	NM_KILLFOCUS(lParam /*unused*/)
		{
		.UnselectAll(.Controller.CurItem)
		return 0
		}

	GotoPath(pathOb)
		{
		item = 0 // root
		for pathElement in pathOb
			for child in .GetChildren(item)
				if pathElement is .GetName(child)
					{
					.ExpandItem(item = child)
					break
					}
		.SelectItem(item)
		}

	Destroy()
		{
		if .instanceImageHandler
			.imageHandler.Destroy()
		super.Destroy()
		}
	}
