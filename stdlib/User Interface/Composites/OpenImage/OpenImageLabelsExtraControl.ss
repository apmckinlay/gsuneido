// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(fieldName = '', type = '')
		{
		super(.layout(fieldName, type))
		}

	layout(fieldName = '', type = '')
		{
		.labels = .setupLabels(fieldName, type)
		return .buildLayout(.labels)
		}

	setupLabels(fieldName, type)
		{
		handlers = Contributions('OpenImageLabelsExtra')
		labels = Object()
		for handler in handlers
			handler(labels, fieldName, type)
		return labels.RemoveIf({ it.Condition() is false }).UniqueValues()
		}

	buildLayout(labels)
		{
		if labels.Empty?()
			return false
		return Object(#Vert,
			#(Static, 'Click the following links to add system labels:')
			#(Skip, small:),
			.columnLayout(labels))
		}

	columnLayout(labels)
		{
		columns = Object()
		labels.Each({ columns.AddUnique(it.Column) })
		columns.Sort!()
		columnLayout = Object(#Horz)
		for col in .. columns.Size()
			{
			vert = Object(#Vert)
			columnLabels = labels.Filter({ it.Column is columns[col] }).Sort!(By(#Label))
			for label in columnLabels
				vert.Add(label.Layout(), #(Skip, small:))
			columnLayout.Add(vert, #(Skip, 50))
			}
		return columnLayout
		}

	Recv(@args)
		{
		call = args[0].RemovePrefix('On_')
		if false is label = .labels.FindIf({ it.Command is call})
			return 0

		.Send('AddLabel', .labels[label].Label)
		}
	}