// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Component
	{
	ComponentName: 'BulletGraphs'
	style: `
		.su-bulletgraph-container {
			position: relative;
			display: flex;
			justify-content: flex-end;
		}
		.su-bulletgraph-axis {
			position: relative;
			font-size: 80%;
			flex-grow: 1;
		}
		.su-bulletgraph-axis-value{
			position: absolute;
		}
		.su-bulletgraph-axis-mark {
			position: absolute;
			border-width: 1px;
		}
		.su-bulletgraph-bar {
			position: absolute;
		}
		.su-bulletgraph-target {
			position: absolute;
			border-width: 2px;
		}
		.su-bulletgraph-range {
			height: 100%;
			width: 100%;
		}`
	className: 'su-bulletgraph'
	New(.data, .satisfactory = 0, .good = 0, .target = false, .range = #(0, 100),
		.color = 0x506363, .width = 128, .height = 32, .rectangle = true,
		outside = 5, .vertical = false, .axis = false, .axisDensity = 5,
		.axisFormat = false, .selectedColor = false)
		{
		.Name = String(.Name) // Handle non-string names
		.Xmin = .width -= outside * 2
		.Ymin = .height -= outside * 2
		LoadCssStyles('bulletgraph-control.css', .style)
		.CreateElement('div', className: .className)
		.SetStyles(Object(padding: outside $ 'px'), .El)
		.graphContainerEl = .createGraphElement(.El, 'container')
		.setupColors()
		.containerStyle = .vertical
			? .verticalGraph()
			: .horizontalGraph()
		.SetStyles(.containerStyle, .graphContainerEl)
		.SetMinSize()
		}

	createGraphElement(parent, classSuffix)
		{
		return CreateElement('div', parent, className: .className $ '-' $ classSuffix)
		}

	setupColors()
		{
		colors = BulletGraphColors(.color, .rgbStr)
		.badColor = colors.bad
		.satisfactoryColor = colors.satisfactory
		.goodColor = colors.good
		.valueColor = colors.value
		if .selectedColor isnt false
			.selectedColor = ToCssColor(.selectedColor)
		}

	rgbStr(rgbOb)
		{
		return 'rgb(' $ rgbOb.r $ ', ' $ rgbOb.g $ ', ' $ rgbOb.b $ ')'
		}

	// ========================== Draw Vertical ==========================
	verticalGraph()
		{
		.graphSize = .width
		.graphAxisEl = .addGraphAxis(.verticalAxisValueStyle, .verticalAxisMarkStyle)
		.graphEl = .addGraph()
		.graphBarEl = .addGraphBar(.verticalBarStyle)
		.graphTargetEl = .addGraphTarget(.verticalTargetStyle)
		.graphRangeEl = .addGraphRange('to top')
		return Object(
			'flex-direction': 'row',
			'margin-bottom': '4px',
			'max-height': .height $ 'px')
		}

	verticalAxisValueStyle(percent, metrics)
		{
		offset = .axisOffset + 2 /*= spacing between value / mark*/
		valueStyle = Object(
			bottom: .axisPosCalc(percent, metrics.height / 2),
			right: offset $ 'px')
		return valueStyle, metrics.width + offset, 0
		}

	verticalAxisMarkStyle(percent, maxWidth)
		{
		markOffset = percent is 0 /*= start of graph*/
			? 1
			: 0
		return Object(
			bottom: .axisPosCalc(percent, markOffset),
			left: maxWidth - .axisOffset $ 'px',
			right: '0px',
			'border-top-style': 'solid')
		}

	verticalBarStyle(barLength, barSize, barOffset, color)
		{
		return Object(
			top: 100 - barLength $ '%',
			height: barLength $ '%',
			width: barSize $ 'px',
			right: barOffset $ 'px',
			'background-color': color)
		}

	verticalTargetStyle(targetPos, targetSize, targetOffset)
		{
		return Object(
			top: 100 - targetPos $ '%',
			width: targetSize $ 'px',
			right: targetOffset $ 'px'
			height: '1px',
			'border-top-style': 'solid')
		}

	// ========================= Draw Horizontal =========================
	horizontalGraph()
		{
		.graphSize = .height
		.graphEl = .addGraph()
		.graphBarEl = .addGraphBar(.horizontalBarStyle)
		.graphTargetEl = .addGraphTarget(.horizontalTargetStyle)
		.graphRangeEl = .addGraphRange('to right')
		.graphAxisEl = .addGraphAxis(.horizontalAxisValueStyle, .horizontalAxisMarkStyle)
		return Object(
			'flex-direction': 'column',
			'max-width': .width $ 'px')
		}

	horizontalBarStyle(barLength, barSize, barOffset, color)
		{
		return Object(
			right: 100 - barLength $ '%',
			width: barLength $ '%',
			height: barSize $ 'px',
			top: barOffset $ 'px',
			'background-color': color)
		}

	horizontalTargetStyle(targetPos, targetSize, targetOffset)
		{
		return Object(
			right: 100 - targetPos $ '%',
			height: targetSize $ 'px',
			top: targetOffset $ 'px',
			width: '1px',
			'border-right-style': 'solid')
		}

	horizontalAxisValueStyle(percent, metrics)
		{
		valueStyle = Object(
			left: .axisPosCalc(percent, metrics.width / 2),
			top: .axisOffset $ 'px')
		return valueStyle, 0, metrics.height + .axisOffset
		}

	horizontalAxisMarkStyle(percent, maxWidth /*unused*/)
		{
		markOffset = percent is 0 /*= start of graph*/
			? 1
			: 0
		return Object(
			left: .axisPosCalc(percent, markOffset),
			bottom: '-' $ .axisOffset $ 'px',
			top: '0px',
			'border-right-style': 'solid')
		}

	// =========================== Draw General ==========================
	addGraph()
		{
		graphEl = .createGraphElement(.graphContainerEl, 'graph')
		graphEl.SetStyle('width', .width $ 'px')
		graphEl.SetStyle('height', .height $ 'px')
		graphEl.AddEventListener('mouseover', .mouseMove)
		graphEl.AddEventListener('mouseup', .mouseUp)
		if .rectangle
			graphEl.SetStyle(`outline`, `1px solid black`)
		return graphEl
		}

	axisOffset: 8
	addGraphAxis(axisValueStyleMethod, axisMarkStyleMethod)
		{
		if .axis is false
			return false
		graphAxisEl = .createGraphElement(.graphContainerEl, 'axis')
		maxWidth = maxHeight = 0
		r = 1 / Max(.axisDensity - 1, 1) //= Increment Rate
		for i in .. .axisDensity
			{
			width, height = .addGraphAxisValue(r * i, graphAxisEl, axisValueStyleMethod)
			maxWidth = Max(maxWidth, width)
			maxHeight = Max(maxHeight, height)
			}
		// Vertical Graphs: Must calculate the maxWidth prior to adding the axis marks
		// Horizontal Graphs: Construction order does not affect the axis mark positioning
		for i in .. .axisDensity
			.addGraphAxisMark(axisMarkStyleMethod, r * i, maxWidth, graphAxisEl)
		.Xmin += maxWidth
		.Ymin += maxHeight
		return graphAxisEl
		}

	addGraphAxisValue(percent, graphAxisEl, styleMethod)
		{
		value = (.range[1] * percent).Round(0) + .range[0]
		axisValue = .createGraphElement(graphAxisEl, 'axis-value')
		axisValue.innerText = .axisText(.axisFormat, value)
		metrics = SuRender().GetTextMetrics(axisValue, axisValue.innerText)
		styles, width, height = (styleMethod)(percent, metrics)
		.SetStyles(styles, axisValue)
		return width, height
		}

	addGraphAxisMark(styleMethod, percent, maxWidth, graphAxisEl)
		{
		.SetStyles((styleMethod)(percent, maxWidth),
			.createGraphElement(graphAxisEl, 'axis-mark'))
		}

	axisText(format, value)
		{
		if format is false
			return String(value)
		text = value.Format(format)
		start = text.FindRx('[0-9]')
		text = start isnt 0
			? text[.. start] $ text[start ..].Replace('^0*', '')
			: text.Replace('^0*', '')
		return text
		}

	axisPosCalc(percent, offset)
		{
		return 'calc(' $ percent.DecimalToPercent() $ '% - ' $ offset $ 'px)'
		}

	barRatio: 0.35
	addGraphBar(barStyle)
		{
		if .data <= .range[0]
			return false
		barLength = .percentage(.data, .range[1])
		barOffset = (.graphSize * (1 - .barRatio)) / 2
		barSize = .graphSize * .barRatio
		graphBarEl = .createGraphElement(.graphEl, 'bar')
		.SetStyles((barStyle)(barLength, barSize, barOffset, .barColor()), graphBarEl)
		return graphBarEl
		}

	percentage(value, max)
		{
		return (value / max).DecimalToPercent()
		}

	targetRatio: 0.55
	addGraphTarget(targetStyle)
		{
		if .target is false
			return false
		targetPos = .percentage(.target, .range[1])
		targetOffset = (.graphSize * (1 - .targetRatio)) / 2
		targetSize = .graphSize * .targetRatio
		graphTargetEl = .createGraphElement(.graphEl, 'target')
		.SetStyles((targetStyle)(targetPos, targetSize, targetOffset), graphTargetEl)
		return graphTargetEl
		}

	barColor()
		{
		return .selectedColor isnt false and .selected
			? .selectedColor
			: .valueColor
		}

	addGraphRange(dir)
		{
		graphRangeEl = .createGraphElement(.graphEl, 'range')
		bad = .percentage(.satisfactory, .range[1])
		satisfactory = .percentage(.good, .range[1])
		styles = Object(
			background: 'linear-gradient(' $ dir $ ', ' $
				.gradientColor(.badColor, 0, bad) $ ', ' $
				.gradientColor(.satisfactoryColor, bad, satisfactory) $ ', ' $
				.gradientColor(.goodColor, satisfactory) $ ')')
		.SetStyles(styles, graphRangeEl)
		return graphRangeEl
		}

	gradientColor(color, rangeLow, rangeHigh = 100)
		{
		return color $ ' ' $ rangeLow $ '%, ' $ color $ ' ' $ rangeHigh $ '%'
		}

	// ======================== User Interaction =========================
	mouseMove(event /*unused*/)
		{
		.Event('MOUSEMOVE')
		}

	mouseUp(event)
		{
		if event.button is 0
			.Event('LBUTTONUP')
		}

	selected: false
	Selected(selected)
		{
		if .graphBarEl is false
			return
		setStyle? = .selected isnt selected
		.selected = selected
		if setStyle?
			.graphBarEl.SetStyle('background-color', .barColor())
		}
	}