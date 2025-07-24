// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// A class to share functionality between Sparkline Control and Format - by Mauro Giubileo
class
	{
	CallClass(@args)
		{
		return new this(@args)
		}

	New(.data, .inside, .thick, .rectangle, .middleLine,
		.allPoints, .firstPoint, .lastPoint, .lines, .minPoint, .maxPoint, .normalRange,
		.normalRangeColor, .borderOnPoints, .circlePoints, .pointLineRatio)
		{
		}

	Draw(.x, .y, .w, .h)
		{
		.PreDraw(.circlePoints, .borderOnPoints, .thick)
		if .rectangle // Draw enclosing rectangle
			.DrawRectangle(.x, .y, .x + .w, .y + .h)

		if .data.Empty?()
			return

		// add inside spacing:
		.w -= .inside * 2
		.h -= .inside * 2
		.x += .inside
		.y += .inside

		.x_incr = .w / (.data.Size() - 1)
		.y_incr = .h / (Abs(.data.Max() - .data.Min()))

		.DrawGraphs(.x, .y, .w, .h)
		}

	PreDraw(@unused) /*usage: dc, circlePoints, borderOnPoints, thick*/
		{ }

	DrawRectangle(@unused) /*usage: x1, y1, x2, y2*/
		{
		throw 'must be implemented by derived class'
		}

	DrawGraphs(.x, .y, .w, .h)
		{
		if Object?(.normalRange)
			.drawNormalBand()

		if .middleLine
			.DrawMiddleLine()

		if .lines
			.DrawLines(.x, .y + .h - .Normalize(.data[0]), .data, .x_incr)

		if .allPoints
			for (pos = 0; pos < .data.Size(); pos++)
				.drawSinglePoint(pos, CLR.Highlight, #pos)

		.drawSpecificPoints()
		}

	drawNormalBand()
		{
		// check that range doesn't exceed graph limits:
		checkedNormalRange = Object(.normalRange[0], .normalRange[1])
		if .normalRange[0] < .data.Min()
			checkedNormalRange[0] = .data.Min()
		if .normalRange[1] > .data.Max()
			checkedNormalRange[1] = .data.Max()

		checkedNormalRange.Map!(.Normalize)
		.DrawNormalBand(checkedNormalRange, .normalRangeColor)
		}

	Normalize(v)
		{
		return (v - .data.Min()) * .y_incr
		}

	DrawNormalBand(@unused) /*usage: checkedNormalRange, normalRangeColor*/
		{
		throw 'must be implemented by derived class'
		}

	DrawMiddleLine()
		{
		throw 'must be implemented by derived class'
		}

	DrawLines(@unused) /*usage: x, y, data, offset*/
		{
		throw 'must be implemented by derived class'
		}

	drawSinglePoint(val, color, search)
		{
		//search can be 'find' or 'pos'
		if not #(find, pos).Has?(search)
			throw `search param has to be 'find' or 'pos'`
		half1 = .thick * .pointLineRatio / 2
		half2 = .thick * .pointLineRatio / 2
		newPos = .calcNewPos(search, val)
		.DrawPoint(newPos.x, newPos.y, half1.Floor(), half2.Ceiling(), color)
		}

	calcNewPos(search, val)
		{
		if search is 'find'
			{
			x = .x + .x_incr * .data.Find(val)
			y = .y + .h - .Normalize(val)
			}
		else
			{
			x = .x + .x_incr * val
			y = .y + .h - .Normalize(.data[val])
			}
		return [:x, :y]
		}

	DrawPoint(@unused) /*usage: x, y, half1, half2, color*/
		{
		throw 'must be implemented by derived class'
		}

	drawSpecificPoints()
		{
		if .minPoint
			.drawSinglePoint(.data.Min(), CLR.Highlight, #find)
		if .maxPoint
			.drawSinglePoint(.data.Max(), CLR.Highlight, #find)
		if .firstPoint
			.drawSinglePoint(0, CLR.red, #pos)
		if .lastPoint
			.drawSinglePoint(.data.Size() - 1, CLR.red, #pos)
		}
	}
