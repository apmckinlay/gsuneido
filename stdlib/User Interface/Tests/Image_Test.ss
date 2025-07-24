// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_imageFromFile()
		{
		mock = Mock(Image)
		mock.When.imageFromFile([anyArgs:]).CallThrough()
		mock.When.cleanupAndThrow([anyArgs:]).CallThrough()
		mock.When.createFileStream([anyArgs:]).Return(-1, 0)
		mock.When.FileSize([anyArgs:]).Return(5.5.Mb())

		Assert({ mock.imageFromFile('notValidImage.pdf', [x: NULL]) }
			throws:  `Image: invalid file extension: "pdf"`)

		Assert({ mock.imageFromFile('validImage.jpg', [x: NULL]) }
			throws:  `Image: validImage.jpg: exceeds max image size: 5.5 mb`)

		// Fails due to createFileStream returning -1 not passing SUCCEED
		mock.When.FileSize([anyArgs:]).Return(1.Mb())
		Assert({ mock.imageFromFile('validImage.jpg', [x: 1]) }
			throws:  `Image: failed to create stream on validImage.jpg`)

		// Fails due to pstm.x being NULL
		Assert({ mock.imageFromFile('validImage.jpg', [x: NULL]) }
			throws:  `Image: failed to create stream on validImage.jpg`)

		// No errors, returns the "tuple" (an object)
		Assert(Object?(mock.imageFromFile('validImage.jpg', [x: 1])))
		}
	}