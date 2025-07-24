// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	try Database("destroy testmanyfields")
	fieldSize = 500
	s = "create testmanyfields ("
	for i in .. fieldSize
		s $= " field" $ i
	s $= ") key(field0)"
	start = Date()
	Database(s)
	Print("create " $ Date().MinusSeconds(start))

	start = Date()
	Transaction(update:)
		{|t|
		q = t.Query("testmanyfields")
		rec = Record()
		for (i = 1; i < fieldSize; i += 2)
			rec["field" $ i] = i
		for (i = 0; i < fieldSize; ++i)
			{
			rec.field0 = i
			q.Output(rec)
			}
		}
	Print("output " $ Date().MinusSeconds(start))

	start = Date()
	QueryApply("testmanyfields")
		{|unused|
		}
	Print("input " $ Date().MinusSeconds(start))
	}