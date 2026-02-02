// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: Canvas
	Title: Canvas
	Xstretch: 1
	Ystretch: 1
	ComponentName: 'Canvas'
	New()
		{
		.items = Object()
		.itemsMap = Object()
		.selected = Object()
		.paste_offset = 0
		}

	SyncSelected(selectedIds)
		{
		.selected = Object()
		for id in selectedIds
			if .itemsMap.Member?(id)
				.selected.Add(.itemsMap[id])
		}

	SyncItem(item, recursive? = false)
		{
		.Act('AfterEdit', item.Id, item.GetSuJSObject()[1..], :recursive?)
		}

	AddItem(item)
		{
		.items.Add(item)
		_canvas = this
		if false isnt ob = item.GetSuJSObject()
			{
			.itemsMap[ob.id] = item
			.Act('AddItem', ob)
			}
		.Send('CanvasChanged')
		}

	AddItemAndSelect(item)
		{
		.items.Add(item)
		_canvas = this
		if false isnt ob = item.GetSuJSObject()
			{
			.itemsMap[ob.id] = item
			.Act('AddItem', ob)
			.Select(.items.Find(item))
			}
		.Send('CanvasChanged')
		}

	MoveToBack(item)
		{
		.items.Remove(item)
		.items.Add(item, at: 0)
		.Act('MoveToBack', item.Id)
		.Select(0)
		.Send('CanvasChanged')
		}

	MoveToFront(item)
		{
		.items.Remove(item)
		.items.Add(item)
		.Act('MoveToFront', item.Id)
		.Select(.items.Size() - 1)
		.Send('CanvasChanged')
		}

	RemoveItem(item)
		{
		.items.Remove(item)
		.itemsMap.Delete(item.Id)
		.selected.Remove(item)
		.Send('CanvasChanged')
		.Act('RemoveItem', item.Id)
		}

	ResetSize(item)
		{
		// check if item is grouped
		if item.Grouped?
			{
			Alert('You cannot resize grouped items.', title: 'Error',
				flags: MB.ICONERROR)
			return
			}
		item.ResetSize()
		.Act('ResetSize', item.Id, item.GetCoordinates())
		.Select(.items.Find(item))
		.Send('CanvasChanged')
		}

	DeleteAll()
		{
		.items = Object()
		.itemsMap = Object()
		.selected = Object()
		.Act('DeleteAll')
		}

	GetSelected()
		{
		return .selected
		}

	GetAllItems()
		{
		return .items
		}

	SelectAll()
		{
		.ClearSelect()
		for i in .items.Members()
			.Select(i)
		}

	ClearSelect()
		{
		.Act('ClearSelect', noSync:)
		.selected = Object()
		}

	Select(i)
		{
		if .selected.FindIf({ Same?(it, .items[i]) }) isnt false
			return
		.selected.Add(.items[i])
		.Act('SelectId', .items[i].Id)
		}

	DeleteSelected()
		{
		for item in .selected
			{
			.items.Remove(item)
			.itemsMap.Delete(item.Id)
			.Act('RemoveItem', item.Id)
			item.Destroy()
			}
		.selected = Object()
		.Send('CanvasChanged')
		}

	Get(items = false)
		{
		if items is false
			items = .items
		return items.Map({ it.Get() }).Copy()
		}

	FormatColor(color)
		{
		".SetColor(0x" $ color.Hex() $ ")"
		}

	FormatLineColor(color)
		{
		".SetLineColor(0x" $ color.Hex() $ ")"
		}

	copied: false
	offset: 15
	CopyItems()
		{
		if .selected.Empty?()
			return
		.paste_offset = .offset
		.copied = .Get(.selected)
		}

	PasteItems()
		{
		.ClearSelect()
		if .copied is false
			return

		_canvas = this
		for itemDef in .copied
			{
			item = Construct(itemDef).SetupScale()
			item.Move(.paste_offset, .paste_offset)
			.AddItemAndSelect(item)
			}
		.paste_offset += .offset
		.Send('CanvasChanged')
		}

	CutItems()
		{
		.CopyItems()
		.DeleteSelected()
		}

	color: 0xffffff
	GetColor()
		{
		return .color
		}
	lin_color: 0
	GetLineColor()
		{
		return .lin_color
		}

	// On_Copy is redirected to here (focus)
	// other keyboard shortcuts are handled in DrawCanvasControl
	On_Copy()
		{ .CopyItems() }

	SetReadOnly(readOnly)
		{
		.ClearSelect()
		super.SetReadOnly(readOnly)
		}

	ToItem(id, args)
		{
		if not .itemsMap.Member?(id) // items in group
			return
		(.itemsMap[id][args[0]])(@+1args)
		.Send('CanvasChanged')
		}

	SetXminYmin(xmin, ymin)
		{
		.Act('SetXminYmin', xmin, ymin)
		}
	}
