// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260211
Test
	{
	Test_one()
		{
		r, w = Pipe()
		s = "hello world"
		Thread(Bind(.writer, w, s))
		Assert(r.Read(99) is: s)
		Assert(r.Read(99) is: false)
		}
	writer(w, s)
		{
		w.Write(s)
		w.Close()
		}
	Test_CopyTo()
		{
		file = .MakeFile()
		s = "now is the time\n" $
			"for all good men\n" $
			"to come to the aid\n" $
			"of their party"

		r, w = Pipe()
		Thread(Bind(.writer, w, s))
		File(file, "w") {|f| r.CopyTo(f) }
		Assert(GetFile(file) is: s)

		r, w = Pipe()
		Thread(Bind(.copy, file, w))
		t = ""
		while false isnt chunk = r.Read(999)
			t $= chunk
		Assert(t is: s)
		}
	copy(file, w)
		{
		File(file) {|f| f.CopyTo(w) }
		w.Close()
		}
	}