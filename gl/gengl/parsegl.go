
// line 1 "parsegl.rl"
// -*-go-*-
//
// To compile:
//
//   ragel -Z -G2 -o atoi.go atoi.rl
//   6g -o atoi.6 atoi.go
//   6l -o atoi atoi.6
//   ./atoi
//
// To show a diagram of your state machine:
//
//   ragel -V -G2 -p -o atoi.dot atoi.rl
//   dot -Tpng -o atoi.png atoi.dot
//   chrome atoi.png
//

package main

import (
	"fmt"
	"strings"
)


// line 28 "parsegl.go"
const parsegl_start int = 193
const parsegl_first_final int = 193
const parsegl_error int = 0

const parsegl_en_main int = 193
const parsegl_en_skiperror int = 192


// line 28 "parsegl.rl"


func parse(data string, header *Header) error {
	var cs, p, pe int
	var ts, te, act, eof int

	pe = len(data)
	eof = len(data)

	_, _, _ = ts, te, act

	//stack := make([]int, 32)
	//top := 0

	var curline = 1

	var m0, m1, m2, m3, m4, m5 int
	var heading string
	var lineblock int
	var ifblock int
	var f Func

	
// line 61 "parsegl.go"
	{
	cs = parsegl_start
	ts = 0
	te = 0
	act = 0
	}

// line 69 "parsegl.go"
	{
	if p == pe {
		goto _test_eof
	}
	switch cs {
	case 193:
		goto st_case_193
	case 0:
		goto st_case_0
	case 1:
		goto st_case_1
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	case 4:
		goto st_case_4
	case 194:
		goto st_case_194
	case 5:
		goto st_case_5
	case 6:
		goto st_case_6
	case 7:
		goto st_case_7
	case 8:
		goto st_case_8
	case 9:
		goto st_case_9
	case 10:
		goto st_case_10
	case 11:
		goto st_case_11
	case 12:
		goto st_case_12
	case 13:
		goto st_case_13
	case 14:
		goto st_case_14
	case 15:
		goto st_case_15
	case 16:
		goto st_case_16
	case 17:
		goto st_case_17
	case 18:
		goto st_case_18
	case 19:
		goto st_case_19
	case 20:
		goto st_case_20
	case 21:
		goto st_case_21
	case 22:
		goto st_case_22
	case 23:
		goto st_case_23
	case 24:
		goto st_case_24
	case 25:
		goto st_case_25
	case 26:
		goto st_case_26
	case 27:
		goto st_case_27
	case 28:
		goto st_case_28
	case 29:
		goto st_case_29
	case 30:
		goto st_case_30
	case 31:
		goto st_case_31
	case 32:
		goto st_case_32
	case 33:
		goto st_case_33
	case 34:
		goto st_case_34
	case 35:
		goto st_case_35
	case 36:
		goto st_case_36
	case 37:
		goto st_case_37
	case 38:
		goto st_case_38
	case 39:
		goto st_case_39
	case 40:
		goto st_case_40
	case 41:
		goto st_case_41
	case 42:
		goto st_case_42
	case 43:
		goto st_case_43
	case 44:
		goto st_case_44
	case 45:
		goto st_case_45
	case 46:
		goto st_case_46
	case 47:
		goto st_case_47
	case 48:
		goto st_case_48
	case 49:
		goto st_case_49
	case 50:
		goto st_case_50
	case 195:
		goto st_case_195
	case 51:
		goto st_case_51
	case 52:
		goto st_case_52
	case 53:
		goto st_case_53
	case 54:
		goto st_case_54
	case 55:
		goto st_case_55
	case 56:
		goto st_case_56
	case 57:
		goto st_case_57
	case 58:
		goto st_case_58
	case 59:
		goto st_case_59
	case 60:
		goto st_case_60
	case 61:
		goto st_case_61
	case 62:
		goto st_case_62
	case 63:
		goto st_case_63
	case 64:
		goto st_case_64
	case 65:
		goto st_case_65
	case 66:
		goto st_case_66
	case 67:
		goto st_case_67
	case 68:
		goto st_case_68
	case 69:
		goto st_case_69
	case 70:
		goto st_case_70
	case 71:
		goto st_case_71
	case 72:
		goto st_case_72
	case 73:
		goto st_case_73
	case 74:
		goto st_case_74
	case 75:
		goto st_case_75
	case 76:
		goto st_case_76
	case 77:
		goto st_case_77
	case 78:
		goto st_case_78
	case 79:
		goto st_case_79
	case 80:
		goto st_case_80
	case 81:
		goto st_case_81
	case 82:
		goto st_case_82
	case 83:
		goto st_case_83
	case 84:
		goto st_case_84
	case 85:
		goto st_case_85
	case 86:
		goto st_case_86
	case 87:
		goto st_case_87
	case 88:
		goto st_case_88
	case 89:
		goto st_case_89
	case 90:
		goto st_case_90
	case 91:
		goto st_case_91
	case 92:
		goto st_case_92
	case 93:
		goto st_case_93
	case 94:
		goto st_case_94
	case 95:
		goto st_case_95
	case 96:
		goto st_case_96
	case 97:
		goto st_case_97
	case 98:
		goto st_case_98
	case 99:
		goto st_case_99
	case 100:
		goto st_case_100
	case 101:
		goto st_case_101
	case 102:
		goto st_case_102
	case 103:
		goto st_case_103
	case 104:
		goto st_case_104
	case 105:
		goto st_case_105
	case 106:
		goto st_case_106
	case 107:
		goto st_case_107
	case 108:
		goto st_case_108
	case 109:
		goto st_case_109
	case 110:
		goto st_case_110
	case 111:
		goto st_case_111
	case 112:
		goto st_case_112
	case 113:
		goto st_case_113
	case 114:
		goto st_case_114
	case 115:
		goto st_case_115
	case 116:
		goto st_case_116
	case 117:
		goto st_case_117
	case 118:
		goto st_case_118
	case 119:
		goto st_case_119
	case 120:
		goto st_case_120
	case 121:
		goto st_case_121
	case 122:
		goto st_case_122
	case 123:
		goto st_case_123
	case 124:
		goto st_case_124
	case 125:
		goto st_case_125
	case 126:
		goto st_case_126
	case 127:
		goto st_case_127
	case 128:
		goto st_case_128
	case 129:
		goto st_case_129
	case 130:
		goto st_case_130
	case 131:
		goto st_case_131
	case 132:
		goto st_case_132
	case 133:
		goto st_case_133
	case 134:
		goto st_case_134
	case 135:
		goto st_case_135
	case 136:
		goto st_case_136
	case 137:
		goto st_case_137
	case 138:
		goto st_case_138
	case 139:
		goto st_case_139
	case 140:
		goto st_case_140
	case 141:
		goto st_case_141
	case 142:
		goto st_case_142
	case 143:
		goto st_case_143
	case 144:
		goto st_case_144
	case 145:
		goto st_case_145
	case 146:
		goto st_case_146
	case 147:
		goto st_case_147
	case 148:
		goto st_case_148
	case 149:
		goto st_case_149
	case 150:
		goto st_case_150
	case 151:
		goto st_case_151
	case 196:
		goto st_case_196
	case 152:
		goto st_case_152
	case 153:
		goto st_case_153
	case 154:
		goto st_case_154
	case 155:
		goto st_case_155
	case 156:
		goto st_case_156
	case 157:
		goto st_case_157
	case 158:
		goto st_case_158
	case 159:
		goto st_case_159
	case 160:
		goto st_case_160
	case 161:
		goto st_case_161
	case 162:
		goto st_case_162
	case 163:
		goto st_case_163
	case 164:
		goto st_case_164
	case 165:
		goto st_case_165
	case 166:
		goto st_case_166
	case 167:
		goto st_case_167
	case 168:
		goto st_case_168
	case 169:
		goto st_case_169
	case 170:
		goto st_case_170
	case 171:
		goto st_case_171
	case 172:
		goto st_case_172
	case 173:
		goto st_case_173
	case 174:
		goto st_case_174
	case 175:
		goto st_case_175
	case 176:
		goto st_case_176
	case 177:
		goto st_case_177
	case 178:
		goto st_case_178
	case 179:
		goto st_case_179
	case 180:
		goto st_case_180
	case 181:
		goto st_case_181
	case 182:
		goto st_case_182
	case 183:
		goto st_case_183
	case 184:
		goto st_case_184
	case 185:
		goto st_case_185
	case 186:
		goto st_case_186
	case 187:
		goto st_case_187
	case 188:
		goto st_case_188
	case 189:
		goto st_case_189
	case 190:
		goto st_case_190
	case 191:
		goto st_case_191
	case 192:
		goto st_case_192
	case 197:
		goto st_case_197
	case 198:
		goto st_case_198
	}
	goto st_out
tr2:
// line 51 "parsegl.rl"

 curline++ 
// line 135 "parsegl.rl"

te = p+1
{
				// Reset heading comment.
				heading = ""

				// Start new line block.
				lineblock++
			}
	goto st193
tr22:
// line 51 "parsegl.rl"

 curline++ 
// line 96 "parsegl.rl"

te = p+1
{
				// Ignore defines within conditional blocks other than the top level ifndef.
				if ifblock == 1 && strings.HasPrefix(data[m0:m1], "GL_") {
					header.Define = append(header.Define, Define{data[m0:m1], data[m2:m3], heading, lineblock})
					heading = ""
				}
			}
	goto st193
tr34:
// line 51 "parsegl.rl"

 curline++ 
// line 87 "parsegl.rl"

te = p+1

	goto st193
tr38:
// line 51 "parsegl.rl"

 curline++ 
// line 82 "parsegl.rl"

te = p+1
{
				ifblock--
			}
	goto st193
tr43:
// line 51 "parsegl.rl"

 curline++ 
// line 78 "parsegl.rl"

te = p+1
{
				ifblock++
			}
	goto st193
tr53:
// line 51 "parsegl.rl"

 curline++ 
// line 71 "parsegl.rl"

te = p+1

	goto st193
tr119:
// line 124 "parsegl.rl"

te = p+1
{
				if (ifblock == 1) {
					if f.Type == "void" && f.Addr == 0 {
						f.Type = ""
					}
					header.Func = append(header.Func, f)
				}
			}
	goto st193
tr179:
// line 51 "parsegl.rl"

 curline++ 
// line 74 "parsegl.rl"

te = p+1

	goto st193
tr198:
// line 114 "parsegl.rl"

te = p+1

	goto st193
tr213:
// line 107 "parsegl.rl"

p = (te) - 1
{
				if (ifblock == 1) {
					header.Type = append(header.Type, Type{data[m2:m3], data[m0:m1], strings.TrimSpace(data[m4:m5])})
				}
			}
	goto st193
tr221:
// line 51 "parsegl.rl"

 curline++ 
// line 107 "parsegl.rl"

te = p+1
{
				if (ifblock == 1) {
					header.Type = append(header.Type, Type{data[m2:m3], data[m0:m1], strings.TrimSpace(data[m4:m5])})
				}
			}
	goto st193
tr266:
// line 68 "parsegl.rl"

te = p
p--

	goto st193
tr268:
// line 51 "parsegl.rl"

 curline++ 
// line 62 "parsegl.rl"

te = p+1
{
				heading = strings.Trim(data[m0:m1], " *")
				lineblock++
			}
	goto st193
tr269:
// line 107 "parsegl.rl"

te = p
p--
{
				if (ifblock == 1) {
					header.Type = append(header.Type, Type{data[m2:m3], data[m0:m1], strings.TrimSpace(data[m4:m5])})
				}
			}
	goto st193
	st193:
// line 1 "NONE"

ts = 0

		if p++; p == pe {
			goto _test_eof193
		}
	st_case_193:
// line 1 "NONE"

ts = p

// line 639 "parsegl.go"
		switch data[p] {
		case 9:
			goto st1
		case 10:
			goto tr2
		case 32:
			goto st1
		case 35:
			goto st5
		case 47:
			goto st48
		case 71:
			goto st51
		case 101:
			goto st117
		case 116:
			goto st129
		case 125:
			goto st128
		}
		goto st0
st_case_0:
	st0:
		cs = 0
		goto _out
	st1:
		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
		switch data[p] {
		case 9:
			goto st1
		case 10:
			goto tr2
		case 32:
			goto st1
		case 47:
			goto st2
		}
		goto st0
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
		if data[p] == 42 {
			goto st3
		}
		goto st0
tr5:
// line 51 "parsegl.rl"

 curline++ 
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
// line 700 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr5
		case 42:
			goto st4
		}
		goto st3
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
		switch data[p] {
		case 10:
			goto tr5
		case 42:
			goto st4
		case 47:
			goto st194
		}
		goto st3
	st194:
		if p++; p == pe {
			goto _test_eof194
		}
	st_case_194:
		switch data[p] {
		case 9:
			goto st194
		case 32:
			goto st194
		}
		goto tr266
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
		switch data[p] {
		case 9:
			goto st6
		case 32:
			goto st6
		case 100:
			goto st7
		case 101:
			goto st19
		case 105:
			goto st35
		case 112:
			goto st43
		}
		goto st0
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
		switch data[p] {
		case 9:
			goto st6
		case 32:
			goto st6
		case 100:
			goto st7
		case 101:
			goto st19
		case 105:
			goto st28
		}
		goto st0
	st7:
		if p++; p == pe {
			goto _test_eof7
		}
	st_case_7:
		if data[p] == 101 {
			goto st8
		}
		goto st0
	st8:
		if p++; p == pe {
			goto _test_eof8
		}
	st_case_8:
		if data[p] == 102 {
			goto st9
		}
		goto st0
	st9:
		if p++; p == pe {
			goto _test_eof9
		}
	st_case_9:
		if data[p] == 105 {
			goto st10
		}
		goto st0
	st10:
		if p++; p == pe {
			goto _test_eof10
		}
	st_case_10:
		if data[p] == 110 {
			goto st11
		}
		goto st0
	st11:
		if p++; p == pe {
			goto _test_eof11
		}
	st_case_11:
		if data[p] == 101 {
			goto st12
		}
		goto st0
	st12:
		if p++; p == pe {
			goto _test_eof12
		}
	st_case_12:
		switch data[p] {
		case 9:
			goto st13
		case 32:
			goto st13
		}
		goto st0
	st13:
		if p++; p == pe {
			goto _test_eof13
		}
	st_case_13:
		switch data[p] {
		case 9:
			goto st13
		case 32:
			goto st13
		case 95:
			goto tr20
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr20
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr20
			}
		default:
			goto tr20
		}
		goto st0
tr20:
// line 92 "parsegl.rl"

 m0 = p; m2, m3 = 0, 0 
// line 92 "parsegl.rl"

 m1 = p+1 
	goto st14
tr23:
// line 92 "parsegl.rl"

 m1 = p+1 
	goto st14
	st14:
		if p++; p == pe {
			goto _test_eof14
		}
	st_case_14:
// line 873 "parsegl.go"
		switch data[p] {
		case 9:
			goto st15
		case 10:
			goto tr22
		case 32:
			goto st15
		case 95:
			goto tr23
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr23
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr23
			}
		default:
			goto tr23
		}
		goto st0
tr25:
// line 94 "parsegl.rl"

 m2 = p 
// line 94 "parsegl.rl"

 m3 = p+1 
	goto st15
	st15:
		if p++; p == pe {
			goto _test_eof15
		}
	st_case_15:
// line 910 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr25
		case 10:
			goto tr22
		case 32:
			goto st15
		}
		goto tr24
tr24:
// line 94 "parsegl.rl"

 m2 = p 
// line 94 "parsegl.rl"

 m3 = p+1 
	goto st16
tr26:
// line 94 "parsegl.rl"

 m3 = p+1 
	goto st16
	st16:
		if p++; p == pe {
			goto _test_eof16
		}
	st_case_16:
// line 938 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr22
		case 32:
			goto st17
		}
		goto tr26
	st17:
		if p++; p == pe {
			goto _test_eof17
		}
	st_case_17:
		switch data[p] {
		case 10:
			goto tr22
		case 32:
			goto st18
		}
		goto tr26
	st18:
		if p++; p == pe {
			goto _test_eof18
		}
	st_case_18:
		switch data[p] {
		case 9:
			goto st18
		case 10:
			goto tr22
		case 32:
			goto st18
		}
		goto st0
	st19:
		if p++; p == pe {
			goto _test_eof19
		}
	st_case_19:
		switch data[p] {
		case 108:
			goto st20
		case 110:
			goto st24
		}
		goto st0
	st20:
		if p++; p == pe {
			goto _test_eof20
		}
	st_case_20:
		switch data[p] {
		case 105:
			goto st21
		case 115:
			goto st23
		}
		goto st0
	st21:
		if p++; p == pe {
			goto _test_eof21
		}
	st_case_21:
		if data[p] == 102 {
			goto st22
		}
		goto st0
	st22:
		if p++; p == pe {
			goto _test_eof22
		}
	st_case_22:
		if data[p] == 10 {
			goto tr34
		}
		goto st22
	st23:
		if p++; p == pe {
			goto _test_eof23
		}
	st_case_23:
		if data[p] == 101 {
			goto st22
		}
		goto st0
	st24:
		if p++; p == pe {
			goto _test_eof24
		}
	st_case_24:
		if data[p] == 100 {
			goto st25
		}
		goto st0
	st25:
		if p++; p == pe {
			goto _test_eof25
		}
	st_case_25:
		if data[p] == 105 {
			goto st26
		}
		goto st0
	st26:
		if p++; p == pe {
			goto _test_eof26
		}
	st_case_26:
		if data[p] == 102 {
			goto st27
		}
		goto st0
	st27:
		if p++; p == pe {
			goto _test_eof27
		}
	st_case_27:
		if data[p] == 10 {
			goto tr38
		}
		goto st27
	st28:
		if p++; p == pe {
			goto _test_eof28
		}
	st_case_28:
		if data[p] == 102 {
			goto st29
		}
		goto st0
	st29:
		if p++; p == pe {
			goto _test_eof29
		}
	st_case_29:
		switch data[p] {
		case 9:
			goto st30
		case 32:
			goto st30
		case 100:
			goto st31
		case 110:
			goto st34
		}
		goto st0
	st30:
		if p++; p == pe {
			goto _test_eof30
		}
	st_case_30:
		if data[p] == 10 {
			goto tr43
		}
		goto st30
	st31:
		if p++; p == pe {
			goto _test_eof31
		}
	st_case_31:
		if data[p] == 101 {
			goto st32
		}
		goto st0
	st32:
		if p++; p == pe {
			goto _test_eof32
		}
	st_case_32:
		if data[p] == 102 {
			goto st33
		}
		goto st0
	st33:
		if p++; p == pe {
			goto _test_eof33
		}
	st_case_33:
		switch data[p] {
		case 9:
			goto st30
		case 32:
			goto st30
		}
		goto st0
	st34:
		if p++; p == pe {
			goto _test_eof34
		}
	st_case_34:
		if data[p] == 100 {
			goto st31
		}
		goto st0
	st35:
		if p++; p == pe {
			goto _test_eof35
		}
	st_case_35:
		switch data[p] {
		case 102:
			goto st29
		case 110:
			goto st36
		}
		goto st0
	st36:
		if p++; p == pe {
			goto _test_eof36
		}
	st_case_36:
		if data[p] == 99 {
			goto st37
		}
		goto st0
	st37:
		if p++; p == pe {
			goto _test_eof37
		}
	st_case_37:
		if data[p] == 108 {
			goto st38
		}
		goto st0
	st38:
		if p++; p == pe {
			goto _test_eof38
		}
	st_case_38:
		if data[p] == 117 {
			goto st39
		}
		goto st0
	st39:
		if p++; p == pe {
			goto _test_eof39
		}
	st_case_39:
		if data[p] == 100 {
			goto st40
		}
		goto st0
	st40:
		if p++; p == pe {
			goto _test_eof40
		}
	st_case_40:
		if data[p] == 101 {
			goto st41
		}
		goto st0
	st41:
		if p++; p == pe {
			goto _test_eof41
		}
	st_case_41:
		if data[p] == 10 {
			goto st0
		}
		goto st42
	st42:
		if p++; p == pe {
			goto _test_eof42
		}
	st_case_42:
		if data[p] == 10 {
			goto tr53
		}
		goto st42
	st43:
		if p++; p == pe {
			goto _test_eof43
		}
	st_case_43:
		if data[p] == 114 {
			goto st44
		}
		goto st0
	st44:
		if p++; p == pe {
			goto _test_eof44
		}
	st_case_44:
		if data[p] == 97 {
			goto st45
		}
		goto st0
	st45:
		if p++; p == pe {
			goto _test_eof45
		}
	st_case_45:
		if data[p] == 103 {
			goto st46
		}
		goto st0
	st46:
		if p++; p == pe {
			goto _test_eof46
		}
	st_case_46:
		if data[p] == 109 {
			goto st47
		}
		goto st0
	st47:
		if p++; p == pe {
			goto _test_eof47
		}
	st_case_47:
		if data[p] == 97 {
			goto st41
		}
		goto st0
	st48:
		if p++; p == pe {
			goto _test_eof48
		}
	st_case_48:
		if data[p] == 42 {
			goto tr58
		}
		goto st0
tr60:
// line 51 "parsegl.rl"

 curline++ 
	goto st49
tr58:
// line 61 "parsegl.rl"

 m0 = p 
	goto st49
	st49:
		if p++; p == pe {
			goto _test_eof49
		}
	st_case_49:
// line 1276 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr60
		case 42:
			goto st50
		}
		goto st49
	st50:
		if p++; p == pe {
			goto _test_eof50
		}
	st_case_50:
		switch data[p] {
		case 10:
			goto tr60
		case 42:
			goto st50
		case 47:
			goto tr62
		}
		goto st49
tr62:
// line 61 "parsegl.rl"

 m1 = p-2 
	goto st195
	st195:
		if p++; p == pe {
			goto _test_eof195
		}
	st_case_195:
// line 1308 "parsegl.go"
		switch data[p] {
		case 9:
			goto st195
		case 10:
			goto tr268
		case 32:
			goto st195
		}
		goto tr266
	st51:
		if p++; p == pe {
			goto _test_eof51
		}
	st_case_51:
		if data[p] == 76 {
			goto st52
		}
		goto st0
	st52:
		if p++; p == pe {
			goto _test_eof52
		}
	st_case_52:
		if data[p] == 65 {
			goto st53
		}
		goto st0
	st53:
		if p++; p == pe {
			goto _test_eof53
		}
	st_case_53:
		if data[p] == 80 {
			goto st54
		}
		goto st0
	st54:
		if p++; p == pe {
			goto _test_eof54
		}
	st_case_54:
		if data[p] == 73 {
			goto st55
		}
		goto st0
	st55:
		if p++; p == pe {
			goto _test_eof55
		}
	st_case_55:
		switch data[p] {
		case 9:
			goto st56
		case 10:
			goto tr68
		case 32:
			goto st56
		}
		goto st0
tr68:
// line 51 "parsegl.rl"

 curline++ 
	goto st56
	st56:
		if p++; p == pe {
			goto _test_eof56
		}
	st_case_56:
// line 1378 "parsegl.go"
		switch data[p] {
		case 9:
			goto st56
		case 10:
			goto tr68
		case 32:
			goto st56
		case 95:
			goto tr69
		case 99:
			goto tr70
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr69
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr69
			}
		default:
			goto tr69
		}
		goto st0
tr69:
// line 117 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st57
	st57:
		if p++; p == pe {
			goto _test_eof57
		}
	st_case_57:
// line 1414 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
tr75:
// line 51 "parsegl.rl"

 curline++ 
	goto st58
tr71:
// line 117 "parsegl.rl"

 m1 = p 
	goto st58
tr72:
// line 117 "parsegl.rl"

 m1 = p 
// line 51 "parsegl.rl"

 curline++ 
	goto st58
	st58:
		if p++; p == pe {
			goto _test_eof58
		}
	st_case_58:
// line 1461 "parsegl.go"
		switch data[p] {
		case 9:
			goto st58
		case 10:
			goto tr75
		case 32:
			goto st58
		case 42:
			goto tr76
		case 65:
			goto st61
		case 71:
			goto st98
		}
		goto st0
tr76:
// line 117 "parsegl.rl"

 m4++ 
	goto st59
	st59:
		if p++; p == pe {
			goto _test_eof59
		}
	st_case_59:
// line 1487 "parsegl.go"
		switch data[p] {
		case 9:
			goto st60
		case 10:
			goto tr80
		case 32:
			goto st60
		case 42:
			goto tr76
		case 65:
			goto st61
		case 71:
			goto st98
		}
		goto st0
tr80:
// line 51 "parsegl.rl"

 curline++ 
	goto st60
	st60:
		if p++; p == pe {
			goto _test_eof60
		}
	st_case_60:
// line 1513 "parsegl.go"
		switch data[p] {
		case 9:
			goto st60
		case 10:
			goto tr80
		case 32:
			goto st60
		case 65:
			goto st61
		case 71:
			goto st98
		}
		goto st0
	st61:
		if p++; p == pe {
			goto _test_eof61
		}
	st_case_61:
		if data[p] == 80 {
			goto st62
		}
		goto st0
	st62:
		if p++; p == pe {
			goto _test_eof62
		}
	st_case_62:
		if data[p] == 73 {
			goto st63
		}
		goto st0
	st63:
		if p++; p == pe {
			goto _test_eof63
		}
	st_case_63:
		if data[p] == 69 {
			goto st64
		}
		goto st0
	st64:
		if p++; p == pe {
			goto _test_eof64
		}
	st_case_64:
		if data[p] == 78 {
			goto st65
		}
		goto st0
	st65:
		if p++; p == pe {
			goto _test_eof65
		}
	st_case_65:
		if data[p] == 84 {
			goto st66
		}
		goto st0
	st66:
		if p++; p == pe {
			goto _test_eof66
		}
	st_case_66:
		if data[p] == 82 {
			goto st67
		}
		goto st0
	st67:
		if p++; p == pe {
			goto _test_eof67
		}
	st_case_67:
		if data[p] == 89 {
			goto st68
		}
		goto st0
	st68:
		if p++; p == pe {
			goto _test_eof68
		}
	st_case_68:
		switch data[p] {
		case 9:
			goto st69
		case 10:
			goto tr89
		case 32:
			goto st69
		}
		goto st0
tr89:
// line 51 "parsegl.rl"

 curline++ 
	goto st69
	st69:
		if p++; p == pe {
			goto _test_eof69
		}
	st_case_69:
// line 1614 "parsegl.go"
		switch data[p] {
		case 9:
			goto st69
		case 10:
			goto tr89
		case 32:
			goto st69
		case 103:
			goto tr90
		}
		goto st0
tr90:
// line 119 "parsegl.rl"

 m2 = p 
	goto st70
	st70:
		if p++; p == pe {
			goto _test_eof70
		}
	st_case_70:
// line 1636 "parsegl.go"
		if data[p] == 108 {
			goto st71
		}
		goto st0
	st71:
		if p++; p == pe {
			goto _test_eof71
		}
	st_case_71:
		if data[p] == 95 {
			goto st72
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st72
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st72
			}
		default:
			goto st72
		}
		goto st0
	st72:
		if p++; p == pe {
			goto _test_eof72
		}
	st_case_72:
		switch data[p] {
		case 9:
			goto tr93
		case 10:
			goto tr94
		case 32:
			goto tr93
		case 40:
			goto tr95
		case 95:
			goto st72
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st72
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st72
			}
		default:
			goto st72
		}
		goto st0
tr97:
// line 51 "parsegl.rl"

 curline++ 
	goto st73
tr93:
// line 119 "parsegl.rl"

 m3 = p; f = Func{data[m2:m3], data[m0:m1], m4, nil} 
	goto st73
tr94:
// line 119 "parsegl.rl"

 m3 = p; f = Func{data[m2:m3], data[m0:m1], m4, nil} 
// line 51 "parsegl.rl"

 curline++ 
	goto st73
	st73:
		if p++; p == pe {
			goto _test_eof73
		}
	st_case_73:
// line 1715 "parsegl.go"
		switch data[p] {
		case 9:
			goto st73
		case 10:
			goto tr97
		case 32:
			goto st73
		case 40:
			goto st74
		}
		goto st0
tr99:
// line 51 "parsegl.rl"

 curline++ 
	goto st74
tr95:
// line 119 "parsegl.rl"

 m3 = p; f = Func{data[m2:m3], data[m0:m1], m4, nil} 
	goto st74
	st74:
		if p++; p == pe {
			goto _test_eof74
		}
	st_case_74:
// line 1742 "parsegl.go"
		switch data[p] {
		case 9:
			goto st74
		case 10:
			goto tr99
		case 32:
			goto st74
		case 95:
			goto tr100
		case 99:
			goto tr101
		case 118:
			goto tr102
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr100
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr100
			}
		default:
			goto tr100
		}
		goto st0
tr100:
// line 121 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st75
	st75:
		if p++; p == pe {
			goto _test_eof75
		}
	st_case_75:
// line 1780 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
tr107:
// line 51 "parsegl.rl"

 curline++ 
	goto st76
tr103:
// line 121 "parsegl.rl"

 m1 = p 
	goto st76
tr104:
// line 121 "parsegl.rl"

 m1 = p 
// line 51 "parsegl.rl"

 curline++ 
	goto st76
	st76:
		if p++; p == pe {
			goto _test_eof76
		}
	st_case_76:
// line 1827 "parsegl.go"
		switch data[p] {
		case 9:
			goto st76
		case 10:
			goto tr107
		case 32:
			goto st76
		case 42:
			goto tr108
		case 95:
			goto tr109
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr109
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr109
			}
		default:
			goto tr109
		}
		goto st0
tr108:
// line 121 "parsegl.rl"

 m4++ 
	goto st77
	st77:
		if p++; p == pe {
			goto _test_eof77
		}
	st_case_77:
// line 1863 "parsegl.go"
		switch data[p] {
		case 9:
			goto st78
		case 10:
			goto tr111
		case 32:
			goto st78
		case 42:
			goto tr108
		case 95:
			goto tr109
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr109
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr109
			}
		default:
			goto tr109
		}
		goto st0
tr111:
// line 51 "parsegl.rl"

 curline++ 
	goto st78
	st78:
		if p++; p == pe {
			goto _test_eof78
		}
	st_case_78:
// line 1899 "parsegl.go"
		switch data[p] {
		case 9:
			goto st78
		case 10:
			goto tr111
		case 32:
			goto st78
		case 95:
			goto tr109
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr109
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr109
			}
		default:
			goto tr109
		}
		goto st0
tr109:
// line 121 "parsegl.rl"

 m2 = p; m5 = 0 
// line 121 "parsegl.rl"

 m3 = p+1 
	goto st79
tr115:
// line 121 "parsegl.rl"

 m3 = p+1 
	goto st79
	st79:
		if p++; p == pe {
			goto _test_eof79
		}
	st_case_79:
// line 1941 "parsegl.go"
		switch data[p] {
		case 9:
			goto st80
		case 10:
			goto tr113
		case 32:
			goto st80
		case 41:
			goto tr114
		case 44:
			goto tr114
		case 91:
			goto st90
		case 95:
			goto tr115
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr115
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr115
			}
		default:
			goto tr115
		}
		goto st0
tr113:
// line 51 "parsegl.rl"

 curline++ 
	goto st80
	st80:
		if p++; p == pe {
			goto _test_eof80
		}
	st_case_80:
// line 1981 "parsegl.go"
		switch data[p] {
		case 9:
			goto st80
		case 10:
			goto tr113
		case 32:
			goto st80
		case 41:
			goto tr114
		case 44:
			goto tr114
		}
		goto st0
tr118:
// line 51 "parsegl.rl"

 curline++ 
	goto st81
tr114:
// line 122 "parsegl.rl"

 f.Param = append(f.Param, Param{data[m2:m3], data[m0:m1], m4, m5}) 
	goto st81
	st81:
		if p++; p == pe {
			goto _test_eof81
		}
	st_case_81:
// line 2010 "parsegl.go"
		switch data[p] {
		case 9:
			goto st81
		case 10:
			goto tr118
		case 32:
			goto st81
		case 59:
			goto tr119
		case 95:
			goto tr100
		case 99:
			goto tr101
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr100
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr100
			}
		default:
			goto tr100
		}
		goto st0
tr101:
// line 121 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st82
	st82:
		if p++; p == pe {
			goto _test_eof82
		}
	st_case_82:
// line 2048 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 111:
			goto st83
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st83:
		if p++; p == pe {
			goto _test_eof83
		}
	st_case_83:
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 110:
			goto st84
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st84:
		if p++; p == pe {
			goto _test_eof84
		}
	st_case_84:
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 115:
			goto st85
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st85:
		if p++; p == pe {
			goto _test_eof85
		}
	st_case_85:
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 116:
			goto st86
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st86:
		if p++; p == pe {
			goto _test_eof86
		}
	st_case_86:
		switch data[p] {
		case 9:
			goto tr124
		case 10:
			goto tr125
		case 32:
			goto tr124
		case 95:
			goto st75
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
tr127:
// line 51 "parsegl.rl"

 curline++ 
	goto st87
tr124:
// line 121 "parsegl.rl"

 m1 = p 
	goto st87
tr125:
// line 51 "parsegl.rl"

 curline++ 
// line 121 "parsegl.rl"

 m1 = p 
	goto st87
	st87:
		if p++; p == pe {
			goto _test_eof87
		}
	st_case_87:
// line 2215 "parsegl.go"
		switch data[p] {
		case 9:
			goto st87
		case 10:
			goto tr127
		case 32:
			goto st87
		case 42:
			goto tr108
		case 95:
			goto tr128
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr128
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr128
			}
		default:
			goto tr128
		}
		goto st0
tr131:
// line 121 "parsegl.rl"

 m3 = p+1 
	goto st88
tr128:
// line 121 "parsegl.rl"

 m0 = p; m4 = 0 
// line 121 "parsegl.rl"

 m2 = p; m5 = 0 
// line 121 "parsegl.rl"

 m3 = p+1 
	goto st88
	st88:
		if p++; p == pe {
			goto _test_eof88
		}
	st_case_88:
// line 2262 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr129
		case 10:
			goto tr130
		case 32:
			goto tr129
		case 41:
			goto tr114
		case 44:
			goto tr114
		case 91:
			goto st90
		case 95:
			goto tr131
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr131
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr131
			}
		default:
			goto tr131
		}
		goto st0
tr133:
// line 51 "parsegl.rl"

 curline++ 
	goto st89
tr129:
// line 121 "parsegl.rl"

 m1 = p 
	goto st89
tr130:
// line 121 "parsegl.rl"

 m1 = p 
// line 51 "parsegl.rl"

 curline++ 
	goto st89
	st89:
		if p++; p == pe {
			goto _test_eof89
		}
	st_case_89:
// line 2315 "parsegl.go"
		switch data[p] {
		case 9:
			goto st89
		case 10:
			goto tr133
		case 32:
			goto st89
		case 41:
			goto tr114
		case 42:
			goto tr108
		case 44:
			goto tr114
		case 95:
			goto tr109
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr109
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr109
			}
		default:
			goto tr109
		}
		goto st0
	st90:
		if p++; p == pe {
			goto _test_eof90
		}
	st_case_90:
		if 48 <= data[p] && data[p] <= 57 {
			goto tr134
		}
		goto st0
tr134:
// line 121 "parsegl.rl"

 m5 = m5*10 + (int(data[p]) - '0') 
	goto st91
	st91:
		if p++; p == pe {
			goto _test_eof91
		}
	st_case_91:
// line 2364 "parsegl.go"
		if data[p] == 93 {
			goto st80
		}
		if 48 <= data[p] && data[p] <= 57 {
			goto tr134
		}
		goto st0
tr102:
// line 121 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st92
	st92:
		if p++; p == pe {
			goto _test_eof92
		}
	st_case_92:
// line 2382 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 111:
			goto st93
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st93:
		if p++; p == pe {
			goto _test_eof93
		}
	st_case_93:
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 105:
			goto st94
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st94:
		if p++; p == pe {
			goto _test_eof94
		}
	st_case_94:
		switch data[p] {
		case 9:
			goto tr103
		case 10:
			goto tr104
		case 32:
			goto tr103
		case 95:
			goto st75
		case 100:
			goto st95
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
	st95:
		if p++; p == pe {
			goto _test_eof95
		}
	st_case_95:
		switch data[p] {
		case 9:
			goto tr138
		case 10:
			goto tr139
		case 32:
			goto tr138
		case 41:
			goto st97
		case 95:
			goto st75
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st75
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st75
			}
		default:
			goto st75
		}
		goto st0
tr142:
// line 51 "parsegl.rl"

 curline++ 
	goto st96
tr138:
// line 121 "parsegl.rl"

 m1 = p 
	goto st96
tr139:
// line 51 "parsegl.rl"

 curline++ 
// line 121 "parsegl.rl"

 m1 = p 
	goto st96
	st96:
		if p++; p == pe {
			goto _test_eof96
		}
	st_case_96:
// line 2521 "parsegl.go"
		switch data[p] {
		case 9:
			goto st96
		case 10:
			goto tr142
		case 32:
			goto st96
		case 41:
			goto st97
		case 42:
			goto tr108
		case 95:
			goto tr109
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr109
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr109
			}
		default:
			goto tr109
		}
		goto st0
tr143:
// line 51 "parsegl.rl"

 curline++ 
	goto st97
	st97:
		if p++; p == pe {
			goto _test_eof97
		}
	st_case_97:
// line 2559 "parsegl.go"
		switch data[p] {
		case 9:
			goto st97
		case 10:
			goto tr143
		case 32:
			goto st97
		case 59:
			goto tr119
		}
		goto st0
	st98:
		if p++; p == pe {
			goto _test_eof98
		}
	st_case_98:
		if data[p] == 76 {
			goto st99
		}
		goto st0
	st99:
		if p++; p == pe {
			goto _test_eof99
		}
	st_case_99:
		if data[p] == 65 {
			goto st61
		}
		goto st0
tr70:
// line 117 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st100
	st100:
		if p++; p == pe {
			goto _test_eof100
		}
	st_case_100:
// line 2599 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 95:
			goto st57
		case 111:
			goto st101
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st101:
		if p++; p == pe {
			goto _test_eof101
		}
	st_case_101:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 95:
			goto st57
		case 110:
			goto st102
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st102:
		if p++; p == pe {
			goto _test_eof102
		}
	st_case_102:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 95:
			goto st57
		case 115:
			goto st103
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st103:
		if p++; p == pe {
			goto _test_eof103
		}
	st_case_103:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 95:
			goto st57
		case 116:
			goto st104
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st104:
		if p++; p == pe {
			goto _test_eof104
		}
	st_case_104:
		switch data[p] {
		case 9:
			goto tr149
		case 10:
			goto tr150
		case 32:
			goto tr149
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
tr152:
// line 51 "parsegl.rl"

 curline++ 
	goto st105
tr149:
// line 117 "parsegl.rl"

 m1 = p 
	goto st105
tr150:
// line 51 "parsegl.rl"

 curline++ 
// line 117 "parsegl.rl"

 m1 = p 
	goto st105
	st105:
		if p++; p == pe {
			goto _test_eof105
		}
	st_case_105:
// line 2766 "parsegl.go"
		switch data[p] {
		case 9:
			goto st105
		case 10:
			goto tr152
		case 32:
			goto st105
		case 42:
			goto tr76
		case 65:
			goto tr153
		case 71:
			goto tr154
		case 95:
			goto tr69
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr69
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr69
			}
		default:
			goto tr69
		}
		goto st0
tr153:
// line 117 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st106
	st106:
		if p++; p == pe {
			goto _test_eof106
		}
	st_case_106:
// line 2806 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 80:
			goto st107
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st107:
		if p++; p == pe {
			goto _test_eof107
		}
	st_case_107:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 73:
			goto st108
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st108:
		if p++; p == pe {
			goto _test_eof108
		}
	st_case_108:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 69:
			goto st109
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st109:
		if p++; p == pe {
			goto _test_eof109
		}
	st_case_109:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 78:
			goto st110
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st110:
		if p++; p == pe {
			goto _test_eof110
		}
	st_case_110:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 84:
			goto st111
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st111:
		if p++; p == pe {
			goto _test_eof111
		}
	st_case_111:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 82:
			goto st112
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st112:
		if p++; p == pe {
			goto _test_eof112
		}
	st_case_112:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 89:
			goto st113
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st113:
		if p++; p == pe {
			goto _test_eof113
		}
	st_case_113:
		switch data[p] {
		case 9:
			goto tr162
		case 10:
			goto tr163
		case 32:
			goto tr162
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
tr165:
// line 51 "parsegl.rl"

 curline++ 
	goto st114
tr162:
// line 117 "parsegl.rl"

 m1 = p 
	goto st114
tr163:
// line 117 "parsegl.rl"

 m1 = p 
// line 51 "parsegl.rl"

 curline++ 
	goto st114
	st114:
		if p++; p == pe {
			goto _test_eof114
		}
	st_case_114:
// line 3063 "parsegl.go"
		switch data[p] {
		case 9:
			goto st114
		case 10:
			goto tr165
		case 32:
			goto st114
		case 42:
			goto tr76
		case 65:
			goto st61
		case 71:
			goto st98
		case 103:
			goto tr90
		}
		goto st0
tr154:
// line 117 "parsegl.rl"

 m0 = p; m4 = 0 
	goto st115
	st115:
		if p++; p == pe {
			goto _test_eof115
		}
	st_case_115:
// line 3091 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 76:
			goto st116
		case 95:
			goto st57
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st116:
		if p++; p == pe {
			goto _test_eof116
		}
	st_case_116:
		switch data[p] {
		case 9:
			goto tr71
		case 10:
			goto tr72
		case 32:
			goto tr71
		case 65:
			goto st106
		case 95:
			goto st57
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st57
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st57
			}
		default:
			goto st57
		}
		goto st0
	st117:
		if p++; p == pe {
			goto _test_eof117
		}
	st_case_117:
		if data[p] == 120 {
			goto st118
		}
		goto st0
	st118:
		if p++; p == pe {
			goto _test_eof118
		}
	st_case_118:
		if data[p] == 116 {
			goto st119
		}
		goto st0
	st119:
		if p++; p == pe {
			goto _test_eof119
		}
	st_case_119:
		if data[p] == 101 {
			goto st120
		}
		goto st0
	st120:
		if p++; p == pe {
			goto _test_eof120
		}
	st_case_120:
		if data[p] == 114 {
			goto st121
		}
		goto st0
	st121:
		if p++; p == pe {
			goto _test_eof121
		}
	st_case_121:
		if data[p] == 110 {
			goto st122
		}
		goto st0
	st122:
		if p++; p == pe {
			goto _test_eof122
		}
	st_case_122:
		switch data[p] {
		case 9:
			goto st123
		case 32:
			goto st123
		}
		goto st0
	st123:
		if p++; p == pe {
			goto _test_eof123
		}
	st_case_123:
		switch data[p] {
		case 9:
			goto st123
		case 32:
			goto st123
		case 34:
			goto st124
		}
		goto st0
	st124:
		if p++; p == pe {
			goto _test_eof124
		}
	st_case_124:
		if data[p] == 67 {
			goto st125
		}
		goto st0
	st125:
		if p++; p == pe {
			goto _test_eof125
		}
	st_case_125:
		if data[p] == 34 {
			goto st126
		}
		goto st0
	st126:
		if p++; p == pe {
			goto _test_eof126
		}
	st_case_126:
		switch data[p] {
		case 9:
			goto st127
		case 32:
			goto st127
		}
		goto st0
	st127:
		if p++; p == pe {
			goto _test_eof127
		}
	st_case_127:
		switch data[p] {
		case 9:
			goto st127
		case 32:
			goto st127
		case 123:
			goto st128
		}
		goto st0
	st128:
		if p++; p == pe {
			goto _test_eof128
		}
	st_case_128:
		switch data[p] {
		case 9:
			goto st128
		case 10:
			goto tr179
		case 32:
			goto st128
		}
		goto st0
	st129:
		if p++; p == pe {
			goto _test_eof129
		}
	st_case_129:
		if data[p] == 121 {
			goto st130
		}
		goto st0
	st130:
		if p++; p == pe {
			goto _test_eof130
		}
	st_case_130:
		if data[p] == 112 {
			goto st131
		}
		goto st0
	st131:
		if p++; p == pe {
			goto _test_eof131
		}
	st_case_131:
		if data[p] == 101 {
			goto st132
		}
		goto st0
	st132:
		if p++; p == pe {
			goto _test_eof132
		}
	st_case_132:
		if data[p] == 100 {
			goto st133
		}
		goto st0
	st133:
		if p++; p == pe {
			goto _test_eof133
		}
	st_case_133:
		if data[p] == 101 {
			goto st134
		}
		goto st0
	st134:
		if p++; p == pe {
			goto _test_eof134
		}
	st_case_134:
		if data[p] == 102 {
			goto st135
		}
		goto st0
	st135:
		if p++; p == pe {
			goto _test_eof135
		}
	st_case_135:
		switch data[p] {
		case 9:
			goto st147
		case 32:
			goto st147
		case 59:
			goto st0
		}
		goto st136
	st136:
		if p++; p == pe {
			goto _test_eof136
		}
	st_case_136:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		}
		goto st136
	st137:
		if p++; p == pe {
			goto _test_eof137
		}
	st_case_137:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 80:
			goto st138
		}
		goto st136
	st138:
		if p++; p == pe {
			goto _test_eof138
		}
	st_case_138:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 73:
			goto st139
		}
		goto st136
	st139:
		if p++; p == pe {
			goto _test_eof139
		}
	st_case_139:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 69:
			goto st140
		}
		goto st136
	st140:
		if p++; p == pe {
			goto _test_eof140
		}
	st_case_140:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 78:
			goto st141
		}
		goto st136
	st141:
		if p++; p == pe {
			goto _test_eof141
		}
	st_case_141:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 84:
			goto st142
		}
		goto st136
	st142:
		if p++; p == pe {
			goto _test_eof142
		}
	st_case_142:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 82:
			goto st143
		}
		goto st136
	st143:
		if p++; p == pe {
			goto _test_eof143
		}
	st_case_143:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 89:
			goto st144
		}
		goto st136
	st144:
		if p++; p == pe {
			goto _test_eof144
		}
	st_case_144:
		switch data[p] {
		case 59:
			goto st0
		case 65:
			goto st137
		case 80:
			goto st145
		}
		goto st136
	st145:
		if p++; p == pe {
			goto _test_eof145
		}
	st_case_145:
		if data[p] == 59 {
			goto st0
		}
		goto st146
	st146:
		if p++; p == pe {
			goto _test_eof146
		}
	st_case_146:
		if data[p] == 59 {
			goto tr198
		}
		goto st146
	st147:
		if p++; p == pe {
			goto _test_eof147
		}
	st_case_147:
		switch data[p] {
		case 9:
			goto st147
		case 32:
			goto st147
		case 59:
			goto st0
		case 65:
			goto tr200
		case 95:
			goto tr199
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr199
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr199
			}
		default:
			goto tr199
		}
		goto st136
tr199:
// line 105 "parsegl.rl"

 m0 = p 
	goto st148
	st148:
		if p++; p == pe {
			goto _test_eof148
		}
	st_case_148:
// line 3526 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
tr201:
// line 105 "parsegl.rl"

 m1 = p 
	goto st149
	st149:
		if p++; p == pe {
			goto _test_eof149
		}
	st_case_149:
// line 3564 "parsegl.go"
		switch data[p] {
		case 9:
			goto st150
		case 32:
			goto st150
		case 59:
			goto st0
		case 65:
			goto tr207
		case 95:
			goto tr206
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr206
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr206
			}
		default:
			goto tr206
		}
		goto st136
tr233:
// line 105 "parsegl.rl"

 m1 = p 
	goto st150
	st150:
		if p++; p == pe {
			goto _test_eof150
		}
	st_case_150:
// line 3600 "parsegl.go"
		switch data[p] {
		case 9:
			goto st150
		case 32:
			goto st150
		case 59:
			goto st0
		case 65:
			goto tr209
		case 95:
			goto tr208
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr208
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr208
			}
		default:
			goto tr208
		}
		goto st136
tr208:
// line 105 "parsegl.rl"

 m2 = p 
	goto st151
	st151:
		if p++; p == pe {
			goto _test_eof151
		}
	st_case_151:
// line 3636 "parsegl.go"
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
tr211:
// line 1 "NONE"

te = p+1

// line 105 "parsegl.rl"

 m3 = p; m4, m5 = 0, 0 
	goto st196
	st196:
		if p++; p == pe {
			goto _test_eof196
		}
	st_case_196:
// line 3672 "parsegl.go"
		switch data[p] {
		case 9:
			goto st152
		case 32:
			goto st152
		case 47:
			goto tr215
		}
		goto tr269
	st152:
		if p++; p == pe {
			goto _test_eof152
		}
	st_case_152:
		switch data[p] {
		case 9:
			goto st152
		case 32:
			goto st152
		case 47:
			goto tr215
		}
		goto tr213
tr215:
// line 106 "parsegl.rl"

 m4 = p+2 
	goto st153
	st153:
		if p++; p == pe {
			goto _test_eof153
		}
	st_case_153:
// line 3706 "parsegl.go"
		if data[p] == 42 {
			goto st154
		}
		goto tr213
tr217:
// line 51 "parsegl.rl"

 curline++ 
	goto st154
	st154:
		if p++; p == pe {
			goto _test_eof154
		}
	st_case_154:
// line 3721 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr217
		case 42:
			goto st155
		}
		goto st154
	st155:
		if p++; p == pe {
			goto _test_eof155
		}
	st_case_155:
		switch data[p] {
		case 10:
			goto tr217
		case 42:
			goto st155
		case 47:
			goto tr219
		}
		goto st154
tr219:
// line 106 "parsegl.rl"

 m5 = p-2 
	goto st156
	st156:
		if p++; p == pe {
			goto _test_eof156
		}
	st_case_156:
// line 3753 "parsegl.go"
		switch data[p] {
		case 9:
			goto st156
		case 10:
			goto tr221
		case 32:
			goto st156
		}
		goto tr213
tr209:
// line 105 "parsegl.rl"

 m2 = p 
	goto st157
	st157:
		if p++; p == pe {
			goto _test_eof157
		}
	st_case_157:
// line 3773 "parsegl.go"
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 80:
			goto st158
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st158:
		if p++; p == pe {
			goto _test_eof158
		}
	st_case_158:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 73:
			goto st159
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st159:
		if p++; p == pe {
			goto _test_eof159
		}
	st_case_159:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 69:
			goto st160
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st160:
		if p++; p == pe {
			goto _test_eof160
		}
	st_case_160:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 78:
			goto st161
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st161:
		if p++; p == pe {
			goto _test_eof161
		}
	st_case_161:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 84:
			goto st162
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st162:
		if p++; p == pe {
			goto _test_eof162
		}
	st_case_162:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 82:
			goto st163
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st163:
		if p++; p == pe {
			goto _test_eof163
		}
	st_case_163:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 89:
			goto st164
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st164:
		if p++; p == pe {
			goto _test_eof164
		}
	st_case_164:
		switch data[p] {
		case 59:
			goto tr211
		case 65:
			goto st157
		case 80:
			goto st165
		case 95:
			goto st151
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st151
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st151
			}
		default:
			goto st151
		}
		goto st136
	st165:
		if p++; p == pe {
			goto _test_eof165
		}
	st_case_165:
		switch data[p] {
		case 59:
			goto tr211
		case 95:
			goto st166
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st166
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st166
			}
		default:
			goto st166
		}
		goto st146
tr247:
// line 105 "parsegl.rl"

 m2 = p 
	goto st166
	st166:
		if p++; p == pe {
			goto _test_eof166
		}
	st_case_166:
// line 4027 "parsegl.go"
		switch data[p] {
		case 59:
			goto tr211
		case 95:
			goto st166
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st166
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st166
			}
		default:
			goto st166
		}
		goto st146
tr206:
// line 105 "parsegl.rl"

 m2 = p 
	goto st167
	st167:
		if p++; p == pe {
			goto _test_eof167
		}
	st_case_167:
// line 4057 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st168:
		if p++; p == pe {
			goto _test_eof168
		}
	st_case_168:
		switch data[p] {
		case 9:
			goto tr233
		case 32:
			goto tr233
		case 59:
			goto st0
		case 65:
			goto st137
		}
		goto st136
tr207:
// line 105 "parsegl.rl"

 m2 = p 
	goto st169
	st169:
		if p++; p == pe {
			goto _test_eof169
		}
	st_case_169:
// line 4111 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 80:
			goto st170
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st170:
		if p++; p == pe {
			goto _test_eof170
		}
	st_case_170:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 73:
			goto st171
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st171:
		if p++; p == pe {
			goto _test_eof171
		}
	st_case_171:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 69:
			goto st172
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st172:
		if p++; p == pe {
			goto _test_eof172
		}
	st_case_172:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 78:
			goto st173
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st173:
		if p++; p == pe {
			goto _test_eof173
		}
	st_case_173:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 84:
			goto st174
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st174:
		if p++; p == pe {
			goto _test_eof174
		}
	st_case_174:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 82:
			goto st175
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st175:
		if p++; p == pe {
			goto _test_eof175
		}
	st_case_175:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 89:
			goto st176
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st176:
		if p++; p == pe {
			goto _test_eof176
		}
	st_case_176:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto tr211
		case 65:
			goto st169
		case 80:
			goto st177
		case 95:
			goto st167
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st167
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st167
			}
		default:
			goto st167
		}
		goto st136
	st177:
		if p++; p == pe {
			goto _test_eof177
		}
	st_case_177:
		switch data[p] {
		case 9:
			goto tr242
		case 32:
			goto tr242
		case 42:
			goto st181
		case 59:
			goto tr211
		case 95:
			goto st180
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st180
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st180
			}
		default:
			goto st180
		}
		goto st146
tr242:
// line 105 "parsegl.rl"

 m1 = p 
	goto st178
	st178:
		if p++; p == pe {
			goto _test_eof178
		}
	st_case_178:
// line 4419 "parsegl.go"
		switch data[p] {
		case 9:
			goto st179
		case 32:
			goto st179
		case 59:
			goto tr198
		case 95:
			goto tr246
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr246
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr246
			}
		default:
			goto tr246
		}
		goto st146
tr248:
// line 105 "parsegl.rl"

 m1 = p 
	goto st179
	st179:
		if p++; p == pe {
			goto _test_eof179
		}
	st_case_179:
// line 4453 "parsegl.go"
		switch data[p] {
		case 9:
			goto st179
		case 32:
			goto st179
		case 59:
			goto tr198
		case 95:
			goto tr247
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto tr247
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto tr247
			}
		default:
			goto tr247
		}
		goto st146
tr246:
// line 105 "parsegl.rl"

 m2 = p 
	goto st180
	st180:
		if p++; p == pe {
			goto _test_eof180
		}
	st_case_180:
// line 4487 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr242
		case 32:
			goto tr242
		case 42:
			goto st181
		case 59:
			goto tr211
		case 95:
			goto st180
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st180
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st180
			}
		default:
			goto st180
		}
		goto st146
	st181:
		if p++; p == pe {
			goto _test_eof181
		}
	st_case_181:
		switch data[p] {
		case 9:
			goto tr248
		case 32:
			goto tr248
		case 59:
			goto tr198
		}
		goto st146
tr200:
// line 105 "parsegl.rl"

 m0 = p 
	goto st182
	st182:
		if p++; p == pe {
			goto _test_eof182
		}
	st_case_182:
// line 4537 "parsegl.go"
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 80:
			goto st183
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st183:
		if p++; p == pe {
			goto _test_eof183
		}
	st_case_183:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 73:
			goto st184
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st184:
		if p++; p == pe {
			goto _test_eof184
		}
	st_case_184:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 69:
			goto st185
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st185:
		if p++; p == pe {
			goto _test_eof185
		}
	st_case_185:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 78:
			goto st186
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st186:
		if p++; p == pe {
			goto _test_eof186
		}
	st_case_186:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 84:
			goto st187
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st187:
		if p++; p == pe {
			goto _test_eof187
		}
	st_case_187:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 82:
			goto st188
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st188:
		if p++; p == pe {
			goto _test_eof188
		}
	st_case_188:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 89:
			goto st189
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st189:
		if p++; p == pe {
			goto _test_eof189
		}
	st_case_189:
		switch data[p] {
		case 9:
			goto tr201
		case 32:
			goto tr201
		case 42:
			goto st168
		case 59:
			goto st0
		case 65:
			goto st182
		case 80:
			goto st190
		case 95:
			goto st148
		}
		switch {
		case data[p] < 66:
			if 48 <= data[p] && data[p] <= 57 {
				goto st148
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st148
			}
		default:
			goto st148
		}
		goto st136
	st190:
		if p++; p == pe {
			goto _test_eof190
		}
	st_case_190:
		switch data[p] {
		case 9:
			goto tr242
		case 32:
			goto tr242
		case 42:
			goto st181
		case 59:
			goto st0
		case 95:
			goto st191
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st191
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st191
			}
		default:
			goto st191
		}
		goto st146
	st191:
		if p++; p == pe {
			goto _test_eof191
		}
	st_case_191:
		switch data[p] {
		case 9:
			goto tr242
		case 32:
			goto tr242
		case 42:
			goto st181
		case 59:
			goto tr198
		case 95:
			goto st191
		}
		switch {
		case data[p] < 65:
			if 48 <= data[p] && data[p] <= 57 {
				goto st191
			}
		case data[p] > 90:
			if 97 <= data[p] && data[p] <= 122 {
				goto st191
			}
		default:
			goto st191
		}
		goto st146
	st192:
// line 1 "NONE"

ts = 0

		if p++; p == pe {
			goto _test_eof192
		}
	st_case_192:
// line 4874 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr259
		case 59:
			goto tr260
		}
		goto st192
tr259:
// line 51 "parsegl.rl"

 curline++ 
// line 145 "parsegl.rl"

 {goto st193 } 
	goto st197
	st197:
		if p++; p == pe {
			goto _test_eof197
		}
	st_case_197:
// line 4895 "parsegl.go"
		goto st0
tr260:
// line 145 "parsegl.rl"

 {goto st193 } 
	goto st198
	st198:
		if p++; p == pe {
			goto _test_eof198
		}
	st_case_198:
// line 4907 "parsegl.go"
		switch data[p] {
		case 10:
			goto tr259
		case 59:
			goto tr260
		}
		goto st192
	st_out:
	_test_eof193: cs = 193; goto _test_eof
	_test_eof1: cs = 1; goto _test_eof
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof
	_test_eof4: cs = 4; goto _test_eof
	_test_eof194: cs = 194; goto _test_eof
	_test_eof5: cs = 5; goto _test_eof
	_test_eof6: cs = 6; goto _test_eof
	_test_eof7: cs = 7; goto _test_eof
	_test_eof8: cs = 8; goto _test_eof
	_test_eof9: cs = 9; goto _test_eof
	_test_eof10: cs = 10; goto _test_eof
	_test_eof11: cs = 11; goto _test_eof
	_test_eof12: cs = 12; goto _test_eof
	_test_eof13: cs = 13; goto _test_eof
	_test_eof14: cs = 14; goto _test_eof
	_test_eof15: cs = 15; goto _test_eof
	_test_eof16: cs = 16; goto _test_eof
	_test_eof17: cs = 17; goto _test_eof
	_test_eof18: cs = 18; goto _test_eof
	_test_eof19: cs = 19; goto _test_eof
	_test_eof20: cs = 20; goto _test_eof
	_test_eof21: cs = 21; goto _test_eof
	_test_eof22: cs = 22; goto _test_eof
	_test_eof23: cs = 23; goto _test_eof
	_test_eof24: cs = 24; goto _test_eof
	_test_eof25: cs = 25; goto _test_eof
	_test_eof26: cs = 26; goto _test_eof
	_test_eof27: cs = 27; goto _test_eof
	_test_eof28: cs = 28; goto _test_eof
	_test_eof29: cs = 29; goto _test_eof
	_test_eof30: cs = 30; goto _test_eof
	_test_eof31: cs = 31; goto _test_eof
	_test_eof32: cs = 32; goto _test_eof
	_test_eof33: cs = 33; goto _test_eof
	_test_eof34: cs = 34; goto _test_eof
	_test_eof35: cs = 35; goto _test_eof
	_test_eof36: cs = 36; goto _test_eof
	_test_eof37: cs = 37; goto _test_eof
	_test_eof38: cs = 38; goto _test_eof
	_test_eof39: cs = 39; goto _test_eof
	_test_eof40: cs = 40; goto _test_eof
	_test_eof41: cs = 41; goto _test_eof
	_test_eof42: cs = 42; goto _test_eof
	_test_eof43: cs = 43; goto _test_eof
	_test_eof44: cs = 44; goto _test_eof
	_test_eof45: cs = 45; goto _test_eof
	_test_eof46: cs = 46; goto _test_eof
	_test_eof47: cs = 47; goto _test_eof
	_test_eof48: cs = 48; goto _test_eof
	_test_eof49: cs = 49; goto _test_eof
	_test_eof50: cs = 50; goto _test_eof
	_test_eof195: cs = 195; goto _test_eof
	_test_eof51: cs = 51; goto _test_eof
	_test_eof52: cs = 52; goto _test_eof
	_test_eof53: cs = 53; goto _test_eof
	_test_eof54: cs = 54; goto _test_eof
	_test_eof55: cs = 55; goto _test_eof
	_test_eof56: cs = 56; goto _test_eof
	_test_eof57: cs = 57; goto _test_eof
	_test_eof58: cs = 58; goto _test_eof
	_test_eof59: cs = 59; goto _test_eof
	_test_eof60: cs = 60; goto _test_eof
	_test_eof61: cs = 61; goto _test_eof
	_test_eof62: cs = 62; goto _test_eof
	_test_eof63: cs = 63; goto _test_eof
	_test_eof64: cs = 64; goto _test_eof
	_test_eof65: cs = 65; goto _test_eof
	_test_eof66: cs = 66; goto _test_eof
	_test_eof67: cs = 67; goto _test_eof
	_test_eof68: cs = 68; goto _test_eof
	_test_eof69: cs = 69; goto _test_eof
	_test_eof70: cs = 70; goto _test_eof
	_test_eof71: cs = 71; goto _test_eof
	_test_eof72: cs = 72; goto _test_eof
	_test_eof73: cs = 73; goto _test_eof
	_test_eof74: cs = 74; goto _test_eof
	_test_eof75: cs = 75; goto _test_eof
	_test_eof76: cs = 76; goto _test_eof
	_test_eof77: cs = 77; goto _test_eof
	_test_eof78: cs = 78; goto _test_eof
	_test_eof79: cs = 79; goto _test_eof
	_test_eof80: cs = 80; goto _test_eof
	_test_eof81: cs = 81; goto _test_eof
	_test_eof82: cs = 82; goto _test_eof
	_test_eof83: cs = 83; goto _test_eof
	_test_eof84: cs = 84; goto _test_eof
	_test_eof85: cs = 85; goto _test_eof
	_test_eof86: cs = 86; goto _test_eof
	_test_eof87: cs = 87; goto _test_eof
	_test_eof88: cs = 88; goto _test_eof
	_test_eof89: cs = 89; goto _test_eof
	_test_eof90: cs = 90; goto _test_eof
	_test_eof91: cs = 91; goto _test_eof
	_test_eof92: cs = 92; goto _test_eof
	_test_eof93: cs = 93; goto _test_eof
	_test_eof94: cs = 94; goto _test_eof
	_test_eof95: cs = 95; goto _test_eof
	_test_eof96: cs = 96; goto _test_eof
	_test_eof97: cs = 97; goto _test_eof
	_test_eof98: cs = 98; goto _test_eof
	_test_eof99: cs = 99; goto _test_eof
	_test_eof100: cs = 100; goto _test_eof
	_test_eof101: cs = 101; goto _test_eof
	_test_eof102: cs = 102; goto _test_eof
	_test_eof103: cs = 103; goto _test_eof
	_test_eof104: cs = 104; goto _test_eof
	_test_eof105: cs = 105; goto _test_eof
	_test_eof106: cs = 106; goto _test_eof
	_test_eof107: cs = 107; goto _test_eof
	_test_eof108: cs = 108; goto _test_eof
	_test_eof109: cs = 109; goto _test_eof
	_test_eof110: cs = 110; goto _test_eof
	_test_eof111: cs = 111; goto _test_eof
	_test_eof112: cs = 112; goto _test_eof
	_test_eof113: cs = 113; goto _test_eof
	_test_eof114: cs = 114; goto _test_eof
	_test_eof115: cs = 115; goto _test_eof
	_test_eof116: cs = 116; goto _test_eof
	_test_eof117: cs = 117; goto _test_eof
	_test_eof118: cs = 118; goto _test_eof
	_test_eof119: cs = 119; goto _test_eof
	_test_eof120: cs = 120; goto _test_eof
	_test_eof121: cs = 121; goto _test_eof
	_test_eof122: cs = 122; goto _test_eof
	_test_eof123: cs = 123; goto _test_eof
	_test_eof124: cs = 124; goto _test_eof
	_test_eof125: cs = 125; goto _test_eof
	_test_eof126: cs = 126; goto _test_eof
	_test_eof127: cs = 127; goto _test_eof
	_test_eof128: cs = 128; goto _test_eof
	_test_eof129: cs = 129; goto _test_eof
	_test_eof130: cs = 130; goto _test_eof
	_test_eof131: cs = 131; goto _test_eof
	_test_eof132: cs = 132; goto _test_eof
	_test_eof133: cs = 133; goto _test_eof
	_test_eof134: cs = 134; goto _test_eof
	_test_eof135: cs = 135; goto _test_eof
	_test_eof136: cs = 136; goto _test_eof
	_test_eof137: cs = 137; goto _test_eof
	_test_eof138: cs = 138; goto _test_eof
	_test_eof139: cs = 139; goto _test_eof
	_test_eof140: cs = 140; goto _test_eof
	_test_eof141: cs = 141; goto _test_eof
	_test_eof142: cs = 142; goto _test_eof
	_test_eof143: cs = 143; goto _test_eof
	_test_eof144: cs = 144; goto _test_eof
	_test_eof145: cs = 145; goto _test_eof
	_test_eof146: cs = 146; goto _test_eof
	_test_eof147: cs = 147; goto _test_eof
	_test_eof148: cs = 148; goto _test_eof
	_test_eof149: cs = 149; goto _test_eof
	_test_eof150: cs = 150; goto _test_eof
	_test_eof151: cs = 151; goto _test_eof
	_test_eof196: cs = 196; goto _test_eof
	_test_eof152: cs = 152; goto _test_eof
	_test_eof153: cs = 153; goto _test_eof
	_test_eof154: cs = 154; goto _test_eof
	_test_eof155: cs = 155; goto _test_eof
	_test_eof156: cs = 156; goto _test_eof
	_test_eof157: cs = 157; goto _test_eof
	_test_eof158: cs = 158; goto _test_eof
	_test_eof159: cs = 159; goto _test_eof
	_test_eof160: cs = 160; goto _test_eof
	_test_eof161: cs = 161; goto _test_eof
	_test_eof162: cs = 162; goto _test_eof
	_test_eof163: cs = 163; goto _test_eof
	_test_eof164: cs = 164; goto _test_eof
	_test_eof165: cs = 165; goto _test_eof
	_test_eof166: cs = 166; goto _test_eof
	_test_eof167: cs = 167; goto _test_eof
	_test_eof168: cs = 168; goto _test_eof
	_test_eof169: cs = 169; goto _test_eof
	_test_eof170: cs = 170; goto _test_eof
	_test_eof171: cs = 171; goto _test_eof
	_test_eof172: cs = 172; goto _test_eof
	_test_eof173: cs = 173; goto _test_eof
	_test_eof174: cs = 174; goto _test_eof
	_test_eof175: cs = 175; goto _test_eof
	_test_eof176: cs = 176; goto _test_eof
	_test_eof177: cs = 177; goto _test_eof
	_test_eof178: cs = 178; goto _test_eof
	_test_eof179: cs = 179; goto _test_eof
	_test_eof180: cs = 180; goto _test_eof
	_test_eof181: cs = 181; goto _test_eof
	_test_eof182: cs = 182; goto _test_eof
	_test_eof183: cs = 183; goto _test_eof
	_test_eof184: cs = 184; goto _test_eof
	_test_eof185: cs = 185; goto _test_eof
	_test_eof186: cs = 186; goto _test_eof
	_test_eof187: cs = 187; goto _test_eof
	_test_eof188: cs = 188; goto _test_eof
	_test_eof189: cs = 189; goto _test_eof
	_test_eof190: cs = 190; goto _test_eof
	_test_eof191: cs = 191; goto _test_eof
	_test_eof192: cs = 192; goto _test_eof
	_test_eof197: cs = 197; goto _test_eof
	_test_eof198: cs = 198; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch cs {
		case 194:
			goto tr266
		case 195:
			goto tr266
		case 196:
			goto tr269
		case 152:
			goto tr213
		case 153:
			goto tr213
		case 154:
			goto tr213
		case 155:
			goto tr213
		case 156:
			goto tr213
		}
	}

	_out: {}
	}

// line 149 "parsegl.rl"


	if p < pe {
		m0, m1 = p, p
		for m0 > 0 && data[m0-1] != '\n' {
			m0--
		}
		for m1 < len(data) && data[m1] != '\n' {
			m1++
		}
		return fmt.Errorf("cannot parse header file:%d:%d: %s\n", curline, p-m0, data[m0:m1])
	}
	return nil
}
