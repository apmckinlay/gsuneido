// Copyright (C) 2026 Axon Development Corporation All rights reserved worldwide.
/* README: The order of the graphs are based on the members of: .data
	By default, this supports:
		- Unspecified object members (IE: Object(val1, val2, val3, ...))
		- Specified object members (IE: Object(0: val1, 2: val2, 5: val3, ...))
		- Date members (IE: Object(#170001: val1, #170002: val2, #170004: val3, ...))
	initDisplayDetails, options: false, 'First', 'Last'
		- This argument controls the initial value for the "displayDetails" control
*/
Controller
	{
	Name: 'BulletGraphs'
	New(.data, .width, .height, .vertical, .axisDensity, .target = false, .min = 0,
		.axisFormat = '#,###,###,###', .heading = '', .initDisplayDetails = false,
		// ========================= Graph Colors =========================
		.good = 0x226322 /*green*/, .satisfactory = 0x666622 /*yellow*/,
		.bad = 0x883322 /*red*/, .inactive = 0x4d4d4d /*gray*/,
		.negative = 0x5e3838 /*muted red*/)
		{
		.graphs = .FindControl('graphs')
		.Send('Data')
		}

	Controls()
		{
		controls = Object('Vert')
		if not .heading.Blank?()
			controls.Add(Object('Heading', .heading))
		return controls.
			Add(Object(.controlContainer, name: 'graphs').Append(.buildGraphs()))
		}

	getter_controlContainer()
		{
		return .controlContainer = .vertical ? 'Vert' : 'Horz'
		}

	buildGraphs()
		{
		.labels = #()
		if .data.Empty?()
			return #()

		graphsBase = .graphsBase(.target, .data, .min, .axisDensity)
		.labels = .data.Members().Sort!()
		axisLabel = .vertical
			? .labels.First()
			: .labels.Last()
		graphs = Object(.graphsContainer)
		for label in .labels
			graphs.Add(.bulletGraph(.graphArgs(graphsBase, label, axisLabel, .data)))
		return Object('Skip', graphs, .displayDetails())
		}

	graphsBase(target, data, min, axisDensity)
		{
		good = satisfactory = 0
		color = .inactive
		if false isnt max = .calcMax(target, data, min)
			{
			good = (data.Sum() / data.Size()).Round(0)
			satisfactory = (good / 2).Round(0)
			color = .determineColor(data, good, satisfactory)
			}
		else
			{
			max = min + 1
			axisDensity = 1
			}
		return Object(range: Object(min, max), :good, :satisfactory, :color, :axisDensity)
		}

	calcMax(target, data, min)
		{
		max = target isnt false
			? Max(target, data.Max())
			: data.Max()
		if max <= min
			return false
		offset = 10.Pow(Max(max.IntDigits() - 2, 1))
		return (max + offset).RoundToNearest(offset)
		}

	getter_graphsContainer()
		{
		return .graphsContainer = .vertical ? 'Horz' : 'Vert'
		}

	determineColor(data, good, satisfactory)
		{
		return data.CountIf({ it >= good }) > data.Size() / 2
			? .good
			: data.CountIf({ it >= satisfactory }) > data.Size() / 2
				? .satisfactory
				: .bad
		}

	graphArgs(graphsBase, label, axisLabel, data)
		{
		graphArgs = graphsBase.Copy()
		graphArgs.label = label
		graphArgs.axis = label is axisLabel
		if graphArgs.range[0] >= graphArgs.value = data[label]
			{
			graphArgs.good = graphArgs.satisfactory = 0
			graphArgs.color = graphArgs.value is graphArgs.range[0]
				? .inactive
				: .negative
			graphArgs.value = 0
			}
		return graphArgs
		}

	bulletGraph(graphArgs)
		{
		return Object('BulletGraph', graphArgs.value,
			satisfactory: graphArgs.satisfactory, good: graphArgs.good,
			target: .target, range: graphArgs.range, color: graphArgs.color,
			vertical: .vertical, width: .width, height: .height,
			outside: 0, axis: graphArgs.axis, axisDensity: graphArgs.axisDensity,
			axisFormat: .axisFormat, selectedColor: CLR.Highlight, name: graphArgs.label)
		}

	Getter_AxisFormat()
		{
		return .axisFormat
		}

	displayDetails()
		{
		value = ' '
		if .initDisplayDetails isnt false and .labels.NotEmpty?()
			{
			label = (.labels[.initDisplayDetails])()
			if 0 isnt display = .displayDetailsSend('construct', label)
				value = display
			}
		return Object('Horz', .displayDetailsSpacer,
			Object('Static', value, justify: 'RIGHT', name: 'displayDetails'))
		}

	displayDetailsSend(event, label)
		{
		return .Send('BulletGraphs_DisplayDetails', event, label, .data[label])
		}

	getter_displayDetailsSpacer()
		{
		return .displayDetailsSpacer = .vertical ? 'Fill' : 'Skip'
		}

	Set(.data)
		{
		.selected = .hover = false
		.graphs.RemoveAll()
		.graphs.AppendAll(.buildGraphs())
		}

	Get()
		{
		return .data
		}

	hover: false
	BulletGraph_Hover(source)
		{
		if .hover isnt source.Name
			.setDisplayDetails('hover', .hover = source.Name)
		return false
		}

	setDisplayDetails(event, label)
		{
		if 0 is display = .displayDetailsSend(event, label)
			return false
		.FindControl('displayDetails').Set(display)
		return true
		}

	selected: false
	BulletGraph_Click(source)
		{
		if .selected isnt false
			{
			if .selected is source.Name
				return false
			.FindControl(.selected).Selected(false)
			}
		if .setDisplayDetails('click', .selected = source.Name)
			.FindControl(.selected).Selected(true)
		return false
		}

	// TEMPORARY: 37414
	// Allows us to develop the suneido.js equivalent and test with select customers
	AllowBulletGraphs?()
		{
		return OptContribution('AllowBulletGraphs', { true })()
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}
