// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Import
	{
	Before()
		{
		do
			line = .Getline()
			while (line isnt false and line isnt "<table>")
		}
	Import1(line)
		{
		if line is "</table>"
			return false
		rec = Record()
		Assert(line is "<record>")
		while ((line = .Getline()) isnt "</record>" and
			line isnt false)
			{
			// line is like: <field type="string">value</field>
			i = line.Find('>')
			j = line.Find(' ')
			field = line[1 :: j - 1]
			field_type = line[j + 7 :: i - (j + 8)]
			endtag = '</' $ field $ '>'
			while (not line.Suffix?(endtag))
				line $= '\n' $ .Getline()
			value = line[i + 1 .. -(j + 2)]
			value = XmlEntityDecode(value)
			switch field_type
				{
			case 'Date':
				rec[field] = Date(value)
			case 'Number':
				rec[field] = Number(value)
			case 'Boolean':
				rec[field] = value is 'true'
			case 'Object':
				rec[field] = value.SafeEval()
			default: // string
				rec[field] = value
				}
			}
		return rec
		}
	DefaultConversion(rec /*unused*/, field /*unused*/)
		{}
	}
