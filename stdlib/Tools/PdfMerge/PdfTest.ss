// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	mergerCl: PdfMerger
		{
		PdfMerger_getFileOb(@unused) { return '' }
		PdfMerger_finish(@unused) { }
		PdfMerger_returnInputFile(@unused) { return }
		CompressImages?: false
		PdfMerger_compressImage(@unused) { return .CompressImages? }
		}
	PdfMerger(parent = false, totalObj = 0, mergedOb = #())
		{
		merger = new .mergerCl(#('test.pdf'), 'test')
		merger.PdfMerger_mergedOb = mergedOb.Copy()
		merger.PdfMerger_parent = parent
		merger.PdfMerger_totalObj = totalObj
		return merger
		}

	ReadPdf(data, merger = false)
		{
		reader = new PdfReader()
		reader.Read(f = FakeFile(data), objs = Object(), trailers = Object())
		return merger isnt false
			? merger.PdfMerger_processObjects(reader, f, objs, trailers)
			: PdfMerger.PdfMerger_processObjects(reader, f, objs, trailers)
		}

	PdfObToString(pdfOb, data)
		{
		s = ''
		for obj in pdfOb.objs
			{
			s $= obj.head
			s $= obj.Member?(#streamStart)
				? data[obj.streamStart .. obj.streamEnd]
				: ''
			s $= obj.tail
			}
		return s
		}

	AssertCleanedUpBody(body, expected, catalogIdx)
		{
		method = PdfMerger.PdfMerger_cleanUpBody
		pdfOb = .ReadPdf(body)
		PdfMerger.PdfMerger_fetchPagesInfo(pdfOb)
		Assert(pdfOb.catalog is: catalogIdx)
		method(pdfOb)
		Assert(.PdfObToString(pdfOb, body) is: expected)
		Assert(pdfOb hasntMember: #catalog)
		}
	}
