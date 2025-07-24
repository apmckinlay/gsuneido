// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	try Database("destroy crosstabtest")
	Database("create crosstabtest (city, person, amount) key(city, person, amount)")
	data = #(
		(city: "Saskatoon", person: "Jim", amount: 456)
		(city: "Vancouver", person: "Jim", amount: 100)
		(city: "Vancouver", person: "Jim", amount: 30)
		(city: "Calgary", person: "Jim", amount: 40)
		(city: "Calgary", person: "Jim", amount: 50)
		(city: "Regina", person: "Jake", amount: 60)
		(city: "Winnipeg", person: "Jake", amount: 70)
		(city: "Winnipeg", person: "Jake", amount: 80)
		(city: "Calgary", person: "Jake", amount: 90)
		(city: "Vancouver", person: "Jake", amount: 100)
		(city: "Vancouver", person: "Jerome", amount: 20)
		(city: "Vancouver", person: "Jerome", amount: 30)
		(city: "Regina", person: "Jerome", amount: 40)
		(city: "Calgary", person: "Jerome", amount: 50)
		(city: "Saskatoon", person: "Jerome", amount: 45)
		(city: "Saskatoon", person: "Jake", amount: 6)
		(city: "Regina", person: "Jerome", amount: 80)
		(city: "Vancouver", person: "Jake", amount: 10)
		)
	for (x in data)
		QueryOutput('crosstabtest', x)

	previewWindow = GetFocus()
	// 0d count
	Params.On_Preview(#(Crosstab crosstabtest func: count), :previewWindow)
	// 0d total
	Params.On_Preview(#(Crosstab crosstabtest func: total value: amount), :previewWindow)
	// 1d count
	Params.On_Preview(#(Crosstab crosstabtest rows: person func: count), :previewWindow)
	// 1d total
	Params.On_Preview(#(Crosstab crosstabtest rows: person func: total value: amount),
		:previewWindow)
	// 2d count
	Params.On_Preview(#(Crosstab crosstabtest rows: city cols: person func: count),
		:previewWindow)
	// 2d total
	Params.On_Preview(#(Crosstab crosstabtest rows: city cols: person func: total
		value: amount, sortcols: #(Jim, Jake, Jerome)),	:previewWindow)

	try Database("destroy crosstabtest")
	}
