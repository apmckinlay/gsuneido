// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// abstract base class - derived classes define:
//		ProcessChunk(text)
//		WordChars: "..."
//		ChunkSize: #	// optional, default 64
ScintillaAddon
	{
	changed: false
	ensureInitialized()
		{
		if .changed is false
			.changed = Object(from: 999999, to: 0)
		}

	Modified(scn)
		{
		.ensureInitialized()
		.changed.from = Min(scn.position, .changed.from)
		.changed.to = Max(scn.position + scn.length, .changed.to)
		}

	IgnoreIfEmpty?: true
	IdleAfterChange()
		{
		.ensureInitialized()
		if .changed.from >= .changed.to
			{
			.changed = Object(from: 999999, to: 0)
			return
			}
		n = .GetLength()
		.changed.from = Min(n, .changed.from)
		.changed.to = Min(n, .changed.to)
		text = .get_changed_text()
		if text.Size() > 0 or .IgnoreIfEmpty? is false
			.ProcessChunk(text, .changed.from)
		.changed.from += text.Size() + 1
		if .changed.from >= .changed.to
			{
			.MarkersChanged()
			.changed = Object(from: 999999, to: 0)
			}
		else
			{
			// if not force running the delay on browser,
			// Suneido.js will keep running the delayed callables on the
			// WS server and only returns the result after all the chunks are processed
			_forceOnBrowser = true
			.Defer(.IdleAfterChange, uniqueID: .uniqueId)
			}
		}

	getter_uniqueId()
		{
		return .uniqueId = Timestamp() // once only
		}

	ChunkSize: 64
	get_changed_text()
		{
		while .WordChars.Has?(.GetAt(.changed.from - 1))
			--.changed.from
		to = Min(.changed.from + .ChunkSize, .changed.to)
		while .WordChars.Has?(.GetAt(to))
			++to
		return .GetRange(.changed.from, to)
		}
	}
