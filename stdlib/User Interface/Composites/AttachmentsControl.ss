// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	ImageDropFileList(files, source)
		{
		i = .controls.Find(source)
		if i + files.Size() > .controls.Size()
			{
			.AlertInfo("Drag and Drop",
				"Too many files to drop here")
			return true
			}
		for file in files
			.controls[i++].Drop(file)
		return true
		}
	getter_controls()
		{ // once only
		.controls = Object()
		.find_controls(this)
		return .controls
		}
	find_controls(control)
		{
		for c in control.GetChildren()
			if c.Base?(OpenImageControl)
				.controls.Add(c)
			else
				.find_controls(c)
		}
	}