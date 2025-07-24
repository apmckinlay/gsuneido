// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (old_table, lang)
	{
	columns = "trlang_from,"
	QueryApply(old_table $ " project trlang_name")
		{|x|
		columns $= "trlang_" $ x.trlang_name $ ","
		}
	columns = columns[.. -1]
	Database("ensure translatelanguage (" $ columns $ ") key(trlang_from)")

	Transaction(update:)
		{|t|
		t.QueryApply(old_table $ " where trlang_name is " $ Display(lang))
			{|x|
			q = t.Query("translatelanguage where trlang_from is " $ Display(x.trlang_from))
			rec = q.Next()
			if (rec is false)
				{
				// output
				newrec = Object(trlang_from: x.trlang_from)
				newrec["trlang_" $ x.trlang_name] = x.trlang_to
				t.QueryOutput("translatelanguage", newrec)
				}
			else if (x.trlang_to isnt rec["trlang_" $ x.trlang_name])
				{
				// update
				if (rec["trlang_" $ x.trlang_name] isnt "")
					Print("Updating " $ rec.trlang_from $
						" from: " $ rec["trlang_" $ x.trlang_name] $
						", to: " $ x.trlang_to)
				rec["trlang_" $ x.trlang_name] = x.trlang_to
				rec.Update()
				}
			}
		}
	}
