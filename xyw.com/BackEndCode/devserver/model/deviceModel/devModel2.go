package deviceModel

import (
	"xutils/xtext"
)

func (p *DeviceModel) IsAlertValue() string {
	alert := p.IsAlert
	if xtext.IsBlank(p.IsAlert) {
		alert = "0"
	}

	return alert
}
