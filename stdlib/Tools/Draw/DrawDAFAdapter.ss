// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
function (x1, y1, canvas)
	{
	if false is result = DrawDAFAsk(canvas, showPrompt?:)
		return false

	_canvas = canvas
	items = Object()
	if result.showPrompt is true
		{
		text = CanvasText(result.prompt $ ':', x1, y1, x1, y1, font: result.font,
			fromDraw:)
		items.Add(text)
		rect = text.BoundingRect()
		x1 = rect.x2 + ScaleWithDpiFactor(10/*=gap*/)
		}
	items.Add(CanvasDAF(x1, y1, x1, y1, result.field, result.font, result.justify,
		fromDraw:))
	return items
	}
