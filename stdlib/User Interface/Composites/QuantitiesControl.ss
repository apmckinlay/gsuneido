// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// contributed by Ajith.R
PassthruController
	{
	Name: 'Quantities'
	New(list, cols_head = #("List", "Number"), .saveEmpty = false)
		{
		super(.controls(cols_head))
		.list = .Vert.List
		for (i in list.Split(",").Sort!())
			{
			row = Object()
			row[.cols[0]] = i.BeforeFirst("(").Trim()
			row[.cols[1]] = Number(i.AfterFirst("(").BeforeLast(")").Trim())
			.list.AddRow(row)
			}
		.initial_val = list
		}
	controls(cols_head)
		{
		if cols_head.Size() < 2
			throw "Invalid cols_head: There should be two members "
		if cols_head.Member?('list')
			.cols = Object(cols_head['list'])
		else
			.cols = Object(cols_head[0])
		if cols_head.Member?('num')
			.cols.Add(cols_head['num'])
		else
			.cols.Add(cols_head[1])
		controls = Object('Vert'
			Object('List' noDragDrop:, columns: .cols, xmin: 210, ymin: 300))
		buttons = #('HorzEqual'
			Fill (Button 'OK') Skip (Button 'Clear') Skip (Button 'Cancel') Fill)
		controls.Add(#(Skip 3), buttons)
		return controls
		}
	List_WantNewRow()
		{
		return false
		}
	List_AllowCellEdit(column, row /*unused*/)
		{
		return column isnt 0
		}
	List_WantEditField(col, row /*unused*/, data /*unused*/)
		{
		if col is 1
			return #(Number justify:'LEFT' width: 10)
		}
	On_Cancel()
		{
		.Window.Result(.initial_val)
		}
	OK()
		{
		result = Object()
		ob = .list.Get()
		for i in ob
			result[i[.cols[0]]] = i[.cols[1]]
		if not .saveEmpty
			result.Remove(0, "")
		return result.Map2({|m, v| m $ "(" $ v $ ")" }).Join(',')
		}

	On_Clear()
		{
		for (cnt = 0; cnt < .list.GetNumRows(); ++cnt)
			{
			ob = .list.GetRow(cnt)
			ob[.cols[1]] = 0
			.list.SetRow(cnt, ob)
			}
		}
	}
