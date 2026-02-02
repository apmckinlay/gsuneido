// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
PdfTest
	{
	Test_getSuneidoFormatKids()
		{
		obj = "\n32 0 obj\n<</Count 4/Kids[37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\r<</Count 4/Kids[37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 4/Kids [37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 4/Kids[ 37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 4/Kids[37 0 R 1 0 R 3 0 R 5 0 R ]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 4/Kids [ 37 0 R 1 0 R 3 0 R 5 0 R ]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 4 /Kids[37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")
		obj = "\n32 0 obj\n<</Count 1/Kids [ 37 0 R ]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R")
		obj = "\n32 0 obj\n<</Count 1/Kids[37 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R")
		obj = "\n32 0 obj\n<</Count 1/Kids [ 3 0 R ]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "3 0 R")
		obj = "\n32 0 obj\n<</Count 1/Kids [ 370 0 R ]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "370 0 R")
		obj = "\n32 0 obj\n<</Count 4 /Kids[374 0 R 12 0 R 334 0 R 5212 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "374 0 R 12 0 R 334 0 R 5212 0 R")
		obj = "\n32 0 obj\n<</Count 4/Kids[        37 0 R 1 0 R 3 0 R 5 0 R]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, "37 0 R 1 0 R 3 0 R 5 0 R")

		kids = Seq(1, 400).Map({ it $ ' 0 R' }).Join(' ')
		obj = "\n32 0 obj\n<</Count 100 /Kids[" $ kids $ "]/Type/Pages>>"
		.checkGetSuneidoFormatKids(obj, kids)
		}

	checkGetSuneidoFormatKids(obj, kids)
		{
		method = PdfMerger.PdfMerger_getSuneidoFormatKids
		pdfOb = Object(
			pagesPos: 1,
			objs: Object(
				Object(head: "\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj"),
				Object(head: obj),
				Object(head: "\n32 0 obj\n<</Count 4/Kids[100 0 R]/Type/Pages>>")))
		Assert(method(pdfOb) is: kids)
		}

	Test_handleLinearized()
		{
		//Test the removal of linearized object and first page xref table
		data = "%PDF-1.3\r%encoded\r\n2 0 obj\r" $
			"<</Linearized 1/L 7023/O 6/E 3607/N 1/T 6824/H [ 636 154]>>\r" $
			"endobj\r" $
			"					   \r\nxref\r\n4 17\r\n" $
			"0000000016 00000 n\r\n0000000790 00000 n\r\n0000000850 00000 n\r\n" $
			"0000001113 00000 n\r\n0000001201 00000 n\r\n0000002574 00000 n\r\n" $
			"0000002667 00000 n\r\n0000002768 00000 n\r\n0000002865 00000 n\r\n" $
			"0000002952 00000 n\r\n0000003044 00000 n\r\n0000003143 00000 n\r\n" $
			"0000003238 00000 n\r\n0000003329 00000 n\r\n0000003419 00000 n\r\n" $
			"0000003515 00000 n\r\n0000000636 00000 n\r\ntrailer\r\n" $
			"<</Size 7/Root 1 0 R/Info 3 0 R/ID[<12E2838BE9810A4CA3E11985E9B8C396>" $
			"<12E2838BE9810A4CA3E11985E9B8C396>]/Prev 6814>>\r\nstartxref\r\n" $
			"0\r\n%%EOF\r\n				 \r\n" $
			"6 0 obj\r<</Filter/FlateDecode/I 89/Length 75/S 43>>stream\r\n" $
			"compressedData\r\nendstream\rendobj\r" $
			"1 0 obj\r<</Metadata 2 0 R/Pages 1 0 R/Type/Catalog>>\rendobj\r" $
			"3 0 obj\r<</Contents 6 0 R/CropBox[0 0 612 792]/MediaBox[0 0 612 792]" $
			"/Parent 1 0 R/Resources<</Font<</F1 4 0 R/F1B 5 0 R>>>>" $
			"/Rotate 0/Type/Page>>\rendobj\r4 0 obj\r" $
			"<</BaseFont/Helvetica/Encoding/WinAnsiEncoding/Subtype/Type1/Type/Font>>\r" $
			"endobj\r\r"
		expected =
			"\n6 0 obj\r<</Filter/FlateDecode/I 89/Length 75/S 43>>stream\r\n" $
			"compressedData\r\nendstream\rendobj\r" $
			"1 0 obj\r<</Metadata 2 0 R/Pages 1 0 R/Type/Catalog>>\rendobj\r" $
			"3 0 obj\r<</Contents 6 0 R/CropBox[0 0 612 792]/MediaBox[0 0 612 792]" $
			"/Parent 1 0 R/Resources<</Font<</F1 4 0 R/F1B 5 0 R>>>>" $
			"/Rotate 0/Type/Page>>\rendobj\r4 0 obj\r" $
			"<</BaseFont/Helvetica/Encoding/WinAnsiEncoding/Subtype/Type1/Type/Font>>\r" $
			"endobj"

		.assertLinearized(data, expected)
		}

	assertLinearized(data, expected)
		{
		pdfOb = .ReadPdf(data)
		Assert(pdfOb.linearized?)
		Assert(.PdfObToString(pdfOb, data) is: expected)
		}

	Test_getBodyLinearized_remove_xref()
		{
		data = '%PDF-1.2

3 0 obj
<<
/E 53850 /H [ 5251 273 ] /L 54082 /Linearized 1 /N 1 /O 6 /T 53972
>>
endobj

xref
3 215
0000000012 00000 n
trailer
<<
/ABCpdf 6115
/ID [ <49E1436FB4AC33954740E89D8BC0A046>
<8BEAD6382AE1BD6AA1D904781518EC9E> ]
/Length 0 /Prev 53962 /Root 4 0 R /Size 218 /Type /XRef
>>
startxref
0
%%EOF
4 0 obj
<<
/OpenAction [ 6 0 R /Fit ]
/Outlines 1 0 R /PageMode /UseNone
/Pages 2 0 R /Type /Catalog
>>
endobj\r'
		expected = '\n4 0 obj
<<
/OpenAction [ 6 0 R /Fit ]
/Outlines 1 0 R /PageMode /UseNone
/Pages 2 0 R /Type /Catalog
>>
endobj'
		.assertLinearized(data, expected)
		}

	Test_getBodyLinearized_remove_streamed_xref()
		{
		data = '%PDF-1.4

32 0 obj
<</Linearized 1/L 200346/O 34/E 48831/N 6/T 199991/H [ 447 152]>>
endobj

38 0 obj
<</DecodeParms<</Columns 3/Predictor 12>>/Filter
/FlateDecode/ID[<0767CFA463EC07408D9BB10FB637F7BD><0767CFA463EC07408D9BB10FB637F7BD>]
/Index[32 9]/Info 31 0 R/Length 44/Prev 199992/Root 33 0 R
/Size 41/Type/XRef/W[1 2 0]>>stream
hello stream
endstream
endobj
startxref
0
%%EOF

33 0 obj
<</Metadata 21 0 R/Pages 30 0 R/Type/Catalog>>
endobj '
		expected = '\n33 0 obj
<</Metadata 21 0 R/Pages 30 0 R/Type/Catalog>>
endobj'
		.assertLinearized(data, expected)
		}

	Test_runWithCatch()
		{
		mock = Mock(PdfMerger)
		mock.When.runWithCatch([anyArgs:]).CallThrough()
		mock.InvalidFiles = Object()

		Assert(mock.runWithCatch('test.pdf', { }) is: false)
		Assert(mock.InvalidFiles is: #())

		mock.InvalidFiles = Object()
		Assert(mock.runWithCatch('test.pdf', { throw 'Secured pdf' }) is: false)
		Assert(mock.InvalidFiles is: #('test.pdf (SECURED)'))

		mock.InvalidFiles = Object()
		Assert(mock.runWithCatch('test.pdf', { throw 'Zlib Uncompress error' }) is: false)
		Assert(mock.InvalidFiles is: #('test.pdf (invalid)'))

		mock.InvalidFiles = Object()
		Assert(mock.runWithCatch('test.pdf', { throw PdfMerger.LimitError }))
		Assert(mock.InvalidFiles is: #('test.pdf (last file attempted, too many files)'))

		mock.InvalidFiles = Object()
		Assert(mock.runWithCatch('test.pdf', { throw PdfReader.LimitError }))
		Assert(mock.InvalidFiles is: #('test.pdf (last file attempted, pdf too complex)'))
		}

	Test_isValidPdf()
		{
		reader = new PdfReader()
		f = FakeFile('%PDF-1.3 valid pdf')
		Assert(PdfMerger.PdfMerger_isValidPdf(reader, f))

		reader = new PdfReader()
		f = FakeFile('badFile %PLF-')
		Assert(PdfMerger.PdfMerger_isValidPdf(reader, f) is: false)

		reader = new PdfReader()
		f = FakeFile('				%PDF-1.3 has space but valid file')
		Assert(PdfMerger.PdfMerger_isValidPdf(reader, f))
		}

	Test_processObjects_SecuredPDF_at_trailer()
		{
s = "%PDF-1.4
%\x25\xe2\xe3\xcf\xd3
1 0 obj
<</Pages 3 0 R/Type /Catalog>>
endobj
2 0 obj
<</CreationDate <DDC934CB75DFA7BA68FA7F6A78FF89A5EFA812C79A0481>/ModDate
endobj
3 0 obj
<</Count 3/Kids [4 0 R 39 0 R 60 0 R]/Type /Pages>>
endobj
4 0 obj
<</Contents [5 0 R 83 0 R 84 0 R]/MediaBox [0 0 612 792]/Parent 3 0 R
/Resources <</ColorSpace 34 0 R/ExtGState 35 0 R/Font 38 0 R/Pattern 36 0 R
/ProcSet [/PDF /ImageB /Text]/XObject 37 0 R>>/Type /Page>>
endobj
5 0 obj
<</BitsPerComponent 1/Decode [1 0]/DecodeParms <</Columns 39/K -1>>/Filter
/CCITTFaxDecode/Height 39/ImageMask true/Length 44/Subtype /Image/Width 39>>
stream
\x98\x4e\x6c\x86\x16\x4f\x63\xed\xef\xb5\x84\xeb\xf4\x1a\x3d\xb8\x48\x6d\xc3\xfd" $
"\xd\xa\x75\x71\x87\xbb\xdd\xa4\x62\x92\x4f\x97\x15\xd2\xdc\x24\x5e\x71\xa2\xc4\x" $
"eb\x7b\x39\xda\x47endstream
endobj
xref
0 90
0000000000 65535 f
0000000016 00000 n
0000000062 00000 n
0000000619 00000 n
0000000686 00000 n
0000524276 00000 n
trailer
<</Encrypt 89 0 R/ID [(\327\177>\354\007\251\230\024I\354\274\260\"\246\006') " $
"(\327\177>\354\007\251\230\024I\354\274\260\"\246\006')]/Info 2 0 R/Root 1 0 R/Size 90>>
startxref
524504
%%EOF"
		.PutFile(f = .TempTableName() $ '.pdf', s)
		.AddTeardown({ DeleteFile(f) })
		Assert({ PdfMerger.PdfMerger_getFileOb(f) }
			throws: "Secured pdf")
		}

	Test_processObjects_SecuredPDF_with_obj()
		{
s = "%PDF-1.4
%\x25\xe2\xe3\xcf\xd3
1 0 obj
<</Pages 3 0 R/Type /Catalog>>
endobj
2 0 obj
<</CreationDate <DDC934CB75DFA7BA68FA7F6A78FF89A5EFA812C79A0481>/ModDate
endobj
3 0 obj
<</Count 3/Kids [4 0 R 39 0 R 60 0 R]/Type /Pages>>
endobj
4 0 obj
<</DecodeParms<</Columns 4/Predictor 12>>/Encrypt 8 0 R/Filter
/FlateDecode/ID[<14149F015591FC3E4E3F5692E48DA149><8A3D520FF4D9834AA21E9B0972AF3AD6>]
/Info 6 0 R/Length 39/Root 9 0 R/Size 7/Type/XRef/W[1 3 0]>>
stream
\x68\xde\x62\x62\x20\x2\x26\x46\xc6\x67\xe7\x99\x18\x18\x78\xd3\x81\x4\x43\x2f\x90\x60" $
"\x3c\xc4\xc4\xf8\x9f\xa7\x16\xc4\x65\x4\x8\x30\x20\x54\xc3\x5\x14
endstream
endobj
5 0 obj
<</Filter/FlateDecode/First 4/Length 49/N 1/Type/ObjStm>>
stream
\x4\xbd\x9f\x6c\x1a\x13\xbf\xf1\xf1\x7f\x37\x55\x62\xd5\xba\xea\xb3\xd\xa\x30\xb4\x15" $
"\x50\x4c\x63\x5b\x73\xc\x70\x95\x68\xe8\x5d\x74\xed\xd7\x5c\xc0\x95\x8\xdc\x3a\x88" $
"\xca\x3\xb9\xdb\x54\x3f\x43
endstream
endobj
xref
0 90
0000000000 65535 f
0000000016 00000 n
0000000062 00000 n
0000000619 00000 n
0000000888 00000 n
0000000999 00000 n
trailer
startxref
524504
%%EOF"
		.PutFile(f2 = .TempTableName() $ '.pdf', s)
		.AddTeardown({ DeleteFile(f2) })
		Assert({ PdfMerger.PdfMerger_getFileOb(f2) }
			throws: "Secured pdf")
		}

	Test_pdfByMicrosoftPrint()
		{
		// Microsoft Print to PDF
		body = "%PDF-1.3\nencoded\n" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"13 0 obj <</Type /Catalog /Pages 2 0 R>>\rendobj\r\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj "
		expected =
			"\n2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj" $
			"\n3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.AssertCleanedUpBody(body, expected, 1)

		body = "%PDF-1.3\n\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj\n" $
			"21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"13 0 obj <</Type /Catalog /Pages 21 0 R>>\rendobj\r"
		expected =  "\n3 0 obj <</Type /Font /Subtype /Type1 /BaseFont" $
			" /Helvetica /Encoding /WinAnsiEncoding >>\rendobj" $
			"\n21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj"
		.AssertCleanedUpBody(body, expected, 2)

		//Test that the catalog can have angle brackets in it
		body = "%PDF-1.3\n\r%encoded\r" $
			"21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"13 0 obj <</MarkInfo<</Marked true>>/Type/Catalog /Pages 21 0 R>>\rendobj" $
			"\r3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj "
		expected =
			"\r21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj" $
			"\r3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.AssertCleanedUpBody(body, expected, 1)

		//Test for XRef stream
		body = "%PDF-1.3\n\r%encoded\r" $
			"21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"13 0 obj <</MarkInfo<</Marked true>>/Type/Catalog /Pages 21 0 R>>\rendobj" $
			"\r3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj\r" $
			"5 0 obj\r<</DecodeParms<</Columns 5/Predictor 12>>/Filter/FlateDecode" $
			"/ID[<><>]/Info 33 0 R/Length 50/Root 35 0 R/Size 34/Type/XRef/W[1 3 1]>>" $
			"stream\rencoded\rendstream\rendobj\r"
		expected =
			"\r21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj" $
			"\r3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj\r" $
			"5 0 obj\r<</DecodeParms<</Columns 5/Predictor 12>>/Filter/FlateDecode" $
			"/ID[<><>]/Info 33 0 R/Length 50/Root 35 0 R/Size 34/Type/XRef/W[1 3 1]>>" $
			"stream\rencoded\rendstream\rendobj"
		.AssertCleanedUpBody(body, expected, 1)
		}

	Test_isJpg?()
		{
		method = PdfMerger.PdfMerger_isJpg?
		Assert(method(`C:\Users\name\Documents\new.jpg`))
		Assert(method(`C:\Users\name\Documents\new.pdf`) is: false)
		Assert(method(`new.jpg`))
		Assert(method(`new.JPG`))
		Assert(method(`new.pdf`) is: false)
		}

	Test_updateNumObj()
		{
		pdfOb = Object(numObj: 5, objs: Object(
			Object(head: '5 0 obj << >>')
			Object(head: '85 0 obj << >>')
			Object(head: '87 0 obj << >>')
			Object(head: '88 0 obj << >>')
			Object(head: '89 0 obj << >>')
		))
		method = PdfMerger.PdfMerger_updateNumObj
		data = #("/Size 10", "/Size 90", "/Size 25")
		data.Each() { method(pdfOb, it)	}
		Assert(pdfOb.numObj is: 89)

		data = #("/Size 10", "/Size 89", "/Size 25")
		pdfOb.numObj = 0
		data.Each() { method(pdfOb, it)	}
		Assert(pdfOb.numObj is: 89)

		data = #("/Size 6")
		pdfOb.numObj = 0
		data.Each() { method(pdfOb, it)	}
		Assert(pdfOb.numObj is: 89)

		pdfOb = Object(numObj: 0, objs: Object(
			Object(head: '1 0 obj << >>')
			Object(head: '2 0 obj << >>')
			Object(head: '3 0 obj << >>')
			Object(head: '4 0 obj << >>')
			Object(head: '5 0 obj << >>')
		))
		data = #("/Size 6")
		pdfOb.numObj = 0
		data.Each() { method(pdfOb, it)	}
		Assert(pdfOb.numObj is: 5)

		pdfOb.numObj = 10
		data.Each() { method(pdfOb, it)	}
		Assert(pdfOb.numObj is: 10)
		}

	Test_imageCompressionLimits()
		{
		merger = .PdfMerger()
		merger.PdfMerger_compress = true
		merger.PdfMerger_maxCompressedFileSize = 5
		merger.PdfMerger_totalImageSize = 10
		merger.PdfMerger_curFile = 'test.pdf'
		data = '%PDF-1.2

1 0 obj
<<
/Test
/Subtype/Image
/Length 8987
/Filter
/DCTDecode
>>
endobj
\r'
		Assert(.ReadPdf(data, merger) is: false)
		Assert(merger.InvalidFiles is: #('test.pdf (compressed file size over maximum)'))
		merger.InvalidFiles.Delete(all:)

		merger.PdfMerger_totalImageSize = 4
		Assert(.ReadPdf(data, merger) is: false)
		Assert(merger.InvalidFiles is: #())

		merger.CompressImages? = true
		Assert(.ReadPdf(data, merger) isnt: false)
		Assert(merger.InvalidFiles is: #())
		}

	Test_checkCompression()
		{
		mock = Mock(PdfMerger)
		mock.When.checkCompression([anyArgs:]).CallThrough()
		mock.InvalidFiles = Object()
		mock.PdfMerger_compress = mock.PdfMerger_maxCompressedFileSize =
			mock.PdfMerger_compressed? = false
		mock.PdfMerger_totalImageSizeReduction = 0
		mock.When.fileSize('emptyImage.jpg').Return(0)
		mock.When.fileSize('smallImage.jpg').Return(500)
		mock.When.fileSize('largeImage.jpg').Return(3.Mb())

		Assert(mock.checkCompression('emptyImage.jpg'))
		Assert(mock.InvalidFiles isSize: 0)
		mock.Verify.Never().fileSize([anyArgs:])

		mock.PdfMerger_compress = true
		Assert(mock.checkCompression('emptyImage.jpg'))
		Assert(mock.InvalidFiles isSize: 0)
		mock.Verify.Never().fileSize([anyArgs:])

		mock.PdfMerger_maxCompressedFileSize = 2.Mb()
		Assert(mock.checkCompression('emptyImage.jpg') is: false)
		Assert(mock.InvalidFiles.PopFirst()
			is: 'emptyImage.jpg (file has nothing compressible)')
		mock.Verify.Never().fileSize([anyArgs:])

		mock.PdfMerger_compressed? = true
		Assert(mock.checkCompression('emptyImage.jpg'))
		Assert(mock.InvalidFiles isSize: 0)
		mock.Verify.fileSize('emptyImage.jpg')

		Assert(mock.checkCompression('smallImage.jpg'))
		Assert(mock.InvalidFiles isSize: 0)
		mock.Verify.fileSize('smallImage.jpg')

		Assert(mock.checkCompression('largeImage.jpg') is: false)
		Assert(mock.InvalidFiles.PopFirst()
			is: 'largeImage.jpg (compressed file size over maximum)')
		mock.Verify.fileSize('largeImage.jpg')

		mock.PdfMerger_totalImageSizeReduction = 1.5.Mb()
		Assert(mock.checkCompression('largeImage.jpg'))
		Assert(mock.InvalidFiles isSize: 0)
		mock.Verify.Times(2).fileSize('largeImage.jpg')
		}

	Test_securedPdf?()
		{
		fn = PdfMerger.PdfMerger_securedPdf?
		Assert(fn(objs = Object(), trailers = Object()) is: false)
		Assert(fn(objs, trailers) is: false)

		objs.Add([head: 'test'])
		trailers.Add([head: 'test'])
		Assert(fn(objs, trailers) is: false)

		objs.Add([head: 'other /Encrypt'])
		Assert(fn(objs, trailers))
		Assert(fn(objs, trailers))
		Assert(fn(objs.Delete(all:), trailers) is: false)

		trailers.Add([head: 'other /Encrypt'])
		Assert(fn(objs, trailers))
		Assert(fn(objs, trailers))
		Assert(fn(objs, trailers.Delete(all:)) is: false)
		}

	Test_compressLimits()
		{
		mock = Mock(PdfMerger)
		mock.When.compressLimit?([anyArgs:]).CallThrough()
		mock.PdfMerger_compress = false
		mock.PdfMerger_maxCompressedFileSize = false
		mock.PdfMerger_totalImageSize = 250

		Assert(mock.compressLimit?() is: false)
		mock.Verify.Never().overCompressionLimit?([anyArgs:])

		mock.PdfMerger_maxCompressedFileSize = 200
		Assert(mock.compressLimit?() is: false)
		mock.Verify.Never().overCompressionLimit?([anyArgs:])

		mock.PdfMerger_compress = true
		Assert(mock.compressLimit?())
		mock.Verify.overCompressionLimit?(250)

		mock.PdfMerger_totalImageSize = 200
		Assert(mock.compressLimit?() is: false)
		mock.Verify.overCompressionLimit?(200)

		mock.PdfMerger_maxCompressedFileSize = 150
		Assert(mock.compressLimit?())
		mock.Verify.Times(2).overCompressionLimit?(200)

		mock.PdfMerger_compress = false
		Assert(mock.compressLimit?() is: false)
		mock.Verify.Times(2).overCompressionLimit?(200)
		}

	Test_failedToCompressImage?()
		{
		mock = Mock(PdfMerger)
		mock.When.failedToCompressImage?([anyArgs:]).CallThrough()
		mock.When.compressImage([anyArgs:]).Return(true, false)
		mock.PdfMerger_compress = false
		mock.PdfMerger_maxCompressedFileSize = false

		s = ''
		obj = reader = f = Object() // not used for test cases
		Assert(mock.failedToCompressImage?(obj, s, reader, f) is: false)
		mock.Verify.Never().compressibleImageObject?(s)

		mock.PdfMerger_compress = true
		Assert(mock.failedToCompressImage?(obj, s, reader, f) is: false)
		mock.Verify.compressibleImageObject?(s)
		mock.Verify.Never().compressImage(obj, s, reader, f)

		s = '/subtype/image,/Filter [/DCTDecode]'
		Assert(mock.failedToCompressImage?(obj, s, reader, f) is: false)
		mock.Verify.compressibleImageObject?(s)
		mock.Verify.compressImage(obj, reader, f, s)

		// compressImage returns false, indicating a failure
		Assert(mock.failedToCompressImage?(obj, s, reader, f) is: false)
		mock.Verify.Times(2).compressibleImageObject?(s)
		mock.Verify.Times(2).compressImage(obj, reader, f, s)

		mock.PdfMerger_maxCompressedFileSize = 100
		Assert(mock.failedToCompressImage?(obj, s, reader, f))
		mock.Verify.Times(3).compressibleImageObject?(s)
		mock.Verify.Times(3).compressImage(obj, reader, f, s)
		}

	Test_InvalidFilesMsg()
		{
		m = PdfMerger.InvalidFilesMsg

		Assert(m(#()) is: '')
		Assert(m(#('invalid.pdf'))
			is: 'Unable to append the following attachments to PDF:\n\n' $
				'invalid.pdf' $
				'\n\nPlease check if they are corrupted or secured.')

		Assert(m(#('invalid1.pdf', 'invalid2.pdf', 'invalid3.pdf'))
			is: 'Unable to append the following attachments to PDF:\n\n' $
				'invalid1.pdf\n' $
				'invalid2.pdf\n' $
				'invalid3.pdf' $
				'\n\nPlease check if they are corrupted or secured.')

		Assert(m(#('invalid1.pdf', 'invalid2.pdf',
			'valid.pdf (last file attempted, too many pdfs)'))
			is: 'Unable to merge files into one PDF:\r\n\r\n' $
				'invalid1.pdf\r\n' $
				'invalid2.pdf\r\n' $
				'valid.pdf (last file attempted, too many pdfs)\r\n\r\n' $
				'Please review the above list and adjust accordingly.')

		Assert(m(#('complex.pdf (last file attempted, pdf too complex)'))
			is: 'Unable to merge files into one PDF:\r\n\r\n' $
				'complex.pdf (last file attempted, pdf too complex)\r\n\r\n' $
				'Please review the above list and adjust accordingly.')
		}
	}