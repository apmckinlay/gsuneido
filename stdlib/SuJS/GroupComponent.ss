// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
HtmlContainer
	{
	Xstretch: ""
	Ystretch: ""
	elements: ()
	New(elements)
		{
		super()
		.CreateElement('div')
		.SetStyles(Object(
			'display': 'inline-flex',
			'flex-direction': .Dir is 'vert' ? 'column' : 'row',
			'align-items': 'baseline'))

		.xmin0 = .Xmin
		.set_xstretch? = .Xstretch is ""
		.ymin0 = .Ymin
		.set_ystretch? = .Ystretch is ""

		.elements = Object()
		for e in elements.Values(list:)
			.elements.Add(.Construct(e))
		.Recalc()
		}

	CalcXminByControls(@args)
		{
		.xmin0 = .DoCalcXminByControls(@args)
		.Recalc()
		}

	Recalc()
		{
		.Xmin = .xmin0
		.Ymin = .ymin0
		if .Dir is 'vert'
			.vertRecalc()
		else
			.horzRecalc()
		.SetMinSize()
		}

	vertRecalc()
		{
		xmin = xstretch = ymin = ystretch = left = right = false
		maxheight = 0
		for el in .elements
			{
			xstretch = Max(xstretch, el.Xstretch)
			if el.Left is 0
				xmin = Max(xmin, el.Xmin)
			else
				{
				left = Max(left, el.Left)
				right = Max(right, el.Xmin - el.Left)
				}
			ymin += el.Ymin
			ystretch += el.Ystretch
			if el.GetDefault(#Shrinkable, false) is false
				el.El.SetStyle('flex-shrink', 0)
			if el.Ystretch isnt false
				el.El.SetStyle('flex-grow', el.Ystretch)
			if el.Xstretch >= 0
				{
				el.El.SetStyle('align-self', 'stretch')
				el.El.SetStyle('width', '')
				}

			maxheight += el.CalcMaxHeight()
			}
		xmin = Max(xmin, left + right)
		.MaxHeight = maxheight
		.setMin(xmin, ymin)
		.setStretch(xstretch, ystretch)
		.setLeft(left)
		}

	setLeft(left)
		{
		.Left = left
		for el in .elements
			{
			if el.Left isnt 0 and 0 < elementLeft = left - el.Left
				el.El.SetStyle('padding-left', elementLeft $ 'px')
			}
		}

	horzRecalc()
		{
		xmin = xstretch = ymin = ystretch = false
		for el in .elements
			{
			ystretch = Max(ystretch, el.Ystretch)
			ymin = Max(ymin, el.Ymin)
			xmin += el.Xmin
			xstretch += el.Xstretch
			if el.GetDefault(#Shrinkable, false) is false
				el.El.SetStyle('flex-shrink', 0)
			if el.Xstretch isnt false
				el.El.SetStyle('flex-grow', el.Xstretch)
			if el.Ystretch >= 0
				{
				el.El.SetStyle('align-self', 'stretch')
				el.El.SetStyle('height', '')
				}
			else if el.Base?(VertComponent)
				{
				// VertControl doesn't set .Top
				el.El.SetStyle('align-self', 'flex-start')
				}
			}
		.setMin(xmin, ymin)
		.setStretch(xstretch, ystretch)
		}

	setMin(xmin, ymin)
		{
		if xmin > .xmin0
			.Xmin = xmin
		if ymin >= .ymin0
			.Ymin = ymin
		}

	setStretch(xstretch, ystretch)
		{
		if .set_xstretch?
			.Xstretch = xstretch
		if .set_ystretch?
			.Ystretch = ystretch
		}

	Insert(i, control)
		{
		_at = [parent: this, at: i]
		.elements.Add(el = .Construct(control), at: i)
		DoStartup(el)
		.WindowRefresh()
		}

	Remove(i)
		{
		if not .elements.Member?(i)
			return
		.elements[i].Destroy()
		.elements.Delete(i)
		.WindowRefresh()
		}

	RemoveAll()
		{
		.elements.Each(#Destroy)
		.elements = Object()
		.WindowRefresh()
		}

	Tally()
		{ return .elements.Size() }

	GetChildren()
		{ return .elements }

	Get()
		{
		ob = Object()
		for (element in .elements)
			ob[element.Name] = element.Get()
		return ob
		}
	}
