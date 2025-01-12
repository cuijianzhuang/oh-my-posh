//go:build darwin

package segments

import (
	"errors"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
)

func TestSpotifyDarwinEnabledAndSpotifyPlaying(t *testing.T) {
	cases := []struct {
		Error    error
		Running  string
		Expected string
		Status   string
		Artist   string
		Track    string
	}{
		{Running: "false", Expected: ""},
		{Running: "false", Expected: "", Error: errors.New("oops")},
		{Running: "true", Expected: "\ue602 Candlemass - Spellbreaker", Status: "playing", Artist: "Candlemass", Track: "Spellbreaker"},
		{Running: "true", Expected: "\uF8E3 Candlemass - Spellbreaker", Status: "paused", Artist: "Candlemass", Track: "Spellbreaker"},
	}
	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("RunCommand", "osascript", []string{"-e", "application \"Spotify\" is running"}).Return(tc.Running, tc.Error)
		env.On("RunCommand", "osascript", []string{"-e", "tell application \"Spotify\" to player state as string"}).Return(tc.Status, nil)
		env.On("RunCommand", "osascript", []string{"-e", "tell application \"Spotify\" to artist of current track as string"}).Return(tc.Artist, nil)
		env.On("RunCommand", "osascript", []string{"-e", "tell application \"Spotify\" to name of current track as string"}).Return(tc.Track, nil)

		s := &Spotify{}
		s.Init(properties.Map{}, env)

		assert.Equal(t, tc.Running == "true", s.Enabled())
		assert.Equal(t, tc.Expected, renderTemplate(env, s.Template(), s))
	}
}
