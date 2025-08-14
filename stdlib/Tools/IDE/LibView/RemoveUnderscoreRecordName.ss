// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (recordName, code, w = false)
	{
	if recordName is false
		return code
	type = LibRecordType(code)
	if type not in ('class', 'function')
		return code
	name = LibraryTags.RemoveTagFromName(recordName)
	pat = '\<(?q)_' $ name $ '(?-q)\>'
	if code.Size() > uname_pos = code.FindRx(pat)
		{
		if w isnt false
			w.Add(uname_pos)
			// CheckCode removes this error if it's incorrect
		code = code.Replace(pat, ' ' $ name) // replacement must be same length
		}
	return code
	}
