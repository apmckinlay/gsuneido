// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
/*
This class converts a DrawItem format object into an SVG/HTML string suitable
	for web display. It works by:

1. Setting CanvasToSvgDriver (an HtmlDriver subclass) as _report so the canvas Paint()
	calls dispatch into a record list.
2. Walking that record list and dispatching each item by method name to render as SVG
	primitives (image, line, rect, ellipse, text).
3. Wrapping the SVG in a <div> with flex justify-content corresponding to the original
	DrawControl.ApplyJustify rule.
*/
class
	{
	CallClass(ob, logo)
		{
		return new this().Process(ob, logo)
		}

	New()
		{
		.Driver = new CanvasToSvgDriver
		}

	Process(ob, logo)
		{
		_report = .Driver
		if false is drawItem = ob.FindOne({ Object?(it) and it[0] is 'DrawItem' })
			return ''

		canvasGroup = drawItem[1]
		// from DrawItemFormat. Used to convert canvas unit to twip
		conversion = 17 / canvasGroup.ScaleBy

		rect = canvasGroup.BoundingRect()

		dx = -rect.x1 * conversion
		dy = -rect.y1 * conversion
		canvasGroup.Scale(conversion, print?:)
		canvasGroup.Move(dx, dy)
		canvasGroup.Paint()

		_env = Object(xMax: 0, yMax: 0)
		s = ''
		for item in .Driver.Page
			s $= (this[item[0]])(@+1item) $ '\n'
		w = Max(_env.xMax, (rect.x2 - rect.x1).Abs() * conversion / .factor)
		h = Max(_env.yMax, (rect.y2 - rect.y1).Abs() * conversion / .factor)

		viewBox = '0 0 ' $ w $ ' ' $ h
		svg = Xml('svg', s,
			xmlns: 'http://www.w3.org/2000/svg',
			:viewBox,
			style: 'width:100%;height:auto;max-width:' $ w $ 'px;')

		return Xml('div', svg,
			style: 'display:flex;justify-content:' $ .justify(ob, drawItem) $ ';',
			'data_logo': logo)
		}

	justify(ob, drawItem)
		{
		// revert stdlib:DrawControl.ApplyJustify
		return ob[1] is drawItem
			? 'flex-start'
			: ob.GetDefault(3/*=pos*/, false) is 'Hfill'
				? 'center'
				: 'flex-end'
		}

	/*
	 * Following methods are called on CanvasToSvgDriver results by .Process
	 * .factor is to convert coordinates in twip to html pixel
	*/
	factor:  15 /*=twipToPixel = 1440 / 96*/
	AddImage(x, y, w, h, data)
		{
		x /= .factor
		y /= .factor
		w /= .factor
		h /= .factor
		.updateEnv(x + w, y + h)
		return Xml('image', href: data,
			:x, :y, width: w, height: h, preserveAspectRatio: 'none')
		}

	AddLine(x, y, x2, y2, thick, color = 0)
		{
		x /= .factor
		y /= .factor
		x2 /= .factor
		y2 /= .factor
		.updateEnv(x, y)
		.updateEnv(x2, y2)
		return Xml('line', x1: x, y1: y, :x2, :y2,
			style: 'stroke-width: ' $ thick $ ';stroke: ' $ ToCssColor(color))
		}

	AddRect(x, y, w, h, thick, fillColor = false, lineColor = false)
		{
		x /= .factor
		y /= .factor
		w /= .factor
		h /= .factor
		.updateEnv(x + w, y + h)
		return Xml('rect', :x, :y, width: w, height: h,
			style: .shapeStyle(thick, fillColor, lineColor))
		}

	AddRoundRect(x, y, w, h, width = 0, height = 0, thick = 1,
		fillColor = false, lineColor = false)
		{
		x /= .factor
		y /= .factor
		w /= .factor
		h /= .factor
		.updateEnv(x + w, y + h)
		return Xml('rect', :x, :y, width: w, height: h,
			rx: width / 2 / .factor, ry: height / 2 / .factor,
			style: .shapeStyle(thick, fillColor, lineColor))
		}

	AddEllipse(x, y, w, h, thick = 1, fillColor = false, lineColor = false)
		{
		x /= .factor
		y /= .factor
		w /= .factor
		h /= .factor
		.updateEnv(x + w, y + h)
		return Xml('ellipse',
			cx: x + w / 2, cy: y + h / 2, rx: w / 2, ry: h / 2,
			style: .shapeStyle(thick, fillColor, lineColor))
		}

	shapeStyle(thick, fillColor, lineColor)
		{
		return 'stroke-width: ' $ thick $
			';fill: ' $ (fillColor isnt false ? ToCssColor(fillColor) : 'none') $
			';stroke: ' $ (lineColor isnt false ? ToCssColor(lineColor) : 'black')
		}

	AddText(data, x, y, w, h, font = #(), justify = 'left',
		ellipsis? = false, color = false, html/*unused*/ = false)
		{
		origX = x
		font = PdfFonts.GetCompatibleFont(font)
		data = PdfFonts.StripInvalidChars(data)
		if ellipsis? is true
			data = PdfDriver.PdfEllipsis(font, data, w)
		textWidth = PdfDriver.GetTextWidth(font, data)
		xAdjust = PdfDriver.JustifyAdjust(textWidth, w, justify)
		x += xAdjust

		x /= .factor
		y /= .factor
		.updateEnv(Max(x + textWidth / .factor, (origX + w) / .factor),
			y + h / .factor)
		style = 'font-size: ' $ font.GetDefault(#size, 10 /*= font size */) $ 'pt' $
			';font-family: ' $ font.GetDefault(#name, 'Arial') $
			';font-weight: ' $ font.GetDefault(#weight, 400 /*= font weight */)
		if font.GetDefault(#italic, false) is true
			style $= ';font-style: italic'
		if font.GetDefault(#underline, false) is true
			style $= ';text-decoration: underline'
		else if font.GetDefault(#strikeout, false) is true
			style $= ';text-decoration: line-through'
		if color isnt false
			style $= ';fill: ' $ ToCssColor(color)
		return Xml('text', data, :x, :y, :style)
		}

	updateEnv(x2, y2)
		{
		if _env.xMax < x2
			_env.xMax = x2
		if _env.yMax < y2
			_env.yMax = y2
		}
	}
