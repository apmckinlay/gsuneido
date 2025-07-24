// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query, metrics)
		{
		if Sys.Client?()
			return ServerEval(#QueryTreeHtml, query, metrics)
		return (new this).Do(query, metrics)
		}
	Do(query, metrics)
		{
		.metrics = metrics.Map(.nameToField)
		WithQuery(query)
			{|q|
			while false isnt q.Next()
				{}
			.q = q.Tree()
			.max1 = 1 + .max(.q, .metrics[0])
			s = .tree(0, .q)
			}
		s = Xml('table', '\n' $
				Xml('tr', metrics.Map({ Xml('th', ' ' $ it) }).Join()) $ '\n' $
				s) $ '\n'
		return Xml('html', .head $ '\n' $ Xml('body', '\n' $ s) $ '\n')
		}
	tree(indent, q)
		{
		si = (q[.metrics[0]] * 10 / .max1).Int() /*= 10 levels */
		Assert(0 <= si and si <= 9) /*= 10 levels */
		vals = .metrics.Map({ .format(q, it) })
		vals[0] = Xml('span', vals[0], class: 's' $ si)
		vals.Map!({ Xml('td', it, align: 'right') })
		ind = '<span class="indent"></span>'.Repeat(indent)
		s = Xml('tr', vals.Join() $ Xml('td', ind $ Xml('span', q.string))) $ '\n'
		++indent
		if q.type is 'view'
			return s $ .tree(indent, q.source)
		switch q.nchild
			{
		case 0:
			return s
		case 1:
			return .tree(indent, q.source) $ s
		case 2:
			return .tree(indent, q.source1) $ s $ .tree(indent, q.source2)
			}
		}
	format(q, name)
		{
		if name is 'timecost'
			return (q.tget / q.cost).RoundToPrecision(2)
		value = q[name]
		if name is 'frac'
			return value.RoundToPrecision(2)
		return value.Format("###,###,###,###")
		}
	max(q, mem)
		{
		switch q.nchild
			{
		case 0:
			return q[mem]
		case 1:
			return Max(q[mem], .max(q.source, mem))
		case 2:
			return Max(q[mem], Max(.max(q.source1, mem), .max(q.source2, mem)))
			}
		}
	nameToField: (
		"Est Cost": cost,
		"Self Cost": costself,
		"Est Rows": nrows,
		"Est Frac": frac,
		"Total Time": tget,
		"Self Time": tgetself,
		"N Gets": ngets,
		"N Sels": nsels,
		"N Looks": nlooks,
		"Time / Cost": timecost,
		)
	head:
`		<head>
		<style>
		body { font-size: 100%; }
		.s0 { font-size: 1em; }
		.s1 { font-size: 1.1em; }
		.s2 { font-size: 1.2em; background-color: #fee5d9; }
		.s3 { font-size: 1.3em; background-color: #fee5d9; }
		.s4 { font-size: 1.4em; background-color: #fcbba1; }
		.s5 { font-size: 1.5em; background-color: #fcbba1; }
		.s6 { font-size: 1.6em; background-color: #fc9272; }
		.s7 { font-size: 1.7em; background-color: #fc9272; }
		.s8 { font-size: 1.9em; background-color: #fb6a4a; font-weight: bold; }
		.s9 { font-size: 2em; background-color: #fb6a4a; font-weight: bold; }
		table { border-collapse: collapse; }
		th { text-align: right; }
		td {
			margin-top: 0;
			margin-bottom: 5px;
			white-space: nowrap;
			border: solid;
			border-style: dotted;
			border-color: blue;
			border-width: 1px 0;
		}
		td span { font-family: monospace; }
		th, td { padding-right: 10px; }
		.indent { margin-right: 40px; border-right: 1px dotted blue; }
		</style>
		</head>`
	}