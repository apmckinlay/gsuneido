// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.width = 90, .tabWidth = 4)
		{
		}

	Render(doc)
		{
		rs = [out: "", col: 0, pend: 0, eolc: false]
		stack = [[0, #break, doc]]
		while not stack.Empty?()
			{
			x = stack.PopLast()
			i = x[0]
			m = x[1]
			d = x[2]
			switch d.t
				{
			case #text:
				.stepText(rs, i, d)
			case #verb:
				.stepVerb(rs, i, d)
			case #str:
				.stepStr(rs, i, d)
			case #cat:
				a = d.a
				for (j = a.Size() - 1; j >= 0; --j)
					stack.Add([i, m, a[j]])
			case #nest:
				stack.Add([i + 1, m, d.d])
			case #root:
				.stepRoot(rs, d, stack)
			case #line:
				.stepLine(rs, i, m, d)
			case #hard, #blank:
				.nl(rs, i)
			case #group:
				.stepGroup(rs, i, m, d, stack)
			case #fill:
				.stepFill(rs, i, m, d, stack)
			case #fillr:
				.moreFill(rs, d.a, i, stack)
				}
			}
		return rs.out
		}

	stepText(rs, i, d)
		{
		if rs.eolc and d.s isnt ""
			.nl(rs, i)
		.emit(rs, d.s)
		if d.hb is true // a // comment: nothing else may join this line
			rs.eolc = true
		}

	stepVerb(rs, i, d)
		{
		if rs.eolc
			.nl(rs, i)
		.emit(rs, d.s)
		if d.s.Has?('\n')
			rs.col = d.s.AfterLast('\n').Size()
		}

	stepStr(rs, i, d)
		{
		if rs.eolc
			.nl(rs, i)
		.emitStr(rs, d.s, i)
		}

	stepRoot(rs, d, stack)
		{
		if rs.pend isnt false // line not started: pull it to the margin
			{
			rs.pend = 0
			rs.col = 0
			}
		stack.Add([0, #break, d.d])
		}

	stepLine(rs, i, m, d)
		{
		if m isnt #flat or rs.eolc
			.nl(rs, i)
		else if rs.pend is false // a flat separator at line start is nothing
			.emit(rs, d.s)
		}

	stepGroup(rs, i, m, d, stack)
		{
		if m is #flat
			stack.Add([i, #flat, d.d])
		else if d.hb or not .fits(.width - rs.col, [i, #flat, d.d], stack)
			stack.Add([i, #break, d.d])
		else
			stack.Add([i, #flat, d.d])
		}

	stepFill(rs, i, m, d, stack)
		{
		if m is #flat
			{
			a = d.a
			for (j = a.Size() - 1; j >= 0; --j)
				stack.Add([i, m, a[j]])
			}
		else
			.startFill(rs, d.a, i, stack)
		}

	emit(rs, s)
		{
		if s is ""
			return
		if rs.pend isnt false
			{
			rs.out $= '\t'.Repeat(rs.pend)
			rs.col = rs.pend * .tabWidth
			rs.pend = false
			}
		rs.out $= s
		rs.col += s.Size()
		}

	// emit a plain quoted literal, splitting after a word boundary with a
	// trailing $ when it cannot fit; the pieces continue one indent deeper
	emitStr(rs, s, i)
		{
		q = s[0]
		forever
			{
			if rs.col + s.Size() <= .width
				{
				.emit(rs, s)
				return
				}
			hi = Min(.width - rs.col - 4, s.Size() - 3)/*= close quote + ' $' */
			for (j = hi;; --j)
				{
				if j < 2 // no split point: let it overflow
					{
					.emit(rs, s)
					return
					}
				if s[j] is ' '
					break
				}
			.emit(rs, s[.. j+1] $ q $ " $")
			.nl(rs, i + 1)
			s = q $ s[j+1 ..]
			}
		}

	nl(rs, i)
		{
		while rs.out.Suffix?(' ') or rs.out.Suffix?('\t')
			rs.out = rs.out[..-1]
		rs.out $= '\n'
		rs.pend = i
		rs.col = i * .tabWidth
		rs.eolc = false
		}

	startFill(rs, a, i, stack)
		{
		if a.Empty?()
			return
		if a.Size() > 1
			stack.Add([i, #break, [t: #fillr, a: a[1..]]])
		m = .fits(.width - rs.col, [i, #flat, a[0]], stack) ? #flat : #break
		stack.Add([i, m, a[0]])
		}

	moreFill(rs, a, i, stack)
		{
		if a.Size() < 2
			return
		if a.Size() > 2
			stack.Add([i, #break, [t: #fillr, a: a[2..]]])
		if .leads(a[1]) is 1 // item supplies its own break: keep separator flat
			{
			stack.Add([i, #break, a[1]])
			stack.Add([i, #flat, a[0]])
			return
			}
		pair = [t: #cat, a: [a[0], a[1]]]
		m = .fits(.width - rs.col, [i, #flat, pair], stack) ? #flat : #break
		stack.Add([i, m, a[1]])
		stack.Add([i, m, a[0]])
		}

	leads(d) // what a doc starts with: 1 hard break, 2 content, 0 nothing
		{
		switch d.t
			{
		case #hard:
			return 1
		case #nest, #group, #root:
			return .leads(d.d)
		case #cat, #fill, #fillr:
			for x in d.a
				if 0 isnt r = .leads(x)
					return r
			return 0
		case #text, #verb, #str:
			return d.s is "" ? 0 : 2
		default: // line, blank
			return 2
			}
		}

	// does the first line fit in w columns? pure: scans the unrendered doc,
	// continuing into the pending stack, until a newline or overflow
	fits(w, item, stack)
		{
		work = [item]
		si = stack.Size()
		forever
			{
			if w < 0
				return false
			if work.Empty?()
				{
				if --si < 0
					return true
				x = stack[si]
				}
			else
				x = work.PopLast()
			i = x[0]
			m = x[1]
			d = x[2]
			switch d.t
				{
			case #text:
				if d.hb is true and m is #flat
					return false
				w -= d.s.Size()
			case #str:
				w -= d.s.Size()
			case #verb:
				return m isnt #flat
			case #cat, #fill, #fillr:
				a = d.a
				for (j = a.Size() - 1; j >= 0; --j)
					work.Add([i, m, a[j]])
			case #nest:
				work.Add([i + 1, m, d.d])
			case #line:
				if m isnt #flat
					return true
				w -= d.s.Size()
			case #hard, #blank, #root:
				return m isnt #flat
			case #group:
				if d.hb is true
					return m isnt #flat
				work.Add([i, #flat, d.d])
				}
			}
		}
	}
