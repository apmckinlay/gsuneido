// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (destinationFile, sourceFile = false)
	{
	DeleteFile(destinationFile $ '.9')
	for (i = 8; i > 0 ; --i)
		if FileExists?(destinationFile $ '.' $ i)
			MoveFile(destinationFile $ '.' $ i, destinationFile $ '.' $ (i + 1))
	if sourceFile isnt false and FileExists?(sourceFile)
		MoveFile(sourceFile, destinationFile $ '.1')
	}