// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
function (contribution, window)
	{
	contribMenus = Object()
	// contrib requires: text, seq, cmd. Optional: condition, root
	for item in GetContributions(contribution).Sort!({ |x,y| x.seq < y.seq })
		{
		if item.Member?('condition') and not (item.condition)(:window)
			continue

		if not item.Member?('root')
			contribMenus.Add(Object(root: item.text, cmd: item.cmd))
		else if false is menu = contribMenus.FindOne({ it.root is item.root })
			{
			options = Object()
			options[item.text] = item.cmd
			contribMenus.Add(Object(root: item.root, :options))
			}
		else
			menu.options[item.text] = item.cmd
		}
	if contribMenus.NotEmpty?()
		window.AddWindowMenuOptions(contribMenus)
	}