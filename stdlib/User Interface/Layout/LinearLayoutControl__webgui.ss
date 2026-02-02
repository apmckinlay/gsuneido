// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Group
	{
	Name:          	"LinearLayout"
	ComponentName:	"LinearLayout"

	ctrls:         	false /* List of all controls, visible or not */
	ctrlSizes:		false /* real control sizes synced by the browser code*/
	sizes:			false /* restored sizes */
	New(@args)
		{
		super(.init(args))
		// Initialize
		.ctrls = .GetChildren()
		.ComponentArgs.dir = .Dir
		}

	init(args)
		{
		.Dir = args.GetDefault("dir", "vert")
		return args
		}

	SyncControlSizes(.ctrlSizes, fromUser? = false)
		{
		if fromUser? isnt true
			return

		.ctrls.Each()
			{
			if it.Method?(#Loaded?) and it.Loaded?() is false
				// this is needed to trigger load when users expand
				// the bottom pane in Dispatching - Full Screen
				it.Resize(0, 0, 0, 0)
			}
		}

	SetControlSizes(.sizes)
		{
		.Act('SetControlSizes', .convertToIndexed(.sizes))
		}

	convertToIndexed(map)
		{
		converted = Object()
		for m in map.Members()
			{
			i = .mapIndex(m)
			converted[i] = map[m]
			}
		return converted
		}

	GetControlSize(indexOrName)
		{
		if .ctrlSizes isnt false
			return .ctrlSizes[.mapIndex(indexOrName)]
		if .sizes isnt  false
			return .sizes[indexOrName]
		return 0
		}

	GetControlPos(indexOrName)
		{
		if .ctrlSizes is false
			return 0
		i = .mapIndex(indexOrName)
		pos = 0
		for (j = 0; j < i and j < .ctrlSizes.Size(); j++)
			pos += .ctrlSizes[j]
		return pos
		}

	mapIndex(indexOrName)
		{
		if (Number?(indexOrName) and indexOrName.Int?() and
			0 <= indexOrName and indexOrName < .Tally(all:))
			return indexOrName
		else if String?(indexOrName)
			return .findCtrlIndex(indexOrName)
		else
			throw "no such control: " $ Display(indexOrName)
		}

	findCtrlIndex(name)
		{
		n = .Tally(all:)
		found = false
		for (i = 0; i < n; ++i)
			{
			if name is .ctrls[i].Name
				if false isnt found
					throw "duplicate control name: " $ Display(name)
				else
					found = i
			}
		return found
		}

	SetControlVisibilities?(@unused)
		{
		return true
		}

	SetControlVisibilities(visMap)
		{
		.Act('SetControlVisibilities', .convertToIndexed(visMap))
		}
	}