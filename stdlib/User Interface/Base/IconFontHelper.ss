// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Fonts: (
		'suneido.ttf': #(
			"arrow_down", "arrow_end", "arrow_home", "arrow_up", "back", "forward",
			"previous", "next", "cross", "plus", "minus", "info", "questionMark_black",
			"triangle-warning", "left", "right", "up", "down", "undo", "redo", "bold",
			"italic", "heading", "underline", "text", "code", "UNUSED_1", "UNUSED_2",
			"UNUSED_3", "UNUSED_4", "UNUSED_5", "UNUSED_6", "UNUSED_7", "UNUSED_8",
			"UNUSED_9", "UNUSED_10", "UNUSED_11", "UNUSED_12", "UNUSED_13", "UNUSED_14",
			"UNUSED_15", "UNUSED_16", "copy",
			"cut", "paste", "save", "delete_item", "edit", "checkmark", "new_folder",
			"parent_folder", "commentline", "commentspan", "custom_screen", "delete",
			"expand", "collapse", "filter", "find", "find_in_folders", "find_next",
			"find_previous", "flag", "folder", "open_folder", "document", "notes", "check"
			"hsplit", "vsplit", "link", "list", "print", "UNUSED_17", "UNUSED_18", "zoom",
			"starRatingImageEmpty", "starRatingImageFull", "starRatingImageHalf",
			"refresh", "previous_flag", "next_flag", "im", "view_form", "view_list",
			"locked", "unlocked", "invalid_lock", "valid_lock", "square", "restore",
			"menu"),
		'suneido2.ttf': #(
			"type", "select", "rounded-rectangle", "rectangle", "line", "image",
			"ellipse", "curve", "location", "close"))


	startCode: 0x20
	InitMap(fontName)
		{
		map = Object()
		i = .startCode
		for icon in .Fonts[fontName]
			map[icon] = i++
		return map
		}

	defaultFont: 0x20
	GetCode(image)
		{
		i = false
		image = image.RemoveSuffix('.emf')
		for fontName in .Fonts.Members()
			if false isnt i = .Fonts[fontName].Find(image)
				break
		return i is false ? .defaultFont : .startCode + i
		}

	GetFontStyles(buildUrlFn)
		{
		styles = Object()
		for fontName in .Fonts.Members()
			{
			styles.Add(`@font-face { ` $
				`font-family: "` $ fontName.RemoveSuffix('.ttf') $ `"; ` $
				`src: url("` $ buildUrlFn(fontName) $ `") format("truetype"); }`)
			}
		return styles
		}
	}