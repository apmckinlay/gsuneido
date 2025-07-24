// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_readFromCache()
		{
		reader = new PdfReader()
		method = reader.PdfReader_readFromCache

		reader.PdfReader_cache = false
		Assert(method(0, 1) is: false)

		reader.PdfReader_cache = Object(segment: "1234567890", start: 10, end: 20)
		Assert(method(9, 15) is: false)
		Assert(method(20, 21) is: false)

		Assert(method(10, 15) is: "12345")
		Assert(method(15, 20) is: "67890")
		Assert(method(10, 20) is: "1234567890")
		Assert(method(10, false) is: "1234567890")

		Assert(method(15, 25) is: "67890")
		}

	Test_fileReadWithCache()
		{
		reader = new PdfReader()
		reader.BytesPerRead = 5
		method = reader.PdfReader_fileReadWithCache
		f = (FakeFile)("1234567890")

		Assert(method(f, 5, 1) is: false)
		Assert(method(f, 0, 3) is: '123')
		Assert(reader.PdfReader_cache is: #(segment: '12345', start: 0, end: 5))
		Assert(method(f, 0, 10) is: '12345')
		Assert(reader.PdfReader_cache is: #(segment: '12345', start: 0, end: 5))
		Assert(method(f, 5, 8) is: '678')
		Assert(reader.PdfReader_cache is: #(segment: '67890', start: 5, end: 10))
		Assert(method(f, 5, 10) is: '67890')
		Assert(method(f, 10, 15) is: false)
		Assert(reader.PdfReader_cache is: #(segment: '67890', start: 5, end: 10))

		reader.BytesPerRead = 20
		reader.PdfReader_cache = false
		Assert(method(f, 0, 15) is: '1234567890')
		Assert(reader.PdfReader_cache is: #(segment: '1234567890', start: 0, end: 10))

		reader = new PdfReader('1234567890')
		f = (FakeFile)("")
		method = reader.PdfReader_fileReadWithCache
		Assert(reader.PdfReader_cache is: #(segment: '1234567890', start: 0, end: 10,
			allRead?:))
		Assert(method(f, 0, 15) is: '1234567890')
		Assert(method(f, 0, false) is: '1234567890')
		Assert(method(f, 10, false) is: false)
		Assert(method(f, 9, false) is: '0')
		Assert(method(f, 11, false) is: false)
		}

	Test_FileRead()
		{
		reader = new PdfReader()
		reader.BytesPerRead = 5
		method = reader.FileRead
		f = (FakeFile)("1234567890abcde")

		Assert(method(f, 0, 3) is: "123")
		Assert(reader.PdfReader_cache is: #(segment: '12345', start: 0, end: 5))
		Assert(method(f, 3, 13) is: "4567890abc")
		Assert(reader.PdfReader_cache is: #(segment: 'abcde', start: 10, end: 15))
		Assert(method(f, 13, 20) is: "de")
		Assert(reader.PdfReader_cache is: #(segment: 'abcde', start: 10, end: 15))
		Assert(method(f, 15, 20) is: false)
		Assert(reader.PdfReader_cache is: #(segment: 'abcde', start: 10, end: 15))

		Assert(method(f, 0, 10) is: "1234567890")
		}

	data: `%PDF-1.6
%encoded
4 0 obj
54589
endobj
7 0 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream
jfs;dajf;wehjq;rfn;ljnv;oasnf;oejwepru342j5;hre09nut394ugn59
endstream
endobj

xref
0 33
0000000000 65535 f
0000000015 00000 n
0000000062 00000 n
trailer <</Size 33/Root 1 0 R>>
startxref
55342
%%EOF`
	Test_find()
		{
		reader = new PdfReader()
		reader.BytesPerRead = 175 // So that the cache will cut the second 'endobj'
		method = reader.PdfReader_find
		f = (FakeFile)(.data)

		Assert(method(f, 'endobj[\x00\x09\x0A\x0C\x0D\x20]+', 0) is: 36)
		Assert(reader.PdfReader_cache is: #(segment: '%PDF-1.6
%encoded
4 0 obj
54589
endobj
7 0 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream
jfs;dajf;wehjq;rfn;ljnv;oasnf;oejwepru342j5;hre09nut394ugn59
endstream
end', start: 0, end: 175))

		Assert(method(f, 'endobj[\x00\x09\x0A\x0C\x0D\x20]+', 42) is: 172)
		Assert(reader.PdfReader_cache is: #(segment: 'obj

xref
0 33
0000000000 65535 f
0000000015 00000 n
0000000062 00000 n
trailer <</Size 33/Root 1 0 R>>
startxref
55342
%%EOF', start: 175, end: 310))

		Assert(method(f, 'endobj[\x00\x09\x0A\x0C\x0D\x20]+', 178) is: false)
		}

	Test_fetch()
		{
		reader = new PdfReader()
		reader.BytesPerRead = 20
		find = reader.PdfReader_find
		method = reader.PdfReader_fetch
		f = (FakeFile)(.data)
		pattern = Object(ob: Object(), start: '[\x00\x09\x0A\x0C\x0D\x20][0-9]+ 0 obj',
			end: 'endobj[\x00\x09\x0A\x0C\x0D\x20]+', size: 6, includeStream: false)

		Assert(method(f, find(f, pattern.start, 0), pattern) is: 42)
		Assert(pattern.ob[0] is: Object(head: '\n4 0 obj
54589
endobj', tail: ''))

		Assert({ method(f, 178, pattern) }
			throws: 'Failed to read pdf: closing tag not found')

		reader.PdfReader_maxRead = 2
		pattern.ob = Object()
		pos = method(f, find(f, pattern.start, 0), pattern)
		Assert(pattern.ob isSize: 1)
		Assert({ method(f, find(f, pattern.start, pos), pattern) }
			throws: 'pdf too complex')
		}

	Test_Read()
		{
		reader = new PdfReader()
		reader.BytesPerRead = 20
		method = reader.Read
		f = (FakeFile)(.data)

		objs = Object()
		trailers = Object()
		method(f, objs, trailers)

		expectedObjs = Object(
			Object(
				head: '\n4 0 obj
54589
endobj',
				tail: ''),
			Object(
				head: '\n7 0 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream',
				tail: 'endstream
endobj',
				streamStart: 97,
				streamEnd: 161))
		expectedTrailers = Object(
			Object(
				head: '\ntrailer <</Size 33/Root 1 0 R>>\r',
				tail: ''))
		Assert(objs is: expectedObjs)
		Assert(trailers is: expectedTrailers)

		// handle non-zero generation number
		data = .data.Replace('7 0 obj', '7 333 obj').Replace('4 0 obj', '4 4 obj')
		f = (FakeFile)(data)

		objs = Object()
		trailers = Object()
		method(f, objs, trailers)

		expectedObjs = Object(
			Object(
				head: '\n4 4 obj
54589
endobj',
				tail: ''),
			Object(
				head: '\n7 333 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream',
				tail: 'endstream
endobj',
				streamStart: 99,
				streamEnd: 163))
		expectedTrailers = Object(
			Object(
				head: '\ntrailer <</Size 33/Root 1 0 R>>\r',
				tail: ''))
		Assert(objs is: expectedObjs)
		Assert(trailers is: expectedTrailers)
		}

	dataNestedTrailer: `%PDF-1.6
%encoded
4 0 obj
54589
endobj
7 0 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream
jfs;dajf;wehjq;rfn;ljnv;oasnf;oejwepru342j5;hre09nut394ugn59
endstream
endobj

xref
0 33
0000000000 65535 f
0000000015 00000 n
0000000062 00000 n
trailer
<<
/DecodeParms <<
/Columns 4
/Predictor 12
>>
/ID [ <5E63E25BA63B334E8BDF7EB95B4BFC0B> <072AEC0EFED43E4FA9A51B6665AA7A6C> ]
/Info 14 0 R
/Root 16 0 R
/Size 96
/Encrypt 95 0 R
>>
startxref
55342
%%EOF`
	Test_Read_withNestedTrailer()
		{

		reader = new PdfReader()
		reader.BytesPerRead = 20
		method = reader.Read
		f = (FakeFile)(.dataNestedTrailer)

		objs = Object()
		trailers = Object()
		method(f, objs, trailers)

		expectedObjs = Object(
			Object(
				head: '\n4 0 obj
54589
endobj',
				tail: ''),
			Object(
				head: '\n7 0 obj
<</Length 8 0 R/Filter/FlateDecode>>
stream',
				tail: 'endstream
endobj',
				streamStart: 97,
				streamEnd: 161))
		expectedTrailers = Object(
			Object(
				head: '\ntrailer
<<
/DecodeParms <<
/Columns 4
/Predictor 12
>>
/ID [ <5E63E25BA63B334E8BDF7EB95B4BFC0B> <072AEC0EFED43E4FA9A51B6665AA7A6C> ]
/Info 14 0 R
/Root 16 0 R
/Size 96
/Encrypt 95 0 R
>>\r',
				tail: ''))
		Assert(objs is: expectedObjs)
		Assert(trailers is: expectedTrailers)
		}
	}
