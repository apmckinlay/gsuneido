// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	New(field)
		{
		super(@.controls(field))
		.FindControl(.buttonName).SetCommandTarget(this)
		.fieldCtrl = .FindControl(field)
		}

	controls(field)
		{
		tip = 'Copy ' $ prompt = Display(Prompt(field))
		.info = 'Copied ' $ prompt $ ' to clipboard'
		return Object(field,
			Object(#EnhancedButton, image: #copy, imagePadding: 0.15, alignTop:,
				command: #CopyField, name: .buttonName = #CopyButton_ $ field, :tip),
				leftAlign:)
		}

	On_CopyField()
		{
		ClipboardWriteString(.fieldCtrl.Get())
		InfoWindowControl(.info, titleSize: 0, marginSize: 7, autoClose: 2)
		}
	}