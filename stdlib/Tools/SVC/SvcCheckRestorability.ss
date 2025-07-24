// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function(svcTable, rec)
	{
	if not svcTable.SvcEnabled?()
		return false
	if false is deleted = svcTable.Get(svcTable.MakeName(rec), deleted:)
		return false
	rec.lib_before_text = deleted.lib_before_text
	rec.lib_before_path = deleted.lib_before_path
	rec.lib_committed = deleted.lib_committed
	if svcTable.Type isnt #lib
		svcTable.GetData(rec)
	svcTable.VerifyModified(rec)
	return true
	}
