// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
#(
	(name: Qc_FunctionComplexity,
		type: 'normal', warning:, noUpdates:),

	(name: Qc_LineSize,
		type: 'normal', warning:),

	(name: Qc_MagicNumbers,
		type: 'normal', warning:, noUpdates:, noTests:),

	(name: Qc_MethodOrFuncSize,
		type: 'normal', warning:, noUpdates:, noTests:),

	(name: Qc_NumParams,
		type: 'normal', warning:, noUpdates:),

	(name: Qc_PrivMethods,
		type: 'normal', warning:, noUpdates:),

	(name: Qc_Coupling,
		type: 'extra', warning:, noUpdates:)

	(name: Qc_RecordLength,
		type: 'extra', warning:, noUpdates:),

	(name: Qc_TestChecker,
		type: 'extra', warning:, noUpdates:),

	(name: Qc_UndefinedPublicMethods,
		type: 'extra', warning:),

	(name: Qc_DuplicateCode,
		type: 'slow', warning:)
)
