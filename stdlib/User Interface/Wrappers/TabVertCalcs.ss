// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
TabCalcs
	{
	ScrollNextImage: 'down'
	ScrollPrevImage: 'up'
	New(@args)
		{
		super(@.processArgs(args))
		.offset = Max(.Ymin - .TabHeight, 0)
		}

	processArgs(args)
		{
		if args.orientation is #right
			.rightTabs()
		else
			.leftTabs()
		return args
		}

	TabDragSpecs(x, y, availableSpace)
		{
		drag? = x >= 0 and x <= .W and y >= 0 and y <= availableSpace
		return [:drag?, check: y, cursor: IDC.VSPLITBAR]
		}

	leftTabs()
		{
		.RenderRect = .renderRectLeft
		.DrawSpecs = .drawSpecsLeft
		.TextPos = .textPosLeft
		.ImagePos = .imagePosLeft
		.LinePoints = .linePointsLeft
		.SelectSize = .calcSelectLeft
		.FontOrientation = 900 /*= vertical, text bottom facing to the right */

		.imageRectMethod = .imageRectLeft
		.extraControlX = -1
		.buttonOffset = 0
		}

	rightTabs()
		{
		.RenderRect = .renderRectRight
		.DrawSpecs = .drawSpecsRight
		.TextPos = .textPosRight
		.ImagePos = .imagePosRight
		.LinePoints = .linePointsRight
		.SelectSize = .calcSelectRight
		.FontOrientation = 2700 /*= vertical, text bottom facing to the left */

		.imageRectMethod = .imageRectRight
		.extraControlX = 1
		.buttonOffset = -2
		}

	renderRectLeft(i, tab, prevEnd)
		{
		renderRect = .baseRenderRect(prevEnd)
		renderRect.end = prevEnd + tab.width
		renderRect.right = .TabHeight
		.calcSelectLeft(i, renderRect)
		return renderRect
		}

	calcSelectLeft(i, renderRect)
		{
		renderRect.left = .SelectOffset(i)
		}

	baseRenderRect(prevEnd)
		{
		rect = Record()
		rect.AttachRule(#top, function(){ this.start })
		rect.AttachRule(#bottom, function(){ this.end })
		rect.AttachRule(#tipX, function(){ this.right })
		rect.AttachRule(#tipY, function(){ this.start })
		rect.start = prevEnd
		return rect
		}

	renderRectRight(i, tab, prevEnd)
		{
		renderRect = .baseRenderRect(prevEnd)
		renderRect.end = prevEnd + tab.width + 1
		.calcSelectRight(i, renderRect)
		renderRect.left = .offset
		return renderRect
		}

	calcSelectRight(i, renderRect)
		{
		renderRect.right = .TabHeight - .SelectOffset(i)
		}

	drawSpecsLeft(wLarge, hLarge)
		{
		specs = .baseDrawSpecs(wLarge, hLarge)
		erase = .eraseSize
		specs.overrideRect = [left: erase, top: 0, right: wLarge, bottom: hLarge]
		specs.overrideFill = [left: erase + 1, top: 1, right: erase, bottom: hLarge - 1]
		return specs
		}

	eraseSize: 30
	ellipseSize: 60
	baseDrawSpecs(wLarge, hLarge)
		{
		rect = Object(left: 0, top: 0, right: wLarge, bottom: hLarge)
		return Object(
			baseRound: rect,
			baseFill: [right: wLarge, bottom: hLarge],
			ellipseSize: .ellipseSize)
		}

	drawSpecsRight(wLarge, hLarge)
		{
		specs = .baseDrawSpecs(wLarge, hLarge)
		erase = wLarge - .eraseSize + 5 /*= padding*/
		specs.overrideRect = [left: 0, top: 0, right: erase, bottom: hLarge]
		specs.overrideFill = [left: -1, top: 1,	right: erase, bottom: hLarge - 1]
		return specs
		}

	linePointsLeft()
		{
		x = .Ymin - 1
		return [x1: x, y1: .PaddingSide, x2: x, y2: .H]
		}

	linePointsRight()
		{
		return [x1: 0, y1: 0, x2: 0, y2: .H]
		}

	textPosLeft(tab, selectedTab)
		{
		pos = Object(
			x: tab.renderRect.left + .PaddingTop - 2
			y: tab.renderRect.bottom - .textPadding(tab))
		if not selectedTab and tab.width is tab.renderWidth
			pos.y -= tab.textBoldOffset
		return pos
		}

	textPadding(tab)
		{
		imageWidth = .ImageWidth(tab.image)
		return imageWidth + .PaddingSide
		}

	textPosRight(tab, selectedTab)
		{
		pos = Object(
			x: tab.renderRect.right - .PaddingTop + 2
			y: tab.renderRect.top + .textPadding(tab))
		if not selectedTab and tab.width is tab.renderWidth
			pos.y += tab.textBoldOffset
		return pos
		}

	Resize?(w /*unused*/, h)
		{
		return h isnt .H
		}

	Getter_TabBarSize()
		{
		return .H
		}

	ResizeExtraControl(extraControl, ctrlPos, ctrlSize, xstretch = false)
		{
		if xstretch
			ProgrammerError(
				'Controls with Xstretch should not be used with vertical tabs',
				caughtMsg: 'developer error')
		extraControl.Resize(.extraControlX, ctrlPos, .W - .PaddingSide, ctrlSize)
		}

	ResizeButton(button, pos)
		{
		button.Resize(.PaddingTop - .offset + .buttonOffset, pos,
			.ButtonSize, .ButtonSize)
		return .ButtonSize
		}

	imageOffset: 2
	imagePosLeft(tab)
		{
		x = tab.renderRect.left + .PaddingTop - .imageOffset
		y = tab.renderRect.bottom - .PaddingSide
		return Object(:x, :y)
		}

	imagePosRight(tab)
		{
		x = tab.renderRect.right - .PaddingTop + .imageOffset
		y = tab.renderRect.top + .PaddingSide
		return Object(:x, :y)
		}

	ImageRect(tab)
		{
		return (.imageRectMethod)(tab)
		}

	imageRectLeft(tab)
		{
		imageSpecs = .CalcImageSpecs(tab)
		dimensions = .ImageDimensions(tab)
		right = dimensions.width + left = imageSpecs.x
		top = (bottom = imageSpecs.y) - dimensions.height
		return Object(:left, :right, :top, :bottom)
		}

	imageRectRight(tab)
		{
		imageSpecs = .CalcImageSpecs(tab)
		dimensions = .ImageDimensions(tab)
		left = (right = imageSpecs.x) - dimensions.width
		bottom = dimensions.height + top = imageSpecs.y
		return Object(:left, :right, :top, :bottom)
		}

	InvalidateRect(i, tab)
		{
		renderRect = (.RenderRect)(i, tab, tab.renderRect.top)
		return Object(
			left: renderRect.left,
			right: .TabHeight,
			top: renderRect.top - 1,
			bottom: renderRect.bottom + 1)
		}
	}