// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Split
	{
	Name:	"Split"
	new? :	true
	New(first, second, .associate, candrag? = true, initialsplit = false,
		splitname = false, splitter = false)
		{
		// Generate Split with HandleSplitterControl or splitter
		super(first, second, .create(associate), :splitter)
		.ctrls = .GetChildren()

		// Set correct Xmin/Ymin value of first/second child item
		switch (associate)
			{
		case "east": 	.min = .ctrls[2].Xmin; .ctrls[2].Xmin = 0
		case "west": 	.min = .ctrls[0].Xmin; .ctrls[0].Xmin = 0
		case "south":	.min = .ctrls[2].Ymin; .ctrls[2].Ymin = 0
		case "north":	.min = .ctrls[0].Ymin; .ctrls[0].Ymin = 0
			}
		// Other intialization
		.candrag? = candrag?
		if initialsplit is false
			.open = .GetSplit()
		else
			.open = initialsplit // object as returned by Split.GetSplit()
		if splitname isnt false
			.ctrls[1].SplitName = splitname
		}
	create(associate)
		{
		.Dir = (associate is "east" or associate is "west") ? "horz" : "vert"
		return true
		}
	Movesplit(n)
		{
		.prevsplit = .GetSplit()
		super.Movesplit(n)
		split = .GetSplit()
		if split[0] isnt 0 and split[1] isnt 0
			.prevsplit = split
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		if .new?
			{
			// Set the initial splitter position
			pos = (.Associate is "south") ? h  - .ctrls[1].Ymin :
				(.Associate is "east") ? w - .ctrls[1].Xmin : 0
			.Movesplit(pos)
			.new? = false
			}
		}
	Getter_CanDrag?()
		{ return .candrag? }
	Getter_Associate()
		{ return .associate }
	Getter_Open?()
		{
		switch (.associate)
			{
		case 'west' :
			return .open?(0, #X)
		case 'east' :
			return .open?(2, #X)
		case 'north' :
			return .open?(0, #Y)
		case 'south' :
			return .open?(2, #Y)
			}
		}
	open?(ctrlNum, axis)
		{
		return .ctrls[ctrlNum][axis $ #stretch] isnt 0 or
			   .ctrls[ctrlNum][axis $ #min] isnt 0
		}
	prevsplit: false
	GetPrevSplit()
		{
		return .prevsplit
		}
	SetPrevSplit(split)
		{
		.prevsplit = split
		}
	Open()
		{
		if .candrag?
			.SetSplit(.open is true ? .prevsplit : .open)
		else
			{
			c = .associate is 'west' or .associate is 'north' ? 0 : 2
			m = .Dir is 'vert' ? 'Ymin' : 'Xmin'
			.ctrls[c][m] = .min
			.SetSplit(.GetSplit())
			}
		.open = true
		}
	Close()
		{
		.open = .GetSplit()
		if not .candrag?
			{
			c = .associate is 'west' or .associate is 'north' ? 0 : 2
			m = .Dir is 'vert' ? 'Ymin' : 'Xmin'
			.ctrls[c][m] = 0
			}
		.new? = true
		.Resize(.Group_x, .Group_y, .Group_w, .Group_h)
		}
	}
