// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.

// Warning: does not stretch to fill page like GridFormat

// Example:
// Params.On_Preview(#(GridGen #(
// 	#( #('Text' 'Cell 1') #('Text' 'Cell 2') #('Text' 'Cell 2'))
// 	#( #('Text' 'Cell 3', span: 1) #('Text' 'Cell 4'))
// 	#( #('Text' 'Cell 4'))
// 	), skip: 35
//))

Generator
	{
	New(.formats, .width = 0, .font = false, skip = false, .access = false)
		{
		.skip = skip is false ? HskipFormat.Size.w : skip
		.colwidths = Object().Set_default(0)
		for rowfmts in formats
			{
			pos = 0
			for col in rowfmts.Members()
				{
				item = rowfmts[col]
				if font isnt false
					item = .getItemWithFont(item, font)
				sizeW = .getSizeW(item)
				if sizeW > .colwidths[pos]
					.colwidths[pos] = sizeW
				++pos
				if item.Member?('span') and item.span > 0
					pos += item.span
				}
			}
		.row = 0
		}
	getItemWithFont(item, font)
		{
		item = Object?(item) ? item.Copy() : Object(item)
		item.font = font
		return item
		}
	getSizeW(item)
		{
		return _report.PlainText?() ? 0 : _report.Construct(item).GetSize().w
		}
	Next()
		{
		if .row >= .formats.Size()
			return false
		rowfmt = .formats[.row]
		if rowfmt.Size() is 1
			fmt = rowfmt[0]
		else
			{
			fmt = Object('Row')
			for col in rowfmt.Members()
				{
				colformat = rowfmt[col]
				fmt.Add(colformat)
				for (i = 0; i < colformat.GetDefault('span', 0); i++)
					fmt.Add(#(Text, ''))
				}
			fmt.widths = .colwidths
			fmt.widths.skip = .skip
			}
		if .access isnt false and .access.Member?(.row)
			fmt.access = .access[.row]
		++.row
		if .font isnt false
			fmt.font = .font
		return _report.Construct(fmt)
		}
	}
