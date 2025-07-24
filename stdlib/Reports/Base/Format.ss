// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	X: 0
	Y: 0
	Xmin: 0
	Ymin: 0
	Xstretch: 0
	Ystretch: 0
	Export: true
	Generator?()
		{ return false; }
	GetSize(data /*unused*/ = false)
		{ return #(h: 0, w: 0, d: 0); }
	OnPage() // called when format is placed on page
		{ }
	Print(x /*unused*/, y /*unused*/, w /*unused*/, h /*unused*/, data /*unused*/ = false)
		{ }
	ExportCSV(data /*unused*/= '')
		{ return ''	}
	Variable?()
		{ return false; }
	Header?: false // used by InputFormat and Report
	// Query outputs
	Header()
		{ return false; }
	PageHeader()
		{ return .Header(); }
	twipFactor: 20 /* twentieth of a point */
	ChangeCoords(coords, pageHeight = 0, to = '')
		{
		Assert( to is 'PDF' or to is 'GDI' )

		if to is 'PDF'
			factor = 1 / .twipFactor
		else
			factor = .twipFactor
		// Change DPI:
		if ( Number?(coords) )
			return coords * factor
		for member in coords.Members()
			coords[member] *= factor

		// Change coordinate system:
		// (in PDF, (0,0) is the (bottom, left) corner of the page)
		if ( coords.Member?('y') )
			coords.y = pageHeight - coords.y

		return coords
		}
	CSVExportLine(str)
		{
		if str is '' or str.Tr(',') is '' or .Export is false
			return ''
		return '\n' $ str[.. -1]
		}
	CSVExportString(str)
		{
		if .Export isnt true
			return ''
		// can't use Display, or else it will escape inner quotes
		return '"' $ String(str).Tr('\r\n', ' ').Replace('"', '""') $ '"'
		}
	Access: false
	InitAccessField(access)
		{
		if access isnt false
			.Access = Object("AccessGoTo", access: access[0],
				goto_field: access[1],
				goto_value: access.GetDefault(2, access[1]))
		}
	Hotspot(x, y, w, h, data, access = false)
		{
		if not .PreviewMode?()
			return

		if this.Member?("Field") and data[.Field] isnt ""
			{
			if data.Member?('_access_' $ .Field) and
				data['_access_' $ .Field] is 'disable'
				return

			.addHotspotByField(data, x, y, w, h)
			}
		else if access isnt false
			.addHotspotPoint(x, y, w, h, access)
		}
	addHotspotByField(data, x, y, w, h)
		{
		if not data.Member?('_access_' $ .Field) and
			.Access isnt false and data.Member?(.Access.goto_value)
			QueryFormat.AddAccessField(data, .Field, .Access)

		if data.Member?('_access_' $ .Field)
			.addHotspotPoint(x, y, w, h, data['_access_' $ .Field])
		}
	unit: 18 // need figure out which code did the unit covert for meta file and screen.
	addHotspotPoint(x, y, w, h, access)
		{
		page = _report.GetGeneratedPages() - 1
		points = _report.AccessPoints
		if not points.Member?(page)
			points[page] = Object()
		left = x / .unit
		top = y / .unit
		points[page].Add(Object(:left, :top, right: left + w / .unit,
			bottom: top + h / .unit, :access))
		}

	// method for adding drill down information into data
	DisableAccessField(data, field)
		{
		data['_access_' $ field] = 'disable'
		}
	AddAccessField(data, field, access)
		{
		if Object?(access) and access[0] is 'AccessGoTo'
			.AddAccessPoint(data, field, access.access,	access.goto_field,
				data[access.goto_value])
		else if String?(access)
			{
			dd = Datadict(access)
			if dd.Control[0] is 'Id'
				.AddAccessPoint(data, field, dd.Control.access,	access, data[access])
			}
		}
	AddAccessPoint(data, field, access, goto_field, goto_value)
		{
		data_access = data['_access_' $ field] = Object()
		data_access.control = 'AccessGoTo'

		if Function?(access)
			access = access(goto_value)

		data_access.access = access
		data_access.goto_field = goto_field
		data_access.goto_value = goto_value
		}
	AddDynamicAccessPoint(data, field, func)
		{
		data_access = data['_access_' $ field] = Object()
		data_access.control = 'DynamicGoTo'
		data_access.data = data
		data_access.func = func
		}
	AddAccessReportPoint(data, field, rpt, set_params)
		{
		data_access = data['_access_' $ field] = Object()
		data_access.control = 'ReportGoTo'
		data_access.report = rpt
		data_access.params = set_params
		}
	EqualParam(value)
		{
		return Object(operation: "equals", :value, value2: "")
		}
	PreviewMode?()
		{
		return _report.Member?('Params') and _report.Params.ReportDestination is 'preview'
		}
	DoWithFont(origFont, block)
		{
		oldfont = _report.SelectFont(origFont)
		font = origFont is false ? _report.GetFont() : origFont
		font = _report.EnsureFont(font, oldfont)
		block(font)
		_report.SelectFont(oldfont)
		}
	}
