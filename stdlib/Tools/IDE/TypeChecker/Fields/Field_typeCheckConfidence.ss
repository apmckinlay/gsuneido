Field_number
	{
	Prompt: "Confidence %"
	Width: 4
	Mask: "###"          // 0-100, no decimals at the UI layer
	Valid?(value)
		{
		return value is "" or (Number?(x = Number(value)) and x >= 0 and x <= 100)
		}
	ValidData?(value)
		{
		return .Valid?(value)
		}
	}

