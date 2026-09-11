// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
// YAML subset:
// - a single document
// - block and inline (flow) mappings
// - block sequences ("- ") and flow sequences ("[...]")
class
	{
	FromObject(ob, indentation = "")
		{
		if Object?(ob) and not ob.HasNamed?()
			{
			if ob.Empty?()
				return Opt(indentation, "[]\n")
			return indentation $ .flow(ob) $ '\n'
			}
		s = ""
		for m in ob.Members().Sort!()
			{
			value = ob[m]
			s $= indentation $ m $ ':'
			if Object?(value)
				if value.HasNamed?()
					s $= " \n" $ .FromObject(value, indentation $ "    ")
				else
					s $= (value.Empty?() ? " []" : ' ' $ .flow(value)) $ '\n'
			else
				s $= ' ' $ Display(value) $ '\n'
			}
		return s
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
