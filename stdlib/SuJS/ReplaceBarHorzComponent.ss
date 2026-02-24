// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
HorzComponent
	{
	New(@args)
		{
		super(@args)

		findBarHorz = .Window.FindControl('FindBar').FindControl('HorzEqualHeight')
		findBarChildren = findBarHorz.GetChildren()
		replaceBarChildren = .GetChildren()

		replaceStatic = replaceBarChildren[0]
		replaceStatic.Orig_xmin = findBarChildren[0].Xmin + // EnhancedButton
			findBarChildren[1].Xmin + // Skip
			findBarChildren[2].Xmin   // FindStatic
		replaceStatic.Recalc()

		replaceBarChildren[2].Xstretch = 1

		afterFind = 0
		for (i = 5; i < findBarChildren.Size(); i++)
			afterFind += findBarChildren[i].Xmin

		afterReplace = 0
		for (i = 3; i < replaceBarChildren.Size(); i++)
			afterReplace += replaceBarChildren[i].Xmin

		lastSkip = replaceBarChildren.Last()
		lastSkip.Xmin += afterFind - afterReplace
		lastSkip.SetMinSize()
		.Recalc()
		}
	}