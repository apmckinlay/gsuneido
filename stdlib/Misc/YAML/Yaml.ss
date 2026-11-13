// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// YAML subset:
// - a single document
// - block and inline (flow) mappings
// - block sequences ("- ") and flow sequences ("[...]")
class
	{
	FromObject(ob, indentation = "", blockStyle? = false)
		{
		if Object?(ob) and not ob.HasNamed?()
			{
			if ob.Empty?()
				return Opt(indentation, "[]\n")
			return blockStyle?
				? .blockSequence(ob, indentation)
				: indentation $ .flow(ob) $ '\n'
			}
		s = ""
		for m in ob.Members().Sort!()
			{
			value = ob[m]
			s $= indentation $ .formatMember(m) $ ':'
			if Object?(value)
				if value.HasNamed?() or (blockStyle? and not value.Empty?())
					s $= '\n' $ .FromObject(value, indentation $ "  ", :blockStyle?)
				else
					s $= (value.Empty?() ? " []" : ' ' $ .flow(value)) $ '\n'
			else
				s $= ' ' $ (blockStyle? ? .blockScalar(value) : Display(value)) $ '\n'
			}
		return s
		}

	blockSequence(ob, indentation)
		{
		s = ""
		for value in ob
			{
			s $= indentation $ '-'
			if Object?(value) and not value.Empty?()
				{
				childIndent = indentation $ "  "
				child = .FromObject(value, childIndent, blockStyle?:)
				s $= ' ' $ child[childIndent.Size() ..]
				}
			else
				s $= ' ' $ .blockScalar(value) $ '\n'
			}
		return s
		}

	blockScalar(value)
		{
		return String?(value) and value !~ "^#" ? value : .flow(value)
		}

	// YAML treats bare numeric-looking strings as numbers/booleans/null,
	// so quote them to preserve the string type
	formatMember(m)
		{
		return String?(m) and m.Number?()
			? Display(m)
			: m
		}

	flow(value)
		{
		if Object?(value)
			{
			if value.HasNamed?()
				{
				s = ""
				for m in value.Members().Sort!()
					s $= .flowKey(m) $ ": " $ .flow(value[m]) $ ", "
				return s is "" ? "{}" : '{' $ s[..-2] $ '}'
				}
			s = ""
			for v in value
				s $= .flow(v) $ ", "
			return s is "" ? "[]" : '[' $ s[..-2] $ ']'
			}
		return Display(value)
		}

	flowKey(m)
		{
		return String?(m) ? Display(m) : String(m)
		}
	}
