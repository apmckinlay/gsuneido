// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// determine the differences between two sequences
// based on "An O(ND) Difference Algorithm and Its Variations" by Eugene Myers
class
	{
	CallClass(a, b)
		{
		if false is vs = .search(a, b)
			throw "Diff failed"
		.trace(vs, a, b)
		}
	search(a, b)
		{
		n = a.Size()
		m = b.Size()
		max = n + m
		vs = Object()
		v = Object()
		v[1] = 0
		for (d = 0; d <= max; ++d)
			{
			vs[d] = v.Copy() // TODO: don't need all entries
			for (k = -d; k <= d; k += 2)
				{
				if (k is -d or (k isnt d and v[k - 1] < v[k + 1]))
					x = v[k + 1]
				else
					x = v[k - 1] + 1
				y = x - k
				while (x < n and y < m and a[x] is b[y])
					{
					++x
					++y
					}
				v[k] = x
				//Print(x ", " y)
				if (x >= n and y >= m)
					return vs
				}
			//Print()
			}
		return false
		}
	trace(vs, a, b)
		{
		n = a.Size()
		m = b.Size()
		x = n
		y = m
		edits = Object()
		for (d = vs.Size() - 1; d > 0; --d)
			{
			//Print("d = " d)
			k = x - y
			//Print("k = " k)
			v = vs[d]
			forever
				{
				//Print(x ", " y)
				if (k < d and v[k + 1] is x)
					{
					--y
					edits.Add(Object(x, "I", b[y]))
					break
					}
				else if (k > -d and v[k - 1] is x - 1)
					{
					--x
					edits.Add(Object(x, "D", a[x]))
					break
					}
				Assert(a[x - 1] is b[y - 1])
				--x
				--y
				}
			//Print(x ", " y)
			//Print()
			}
		return edits // Note: reverse order (easiest for Apply)
		}
	Apply(a, edits)
		{
		// edit format is #(at, I or D, value)
		b = a.Copy()
		for (e in edits)
			{
			if e[1] is "D"
				b.Delete(e[0])
			else //if (e[1] is "I")
				b.Add(e[2] at: e[0])
			}
		return b
		}
	SideBySide(a, b)
		{
		edits = Diff(a, b)
		ai = bi = 0
		prev = #('', '')
		result = Object()
		for (ei = edits.Size() - 1; ei >= 0; --ei)
			{
			e = edits[ei]
			while (ai isnt e[0])
				{
				Assert(a[ai] is b[bi], "equal")
				result.Add(prev = Object(a[ai++], "", b[bi++]))
				}
			if (e[1] is "D")
				{
				Assert(e[2] is a[ai], "delete")
				if (prev[1] is ">")
					{
					for (i = result.Size(); i > 0 and result[i-1][1] is ">"; --i)
						{}
					r = result[i]
					r[0] = a[ai++]
					r[1] = "#"
					}
				else
					result.Add(prev = Object(a[ai++], "<", ""))
				}
			else // e[1] is "I"
				{
				Assert(e[2] is b[bi], "insert")
				if (prev[1] is "<")
					{
					for (i = result.Size(); i > 0 and result[i-1][1] is "<"; --i)
						{}
					r = result[i]
					r[1] = "#"
					r[2] = b[bi++]
					}
				else
					result.Add(prev = Object("", ">", b[bi++]))
				}
			}
		while (ai < a.Size())
			result.Add(Object(a[ai++], "", b[bi++]))
		return result
		}
	Three(base, a, b) // base vs. merged, side by side
		{
		ae = Diff(base, a)
		be = Diff(base, b)
		result = Object()
		for (i in base.Members())
			result[i] = Object(base[i], '', base[i])
		for (ai = bi = 0; ai < ae.Size() or bi < be.Size(); )
			{
			if (bi >= be.Size() or (ai < ae.Size() and ae[ai][0] > be[bi][0]))
				{
				e = ae[ai++]
				which = 1
				}
			else if (ai >= ae.Size() or be[bi][0] > ae[ai][0])
				{
				e = be[bi++]
				which = 2
				}
			else if (ae[ai] is be[bi])
				{
				e = ae[ai]
				which = ''
				++ai
				++bi
				}
			else // conflict
				{
				Assert(ae[ai][0] is be[bi][0])
				at = ae[ai][0]
				for (; bi < be.Size() and be[bi][0] is at; ++bi)
					if (be[bi][1] is "I")
						result.Add(Object('', "+2", be[bi][2]), :at)
					else // "D"
						{
						result[at][1] = "-2"
						result[at][2] = ""
						}
				for (; ai < ae.Size() and ae[ai][0] is at; ++ai)
					if (ae[ai][1] is "I")
						result.Add(Object('', "+1", ae[ai][2]), :at)
					else // "D"
						{
						result[at][1] = "-1"
						result[at][2] = ""
						}
				continue
				}
			if (e[1] is "D")
				{
				r = result[e[0]]
				r[1] = "-" $ which
				r[2] = ""
				}
			else //if (e[1] is "I")
				{
				result.Add(Object('', "+" $ which, e[2]) at: e[0])
				}
			}
		// combine + and -
		for (i = 0; i < result.Size(); ++i)
			{
			a = result[i]
			if a[1] =~ "+" and i > 0 and result[i-1][1] is (w = "-" $ a[1][1..])
				{
				for (j = i; j > 0 and result[j-1][1] is w; --j)
					{}
				d = result[j]
				d[1] = "#" $ w[1..]
				d[2] = a[2]
				result.Delete(i)
				--i
				}
			}
		return result
		}
	Merge(base, a, b)
		{
		t = .Three(base, a, b)
		result = Object()
		for (x in t)
			if (x[1][0] isnt '-')
				result.Add(x[2])
		return result
		}
	}
