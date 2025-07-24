// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	CssStyle()
		{
		return '<style>' $
			'body {' $ 	StdFonts.GetCSSFont(sizeFactor: 0.8) $
				'background-color: #b7d5ff;' $
			'}' $
			'.msg-container {
				display: flex;
				flex-direction: column;
				max-width: 75%;
				width: fit-content;
				word-wrap: break-word;
				white-space: pre-wrap;
				margin: 5px 0px;
				border: 2px;
				padding: 10px;
			}
			' $  .outgoingMsgs() $ '
			.msg-in {
				align-self: flex-start;
				border-radius: 10px 25px;
				background-color: #f9f9fb;
				}
			.msg-body {
				text-align: left;
			}
			.container {
				display: flex;
				flex-direction: column;
				}
			</style>'
		}

	outgoingMsgs()
		{
		return '.msg-out {
				align-self: flex-end;
				border-radius: 25px 10px;
				color: white;
				background-color: #004895;
				}
			.msg-out :link {
				color: #00ffff;
				}
			.msg-out :visited {
				color: #FF00FF;
				}'
		}

	IsTyping?()
		{
		return false
		}

	FocusEditor()
		{
		return
		}
	}