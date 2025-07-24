// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	FolderClosed:				0
	FolderClosedDark:			1
	FolderClosedModified:		2
	FolderClosedModifiedDark:	3

	FolderOpen:					4
	FolderOpenDark:				5
	FolderOpenModified:			6
	FolderOpenModifiedDark:		7

	Document:					8
	DocumentDark:				9

	CloseButton:				10
	CloseButtonDark:			11

	Modified:					12
	ModifiedDark:				13

	Invalid:					14
	InvalidDark:				15
	New()
		{
		.init()
		}

	init()
		{
		imageMap = .imageMap()
		.initImageResources(imageMap)
		.initImageList(imageMap)
		.invalidImages = [.Invalid, .InvalidDark]
		.modifiedImages = [.Modified, .ModifiedDark, .FolderClosedModified,
			.FolderClosedModifiedDark, .FolderOpenModified, .FolderOpenModifiedDark]
		}

	imageMap()
		{
		return [
			// Closed folders
			[image: 'folder.emf', 	color: CLR.black],
			[image: 'folder.emf', 	color: CLR.darkorange],
			// Open folders
			[image: 'open_folder.emf', color: CLR.black],
			[image: 'open_folder.emf', color: CLR.orange],
			// Documents / Close
			[image: 'document.emf', color: CLR.black],
			[image: 'delete.emf',   color: CLR.black, padding: 0.2],
			[image: 'document.emf', color: CLR.orange],
			[image: 'document.emf', color: CLR.red],
			]
		}

	initImageResources(imageMap)
		{
		.ImageResources = Object()
		imageMap.Each()
			{
			if not Sys.SuneidoJs?()
				{
				.ImageResources.Add([ImageResource(it.image), it.color])
				.ImageResources.Add([ImageResource(it.image), it.color])
				}
			else
				{
				codeOb = IconFont().MapToCharCode(it.image)
				.ImageResources.Add([char: codeOb.char, font: codeOb.font,
					color: it.color])
				.ImageResources.Add([char: codeOb.char, font: codeOb.font,
					color: it.color])
				}
			}
		}

	initImageList(imageMap)
		{
		x = y = 16
		.ImageList = list = CreateImageList(x, y)
		imageMap.Each()
			{
			padding = it.GetDefault('padding', 0)
			ImageList_AddVectorImage(list, it.image, it.color, x, y, :padding)
			ImageList_AddVectorImage(list, it.image, it.color, x, y, dark?:, :padding)
			}
		}

	getIcon(icon)
		{
		return this[icon $ 'Theme'] = IDE_ColorScheme.IsDark?()
			? this[icon $ 'Dark']
			: this[icon]
		}

	Getter_(member)
		{
		if member.Suffix?('Theme')
			return .getIcon(member.RemoveSuffix('Theme'))
		throw 'Invalid theme member - ' $ member
		}

	SetTheme(data, theme)
		{
		data.theme = this[theme $ 'Theme']
		data.image = this[theme]
		}

	ResetTheme()
		{
		.Delete(#FolderOpenTheme)
		.Delete(#FolderOpenModifiedTheme)
		.Delete(#FolderClosedTheme)
		.Delete(#FolderClosedModifiedTheme)
		.Delete(#DocumentTheme)
		.Delete(#CloseButtonTheme)
		.Delete(#ModifiedTheme)
		.Delete(#InvalidTheme)
		.Delete(#EditTheme)
		}

	Destroy()
		{
		ImageList_Destroy(.ImageList)
		if not Sys.SuneidoJs?()
			.ImageResources.Each({ it[0].Close() })
		}
	}
