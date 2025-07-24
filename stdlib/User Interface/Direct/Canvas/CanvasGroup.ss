// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	Grouped?: true
	New(.items, .name = '')
		{
		if items.NotEmpty?() and not Instance?(items[0])
			.items = items.Map(Construct).Copy()
		}
	Paint(data = false)
		{
		for (item in .items)
			item.Paint(:data)
		}
	BoundingRect()
		{
		if .items.Size() is 0
			return super.BoundingRect()
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
	GetItems()
		{ return .items }
	SetItems(items)
		{
		.items = items
		}
	ToString()
		{
		return 'CanvasGroup(Object(' $ .items.Map(
			{ |x|
			x.ToString() $ CanvasControl.FormatColor(x.GetColor()) $
				CanvasControl.FormatLineColor(x.GetLineColor())
			}).Join(', ') $ '))'
		}
	StringToSave()
		{
		return .ToString()
		}
	ObToSave()
		{
		return Object('CanvasGroup', .items.Map(#ObToSave))
		}
	Move(dx, dy)
		{
		for item in .items
			item.Move(dx, dy)
		}
	GetResource()
		{
		resOb = Object()
		.items.Each() { resOb.MergeUnion(it.GetResource()) }
		return resOb
		}
	GetName()
		{ return .name }
	GetSize()
		{
		b = .BoundingRect()
		return Object(x1: b.x1, y1: b.y1, w: b.x2 - b.x1, h: b.y2 - b.y1)
		}
	SetSize(x1, y1, x2, y2)
		{
		.x1 = x1
		.y1 = y1
		.x2 = x2
		.y2 = y2
		}
	Resize(origx /*unused*/, origy /*unused*/, x /*unused*/, y /*unused*/)
		{
		}
	Scale(by, print? = false)
		{
		for item in .items
			item.Scale(by, :print?)
		}
	ReverseScale(by, print? = false)
		{
		for item in .items
			item.ReverseScale(by, :print?)
		}

	SetThick(thick)
		{
		for item in .items
			item.SetThick(thick)
		}

	GetSuJSObject()
		{
		items = Object()
		for item in .items
			items.Add(item.GetSuJSObject())
		return Object('SuCanvasGroup', items, id: .Id)
		}
	}