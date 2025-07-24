// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	jpg_compress_options: (
		'-sampling-factor 4:2:0',
		'-resize 2200x2200>',
		'-quality 25',
		'-interlace JPEG',
		'-colorspace RGB'
		)

	png_compress_options: (
		'-define png:compression-filter=5',
		'-define png:compression-level=9'
		)

	bmp_compress_options: (
		'-depth 8',
		'-compress RLE',
		'-colors 256'
		)

	tif_compress_options: (
		'-depth 16',
		'-resize 2200x2200>',
		'-compress lzw'
		)

	BuildCompressionCmd(type, options)
		{
		cmd = Object('convert', 'src.' $ type)
		for option in options
			cmd.Add(@option.Split(' '))
		cmd.Add('result.' $ type)
		return cmd
		}

	GetOption(type, extraOptions = #(colorSpace: 'rgb'))
		{
		if false is options = .getOption(type)
			return false
		return .processJPEGOptions(type, options, extraOptions)
		}

	getOption(type)
		{
		switch (type)
			{
		case 'jpg', 'jpeg' :
			return .jpg_compress_options
		case 'png' :
			return .png_compress_options
		case 'bmp' :
			return .bmp_compress_options
		case 'tif', 'tiff' :
			return .tif_compress_options
		default :
			return false
			}
		}

	normalRGB: '-colorspace RGB'
	forceRGB: '-set colorspace sRGB'
	processJPEGOptions(type, options, extraOptions)
		{
		if not #(jpg, jpeg).Has?(type) or extraOptions.colorSpace.Lower().Has?('rgb')
			return options

		return options.Copy().Remove(.normalRGB).Add(.forceRGB)
		}

	// Convert returns true if the file was successfully compressed, otherwise false
	Convert(file, outFileName)
		{
		if false is cmd = .buildCommand(file, outFileName, #())
			return false
		return 0 is RunPiped(cmd).ExitValue()
		}

	// Compress returns true if the file was successfully compressed, otherwise false
	Compress(file, outFileName, extraOptions = #(colorSpace: 'rgb'))
		{
		type = file.AfterLast('.').Lower()
		if false is options = .GetOption(type, extraOptions)
			return false
		if false is cmd = .buildCommand(file, outFileName, options)
			return false
		exit_code = RunPipedOutput.WithExitValue(cmd)
		if 0 is exit_code.exitValue
			return true

		// for certain errors, the output of the command could be helpful but for
		// now we are just logging return codes. RunPipedOutput would return the output
		//
		// known return codes
		// 1: corrupted file
		// -1073741515: missing dependency (could be different depending on the missing file)
		// supress errors where user changed extension. Might need to suppress more errors that we can't fix
		if not exit_code.output.Prefix?('convert: Not a')
			SuneidoLog('WARNING: Compressing file ' $ file $ ' failed.', calls:,
				params: exit_code)
		return false
		}

	buildCommand(file, outFile = false, options = false, forceOptions = false)
		{
		cmd = .getApp()
		if cmd is false
			return false

		if forceOptions isnt false
			cmd $= ' ' $ forceOptions.Join(' ')
		else
			cmd $= 'convert "' $ file $ '" ' $ options.Join(' ')
		if outFile isnt false
			cmd $= ' "' $ outFile $ '"'
		return cmd
		}

	// for jpeg
	//	1 = Horizontal (normal) - TopLeft
	//	2 = Mirror horizontal - TopRight
	//	3 = Rotate 180 - BottomRight
	//	4 = Mirror vertical - BottomLeft
	//	5 = Mirror horizontal and rotate 270 CW - LeftTop
	//	6 = Rotate 90 CW - RightTop
	//	7 = Mirror horizontal and rotate 90 CW - RightBottom
	//	8 = Rotate 270 CW - LeftBottom
	orientations: #(TopRight, BottomRight, BottomLeft,
		LeftTop, RightTop, RightBottom, LeftBottom)
	Orientation(file)
		{
		forceOptions = Object('identify', '-format', '"%[orientation]"',
			'"' $ file $ '"')
		if false is cmd = .buildCommand(file, :forceOptions)
			return false

		exit_code = .runWithErrorChecking({ RunPipedOutput.WithExitValue(cmd) })
		output = exit_code.output
		return exit_code.exitValue is 0 and .orientations.Has?(output)
			? output
			: false
		}

	AutoRotate(file, outFile)
		{
		options = Object('-auto-orient')
		if false is cmd = .buildCommand(file, outFile, options)
			return false
		exit_code = RunPipedOutput.WithExitValue(cmd)
		return 0 is exit_code.exitValue
		}

	ConvertToPdf(file, outFile)
		{
		if outFile is false
			outFile = file.BeforeLast('.') $ '.pdf'

		Assert(outFile.Suffix?('.pdf'), msg: '.pdf extension required')

		pageSize = .letterSize(file)
		if false is cmd = .buildCommand(file, outFile, Object(' -page ' $ pageSize))
			return false

		exit_code = RunPipedOutput.WithExitValue(cmd)
		return 0 is exit_code.exitValue
		}

	ConvertToPdfCmd(pageSize)
		{
		return Object('convert', 'src.jpg', '-page', pageSize, 'result.pdf')
		}

	GetResolutionCmd()
		{
		return Object('identify', '-format', '"%x %y\\n"', 'src.jpg')
		}

	letterSize(file)
		{
		forceOptions = Object('identify', '-format', '"%x %y"', '"' $ file $ '"')
		if false is cmd = .buildCommand(file, :forceOptions)
			return '612x792'

		if 0 is exit_code = RunPipedOutput.WithExitValue(cmd)
			return '612x792'

		.LetterSize(exit_code.output)
		}

	LetterSize(output)
		{
		dims = output.Split(' ')
		if dims.Size() isnt 2 or not dims.Every?({ it.Number?() })
			return '612x792'

		xDpi = Number(dims[0])
		yDpi = Number(dims[1])

		resize = String((8.5 * xDpi).Round(0)) $ 'x' $ /*= letter size width*/
			String( (11 * yDpi).Round(0)) /*= letter size height*/
		return resize
		}

	/* https://nono.ma/image-width-height-imagemagick-identify */
	GetWidthHeight(file)
		{
		result = ''
		size = Object(w: 0, h: 0)
		cmd = .getApp()
		if cmd is false
			return size

		cmd $= 'identify -format "%wx%h" '
		result = .runWithErrorChecking({ .runCommand(cmd, file) })
		if Object?(result) and result.exitValue is 0
			{
			output = result.output.Split('\r\n').Filter({ not it.Has?('identify:') })
			if output.Empty?()
				return size
			output = output[0].Split('x')
			size.w = Number(output[0])
			size.h = Number(output[1])
			}
		return size
		}
	getApp()
		{
		// We use an older version of ImageMagick on Linux which doesn't use magick.exe.
		// On Linux, the command is just "convert" which gets added in the buildCommand.
		// When we update the version being used on Linux, this will likely change.
		if not Sys.Windows?()
			return ''

		cmd = ExternalApp('imagemagick')
		return cmd is false ? false : '"' $ cmd $ '" '
		}

	GenerateThumbnail(file)
		{
		forceOptions = Object('convert', '"' $ file $ '"', '-thumbnail', '400x400','-')
		if false is cmd = .buildCommand(file, :forceOptions)
			return false
		res = .runWithErrorChecking({ RunPipedOutput.WithExitValue(cmd) })
		if res.exitValue isnt 0
			return false
		return res.output
		}

	runCommand(cmd, file)
		{
		result = ''
		if not Paths.IsValid?(file) // its a data stream
			result = RunPipedOutput.WithExitValue(cmd $ '-', file)
		else
			result = RunPipedOutput.WithExitValue(cmd $ '"' $ file $ '"')
		return result
		}

	runWithErrorChecking(block)
		{
		errPattern = '(?i)(Access is denied|The user name or password is incorrect|' $
			'Please try again later)'
		try
			{
			res = block()
			Suneido.Delete('ImageMagickAccessError')
			return res
			}
		catch (err)
			{
			if Sys.Client?() and (err =~ errPattern)
				{
				errRes = Object(exitValue: 1, output: err)
				if Suneido.Member?('ImageMagickAccessError') and
					Suneido.ImageMagickAccessError is true
					return errRes
				Suneido.ImageMagickAccessError = true
				msg = 'There was a problem accessing ImageMagick \r\n\r\n' $
					'Please contact your system administrator'
				Alert(msg, title: 'Access Denied', flags: MB.ICONERROR)
				SuneidoLog('ERROR: (CAUGHT) ' $ err, caughtMsg: 'Alerted user', call:)
				return errRes
				}
			throw err
			}
		}
	}
