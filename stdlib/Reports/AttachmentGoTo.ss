// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
function (file, hwnd)
	{
	file = Paths.ToLocal(file)
	if not FileExists?(file)
		file = OpenImageSettings.Copyto() $ file
	ShellExecute(hwnd, NULL, file)
	}