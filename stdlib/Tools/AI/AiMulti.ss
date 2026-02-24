// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// find symbols with 2 or 3 of code, test, help
function ()
	{
	triple = Object()
	test = Object()
	help = Object()
	QueryApply("stdlib where group = -1")
		{|x|
		t = not QueryEmpty?("stdlib where group = -1 and
			name =~ `^" $ x.name $ "_?Test$`")
		h = not QueryEmpty?("suneidoc", name: x.name)
		if t and h
			triple.Add(x.name)
		else if t
			test.Add(x.name)
		else if h
			help.Add(x.name)
		}
	return triple, test, help
	}