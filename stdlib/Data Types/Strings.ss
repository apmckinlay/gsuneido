// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	As(e) // This is to handle cases where strings are being treated as exceptions
		{
		return e
		}
	LineCount()
		{
		return this is "" ? 0 : this.Count('\n') +
			(this.Has?('\n') and this.AfterLast('\n') is "" ? 0 : 1)
		}
	FirstLine()
		{
		i = .Find('\n')
		return this[..i].RightTrim('\r')
		}
	Lines()
		{
		return Lines(this)
		}
	RemoveBlankLines()
		{
		return .Replace('^[ \t]*\r?\n')
		}
	RemovePrefix(prefix)
		{
		return this.Prefix?(prefix) ? this[prefix.Size()..] : this
		}
	RemoveSuffix(suffix)
		{
		return this.Suffix?(suffix) and suffix isnt "" ? this[.. -suffix.Size()] : this
		}
	ChangeEol(eol)
		{
		return eol is '\n'
			? .Tr('\r')
			: .Replace('\r?\n', eol)
		}
	LineFromPosition(pos)
		{
		return this[.. pos].Count('\n')
		}
	LineAtPosition(pos)
		{
		if pos >= .Size()
			return ""
		org = .FindLast('\n', pos)
		if org is false
			org = 0
		else
			++org
		end = .Find('\n', pos)
		this[org .. end].RightTrim('\r')
		}
	StartPositionOfLine(line)
		{
		pos = -1
		while line > 0 and pos < .Size()
			{
			pos = .Find('\n', pos + 1)
			--line
			}
		return pos is .Size() ? pos : pos + 1
		}
	Capitalize()
		{ return this[0].Upper() $ this[1..].Lower() }
	UnCapitalize()
		{ return this[0].Lower() $ this[1..] }
	CapitalizeWords(lower = true)
		{
		str = lower ? .Lower() : this
		return str.
			Replace("\<(po|ne|nw|se|sw|usa|llc)\>", "\U&").
			Replace("([^[:alnum:]'])([[:lower:]])", '\1\u\2').
			Replace('^[[:lower:]]', '\u&')
		}
	Trim(chars = " \t\r\n")
		{
		chars = "^" $ chars
		first = .Find1of(chars)
		last = .FindLast1of(chars)
		return this[first .. last + 1]
		}
	LeftTrim(chars = " \t\r\n")
		{
		first = .Find1of("^" $ chars)
		return this[first..]
		}
	RightTrim(chars = " \t\r\n")
		{
		last = .FindLast1of("^" $ chars)
		return last is false ? '' : this[.. last + 1]
		}
	SplitCSV(fields = false, string_vals = false) // TODO extract code to separate record
		{
		line = this
		values = Object()
		for (i = first = quotes = 0, n = line.Size(); i <= n; ++i)
			{
			if ((i is n) or ((line[i] is ',') and ((quotes % 2) is 0)))
				{
				x = line[first .. i]
				if (x.Number?() and not string_vals)
					x = Number(x)
				else
					{
					if (x[0] is '"' and x[-1] is '"')
						x = x[1 .. -1]
					if (x.Has?('"'))
						x = x.Replace('""', '"')
					}
				values.Add(x)
				while line[i + 1].White?()
					++i
				first = i + 1
				}
			else if (line[i] is '"')
				++quotes
			}
		if (fields isnt false)
			{
			record = Object()
			for (i = 0, n = Min(values.Size(), fields.Size()); i < n; ++i)
				record[fields[i]] = values[i]
			values = record
			}
		return values
		}

	SplitFixedLength(map)
		{
		return FixedLength.Split(this, map)
		}

	ReplaceSubstr(i, n, s)
		// pre: i and n are integers >= 0 AND s is a string
		{
		return this[.. i] $ s $ this[i + n ..]
		}

	LeftFill(minSize, char = " ")
		// pre:		minSize is a positive integer AND
		//			char is a string of length 1
		// post:	returns a string of at least minSize characters
		//			with leading char's added if necessary
		{
		Assert(char.Size() is 1)
		return char.Repeat(minSize - .Size()) $ this
		}

	TruncateLeftFill(size, char = ' ')
		{
		return this[:: size].LeftFill(size, char)
		}

	RightFill(minSize, char = " ")
		// pre:		minSize is a positive integer AND
		//			char is a string of length 1
		// post:	returns a string of at least minSize characters
		//			with trailing char's added if necessary
		{
		Assert(char.Size() is 1)
		return this $ char.Repeat(minSize - .Size())
		}

	TruncateRightFill(size, char = ' ')
		{
		return this[:: size].RightFill(size, char)
		}

	Center(minSize, char = " ")
		{
		Assert(char.Size() is 1)
		fill = minSize - .Size()
		left = (fill / 2).Int()
		right = fill - left
		return char.Repeat(left) $ this $ char.Repeat(right)
		}

	SplitOnFirst(delimiter = ' ')
		{
		i = .Find(delimiter)
		return Object(this[..i], this[i + delimiter.Size() ..])
		}
	SplitOnLast(delimiter = ' ')
		{
		i = .FindLast(delimiter)
		if (i is false)
			i = .Size()
		return Object(this[..i], this[i + delimiter.Size() ..])
		}

	BeforeFirst(delimiter)
		{
		return this[.. .Find(delimiter)]
		}
	AfterFirst(delimiter)
		{
		return this[.Find(delimiter) + delimiter.Size() ..]
		}
	BeforeLast(delimiter)
		{
		return this[.. .FindLast(delimiter)]
		}
	AfterLast(delimiter)
		{
		i = .FindLast(delimiter)
		return i is false ? "" : this[i + delimiter.Size() ..]
		}

//	Cut(sep)
//		{
//		i = this.Find(sep)
//		if i is -1
//			return this, ""
//		else
//			return this[..i], this[i + sep.Size() ..]
//		}

	Capitalized?()
		{
		return this =~ `\A[[:upper:]]`
		}
	FindRx(rx)
		{
		m = .Match(rx)
		return m is false ? .Size() : m[0][0]
		}
	FindRxLast(rx)
		{
		for (i = .Size() - 1; i >= 0; --i)
			if this[i..] =~ rx
				return i
		return false
		}
	ExtractAll(pattern)
		{
		if false is matches = this.Match(pattern)
			return false
		return matches.Map!( { this[it[0]::it[1]] })
		}

	WrapLines(width) // TODO extract code to separate record
		{
		lines = .Lines()
		wraplines = Object()
		for line in lines
			{
			if (line.Size() < width)
				wraplines.Add(line)
			else
				{
				words = line.Tr("\t", " ").Split(" ")
				sizedwords = Object()
				for word in words
					for part in word.Divide(width)
						sizedwords.Add(part)
				newline = ""
				for word in sizedwords
					{
					word = word.Trim()
					if word is ""
						continue
					line = newline is ""
						? word : newline $ ' ' $ word
					if (line.Size() <= width)
						newline $= newline is "" ? word : ' ' $ word
					else
						{
						wraplines.Add(newline.Trim())
						newline = word
						}
					}
				wraplines.Add(newline.Trim())
				}
			}
		return wraplines
		}
	Divide(n = 1)
		{
		ob = Object()
		.MapN(n, { ob.Add(it); '' })
		return ob
		}
	White?() // like Blank? but returns false for ""
		{
		return .Size() > 0 and .Find1of("^ \t\r\n") >= .Size()
		}
	Blank?() // like White? but returns true for ""
		{
		return .Find1of("^ \t\r\n") >= .Size()
		}
	In?(x)
		{
		return x.Has?(this)
		}
	ToHex()
		{
		return .Map({ it.Asc().Hex().LeftFill(2, '0') })
		}
	FromHex()
		{
		return .MapN(2, { ('0x' $ it).Compile().Chr() })
		}

	Base64Encode()
		{
		return Base64.Encode(this)
		}
	Base64Decode()
		{
		return Base64.Decode(this)
		}

	Xor(key)
		{
		return StringXor(this, key)
		}

	GlobalName?()
		{
		return this =~ `\A[[:upper:]]\w*[!?]?\Z`
		}
	LocalName?()
		{
		return this =~ `\A[[:lower:]]\w*[!?]?\Z`
		}
	Identifier?()
		{
		return this =~ `\A_?[[:alpha:]]\w*[!?]?\Z`
		}
	DynamicName?()
		{
		return this[0] is '_' and this[1..].LocalName?()
		}

	ForEachMatch(pat, block) // non-overlapping matches
		{
		for (n = .Size(), i = 0; i < n; )
			{
			match = .Match(pat, i)
			if match is false
				return
			try
				block(match)
			catch (e, "block:")
				if e is "block:break"
					break
				// else block:continue ... so continue
			i = match[0][0] + Max(1, match[0][1])
			}
		}
	ForEach1of(chars, block)
		{
		i = 0
		n = .Size()
		while n isnt i = .Find1of(chars, i)
			block(i++)
		}
	SafeEval()
		{
		if this is ""
			return ""
		try
			return .Compile()
		catch
			throw "invalid SafeEval: " $ Display(this)
		}
	Map(block)
		{
		return .MapN(1, block)
		}
	Escape()
		{
		return .Map(function (c)
			{
			' ' <= c and c <= '~' and c isnt `\` and c isnt '"'
				? c
				: #('\t': '\\t', '\r': '\\r', '\n': '\\n', '"': `\"`, `\`: `\\`).
					GetDefault(c, '\\x' $ c.Asc().Hex().LeftFill(2, '0'))
			})
		}
	Ellipsis(maxLength, atEnd = false)
		{
		if .Size() <= maxLength
			return this
		halfVal = (maxLength / 2).Floor()
		return not atEnd
			? this[.. halfVal] $ '...' $ this[-halfVal ..]
			: this[.. maxLength] $ '...'
		}
	Has1of?(chars)
		{
		return this.Find1of(chars) < .Size()
		}
	UniqueChars()
		{
		set = ''
		for c in this
			if not set.Has?(c)
				set $= c
		return set
		}
	StripInvalidChars()
		{
		notValid = '^' $ // not
			'\t' $ // tab
			'\n' $ // newline
			'\r' $ // return
			'\x20-\x7f' $ // printableAscii
			'\xC0-\xFF' // printableAnsii
		// Remove invalid chars
		return this.Tr(notValid)
		}
	RandChar()
		{
		return this[Random(this.Size())]
		}
	Shuffle()
		{
		return .Split().Shuffle!().Join()
		}
	}
