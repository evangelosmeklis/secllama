package commontray

var (
	Title   = "SecLlama"
	ToolTip = "SecLlama - Secure AI"

	UpdateIconName = "tray_upgrade"
	IconName       = "tray"
)

type Callbacks struct {
	Quit       chan struct{}
	Update     chan struct{}
	DoFirstUse chan struct{}
	ShowLogs   chan struct{}
}

type OllamaTray interface { // Note: kept as OllamaTray for code compatibility
	GetCallbacks() Callbacks
	Run()
	UpdateAvailable(ver string) error
	DisplayFirstUseNotification() error
	Quit()
}
