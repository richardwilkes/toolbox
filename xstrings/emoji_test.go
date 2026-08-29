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

	c.True(xstrings.SelectEmojiRegex.MatchString("😀"))
	c.True(xstrings.SelectEmojiRegex.MatchString("😊"))
	c.True(xstrings.SelectEmojiRegex.MatchString("🎉"))
	c.True(xstrings.SelectEmojiRegex.MatchString("🚀"))
	c.True(xstrings.SelectEmojiRegex.MatchString("❤️"))
	c.True(xstrings.SelectEmojiRegex.MatchString("👍"))

	c.True(xstrings.SelectEmojiRegex.MatchString("Hello 😀 World"))
	c.True(xstrings.SelectEmojiRegex.MatchString("Great job! 👍"))
	c.True(xstrings.SelectEmojiRegex.MatchString("🎉 Celebration time! 🎊"))

	c.False(xstrings.SelectEmojiRegex.MatchString("Hello World"))
	c.False(xstrings.SelectEmojiRegex.MatchString("No emojis here"))
	c.False(xstrings.SelectEmojiRegex.MatchString("123 ABC"))
	c.False(xstrings.SelectEmojiRegex.MatchString(""))

	c.False(xstrings.SelectEmojiRegex.MatchString("@#$%^&*()"))
	c.False(xstrings.SelectEmojiRegex.MatchString("+-=[]{}|;':\"<>?"))
	c.False(xstrings.SelectEmojiRegex.MatchString("àáâãäåæçèéêë"))
}

func TestSelectEmojiRegexFindAllString(t *testing.T) {
	c := check.New(t)

	matches := xstrings.SelectEmojiRegex.FindAllString("Hello 😀 World", -1)
	c.Equal(1, len(matches))
	c.Equal("😀", matches[0])

	matches = xstrings.SelectEmojiRegex.FindAllString("😀😊🎉", -1)
	c.Equal(3, len(matches))
	c.Equal("😀", matches[0])
	c.Equal("😊", matches[1])
	c.Equal("🎉", matches[2])

	matches = xstrings.SelectEmojiRegex.FindAllString("Start 🚀 middle 👍 end 🎊", -1)
	c.Equal(3, len(matches))
	c.Equal("🚀", matches[0])
	c.Equal("👍", matches[1])
	c.Equal("🎊", matches[2])

	matches = xstrings.SelectEmojiRegex.FindAllString("No emojis here", -1)
	c.Equal(0, len(matches))

	matches = xstrings.SelectEmojiRegex.FindAllString("", -1)
	c.Equal(0, len(matches))
}

func TestSelectEmojiRegexCommonEmojis(t *testing.T) {
	c := check.New(t)

	c.True(xstrings.SelectEmojiRegex.MatchString("😀")) // grinning face
	c.True(xstrings.SelectEmojiRegex.MatchString("😂")) // face with tears of joy
	c.True(xstrings.SelectEmojiRegex.MatchString("😍")) // smiling face with heart-eyes
	c.True(xstrings.SelectEmojiRegex.MatchString("🤔")) // thinking face
	c.True(xstrings.SelectEmojiRegex.MatchString("😭")) // loudly crying face

	c.True(xstrings.SelectEmojiRegex.MatchString("👍")) // thumbs up
	c.True(xstrings.SelectEmojiRegex.MatchString("👎")) // thumbs down
	c.True(xstrings.SelectEmojiRegex.MatchString("👋")) // waving hand
	c.True(xstrings.SelectEmojiRegex.MatchString("🤝")) // handshake
	c.True(xstrings.SelectEmojiRegex.MatchString("👏")) // clapping hands

	c.True(xstrings.SelectEmojiRegex.MatchString("🚗")) // automobile
	c.True(xstrings.SelectEmojiRegex.MatchString("🏠")) // house
	c.True(xstrings.SelectEmojiRegex.MatchString("📱")) // mobile phone
	c.True(xstrings.SelectEmojiRegex.MatchString("💻")) // laptop computer
	c.True(xstrings.SelectEmojiRegex.MatchString("📝")) // memo

	c.True(xstrings.SelectEmojiRegex.MatchString("🌳"))  // deciduous tree
	c.True(xstrings.SelectEmojiRegex.MatchString("🌸"))  // cherry blossom
	c.True(xstrings.SelectEmojiRegex.MatchString("☀️")) // sun
	c.True(xstrings.SelectEmojiRegex.MatchString("🌙"))  // crescent moon
	c.True(xstrings.SelectEmojiRegex.MatchString("⭐"))  // star

	c.True(xstrings.SelectEmojiRegex.MatchString("🍕")) // pizza
	c.True(xstrings.SelectEmojiRegex.MatchString("🍔")) // hamburger
	c.True(xstrings.SelectEmojiRegex.MatchString("🍎")) // red apple
	c.True(xstrings.SelectEmojiRegex.MatchString("🍰")) // shortcake
	c.True(xstrings.SelectEmojiRegex.MatchString("☕")) // hot beverage
}

func TestSelectEmojiRegexComplexEmojis(t *testing.T) {
	c := check.New(t)

	c.True(xstrings.SelectEmojiRegex.MatchString("👋🏻")) // waving hand: light skin tone
	c.True(xstrings.SelectEmojiRegex.MatchString("👋🏽")) // waving hand: medium skin tone
	c.True(xstrings.SelectEmojiRegex.MatchString("👋🏿")) // waving hand: dark skin tone

	c.True(xstrings.SelectEmojiRegex.MatchString("👨‍👩‍👧‍👦")) // family: man, woman, girl, boy
	c.True(xstrings.SelectEmojiRegex.MatchString("👩‍💻"))     // woman technologist
	c.True(xstrings.SelectEmojiRegex.MatchString("👨‍🚀"))     // man astronaut

	c.True(xstrings.SelectEmojiRegex.MatchString("🇺🇸")) // United States flag
	c.True(xstrings.SelectEmojiRegex.MatchString("🇬🇧")) // United Kingdom flag
	c.True(xstrings.SelectEmojiRegex.MatchString("🇯🇵")) // Japan flag

	c.True(xstrings.SelectEmojiRegex.MatchString("❤️")) // red heart with variation selector
	c.True(xstrings.SelectEmojiRegex.MatchString("✨"))  // sparkles
	c.True(xstrings.SelectEmojiRegex.MatchString("⚡"))  // high voltage
}

func TestSelectEmojiRegexReplaceAllString(t *testing.T) {
	c := check.New(t)

	result := xstrings.SelectEmojiRegex.ReplaceAllString("Hello 😀 World 👍", "[EMOJI]")
	c.Equal("Hello [EMOJI] World [EMOJI]", result)

	result = xstrings.SelectEmojiRegex.ReplaceAllString("Start 🚀 middle 👍 end 🎊", "")
	c.Equal("Start  middle  end ", result)

	result = xstrings.SelectEmojiRegex.ReplaceAllString("No emojis here", "[EMOJI]")
	c.Equal("No emojis here", result)

	result = xstrings.SelectEmojiRegex.ReplaceAllString("😀😊🎉", "[E]")
	c.Equal("[E][E][E]", result)

	result = xstrings.SelectEmojiRegex.ReplaceAllString("Great work! 👏 Keep it up! 💪 You're awesome! 🌟", " ")
	c.Equal("Great work!   Keep it up!   You're awesome!  ", result)
}

func TestSelectEmojiRegexSplit(t *testing.T) {
	c := check.New(t)

	parts := xstrings.SelectEmojiRegex.Split("Hello😀World👍End", -1)
	c.Equal(3, len(parts))
	c.Equal("Hello", parts[0])
	c.Equal("World", parts[1])
	c.Equal("End", parts[2])

	parts = xstrings.SelectEmojiRegex.Split("Hello 😀 World 👍 End", -1)
	c.Equal(3, len(parts))
	c.Equal("Hello ", parts[0])
	c.Equal(" World ", parts[1])
	c.Equal(" End", parts[2])

	parts = xstrings.SelectEmojiRegex.Split("No emojis here", -1)
	c.Equal(1, len(parts))
	c.Equal("No emojis here", parts[0])

	parts = xstrings.SelectEmojiRegex.Split("😀Hello World", -1)
	c.Equal(2, len(parts))
	c.Equal("", parts[0])
	c.Equal("Hello World", parts[1])

	parts = xstrings.SelectEmojiRegex.Split("Hello World😀", -1)
	c.Equal(2, len(parts))
	c.Equal("Hello World", parts[0])
	c.Equal("", parts[1])
}

func TestSelectEmojiRegexEdgeCases(t *testing.T) {
	c := check.New(t)

	c.False(xstrings.SelectEmojiRegex.MatchString("Ñ"))
	c.False(xstrings.SelectEmojiRegex.MatchString("ü"))
	c.False(xstrings.SelectEmojiRegex.MatchString("ñoño"))
	c.False(xstrings.SelectEmojiRegex.MatchString("über"))
	c.False(xstrings.SelectEmojiRegex.MatchString("café"))

	c.False(xstrings.SelectEmojiRegex.MatchString("∑"))
	c.False(xstrings.SelectEmojiRegex.MatchString("∆"))
	c.False(xstrings.SelectEmojiRegex.MatchString("∞"))
	c.False(xstrings.SelectEmojiRegex.MatchString("π"))

	c.False(xstrings.SelectEmojiRegex.MatchString("$"))
	c.False(xstrings.SelectEmojiRegex.MatchString("€"))
	c.False(xstrings.SelectEmojiRegex.MatchString("£"))
	c.False(xstrings.SelectEmojiRegex.MatchString("¥"))

	c.False(xstrings.SelectEmojiRegex.MatchString(" "))
	c.False(xstrings.SelectEmojiRegex.MatchString("\t"))
	c.False(xstrings.SelectEmojiRegex.MatchString("\n"))
	c.False(xstrings.SelectEmojiRegex.MatchString("\r"))

	c.False(xstrings.SelectEmojiRegex.MatchString("\u0000"))
	c.False(xstrings.SelectEmojiRegex.MatchString("\u0001"))
	c.False(xstrings.SelectEmojiRegex.MatchString("\u001F"))
}
