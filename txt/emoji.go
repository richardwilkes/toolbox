// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package txt

import "regexp"

// SelectEmojiRegex identifies emoji runs.
var SelectEmojiRegex = regexp.MustCompile(`[\x{200D}\x{203C}\x{2049}\x{20E3}\x{2122}\x{2139}\x{2194}-\x{2199}` +
	`\x{21A9}-\x{21AA}\x{231A}-\x{231B}\x{2328}\x{2388}\x{23CF}\x{23E9}-\x{23F3}\x{23F8}-\x{23FA}\x{24C2}` +
	`\x{25AA}-\x{25AB}\x{25B6}\x{25C0}\x{25FB}-\x{25FE}\x{2600}-\x{2605}\x{2607}-\x{2612}\x{2614}-\x{2705}` +
	`\x{2708}-\x{2712}\x{2714}\x{2716}\x{271D}\x{2721}\x{2728}\x{2733}-\x{2734}\x{2744}\x{2747}\x{274C}\x{274E}` +
	`\x{2753}-\x{2755}\x{2757}\x{2763}-\x{2767}\x{2795}-\x{2797}\x{27A1}\x{27B0}\x{27BF}\x{2934}-\x{2935}` +
	`\x{2B05}-\x{2B07}\x{2B1B}-\x{2B1C}\x{2B50}\x{2B55}\x{3030}\x{303D}\x{3297}\x{3299}\x{FE00}-\x{FE0F}` +
	`\x{1F000}-\x{1F0FF}\x{1F10D}-\x{1F10F}\x{1F12F}\x{1F16C}-\x{1F171}\x{1F17E}-\x{1F17F}\x{1F18E}` +
	`\x{1F191}-\x{1F19A}\x{1F1AD}-\x{1F1FF}\x{1F201}-\x{1F20F}\x{1F21A}\x{1F22F}\x{1F232}-\x{1F23A}` +
	`\x{1F23C}-\x{1F23F}\x{1F249}-\x{1F62B}\x{1F62C}-\x{1F64F}\x{1F680}-\x{1F6FF}\x{1F774}-\x{1F77F}` +
	`\x{1F7D5}-\x{1F7FF}\x{1F80C}-\x{1F80F}\x{1F848}-\x{1F84F}\x{1F85A}-\x{1F85F}\x{1F888}-\x{1F88F}` +
	`\x{1F8AE}-\x{1F93A}\x{1F93C}-\x{1F945}\x{1F947}-\x{1F9FF}\x{E0020}-\x{E007F}]`)
