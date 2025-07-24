// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	disabled: false
	label: false
	New(.text = "", lefttext = false, .font = "", .size = "", .weight = "",
		.readonly = false, tip = "", hidden = false, tabover = false, .align = 'center')
		{
		if .text isnt ""
			{
			.CreateElement('div')
			.input = CreateElement('input', .El)
			id = .text $ Suneido.GetInit(#SuneidoJS_UI_Num, 0)
			Suneido.SuneidoJS_UI_Num++
			.input.id = id
			.label = CreateElement("label", .El, at: lefttext ? 0 : 1)
			.label.innerHTML = .text
			.label.htmlFor = id
			}
		else
			{
			.CreateElement('input')
			.input = .El
			.input.Control(this)
			.input.Window(.Window)
			.input.SetStyle('align-self', .align)
			}
		.input.type = .InputType
		.setMargin(lefttext)
		.SetHidden(hidden)
		if tabover isnt false
			.input.tabIndex = "-1"
		if size isnt ""
			{
			size = .ConvertSize(size)
			.input.SetStyle('height', size)
			.input.SetStyle('width', size)
			}
		.updateStyle()
		.Recalc()
		.input.AddEventListener('click', .Toggle)
		.AddToolTip(tip)
		}

	Resize(w, h)
		{
		super.Resize(w, h)
		// stretch the input so that it will cover the whole listedit area and capture
		// the focus and click events
		.input.SetStyle('width', w $ 'px')
		.input.SetStyle('height', h * .6/*=smaller*/ $ 'px')
		.input.SetStyle('vertical-align', 'text-bottom')
		}

	setMargin(lefttext)
		{
		.input.SetStyle('margin-top', '0')
		.input.SetStyle('margin-bottom', '0')
		.input.SetStyle(lefttext ? 'margin-left' : 'margin-right', '5px')
		.input.SetStyle(lefttext ? 'margin-right' : 'margin-left', '0')
		}

	Recalc()
		{
		.SetFont(.font, .size, .weight)
		metrics = SuRender().GetTextMetrics(.input, 'M')
		.Xmin = metrics.width + 5/*=gap*/
		.Ymin = metrics.height
		if .label isnt false
			{
			metrics = SuRender().GetTextMetrics(.label, .text)
			.Xmin += metrics.width
			.Ymin = Max(metrics.height, .Ymin + metrics.height - metrics.ascent)
			.SetMinSize()
			}
		}

	value: ""
	Toggle(event)
		{
		if .readonly or .set_readonly or .disabled
			{
			event.PreventDefault()
			return
			}
		.RunWhenNotFrozen(.DoToggle)
		}

	DoToggle()
		{
		.value = .value isnt true
		.EventWithOverlay('Toggle')
		.input.checked = .value is true
		}

	readonly: false
	set_readonly: false
	GetReadOnly()
		{
		return .readonly or .set_readonly
		}
	SetReadOnly(ro)
		{
		.set_readonly = ro
		.updateStyle()
		}
	SetEnabled(enabled)
		{
		.disabled = not enabled
		.updateStyle()
		}

	styleDisabled?: false
	updateStyle()
		{
		disabled? = .readonly or .set_readonly or .disabled
		if disabled? isnt .styleDisabled?
			.input.SetStyle('opacity', .styleDisabled? = disabled? ? '0.5' : '')
		}

	SetFocus()
		{
		.input.Focus()
		}
	ClearFocus()
		{
		.input.Blur()
		}
	HasFocus?()
		{
		return SuUI.GetCurrentDocument().activeElement is .input
		}

	Get()
		{
		return .value is true
		}
	Set(value)
		{
		if value is .value
			return

		.value = value
		.input.checked = value is true
		}

	SetColor(color)
		{
		.El.SetStyle('color', ToCssColor(color))
		}

	GetText()
		{
		return .text
		}

	// for HighlightControl
	FindControl(name)
		{
		return name is 'Static' ? this : false
		}
	}
