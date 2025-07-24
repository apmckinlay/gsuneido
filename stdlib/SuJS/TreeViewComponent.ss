// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Name: 		'TreeView'
	Xmin:		50
	Ymin:		100
	Xstretch:	1
	Ystretch:	1
	styles: `
		.su-treeview-container {
			position: relative;
			overflow: auto;
			outline-offset: -1px;
			outline: 1px solid black;
		}
		.su-treeview-container:focus {
			outline: none;
		}
		.su-treeview {
			position: absolute;
			top: 0px;
			left: 0px;
			width: 100%;
			height: 100%;
			table-layout: fixed;
			user-select: none;
			border: none;
			margin: 0px;
			border-spacing: 0px;
			padding: 5px;
			box-sizing: border-box;
			padding-inline-start: 0px;
			list-style-type: none;
		}
		.su-treeview-item {
			user-select: none;
			white-space: nowrap;
		}
		.su-treeview-item-selected>.su-treeview-label{
			background-color: lightgrey;
		}
		.su-treeview-subtree {
			padding-inline-start: 1em;
			list-style-type: none;
		}
		.su-treeview-item-folded .su-treeview-item {
			display: none;
		}
		.su-treeview-button:before {
			content: '\0229F';
		}
		.su-treeview-item-folded .su-treeview-button:before {
			content: '\0229E';
		}
		.su-treeview-button {
			display: inline-block;
			margin-right: 3px;
		}
		.su-treeview-edit {
			position: absolute;
			left: 0;
			top: -2px;
			outline: none;
			padding: 1px 2px;
			border: 1px black solid;
		}
		`

	New(@unused)
		{
		LoadCssStyles('treeview-control.css', .styles)
		.CreateElement('div', className: 'su-treeview-container')
		.root = CreateElement('ul', .El, className: 'su-treeview')
		.trees = Object(.root)
		.items = Object()
		}

	AddItem(parent, id, name, image, container?)
		{
		parentEl = .trees[parent]
		.items[id] = item = CreateElement('li', parentEl,
			className: 'su-treeview-item su-treeview-item-folded')
		if container? is true
			{
			button = CreateElement('span', item, className: 'su-treeview-button')
			button.AddEventListener('click', .buttonFactory(item, id))
			}
		textEl = CreateElement('span', item, className: 'su-treeview-label')
		textEl.innerText = name
		textEl.title = name
		textEl.AddEventListener('click', .labelFactory(id))
		if container? is true
			{
			subList = CreateElement('ul', item, className: 'su-treeview-subtree')
			.trees[id] = subList
			}
		.addImageEl(textEl, image)
		}

	addImageEl(el, image)
		{
		if image is false
			return
		imageEl = CreateElement('span', el, at: 0)
		imageEl.SetAttribute('translate', 'no')
		imageEl.textContent = image[0].char
		.SetStyles(Object(
			'font-family': image[0].font,
			'font-style': 'normal',
			'font-weight': 'normal',
			'margin-right': '5px',
			'user-select': 'none',
			'color': ToCssColor(image.GetDefault(1, #inherit))), imageEl)
		}

	SetImage(id, image)
		{
		if not .items.Member?(id)
			return
		textEl = .items[id].GetElementsByClassName('su-treeview-label').Item(0)
		if textEl.children.length > 0
			textEl.children.Item(0).Remove()
		.addImageEl(textEl, image)
		}

	buttonFactory(item, id)
		{
		return { |event/*unused*/|
			item.classList.Toggle('su-treeview-item-folded')
			.EventWithOverlay('TREEVIEW_TOGGLE',
				item.classList.Contains('su-treeview-item-folded'), id)
			}
		}

	labelFactory(id)
		{
		return { |event/*unused*/|
			if id isnt .selected
				{
				oldSelect = .selected
				.SelectItem(id)
				that = this
				.RunWhenNotFrozen({
					that.EventWithOverlay('TVN_SELCHANGED', oldSelect, id) })
				}
			}
		}

	selected: 0
	SelectItem(id)
		{
		if .selected is id
			return
		if .selected isnt 0 and .items.Member?(.selected)
			.items[.selected].classList.Remove('su-treeview-item-selected')
		if 0 isnt .selected = id
			{
			.items[.selected].classList.Add('su-treeview-item-selected')
			.ensureVisible(.items[.selected])
			}
		}

	ensureVisible(item)
		{
		parentTree = item.parentNode
		while parentTree isnt .root
			{
			parentTree.parentNode.classList.Remove('su-treeview-item-folded')
			parentTree = parentTree.parentNode.parentNode
			}

		rowHeight = item.offsetHeight
		rowOffsetTop = item.offsetTop
		scrollTop = .El.scrollTop
		scrollHeight = .El.clientHeight

		if scrollTop > rowOffsetTop
			.El.scrollTop = rowOffsetTop
		else if scrollTop + scrollHeight < rowOffsetTop + rowHeight
			.El.scrollTop = rowOffsetTop + rowHeight - scrollHeight
		}

	ExpandItem(id, collapse = false)
		{
		if collapse isnt true
			.items[id].classList.Remove('su-treeview-item-folded')
		else
			.items[id].classList.Add('su-treeview-item-folded')
		}

	ReorderChildren(parent, newOrders)
		{
		parentEl = .trees[parent]
		Assert(parentEl.children.length is: newOrders.Size())
		children = Object().AddMany!(0, newOrders.Size())
		for (i = 0; i < parentEl.children.length; i++)
			children[newOrders[i]] = parentEl.children.Item(i)

		for el in children
			parentEl.AppendChild(el)
		}

	DeleteItem(id, toDelete)
		{
		.items[id].Remove()
		.trees.Erase(@toDelete)
		.items.Erase(@toDelete)
		}

	edit: false
	editParent: false
	editId: false
	EditLabel(id)
		{
		if not .items.Member?(id)
			return
		el = .items[id]
		if .edit is false
			{
			.edit = CreateElement('input', el, className: 'su-treeview-edit')
			.edit.AddEventListener('blur', .onEditBlur)
			.edit.AddEventListener('keydown', .onEditKeydown)
			}
		el.AppendChild(.edit)
		el.SetStyle('position', 'relative')
		.editParent = el
		.editId = id
		.edit.value = el.GetElementsByClassName('su-treeview-label').Item(0).title
		.edit.Focus()
		.ensureVisible(el)
		}

	onEditBlur()
		{
		if .editId is false
			return
		id = .editId
		text = .edit.value
		.closeEdit()
		.Event('EndLabelEdit', id, text)
		}

	onEditKeydown(event)
		{
		if event.key is #Enter
			.onEditBlur()
		else if event.key is #Escape
			.closeEdit()
		}

	closeEdit()
		{
		if .editId is false
			return
		parent = .editParent
		.editId = .editParent = false
		.edit.Remove()
		parent.SetStyle('position', '')
		.ensureVisible(parent)
		}

	SetName(id, name)
		{
		if not .items.Member?(id)
			return
		labelEl = .items[id].GetElementsByClassName('su-treeview-label').Item(0)
		labelEl.title = name
		textNode = false
		for (i = 0; i < labelEl.childNodes.length; i++)
			{
			node = labelEl.childNodes.Item(i)
			if node.nodeType is 3 /*=Node.TEXT_NODE*/
				{
				textNode = node
				break
				}
			}
		if textNode isnt false
			textNode.nodeValue = name
		}
	}
