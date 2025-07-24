// Copyright (C) 2000 Suneido Software Corp.
Test
	{
	Test_main()
		{
		Cursor("(" $ .MakeTable('(a, b, b_deps) key(a)') $ " rename b to c)
			union
			(" $ .MakeTable('(a) key(a)') $ " extend c = 1, c_deps = '')
			 /* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */").Close()
		}
	}