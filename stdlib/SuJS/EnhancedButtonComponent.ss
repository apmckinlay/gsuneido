// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Component
	{
	styles: `
		.su-enhanced-button {
			padding: 0;
			margin: 0px;
			border-width: 1px;
			border-style: solid;
			border-radius: 0.3em;
			background-color: lightgrey;
			user-select: none;
		}
		.su-enhanced-button-no-button-style {
			padding: 0;
			border: 1px solid transparent;
			border-radius: 0.3em;
			background-color: initial;
			user-select: none;
		}
		.su-enhanced-button-mouse-effect:hover {
			background-color: azure;
			border-color: deepskyblue;
		}
		.su-enhanced-button-mouse-effect:active,
		.su-enhanced-button-pushed {
			background-color: lightblue;
			border-color: deepskyblue;
		}
		.su-enhanced-button-mouse-effect[data-highlight=true] {
			border-color: deepskyblue;
		}
		.su-enhanced-button-mouse-effect:focus {
			outline-offset: -3px;
			outline: 1px dashed black;
		}`
	ContextMenu: true
	New(text = false,
		tabover = false, .defaultButton = false,
		tip = false, .pad = false,
		font = "", size = "", weight = "", textColor = false,
		.width = false, .buttonWidth = false, .buttonHeight = false,
		italic = false, underline = false, strikeout = false,
		image = false, mouseOverImage = false, mouseDownImage = false,
		imageColor = false, mouseOverImageColor = false, .book = 'imagebook',
		.mouseEffect = false, .imagePadding = 0, .buttonStyle = false,
		.enlargeOnHover = false)
		{
		LoadCssStyles('su-enhanced-button.css', .styles)
		.CreateElement('button')
		.initCssClass()
		.xmin_orig = .Xmin
		.textEl = CreateTextNode('', .El)
		.SetText(text)
		.SetTextColor(textColor)
		.SetImageColor(imageColor, mouseOverImageColor)
		.SetImage(image, mouseOverImage, mouseDownImage)
		.SetFont(font, size, weight, underline, italic, strikeout)
		.Recalc()
		if .buttonStyle is false
			.El.classList.Add('su-no-button-style')
		if tabover is true or .buttonStyle is false
			.El.tabIndex = "-1"
		if defaultButton is true
			.El.Focus()

		.El.AddEventListener('click', .CLICKED)
		.El.AddEventListener('mouseenter', .onMouseenter)
		.El.AddEventListener('mouseleave', .onMouseleave)
		.El.AddEventListener('mousedown', .onMousedown)
		.El.AddEventListener('focus', .focus)
		.El.AddEventListener('blur', .blur)

		.AddToolTip(tip)
		}

	initCssClass()
		{
		classList = .El.classList
		classList.Add(.buttonStyle
			? 'su-enhanced-button'
			: 'su-enhanced-button-no-button-style')
		if .mouseEffect is true
			classList.Add('su-enhanced-button-mouse-effect')
		}

	SetMouseEffect(.mouseEffect)
		{
		if .mouseEffect is true
			.El.classList.Add('su-enhanced-button-mouse-effect')
		else
			.El.classList.Remove('su-enhanced-button-mouse-effect')
		}

	Recalc()
		{
		text = .text is false ? 'M' : .text
		metrics = SuRender().GetTextMetrics(.El, text)
		.calcHeight(metrics)
		tw = .text is false ? 0 : metrics.width + .pad
		w = .buttonWidth isnt false
			? .buttonWidth
			: .width isnt false
				? SuRender().GetTextMetrics(.El, 'M'.Repeat(.width)).width
				: 0
		h = .buttonHeight isnt false ? .buttonHeight : .Ymin
		.Ymin = h
		.recalcImages()
		iw = .imageChar isnt false ? h : 0
		.Xmin = Max(Max(iw + 5/*=gap*/ + tw, w), .xmin_orig)
		.SetMinSize()
		}

	calcHeight(metrics)
		{
		.Ymin = metrics.height
		if .pad is false
			.pad = .Ymin + .Ymin % 2

		.Ymin += .mouseEffect is false ? 8 : 10 /* = y padding + border*/
		.descent = metrics.descent + 5/*=padding*/
		}

	oldYmin: false
	recalcImages()
		{
		.oldYmin = .Ymin
		.recalcImageSize(.imageEl)
		.recalcImageSize(.mouseOverImageEl)
		.recalcImageSize(.mouseDownImageEl)
		}

	SetMinSize()
		{
		if .enlargeOnHover isnt false
			return
		super.SetMinSize()
		}

	Highlight(highlight?)
		{
		if .El is false
			return
		.El.dataset.highlight = highlight? is true
		}

	focus(event)
		{
		.Window.HighlightDefaultButton(false)

		el = false
		try
			el = event.relatedTarget
		if false isnt control = .GetControlFromEl(el)
			.Event('SyncPrevFocus', control.UniqueId)
		}

	blur()
		{
		.Window.HighlightDefaultButton(true)
		}

	imageEl: false
	imageChar: false
	mouseOverImageEl: false
	mouseDownImageEl: false
	svgImages: false
	SetImage(image, mouseOverImage = false, mouseDownImage = false)
		{
		.clear()
		if image isnt false
			{
			.imageEl = .mouseOverImageEl = .mouseDownImageEl = .createEl(image)
			.imageChar = .isTwoImages?(image)
				? image[0].char $ image[1].char
				: .isFontImage?(image)
					? image.char
					: false
			}
		if mouseOverImage isnt false
			.mouseOverImageEl = .mouseDownImageEl = .createEl(mouseOverImage)
		if mouseDownImage isnt false
			.mouseDownImageEl = .createEl(mouseDownImage)
		.paint()
		.recalcImages()
		}

	clear()
		{
		.El.innerHTML = ""
		if .textEl isnt false
			AttachElement(.textEl, .El, false)
		.imageEl = .imageChar = .mouseOverImageEl = .mouseDownImageEl = false
		.svgImages = Object()
		}

	createEl(image)
		{
		el = false
		if .isTwoImages?(image)
			el = .createTwoImagesEl(image)
		else if .isFontImage?(image)
			el = .createSpan(image.char, image.font, .El)
		else
			el = .createImg(image)

		.styleImage(el)
		return el
		}

	isTwoImages?(image)
		{
		return Object?(image) and image.Members(named:).Empty?()
		}

	isFontImage?(image)
		{
		return Object?(image) and image.Member?(#font)
		}

	createTwoImagesEl(image)
		{
		span = .createSpan('', 'suneido', .El,
			Object('letter-spacing': image[3/*=gap*/] $ 'px'))
		highlightStyle = #('color': 'blue')
		.createSpan(image[0].char, image[0].font, span,
			image[2/*=highlighted*/] is 0 ? highlightStyle : #())
		.createSpan(image[1].char, image[1].font, span,
			image[2/*=highlighted*/] is 1 ? highlightStyle : #())
		return span
		}

	 createSpan(ch, font, parent, styles = #())
		{
		el = CreateElement('span', parent, at: 0)
		el.SetAttribute('translate', 'no')
		if ch is ' '
			el.innerHTML = '&nbsp;'
		else
			el.textContent = ch
		.SetStyles(Object('font-family': font,
			'font-style': 'normal',
			'font-weight': 'normal'), el)
		.SetStyles(styles, el)
		return el
		}

	createImg(image)
		{
		img = CreateElement('img', .El)
		img.src = 'data:image/svg+xml;base64,' $ Base64.Encode(image)
		.SetStyles(Object('font-family': 'suneido',
			'font-style': 'normal',
			'font-weight': 'normal'), img)
		return img
		}

	descent: false
	recalcImageSize(imageEl)
		{
		if imageEl is false
			return

		imagePadding = Object?(.imagePadding) ? .imagePadding.h : .imagePadding
		size = ((.Ymin - 2) * (1 - imagePadding)).Floor()
		if imageEl.tagName is 'IMG'
			imageEl.SetStyle('height', size $ 'px')
		imageEl.SetStyle('font-size', size $ 'px')
		if .descent isnt false and .buttonHeight is false
			{
			metrics = imageEl.tagName is 'IMG'
				? Object(height: size, descent: 0)
				: SuRender().GetTextMetrics(imageEl, 'M')
			// align images to text's baseline
			imageEl.SetStyle('vertical-align',
				(((.Ymin - metrics.height + 1) >> 1) + metrics.descent - .descent) $ 'px')
			}
		}

	styleImage(el)
		{
		if el is false
			return
		if .text not in (false, '')
			{
			el.SetStyle('padding-right', '5px')
			}
		else
			{
			el.SetStyle('padding-right', 'initial')
			}
		el.draggable = false
		}

	mouseover?: false
	onMouseenter()
		{
		if .El is false
			return
		.mouseover? = true
		if .enlargeOnHover isnt false
			{
			super.SetMinSize()
			.El.classList.Add('su-no-button-style', 'su-enhanced-button-mouse-effect')
			}
		.paint()
		}
	onMouseleave()
		{
		if .El is false
			return
		.mouseover? = false
		if .enlargeOnHover isnt false
			{
			.El.SetStyle('min-width', '')
			.El.SetStyle('min-height', '')
			.El.classList.remove('su-no-button-style', 'su-enhanced-button-mouse-effect')
			}
		.paint()
		}
	mousedown?: false
	onMousedown(event)
		{
		if event.button is 2
			{
			r = SuRender.GetClientRect(.El)
			.RunWhenNotFrozen({ .Event('RBUTTONDOWN', r) })
			return
			}
		if event.button isnt 0
			return
		.mousedown? = true
		.StartMouseTracking(.onMouseup)
		.paint()
		}
	onMouseup(@unused)
		{
		.mousedown? = false
		.StopMouseTracking()
		.paint()
		}

	paint()
		{
		.drawImage()
		}

	drawImage()
		{
		showEl = .pushed or .mousedown?
			? .mouseDownImageEl
			: .mouseover?
				? .mouseOverImageEl
				: .imageEl

		if showEl is false
			return
		.imageEl.SetStyle(#display, .imageEl is showEl ? 'initial' : 'none')
		.mouseOverImageEl.SetStyle(#display,
			.mouseOverImageEl is showEl ? 'initial' : 'none')
		.mouseDownImageEl.SetStyle(#display,
			.mouseDownImageEl is showEl ? 'initial' : 'none')
		.showColor()
		}

	showColor()
		{
		if .imageEl is false
			return
		color = .mouseover? or .mousedown? or .pushed ? .mouseOverImageColor : .imageColor
		if color is false
			color = 'black'

		color = ToCssColor(color)
		if .imageEl.tagName is 'IMG'
			.setSvgImageColor(color, .imageEl)
		else
			.imageEl.SetStyle(#color, color)
		.mouseOverImageEl.SetStyle(#color, color)
		.mouseDownImageEl.SetStyle(#color, color)
		}

	setSvgImageColor(color, imageEl)
		{
		if color is 'black' // unset
			return

		if .svgImages isnt false and .svgImages.Member?(color)
			imageEl.src = .svgImages[color]
		else
			{
			image = Base64.Decode(imageEl.src.RemovePrefix('data:image/svg+xml;base64,'))
			image = image.Replace('fill:#[\da-f]+', 'fill:' $ color)
			.svgImages[color] = imageEl.src = 'data:image/svg+xml;base64,' $
				Base64.Encode(image)
			}
		}

	imageColor: false
	mouseOverImageColor: false
	SetImageColor(.imageColor = false, .mouseOverImageColor = false)
		{
		.showColor()
		}

	textEl: false
	SetText(.text)
		{
		if String?(.text)
			.text = ButtonComponent.RemoveAmpersand(.text)
		.textEl.nodeValue = .text isnt false ? .text : ''

		.styleImage(.imageEl)
		.styleImage(.mouseOverImageEl)
		.styleImage(.mouseDownImageEl)
		.WindowRefresh()
		}

	SetTextColor(.textColor)
		{
		if .textEl is false
			return
		.El.SetStyle('color', .textColor is false ? 'black' : ToCssColor(.textColor))
		}

	CLICKED()
		{
		.RunWhenNotFrozen({ .EventWithOverlay('CLICKED') })
		}

	pushed: false
	Pushed?(state = -1)
		{
		if state isnt -1 and state isnt .pushed
			{
			.pushed = state
			if .pushed is true
				.El.classList.Add('su-enhanced-button-pushed')
			else
				.El.classList.Remove('su-enhanced-button-pushed')
			.paint()
			}
		return .pushed
		}
	}
