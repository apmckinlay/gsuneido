// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (folder)
	{
	exists? = false
	res = CatchFileAccessErrors(folder,
		{ exists? = DirExists?(folder.RightTrim('\\/')) })
	return res isnt true ? res : exists?
	}