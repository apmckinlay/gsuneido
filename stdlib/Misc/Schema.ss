// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(table)
		{
		Transaction(read:)
			{ |t|
			if t.QueryEmpty?("tables", :table)
				throw "Schema: non-existent table: " $ table
			s = table $ "\r\n    ("
			s $= .columnsDisplay(t, table)
			s $= .indexesDisplay(t, table)
			}
		return s
		}

	maxLineSize: 40
	columnsDisplay(t, table)
		{
		display = ""
		line = ""
		t.QueryApply("columns where table = " $ Display(table) $ " sort column")
			{ |c|
			if line.Size() > .maxLineSize
				{
				display $= line $ "\r\n"
				line = "        "
				}
			line $= (c.field isnt -1 or c.column.Suffix?('_lower!')
				? c.column : c.column.Capitalize()) $ ", "
			}
		display $= line[.. -2] $ ")\r\n"
		return display
		}

	indexesDisplay(t, table)
		{
		display = ""
		t.QueryApply("indexes", :table)
			{|x|
			display $= "    " $ (x.key is true ? "key " :
				(x.key is false ? "index " : "index unique ")) $
				"(" $ x.columns $ ")"

			if x.fktable > ""
				{
				display $= " in " $ x.fktable
				if x.fkcolumns isnt x.columns
					display $= "(" $ x.fkcolumns $ ")"
				if x.fkmode isnt 0
					display $= " cascade"
				if x.fkmode is 1
					display $= " update"
				}
			display $= "\r\n"
			}
		return display
		}
	}