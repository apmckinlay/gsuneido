// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Asup("hello <$ 123 + 456 $> world <$ Date.End().Format('d/M/yyyy')
			$> the end") is: "hello 579 world 1/1/3000 the end")
		}
	}