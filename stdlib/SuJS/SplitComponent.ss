// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
GroupComponent
	{
	styles: '
		.su-splitline-vert {
			position: absolute;
			border-top: 2px dashed black;
			left: 0;
			right: 0;
		}
		.su-splitline-horz {
			position: absolute;
			border-left: 2px dashed black;
			top: 0;
			bottom: 0;
		}'
	New(@elements)
		{
		super(elements)
		.LoadCssStyles()

		.SetStyles(#('position': 'relative'))
		ctrls = .GetChildren()
		.left = ctrls[0]
		.right = ctrls[2]
		if SuUI.GetCurrentWindow().
			GetComputedStyle(.left.El).
			GetPropertyValue('overflow') isnt 'auto'
			.left.SetStyles(#('overflow': 'hidden'))
		if SuUI.GetCurrentWindow().
			GetComputedStyle(.right.El).
			GetPropertyValue('overflow') isnt 'auto'
			.right.SetStyles(#('overflow': 'hidden'))
		}

	LoadCssStyles()
		{
		LoadCssStyles('split.css', .styles)
		}

	getter_splitterSize()
		{
		ctrls = .GetChildren()
		return .Dir is #vert ? ctrls[1].Ymin : ctrls[1].Xmin
		}

	// enum - None, Normal, MaxSecond
	splitMethod: #None
	SetDefaultSplit()
		{
		// Set a dummy width/height so that the element doesn't determine its width/height
		// based on its content dynamically. The elements will still stretch.
		style = .Dir is #vert ? 'height' : 'width'
		.left.El.SetStyle(style, '1px')
		.right.El.SetStyle(style, '1px')
		}

	n: #(.5, .5)
	GetSplit()
		{
		return .n
		}

	SetSplit(.n)
		{
		split = .calcSplit(n)
		.doResize()
			{
			if .Dir is 'vert'
				{
				.left.El.SetStyle('height', .getCalc(n[0], split[0], .right.Ymin))
				.right.El.SetStyle('height', .getCalc(n[1], split[1], .left.Ymin))
				}
			else
				{
				.left.El.SetStyle('width', .getCalc(n[0], split[0], .right.Xmin))
				.right.El.SetStyle('width', .getCalc(n[1], split[1], .left.Xmin))
				}
			}
		.splitMethod = #Normal
		}

	splitLine: false
	splitPos: false
	Splitter_mousedown(pos)
		{
		.splitLine = CreateElement('div', .El,
			className: .Dir is 'vert' ? 'su-splitline-vert' : 'su-splitline-horz')
		containerRect = .El.GetBoundingClientRect()
		if .Dir is 'vert'
			.splitLine.SetStyle('top', (.splitPos = pos - containerRect.y) $ 'px' )
		else
			.splitLine.SetStyle('left', (.splitPos = pos - containerRect.x) $ 'px' )
		}

	Splitter_mousemove(pos)
		{
		containerRect = .El.GetBoundingClientRect()
		halfSplit = .splitterSize / 2
		if .Dir is 'vert'
			{
			if pos - containerRect.y - halfSplit >= .left.Ymin and
				containerRect.y + containerRect.height - pos - halfSplit >= .right.Ymin
				{
				.splitLine.SetStyle('top', (.splitPos = pos - containerRect.y) $ 'px' )
				}
			}
		else
			{
			if pos - containerRect.x - halfSplit >= .left.Xmin and
				containerRect.x + containerRect.width - pos - halfSplit >= .right.Xmin
				{
				.splitLine.SetStyle('left', (.splitPos = pos - containerRect.x) $ 'px' )
				}
			}
		}

	Splitter_mouseup()
		{
		containerRect = .El.GetBoundingClientRect()
		if .Dir is 'vert'
			{
			topPercent = .splitPos / containerRect.height
			bottomPercent = 1 - topPercent
			split = .calcSplit(Object(topPercent, bottomPercent))
			.doResize()
				{
				.left.El.SetStyle('height', .getCalc(topPercent, split[0], .right.Ymin))
				.right.El.SetStyle('height',
					.getCalc(bottomPercent, split[1], .left.Ymin))
				}
			.Event('UpdateSplit', .n = [topPercent, bottomPercent])
			}
		else
			{
			leftPercent = .splitPos / containerRect.width
			rightPercent = 1 - leftPercent
			split = .calcSplit(Object(leftPercent, rightPercent))
			.doResize()
				{
				.left.El.SetStyle('width', .getCalc(leftPercent, split[0], .right.Xmin))
				.right.El.SetStyle('width', .getCalc(rightPercent, split[1], .left.Xmin))
				}
			.Event('UpdateSplit', .n = [leftPercent, rightPercent])
			}

		.splitLine.Remove()
		.splitLine = .splitPos = false
		.splitMethod = #Normal
		}

	getCalc(percent, split, min)
		{
		return 'min(' $
			percent.DecimalToPercent(2) $ '%' $ ' - ' $ split $ 'px, ' $
			'100% - ' $ .splitterSize $ 'px - ' $ min $ 'px)'
		}

	calcSplit(n)
		{
		size = .splitterSize
		return n.Map({ (it * size).Round(2) })
		}

	MaximizeSecond()
		{
		Assert(.Dir is: 'vert')
		// This code is specific for VirtualListViewControl
		// It assume the left has only one child
		// Cannot get Ymin from left directly because its height is set to 0 explicitly
		topSize = .left.GetChildren()[0].Ymin
		.doResize()
			{
			.left.El.SetStyle('height', topSize $ 'px')
			.right.El.SetStyle('height', .getCalc(1, topSize + .splitterSize, topSize))
			}
		.splitMethod = #MaxSecond
		}

	Recalc()
		{
		super.Recalc()
		if .refreshTimer is false
			.refreshTimer = SuDelayed(0, .refresh)
		}

	refreshTimer: false
	refresh()
		{
		.refreshTimer = false
		switch (.splitMethod)
			{
		case #Normal:
			.SetSplit(.n)
		case #MaxSecond:
			.MaximizeSecond()
		default:
			}
		}

	doResize(block)
		{
		.TopDown(#BeforeResize)
		block()
		.TopDown(#AfterResize)
		}

	Destroy()
		{
		if .refreshTimer isnt false
			.refreshTimer.Kill()
		super.Destroy()
		}
	}
