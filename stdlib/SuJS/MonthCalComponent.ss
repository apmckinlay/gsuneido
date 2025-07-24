// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: "MonthCal"
	Xstretch: 1
	Ystretch: 1
	styles: `
		.su-datepicker {
			display: inline-block;
			user-select: none;
			background-color: white;
			color: #454545;
		}
		.su-datepicker-prev-next {
			border: 1px solid lightgrey;
			border-radius: 3px;
			padding: 3px;
			font-size: 15pt;
			font-weight: bold;
			cursor: default;
		}
		.su-datepicker-prev-next:hover,
		.su-datepicker-prev-next:focus {
			cursor: pointer;
			border-color: grey;
			text-decoration: none;
		}
		.su-datepicker-prev-next:active {
			background-color: grey
		}
		.su-datepicker-prev {
			float: left;
			margin-left: 5px;
		}
		.su-datepicker-next {
			float: right;
			margin-right: 5px;
		}
		.su-datepicker-title {
			flex-grow: 1;
			text-align: center;
			padding: 3px;
			font-size: 15pt;
			font-weight: bold;
			cursor: default;
		}
		.su-datepicker-header {
			display: flex;
			align-items: baseline;
			justify-content: space-between;
			background-color: lightgrey;
			border: 1px solid darkgrey;
			border-radius: 3px;
		}
		.su-datepicker table {
			width: 100%;
		}
		.su-datepicker-day-td,
		.su-datepicker th {
			text-align: right;
			width: 2em;
		}
		.su-datepicker-td {
			border: 1px  solid #c5c5c5;
			background-color: #f6f6f6;
		}
		.su-datepicker-td:hover {
			background-color: lightblue;
			cursor: pointer;
		}
		.su-datepicker-grey-td {
			color: lightgrey;
		}
		.su-datepicker-selected-td {
			background-color: #cde8ff;
			border: 1px dotted black;
			cursor: pointer;
		}
		.su-datepicker-today-td {
			border: 1px solid #0066cc;
			color: #0066cc;
		}
		.su-datepicker-state-date-td {
			font-weight: bold;
		}
		.su-datepicker-month-td {
			text-align: center;
			padding: 1em;
		}
		.su-datepicker-year-td {
			text-align: center;
			padding: 0.5em 1em;
		}`
	displayState: 'date'
	New()
		{
		LoadCssStyles('su_datepicker.css', .styles)

		.CreateElement('div')
		.El.className = "su-datepicker"
		.El.SetAttribute('translate', 'no')
		.renderHeader()

		.SetMinSize()
		}

	Set(date)
		{
		if not Date?(date)
			date = Date()
		.selectedDate = date.NoTime()
		.switchState('date')
		}

	Get()
		{
		return .selectedDate
		}

	SELECT()
		{
		.Event('SELECT', .Get())
		}

	anchorMax: #29991201
	switchState(newState)
		{
		.displayState = newState

		.anchor = .selectedDate > .anchorMax
			? .anchorMax
			: Date(year: .selectedDate.Year(), month: .selectedDate.Month(), day: 1)

		switch (newState)
			{
		case 'date':
			.hideTable(.pickMonthTable, .pickYearTable)
			.showTable(.pickDateTable, .renderPickDateTable)
		case 'month':
			.hideTable(.pickDateTable, .pickYearTable)
			.showTable(.pickMonthTable, .renderPickMonthTable)
		case 'year':
			.hideTable(.pickDateTable, .pickMonthTable)
			.showTable(.pickYearTable, .renderPickYearTable)
			}
		}

	hideTable(@tables)
		{
		tables.Each()
			{ |table|
			if table isnt false
				table.SetStyle('display', 'none')
			}
		}

	showTable(table, tableRender)
		{
		if table is false
			tableRender()
		else
			table.SetStyle('display', 'table')
		.refresh()
		}

	renderHeader()
		{
		.header = CreateElement('div', .El, 'su-datepicker-header')
		.prev = CreateElement('div', .header,
			'su-datepicker-prev su-datepicker-prev-next')
		.prev.innerHTML = "<"
		.title = CreateElement('div', .header, 'su-datepicker-title')
		.next = CreateElement('div', .header,
			'su-datepicker-next su-datepicker-prev-next')
		.next.innerHTML = ">"

		.prev.AddEventListener('click', .On_Prev)
		.next.AddEventListener('click', .On_Next)
		.title.AddEventListener('click', .On_Title)
		}

	On_Prev()
		{
		.offsetAnchor(-1)
		.refresh()
		}

	On_Next()
		{
		.offsetAnchor(1)
		.refresh()
		}

	offsetAnchor(offset)
		{
		switch (.displayState)
			{
		case 'date':
			.anchor = .getValidMonth(.anchor, offset)
		case 'month':
			.anchor = .getValidYear(.anchor, offset)
		case 'year':
			.anchor = .getValidYear(.anchor, .yearsPerPage * offset)
			}
		}

	getValidMonth(date, offset)
		{
		return date > .anchorMax
			? .anchorMax
			: Min(.anchorMax, date.Plus(months: offset))
		}

	getValidYear(date, offset)
		{
		// Purposely stoping prior to Date.End to ensure we dont get 'bad date'
		if date > #28910101
			return #28910101
		newDate = date.Plus(years: offset)
		return newDate < #19000101
			? #19000101
			: newDate > #28910101
				? #28910101
				: newDate
		}

	refreshVersion: 0
	refresh()
		{
		++.refreshVersion
		switch (.displayState)
			{
		case 'date':
			.fillDates()
		case 'month':
			.fillMonths()
		case 'year':
			.fillYears()
			}
		.setTitle()
		}

	On_Title()
		{
		switch (.displayState)
			{
			case 'date':
				.switchState('month')
			case 'month':
				.switchState('year')
			default:
			}
		}

	pickDateTable: false
	renderPickDateTable()
		{
		.pickDateTable = CreateElement('table', .El)
		.renderTableHeader()
		.renderTableCells()
		}

	weekDays: #(Su, Mo, Tu, We, Th, Fr, Sa)
	weekDaySize: 7
	renderTableHeader()
		{
		.thead = CreateElement('thead', .pickDateTable)
		tr = CreateElement('tr', .thead)
		for weekday in .weekDays
			{
			head = CreateElement('th', tr)
			span = CreateElement('span', head)
			span.innerHTML = weekday
			}
		}

	dateCells: false
	maxDateRows: 6
	renderTableCells()
		{
		.dateCells = Object()
		tbody = CreateElement('tbody', .pickDateTable)
		for row in .. .maxDateRows
			{
			tr = CreateElement('tr', tbody)
			for col in .. .weekDaySize
				{
				td = CreateElement('td', tr)
				td.AddEventListener('click',
					.selectFactory(row, col, .weekDaySize, .selectDate))
				.dateCells.Add(td)
				}
			}
		}
	selectFactory(row, col, rowSize, fn)
		{
		return { fn(row * rowSize+ col) }
		}
	selectDate(index)
		{
		if .startDate is .anchorMax and index > 31 /*= jan 1 3000*/
			return
		.selectedDate = .startDate.Plus(days: index)
		.fillDates()
		.SELECT()
		}
	fillDates()
		{
		month = .anchor.Month()
		.initStartDate()
		today = Date().NoTime()

		months = [.startDate.Month()]
		hide = false
		curDate = .startDate
		for i in .. .dateCells.Size()
			{
			if curDate < Date.End()
				curDate = .startDate.Plus(days: i)
			else
				{
				.clearJan3000Days()
				break
				}
			isntCurMonth? = curDate.Month() isnt month
			if i % .weekDaySize is 0
				{
				hide = i isnt 0 and isntCurMonth?
				.dateCells[i].parentElement.SetStyle('display',
					hide ? 'none' : 'table-row')
				}
			if hide is true
				break
			if months.Last() isnt curDate.Month()
				months.Add(curDate.Month())
			.dateCells[i].innerHTML = curDate.Day()
			.dateCells[i].className = .dateClassName(curDate, today, isntCurMonth?)
			}
		.Event('GETDAYSTATE', .startDate, months.Size(), .refreshVersion)
		}
	initStartDate()
		{
		startDate = Min(.anchorMax, .anchor.Plus(days: -.anchor.WeekDay()))
		while startDate.Month() is .anchor.Month() and startDate.Day() isnt 1
			startDate = startDate.Plus(days: -7)
		.startDate = startDate
		}

	clearJan3000Days()
		{
		// expects last month displayed to be Dec 2999, these idx are Jan 2 / 3 / 4, 3000
		for invalidDate in #(32, 33, 34)
			.dateCells[invalidDate].innerHTML = .dateCells[invalidDate].className = ''
		}

	dateClassName(curDate, today, isntCurMonth?)
		{
		dateClassName = curDate is .selectedDate
			? "su-datepicker-td su-datepicker-day-td su-datepicker-selected-td"
			: isntCurMonth?
				? "su-datepicker-td su-datepicker-day-td su-datepicker-grey-td"
				: "su-datepicker-td su-datepicker-day-td"
		if today is curDate
			dateClassName $= ' su-datepicker-today-td'
		return dateClassName
		}

	SetDayState(states, refreshVersion)
		{
		if refreshVersion isnt .refreshVersion
			return

		cells = .pickDateTable.QuerySelectorAll('.su-datepicker-state-date-td')
		cells.ForEach(function (@args) {
			args[0].classList.Remove('su-datepicker-state-date-td') })

		for i in states.Members()
			{
			date = .startDate
			for ..i
				date = date.EndOfMonth().Plus(days: 1)
			day = 1
			while states[i] isnt 0
				{
				if ((states[i] & 1) is 1)
					{
					date = date.Replace(:day)
					n = date.MinusDays(.startDate)
					if .dateCells.Member?(n)
						.dateCells[n].classList.Add('su-datepicker-state-date-td')
					}
				day++
				states[i] >>= 1
				}
			}
		}

	pickMonthTable: false
	monthCells: false
	monthRow: 3
	monthsPerRow: 4
	months: #('Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug',
		'Sep', 'Oct', 'Nov', 'Dec')
	renderPickMonthTable()
		{
		.pickMonthTable = CreateElement('table', .El)
		.monthCells = Object()
		tbody = CreateElement('tbody', .pickMonthTable)
		for row in .. .monthRow
			{
			tr = CreateElement('tr', tbody)
			for col in .. .monthsPerRow
				{
				td = CreateElement('td', tr)
				td.innerHTML = .months[row * .monthsPerRow + col]
				td.AddEventListener('click',
					.selectFactory(row, col, .monthsPerRow, .selectMonth))
				.monthCells.Add(td)
				}
			}
		}
	selectMonth(index)
		{
		newDate = Date(year: .anchor.Year(), month: index + 1, day: 1)
		newDay = newDate.EndOfMonthDay() > .selectedDate.Day()
			? .selectedDate.Day()
			: newDate.EndOfMonthDay()
		.selectedDate = newDate.Plus(days: newDay - 1)
		.switchState('date')
		}
	fillMonths()
		{
		inSelectedYear? = .anchor.Year() is .selectedDate.Year()
		selectedMonth = .selectedDate.Month()
		for i in .. .monthRow * .monthsPerRow
			{
			.monthCells[i].className = inSelectedYear? and i + 1 is selectedMonth
				? 'su-datepicker-td su-datepicker-month-td su-datepicker-selected-td'
				: 'su-datepicker-td su-datepicker-month-td'
			}
		}

	pickYearTable: false
	yearCells: false
	yearRow: 4
	yearsPerRow: 4
	yearsPerPage: 16
	renderPickYearTable()
		{
		.pickYearTable = CreateElement('table', .El)
		.yearCells = Object()
		tbody = CreateElement('tbody', .pickYearTable)
		for row in .. .yearRow
			{
			tr = CreateElement('tr', tbody)
			for col in .. .yearsPerRow
				{
				td = CreateElement('td', tr)
				td.AddEventListener('click',
					.selectFactory(row, col, .yearsPerRow, .selectYear))
				.yearCells.Add(td)
				}
			}
		}
	selectYear(index)
		{
		newDate = Date(year: .startYear.Year() + index, month: .anchor.Month(), day: 1)
		newDay = newDate.EndOfMonthDay() > .selectedDate.Day()
			? .selectedDate.Day()
			: newDate.EndOfMonthDay()
		.selectedDate = newDate.Plus(days: newDay - 1)
		.switchState('month')
		}
	fillYears()
		{
		.initStartYear()
		for i in .. .yearCells.Size()
			{
			curYear = .startYear.Year() + i
			.yearCells[i].innerHTML = curYear
			.yearCells[i].className = curYear is .selectedDate.Year()
				? "su-datepicker-td su-datepicker-year-td su-datepicker-selected-td"
				: "su-datepicker-td su-datepicker-year-td"
			}
		}
	initStartYear()
		{
		.startYear = .getValidYear(.anchor, - 6 /*=offset to center anchor year*/)
		}

	setTitle()
		{
		switch (.displayState)
			{
		case 'date':
			.title.innerHTML = .anchor.FormatEn("MMM yyyy")
		case 'month':
			.title.innerHTML = .anchor.Year()
		case 'year':
			.title.innerHTML = .startYear.Year() $ " ~ " $
				(.startYear.Year() + .yearsPerPage - 1)
			}
		}

	// TODO:
	SetRange(@unused) {}
	}
