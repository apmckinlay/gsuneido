// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
/* e.g.
grammar = #(
	expr: (term expr2)
	expr2: ('+' term or '-' term or)
	term: (factor term2)
	term2: ('*' factor or '/' factor or)
	factor: (number or id)
	)
Parse(grammar, 'expr', '2 * x + 4')
*/
class
	{
	CallClass(grammar, rule, text)
		{
		p = .parse(grammar, rule, text)
		if String?(p)
			throw p
		if p[1] isnt ''
			throw 'failure'
		return p[0]
		}
	parse(grammar, rule, text)
		{
		def = grammar[rule]
		for a in .alts(def)
			{
//Print(trying: a)
			p = .parseSeq(grammar, a, text)
			if Object?(p)
				return p // return first one that succeeds
			}
		return "failure" // none of alternatives matched
		}
	alts(rule)
		{
		alt = []
		alts = []
		for x in rule
			if x is 'or'
				{
				alts.Add(alt)
				alt = []
				}
			else
				alt.Add(x)
		alts.Add(alt)
		return alts
		}
	parseSeq(grammar, def, text)
		{
		result = []
		for x in def
			{
			p = .element(grammar, x, text)
			if String?(p)
				return p
			text = p[1]
			result.Add(p[0])
			}
		return [def.Size() is 1 ? result[0] : result, text]
		}
	element(grammar, x, text)
		{
//Print('element', Display(x), Display(text))
		text = text.LeftTrim()
		switch x
			{
		case 'id':
			return .id(text)
		case 'number':
			return .number(text)
		default:
			if grammar.Member?(x)
				return .parse(grammar, x, text)
			else
				return .literal(x, text)
			}
		}
	id(text)
		{
		id = text.Extract('^[a-z]+')
//Print('id', Display(text), "=>", id)
		return id is false ? "failure" : [id, text[id.Size()..]]
		}
	number(text)
		{
		n = text.Extract('^-?\d+')
//Print('number', Display(text), "=>", n)
		return n is false ? "failure" : [n, text[n.Size()..]]
		}
	literal(x, text)
		{
		return text.Prefix?(x) ? [x, text[x.Size()..]] : "failure"
		}
	}