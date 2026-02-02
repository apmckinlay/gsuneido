// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	chunkSize: 2000
	fileExceptionMsg: "SHOW: Unable to access file"
	CallClass(file)
		{
		if false is firstChunk = .getFile(file, limit: .chunkSize)
			throw .fileExceptionMsg
		return .firstChunkValid?(firstChunk)
		}

	// extracted for testing
	getFile(filename, limit = false)
		{
		filename = FileStorage.GetAccessibleFilePath(filename)
		return GetFile(filename, limit)
		}

	unPrintableChars: '[^[:print:][:space:]]'

	utf8Marker: '\xEF\xBB\xBF'
	nullMarker: '\x00'

	firstChunkValid?(firstChunk)
		{
		chunkIsWholeFile = firstChunk.Size() < .chunkSize
		if chunkIsWholeFile and firstChunk.Suffix?('\x1a')
			firstChunk = firstChunk[..-1] // exclude any potential EOF char

		// ignore extended ascii chars
		lower = 128
		upper = 255
		extendedAsciiRange = lower.Chr() $ "-" $ upper.Chr()
		firstChunk = firstChunk.RemovePrefix(.utf8Marker).
			Tr(extendedAsciiRange).Tr(.nullMarker)
		return not firstChunk.Blank?() and firstChunk !~ .unPrintableChars
		}
	}