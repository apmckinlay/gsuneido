// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
/*
	text: literal s (no newlines); verb: verbatim s (multiline token);
	str: plain quoted literal the renderer may split with $ to fit;
	line: newline or s when flat (Line ' ', Soft '', Semi '; ');
	hard/blank: newline, never flat; cat: sequence a; nest: +1 tab on breaks;
	group: flat if it fits; fill: pack a, breaking only where needed;
	root: absolute column 0 on its own line, never flattens.
 */
class
	{
	// t is the tag, s is string payload, hb is has break
	Line: (t: line, s: ' ', hb: false)
	Soft: (t: line, s: "", hb: false)
	Semi: (t: line, s: "; ", hb: false)
	Hard: (t: hard, hb:)
	Blank: (t: blank, hb:)

	Text(s)
		{
		return [t: #text, :s, hb: false]
		}

	Tok(s)
		{
		if s.Has?('\n')
			return [t: #verb, :s, hb:]
		return [t: #text, :s, hb: false]
		}

	Tokc(s)
		{
		if s.Has?('\n')
			return [t: #verb, :s, hb:]
		return [t: #text, :s, hb: s.Prefix?("//")]
		}

	Cat(@docs)
		{
		return .Catl(docs)
		}

	Catl(docs)
		{
		a = docs.Filter(
			{
			it isnt false
			}).Map!({ String?(it) ? [t: #text, s: it, hb: false] : it })
		return [t: #cat, :a, hb: a.Any?({ it.hb is true })]
		}

	Nest(doc)
		{
		return [t: #nest, d: doc, hb: doc.hb]
		}

	Root(doc) // absolute column 0 on its own line; never flattens
		{
		return [t: #root, d: doc, hb:]
		}

	Str(s) // a plain quoted literal the renderer may split with $ to fit
		{
		return [t: #str, :s, hb: false]
		}

	Group(doc)
		{
		return [t: #group, d: doc, hb: doc.hb]
		}

	Fill(a)
		{
		return [t: #fill, :a, hb: a.Any?({ it.hb is true })]
		}

	Seq(docs, sep)
		{
		return .Catl(.Interleave(docs, sep))
		}

	Fillsep(docs, sep)
		{
		return .Fill(.Interleave(docs, sep))
		}

	Interleave(docs, sep)
		{
		a = Object()
		for (j = 0; j < docs.Size(); ++j)
			{
			if j > 0
				a.Add(sep)
			a.Add(docs[j])
			}
		return a
		}
	}
