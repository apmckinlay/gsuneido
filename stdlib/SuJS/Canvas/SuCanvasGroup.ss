// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	New(items)
		{
		.items = Object()
		.build(items)
		}

	build(items)
		{
		_canvas = .Canvas
		for item in items
			{
			_spec = item
			.items.Add(Construct(item))
			}
		}

	BoundingRect()
		{
		r = .items[0].BoundingRect()
		for (item in .items)
			{
			br = item.BoundingRect()
			r.x1 = Min(r.x1, br.x1)
			r.x2 = Max(r.x2, br.x2)
			r.y1 = Min(r.y1, br.y1)
			r.y2 = Max(r.y2, br.y2)
			}
		return r
		}
	Resize(origx /*unused*/, origy /*unused*/, x /*unused*/, y /*unused*/)
		{
		}

	Move(dx, dy)
		{
		for item in .items
			item.Move(dx, dy)
		super.Move(dx, dy)
		}

	Remove()
		{
		for item in .items
			item.Remove()
		super.Remove()
		}

	MoveToBack()
		{
		for item in .items
			item.MoveToBack()
		}

	MoveToFront()
		{
		for item in .items
			item.MoveToFront()
		}

	AfterEdit(items)
		{
		.Remove()
		.items = Object()
		.build(items)
		super.AfterEdit()
		}

	GetItems()
		{
		return .items
		}
	}
