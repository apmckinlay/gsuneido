// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(datalist, from, to, span, format = '-###,###,###')
		{
		return .controls(datalist, from, to, span, format)
		}
	controls(datalist, from, to, span, format)
		{
		cells = []
		cells.Add(.header(from, to, span))
		i = 0
		for data in datalist
			cells.Add(.row(data, format, i++))
		return ['Grid', cells]
		}
	header(from, to, span)
		{
		return [
			['Static', from, size: '-2'],
			['Static', span, size: '-2', xstretch: 1, justify: 'CENTER'],
			['Static', to, size: '-2', xstretch: 1, justify: 'CENTER'],
			#(Static, 'low', size: '-2', xstretch: 1, justify: 'CENTER')
			#(Static, 'high', size: '-2', xstretch: 1, justify: 'CENTER')
			]
		}
	row(data, format, i)
		{
		red = CLR.RED
		blue = CLR.BLUE
		return [
			['Static', data.First().Format(format), color: red,
				xstretch: 1, justify: 'RIGHT', name: "Front " $ i],
			['Sparkline', data,
				rectangle: false, circlePoints:,
				firstPoint:, lastPoint:, minPoint:, maxPoint:],
			['Static', data.Last().Format(format), color: red,
				xstretch: 1, justify: 'RIGHT', name: "Back " $ i],
			['Static', data.Min().Format(format), color: blue,
				xstretch: 1, justify: 'RIGHT', name: "Min " $ i],
			['Static', data.Max().Format(format), color: blue,
				xstretch: 1, justify: 'RIGHT', name: "Max " $ i]
			]
		}
	}