// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	BytesPerRead: 100000 // cache size 100KB
	tailLength: 10
	ws: '[\x00\x09\x0A\x0C\x0D\x20]' // Whitespace: see pdf reference 1.7 > table 3.1
	New(fileData = false)
		{
		.maxRead = Objects.SizeLimit
		.cache = false
		if fileData isnt false
			.cache = Object(segment: fileData, start: 0, end: fileData.Size(), allRead?:)
		}

	Read(f, objs, trailers)
		{
		start = 0
		tagSize = 9
		pat = Object(
			obj: Object(
				ob: objs,
				start: .ws $ '[0-9]+' $ .ws $ '+[0-9]+' $ .ws $ '+obj',
				end: 'endobj' $ .ws $ '+?',
				size: 6,
				includeStream: false)
			trailer: Object(
				ob: trailers,
				start: .ws $ 'trailer' $ .ws $  '*<<',
				size: 3,
				includeStream:))

		regex = pat.obj.start $ '|' $ pat.trailer.start
		while false isnt startPos = .find(f, regex, start)
			{
			s = .fileReadWithCache(f, startPos, startPos + tagSize)
			obToUse = s.Prefix?('trailer', 1) ? 'trailer' : 'obj'
			start = .fetch(f, startPos, pat[obToUse], obToUse is 'trailer')
			}
		}

	newLine: '\r\n'
	ExtractStreamToJPG(f, obj, outFileName)
		{
		// could use FileRead but that reads the entire stream into memory at once which
		// may not be desirable if the image is large
		chunkSize = 64.Kb()
		size = obj.streamEnd - obj.streamStart
		remaining = size
		f.Seek(obj.streamStart)
		File(outFileName, 'w')
			{ |out_file|
			firstChunk? = true
			while remaining > 0
				{
				readSize = Min(chunkSize, remaining)
				// trim should only really be needed on first and last chunk
				imgChunk = f.Read(readSize)
				if firstChunk?
					{
					firstChunk? = false
					imgChunk = imgChunk.LeftTrim(.newLine)
					}
				else if remaining <= chunkSize
					imgChunk = imgChunk.RightTrim(.newLine)
				remaining -= readSize
				out_file.Write(imgChunk)
				}
			}
		return true
		}

	ExtractStream(obj)
		{
		return .cache.segment[obj.streamStart .. obj.streamEnd].
			LeftTrim(.newLine).
			RightTrim(.newLine)
		}

	LimitError: 'pdf too complex'
	fetch(f, startPos, pattern, trailer = false)
		{
		if trailer
			{
			if false is endPos = .findTrailer(f, startPos)
				throw 'Failed to read pdf: cannot find end of trailer'
			}
		else if false is endPos = .find(f, pattern.end, startPos)
			throw 'Failed to read pdf: closing tag not found'
		if .maxRead <= pattern.ob.Size() + 1
			throw .LimitError
		pattern.ob.Add(.generateObj(f, startPos, endPos + pattern.size,
			pattern.includeStream))
		return endPos + pattern.size
		}

	find(f, pat, start = 0, end = false)
		{
		tail = ''
		pos = start
		while false isnt segment = .fileReadWithCache(f, pos, end)
			{
			segment = tail $ segment
			if segment.Size() isnt loc = segment.FindRx(pat)
				return pos + loc - tail.Size()

			pos += segment.Size() - tail.Size()
			tail = segment[-.tailLength..]
			}
		return false
		}

	findTrailer(f, start)
		{
		tail = ''
		pos = start
		while false isnt segment = .fileReadWithCache(f, pos)
			{
			segment = tail $ segment
			if false isnt idx = .parseNestedTrailer(segment)
				return pos + idx - tail.Size()
			pos += segment.Size() - tail.Size()
			tail = segment
			}
		return false
		}

	parseNestedTrailer(segment)
		{
		firstMatch = false
		idx = 0
		count = 0
		segment.ForEachMatch("<<|>>")
			{ |m|
			idx = m[0][0]
			if firstMatch is false
				firstMatch = idx
			if segment[idx::2] is "<<"
				++count
			else if segment[idx::2] is ">>"
				--count
			if count is 0
				return idx
			}
		if count < 0
			throw 'Should not have unmatched closing tag (>>)'
		return false
		}

	streamStartPat: 'stream'
	streamEndPat: 'endstream'
	lookBackBytes: 100
	generateObj(f, startPos, endPos, includeStream = false)
		{
		if ((not includeStream) and
			(false isnt streamEnd = .find(f, .streamEndPat,
				Max(0, endPos - .lookBackBytes, startPos), endPos)) and
			(false isnt streamStart = .find(f, .streamStartPat, startPos, endPos)))
			{
			streamStart += .streamStartPat.Size()
			head = .FileRead(f, startPos, streamStart)
			tail = .FileRead(f, streamEnd, endPos)
			if head is false or tail is false
				throw "Failed to read pdf: head/tail not found"
			return Object(:streamStart, :streamEnd, :head, :tail)
			}

		return Object(head: .FileRead(f, startPos, endPos), tail: "")
		}

	FileRead(f, startPos, endPos)
		{
		if false is res = .fileReadWithCache(f, startPos, endPos)
			return false
		size = res.Size()
		while size < endPos - startPos
			{
			if false is s = .fileReadWithCache(f, startPos + size, endPos)
				break
			size += s.Size()
			res $= s
			}
		return res
		}

	fileReadWithCache(f, startPos, endPos = false)
		{
		if endPos isnt false and startPos >= endPos
			return false

		if false isnt s = .readFromCache(startPos, endPos)
			return s

		if .noneToRead?(startPos, endPos)
			return false

		f.Seek(startPos)

		if false is s = f.Read(.BytesPerRead)
			return false

		.cache = Object(segment: s, start: startPos, end: startPos + s.Size())
		return .readFromCache(startPos, endPos)
		}

	noneToRead?(startPos, endPos)
		{
		return .cache isnt false and .cache.GetDefault('allRead?', false) is true and
			endPos is false and startPos is .cache.end
		}

	readFromCache(startPos, endPos)
		{
		if .cache is false or startPos < .cache.start or startPos >= .cache.end
			return false

		if endPos is false or endPos > .cache.end
			return .cache.segment[startPos-.cache.start .. ]

		return .cache.segment[startPos-.cache.start..endPos-.cache.start]
		}
	}
