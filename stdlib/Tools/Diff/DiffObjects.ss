// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(ob1, ob2)
		{
		s1 = .format(ob1)
		s2 = .format(ob2)
		Diff2Control('DiffObjects', s1, s2, 'first', 'second')
		}

	format(ob, indent = 0)
		{
		s = ''
		if ob.Size() > 10000 /*=limit*/
			{
			for m, vv in ob
				s $= .formatValue(m, vv, indent)
			}
		else
			{
			for m in ob.Members().Sort!()
				s $= .formatValue(m, ob[m], indent)
			}
		return s.RemoveSuffix('\r\n')
		}

	formatValue(m, vv, indent)
		{
		v = Object?(vv)
			? 'OBJECT\r\n' $ .format(vv, indent+1)
			: Instance?(vv)
				? 'INSTANCE\r\n' $ .format(vv, indent+1)
				: Display(vv)
		return '\t'.Repeat(indent) $ m $ ': ' $ v $ '\r\n'
		}
	}