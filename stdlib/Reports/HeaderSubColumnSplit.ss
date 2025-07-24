// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
// Allows block to give a header per x amount of subcolumns for a report
function (numColumns, colHeads, block)
	{
	widths = colHeads[2]
	fmt = Object('Vert')
	hdr = Object('Horz')
	skip = ['Hskip', (widths.skip / 1.InchesInTwips()) * numColumns]

	hdr.Add(Object('Text', '', w: widths[0]))
	if false isnt (font = colHeads.GetDefault('font', false))
		{
		font = font.Copy()
		font.weight = 'bold'
		}

	subColumnWidth = 0
	for (i=1; i<=numColumns; i++)
		subColumnWidth += widths[i]

	block(hdr, skip, subColumnWidth, font)

	fmt.Add(hdr)
	fmt.Add(colHeads)
	return fmt
	}