// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddonForChanges
	{
	// Need this to handle deleting trailing characters in which case the text in
	// .IdleAfterChange is empty
	IgnoreIfEmpty?: false
	styleLevel: 100
	ChunkSize: 1024
	WordChars: ""
	Init()
		{
		.indic_error = .IndicatorIdx(level: .styleLevel)
		.marker_error = .MarkerIdx(level: .styleLevel)
		}

	Styling()
		{
		return [
			[level: .styleLevel,
				marker: [SC.MARK_ROUNDRECT, back: CLR.red]
				indicator: [INDIC.BOX, fore: CLR.red]]]
		}

	ProcessChunk(text, pos)
		{
		badChars = "\x80-\xff"
		n = text.Size()
		if text.Find1of(badChars) is n and .MarkerNext(0, .markerMask()) is -1
			return
		.ClearIndicator(.indic_error, pos, n)
		.adjustMarkers(pos, n)
		text.ForEach1of(badChars)
			{ .mark(pos + it) }
		}

	markerMask()
		{
		return 1 << .marker_error
		}

	// We need to clear all markers on a line every time. Even if we only apply 1
	// mark ourselves, scintilla may still force multiple of the same marker onto 1 line
	// when deleting multiple lines that are marked.
	adjustMarkers(pos, n)
		{
		org = .LineFromPosition(pos)
		end = .LineFromPosition(pos + n)
		for (line = org; line <= end; ++line)
			{
			if not .indicFound?(line)
				{
				while .lineMarked?(line)
					.MarkerDelete(line, .marker_error)
				}
			// need to add marker when changed text is only a newline and it forces
			// the invalid character onto a new line - indicator remains but we need mark
			else if not .lineMarked?(line)
				.MarkerAdd(line, .marker_error)
			}
		}

	indicFound?(line)
		{
		// IndicatorStart searches from end to start for an indicator
		linEnd = .GetLineEndPosition(line)
		linStart = .PositionFromLine(line)
		if linEnd is linStart
			return false

		indicStart = .IndicatorStart(.indic_error, linEnd)
		return indicStart isnt 0 and indicStart >= linStart
		}

	lineMarked?(line)
		{
		// .MarkerGet returns a 32-bit integer that indicates which markers were
		// present on the line. Bit 0 is set if marker 0 is present, bit 1 for
		// marker 1 and so on; can be multiple of the marker present on the line
		mask =  1 << .marker_error
		markers = .MarkerGet(line)
		return (markers & mask) isnt 0
		}

	mark(pos)
		{
		.SetIndicator(.indic_error, pos, 1)
		line = .LineFromPosition(pos)
		.MarkerAdd(line, .marker_error)
		}

	MarkersChanged()
		{
		// Scintilla doesn't refresh itself when the last character gets deleted
		// by a delete/backspace, so a drawn error marker is not cleared even though
		// it is removed by Scintilla
		// any manual refresh is valid, e.g. user switching windows
		if .GetLength() is 0
			.Repaint()
		}
	}
