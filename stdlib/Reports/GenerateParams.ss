// Copyright (C) 2010 Axon Development Corporation. All rights reserved worldwide.
Params
	{
	ExtraButtons: ('Generate',	'Generate && Print', 'Generate && PDF')
	New(@report)
		{
		super(@.setup(report))
		}

	setup(report)
		{
		// Need this to prevent users from leaving report screen while preview is open
		// because generate from preview requires the report screen
		report.previewDialog = true
		return report
		}

	action: false
	On_Generate(@report)
		{
		if false is .generateValid()
			return

		if false is .runNoOutput(report)
			return
		.updateResult(report)
		}

	On_Generate_Result(@report)
		{
		if false is .generateValid()
			return false

		result = .RunWithNoOutput(report)
		.unlockGenerate()
		if result is false
			return false
		.updateResult(report)
		return result
		}

	runNoOutput(@report)
		{
		result = .RunWithNoOutput(report)
		.unlockGenerate()
		return result
		}

	updateResult(report)
		{
		msgctrl = .FindControl('generate_msg')
		msgctrl.Remove(1)
		msgctrl.Insert(1, Object('Heading1' 'Complete'))
		.CloseDialog(report)
		.action = false
		}

	On_Generate__Print(@report)
		{
		if false isnt .generateValid()
			{
			this.On_Print(report)
			.unlockGenerate()
			}
		.action = false
		}

	On_Generate__PDF_Save_to_file(@report)
		{
		if false isnt .generateValid()
			{
			_fromGenerateParams = true
			this.On_PDF_Save_to_file(@report)
			.unlockGenerate()
			}
		.action = false
		}

	On_Generate__PDF_Download(@report)
		{
		.On_Generate__PDF_Save_to_file(@report)
		}

	On_Generate__PDF_Email_as_attachment(@report)
		{
		if false isnt .generateValid()
			{
			_fromGenerateParams = true
			this.On_PDF_Email_as_attachment(@report)
			.unlockGenerate()
			}
		.action = false
		}

	SetExtraParamsData(paramsData)
		{
		paramsData.action = .action
		}

	generateValid()
		{
		msg = ""
		params = .Vert.Data.Get()
		if "" isnt msg = params[.Params_report.generateValidField]
			{
			.AlertWarn("Invalid Parameter(s)", msg)
			return false
			}
		if true isnt lock = ServerEval("Biz_LockProcess.Lock", .Name, Suneido.User)
			{
			.AlertWarn("Process Locked", "User " $ lock.user $
				" is already generating, started " $ lock.asof.ShortDateTime() $
				"\nIf the Generate Process failed to complete, it will remain " $
				"locked for 30 minutes.")
			return false
			}

		.action = 'Generate'
		return true
		}

	unlockGenerate()
		{
		ServerEval("Biz_LockProcess.UnLock", .Name)
		}
	}
