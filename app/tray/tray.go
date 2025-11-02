package tray

import (
	"fmt"
	"runtime"

	"github.com/ollama/ollama/app/assets"
	"github.com/ollama/ollama/app/tray/commontray"
)

func NewTray() (commontray.OllamaTray, error) {
	// Use .ico for all platforms as that's what we have in assets
	extension := ".ico"
	iconName := commontray.UpdateIconName + extension
	updateIcon, err := assets.GetIcon(iconName)
	if err != nil {
		return nil, fmt.Errorf("failed to load icon %s: %w", iconName, err)
	}
	iconName = commontray.IconName + extension
	icon, err := assets.GetIcon(iconName)
	if err != nil {
		return nil, fmt.Errorf("failed to load icon %s: %w", iconName, err)
	}

	return InitPlatformTray(icon, updateIcon)
}
