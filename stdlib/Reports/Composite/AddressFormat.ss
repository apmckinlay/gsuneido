// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// WARNING: must pass data to New
VertFormat
	{
	New(data, prefix = "", suffix = "", justify = 'left', font = false,  w = false)
		{
		super(@.layout(data, prefix, suffix, justify, font, w))
		}
	csv: ''
	layout(data, prefix, suffix, justify, font, w)
		{
		wrap = justify is 'center'
			? function (item) { Object('Horz', 'Hfill', item, 'Hfill') }
			: function (item) { item }

		layout = Object('Vert')
		addr1 = data[prefix $ 'address1' $ suffix]
		addr2 = data[prefix $ 'address2' $ suffix]
		if (addr1 isnt "")
			layout.Add(wrap(Object('Text', addr1, :font, :w)))
		if (addr2 isnt "")
			layout.Add(wrap(Object('Text', addr2, :font, :w)))

		city = data[prefix $ 'city' $ suffix]
		state_prov = data[prefix $ 'state_prov' $ suffix]
		zip_postal = data[prefix $ 'zip_postal' $ suffix]
		s = Join('  ', Join(' ', city, state_prov), zip_postal)
		if (s isnt "")
			layout.Add(wrap(Object('Text', s, :font, :w)))
		.csv = Join(', ', addr1, addr2, s)

		return layout
		}
	ExportCSV(data /*unused*/= '')
		{
		return .CSVExportString(.csv)
		}
	}
