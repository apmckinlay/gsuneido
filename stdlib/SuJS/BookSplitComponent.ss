// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
SplitComponent
	{
	Name:	"Split"
	ctrls: 	false
	New(@args)
		{
		super(@.create(args))
		.ctrls = .GetChildren()

		if args.GetDefault(#splitname, false) isnt false
			.ctrls[1].SplitName = args.splitname

		.Recalc()
		.open = .GetSplit()
		}

	Recalc()
		{
		super.Recalc()
		if .ctrls is false
			return
		// Set correct Xmin/Ymin value of first/second child item
		switch (.associate)
			{
		case "east": 	.min = .ctrls[2].Xmin; .ctrls[2].Xmin = 0; .ctrls[2].SetMinSize()
		case "west": 	.min = .ctrls[0].Xmin; .ctrls[0].Xmin = 0; .ctrls[0].SetMinSize()
		case "south":	.min = .ctrls[2].Ymin; .ctrls[2].Ymin = 0; .ctrls[2].SetMinSize()
		case "north":	.min = .ctrls[0].Ymin; .ctrls[0].Ymin = 0; .ctrls[0].SetMinSize()
			}
		}

	create(args)
		{
		.associate = args.associate
		.Dir = .associate in (#east, #west) ? #horz : #vert
		return args
		}

	Splitter_mouseup()
		{
		.prevsplit = .GetSplit()
		super.Splitter_mouseup()
		split = .GetSplit()
		if split[0] isnt 0 and split[1] isnt 0
			.prevsplit = split
		.updateHandleSplitter()
		}

	prevsplit: false
	GetPrevSplit()
		{
		return .prevsplit
		}

	SetPrevSplit(split)
		{
		.prevsplit = split
		.updateHandleSplitter()
		}

	Getter_Associate()
		{
		return .associate
		}

	Getter_Open?()
		{
		n = .GetSplit()
		return .associate in (#east, #south) ? n[1] isnt 0 : n[0] isnt 0
		}

	Open()
		{
		n = .open is true ? .prevsplit : .open
		.SetSplit(n)
		.Event('UpdateSplit', n)
		.open = true
		}

	Close()
		{
		n = .close()
		.Event('UpdateSplit', n)
		}

	close()
		{
		n = .associate in (#east, #south) ? [1, 0] : [0, 1]
		.open = .GetSplit()
		.SetSplit(n)
		return n
		}

	SetSplit(n)
		{
		super.SetSplit(n)
		.updateHandleSplitter()
		}

	updateHandleSplitter()
		{
		if not .ctrls[1].Method?(#UpdateButton)
			return
		.ctrls[1].UpdateButton()
		}
	}
