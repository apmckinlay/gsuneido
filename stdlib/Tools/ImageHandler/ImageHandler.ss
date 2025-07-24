// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Available?()
		{
		return not Sys.Windows?() or false isnt ExternalApp('imagemagick')
		}

	Convert(file, outFile)
		{
		processor = .getImageProcessor()
		return processor.Convert(file, outFile)
		}

	Compress(file, outFile, extraOptions = #(colorSpace: 'rgb'))
		{
		processor = .getImageProcessor()
		return processor.Compress(file, outFile, :extraOptions)
		}

	Orientation(file)
		{
		processor = .getImageProcessor()
		return processor.Orientation(file)
		}

	AutoRotate(file, outFile)
		{
		processor = .getImageProcessor()
		return processor.AutoRotate(file, outFile)
		}

	ConvertToPdf(file, outFile = false)
		{
		processor = .getImageProcessor()
		return processor.ConvertToPdf(file, outFile)
		}

	GetWidthHeight(file)
		{
		processor = .getImageProcessor()
		return processor.GetWidthHeight(file)
		}

	GenerateThumbnail(file)
		{
		processor = .getImageProcessor()
		return processor.GenerateThumbnail(file)
		}

	getImageProcessor()
		{
		return ImageMagick
		}
	}