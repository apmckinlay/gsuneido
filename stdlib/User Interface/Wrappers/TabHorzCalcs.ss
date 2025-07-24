// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
TabCalcs
	{
	ScrollNextImage: 'right'
	ScrollPrevImage: 'left'
	FontOrientation: 0
	New(@args)
		{
		super(@args)
		if args.orientation is #top
			.topTabs()
		else
			.bottomTabs()
		}

	TabDragSpecs(x, y, availableSpace)
		{
		drag? = x >= 0 and x <= availableSpace and y >= 0 and y <= .H
		return [:drag?, check: x, cursor: IDC.HSPLITBAR]
		}

	topTabs()
		{
		.RenderRect = .renderRectTop
		.DrawSpecs = .drawSpecsTop
		.TextPos = .textPosTop
		.ImagePos = .imagePosTop
		.LinePoints = .linePointsTop
		.SelectSize = .calcSelectTop

		.extraControlY = 0
		.offset = Max(.Ymin - .TabHeight, 0)
		}

	bottomTabs()
		{
		.RenderRect = .renderRectBottom
		.DrawSpecs = .drawSpecsBottom
		.TextPos = .textPosBottom
		.ImagePos = .imagePosBottom
		.LinePoints = .linePointsBottom
		.SelectSize = .calcSelectBottom

		.extraControlY = 1
		.offset = -4
		}

	renderRectTop(i, tab, prevEnd)
		{
		renderRect = .baseRenderRect(prevEnd, tab.width)
		// Must calculate each time to ensure we position the tabs correctly
		extraControlOffset = Max(.Ymin - .TabHeight, 0)
		.calcSelectTop(i, renderRect, extraControlOffset)
		renderRect.bottom = .TabHeight + extraControlOffset
		return renderRect
		}

	calcSelectTop(i, renderRect, extraControlOffset = false)
		{
		if false is extraControlOffset
			extraControlOffset = .extraControlOffset
		renderRect.top = .SelectOffset(i) + extraControlOffset
		}

	getter_extraControlOffset()
		{
		// Must calculate each time to ensure we position the tabs correctly
		return Max(.Ymin - .TabHeight, 0)
		}

	baseRenderRect(prevEnd, width)
		{
		rect = Record()
		rect.AttachRule(#left, function(){ this.start })
		rect.AttachRule(#right, function(){ this.end })
		rect.AttachRule(#tipX, function(){ this.start })
		rect.AttachRule(#tipY, function(){ this.bottom })
		rect.start = prevEnd
		rect.end = prevEnd + width
		return rect
		}

	renderRectBottom(i, tab, prevEnd)
		{
		renderRect = .baseRenderRect(prevEnd, tab.width)
		renderRect.top = 0
		.calcSelectBottom(i, renderRect)
		return renderRect
		}

	calcSelectBottom(i, renderRect)
		{
		renderRect.bottom = .TabHeight - .SelectOffset(i)
		}

	drawSpecsTop(wLarge, hLarge)
		{
		specs = .baseDrawSpecs(wLarge, hLarge)
		erase = specs.eraseSize
		specs.overrideRect = [left: 0, top: erase, right: wLarge, bottom: hLarge]
		specs.overrideFill = [left: 1, right: wLarge - 1, top: erase, bottom: erase + 1]
		return specs
		}

	rectCurve: 0.5
	baseDrawSpecs(wLarge, hLarge)
		{
		rect = Object(left: 0, top: 0, right: wLarge, bottom: hLarge)
		return Object(
			baseRound: rect,
			baseFill: [right: rect.right, bottom: rect.bottom]
			ellipseSize: ellipse = (hLarge * .rectCurve).Ceiling(),
			eraseSize: ellipse >> 1,
			)
		}

	drawSpecsBottom(wLarge, hLarge)
		{
		specs = .baseDrawSpecs(wLarge, hLarge)
		specs.overrideRect = [left: 0, top: 0, right: wLarge, bottom: specs.eraseSize]
		specs.overrideFill = [
			left: 1,
			right: wLarge - 1,
			top: specs.eraseSize - 1,
			bottom: specs.eraseSize]
		return specs
		}

	linePointsTop()
		{
		return [x1: 0, y1: .Ymin - 1, x2: .W -1, y2: .Ymin - 1]
		}

	linePointsBottom()
		{
		return [x1: 0, y1: 0, x2: .W -1, y2: 0]
		}

	textPosTop(tab, selectedTab)
		{
		specs = .baseTextPos(tab, selectedTab)
		specs.y += .PaddingTop - 1
		return specs
		}

	baseTextPos(tab, selectedTab)
		{
		padding = .ImageWidth(tab.image) + .PaddingSide
		pos = Object(
			x: tab.renderRect.left + padding
			y: tab.renderRect.top)
		if not selectedTab and tab.renderWidth is tab.width
			pos.x += tab.textBoldOffset
		return pos
		}

	textPosBottom(tab, selectedTab)
		{
		textRect = .baseTextPos(tab, selectedTab)
		textRect.y += tab.renderRect.bottom - tab.height + .PaddingTop + 2
		return textRect
		}

	Resize?(w, h /*unused*/)
		{
		return w isnt .W
		}

	Getter_TabBarSize()
		{
		return .W
		}

	ResizeExtraControl(extraControl, ctrlPos, ctrlSize)
		{
		extraControl.Resize(ctrlPos, .extraControlY, ctrlSize, .H - 1)
		}

	ResizeButton(button, pos)
		{
		button.Resize(pos, .PaddingTop + .offset, .ButtonSize, .ButtonSize)
		return .ButtonSize
		}

	imagePosTop(tab)
		{
		x = tab.renderRect.left + .PaddingSide
		y = tab.renderRect.top + .PaddingTop - 2
		return Object(:x, :y)
		}

	imagePosBottom(tab)
		{
		x = tab.renderRect.left + .PaddingSide
		y = tab.renderRect.bottom - tab.height + .PaddingTop + 2
		return Object(:x, :y)
		}

	ImageRect(tab)
		{
		imageSpecs = .CalcImageSpecs(tab)
		dimensions = .ImageDimensions(tab)
		right = dimensions.width + left = imageSpecs.x
		bottom = dimensions.height + top = imageSpecs.y
		return Object(:left, :right, :top, :bottom)
		}

	InvalidateRect(i, tab)
		{
		renderRect = (.RenderRect)(i, tab, tab.renderRect.left)
		return Object(
			left: renderRect.left,
			right: renderRect.right,
			top: 0,
			bottom: .Ymin)
		}
	}