// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// 	map is an object with field names for members and the values
	// 	are objects containing the start position and length ("pos" and "len")
	// 	example: #(field1: (pos: 0, len: 5), field2: (pos: 5, len: 3),
	//		field3: (pos: 8, len: 20))

	// 	for the build method, a 'justify' and 'padChar' member can be added.
	//	'justify' to control which side of the value gets padded to satisfy length
	//	'padChar' will specify the character to pad the value with.
	//	Justification will default to left and padChar will default to a space

	Split(line, map)
		{
		rec = Record()
		for field in map.Members()
			{
			fmap = map[field]
			rec[field] = .formatField(fmap, line[fmap.pos::fmap.len])
			}
		return rec
		}

	formatField(fmap, value)
		{
		if not fmap.Member?('type')
			return value

		switch fmap.type
			{
		case 'string' :
			{
			return fmap.GetDefault('trim', false) ? value.Trim() : value
			}
		case 'date' :
			{
			if false is d = Date(value, fmap.GetDefault('datefmt', 'yMd'))
				throw 'invalid date'
			return d
			}
		case 'number' :
			{
			n = Number(value)
			if fmap.Member?('precision')
				n = n * .1.Pow(fmap.precision) /* = .1 is the base for precision calc*/
			return n
			}}
		}

	Build(record, map)
		{
		strSpec = Object()
		for field in map.Members()
			{
			fmap = map[field]
			val = String(record[field])
			if val.Size() > fmap.len
				val = val[::fmap.len]
			else if val.Size() < fmap.len
				{
				padChar = fmap.GetDefault('padChar', ' ')
				just = fmap.GetDefault('justify', 'left')
				if just isnt 'left' and just isnt 'right'
					throw 'FixedLength: invalid justify in map'
				fillMethod = just is 'left' ? 'RightFill' : 'LeftFill'
				val = val[fillMethod](fmap.len, padChar)
				}
			strSpec.Add(Object(pos: fmap.pos, :val))
			}
		str = ""
		for ob in strSpec.Sort!({ |x,y| x.pos < y.pos })
			str $= ob.val
		return str
		}
	}