// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(name, code, search, hint = false)
		{
		lines = Object()
		if name =~ '\.(js|css)$'
			return lines

		code = RemoveUnderscoreRecordName(name, code)

		skipFn? = hint is false
			? false
			: { |node, parents/*unused*/|
				node.pos not in (0, false) and not code[node.pos..node.end].Has?(hint) }
		results = AstSearch(code, search, :skipFn?)
		if String?(results)
			{
			if results.Prefix?('Parse Search text')
				throw results
			return Object(results)
			}

		results.Each()
			{
			from = code.LineFromPosition(it.pos)
			to = code.LineFromPosition(it.end - 1)
			lines.Add(Seq(from, to + 1))
			}
		return lines
		}
	}