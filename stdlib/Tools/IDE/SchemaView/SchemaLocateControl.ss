// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: SchemaLocate
	New()
		{
		super(['AutoChoose', SchemaLocateControl.TableNames, width: 20,
			cue: 'Locate by name', allowOther:])
		}
	TableNames(prefix)
		{
		tables = .tables()
		extract = function (x) { return [x.name] }
		names = tables.FlatMap(extract)
		if prefix is ""
			return names
		match = '\<(?i)(?q)' $ prefix
		list = names.Filter({ it =~ match }) // exact matches
		others = tables.Filter({ it.trimmed =~ match })
		return list.Copy().MergeUnion(others.FlatMap(extract))
		}
	tables()
		{
		tables = Object()
		QueryApply('views')
			{ |x|
			tables.Add([name: x.view_name, trimmed: x.view_name.Tr("_")])
			}
		QueryApply('tables')
			{|x|
			tables.Add([name: x.table, trimmed: x.table.Tr("_")])
			}
		return tables.Sort!(By(#name))
		}
	NewValue(value)
		{
		.value = value
		.Send('Locate', .value)
		}
	value: ''
	Get()
		{
		return .value
		}
	}