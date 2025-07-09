package txt

import (
	"testing"

	"github.com/richardwilkes/toolbox/v2/check"
)

func TestSelectEmojiRegex(t *testing.T) {
	c := check.New(t)

	// Test that the regex is properly compiled
	c.NotNil(SelectEmojiRegex)

	// Test basic emoji detection
	c.True(SelectEmojiRegex.MatchString("😀"))
	c.True(SelectEmojiRegex.MatchString("😊"))
	c.True(SelectEmojiRegex.MatchString("🎉"))
	c.True(SelectEmojiRegex.MatchString("🚀"))
	c.True(SelectEmojiRegex.MatchString("❤️"))
	c.True(SelectEmojiRegex.MatchString("👍"))

	// Test emoji in text
	c.True(SelectEmojiRegex.MatchString("Hello 😀 World"))
	c.True(SelectEmojiRegex.MatchString("Great job! 👍"))
	c.True(SelectEmojiRegex.MatchString("🎉 Celebration time! 🎊"))

	// Test non-emoji text
	c.False(SelectEmojiRegex.MatchString("Hello World"))
	c.False(SelectEmojiRegex.MatchString("No emojis here"))
	c.False(SelectEmojiRegex.MatchString("123 ABC"))
	c.False(SelectEmojiRegex.MatchString(""))

	// Test symbols that might be confused with emojis
	c.False(SelectEmojiRegex.MatchString("@#$%^&*()"))
	c.False(SelectEmojiRegex.MatchString("+-=[]{}|;':\"<>?"))
	c.False(SelectEmojiRegex.MatchString("àáâãäåæçèéêë"))
}

func TestSelectEmojiRegexFindAllString(t *testing.T) {
	c := check.New(t)

	// Test finding single emoji
	matches := SelectEmojiRegex.FindAllString("Hello 😀 World", -1)
	c.Equal(1, len(matches))
	c.Equal("😀", matches[0])

	// Test finding multiple emojis
	matches = SelectEmojiRegex.FindAllString("😀😊🎉", -1)
	c.Equal(3, len(matches))
	c.Equal("😀", matches[0])
	c.Equal("😊", matches[1])
	c.Equal("🎉", matches[2])

	// Test finding emojis in text
	matches = SelectEmojiRegex.FindAllString("Start 🚀 middle 👍 end 🎊", -1)
	c.Equal(3, len(matches))
	c.Equal("🚀", matches[0])
	c.Equal("👍", matches[1])
	c.Equal("🎊", matches[2])

	// Test text without emojis
	matches = SelectEmojiRegex.FindAllString("No emojis here", -1)
	c.Equal(0, len(matches))

	// Test empty string
	matches = SelectEmojiRegex.FindAllString("", -1)
	c.Equal(0, len(matches))
}

func TestSelectEmojiRegexCommonEmojis(t *testing.T) {
	c := check.New(t)

	// Test face emojis
	c.True(SelectEmojiRegex.MatchString("😀")) // grinning face
	c.True(SelectEmojiRegex.MatchString("😂")) // face with tears of joy
	c.True(SelectEmojiRegex.MatchString("😍")) // smiling face with heart-eyes
	c.True(SelectEmojiRegex.MatchString("🤔")) // thinking face
	c.True(SelectEmojiRegex.MatchString("😭")) // loudly crying face

	// Test hand gestures
	c.True(SelectEmojiRegex.MatchString("👍")) // thumbs up
	c.True(SelectEmojiRegex.MatchString("👎")) // thumbs down
	c.True(SelectEmojiRegex.MatchString("👋")) // waving hand
	c.True(SelectEmojiRegex.MatchString("🤝")) // handshake
	c.True(SelectEmojiRegex.MatchString("👏")) // clapping hands

	// Test objects and symbols
	c.True(SelectEmojiRegex.MatchString("🚗")) // automobile
	c.True(SelectEmojiRegex.MatchString("🏠")) // house
	c.True(SelectEmojiRegex.MatchString("📱")) // mobile phone
	c.True(SelectEmojiRegex.MatchString("💻")) // laptop computer
	c.True(SelectEmojiRegex.MatchString("📝")) // memo

	// Test nature emojis
	c.True(SelectEmojiRegex.MatchString("🌳"))  // deciduous tree
	c.True(SelectEmojiRegex.MatchString("🌸"))  // cherry blossom
	c.True(SelectEmojiRegex.MatchString("☀️")) // sun
	c.True(SelectEmojiRegex.MatchString("🌙"))  // crescent moon
	c.True(SelectEmojiRegex.MatchString("⭐"))  // star

	// Test food emojis
	c.True(SelectEmojiRegex.MatchString("🍕")) // pizza
	c.True(SelectEmojiRegex.MatchString("🍔")) // hamburger
	c.True(SelectEmojiRegex.MatchString("🍎")) // red apple
	c.True(SelectEmojiRegex.MatchString("🍰")) // shortcake
	c.True(SelectEmojiRegex.MatchString("☕")) // hot beverage
}

func TestSelectEmojiRegexComplexEmojis(t *testing.T) {
	c := check.New(t)

	// Test emojis with skin tone modifiers
	c.True(SelectEmojiRegex.MatchString("👋🏻")) // waving hand: light skin tone
	c.True(SelectEmojiRegex.MatchString("👋🏽")) // waving hand: medium skin tone
	c.True(SelectEmojiRegex.MatchString("👋🏿")) // waving hand: dark skin tone

	// Test multi-part emojis (family, couples, etc.)
	c.True(SelectEmojiRegex.MatchString("👨‍👩‍👧‍👦")) // family: man, woman, girl, boy
	c.True(SelectEmojiRegex.MatchString("👩‍💻"))     // woman technologist
	c.True(SelectEmojiRegex.MatchString("👨‍🚀"))     // man astronaut

	// Test flag emojis
	c.True(SelectEmojiRegex.MatchString("🇺🇸")) // United States flag
	c.True(SelectEmojiRegex.MatchString("🇬🇧")) // United Kingdom flag
	c.True(SelectEmojiRegex.MatchString("🇯🇵")) // Japan flag

	// Test variation selectors
	c.True(SelectEmojiRegex.MatchString("❤️")) // red heart with variation selector
	c.True(SelectEmojiRegex.MatchString("✨"))  // sparkles
	c.True(SelectEmojiRegex.MatchString("⚡"))  // high voltage
}

func TestSelectEmojiRegexReplaceAllString(t *testing.T) {
	c := check.New(t)

	// Test replacing emojis with placeholder
	result := SelectEmojiRegex.ReplaceAllString("Hello 😀 World 👍", "[EMOJI]")
	c.Equal("Hello [EMOJI] World [EMOJI]", result)

	// Test removing emojis entirely
	result = SelectEmojiRegex.ReplaceAllString("Start 🚀 middle 👍 end 🎊", "")
	c.Equal("Start  middle  end ", result)

	// Test text without emojis
	result = SelectEmojiRegex.ReplaceAllString("No emojis here", "[EMOJI]")
	c.Equal("No emojis here", result)

	// Test multiple consecutive emojis
	result = SelectEmojiRegex.ReplaceAllString("😀😊🎉", "[E]")
	c.Equal("[E][E][E]", result)

	// Test complex text with multiple emojis
	result = SelectEmojiRegex.ReplaceAllString("Great work! 👏 Keep it up! 💪 You're awesome! 🌟", " ")
	c.Equal("Great work!   Keep it up!   You're awesome!  ", result)
}

func TestSelectEmojiRegexSplit(t *testing.T) {
	c := check.New(t)

	// Test splitting text by emojis
	parts := SelectEmojiRegex.Split("Hello😀World👍End", -1)
	c.Equal(3, len(parts))
	c.Equal("Hello", parts[0])
	c.Equal("World", parts[1])
	c.Equal("End", parts[2])

	// Test splitting with spaces around emojis
	parts = SelectEmojiRegex.Split("Hello 😀 World 👍 End", -1)
	c.Equal(3, len(parts))
	c.Equal("Hello ", parts[0])
	c.Equal(" World ", parts[1])
	c.Equal(" End", parts[2])

	// Test text without emojis
	parts = SelectEmojiRegex.Split("No emojis here", -1)
	c.Equal(1, len(parts))
	c.Equal("No emojis here", parts[0])

	// Test starting with emoji
	parts = SelectEmojiRegex.Split("😀Hello World", -1)
	c.Equal(2, len(parts))
	c.Equal("", parts[0])
	c.Equal("Hello World", parts[1])

	// Test ending with emoji
	parts = SelectEmojiRegex.Split("Hello World😀", -1)
	c.Equal(2, len(parts))
	c.Equal("Hello World", parts[0])
	c.Equal("", parts[1])
}

func TestSelectEmojiRegexEdgeCases(t *testing.T) {
	c := check.New(t)

	// Test unicode characters that are not emojis
	c.False(SelectEmojiRegex.MatchString("Ñ"))
	c.False(SelectEmojiRegex.MatchString("ü"))
	c.False(SelectEmojiRegex.MatchString("ñoño"))
	c.False(SelectEmojiRegex.MatchString("über"))
	c.False(SelectEmojiRegex.MatchString("café"))

	// Test mathematical symbols
	c.False(SelectEmojiRegex.MatchString("∑"))
	c.False(SelectEmojiRegex.MatchString("∆"))
	c.False(SelectEmojiRegex.MatchString("∞"))
	c.False(SelectEmojiRegex.MatchString("π"))

	// Test currency symbols
	c.False(SelectEmojiRegex.MatchString("$"))
	c.False(SelectEmojiRegex.MatchString("€"))
	c.False(SelectEmojiRegex.MatchString("£"))
	c.False(SelectEmojiRegex.MatchString("¥"))

	// Test whitespace characters
	c.False(SelectEmojiRegex.MatchString(" "))
	c.False(SelectEmojiRegex.MatchString("\t"))
	c.False(SelectEmojiRegex.MatchString("\n"))
	c.False(SelectEmojiRegex.MatchString("\r"))

	// Test control characters
	c.False(SelectEmojiRegex.MatchString("\u0000"))
	c.False(SelectEmojiRegex.MatchString("\u0001"))
	c.False(SelectEmojiRegex.MatchString("\u001F"))
}
