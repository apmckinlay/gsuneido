// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (x1, y1)
	{
	result = DrawTextAsk(title: 'Text')
	return result is false or result.text is ""
		? false : CanvasText(result.text, x1, y1, x1, y1, result.font,
			justify: result.justify, fromDraw:)
	}