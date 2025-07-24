// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_validFileToOpen()
		{
		mockImage = Mock(ImageControl)
		mockImage.When.ImageControl_fileSize([anyArgs:]).Return(1048577,
			{throw 'FileSize: image.jpg does not exist'},
			{throw 'FileSize: The account is not authorized to log in from this station'},
			11)
		mockImage.When.ImageControl_validFileToOpen([anyArgs:]).CallThrough()
		Assert(mockImage.ImageControl_validFileToOpen('garbage') is: 'InvalidFile',
			msg: 'non image file')

		Assert(mockImage.ImageControl_validFileToOpen('image.jpg') is: 'InvalidFile',
			msg: 'over max size')

		Assert(mockImage.ImageControl_validFileToOpen('image.jpg')
			is: 'InvalidFile: Can not find file')

		Assert(mockImage.ImageControl_validFileToOpen('image.jpg')
			is: 'InvalidFile: Can not open file')

		Assert(mockImage.ImageControl_validFileToOpen('image.jpg'), msg: 'valid file')
		}
	}