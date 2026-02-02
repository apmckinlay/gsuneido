// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Group
	{
	ComponentName: 'Split'
	New(first, second, handle? = false, splitter = false, .splitSaveName = false,
		.splitSaveNameSuffix = ' - Split Position')
		{
		super(Object(first, .getSplitter(handle?, splitter), second))
		.moveObservers = Object()
		.loadSplit()
		}

	SetSplitSaveName(.splitSaveName, suffix = '')
		{
		if suffix isnt ''
			.splitSaveNameSuffix = suffix
		return .loadSplit()
		}

	splitSaveName: false
	loadSplit()
		{
		if .Destroyed?()
			return false
		if .splitSaveName is false
			{
			.Act('SetDefaultSplit')
			return false
			}
		prevSplit = UserSettings.Get(.splitSaveName $ .splitSaveNameSuffix, false)
		if prevSplit isnt false
			{
			.SetSplit(prevSplit)
			return true
			}
		return false
		}

	getSplitter(handle?, splitter)
		{
		if splitter isnt false
			return splitter
		return Object(handle? ? "HandleSplitter" : "Splitter")
		}

	moveObservers: ()
	AddMoveObserver(fn)
		{
		.moveObservers.Add(fn)
		}

	callMoveObservers()
		{
		for fn in .moveObservers
			fn()
		}

	n: #(0.5, 0.5)
	SetSplit(.n)
		{
		total = .n[0] + .n[1]
		.n = .n.Map({ it / total })
		.Act('SetSplit', .n)
		}

	GetSplit()
		{
		return .n
		}

	splitChanged: false
	UpdateSplit(.n)
		{
		.splitChanged = true
		.callMoveObservers()
		}

	SaveSplit()
		{
		if .splitSaveName isnt false and .splitChanged
			UserSettings.Put(.splitSaveName $ .splitSaveNameSuffix, .GetSplit())
		}

	UpdateSplitter(remove = false)
		{
		.Remove(1)
		.Insert(1, remove ? #Vert : 'Splitter')
		}

	MaximizeSecond()
		{
		.Act(#MaximizeSecond)
		}

	Destroy()
		{
		.SaveSplit()
		super.Destroy()
		}
	}
