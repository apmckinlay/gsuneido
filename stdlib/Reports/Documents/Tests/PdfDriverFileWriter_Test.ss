// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	testWriter: PdfDriverFileWriter
		{
		PdfDriverFileWriter_file(filename/*unused*/)
			{
			return FakeFile('')
			}
		}
	Test_one()
		{
		writer = (.testWriter)('test')

		writer.AddObject('<<Object 0>>')
		writer.AddObject('<<Object 1>>')
		writer.Reserve()
		Assert(writer.NextObjId is: 3)

		writer.AddObject('<<Object 3>>')
		writer.AddObject('<<Object 4>>')
		writer.AddObject('<<Object 2>>', id: 2)
		Assert(writer.NextObjId is: 5)
		Assert(writer.Locations is: #(0, 12, 72, 32, 52))
		Assert(writer.TotalLength is: 92)

		writer.Flush()
		writer.Write('END')
		Assert(writer.NextObjId is: 5)
		Assert(writer.Locations is: #(0, 12, 72, 32, 52))
		Assert(writer.TotalLength is: 95)

		writer.Flush()
		Assert(writer.PdfDriverFileWriter_f.Get()
			is: "<<Object 0>>1 0 obj <<Object 1>>3 0 obj <<Object 3>>" $
				"4 0 obj <<Object 4>>2 0 obj <<Object 2>>END")
		}
	}