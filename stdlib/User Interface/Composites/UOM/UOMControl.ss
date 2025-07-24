// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'UOM'
	New(numctrl, uomctrl, readonly = false, div = '', .uom_optional = false,
		tabover = false, .mandatory = false, .hidden = false)
		{
		super(.layout(numctrl, uomctrl, div, tabover))
		.value = .Horz.Value
		.uom = .Horz.Uom
		.Left = .Horz.Value.Left
		.Top = .Horz.Top
		// initial SetReadOnly call has to happen before .readonly is set because the
		// SetReadOnly method checks .readonly and doesn't do anything if true
		.SetReadOnly(readonly)
		.readonly = readonly
		.Send('Data')
		}
	layout(numctrl, uomctrl, div, tabover)
		{
		numctrl = Object?(numctrl) ? numctrl.Copy() : Object(numctrl)
		numctrl.name = 'Value'
		numctrl.tabover = tabover
		numctrl.mandatory = .mandatory
		numctrl.hidden = .hidden
		uomctrl = Object?(uomctrl) ? uomctrl.Copy() : Object(uomctrl)
		uomctrl.name = 'Uom'
		uomctrl.tabover = tabover
		uomctrl.hidden = .hidden
		ctrl = Object('Horz', numctrl, uomctrl, overlap:)
		if div isnt '' and not .hidden
			{
			skip = #(Skip medium:)
			ctrl.Add(skip, Object('Static', div, weight: 'bold'), skip, at: 2)
			}
		return ctrl
		}
	NewValue(value/*unused*/)
		{
		.flat_amt = false
		.Send('NewValue', .Get())
		}
	// flat_amt member is used for when flat rates are set (will be a number)
	flat_amt: false
	Set(x)
		{
		if Number?(x)
			{
			.flat_amt = x
			.value.Set('')
			.uom.Set('')
			return
			}
		ob = Split_UOM(x)
		.value.Set(ob.value)
		.uom.Set(ob.uom)
		}
	Get()
		{
		val = .value.Get()
		uom = .uom.Get()

		if val is '' and uom is '' and .flat_amt isnt false
			return .flat_amt
		return val $ Opt(' ', uom)
		}
	Dirty?(dirty = "")
		{
		return .value.Dirty?(dirty) or .uom.Dirty?(dirty)
		}

	readonly: false
	SetReadOnly(on = true)
		{
		if .readonly
			return

		.value.SetReadOnly(on)
		.uom.SetReadOnly(on)
		}
	HandleTab()
		{
		if .value.HasFocus?()
			{
			.uom.SetFocus()
			return true
			}
		return false
		}

	Valid?()
		{
		if .flat_amt isnt false
			return true
		if not .value.Valid?() or not .uom.Valid?()
			return false
		data = .Get().Trim()
		return .validCheck?(data, .mandatory, .uom_optional)
		}
	validCheck?(data, mandatory, uom_optional)
		{
		if data is ''
			return true

		if mandatory is true and .ConsideredEmpty?(data)
			return false

		ob = Split_UOM(data)
		ob.value = String(ob.value)
		if not uom_optional and (ob.value.Blank?() or ob.uom.Blank?())
			return false
		return true
		}

	ValidData?(@args)
		{
		value = args[0]
		uom_optional = args.GetDefault('uom_optional', false)

		if Number?(value) and uom_optional
			return true

		ob = Split_UOM(value)
		ob.value = String(ob.value)
		if ob.value isnt "" and not ob.value.Number?()
			return false
		return .validCheck?(value, args.GetDefault('mandatory', false), uom_optional)
		}

	ConsideredEmpty?(data)
		{
		value = String(Split_UOM(data).value)
		return value.Trim() is "0" or value.Blank?()
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}