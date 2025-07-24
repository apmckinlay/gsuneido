// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (indexTable, lib, libRecordName, libRec = false)
	{
	try Database("destroy " $ indexTable)
	Database("ensure " $ indexTable $ " (name, howtos) key(name)")

	if libRec is false
		{
		libRec = Query1(lib $
			" where name is " $ Display(libRecordName) $ " and group is -1")
		if libRec is false
			{
			SuneidoLog("ERROR: Couldn't find library record to create " $
				"Book HowToIndex from - INDEX TABLE NOT CREATED")
			return
			}
		ob = libRec.text.SafeEval()
		}
	else
		ob = libRec.SafeEval()

	if not Object?(ob)
		{
		SuneidoLog("ERROR: Couldn't load object from library record to create " $
			"Book HowToIndex - INDEX TABLE NOT CREATED")
		return
		}
	for rec in ob
		QueryOutput(indexTable, rec)
	}