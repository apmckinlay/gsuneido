// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
AstFmtDoc
	{
	New(src)
		{
		toks = Object()
		scan = Scanner(src)
		while scan isnt scan.Next2()
			toks.Add([kind: scan.Type(), text: scan.Text(), end: scan.Position()])
		.toks = toks
		.spans = toks.Filter(
			{
			it.kind is #COMMENT
			}).Map!({ [start: it.end - it.text.Size(), end: it.end] })
		}

	CommentIn?(from, to)
		{
		return .spans.Any?({ it.end > from and it.start < to })
		}

	// Leading consumes trivia up to pos; comments come out on their own line
	// when the source had a newline after them, else attached with a space
	Leading(curr, pos, unusedParam = false)
		{
		docs = Object()
		pend = false
		while curr.done < pos and curr.i < .toks.Size()
			{
			tok = .toks[curr.i]
			++curr.i
			curr.done = tok.end
			if tok.kind is #COMMENT
				{
				if pend isnt false
					{
					docs.Add(pend)
					docs.Add(.Text(' '))
					}
				pend = .suppress?(tok.text, unusedParam) ? false : .Tokc(tok.text)
				}
			else if tok.kind is #NEWLINE and pend isnt false
				{
				docs.Add(pend)
				docs.Add(.Hard)
				pend = false
				}
			}
		if pend isnt false
			{
			docs.Add(pend)
			docs.Add(.Text(' '))
			}
		return .Catl(docs)
		}

	// Trailing consumes comments before pos, and after pos up to the next
	// code token; a blank line in the source sets curr.blank for the next item
	Trailing(curr, pos, parentEnd = 9999999, allowBlank = false, unusedParam = false)
		{
		docs = Object()
		nl = false
		while curr.done < parentEnd and curr.i < .toks.Size()
			{
			tok = .toks[curr.i]
			++curr.i
			curr.done = tok.end
			if tok.kind is #COMMENT
				{
				if not .suppress?(tok.text, unusedParam)
					{
					// annotations hug their token: x/*unused*/, 2/*=EtchedLine*/
					hug = tok.text is "/*unused*/" or tok.text.Prefix?("/*=")
					if nl
						docs.Add(.Hard)
					else if not hug
						docs.Add(.Text(' '))
					docs.Add(.Tokc(tok.text))
					}
				}
			else if tok.kind is #NEWLINE
				{
				if allowBlank and .blank?(tok.text)
					curr.blank = true
				nl = true
				if curr.done >= pos
					return .Catl(docs)
				}
			else if tok.kind isnt #WHITESPACE
				if curr.done > pos
					return .Catl(docs)
			}
		return .Catl(docs)
		}

	SkipTo(curr, pos) // advance the cursor silently, e.g. past verbatim slices
		{
		while curr.done < pos and curr.i < .toks.Size()
			{
			curr.done = .toks[curr.i].end
			++curr.i
			}
		}

	blank?(s)
		{
		return s.Find('\n', s.Find('\n') + 1) < s.Size()
		}

	suppress?(text, unusedParam)
		{
		return unusedParam is true and text is "/*unused*/"
		}
	}
