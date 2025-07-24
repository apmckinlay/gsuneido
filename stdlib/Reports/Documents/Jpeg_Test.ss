// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getImageInfo()
		{
		fn = Jpeg.Jpeg_getImageInfo
		Assert({ fn('') } throws: 'invalid jpg file format')

		Assert({ fn('\xff\xd8\00\00\00') } throws: 'cannot find SOI')

		Assert(fn('\xff\xd8' $
			'\xff\xc0\x00\x08\x00\x00\x01\x01\x00\x00\x00')
			is: #(index: 2, width: 256, height: 1))

		Assert(fn('\xff\xd8' $
			'\xff\xc0\x00\x09\x00\x00\x01\x01\x00\x00\x00' $
			'\xff\xc1\x00\x09\x00\x00\x02\x02\x00\x00\xff')
			is: #(index: 13, width: 512, height: 2))

		Assert(fn('\xff\xd8' $
			'\xff\xc0\x00\x09\x00\x00\x10\x10\x00\x00\x00' $
			'\xff\xc1\x00\x09\x00\x00\x02\x02\x00\x00\xff')
			is: #(index: 2, width: 4096, height: 16))

		// with '\xff\xc0' in APP segment
		Assert(fn('\xff\xd8' $
			'\xff\xe0\x00\x04\xff\xc0' $ // APP
			'\xff\xc0\x00\x09\x00\x00\x10\x10\x00\x00\x00' $ // SOF
			'\xff\xc1\x00\x09\x00\x00\x02\x02\x00\x00\xff' $ // SOF
			'\xff\x00\xc0\x00\x01\x02\xff\x00\x32') // entropy-coded data
			is: #(index: 8, width: 4096, height: 16))
		}

	Test_validExtension()
		{
		fn = Jpeg.ValidExtension?
		Assert(fn('') is: false)
		Assert(fn('xyz abc.txt') is: false)
		Assert(fn('xyz abc.pdf') is: false)
		Assert(fn('xyz.test.abc.test') is: false)
		Assert(fn('xyz.test.abc.png') is: false)
		Assert(fn('xyz.test.abc.jpg'))
		Assert(fn('xyz.test.abc.jpeg'))
		}
	}