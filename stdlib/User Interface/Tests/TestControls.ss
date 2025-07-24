// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "TestControls"
	Controls: (Record (Vert
		zip_postal
		phone
		date
		(Pair (Static ChooseDates) (ChooseDates))
		(Pair (Static Key)(Key stdlib name columns: (name) prefixColumn: 'name'))
		(Pair (Static ChooseMany)(ChooseMany #(one two three)))
		(Pair (Static Number)(Number name: Number) name: 'Pair1')
		(Pair (Static Hours)(Hours))
		dollar
		(Pair (Static Field) Field)
		(Pair (Static Readonly) (Field readonly:, name: Readonly))
		email
		web
		(Pair (Static Editor) (Editor))
		(Pair (Static ParamsChooseList) (ParamsChooseList, date, values: #(),
			readonly: false))
		(UOM (Number mask: '###,###.##') Field, name: uom)
		(Pair (Static 'Number without mask') (Number mask: false, name: nomask))
		(Info, name: 'test_info')
		(Override override_value override
			valid: function (value)
				{
				return Number(value) < 100 /*= override max*/
					? ''
					: 'must be < 100'
				})
		(Pair
			(Static ChooseList)
			(ChooseList #('one - 1' 'two - 2' 'THREE - 3' 'Four - 4' 'Five',
				'six', 'seven', eight, nine), width: 10))
		(Pair
			(Static Spinner)
			(Spinner 1, 4))
		(Pair
			(Static FieldHistory)
			(FieldHistory, name: 'fieldhistory'))
		(Pair
			(Static listField)
			(ChooseList listField: Field, width: 10))
		state_prov
		(Pair
			(Static Option)
			(RadioButtons one two three name: radio))
		FirstLastNameControl
		(Pair (Static Encrypt) EncryptControl)
		(Pair (Static 'Mandatory ScintillaAddonsEditor')
			(ScintillaAddonsEditor mandatory:))
		Skip
		Skip
		(Horz (Button Inspect) (Button Set) Skip (Button Get)
			Skip (Static 'F9 = Dirty?') Skip (Static 'F4 = Inspect'))
		Skip
		Statusbar
		))
	Commands: ((Dirty, F9)(Inspect, F4))
	Status(status)
		{
		.Data.Vert.Status.Set(status)
		}
	On_Inspect()
		{
		Inspect(.Data.Get())
		}
	On_Set()
		{
		.Data.Set(Record(
			state_prov: 'BC',
			Number: 345,
			dollar: 125.58,
			uom: '123456 lb',
			nomask: 1.2365478545,
			radio: 'two',
			FirstLastName: 'Jones, Joe'
			Encrypt: "a string".Xor(EncryptControlKey())
			))
		}
	On_Get()
		{
		Alert(.Data.Vert.Pair1.Number.Get(), title: 'Get',
			flags: MB.ICONINFORMATION)
		}
	On_Dirty()
		{
		Alert(.Data.Dirty?() ? "Dirty!" : "Clean", title: 'Dirty',
			flags: MB.ICONINFORMATION)
		}
	}
