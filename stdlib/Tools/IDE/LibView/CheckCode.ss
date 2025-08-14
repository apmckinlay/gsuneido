// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MaxLineLength: 90

	// returns true if ok, false if errors, "" if warnings
	CallClass(code, name = false, lib = false, results = false)
		{
		// don't check Suneido versions of functions that are built-in on current exe
		if BuiltinNames().BinarySearch?(name)
			return true
		if not CodeTags.Matches(code)
			return true
		if CheckLibrary.BuiltDate_skip?(code)
			return true
		return (new this(lib, name, code, results))()
		}
	New(.lib, .name, .code, .results)
		{
		.exclude = Sys.Win32?() ? #() : Win32BuiltinNames
		}
	Call()
		{
		if .skip?()
			return true
		.compile_checks()
		.extra_checks()
		if .errors?
			return false
		else if .warnings?
			return ""
		else
			return true
		}
	skip?()
		{
		return .code is "" or .IsWeb?(.name)
		}
	IsWeb?(name)
		{
		return String?(name) and name =~ "(?i)\.(js|css)$"
		}

	compile_checks()
		{
		results = Object()
		// need to handle _Name because string.Compile won't accept it
		code = RemoveUnderscoreRecordName(.name, .code, results)
		try
			code.Compile(results)
		catch (err)
			{
			if err.Prefix?("can't LoadLibrary")
				{
				.result("WARNING: " $ err)
				return
				}
			line = false
			if err.Prefix?("syntax error at line ")
				{
				line = Number(err.Extract("syntax error at line ([0-9]+)")) - 1
				pos = code.StartPositionOfLine(line)
				}
			else // gSuneido
				pos = Number(err.Extract("(syntax|compile) error @([0-9]+)", 2))
			.result("ERROR: " $ err, pos, :line)
			return
			}
		.process_compile_checks(results)
		}

	errors?: false
	warnings?: false
	result(msg, pos = 0, len = 0, line = false)
		{
		if line is false
			line = .code.LineFromPosition(pos)
		if msg.Prefix?("ERROR:")
			.errors? = true
		else
			.warnings? = true
		if .results isnt false
			.results.Add(Object(:msg, :pos, :len, :line))
		}

	process_compile_checks(results)
		{
		libInUse = .libraries.Has?(.lib)
		for w in results
			{
			if String?(w)
				{
				if .gsuneidoSkipWarning?(w, libInUse)
					continue
				pos = Number(w.AfterLast('@'))
				token = w.AfterLast(':').BeforeFirst('@').Trim()
				msg = w.BeforeFirst(' @')
				}
			else // number from RemoveUnderscoreRecordName
				{
				pos = w.Abs()
				scan = Scanner(.code[pos..])
				token = scan.Next()
				if token is "unused" or // built-in on gSuneido
					.code[pos + token.Size() ..].LeftTrim().Prefix?('/*unused*/')
					continue
				if token.Capitalized?()
					{
					if "" is msg = .check_global(token, libInUse)
						continue
					}
				else
					msg = (w < 0 // make messages consistent with gSuneido
						? "WARNING: initialized but not used: "
						: "ERROR: used but not initialized: ") $ token
				}
			if token[0] is '_'
				{ // from RemoveUnderscoreRecordName
				if libInUse and not .check_underscore_name(token)
					.result('ERROR: invalid use of: ' $ token, pos, token.Size())
				}
			else
				.result(msg, pos, token.Size())
			}
		}

	gsuneidoSkipWarning?(w, libInUse)
		{
		if not w.Prefix?("ERROR: can't find:")
			return false
		return libInUse
			? .exclude.Any?({ w.Has?(it) })
			: true
		}

	check_global(token, libInUse)
		{
		try
			Global(token)
		catch(e)
			{
			if e.Has?("can't LoadLibrary")
				return "WARNING: " $ token $ " can't load library"
			else if e.Prefix?("error loading")
				return "ERROR: " $ token $ " has a syntax error"
			else if e.Prefix?("can't find")
				return libInUse
					? "ERROR: can't find: " $ token
					: "" // ignore if lib not in use
			else
				return "ERROR: " $ token $ e
			}
		return "ERROR: can't find: " $ token
		}

	check_underscore_name(token)
		{
		if .name is false
			return false

		pureName = LibraryTags.RemoveTagFromName(.name)
		if token[1..] isnt pureName
			return false

		// record with library tags can only underscore override a record in the same lib
		if pureName isnt .name
			return not QueryEmpty?(.lib, group: -1, name: pureName)

		return .prev_def(.name, .lib, .libraries)
		}
	getter_libraries()
		{
		return .libraries = Libraries() // once only
		}
	prev_def(name, lib, libs)
		{
		if false is i = libs.Find(lib)
			return false
		for (j = 0; j < i; ++j)
			if not QueryEmpty?(libs[j], group: -1, :name)
				return true // found def
		return false
		}

	extra_checks()
		{
		.regex()
		.ast()
		.lineEnds()
		}
	lineEnds()
		{
		if .lib is false or .name is false
			return

		if .HasInvalidLineEnd?(.lib, .name)
			.result('WARNING: saved record uses non-standard line ending characters')
		}

	HasInvalidLineEnd?(lib, name)
		{
		rec = Query1(lib, group: -1, :name)
		return rec isnt false and rec.text =~ '[^\r][\n]'
		}

	regex()
		{
		code = .remove_ignored(.code)
		for bad in .bad
			code.ForEachMatch(bad[0])
				{
				pos = it[0][0]
				len = it[0][1]
				if bad.Member?('filter') and (bad.filter)(:code, :pos, :len, name: .name)
					continue
				msg = (bad.GetDefault(#warning, false) ? "WARNING" : "ERROR") $
					": use of " $ code[pos :: len].Tr(' \t\r\n', ' ').Trim() $
					Opt(' - ', bad[1])
				.result(msg, pos, len)
				}
		}
	remove_ignored(code)
		{
		// need to keep exact same code size so offsets match original
		// need to keep newlines within comments and strings to maintain line numbers
		// can't use Tr instead of Replace because of collapsing
		return ScannerMap(code)
			{|prev2/*unused*/, prev/*unused*/, token, next/*unused*/|
			if token.Prefix?('//')
				token = '//' $ token[2..].Replace('\S', '-')
			else if token.Prefix?('/*') and token.Tr(' ') isnt '/*unused*/'
				token = '/*' $ token[2..-2].Replace('\S', '-') $ '*/'
			else if token[0] in ('"', "'", "`")
				token = token[0] $ token[1..-1].Replace('\S', '-') $ token[-1]
			token
			}
		}
	// Warning: if using backquotes, you can't use \t, \n, etc.
	bad: (
		(`Date\(\).(Begin|End)\>`, 				'use Date.Begin/End()'),
		(`\.Trim\(\) is(nt)? (''|"")`, 			'use .Blank?()'),
		('false is(nt)? [^={};\n]*?\.Find\>', 	'use .Has?')
		('false is(nt)? [^={};\n]*?\.FindIf\>', 'use .Any?')
		(`\<QueryAccum\(`, 						'', 							warning:)
		(`^((Trace|Server)?Print|StackTrace|Inspect|TraceCallStack)\(`,
												'debugging code',				warning:)
		('DoWithTran\(false',					'use Transaction',				warning:)
		(`\<weight:\s\d+`,			'use named weights, e.g. weight: "bold"',	warning:)
		('RetryTransaction\(update:',			'RetryTransaction does not take update:')
		('/\*[ \t\r\n]*?\*/',					'empty comment',				warning:)
		(`\+\s*\<0\>|\<0\>\s*\+[^+]`,			'use Number()',					warning:)
		(`Curry\(`,							'deprecated, use a block instead',	warning:)
		('\?[ \t\r\n]*?true[ \t\r\n]*?:[ \t\r\n]*?false', "useless, remove", 	warning:)
		(`\[0 ?(\.\.|::)[^\]]`,					"omit 0 default",				warning:)
		(`\<class ?:? [A-Z]`,					"omit 'class :'",				warning:)
		(`[^.]\<getter_[A-Z]\w+?\(`				"invalid getter")
		(`[^.]\<Getter_[a-z]\w+?\(`				"invalid getter",				warning:)
		(`^\s*?(\$|and|or|is|isnt)\>[^:]`, "should be at end of previous line",	warning:)
		(`\<super.New\>`,						"should probably be just super(...)")
		(`[^.]\<[a-z]\w*: ?[a-z]\w*\>([^:?.\[\(]|$)`,	"use :name shortcut",	warning:,
			filter: function (code, pos, len)
				{
				s = code[pos :: len][1..]
				name = s.BeforeFirst(':')
				return s.Extract(`: ?([a-z]\w*)`) isnt name or
					code[pos + len - 1 ..].FirstLine() !~ `^\s*(,|\)|$)`
				})
		(`catch ?\(unused\)`,					"omit (unused)", 				warning:)
		(`:[a-zA-Z]\w*`,						"must be preceded by ( or comma",
			filter: function (code, pos)
				{
				if code[pos-1] is ':' or code[..pos][-80..].Has?('dll') /*= line length */
					return true
				for (i = pos-1; i >= 0; i--)
					if code[i] in ('(', '[', ',')
						return true
					else if code[i] not in (' ', '\t', '\r', '\n')
						break
				return false
				})
		(`LocalCmds\(\)`, 						'instance not needed, omit ()')
		('default ?:[ \t\r\n]*throw', 			'omit default that just throws', warning:)
		(`false is(nt)? (CopyFile|MoveFile|DeleteFile(Api)?|CreateDir|DeleteDir)\(`,
												'use true is(nt)')
		(`(if|not) (CopyFile|MoveFile|DeleteFile(Api)?|CreateDir|DeleteDir)\(`,
												'use true is(nt)')
		(`Sort!\(.*unused`, 		'sort functions should use both arguments')
		(`false isnt Query(1|First|Last)\(`,	'use not QueryEmpty?', 			warning:)
		(`false is Query(1|First|Last)\(`,		'use QueryEmpty?', 				warning:)
		(`\<ResourceCounts\(`,					'use Suneido.Info',				warning:)
		(`[^.]\<Sleep\(`,						'use Thread.Sleep')
		(`\<CreateDir\(`,						'use EnsureDir',				warning:)
		(`\.ServerEval\(`,						'use ServerEval(...)')
		(`{ *(([A-Z]\w*)?[.])?[A-Z][a-zA-Z0-9_?!]*\(it\) *}`,
												'unnecessary block',			warning:)
		(`[.]SpyOn[(]["'#]`,					'use the value, not a string',	warning:)
		(`\.Times\(`,							'use for ..',					warning:,
			filter: function (code, pos) { return code[::pos].Suffix?(".Verify") })
		(`\.Merge\((Object\(|Record\(|#\(|#{|\[)`,	'use assignments', 			warning:)
		("\t[.]\w+ = ",				'rules should not have side effects',
			filter: function (name)
				{
				return not String?(name) or
					not name.Prefix?("Rule_") or name.Suffix?("Test")
				})
		(`in Seq\(`,							'use for in ..',	warning:)
		(`Http(s)?\(["'](Get|Put|Post)`, 'use Http(s).Get|Post|Put to check result code',
			warning:)
		(`(?q)Thread.Main?()`,					'use Sys.MainThread?',			warning:,
			filter: function (name) { return name is 'Sys' })
		(`\<Each2\(`,							'use: for m, v',				warning:)
		)

	ast()
		{
		code = RemoveUnderscoreRecordName(.name, .code)
		if false is codeAst = .parseSource(code)
			return

		for bad in .badAst
			{
			hint = AstSearch.GetHint(bad[0])
			skipFn? = hint is false
				? false
				: { |node, parents/*unused*/|
					node.pos not in (0, false) and
						not code[node.pos..node.end].Has?(hint) }
			if String?(results = AstSearch(codeAst, bad[0], :skipFn?))
				{
				.result(results, 0, 0)
				return
				}

			results.Each()
				{
				len = it.end - it.pos
				if bad.Member?('filter') and (bad.filter)(.code, it.pos, len)
					continue
				msg = (bad.GetDefault(#warning, false) ? "WARNING" : "ERROR") $
					": use of " $ .code[it.pos .. it.end].Tr(' \t\r\n', ' ').Trim() $
					Opt(' - ', bad[1])
				.result(msg, it.pos, len)
				}
			}
		}

	parseSource(code)
		{
		try
			codeAst = Suneido.Parse(code)
		catch
			return false
		return Type(codeAst) isnt 'AstNode' ? false : codeAst
		}

	badAst: (
		("QueryApplyMulti(a, block: b, update: false)",
			"QueryApplyMulti should be update"),
		("Assert(x is: true)",
			"just use Assert(x)", warning:),
		("ServerEval('Sys.Linux?')",
			"use Sys.LinuxServer?()"),
		("ServerEval('Timestamp')",
			"unnecessary"),
		("QueryFirst(q) is false",
			"use QueryEmpty?", warning:),
		("QueryFirst(q) isnt false",
			"use QueryEmpty?", warning:),
		("QueryLast(q) is false",
			"use QueryEmpty?", warning:),
		("QueryLast(q) isnt false",
			"use QueryEmpty?", warning:),
		("DeleteFile(f) is false",
			"use: isnt true"),
		("DeleteFile(f) isnt false",
			"use: is true"),
		)
	}