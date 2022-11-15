package oceancd

var (
	PromoteAction     = "promote"
	PromoteFullAction = "promoteFull"
	PauseAction       = "pause"
	AbortAction       = "abort"
	RetryAction       = "retry"
	RestartAction     = "restart"
	RollbackAction    = "rollback"
)

type QueryParams map[string]string

type PathParams map[string]string
