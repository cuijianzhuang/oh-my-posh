package segments

import (
	"errors"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
)

func TestSysInfo(t *testing.T) {
	cases := []struct {
		Error          error
		Case           string
		ExpectedString string
		Template       string
		SysInfo        runtime.SystemInfo
		Precision      int
		ExpectDisabled bool
	}{
		{
			Case:           "Error",
			ExpectDisabled: true,
			Error:          errors.New("error"),
		},
		{
			Case:           "physical mem",
			ExpectedString: "50",
			SysInfo: runtime.SystemInfo{
				Memory: runtime.Memory{
					PhysicalPercentUsed: 50,
				},
			},
			Template: "{{ round .PhysicalPercentUsed .Precision }}",
		},
		{
			Case:           "physical mem 2 digits",
			ExpectedString: "60.51",
			SysInfo: runtime.SystemInfo{
				Memory: runtime.Memory{
					PhysicalPercentUsed: 60.51,
				},
			},
			Precision: 2,
			Template:  "{{ round .PhysicalPercentUsed .Precision }}",
		},
		{
			Case:           "physical meme rounded",
			ExpectedString: "61",
			SysInfo: runtime.SystemInfo{
				Memory: runtime.Memory{
					PhysicalPercentUsed: 61,
				},
			},
			Template: "{{ round .PhysicalPercentUsed .Precision }}",
		},
		{
			Case:           "load",
			ExpectedString: "0.22 0.12 0",
			Precision:      2,
			Template:       "{{ round .Load1 .Precision }} {{round .Load5 .Precision }} {{round .Load15 .Precision }}",
			SysInfo:        runtime.SystemInfo{Load1: 0.22, Load5: 0.12, Load15: 0},
		},
		{
			Case:           "not enabled",
			ExpectDisabled: true,
			SysInfo: runtime.SystemInfo{
				Memory: runtime.Memory{
					PhysicalPercentUsed: 0,
					SwapPercentUsed:     0,
				},
			},
		},
	}

	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("SystemInfo").Return(&tc.SysInfo, tc.Error)
		sysInfo := &SystemInfo{}
		props := properties.Map{
			Precision: tc.Precision,
		}
		sysInfo.Init(props, env)
		enabled := sysInfo.Enabled()
		if tc.ExpectDisabled {
			assert.Equal(t, false, enabled, tc.Case)
		} else {
			assert.Equal(t, tc.ExpectedString, renderTemplate(env, tc.Template, sysInfo), tc.Case)
		}
	}
}
