// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	styles: `
		.su-preview-access-point {
			fill: white;
			fill-opacity: 0;
		}
		.su-preview-access-point:hover {
			stroke: black;
			stroke-dasharray: 2;
			cursor: pointer;
		}
		.su-svg-image {
			pointer-events: none;
		}`
	free3of9: `
		@import url('https://fonts.googleapis.com/css2?` $
			`family=Libre+Barcode+39&display=swap');`
	free3of9x: `
		@import url('https://fonts.googleapis.com/css2?` $
			`family=Libre+Barcode+39+Extended&display=swap');`

	Namespace: 'http://www.w3.org/2000/svg'
	pixelPerInch: 96
	// el should be a <svg> element
	New(.el, dimens, .factor = 15 /*=twipToPixel = 1440 / 96*/)
		{
		LoadCssStyles('su-svg-driver', .styles)
		.el.innerHTML = ''
		.el.SetAttribute('viewBox', '0 0 ' $
			.ConvertToPixcel(dimens.width) $ ' ' $
			.ConvertToPixcel(dimens.height))
		}

	ConvertToPixcel(n, scale = 1)
		{
		return n * .pixelPerInch * scale
		}

	AddText(data, x, y, w, h/*unused*/, font, justify = 'left',
		ellipsis? = false, color = false, html = false)
		{
		w = w / .factor
		textEl = CreateElement('text', .el, namespace: .Namespace)
		textEl.style['white-space'] = 'pre'

		fontName = .handleFont(font)
		textEl.style['font-size'] = font.size $ 'pt'
		textEl.style['font-family'] = fontName
		textEl.style['font-weight'] = font.GetDefault(#weight, 400) /*= font weight */
		textEl.style['font-style'] =
			font.GetDefault(#italic, false) is true ? 'italic' : 'normal'
		if font.GetDefault(#underline, false) is true
			textEl.style['text-decoration'] = 'underline'
		else if font.GetDefault(#strikeout, false) is true
			textEl.style['text-decoration'] = 'line-through'
		if color isnt false
			textEl.SetAttribute('fill', ToCssColor(color))
		if ellipsis?
			data = .ellipsis(textEl, data, w)
		textWidth = SuRender().GetTextMetrics(textEl, data).width
		xAdjust = .justifyAdjust(textWidth, w, justify)
		textEl.SetAttribute('x', x / .factor + xAdjust )
		textEl.SetAttribute('y', y / .factor)
		if html is true
			textEl.innerHTML = data
		else
			textEl.textContent = data
		return textEl
		}

	handleFont(font)
		{
		fontName = font.name
		if fontName is 'Free 3 of 9 Regular'
			{
			fontName = '"Libre Barcode 39"'
			LoadCssStyles('su-svg-driver-free3of9', .free3of9, media: 'both')
			}
		else if fontName is 'Free 3 of 9 Extended Regular'
			{
			fontName = '"Libre Barcode 39 Extended"'
			LoadCssStyles('su-svg-driver-free3of9x', .free3of9x, media: 'both')
			}
		return fontName
		}

	justifyAdjust(textSize, maxWidth, justify)
		{
		if justify is 'left'
			return 0
		else if justify is 'right'
			return maxWidth - textSize
		else if justify is 'center'
			return (maxWidth / 2) - (textSize / 2)
		else
			return 0
		}

	ellipsis(textEl, data, w)
		{
		// .Ceiling is needed because SuRender().GetTextMetrics does .Cailing on width
		w = w.Ceiling()
		do
			{
			textSize = SuRender().GetTextMetrics(textEl, data).width
			if textSize <= w
				return data
			if not data.Suffix?('...')
				data = data[..-2] $ '...' // remove 2 chars to make up for ...
			data = data[..-4] $ '...' // remove 1 char
			} while data.Size() > 3
		return data
		}

	MoveText(textEl, dx, dy)
		{
		textEl.SetAttribute('x', textEl.GetAttribute('x') + dx / .factor)
		textEl.SetAttribute('y', textEl.GetAttribute('y') + dy / .factor)
		}

	ResizeText(textEl, x, y, w, h/*unused*/, justify = 'left', ellipsis? = false)
		{
		data = textEl.textContent
		textWidth = SuRender().GetTextMetrics(textEl, data).width
		xAdjust = .justifyAdjust(textWidth, w, justify)
		if ellipsis?
			data = .ellipsis(textEl, data, w)
		textEl.SetAttribute('x', x / .factor + xAdjust)
		textEl.SetAttribute('y', y / .factor)
		}

	AddLine(x, y, x2, y2, thick, color = 0x00000000)
		{
		lineEl = CreateElement('line', .el, namespace: .Namespace)
		.ResizeLine(lineEl, x, y, x2, y2)
		lineEl.style['stroke-width'] = thick / .factor
		lineEl.style.stroke = ToCssColor(color)
		return lineEl
		}

	MoveLine(lineEl, dx, dy)
		{
		lineEl.SetAttribute('x1', lineEl.GetAttribute('x1') + dx / .factor)
		lineEl.SetAttribute('y1', lineEl.GetAttribute('y1') + dy / .factor)
		lineEl.SetAttribute('x2', lineEl.GetAttribute('x2') + dx / .factor)
		lineEl.SetAttribute('y2', lineEl.GetAttribute('y2') + dy / .factor)
		}

	ResizeLine(lineEl, x, y, x2, y2)
		{
		lineEl.SetAttribute('x1', x / .factor)
		lineEl.SetAttribute('y1', y / .factor)
		lineEl.SetAttribute('x2', x2 / .factor)
		lineEl.SetAttribute('y2', y2 / .factor)
		}

	AddRect(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		rectEl = CreateElement('rect', .el, namespace: .Namespace)
		.ResizeRect(rectEl, x, y, w, h)
		rectEl.style['stroke-width'] = thick / .factor

		rectEl.style.fill = fillColor isnt false ? ToCssColor(fillColor) : 'none'
		rectEl.style.stroke = lineColor isnt false ? ToCssColor(lineColor) : 'black'
		return rectEl
		}

	MoveRect(rectEl, dx, dy)
		{
		rectEl.SetAttribute('x', rectEl.GetAttribute('x') + dx / .factor)
		rectEl.SetAttribute('y', rectEl.GetAttribute('y') + dy / .factor)
		}

	ResizeRect(rectEl,x, y, w, h)
		{
		rectEl.SetAttribute('x', x / .factor)
		rectEl.SetAttribute('y', y / .factor)
		rectEl.SetAttribute('width', w / .factor)
		rectEl.SetAttribute('height', h / .factor)
		}

	AddImage(x, y, w, h, data)
		{
		imageEl = CreateElement('image', .el, namespace: .Namespace)
		imageEl.SetAttribute('class', 'su-svg-image')
		.ResizeImage(imageEl, x, y, w, h)
		imageEl.SetAttribute('href', data)
		imageEl.SetAttribute('preserveAspectRatio', 'none')
		return imageEl
		}

	MoveImage(imageEl, dx, dy)
		{
		imageEl.SetAttribute('x', imageEl.GetAttribute('x') + dx / .factor)
		imageEl.SetAttribute('y', imageEl.GetAttribute('y') + dy / .factor)
		}

	ResizeImage(imageEl, x, y, w, h)
		{
		imageEl.SetAttribute('x', x / .factor)
		imageEl.SetAttribute('y', y / .factor)
		imageEl.SetAttribute('width', w / .factor)
		imageEl.SetAttribute('height', h / .factor)
		}

	AddRoundRect(x, y, w, h, width = 0, height = 0, thick = 1,
		fillColor = false, lineColor = false)
		{
		el = .AddRect(x, y, w, h, thick, fillColor, lineColor)
		el.SetAttribute('rx', width / 2 / .factor)
		el.SetAttribute('ry', height / 2 / .factor)
		return el
		}

	MoveRoundRect(roundRectEl, dx, dy)
		{
		.MoveRect(roundRectEl, dx, dy)
		}

	ResizeRoundRect(roundRectEl, x, y, w, h, width = 0, height = 0)
		{
		.ResizeRect(roundRectEl, x, y, w, h)
		roundRectEl.SetAttribute('rx', width / 2 / .factor)
		roundRectEl.SetAttribute('ry', height / 2 / .factor)
		}

	AddEllipse(x, y, w, h, thick = 1, fillColor = false, lineColor = false)
		{
		ellipseEl = CreateElement('ellipse', .el, namespace: .Namespace)
		.ResizeEllipse(ellipseEl, x, y, w, h)

		ellipseEl.style['stroke-width'] = thick / .factor

		ellipseEl.style.fill = fillColor isnt false ? ToCssColor(fillColor) : 'none'
		ellipseEl.style.stroke = lineColor isnt false ? ToCssColor(lineColor) : 'black'
		return ellipseEl
		}

	MoveEllipse(ellipseEl, dx, dy)
		{
		ellipseEl.SetAttribute('cx', ellipseEl.GetAttribute('cx') + dx / .factor)
		ellipseEl.SetAttribute('cy', ellipseEl.GetAttribute('cy') + dy / .factor)
		}

	ResizeEllipse(ellipseEl, x, y, w, h)
		{
		ellipseEl.SetAttribute('cx', (x + w / 2) / .factor)
		ellipseEl.SetAttribute('cy', (y + h / 2) / .factor)
		ellipseEl.SetAttribute('rx', w / 2 / .factor)
		ellipseEl.SetAttribute('ry', h / 2 / .factor)
		}

	AddArc(left, top, right, bottom,
		xStartArc = 0, yStartArc = 0, xEndArc = 0, yEndArc = 0,
		thick = 1, lineColor = false)
		{
		points = .getClipPoints(left, top, right, bottom,
			xStartArc, yStartArc, xEndArc, yEndArc)
		.drawWithinClip(points)
			{ |clipId|
			ellipseEl = .AddEllipse(left, top, right - left, bottom - top,
				:thick, :lineColor)
			ellipseEl.SetAttribute('clip-path', 'url(#' $ clipId $ ')')
			}
		return ellipseEl
		}

	getClipPoints(left, top, right, bottom,
		xStartArc = 0, yStartArc = 0, xEndArc = 0, yEndArc = 0)
		{
		points = Object(curPoint = Object(xStartArc, yStartArc))
		endPoint = Object(xEndArc, yEndArc)
		do
			{
			curPoint = .nextPoint(curPoint, endPoint, left, top, right, bottom)
			points.Add(curPoint)
			}
		while (curPoint isnt endPoint)
		return points
		}

	nextPoint(curPoint, endPoint, left, top, right, bottom)
		{
		if .onSameEdge?(curPoint, endPoint)
			return endPoint
		if curPoint[0] is left and curPoint[1] isnt bottom
			return Object(left, bottom)
		if curPoint[1] is bottom and curPoint[0] isnt right
			return Object(right, bottom)
		if curPoint[0] is right and curPoint[1] isnt top
			return Object(right, top)
		return Object(left, top)
		}

	onSameEdge?(point1, point2)
		{
		return point1[0] is point2[0] or point1[1] is point2[1]
		}

	MoveArc(arcEl, dx, dy)
		{
		clipEl = .el.GetElementById(arcEl.GetAttribute('clip-path')[5/*=start pos*/..-1])
		polygon = clipEl.firstChild
		newPoints = polygon.GetAttribute('points').
			Split(' ').
			Map({ it.Split(',') }).
			Map({ Object(it[0] + dx / .factor, it[1] + dy / .factor) }).
			Map({ it[0] $ ',' $ it[1] }).
			Join(' ')
		polygon.SetAttribute('points', newPoints)
		.MoveEllipse(arcEl, dx, dy)
		}

	ResizeArc(arcEl, left, top, right, bottom,
		xStartArc = 0, yStartArc = 0, xEndArc = 0, yEndArc = 0)
		{
		points = .getClipPoints(left, top, right, bottom,
			xStartArc, yStartArc, xEndArc, yEndArc)
		clipEl = .el.GetElementById(arcEl.GetAttribute('clip-path')[5/*=start pos*/..-1])
		polygon = clipEl.firstChild
		factor = .factor
		polygon.SetAttribute('points',
			points.Map({ (it[0] / factor) $ ',' $ (it[1] / factor) }).Join(' '))
		.ResizeEllipse(arcEl, left, top, right - left, bottom - top)
		}

	AddAccessPoint(left, top, width, height, idx, callFn = false)
		{
		rectEl = CreateElement('rect', .el, namespace: .Namespace)
		rectEl.SetAttribute('class', 'su-preview-access-point')
		rectEl.SetAttribute('x', left / .factor)
		rectEl.SetAttribute('y', top / .factor)
		rectEl.SetAttribute('width', width / .factor)
		rectEl.SetAttribute('height', height / .factor)
		rectEl.SetAttribute('data-idx', idx)
		if callFn isnt false
			rectEl.AddEventListener('dblclick', callFn)
		return rectEl
		}

	Remove(el)
		{
		el.Remove()
		}

	MoveToBack(el)
		{
		.el.InsertBefore(el, .el.firstChild)
		}

	MoveToFront(el)
		{
		.el.AppendChild(el)
		}

	drawWithinClip(points, block)
		{
		clipPath = CreateElement('clipPath', .el, namespace: .Namespace)
		clipId = 'suClip' $ Suneido.GetInit('SuClipId', 0)
		Suneido.SuClipId++
		clipPath.SetAttribute('id', clipId)
		polygon = CreateElement('polygon', clipPath, namespace: .Namespace)
		factor = .factor
		p = points.Map({ (it[0] / factor) $ ',' $ (it[1] / factor) }).Join(' ')
		polygon.SetAttribute('points', p)
		block(clipId)
		}
	}
