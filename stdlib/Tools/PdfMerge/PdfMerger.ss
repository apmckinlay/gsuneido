// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
class
	{
	//The extensions for the types of files that can be merged
	extensions: #("pdf", "jpg", "jpeg")
	compress: false
	compressed?: false

	/* Creates a PDF with the given filename by merging together the given list of
	*  attachments which is an object.
	*  Params:
	*		saveFilename: The filename of the PDF that will be created from the merge
	*		files: An object that contains a list of files to be merged
	*		compress: if true, the merger will attempt to compress the input file.
	*		maxCompressedFileSizeInMb: only used if compress is true.
	*			If given a value PdfMerger won't rebuild the PDF if the compressed size
	*			would be over this specification;
	*			If kept false PdfMerger will compress and merge all files passed in
	*  Returns: An object that contains a list of files that couldn't be merged/compressed
	*/
	afterMergedAsync: false
	CallClass(files, saveFileName, compress = false, maxCompressedFileSizeInMb = false,
		filesData = false, afterMergedAsync = false)
		{
		if files.Size() is 0 or (files.Size() is 1 and files[0] is saveFileName)
			return Object()
		if compress and maxCompressedFileSizeInMb isnt false
			Assert(files.Size() is 1)
		merger = new this(files, saveFileName, compress, maxCompressedFileSizeInMb,
			filesData, afterMergedAsync)
		return merger.InvalidFiles
		}

	/* Creates a PdfMerger object that will write to the given file name. The Merge
	*  function is used to merge a set of files together into one file with the name given
	*  to the constructor. The Finish function outputs the merged data to the file
	*  Params:
	*		.newfile: The filename of the new file being created
	*/
	New(.files, .newfile, .compress = false, maxCompressedFileSizeInMb = false,
		.filesData = false, .afterMergedAsync = false)
		{
		.maxMerge = Objects.SizeLimit
		.maxCompressedFileSize = maxCompressedFileSizeInMb isnt false
			? maxCompressedFileSizeInMb.Mb()
			: false
		if .compress
			.curFile = files[0]
		.totalImageSizeReduction = 0
		.totalImageSize = 0
		.InvalidFiles = Object()
		.cleanupFiles = Object()
		if files.Size() is 1 and not compress
			{
			Assert(.afterMergedAsync is: false)
			// If we reach this point, it is basically a file rename
			.returnInputFile(files[0])
			return
			}

		.totalObj = 0
		.parent = false
		.mergedOb = Object()
		if .afterMergedAsync isnt false
			{
			.mergeAsync()
			return
			}

		.merge(files)

		.finish()
		.cleanUp()
		}

	returnInputFile(file)
		{
		if .isJpg?(file)
			.jpgToPdf(file, .newfile)
		else
			.copyFile(file, .newfile)
		}

	/* isJpg?
	*  Returns if the file is a jpg.
	*  Params:
	*		filename: A string that is the name of the file
	*  Returns: true if the file is a jpg, false otherwise
	*/
	isJpg?(filename)
		{
		name = filename.Lower()
		return name.Suffix?(".jpg") or name.Suffix?(".jpeg")
		}

	/* jpgToPdf
	*  Creates a letter sized temporary PDF file from a JPG that can be merged
	*  with other PDF files. If the dimensions of the image are larger than letter format
	*  the image will be scaled down to fit on the page maintaining its aspect ratio.
	*  If a dimension is less than letter format then the image will be centered on the
	*  coresponding axis. The temporary file created will have the prefix "pdf".
	*  Params:
	*		imagename: The filename of the image to be converted, assumed to be a jpg
	*  Returns: The full filename of the temporary file that contains the PDF
	*/
	jpgToPdf(imagename, filename = false)
		{
		if filename is false
			filename = GetAppTempFullFileName("pdf").BeforeFirst(".tmp") $ ".pdf"

		newImage = .fileClass(imagename)
		if ImageHandler.ConvertToPdf(newImage, filename)
			return filename
		.trackInvalidFile(imagename)
		return false
		}

	trackInvalidFile(file, reason = 'invalid')
		{
		.InvalidFiles.Add(file $ Opt(' (', reason, ')'))
		}

	copyFile(from, to)
		{
		.runWithCatch(to)
			{
			if true isnt CopyFile(from, to, false)
				throw "CopyFile failed"
			}
		}

	/* merge
	*  Responsible for initiating the merging process. Merges multiple PDF files
	*  into one file using an array of file names. Files are merged in the order
	*  of the array. The contents of the PDF files are not changed so any page numbering
	*  is not updated.
	*  Params:
	*  		files: An array of PDF file names that are to be merged into a single file
	*/
	merge(files)
		{
		processFile = .processFile
		files.Any?({|file| .runWithCatch(file, { processFile(file) }) })
		if .InvalidFiles.NotEmpty?()
			return
		.updatePages()
		.xRef = '\n' $ .buildXRef()
		}

	mergeAsync()
		{
		.fileIndex = 0
		.mergeOneAsync()
		}

	mergeOneAsync()
		{
		// use timeout to avoid stack overflow
		SuUI.GetCurrentWindow().SetTimeout({ .mergeOneFileAsync() }, 0)
		}

	mergeOneFileAsync()
		{
		try
			{
			filename = .files[.fileIndex]
			if .fileSize(filename) is 0
				{
				.trackInvalidFile(filename, 'empty file')
				return false
				}
			fileData = .getFileContent(filename)
			if .isJpg?(filename)
				{
				.jpgToPdfAsync(filename, fileData)
				return
				}
			.getFileObAsync(filename, fileData)
			}
		catch (err)
			{
			.processErrorAndContinue(err)
			}
		}

	processErrorAndContinue(err)
		{
		if .files.Member?(.fileIndex)
			.processError(.files[.fileIndex], err)
		else
			throw err
		.afterOneFileMerged()
		}

	jpgToPdfAsync(file, fileData)
		{
		cmd = ImageMagick.GetResolutionCmd()
		sourceBytes = SuUI.GetCurrentWindow().Uint8Array(fileData)
		imFiles = Object(Object(name: 'src.jpg', content: sourceBytes))
		SuUI.GetCurrentWindow().Magick(imFiles, cmd).Then(
			{|result1|
			if result1.exitCode is 0
				.convertJpgToPdfAsync(result1, imFiles, file)
			else
				.processErrorAndContinue(result1.stderr[0])
			})
		return false
		}

	convertJpgToPdfAsync(result1, imFiles, file)
		{
		output = result1.stdout[0].Tr('"')
		letterSize = ImageMagick.LetterSize(output)
		cmd = ImageMagick.ConvertToPdfCmd(letterSize)
		SuUI.GetCurrentWindow().Magick(imFiles, cmd).Then(
			{ |result2|
			if result2.exitCode is 0
				.afterJpgToPdf(result2.outputFiles[0], file)
			else
				.processErrorAndContinue(result2.stderr[0])
			})
		}

	afterJpgToPdf(output, file)
		{
		resultFile = SuUI.GetCurrentWindow().File(output, MimeTypes.pdf)
		resultFile.ArrayBuffer().Then({|arrayBuffer|
			sourceBytes = SuUI.GetCurrentWindow().Uint8Array(arrayBuffer)
			convertPDF = SuUI.GetCurrentWindow().ArrayToString(sourceBytes)
			.filesData[.fileIndex].fileData = convertPDF
			.getFileObAsync(file, convertPDF)
			})
		}

	getFileObAsync(file, fileData)
		{
		objs = Object()
		trailers = Object()
		reader = new PdfReader(:fileData)
		if fileData.Size() is 0
			{
			.processErrorAndContinue('empty file')
			return
			}
		.fileClass(file, 'r', :fileData)
			{|f|
			if not .isValidPdf(reader, f)
				{
				.processErrorAndContinue('Invalid pdf')
				return
				}

			reader.Read(f, objs, trailers)
			.processObjectsAsync(reader, f, objs, trailers)
			}
		}

	processObjectsAsync(reader, f, objs, trailers)
		{
		pdfOb = .initializePdfOb()
		.pdfObIndex = 0
		.processOneObjectAsync(pdfOb, reader, f, objs, trailers)
		}

	processOneObjectAsync(pdfOb, reader, f, objs, trailers)
		{
		obj = objs[.pdfObIndex]
		if .compressLimit?()
			{
			.imageSizeOverLimit = true
			.trackInvalidFile(.files[0], .maxCompressedError)
			.afterAllFiles()
			return
			}

		head = obj.head
		if 0 is Number(head.Extract('(^|[\[' $ .ws $ '])([0-9]+)', 2))
			{
			.afterEachObj(head, pdfOb, reader, f, obj, objs, trailers)
			return
			}

		.updateLinearized(head, pdfOb)
		.updateSuneidoFormat(head, pdfOb)

		if .compress is true and .compressibleImageObject?(head)
			.compressImageAsync(head, pdfOb, reader, f, obj, objs, trailers)
		else
			.afterEachObj(head, pdfOb, reader, f, obj, objs, trailers)
		}

	compressImageAsync(head, pdfOb, reader, f, obj, objs = #(), trailers = #())
		{
		type = 'jpg'
		imageData = reader.ExtractStream(obj)
		options = ImageMagick.GetOption(type,
			extraOptions: Object(colorSpace: .imageObjectColorSpace(head)))
		cmd = ImageMagick.BuildCompressionCmd(type, options)
		sourceBytes = SuUI.GetCurrentWindow().Uint8Array(imageData)
		imFiles = Object(Object(name: 'src.jpg', content: sourceBytes))
		afterImageCompressed = .afterImageCompressed
		afterEachObj = .afterEachObj
		SuUI.GetCurrentWindow().Magick(imFiles, cmd).Then(
			{|result|
			if result.exitCode is 0
				{
				output = result.outputFiles[0]
				resultFile = SuUI.GetCurrentWindow().File(output, "image/jpeg")
				resultFile.ArrayBuffer().Then({|arrayBuffer|
					sourceBytes = SuUI.GetCurrentWindow().Uint8Array(arrayBuffer)
					compressImg = SuUI.GetCurrentWindow().ArrayToString(sourceBytes)
					if false is afterImageCompressed(compressImg, imageData, head, obj)
						.processErrorAndContinue(.maxCompressedError)
					else
						afterEachObj(head, pdfOb, reader, f, obj, objs, trailers)
					})
				}
			else
				.processErrorAndContinue(.maxCompressedError)
			})
		}

	afterEachObj(head, pdfOb, reader, f, obj, objs = #(), trailers = #())
		{
		.afterObjCompressed(head, pdfOb, reader, f, obj, objs, trailers)
		.pdfObIndex++
		if .pdfObIndex isnt objs.Size()
			.processOneObjectAsync(pdfOb, reader, f, objs, trailers)
		else
			.afterAllObjectProcessed(pdfOb, trailers)
		}

	afterOneMerged(pdfOb)
		{
		if .afterMergedAsync is false
			return

		pdfOb.filename = .files[.fileIndex]
		if false is .checkCompression(.files[.fileIndex])
			{
			.afterAllFiles()
			return
			}

		try
			{
			.processFileAfterCompressed(pdfOb)
			}
		catch (err)
			{
			.trackInvalidFile(pdfOb.filename, err)
			.afterAllFiles()
			return
			}

		.afterOneFileMerged()
		}

	afterOneFileMerged()
		{
		.fileIndex++
		if .fileIndex isnt .files.Size()
			{
			.mergeOneAsync()
			return
			}
		.afterAllFiles()
		}

	afterAllFiles()
		{
		if .InvalidFiles.NotEmpty?()
			{
			(.afterMergedAsync)(invalidFiles: .InvalidFiles)
			return
			}
		.updatePages()
		.xRef = '\n' $ .buildXRef()
		.finish()
		.cleanUp()
		(.afterMergedAsync)(invalidFiles: .InvalidFiles)
		}

	runWithCatch(file, block)
		{
		try
			block()
		catch (err)
			return .processError(file, err)
		return false
		}

	LimitError: 'too many files'
	processError(file, err)
		{
		reason = 'invalid'
		if stop? = (err is .LimitError or err is PdfReader.LimitError)
			reason = 'last file attempted, ' $ err
		else if err is 'Secured pdf'
			reason = 'SECURED'
		else if not err.Prefix?('File:')
			SuneidoLog('ERRATIC: (CAUGHT) Unable to merge - ' $ err, params: [:file],
				caughtMsg: 'internal')
		.trackInvalidFile(file, reason)
		return stop?
		}

	processFile(file)
		{
		if false is pdfOb = .getBody(file)
			return
		.processFileAfterCompressed(pdfOb)
		}

	processFileAfterCompressed(pdfOb)
		{
		if .mergedOb.Empty?()
			.initializeBodyVariables(pdfOb)
		else
			.appendBody(pdfOb)
		if .maxMerge <= .totalObj
			throw .LimitError
		}

	maxCompressedError: 'compressed file size over maximum'
	getBody(filename)
		{
		if .isJpg?(filename) and false is filename = .jpgToPdf(filename)
			return false
		if .fileSize(filename) is 0
			{
			.trackInvalidFile(filename, 'empty file')
			return false
			}

		if false is pdfOb = .getFileOb(filename, .getFileContent(filename))
			{
			// No need to add to the InvalidFiles again if pdf is over the max compressed
			// size limit; If code gets here then it quits rebuilding pdf
			if not .imageSizeOverLimit
				.trackInvalidFile(filename, .maxCompressedError)
			return false
			}
		if false is .checkCompression(filename)
			return false
		return pdfOb
		}

	fileSize(filename)
		{
		if .filesData isnt false
			return .getFileContent(filename).Size()

		return .filestorage.FileSize(filename)
		}

	getFileContent(filename)
		{
		if .filesData is false
			return false
		return .filesData.FindOne({ filename.Suffix?(it.path) }).fileData
		}

	checkCompression(filename)
		{
		if .compress isnt true or .maxCompressedFileSize is false
			return true
		if .compressed? isnt true
			{
			.trackInvalidFile(filename, 'file has nothing compressible')
			return false
			}
		if .overCompressionLimit?(.fileSize(filename) - .totalImageSizeReduction)
			{
			.trackInvalidFile(filename, .maxCompressedError)
			return false
			}
		return true
		}

	overCompressionLimit?(size)
		{
		return .maxCompressedFileSize isnt false and size > .maxCompressedFileSize
		}

	getFileOb(file, fileData = false)
		{
		pdfOb = Object()
		objs = Object()
		trailers = Object()
		reader = new PdfReader(:fileData)
		newFile = .fileClass(file, 'r', :fileData)
			{|f|
			if not .isValidPdf(reader, f)
				throw "Invalid pdf"
			reader.Read(f, objs, trailers)
			if false is pdfOb = .processObjects(reader, f, objs, trailers)
				return false
			}
		pdfOb.filename = newFile
		return pdfOb
		}

	filesData: false
	fileClass(file, mode = 'r', fileData = false, block = function(unused){})
		{
		if .filesData isnt false
			{
			(.memoryFile)(file, :mode, :fileData, :block)
			return file
			}

		f = .filestorage.GetAccessibleFilePath(file)
		if f isnt file
			.cleanupFiles.Add(f)
		File(f, :mode, :block)
		return f
		}

	getter_filestorage()
		{
		// so suneido.js does not need to load FileStorage record onto browser
		return Global('FileStorage')
		}

	memoryFile: class
		{
		CallClass(file, mode /*unused*/ = '', fileData = false, block = false)
			{
			f = new this(file, fileData)
			block(f)
			}

		New(.file, .fileData)
			{
			}

		Write(s)
			{
			.fileData $= s
			}

		Get()
			{
			return .fileData
			}
		}

	isValidPdf(reader, f)
		{
		headerLength = 100
		if false is file = reader.FileRead(f, 0, headerLength)
			return false
		return file.Trim().Prefix?("%PDF")
		}

	imageSizeOverLimit: false
	processObjects(reader, f, objs, trailers)
		{
		pdfOb = .initializePdfOb()
		for obj in objs
			{
			if .compressLimit?()
				{
				.imageSizeOverLimit = true
				.trackInvalidFile(.curFile, .maxCompressedError)
				return false
				}

			s = obj.head

			if 0 is Number(s.Extract('(^|[\[' $ .ws $ '])([0-9]+)', 2))
				continue

			.updateLinearized(s, pdfOb)
			.updateSuneidoFormat(s, pdfOb)

			if .failedToCompressImage?(obj, s, reader, f)
				return false

			.afterObjCompressed(s, pdfOb, reader, f, obj, objs, trailers)
			}

		.afterAllObjectProcessed(pdfOb, trailers)

		return pdfOb
		}

	afterAllObjectProcessed(pdfOb, trailers)
		{
		.updateNumObj(pdfOb)

		.processTrailers(trailers, pdfOb)

		.handleLinearized(pdfOb)

		.afterOneMerged(pdfOb)
		}

	afterObjCompressed(s, pdfOb, reader, f, obj, objs = #(), trailers = #())
		{
		if s.Has?('ObjStm')
			try
				pdfOb.objs.Append(.convertStream(s,
					reader.FileRead(f, obj.streamStart, obj.streamEnd)))
			catch (e)
				throw .securedPdf?(e, objs, trailers)
					? 'Secured pdf'
					: e
		else
			pdfOb.objs.Add(obj)
		}

	initializePdfOb()
		{
		return Object(
			numObj: 0,
			linearized?: false,
			suneidoFormat: false,
			objs: Object())
		}

	compressLimit?()
		{
		return .compress is true and .overCompressionLimit?(.totalImageSize)
		}

	failedToCompressImage?(obj, s, reader, f)
		{
		return not .handleImageCompression(obj, s, reader, f) and
			.maxCompressedFileSize isnt false
		}

	securedPdf?(e, objs, trailers)
		{
		return e.Prefix?('Zlib') and
			(objs.Any?({ it.head.Has?('/Encrypt') }) or
				trailers.Any?({ it.head.Has?('/Encrypt') }))
		}

	updateNumObj(pdfOb, s = '')
		{
		largestObj = definedSize = 0
		pdfOb.objs.Each({ largestObj = Max(.objSize(it), largestObj) })
		if false isnt result = s.Match("/Size[" $ .ws $ "]+([0-9]+)")
			definedSize = Number(s[result[1][0] :: result[1][1]]) - 1
		pdfOb.numObj = Max(pdfOb.numObj, definedSize, largestObj, pdfOb.objs.Size())
		}

	objSize(obj)
		{
		return Number(obj.head.Tr(.ws, ' ').BeforeFirst(' 0'))
		}

	updateLinearized(s, pdfOb)
		{
		if pdfOb.linearized? isnt false
			return
		pdfOb.linearized? = pdfOb.objs.Empty?() and s.Has?("/Linearized")
		}

	updateSuneidoFormat(s, pdfOb)
		{
		pdfOb.suneidoFormat = pdfOb.suneidoFormat or
			s.Has?("/Producer (Suneido PDF Generator)")
		}

	handleImageCompression(obj, head, reader, f)
		{
		if .compress is true and
			.compressibleImageObject?(head) and
			false is .compressImage(obj, reader, f, head)
			return false
		return true
		}

	compressibleImageObject?(head)
		{
		head = head.Tr(" \t\r\n").Lower()
		filter = head.AfterFirst('/filter')
		if filter.Prefix?("[")
			filter = filter.BeforeFirst("]")
		else
			{
			filter = filter.AfterFirst("/")
			if filter.Has?('>>')
				filter = filter.BeforeFirst(">>")
			if filter.Has?('/')
				filter = filter.BeforeFirst("/")
			}
		return head.Has?('/subtype/image') and filter.Tr("^a-z") is 'dctdecode'
		}

	additionalEntries: #(SMask)
	compressImage(obj, reader, f, head)
		{
		baseImageFileName = GetAppTempPath() $ Display(Timestamp()).Tr('#.')
		imgFileName = baseImageFileName $ '.jpg'
		compressedImgFileName = baseImageFileName $ '_compressed.jpg'
		.cleanupFiles.Add(imgFileName)
		.cleanupFiles.Add(compressedImgFileName)
		reader.ExtractStreamToJPG(f, obj, imgFileName)
		if false is ImageHandler.Compress(imgFileName, compressedImgFileName,
			extraOptions: Object(colorSpace: .imageObjectColorSpace(head)))
			return false

		return .afterImageCompressed(compressedImgFileName, imgFileName, head, obj)
		}

	afterImageCompressed(compressedImg, origImg, head, obj)
		{
		if .afterMergedAsync isnt false
			{
			newLength = compressedImg.Size()
			originalLength = origImg.Size()
			}
		else
			{
			newLength = .fileSize(compressedImg)
			originalLength = .fileSize(origImg)
			}
		if newLength >= originalLength // image did not compress, retain original
			{
			.totalImageSize += originalLength
			return true
			}
		.totalImageSizeReduction += originalLength - newLength
		compressed = .afterMergedAsync isnt false ? compressedImg : GetFile(compressedImg)
		if not Jpeg.RunWithCatch({ jpeg = Jpeg(compressed) })
			return false

		.totalImageSize += newLength
		.compressed? = true
		height = jpeg.GetHeight()
		width = jpeg.GetWidth()
		colorSpace = jpeg.GetColorSpace()
		head = PdfDriver.BuildImageObjectHead(width, height, colorSpace, newLength,
			additionalEntries: .getAdditionalEntries(head))
		obj.head = obj.head.BeforeFirst('<<') $ head $ obj.head.AfterLast('>>')
		// there must be a newline after "stream"
		if obj.head.Suffix?('stream')
			obj.head $= '\n'
		obj.Delete('streamStart')
		obj.Delete('streamEnd')
		if .afterMergedAsync isnt false
			obj.streamData = compressedImg
		else
			obj.streamFile = compressedImg
		obj.streamSize = newLength
		return true
		}

	getAdditionalEntries(head)
		{
		result = Object()
		for entry in .additionalEntries
			if head.Has?(entry)
				result.Add(.findEntry(head, entry))

		return result
		}

	findEntry(head, entry)
		{
		find = head[head.Find(entry) ..]
		if find isnt result = find.BeforeFirst(`/`)
			return result.Trim(" \t\r\n")

		return find.BeforeFirst(`>>`).Trim(" \t\r\n")
		}

	imageObjectColorSpace(head)
		{
		cs = head.Tr('\r\n').Lower().AfterFirst('colorspace').Trim()
		firstChar = cs[0]
		if firstChar is `/`
			return .getColorSpace1(cs.AfterFirst(`/`))
		else if Numberable?(firstChar)
			return .getColorSpace1(cs)
		else if firstChar is `[`
			return .getColorSpace2(cs, `[`, `]`)
		else if cs.Prefix?(`<<`)
			return .getColorSpace2(cs, `<<`, `>>`)

		return 'nomatch'
		}

	getColorSpace1(cs)
		{
		if cs isnt result = cs.BeforeFirst(`/`)
			return result

		return cs.BeforeFirst(`>>`)
		}

	getColorSpace2(str, after, before)
		{
		return str.AfterFirst(after).BeforeFirst(before)
		}

	processTrailers(trailers, pdfOb)
		{
		for trailer in trailers
			{
			if trailer.head.Has?('/Encrypt')
				throw "Secured pdf"
			.updateNumObj(pdfOb, trailer.head)
			}
		}

	convertStream(segment, stream)
		{
		start = segment.Match("/Type[" $ .ws $ "]*/ObjStm")[0][0]
		headerStart = segment.FindLast('<<', start)
		headerEnd = segment.Find('>>', start)
		header = segment[headerStart .. headerEnd]
		first = Number(header.Extract("/First[" $ .ws $ "]+([0-9]+)"))
		n = Number(header.Extract("/N[" $ .ws $ "]+([0-9]+)"))
		// pdf format is stream begins after first \n after keyword stream
		// Zlib.Uncompress only handles "/Filter /FlateDecode"
		stream = header.Has?('/Filter')
			? .unzip(stream.AfterFirst("\n"))	// don't need to trim
			: stream.AfterFirst("\n").RemoveSuffix("\n")
		indexes = stream[.. first].Tr(.ws, " ").Split(" ").Map(Number)
		indexes.Add(0) // Unused, just a placeholder
		indexes.Add(stream.Size())
		objects = Object()
		for (i = 0; i < n; i++)
			{
			objNum = indexes[i * 2]
			objStart = indexes[i * 2 + 1] + first
			objEnd = indexes[i * 2 + 3/*=offset*/] + first // Start of next object
			objects[i] = Object(
				head: '\n' $ objNum $ " 0 obj\n" $ stream[objStart .. objEnd] $
					"\nendobj\n"
				tail: '')
			}
		return objects
		}

	unzip(str)
		{
		if .afterMergedAsync isnt false
			{
			// Pako is sensitive on the last new line characters
			// TODO: read the length property from stream object directly,
			// TODO: need to handle when the length property is pointing to another object
			try
				return SuUI.GetCurrentWindow().PakoInflate(str.RightTrim('\r\n'))
			catch
				try
					return SuUI.GetCurrentWindow().PakoInflate(str.RightTrim('\n'))
				catch // when stream ends with \n\r
					return SuUI.GetCurrentWindow().PakoInflate(str.RightTrim('\r'))
			}
		return Zlib.Uncompress(str)
		}

	handleLinearized(pdfOb)
		{
		// always remove first object remove second obj if it has xref
		// can move this to within the file block?
		if pdfOb.linearized? is false
			return

		deleteList = Object()
		if pdfOb.objs[1].head.Has?("/XRef")
			deleteList.Add(1)
		deleteList.Add(0)
		pdfOb.objs.Delete(@deleteList)
		}

	initializeBodyVariables(pdfOb)
		{
		if pdfOb.suneidoFormat
			{
			.fetchPagesInfo(pdfOb)
			.totalObj = pdfOb.numObj
			.kids = .getSuneidoFormatKids(pdfOb)
			.totalPages = pdfOb.pageCount
			.parent = `/Parent ` $ pdfOb.pages.BeforeFirst(" ")
			.mergedOb.Add(pdfOb)
			}
		else
			.convertBodyToSuneidoFormat(pdfOb)
		}

	fetchPagesInfo(pdfOb)
		{
		if false is catalog = pdfOb.objs.FindIf({|obj| obj.head.Has?("/Catalog")})
			throw "Root page node not found"
		pdfOb.catalog = catalog
		pdfOb.pages = pdfOb.objs[catalog].head.
			Extract("/Pages[" $ .ws $ "]+([0-9]+[" $ .ws $ "]+0[" $ .ws $ "]+R)")
		pagesObjNum = pdfOb.pages.Extract('^[0-9]+')

		for (i = pdfOb.objs.Size() - 1; i >= 0; i--)
			if pdfOb.objs[i].head =~ '^[' $ .ws $ ']' $ pagesObjNum $
				'[' $ .ws $ ']+0[' $ .ws $ ']+obj'
				break
		if i < 0
			throw 'PDFMerger - cannot find the start of page tree object'
		pdfOb.pagesPos = i

		pdfOb.pageCount =
			Number(pdfOb.objs[i].head.Extract("/Count[" $ .ws $ "]+([0-9]+)"))
		}

	/* getSuneidoFormatKids
	*  Gets the list of objects that are pages for the PDF. Gets the contents of
	*  the Kids array which is a list of all pages. Returns in the format of
	*  16 0 R 18 0 R ...
	*  Params:
	*		body: A Object that contains the body of the PDF
	*  Returns: A string containg the contents of the Kids array in the format of
	*		14 0 R 16 0 R 20 0 R ...
	*/
	getSuneidoFormatKids(pdfOb)
		{
		return pdfOb.objs[pdfOb.pagesPos].head.
			AfterFirst('/Kids').AfterFirst('[').BeforeFirst(']').Trim()
		}

	// The basic format for the start of Suneido Pdf files
	headerFormat: "%PDF-1.3\n%\xe9\xe9\xe9\xe9\n"
	bodyFormat: #("\n1 0 obj <</Type /Catalog /Pages 2 0 R>>\nendobj",
		"\n2 0 obj <</Type /Pages /Kids [] /Count 0>>\nendobj",
		"\n3 0 obj << /Producer (Suneido PDF Generator) >>\nendobj")
	convertBodyToSuneidoFormat(pdfOb)
		{
		firstPdf = Object(
			numObj: 3,
			linearized?: false,
			suneidoFormat: true,
			catalog: 0,
			pages: '2 0 R',
			pagesPos: 1,
			objs: .bodyFormat.Map({ Object(head: it, tail: "") }))
		.mergedOb.Add(firstPdf)
		.totalObj = 3
		.kids = ""
		.totalPages = 0
		.parent = "/Parent 2"
		.appendBody(pdfOb)
		}

	/* appendBody
	*  Responsible for building the body of the new merged pdf, which is an array of pdf
	*  objects. It offsets object numbers so there is no overlap between the
	*  two input files.
	*  Params:
	*		pdfOb: An array of objects that is a list of objects in second pdf
	*/
	appendBody(pdfOb)
		{
		.offsetObjNums(pdfOb, .totalObj)
		.fetchPagesInfo(pdfOb)
		.updateParentRef(pdfOb)

		.cleanUpBody(pdfOb)
		.totalObj += pdfOb.numObj
		.totalPages += pdfOb.pageCount
		.kids $= " " $ pdfOb.pages
		.mergedOb.Add(pdfOb)
		}

	/* offsetObjNums
	*  Responsible for offsetting all object numbers in the body of a PDF so that
	*  they are a continuation of the object numbers of the PDF they are being appended
	*  to. Ensures that there is no overlap in the object numbers between the two
	*  files with the exception of references to the page tree node and the root,
	*  which are made common between files.
	*  Params:
	*		pdfOb: A object containing the body of the PDF being appended onto another
	*		numObj1: The number of objects in the first PDF
	*  Returns: A string containing the modified body
	*/
	offsetObjNums(pdfOb, numObj1)
		{
		ws = .ws
		for obj in pdfOb.objs
			{
			obj.head = obj.head.Replace(
				"(^|[\[" $ .ws $ "])[0-9]+[" $ .ws $ "]+0[" $ .ws $ "]+(obj|\<R\>)",
				{|s|
				idx = Number(s.BeforeLast('0').Trim(ws $ "["))
				idx += numObj1
				str = s[0].Number?() ? '' : s[0]
				str $= idx $ ' 0' $ s.AfterLast('0')
				})
			}
		}

	ws: '\x00\x09\x0A\x0C\x0D\x20' // see pdf reference 1.7 > table 3.1
	updateParentRef(pdfOb)
		{
		oldPagesObj = pdfOb.objs[pdfOb.pagesPos].head
		pdfOb.objs[pdfOb.pagesPos].head = oldPagesObj.BeforeFirst("<<") $ "<<" $ .parent $
			" 0 R" $ oldPagesObj.AfterFirst("<<")
		}

	/* cleanUpBody
	*  Completes the process of removing unnecessary parts of the second PDF.
	*  Completely removes the PDF version header, Catalog and Pages objects.
	*  It assumes that the information in the Pages object is no longer needed.
	*  Params:
	*		pdfOb: A object that is the body that contains the objects being removed
	*/
	cleanUpBody(pdfOb)
		{
		pdfOb.objs.Delete(pdfOb.catalog)
		pdfOb.Delete(#catalog)
		}

	updatePages()
		{
		pagesPos = .mergedOb[0].pagesPos
		// only replace the root node 'kids'
		str = .mergedOb[0].objs[pagesPos].head.
			Replace("/Count [0-9]+", "/Count " $ .totalPages, 1)
		str = str.BeforeFirst('/Kids') $
			'/Kids [' $ .kids $ ']' $ str.AfterFirst('/Kids').AfterFirst(']')
		.mergedOb[0].objs[pagesPos].head = str
		}

	/* buildXRef
	*  Responsible for building the XRef portion of the PDF based on the body of the
	*  PDF. Builds the table from scratch, not relying on the previous XRef tables
	*  from the source PDF files.
	*  Returns: A string containing the XRef table as well as the trailer and end of file
	*/
	buildXRef()
		{
		locs = .calcLocations(.mergedOb)
		firstPdf = .mergedOb[0]
		rootIdx = firstPdf.catalog
		s = firstPdf.objs[rootIdx].head
		root = s[1..s.Find(' ', 1)]
		return PdfDriver.BuildXRef(.totalObj + 1, locs, root, .totalLength)
		}

	/* calcLocations
	*  Responsible for calculating the byte locations of each object in the body
	*  for use in the XRef table. Locations are the number of bytes from the start
	*  of the PDF file to the object declaration "n 0 obj" where n is the object
	*  number. Assumes that the highest object number is equal to the number of objects.
	*  Params:
	*		mergedOb: A object containing the merged body of PDF files
	*  Returns: An object containing the byte locations of each object
	*/
	calcLocations(mergedOb)
		{
		pos = .headerFormat.Size() + 1
		objectLocs = Object()
		objectLocs.Add(0) //entry for object 0
		for pdfOb in mergedOb
			{
			for obj in pdfOb.objs
				{
				s = obj.head
				idx = Number(s[1..s.Find(' ', 1)])
				objectLocs[idx] = pos
				streamSize = .getStreamSize(obj)
				pos += obj.head.Size() + obj.tail.Size() + streamSize
				}
			}
		.totalLength = pos
		for (i = 1; i <= .totalObj; i++)
			if not objectLocs.Member?(i)
				objectLocs[i] = .totalLength
		return objectLocs
		}

	getStreamSize(obj)
		{
		if obj.Member?(#streamSize)
			return obj.streamSize
		else if obj.Member?(#streamStart)
			return obj.streamEnd - obj.streamStart
		return 0
		}

	/* finish
	*  Outputs the current merged PDF to a file with the name specified in the
	*  constructor.
	*/
	buf: ''
	bufSize: 0
	finish()
		{
		if not .InvalidFiles.Empty?()
			return
		overwrite = .filesData is false and .mergedOb.HasIf?(
			{|pdfOb|
			pdfOb.Member?(#filename) and
			Paths.Basename(pdfOb.filename) is Paths.Basename(.newfile)
			})
		tmp = overwrite ? GetAppTempFullFileName("pdf") : .newfile

		outputToTmp = .outputToTmp
		.runWithCatch(tmp) { outputToTmp(tmp) }
		if .InvalidFiles.Empty?() and overwrite
			.copyFile(tmp, .newfile)
		}

	outputToTmp(tmp)
		{
		.fileClass(tmp, 'w', fileData: '')
			{|f|
			.writeWithBuf(f, .headerFormat)
			for pdfOb in .mergedOb
				.write(f, pdfOb)
			.writeWithBuf(f, .xRef)
			.flush(f)
			if .filesData isnt false
				.newfile.fileData = f.Get()
			}
		}

	write(outFile, pdfOb)
		{
		if not pdfOb.Member?(#filename)
			{
			for obj in pdfOb.objs
				.writeWithBuf(outFile, obj.head $ obj.tail)
			return
			}

		.fileClass(pdfOb.filename, 'r', fileData: .getFileContent(pdfOb.filename))
			{|inFile|
			reader = new PdfReader(fileData: .getFileContent(pdfOb.filename))
			for obj in pdfOb.objs
				{
				.writeWithBuf(outFile, obj.head)
				if obj.Member?(#streamFile)
					.writeStreamFile(obj, outFile)
				else if obj.Member?(#streamData)
					.writeWithBuf(outFile, obj.streamData)
				else if obj.Member?(#streamStart)
					{
					from = obj.streamStart
					while from < obj.streamEnd
						{
						to = Min(from + PdfReader.BytesPerRead, obj.streamEnd)
						.writeWithBuf(outFile, reader.FileRead(inFile, from, to))
						from = to
						}
					}
				.writeWithBuf(outFile, obj.tail)
				}
			}
		}

	writeStreamFile(obj, outFile)
		{
		.fileClass(obj.streamFile, 'r',
			fileData: .getFileContent(obj.streamFile))
			{|f2|
			while false isnt s = f2.Read(PdfReader.BytesPerRead)
				.writeWithBuf(outFile, s)
			}
		}

	writeWithBuf(f, s)
		{
		if .bufSize + s.Size() > PdfReader.BytesPerRead
			.flush(f)
		.buf $= s
		.bufSize += s.Size()
		}

	flush(f)
		{
		f.Write(.buf)
		.buf = ""
		.bufSize = 0
		}

	InvalidFilesMsg(invalidFiles)
		{
		if invalidFiles.Empty?()
			return ''

		lastFile = invalidFiles.Last()
		if lastFile.Has?(' (last file attempted, ')
			return 'Unable to merge files into one PDF:\r\n\r\n' $
				invalidFiles.Join('\r\n') $
				'\r\n\r\nPlease review the above list and adjust accordingly.'

		return 'Unable to append the following attachments to PDF:\n\n' $
			invalidFiles.Join('\n') $
			'\n\nPlease check if they are corrupted or secured.'
		}

	/* filterFiles
	*  Filters the object of input files based on the types of files that can be merged.
	*  Params:
	*		files: An object containing file names
	*  Returns: An object containing the file names that can be merged
	*/
	FilterFiles(files)
		{
		validFiles = Object()
		for file in files
			if .extensions.Has?(file.Trim().AfterLast('.').Lower())
				validFiles.Add(file.Trim())
		return validFiles
		}

	cleanUp()
		{
		if .filesData isnt false
			return
		for f in .cleanupFiles
			DeleteFile(f)
		}
	}
