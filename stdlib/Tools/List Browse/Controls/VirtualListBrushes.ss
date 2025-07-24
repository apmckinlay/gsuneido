// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{
		.brushes = Object()
		.highlights = Object()
		.highlightRecords = Object()
		}

	HighlightValues(member, values, color)
		{
		if not .highlights.Member?(member)
			.highlights[member] = Object()

		for val in values
			.highlights[member][val] = .createBrush(color)
		}

	HighlightRecords(records, color)
		{
		for rec in records
			{
			.ClearHighlightRecord(rec)
			.highlightRecords.Add([:rec, color: .createBrush(color)])
			}
		}

	ClearHighlightRecord(rec)
		{
		.highlightRecords.RemoveIf({ rec is it.rec })
		}

	createBrush(color)
		{
		if not .brushes.Member?(color)
			.brushes[color] = .createSystemBrush(color)
		return .brushes[color]
		}

	createSystemBrush(color)
		{
		CreateSolidBrush(color)
		}

	GetBrush(row)
		{
		for member in .highlights.Members()
			return .highlights[member].GetDefault(row[member], false)
		if false isnt found = .highlightRecords.FindOne({ it.rec is row })
			return found.color
		return false
		}

	Destroy()
		{
		.brushes.Filter({ it isnt false }).Each(DeleteObject)
		}
	}
