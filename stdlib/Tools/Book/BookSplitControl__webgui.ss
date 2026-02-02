// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Split
	{
	Name: "Split"
	ComponentName: 'BookSplit'
	New(first, second, .associate, candrag?/*unused*/ = true,
		initialsplit/*unused*/ = false, splitname = false, splitter = false)
		{
		super(first, second, handle?:, :splitter)
		.ctrls = .GetChildren()

		if splitname isnt false
			.ctrls[1].SplitName = splitname

		.ComponentArgs.associate = associate
		.ComponentArgs.splitname = splitname
		.SetSplit(.associate in ("south", "east") ? [1, 0] : [0, 1])
		}

	Getter_Open?()
		{
		n = .GetSplit()
		return .associate in (#east, #south) ? n[1] isnt 0 : n[0] isnt 0
		}

	Open()
		{
		.Act('Open')
		}

	Close()
		{
		.Act('Close')
		}

	UpdateSplit(n)
		{
		.prevsplit = .GetSplit()
		super.UpdateSplit(n)
		split = .GetSplit()
		if split[0] isnt 0 and split[1] isnt 0
			.prevsplit = split
		}

	prevsplit: #(0.5, 0.5)
	GetPrevSplit()
		{
		return .prevsplit
		}
	SetPrevSplit(split)
		{
		.prevsplit = split
		.Act(#SetPrevSplit, split)
		}
	}