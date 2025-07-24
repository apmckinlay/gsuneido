// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	// user might not have powershell in the PATH
	if '' is systemRoot = Getenv('SystemRoot')
		systemRoot = 'C:/Windows'
	// 64 bit Powershell executable
	return '"' $
		Paths.Combine(systemRoot, "system32/WindowsPowerShell/v1.0/powershell.exe") $ '"'
	}