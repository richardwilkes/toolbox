// Copyright (c) 2016-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xstrings_test

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
	"github.com/richardwilkes/toolbox/v2/xstrings"
)

func TestSelectEmojiRegex(t *testing.T) {
	c := check.New(t)
	c.NotNil(xstrings.SelectEmojiRegex)
	table := []struct {
		In    string
		Match bool
	}{
		{In: "😀", Match: true},
		{In: "😊", Match: true},
		{In: "🎉", Match: true},
		{In: "🚀", Match: true},
		{In: "❤️", Match: true},
		{In: "👍", Match: true},
		{In: "😂", Match: true},
		{In: "😍", Match: true},
		{In: "🤔", Match: true},
		{In: "😭", Match: true},
		{In: "👎", Match: true},
		{In: "👋", Match: true},
		{In: "🤝", Match: true},
		{In: "👏", Match: true},
		{In: "🚗", Match: true},
		{In: "🏠", Match: true},
		{In: "📱", Match: true},
		{In: "💻", Match: true},
		{In: "📝", Match: true},
		{In: "🌳", Match: true},
		{In: "🌸", Match: true},
		{In: "☀️", Match: true},
		{In: "🌙", Match: true},
		{In: "⭐", Match: true},
		{In: "🍕", Match: true},
		{In: "🍔", Match: true},
		{In: "🍎", Match: true},
		{In: "🍰", Match: true},
		{In: "☕", Match: true},
		{In: "✨", Match: true},
		{In: "⚡", Match: true},
		// Skin tone modifiers, ZWJ sequences and regional indicator flags.
		{In: "👋🏻", Match: true},
		{In: "👋🏽", Match: true},
		{In: "👋🏿", Match: true},
		{In: "👨‍👩‍👧‍👦", Match: true},
		{In: "👩‍💻", Match: true},
		{In: "👨‍🚀", Match: true},
		{In: "🇺🇸", Match: true},
		{In: "🇬🇧", Match: true},
		{In: "🇯🇵", Match: true},

		{In: "Hello 😀 World", Match: true},
		{In: "Great job! 👍", Match: true},
		{In: "🎉 Celebration time! 🎊", Match: true},
		{In: "Hello World"},
		{In: "No emojis here"},
		{In: "123 ABC"},
		{In: ""},
		{In: "@#$%^&*()"},
		{In: "+-=[]{}|;':\"<>?"},
		{In: "àáâãäåæçèéêë"},
		{In: "Ñ"},
		{In: "ü"},
		{In: "ñoño"},
		{In: "über"},
		{In: "café"},
		{In: "∑"},
		{In: "∆"},
		{In: "∞"},
		{In: "π"},
		{In: "$"},
		{In: "€"},
		{In: "£"},
		{In: "¥"},
		{In: " "},
		{In: "\t"},
		{In: "\n"},
		{In: "\r"},
		{In: "\u0000"},
		{In: "\u0001"},
		{In: "\u001F"},
	}
	for _, one := range table {
		c.Equal(one.Match, xstrings.SelectEmojiRegex.MatchString(one.In), "%q", one.In)
	}
}

func TestSelectEmojiRegexFindAllString(t *testing.T) {
	table := []struct {
		In  string
		Out []string
	}{
		{In: "Hello 😀 World", Out: []string{"😀"}},
		{In: "😀😊🎉", Out: []string{"😀", "😊", "🎉"}},
		{In: "Start 🚀 middle 👍 end 🎊", Out: []string{"🚀", "👍", "🎊"}},
		{In: "Hello 😀 World 👍", Out: []string{"😀", "👍"}},
		{In: "Great work! 👏 Keep it up! 💪 You're awesome! 🌟", Out: []string{"👏", "💪", "🌟"}},
		{In: "Hello😀World👍End", Out: []string{"😀", "👍"}},
		{In: "😀Hello World", Out: []string{"😀"}},
		{In: "Hello World😀", Out: []string{"😀"}},
		{In: "No emojis here"},
		{In: ""},
	}
	c := check.New(t)
	for _, one := range table {
		c.Equal(one.Out, xstrings.SelectEmojiRegex.FindAllString(one.In, -1), "%q", one.In)
	}
}
