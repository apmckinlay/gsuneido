// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
// more test cases in PdfMerger_special_Test
PdfTest
	{
	Test_cleanUpBody()
		{
		//Test that it removes our normal format properly
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\r"
		expected = "\n2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] " $
			"/Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj"
		.AssertCleanedUpBody(body, expected, 0)

		//Test that it works with carriage returns
		body = "%PDF-1.3\n\r%encoded\r1 0 obj <</Type /Catalog /Pages 2 0 R>>\rendobj\r" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj\r"
		expected = "\r2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] " $
			"/Count 3>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.AssertCleanedUpBody(body, expected, 0)

		//Test that order of catalog and pages doessn't matter
		body = "%PDF-1.3\n\r%encoded\r" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"1 0 obj <</Type /Catalog /Pages 2 0 R>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj "
		expected = "\r2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] " $
			"/Count 3>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.AssertCleanedUpBody(body, expected, 1)

		//Test that the object numbers of catalog and pages doesn't matter
		body = "%PDF-1.3\n\r%encoded\r" $
			"21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\rendobj\r" $
			"13 0 obj <</Type /Catalog /Pages 21 0 R>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj "
		expected = "\r21 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] " $
			"/Count 3>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.AssertCleanedUpBody(body, expected, 1)

		// Test object number 0 is handled
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0  R>>\nendobj\n" $
			"0 0 obj\n<</Size 58/Root 33 0 R/Info " $
				"1 0 R/ID[<7058466DE1806F2DAC7A64228AFDFF39>" $
				"<7058466DE1806F2DAC7A64228AFDFF39>]>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] /Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\r"
		expected = "\n2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0 R ] " $
			"/Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj"
		.AssertCleanedUpBody(body, expected, 0)
		}

	Test_offsetObjNums()
		{
		//Test that it changes objects and references
		//Test that missing object numbers don't matter
		numObj1 = 5
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [18 0 R 22 0 R 31 0  R ] /Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj "
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [23 0 R 27 0 R 36 0  R ] /Count 3>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj"

		.assertOffsetObjNums(body, expected, numObj1)

		//Test that it changes both the object and reference
		numObj1 = 5
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj "
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)

		body = "%PDF-1.3\n%encoded\n01 0 obj <</Type /Catalog /Pages 02 0 R>>\nendobj\n" $
			"02 0 obj <</Type /Pages /Kids [01 0 R 02 0 R 03 0 R ] /Count 3>>\nendobj\n" $
			"03 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj "
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)

		//Test that it works with carriage returns
		body = "%PDF-1.3\r%encoded\r1 0 obj <</Type /Catalog /Pages 2 0 R>>\rendobj\r" $
			"2 0 obj <</Type /Pages /Kids [1 0 R 2 0 R 3 0 R ] /Count 3>>\rendobj\r" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj "
		expected = "\r6 0 obj <</Type /Catalog /Pages 7 0 R>>\rendobj\r" $
			"7 0 obj <</Type /Pages /Kids [6 0 R 7 0 R 8 0 R ] /Count 3>>\rendobj\r" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj"
		.assertOffsetObjNums(body, expected, numObj1)

		//Test that order of objects and whitespace in front of references doesn't matter
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n"
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)
		}

	assertOffsetObjNums(body, expected, numObj1)
		{
		method = PdfMerger.PdfMerger_offsetObjNums
		pdfOb = .ReadPdf(body)
		method(pdfOb, numObj1)
		Assert(.PdfObToString(pdfOb, body) is: expected)
		}

	Test_offsetObjNums2()
		{
		numObj1 = 5
		//Test to ensure RG is not changed by this code
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
			"19 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(Scanned by \r\nTEST) Tj ET\nendstream\nendobj\n"
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
			"24 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(Scanned by \r\nTEST) Tj ET\nendstream\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)

		//Testing stream/endstream inside of text stream
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
			"19 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(stream\nScanned by \r\nTEST\nendstream) Tj ET\nendstream\nendobj\n"
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
			"24 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(stream\nScanned by \r\nTEST\nendstream) Tj ET\nendstream\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)

		//Testing stream/endstream inside of text stream
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Pages 2 0 R/Type /Catalog>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
			"19 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(stream\nScanned by \r\nTEST\nendstream) Tj ET\nendstream\nendobj\n"
		expected = "\n6 0 obj <</Pages 7 0 R/Type /Catalog>>\nendobj\n" $
			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
			"24 0 obj\nstream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
			"(stream\nScanned by \r\nTEST\nendstream) Tj ET\nendstream\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)
		// 34690: PdfReader not handling "# # obj" in text stream
		//Testing "10 0 obj" inside of text stream
//		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
//			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
//			"/WinAnsiEncoding >>\nendobj\n" $
//			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
//			"19 0 obj\n<</Length 20 0 R>>\n" $
//			"stream\n0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
//			"(endstream stream\nScanned endobj 10 0 obj by \r\nTEST\n" $
//			"endstream) Tj ET\nendstream\nendobj\n"
//		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\n" $
//			"endobj\n" $
//			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
//			"/WinAnsiEncoding >>\nendobj\n" $
//			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
//			"24 0 obj\n<</Length 25 0 R>>\nstream\n" $
//			"0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
//			"(endstream stream\nScanned endobj 10 0 obj by \r\nTEST\n" $
//			"endstream) Tj ET\nendstream\nendobj"
//		.assertOffsetObjNums(body, expected, numObj1)

//		//Testing "10 0 obj" inside of text stream
//		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
//			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
//			"/WinAnsiEncoding >>\nendobj\n" $
//			"2 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n" $
//			"19 0 obj\n<</Length 20 0 R>>\nstream\n" $
//			"0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
//			"(endstream stream\n>>Scanned endobj 10 0 obj by \r\nTEST\n" $
//			"endstream) Tj ET\nendstream\nendobj\n"
//		expected = "%PDF-1.3\n%encoded\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\n" $
//			"endobj\n" $
//			"8 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
//			"/WinAnsiEncoding >>\nendobj\n" $
//			"7 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj\n" $
//			"24 0 obj\n<</Length 25 0 R>>\nstream\n" $
//			"0 0 0 RG\n0 Tr\nBT /F1 14.837 Tf 36 746.69 Td " $
//			"(endstream stream\n>>Scanned endobj 10 0 obj by \r\nTEST\n" $
//			"endstream) Tj ET\nendstream\nendobj\n"
//		.assertOffsetObjNums(body, expected, numObj1)
		}

	Test_whitespace()
		{
		numObj1 = 5
		// Newlines followed by Null characters before object numbers
		body = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"\x003 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"\x002 0 obj <</Type /Pages /Kids [ 1 0 R 2 0 R 3 0 R ] /Count 3>>\nendobj\n"
		expected = "\n6 0 obj <</Type /Catalog /Pages 7 0 R>>\nendobj" $
			"\x008 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj" $
			"\x007 0 obj <</Type /Pages /Kids [ 6 0 R 7 0 R 8 0 R ] /Count 3>>\nendobj"
		.assertOffsetObjNums(body, expected, numObj1)
		}

	body1: "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [4 0 R 6 0 R ] /Count 2>>\nendobj\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\nendobj" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"4 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"5 0 obj\n<</Length 6 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"6 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"7 0 obj\n<</Length 7 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n%comments\ntrailer <</Size 9/Root 1 0 R>> "

	body2: "%PDF-1.3\n%encoded\n4 0 obj <</Type /Catalog /Pages 3 0 R>>\nendobj\n" $
			"3 0 obj <</Type /Pages /Kids [2 0 R 6 0 R ] /Count 2>>\nendobj\n" $
			"1 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"2 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\n/Font <<\n/F1 1 0 R>>>>>>\nendobj\n" $
			"5 0 obj\n<</Length 6 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"6 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\n/Font <<\n/F1 1 0 R>>>>>>\nendobj\n" $
			"7 0 obj\n<</Length 7 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"trailer <</Size 9/Root 4 0 R>> "


	Test_appendBody()
		{
		//Tests when the documents have different order of catalog and pages
		//and both have multiple pages
		expected = "\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\n" $
			"endobj\n" $
			"2 0 obj <</Type /Pages /Kids [4 0 R 6 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\nendobj" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"4 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"5 0 obj\n<</Length 6 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"6 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"7 0 obj\n<</Length 7 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj" $
			"\n11 0 obj <</Parent 2 0 R/Type /Pages /Kids [10 0 R 14 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"9 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"10 0 obj<</Type /Page /Parent 11 0 R /MediaBox [0 0 612 792] " $
			"/Contents [13 0 R] /Resources<<\n/Font <<\n/F1 9 0 R>>>>>>\nendobj\n" $
			"13 0 obj\n<</Length 14 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"14 0 obj<</Type /Page /Parent 11 0 R /MediaBox [0 0 612 792] " $
			"/Contents [15 0 R] /Resources<<\n/Font <<\n/F1 9 0 R>>>>>>\nendobj\n" $
			"15 0 obj\n<</Length 15 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"16 0 obj << /Producer (Suneido PDF Generator) >>\nendobj"

		//NOTE:Test can easily be broken by changing the format of how things get replaced
		//For example adding or removing trailing or leading white space or newlines
		merger = .PdfMerger()
		merger.PdfMerger_initializeBodyVariables(.ReadPdf(.body1, merger))
		merger.PdfMerger_appendBody(.ReadPdf(.body2, merger))

		s = .PdfObToString(merger.PdfMerger_mergedOb[0], .body1)
		s $= .PdfObToString(merger.PdfMerger_mergedOb[1], .body2)
		Assert(s like: expected)
		}

	body2R: "%PDF-1.3\r%encoded\r4 0 obj <</Type /Catalog /Pages 3 0 R>>\rendobj\r" $
			"3 0 obj <</Type /Pages /Kids [2 0 R 6 0 R ] /Count 2>>\rendobj\r" $
			"1 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\rendobj\r" $
			"2 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\r/Font <<\r/F1 1 0 R>>>>>>\rendobj\r" $
			"5 0 obj\r<</Length 6 0 R /Filter /FlateDecode>>\rstream\rcompressed\r" $
			"endstream\rendobj\r" $
			"6 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\r/Font <<\r/F1 1 0 R>>>>>>\rendobj\r" $
			"7 0 obj\r<</Length 7 0 R /Filter /FlateDecode>>\rstream\rcompressed\r" $
			"endstream\rendobj\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"trailer <</Size 9/Root 4 0 R>> "
	body2RN: "%PDF-1.3\r\n%encoded\r\n4 0 obj <</Type /Catalog /Pages 3 0 R>>\r\n" $
			"endobj\r\n" $
			"3 0 obj <</Type /Pages /Kids [2 0 R 6 0 R ] /Count 2>>\r\nendobj\r\n" $
			"1 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\r\nendobj\r\n" $
			"2 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\r\n/Font <<\r\n/F1 1 0 R>>>>>>\r\n" $
			"endobj\r\n" $
			"5 0 obj\r\n<</Length 6 0 R /Filter /FlateDecode>>\r\nstream\r\n" $
			"compressed\r\n" $
			"endstream\r\nendobj\r\n" $
			"6 0 obj<</Type /Page /Parent 3 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\r\n/Font <<\r\n/F1 1 0 R>>>>>>\r\n" $
			"endobj\r\n" $
			"7 0 obj\r\n<</Length 7 0 R /Filter /FlateDecode>>\r\nstream\r\n" $
			"compressed\r\n" $
			"endstream\r\nendobj\r\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\r\nendobj\r\n" $
			"trailer <</Size 9/Root 4 0 R>>\r\n"

	Test_calcLocations()
		{
		merger = .PdfMerger()
		method = merger.PdfMerger_calcLocations

		//Test out of order objects and when there is a reference to the object
		//before the object
		expected = Object(0, 125, 221, 63, 16, 343, 425, 547, 629)
		Assert(method(Object(.ReadPdf(.body2, merger))) is: expected)

		//Test with carriage returns
		expected = Object(0, 125, 221, 63, 16, 343, 425, 547, 629)
		Assert(method(Object(.ReadPdf(.body2R, merger))) is: expected)

		//Test with carriage returns and newlines
		//The merged pdf will only keep one close to opening tag
		expected = Object(0, 127, 224, 64, 16, 349, 436, 561, 648)
		Assert(method(Object(.ReadPdf(.body2RN, merger))) is: expected)
		}

	Test_buildXRef()
		{
		//Test when objects are not ordered
		expected = "xref\n0 9\n" $
			"0000000000 65535 f \n0000000125 00000 n \n" $
			"0000000221 00000 n \n0000000063 00000 n \n0000000016 00000 n \n" $
			"0000000343 00000 n \n0000000425 00000 n \n0000000547 00000 n \n" $
			"0000000629 00000 n \n" $
			"trailer <</Size 9/Root 4 0 R>>\nstartxref\n684\n%%EOF"
		.assertBuiltXRef(.body2, expected)

		//Test with carriage returns
		expected = "xref\n0 9\n" $
			"0000000000 65535 f \n0000000125 00000 n \n" $
			"0000000221 00000 n \n0000000063 00000 n \n0000000016 00000 n \n" $
			"0000000343 00000 n \n0000000425 00000 n \n0000000547 00000 n \n" $
			"0000000629 00000 n \n" $
			"trailer <</Size 9/Root 4 0 R>>\nstartxref\n684\n%%EOF"
		.assertBuiltXRef(.body2R, expected)

		//Test with carriage returns and newlines
		expected = "xref\n0 9\n" $
			"0000000000 65535 f \n0000000127 00000 n \n" $
			"0000000224 00000 n \n0000000064 00000 n \n0000000016 00000 n \n" $
			"0000000349 00000 n \n0000000436 00000 n \n0000000561 00000 n \n" $
			"0000000648 00000 n \n" $
			"trailer <</Size 9/Root 4 0 R>>\nstartxref\n704\n%%EOF"
		.assertBuiltXRef(.body2RN, expected)
		}

	assertBuiltXRef(body, expected)
		{
		merger = .PdfMerger()
		merger.PdfMerger_initializeBodyVariables(.ReadPdf(body, merger))
		Assert(merger.PdfMerger_buildXRef() like: expected)
		}

	Test_filterFiles()
		{
		method = PdfMerger.FilterFiles
		files = #()
		expected = #()
		Assert(method(files) is: expected)
		files = #(`C:\Users\name\Documents\new.pdf`,
			`C:\Users\name\Documents\new.jpg`, `new.jpg`)
		Assert(method(files) is: files)
		files = #(`C:\Users\name\Documents\new.PdF`,
			`C:\Users\name\Documents\new.jPg`, `new.JPG`)
		Assert(method(files) is: files)
		files = #(`C:\Users\name\Documents\new.PDF`,
			`C:\Users\name\Documents\new.jpg`, `new.JPG`, `\Pictures\old.png`)
		expected = #(`C:\Users\name\Documents\new.PDF`,
			`C:\Users\name\Documents\new.jpg`, `new.JPG`)
		Assert(method(files) is: expected)
		files = #(` C:\Users\name\Documents\data.csv  `,
			`C:\Users\name\Documents\new.PDF`, `C:\Users\name\Documents\new.jpg`,
			 `new.JPG`, `\Pictures\old.png`)
		expected = #(`C:\Users\name\Documents\new.PDF`,
			`C:\Users\name\Documents\new.jpg`, `new.JPG`)
		Assert(method(files) is: expected)
		}

	Test_convertObjectStreams()
		{
		method = PdfMerger.PdfMerger_convertStream
		// Data that ends with a tab
		data = "\n36 0 obj\r<</Filter/FlateDecode/First 5/Length 60/N 1/Type/ObjStm>>" $
			"stream" $
		stream = "\r\n" $
			Zlib.Compress("39 0 <</ProcSet[/PDF/ImageC]/XObject<</Im0 37 0 R>>>>")
		expected = "\n39 0 obj\n<</ProcSet[/PDF/ImageC]/XObject<</Im0 37 0 R" $
			">>>>\nendobj\n"
		Assert(method(data, stream, #()) is: Object(Object(head: expected, tail:'')))

		// object with long property text
		data = "\n36 0 obj\r<</Filter/FlateDecode/First 5/Length 60/N 1/Type/ObjStm>>" $
			"stream\r\n" $
			Zlib.Compress("39 0 <</ProcSet[/PDF/ImageC]/XObject<</Im0 37 0 R>>>>") $
			"\nendstream\nendobj"
		property = 'a'.Repeat(2000)
		lastObj = "\n40 0 obj<</Title " $ property $ ">> endobj"
		data $= lastObj $ '\n'
		expected = "\n39 0 obj\n<</ProcSet[/PDF/ImageC]/XObject<</Im0 37 0 R" $
			">>>>\nendobj\n" $ lastObj
		Assert(.PdfObToString(.ReadPdf(data), data) like: expected)

		// Data that ends with a \r
		data = "\n12 0 obj\n" $
			"<</Filter/FlateDecode/First 5/Length 57/N 1/Type/ObjStm>>stream\r\n" $
			Zlib.Compress("24 0 <</Count 3/Kids[28 0 R 1 0 R 4 0 R]/Type/Pages>>") $
			"\nendstream\nendobj\n"
		expected = "\n24 0 obj\n<</Count 3/Kids[28 0 R 1 0 R 4 0 R]/Type/Pages>>" $
			"\nendobj"
		Assert(.PdfObToString(.ReadPdf(data), data) like: expected)

		// Object number info has newline delimiters
		data = "\n2 0 obj\n<</Length 12 0 R/Filter/FlateDecode/Type/ObjStm/N 4" $
			"/First 20>>\nstream\n" $
			Zlib.Compress("1 0\n9 32\n6 71\n3 161\n" $
			"[/PDF/ImageB/ImageC/ImageI/Text]<</ProcSet 1 0 R/XObject<</I0 4 0 R>>>>" $
			"<</Type/Page/Parent 3 0 R/Contents 7 0 R/Resources 9 0 R" $
			"/MediaBox[0 0 609.8824 799.3469]>><</Type/Pages/Count 1/Kids[ 6 0 R]>>") $
			"\nendstream\nendobj\n"
		expected = "\n1 0 obj\n[/PDF/ImageB/ImageC/ImageI/Text]\nendobj\n" $
			"\n9 0 obj\n<</ProcSet 1 0 R/XObject<</I0 4 0 R>>>>\nendobj\n" $
			"\n6 0 obj\n<</Type/Page/Parent 3 0 R/Contents 7 0 R/Resources 9 0 R" $
			"/MediaBox[0 0 609.8824 799.3469]>>\nendobj\n" $
			"\n3 0 obj\n<</Type/Pages/Count 1/Kids[ 6 0 R]>>\nendobj"
		Assert(.PdfObToString(.ReadPdf(data), data) like: expected)

		data = "\n2 0 obj\n<</Length 12 0 R /Type/ObjStm/N 4" $
			"/First 20>>\nstream\n" $
			"1 0\n9 32\n6 71\n3 161\n" $
			"[/PDF/ImageB/ImageC/ImageI/Text]<</ProcSet 1 0 R/XObject<</I0 4 0 R>>>>" $
			"<</Type/Page/Parent 3 0 R/Contents 7 0 R/Resources 9 0 R" $
			"/MediaBox[0 0 609.8824 799.3469]>><</Type/Pages/Count 1/Kids[ 6 0 R]>>" $
			"\nendstream\nendobj\n"
		expected = "\n1 0 obj\n[/PDF/ImageB/ImageC/ImageI/Text]\nendobj\n" $
			"\n9 0 obj\n<</ProcSet 1 0 R/XObject<</I0 4 0 R>>>>\nendobj\n" $
			"\n6 0 obj\n<</Type/Page/Parent 3 0 R/Contents 7 0 R/Resources 9 0 R" $
			"/MediaBox[0 0 609.8824 799.3469]>>\nendobj\n" $
			"\n3 0 obj\n<</Type/Pages/Count 1/Kids[ 6 0 R]>>\nendobj"
		Assert(.PdfObToString(.ReadPdf(data), data) like: expected)
		}

	Test_removeXRefs()
		{
		// XRefs are removed automatically when read
		data = "\n3 0 obj\n<</Type/XRef>>stream\nencoded\nendstream\nendobj\nstartxref" $
			"\n216\n%%EOF\n92 0 obj\n<</Size 93>>stream\nencoded\nendstream\nendobj\n" $
			"startxref\n281557\n%%EOF"
		expected = "\n3 0 obj\n<</Type/XRef>>stream\nencoded\nendstream\nendobj" $
			"\n92 0 obj\n<</Size 93>>stream\nencoded\nendstream\nendobj"
		pdfOb = .ReadPdf(data)
		Assert(.PdfObToString(pdfOb, data) is: expected)

		// When there is a XRef table in the middle of the document
		data = "\n3 0 obj\n<</Type/XRef>>stream\nencoded\nendstream\nendobj\n" $
			"xref\n4 17\n0000000016 00000 n\n0000000790 00000 n\n0000000850 00000 n\n" $
			"trailer\n<</Size 21/Root 5 0 R/Info 3 0 R]/Prev 6814>>\nstartxref\n" $
			"0\n%%EOF\n92 0 obj\n<</Size 93>>stream\nencoded\nendstream\nendobj\n" $
			"startxref\n281557\n%%EOF"
		expected = "\n3 0 obj\n<</Type/XRef>>stream\nencoded\nendstream\nendobj" $
			"\n92 0 obj\n<</Size 93>>stream\nencoded\nendstream\nendobj"
		pdfOb = .ReadPdf(data)
		Assert(.PdfObToString(pdfOb, data) is: expected)
		}

	Test_convertBodyToSuneidoFormat()
		{
		data = "%PDF-1.3\n%encoded\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [4 0 R 6 0 R ] /Count 2>>\nendobj\n" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"4 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"5 0 obj\n<</Length 6 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"6 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"7 0 obj\n<</Length 7 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n%comments\ntrailer <</Size 8/Root 1 0 R>> "
		expected =
			"\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [] /Count 0>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"6 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"7 0 obj<</Type /Page /Parent 5 0 R /MediaBox [0 0 612 792] " $
			"/Contents [8 0 R] /Resources<<\n/Font <<\n/F1 6 0 R>>>>>>\nendobj\n" $
			"8 0 obj\n<</Length 9 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"9 0 obj<</Type /Page /Parent 5 0 R /MediaBox [0 0 612 792] " $
			"/Contents [10 0 R] /Resources<<\n/Font <<\n/F1 6 0 R>>>>>>\nendobj\n" $
			"10 0 obj\n<</Length 10 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj"

		merger = .PdfMerger()
		pdfOb = .ReadPdf(data)
		merger.PdfMerger_convertBodyToSuneidoFormat(pdfOb)

		s = .PdfObToString(merger.PdfMerger_mergedOb[0], '')
		s $= .PdfObToString(merger.PdfMerger_mergedOb[1], data)
		Assert(s like: expected)
		Assert(merger.PdfMerger_totalObj is: 10)
		Assert(merger.PdfMerger_kids is: ' 5 0 R')
		Assert(merger.PdfMerger_totalPages is: 2)
		Assert(merger.PdfMerger_parent is: "/Parent 2")
		}

	Test_updateParentRef()
		{
		merger = .PdfMerger('/Parent 10')
		body = "%PDF-1.3\n%\xe9\xe9\xe9\xe9\n" $
			"1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [] /Count 0>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n"
		expected = "\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"2 0 obj <</Parent 10 0 R/Type /Pages /Kids [] /Count 0>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj"
		.assertParentRefUpdates(body, merger, expected)

		//Test that the object number of Pages doesn't matter
		body = "%PDF-1.3\n%\xe9\xe9\xe9\xe9\n" $
			"1 0 obj <</Type /Catalog /Pages 45 0 R>>\nendobj\n" $
			"45 0 obj <</Type /Pages /Kids [5 0 R] /Count 1>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 45 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n"
		expected = "\n1 0 obj <</Type /Catalog /Pages 45 0 R>>\nendobj\n" $
			"45 0 obj <</Parent 10 0 R/Type /Pages /Kids [5 0 R] /Count 1>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 45 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj"
		.assertParentRefUpdates(body, merger, expected)

		//Test that order of the pages tables doesnt matter
		body = "%PDF-1.3\n%\xe9\xe9\xe9\xe9\n" $
			"1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"2 0 obj <</Type /Pages /Kids [5 0 R] /Count 1>>\nendobj\n"
		expected = "\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"2 0 obj <</Parent 10 0 R/Type /Pages /Kids [5 0 R] /Count 1>>\nendobj"
		.assertParentRefUpdates(body, merger, expected)
		}

	assertParentRefUpdates(body, merger, expected)
		{
		merger.PdfMerger_mergedOb = Object()
		pdfOb = .ReadPdf(body, merger)
		merger.PdfMerger_fetchPagesInfo(pdfOb)
		merger.PdfMerger_updateParentRef(pdfOb)
		Assert(.PdfObToString(pdfOb, body) like: expected)
		}

	Test_updateParentRef_overriding_page_count()
		{
		merger = .PdfMerger('/Parent 10')
		//Test the overriding page count
		body = "%PDF-1.3\n%\xe9\xe9\xe9\xe9\n" $
			"1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"2 0 obj <</Type /Pages /Kids [5 0 R] /Count 1>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"2 0 obj <</Type /Pages /Kids [5 0 R 14 0 R] /Count 3>>\nendobj\n"
		expected = "\n" $
			"1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"5 0 obj <</Parent 2 0 R/Type /Pages /Kids [7 0 R 9 0 R ] /Count 2>>\n" $
			"endobj\n" $
			"2 0 obj <</Type /Pages /Kids [5 0 R] /Count 1>>\nendobj\n" $
			"3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj\n" $
			"2 0 obj <</Parent 10 0 R/Type /Pages /Kids [5 0 R 14 0 R] /Count 3>>" $
			"\nendobj"
		.assertParentRefUpdates(body, merger, expected)
		}

	Test_updatePages()
		{
		bodyBeforeKids = "\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\n" $
			"endobj\n2 0 obj <</Type /Pages /Kids ["

		bodyAfterKids = "endobj\n" $
			"8 0 obj << /Producer (Suneido PDF Generator) >>\nendobj" $
			"3 0 obj <</Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding " $
			"/WinAnsiEncoding >>\nendobj\n" $
			"4 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [5 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"5 0 obj\n<</Length 6 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj\n" $
			"6 0 obj<</Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] " $
			"/Contents [7 0 R] /Resources<<\n/Font <<\n/F1 3 0 R>>>>>>\nendobj\n" $
			"7 0 obj\n<</Length 7 0 R /Filter /FlateDecode>>\nstream\ncompressed\n" $
			"endstream\nendobj"

		kids = "4 0 R 6 0 R 8 0 R 10 0 R"
		expected = bodyBeforeKids $ kids $ "] /Count 4>>\n" $ bodyAfterKids
		.assertUpdatesPages(kids, 4, expected)

		kids = Seq(1, 400).Map({ it $ ' 0 R' }).Join(' ')
		expected = bodyBeforeKids $ kids $ "] /Count 400>>\n" $ bodyAfterKids
		.assertUpdatesPages(kids, 400, expected)

		kids  = Seq(1, 420).Map({ it $ ' 0 R' }).Join(' ')
		expected = bodyBeforeKids $ kids $ "] /Count 420>>\n" $ bodyAfterKids
		.assertUpdatesPages(kids, 420, expected)
		}

	assertUpdatesPages(kids, totalPages, expected)
		{
		merger = .PdfMerger()
		merger.PdfMerger_initializeBodyVariables(.ReadPdf(.body1, merger))
		merger.PdfMerger_appendBody(.ReadPdf(.body1, merger))
		merger.PdfMerger_totalPages = totalPages
		merger.PdfMerger_kids = kids
		merger.PdfMerger_updatePages()

		Assert(.PdfObToString(merger.PdfMerger_mergedOb[0], .body1) is: expected)
		}

	Test_isJpg()
		{
		fn = PdfMerger.PdfMerger_isJpg?

		Assert(fn('not_a_jpg.pdf') is: false)
		Assert(fn('not_a.jpg.pdf') is: false)
		Assert(fn('not_a.jpeg.pdf') is: false)
		Assert(fn('not_a.gif') is: false)
		Assert(fn('is_a.jpg'))
		Assert(fn('is_a.jpeg'))
		Assert(fn('is_a.pdf.jpg'))
		Assert(fn('is_a.pdf.jpeg'))
		}

	Test_empty_file()
		{
		file1 = .TempTableName() $ '.pdf'
		file2 = .TempTableName() $ '.pdf'
		invalidFiles = PdfMerger([file1, file2], .TempTableName())
		Assert(invalidFiles[0] endsWith: '(invalid)')
		Assert(invalidFiles[1] endsWith: '(invalid)')
		}

	Test_compressibleImageObject?()
		{
		m = PdfMerger.PdfMerger_compressibleImageObject?
		Assert(m("") is: false)
		Assert(m(" 	") is: false)
		Assert(m("ImageTest More Stuff") is: false)
		Assert(m("/Subtype/Image") is: false)
		Assert(m("/Subtype/Image /Filter/DCTDecode"))
		Assert(m("/Test /Subtype /Image /Length 8987 /Filter/DCTDecode"))
		Assert(m("/Test /Subtype /Image /Length 8987 /Filter/FlateDecode") is: false)
		Assert(m("/Test /Subtype/Image /Filter/DCTDecode /Length 8987"))
		Assert(m("/Test
			/Subtype/Image
			/Length 8987
			/Filter
			/DCTDecode"))
		Assert(m("<</Test
			/Subtype/Image
			/Length 8987
			/Filter
			/DCTDecode
			>>"))
		Assert(m("/Test
			/Subtype
			/Image
			/Filter
			/DCTDecode
			/Length 8987"))
		Assert(m("/Subtype/Image /Filter[/FlateDecode/DCTDecode]/Length") is: false)
		Assert(m("/Subtype/Image /Filter/FlateDecode/Length 7987") is: false)
		Assert(m("/Subtype/Image /Filter [ /DCTDecode ]"))
		Assert(m("<< /Length 8 0 R /Type /XObject /Subtype /Image /Width 1654
/Height 2338 /Interpolate true /ColorSpace 9 0 R /Intent /Perceptual
/BitsPerComponent 8 /Filter /DCTDecode
>>
stream"))
		Assert(m("<< /Length 8 0 R /Type /XObject /Subtype /Image /Width 1654
/Height 2338 /Interpolate true /ColorSpace 9 0 R /Intent /Perceptual
/BitsPerComponent 8 /Filter /FlateDecode/DCTDecode/otherDecode/other2Decode
>>
stream") is: false)
		}

	Test_getStreamSize()
		{
		m = PdfMerger.PdfMerger_getStreamSize
		Assert(m(#()) is: 0)
		Assert(m(#(streamSize: 8767, otherStuff: 'stuff')) is: 8767)
		Assert(m(#(a:'a', streamStart: 111, b: 'b', streamEnd: 156)) is: 45)
		}

	Test_imageObjectColorSpace()
		{
		fun = PdfMerger.PdfMerger_imageObjectColorSpace
		header = `<< /Width 76
/Height 99
/ColorSpace /DeviceRGB
/BitsPerComponent 8
/Length 13 0 R
/Filter [/ASCII85Decode /DCTDecode]
>>`
		Assert(fun(header) is: 'devicergb')

		header = `<<
/Type /XObject
/Subtype /Image
/Filter /DCTDecode
/Width 2479
/Height 3229
/Length 838617
/BitsPerComponent 8
/ColorSpace /DeviceRGB
>>`
		Assert(fun(header) is: 'devicergb')

		header = `<< % Attributes dictionary
/Subtype /NChannel
/Process
<< /ColorSpace [/ICCBased CMYK_ICC profile ]
/Components [/Cyan /Magenta /Yellow /Black]
>>
/Colorants
<< /Spot1 [/Separation /Spot1 alternateSpace tintTransform2]
/Spot2 [/Separation /Spot2 alternateSpace tintTransform3]
>>
>>`
		Assert(fun(header) is: `/iccbased cmyk_icc profile `)

		header = `<< /Width 76
/Height 99
/ColorSpace << /Cs12 12 0 R >>
/Pattern << /P1 15 0 R >>
/BitsPerComponent 8
/Length 13 0 R
/Filter [/ASCII85Decode /DCTDecode]
>>`
		Assert(fun(header) is: ' /cs12 12 0 r ')

		header = `<< /Width 76
/Height 99
/Pattern << /P1 15 0 R >>
/BitsPerComponent 8
/Length 13 0 R
/Filter [/ASCII85Decode /DCTDecode]
/ColorSpace << /Cs12 12 0 R >>
>>`
		Assert(fun(header) is: ' /cs12 12 0 r ')

		header = `<< /Type /XObject
/Subtype /Image
/Width 288
/Height 288
/ColorSpace 10 0 R
/BitsPerComponent 8
/Length 105278
/Filter /ASCII85Decode
>>`
		Assert(fun(header) is: '10 0 r')

		header = `<< /Type /XObject
/Subtype /Image
/Width 288
/Height 288
/BitsPerComponent 8
/Length 105278
/Filter /ASCII85Decode
/ColorSpace 10 0 R
>>`
		Assert(fun(header) is: '10 0 r')
		}

	Test_getStreamLength()
		{
		fn = PdfMerger.PdfMerger_getStreamLength
		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length 152/` $
			`Matrix[0.24 0 0 -0.24 18 -223.56]/PaintType 2/PatternType 1/` $
			`Resources<</ProcSet[/PDF /ImageB]>>/TilingType 1/Type/` $
			`Pattern/XStep 16/YStep 16`
		Assert(fn(header, #()) is: 152)

		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length 152`
		Assert(fn(header, #()) is: 152)

		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length hello/Pattern/XStep 16`
		Assert({ fn(header, #()) } throws: 'invalid length object')

		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length 152 0 R`
		Assert(fn(header, Object(Object(head: '\n152 0 obj\n333\nendobj'))) is: 333)

		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length 152 0 R/Pattern/XStep 16`
		Assert(fn(header, Object(Object(head: '\n152 0 obj\n333\nendobj'))) is: 333)

		header = `/BBox[0 0 16 16]/Filter/FlateDecode/Length 152 0 R/Pattern/XStep 16`
		Assert({ fn(header, Object(Object(head: '\n1 0 obj\n333\nendobj'))) }
			throws: 'cannot find length object 152')
		}
	}
