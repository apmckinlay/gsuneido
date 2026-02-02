// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		r, w = Pipe()
		Thread({ w.Write("hello world"); w.Close() })
		Assert(r.Read(99) is: "hello world")
		Assert(r.Read(99) is: false)
		}
	Test_CopyTo()
		{
		file = .MakeFile()
		s = "now is the time\n" $
			"for all good men\n" $
			"to come to the aid\n" $
			"of their party"

		r, w = Pipe()
		Thread({ w.Write(s); w.Close() })
		File(file, "w") {|f| r.CopyTo(f); f = 0 }
		Assert(GetFile(file) is: s)

		r, w = Pipe()
		Thread()
			{
			File(file) {|f| f.CopyTo(w); f = 0 }
			w.Close()
			}
		t = ""
		while false isnt chunk = r.Read(999)
			t $= chunk
		Assert(t is: s)
		}
	}