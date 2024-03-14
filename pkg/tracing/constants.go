package tracing

const (
	BindType        = "error.bind"
	QueryType       = "error.query"
	TimeFormatType  = "error.query-time-format"
	MakeRequestType = "error.make-request"

	CreateMatrixType   = "error.create-matrix"
	GetHistoryType     = "error.get-history"
	GetDifferenceType  = "error.get-different"
	GetTendencyType    = "error.get-tendency"
	GetMatrixType      = "error.get-matrix"
	GetMatrixPagesType = "error.get-matrix-pages"

	PrepareAndSendStorageType = "error.prepare-storage"
)

const (
	CallToService = "Call to service"

	CreateMatrix          = "CreateMatrixWithoutParent Matrix"
	GetHistory            = "Get history"
	GetDifference         = "Get difference"
	GetTendency           = "Get tendency"
	GetMatrix             = "Get matrix"
	GetMatrixPages        = "Get matrix pages"
	GetMatricesByDuration = "Get matrices by duration"

	PrepareAndSendStorage = "Prepare and send storage"
	SwitchStorageToNext   = "Switch storage to next"
	GetCurrentStorage     = "Get current storage"
)

const (
	SuccessfulCompleting = "Operation completed successfully"
)
