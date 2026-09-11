// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	sig_As: "(e) :string"
	As(e) // This is to handle cases where strings are being treated as exceptions
		{
		return e
		}

	sig_LineCount: "() :number"
	LineCount()
		{
		return this is ""
			? 0
			: .Count('\n') + (.Has?('\n') and .AfterLast('\n') is "" ? 0 : 1)
		}

	sig_FirstLine: "() :string"
	FirstLine()
		{
		i = .Find('\n')
		return this[..i].RightTrim('\r')
		}

	sig_Lines: "() :object"
	Lines()
		{
		return Lines(this)
		}

	sig_RemoveBlankLines: "() :string"
	RemoveBlankLines()
		{
		return .Replace("^[ \t]*\r?\n")
		}

	sig_RemovePrefix: "(prefix) :string"
	RemovePrefix(prefix)
		{
		return .Prefix?(prefix) ? this[prefix.Size() ..] : this
		}

	sig_RemoveSuffix: "(suffix) :string"
	RemoveSuffix(suffix)
		{
		return .Suffix?(suffix) and suffix isnt "" ? this[.. -suffix.Size()] : this
		}

	sig_ChangeEol: "(eol) :string"
	ChangeEol(eol)
		{
		return eol is '\n' ? .Tr('\r') : .Replace("\r?\n", eol)
		}

	sig_LineFromPosition: "(pos) :number"
	LineFromPosition(pos)
		{
		return this[..pos].Count('\n')
		}

	sig_LineAtPosition: "(pos) :string"
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
		this[org..end].RightTrim('\r')
		}

	sig_StartPositionOfLine: "(line) :number"
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

	sig_Capitalize: "() :string"
	Capitalize()
		{
		return this[0].Upper() $ this[1..].Lower()
		}

	sig_UnCapitalize: "() :string"
	UnCapitalize()
		{
		return this[0].Lower() $ this[1..]
		}

	sig_CapitalizeWords: "(lower = true) :string"
	CapitalizeWords(lower = true)
		{
		str = lower ? .Lower() : this
		return str.
			Replace("\<(po|ne|nw|se|sw|usa|llc)\>", "\U&").
			Replace("([^[:alnum:]'])([[:lower:]])", "\1\u\2").
			Replace("^[[:lower:]]", "\u&")
		}

	sig_Trim: "(chars = ' \t\r\n') :string"
	Trim(chars = " \t\r\n")
		{
		chars = '^' $ chars
		first = .Find1of(chars)
		last = .FindLast1of(chars)
		return this[first .. last+1]
		}

	sig_LeftTrim: "(chars = ' \t\r\n') :string"
	LeftTrim(chars = " \t\r\n")
		{
		first = .Find1of('^' $ chars)
		return this[first..]
		}

	sig_RightTrim: "(chars = ' \t\r\n') :string"
	RightTrim(chars = " \t\r\n")
		{
		last = .FindLast1of('^' $ chars)
		return last is false ? "" : this[.. last+1]
		}

	sig_SplitCSV: "(fields = false, string_vals = false) :object"
	SplitCSV(fields = false, string_vals = false) // TODO extract code to separate record
		{
		line = this
		values = Object()
		for (i = first = quotes = 0, n = line.Size(); i <= n; ++i)
			if ((i is n) or ((line[i] is ',') and ((quotes % 2) is 0)))
				{
				x = line[first..i]
				if x.Number?() and not string_vals
					x = Number(x)
				else
					{
					if x[0] is '"' and x[-1] is '"'
						x = x[1..-1]
					if x.Has?('"')
						x = x.Replace('""', '"')
					}
				values.Add(x)
				while line[i+1].White?()
					++i
				first = i + 1
				}
			else if line[i] is '"'
				++quotes
		if fields isnt false
			{
			record = Object()
			for (i = 0, n = Min(values.Size(), fields.Size()); i < n; ++i)
				record[fields[i]] = values[i]
			values = record
			}
		return values
		}

	sig_SplitFixedLength: "(map) :object"
	SplitFixedLength(map)
		{
		return FixedLength.Split(this, map)
		}

	sig_ReplaceSubstr: "(i, n, s) :string"
	ReplaceSubstr(i, n, s)
		{
		return this[..i] $ s $ this[i+n ..]
		}

	sig_LeftFill: "(minSize, char = ' ') :string"
	LeftFill(minSize, char = ' ')
		{
		Assert(char.Size() is 1)
		return char.Repeat(minSize - .Size()) $ this
		}

	sig_TruncateLeftFill: "(size, char = ' ') :string"
	TruncateLeftFill(size, char = ' ')
		{
		return this[::size].LeftFill(size, char)
		}

	sig_RightFill: "(minSize, char = ' ') :string"
	RightFill(minSize, char = ' ')
		{
		Assert(char.Size() is 1)
		return this $ char.Repeat(minSize - .Size())
		}

	sig_TruncateRightFill: "(size, char = ' ') :string"
	TruncateRightFill(size, char = ' ')
		{
		return this[::size].RightFill(size, char)
		}

	sig_Center: "(minSize, char = ' ') :string"
	Center(minSize, char = ' ')
		{
		Assert(char.Size() is 1)
		fill = minSize - .Size()
		left = (fill / 2).Int()
		right = fill - left
		return char.Repeat(left) $ this $ char.Repeat(right)
		}

	sig_SplitOnFirst: "(delimiter = ' ') :object"
	SplitOnFirst(delimiter = ' ')
		{
		i = .Find(delimiter)
		return [this[..i], this[i + delimiter.Size() ..]]
		}

	sig_SplitOnLast: "(delimiter = ' ') :object"
	SplitOnLast(delimiter = ' ')
		{
		i = .FindLast(delimiter)
		if i is false
			i = .Size()
		return [this[..i], this[i + delimiter.Size() ..]]
		}

	sig_BeforeFirst: "(delimiter) :string"
	BeforeFirst(delimiter)
		{
		return this[.. .Find(delimiter)]
		}

	sig_AfterFirst: "(delimiter) :string"
	AfterFirst(delimiter)
		{
		return this[.Find(delimiter) + delimiter.Size() ..]
		}

	sig_BeforeLast: "(delimiter) :string"
	BeforeLast(delimiter)
		{
		return this[.. .FindLast(delimiter)]
		}

	sig_AfterLast: "(delimiter) :string"
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
	sig_Capitalized?: "() :boolean"
	Capitalized?()
		{
		return this =~ `\A[[:upper:]]`
		}

	sig_FindRx: "(rx) :number"
	FindRx(rx)
		{
		m = .Match(rx)
		return m is false ? .Size() : m[0][0]
		}

	sig_FindRxLast: "(rx) :false|number"
	FindRxLast(rx)
		{
		for (i = .Size() - 1; i >= 0; --i)
			if this[i..] =~ rx
				return i
		return false
		}

	sig_ExtractAll: "(pattern) :false|object"
	ExtractAll(pattern)
		{
		if false is matches = .Match(pattern)
			return false
		return matches.Map!({ this[it[0] :: it[1]] })
		}

	sig_WrapLines: "(width) :object"
	WrapLines(width) // TODO extract code to separate record
		{
		lines = .Lines()
		wraplines = Object()
		for line in lines
			if line.Size() < width
				wraplines.Add(line)
			else
				{
				words = line.Tr('\t', ' ').Split(' ')
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
					line = newline is "" ? word : newline $ ' ' $ word
					if line.Size() <= width
						newline $= newline is "" ? word : ' ' $ word
					else
						{
						wraplines.Add(newline.Trim())
						newline = word
						}
					}
				wraplines.Add(newline.Trim())
				}
		return wraplines
		}

	sig_Divide: "(n = 1) :object"
	Divide(n = 1)
		{
		ob = Object()
		.MapN(n, { ob.Add(it); "" })
		return ob
		}

	sig_White?: "() :boolean"
	White?() // like Blank? but returns false for ""
		{
		return .Size() > 0 and .Find1of("^ \t\r\n") >= .Size()
		}

	sig_Blank?: "() :boolean"
	Blank?() // like White? but returns true for ""
		{
		return .Find1of("^ \t\r\n") >= .Size()
		}

	sig_In?: "(x) :boolean"
	In?(x)
		{
		return x.Has?(this)
		}

	sig_ToHex: "() :string"
	ToHex()
		{
		return .Map({ it.Asc().Hex().LeftFill(2, '0') })
		}

	sig_FromHex: "() :string"
	FromHex()
		{
		return .MapN(2, { Suneido.Compile(("0x" $ it)).Chr() })
		}

	sig_Base64Encode: "() :string"
	Base64Encode()
		{
		return Base64.Encode(this)
		}

	sig_Base64Decode: "() :string"
	Base64Decode()
		{
		return Base64.Decode(this)
		}

	sig_Xor: "(key) :string"
	Xor(key)
		{
		return StringXor(this, key)
		}

	sig_GlobalName?: "() :boolean"
	GlobalName?()
		{
		return this =~ `\A[[:upper:]]\w*[!?]?\Z`
		}

	sig_LocalName?: "() :boolean"
	LocalName?()
		{
		return this =~ `\A[[:lower:]]\w*[!?]?\Z`
		}

	sig_Identifier?: "() :boolean"
	Identifier?()
		{
		return this =~ `\A_?[[:alpha:]]\w*[!?]?\Z`
		}

	sig_DynamicName?: "() :boolean"
	DynamicName?()
		{
		return this[0] is '_' and this[1..].LocalName?()
		}

	sig_ForEachMatch: "(pat, block) :unknown"
	ForEachMatch(pat, block) // non-overlapping matches
		{
		for (n = .Size(), i = 0; i < n;)
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

	sig_ForEach1of: "(chars, block) :unknown"
	ForEach1of(chars, block)
		{
		i = 0
		n = .Size()
		while n isnt i = .Find1of(chars, i)
			block(i++)
		}

	sig_SafeEval: "() :unknown"
	SafeEval()
		{
		if this is ""
			return ""
		try
			return Suneido.Compile(this)
		catch
			throw "invalid SafeEval: " $ Display(this)
		}

	sig_Map: "(block) :string"
	Map(block)
		{
		return .MapN(1, block)
		}

	sig_Escape: "() :string"
	Escape()
		{
		return .Map(function(c)
			{
			' ' <= c and c <= '~' and c isnt `\` and c isnt '"'
				? c
				: #('\t': `\t`, '\r': `\r`, '\n': `\n`, '"': `\"`,
					`\`: `\\`).GetDefault(c, "\\x" $ c.Asc().Hex().LeftFill(2, '0'))
			})
		}

	sig_Ellipsis: "(maxLength, atEnd = false) :string"
	Ellipsis(maxLength, atEnd = false)
		{
		if .Size() <= maxLength
			return this
		halfVal = (maxLength / 2).Floor()
		return not atEnd
			? this[..halfVal] $ "..." $ this[-halfVal ..]
			: this[..maxLength] $ "..."
		}

	sig_Has1of?: "(chars) :boolean"
	Has1of?(chars)
		{
		return .Find1of(chars) < .Size()
		}

	sig_UniqueChars: "() :string"
	UniqueChars()
		{
		set = ""
		for c in this
			if not set.Has?(c)
				set $= c
		return set
		}

	sig_StripInvalidChars: "() :string"
	StripInvalidChars()
		{
		notValid = '^' $
			// not
			'\t' $
			// tab
			'\n' $
			// newline
			'\r' $
			// return
			"\x20-\x7f" $
			// printableAscii
			"\xC0-\xFF" // printableAnsii
		// Remove invalid chars
		return .Tr(notValid)
		}

	sig_RandChar: "() :string"
	RandChar()
		{
		return this[Random(.Size())]
		}

	sig_Shuffle: "() :string"
	Shuffle()
		{
		return .Split().Shuffle!().Join()
		}
	}
