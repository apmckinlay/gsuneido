// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 'Canvas'
	Xstretch: 1
	Ystretch: 1
	Xmin: 50
	Ymin: 50
	styles: `
		.su-canvas {
			position: relative;
			background-color: white;
		}
		.su-canvas:focus {
			outline: none;
		}
		.su-canvas-container {
			position: absolute;
			width: 100%;
			height: 100%;
			overflow: hidden;
		}`
	New()
		{
		LoadCssStyles('su-canvsa.css', .styles)
		.CreateElement('div', className: 'su-canvas')
		.El.tabIndex = "0"
		.container = CreateElement('div', .El, className: 'su-canvas-container')
		.canvas = CreateElement('svg', .container, namespace: SvgDriver.Namespace)

		.driver = SvgDriver(.canvas, Object(width: 999, height: 999), factor: 1)
		svgSize = .driver.ConvertToPixcel(999/*=a big dimension*/)
		.canvas.style.width = svgSize $ 'px'
		.canvas.style.height = svgSize $ 'px'

		.SetMinSize()
		if .Xstretch is false
			.El.SetStyle('width', .Xmin $ 'px')
		.items = Object()
		.selected = Object()
		}

	AddItem(ob)
		{
		_canvas = this
		_spec = ob
		.items.Add(Construct(ob))
		}

	MoveToBack(id)
		{
		if false isnt item = .findItemById(id)
			{
			.items.Remove(item)
			.items.Add(item, at: 0)
			item.MoveToBack()
			}
		}

	MoveToFront(id)
		{
		if false isnt item = .findItemById(id)
			{
			.items.Remove(item)
			.items.Add(item)
			item.MoveToFront()
			}
		}

	RemoveItem(id)
		{
		if false isnt item = .findItemById(id)
			{
			item.Remove()
			.items.RemoveIf({ it.Id is id })
			.selected.RemoveIf({  it.Id is id })
			}
		}

	ResetSize(id, coordinates)
		{
		if false isnt item = .findItemById(id)
			{
			item.ResetSize(@coordinates)
			}
		}

	AfterEdit(id, args, recursive? = false)
		{
		.syncSelected()
		if false isnt item = .findItemById(id, recursive?)
			{
			item.AfterEdit(@args)
			}
		.syncSelected()
		}

	DeleteAll()
		{
		.canvas.innerHTML = ''
		.items = Object()
		.selected = Object()
		}

	GetSelected()
		{
		return .selected
		}

	GetAllItems()
		{
		return .items
		}

	SelectPoint(x, y)
		{
		.ClearSelect()
		if false isnt i = .ItemAtPoint(x, y)
			.Select(i)
		}

	ItemAtPoint(x, y)
		{
		for item in .selected
			if item.Contains(x, y)
				return .items.Find(item)
		for (i = .items.Size() - 1; i >= 0; --i)
			if .items[i].Contains(x, y)
				return i
		return false
		}

	SelectRect(x1, y1, x2, y2)
		{
		.MaybeClearSelect()
		for item in .items
			if item.Overlaps?(x1, y1, x2, y2)
				{
				item.Select()
				.selected.Add(item)
				}
		.syncSelected()
		}

	MaybeClearSelect(_event = false)
		{

		if event isnt false and event.ctrlKey is false and event.shiftKey is false
			{
			.ClearSelect()
			return true
			}
		else
			return false
		}

	SelectAll()
		{
		.ClearSelect()
		for item in .items
			{
			item.Select()
			.selected.Add(item)
			}
		.syncSelected()
		}

	ClearSelect(noSync = false)
		{
		for item in .selected
			item.Unselect()
		.selected = Object()
		if noSync is false
			.syncSelected()
		}

	Select(i)
		{
		.items[i].Select()
		.selected.AddUnique(.items[i])
		.syncSelected()
		}

	// from server, no need to sync
	SelectId(id)
		{
		if false isnt item = .findItemById(id)
			{
			item.Select()
			.selected.AddUnique(item)
			}
		}

	UnSelect(i)
		{
		.items[i].Unselect()
		.selected.Remove(.items[i])
		.syncSelected()
		}

	MoveSelected(dx, dy)
		{
		if .selected.Empty?()
			return
		rects = Object()
		for item in .selected
			{
			r = item.BoundingRect()
			rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2, :item)
			rects.Add(rect)
			}
		move = DrawSelectTracker_CalcNextMove(dx, dy, rects, this)
		for item in .selected
			item.Move(move.x, move.y)
		}

	DeleteSelected()
		{
		.Event('DeleteSelected')
		}

	GetWidth()
		{
		return .El.clientWidth
		}

	GetHeight()
		{
		return .El.clientHeight
		}

	Getter_Driver()
		{
		return .driver
		}

	readOnly: false
	SetReadOnly(.readOnly)
		{
		.ClearSelect()
		}

	GetReadOnly()
		{
		return .readOnly
		}

	syncSelected()
		{
		.Event('SyncSelected', .selected.Map({ it.Id }))
		}

	findItemById(id, recursive? = false)
		{
		if recursive? is true
			return .find(.items, id)
		return .items.FindOne({ it.Id is id })
		}

	find(items, id)
		{
		for item in items
			{
			if item.Id is id
				return item
			if item.Base?(SuCanvasGroup) and false isnt found = .find(item.GetItems(), id)
				return found
			}
		return false
		}

	SetXminYmin(.Xmin, .Ymin)
		{
		.SetMinSize()
		.WindowRefresh()
		}
	}
