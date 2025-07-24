// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s, recordName, name, orig_name, textSearch? = false)
		{
		if not .fastCheck(s, name)
			return false

		if textSearch? is false and .useAST?(name, orig_name)
			{
			results = .findByAST(s, recordName, name, orig_name, findFirst?:)
			return results.Empty?() ? false : results[0]
			}

		m = s.Match(.buildRegex(name, orig_name))
		return m is false ? false : m.Last()[0]
		}

	fastCheck(s, name)
		{
		return s.Has?(name) // initial fast check
		}

	FindAllPos(s, recordName, name, orig_name, textSearch? = false)
		{
		if not .fastCheck(s, name)
			return #()

		if textSearch? is false and .useAST?(name, orig_name)
			return .findByAST(s, recordName, name, orig_name)

		results = Object()
		s.ForEachMatch(.buildRegex(name, orig_name))
			{ |m|
			results.Add(m.Last()[0])
			}
		return results
		}

	useAST?(name, orig_name/*unused*/)
		{
		return name[0].Upper?()
		}

	findByAST(s, recordName, name, orig_name, findFirst? = false)
		{
		s = RemoveUnderscoreRecordName(recordName, s)

		searches = Object(name)
		if orig_name isnt name
			searches.Add(orig_name)

		skipFn? = { |node, parents| .skipFind?(node, parents, s, searches) }
		strCompare = { |target| .matchStrReference?(target, searches) }
		if String?(results = AstSearch(s, searches, [strCompare], :findFirst?, :skipFn?))
			return #()

		return results.Map({ it.pos }).Sort!()
		}

	skipFind?(node, parents, s, searches)
		{
		if node.pos not in (0, false) and
			not searches.Any?({ s[node.pos..node.end].Has?(it) })
			return true

		// to avoid finding ".Member"
		if parents.Size() > 0 and
			parents.Last().type is 'Mem' and
			Same?(parents.Last().mem, node)
			return true

		return false
		}

	matchStrReference?(value, searches)
		{
		if not String?(value)
			return false

		return searches.Any?({ |search| value is search or
			value.Prefix?(search $ '.') and value.AfterFirst('.').Identifier?() or
			value.Prefix?(search $ '(') and value.Suffix?(')') })
		}

	buildRegex(name, orig_name)
		{
		return name isnt orig_name
			? "\<(" $ name $ "|" $ orig_name $ ")\>"
			: name[0].Upper?()
				? "(^|[^.])\<(" $ name $ ")\>"
				: "\<(" $ name $ ")\>"
		}
	}