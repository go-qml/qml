//
// To compile:
//
//   ragel -Z -G2 -o parsegl.go parsegl.rl
//
// To show a diagram of the state machine:
//
//   ragel -V -G2 -p -o parsegl.dot parsegl.rl
//   dot -Tsvg -o parsegl.svg parsegl.dot
//   chrome parsegl.svg
//

package main

import (
	"fmt"
	"strings"
)

%%{
	machine parsegl;

	write data;
}%%

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

	%%{
		nl = '\n' @{ curline++ };
		cd = [^\n];
		sp = [ \t];
		id = [A-Za-z0-9_]+;
		spnl = ( sp | nl );

		comment = '/*' ( cd | nl )* :>> '*/';

		main := |*
			# Track heading comments.
			'/*' @{ m0 = p } ( [^\n] | nl )* :>> '*/' @{ m1 = p-2 } sp* nl
			{
				heading = strings.Trim(data[m0:m1], " *")
				lineblock++
			};

			# Ignore other comments.
			sp* comment sp*;

			# Ignore pragmas and includes.
			( '#pragma' | '#include' ) cd+ nl;

			# Ignore the extern declaration.
			( 'extern' sp+ '"C"' sp+ '{' | '}' ) sp* nl;

			# Track conditional blocks, potentially nested.
			'#' sp* 'if' ('def' | 'ndef' | '') sp [^\n]* nl
			{
				ifblock++
			};
			'#' sp* 'endif' [^\n]* nl
			{
				ifblock--
			};

			# Ignore elses. Both sides are within the ifblock.
			'#' sp* ( 'else' | 'elif' ) cd* nl;

			# Record defines.
			'#' sp* 'define' sp+
				# Key
				id >{ m0 = p; m2, m3 = 0, 0 } @{ m1 = p+1 }
				# Value
				( sp+ ( [^\n ] ( sp? [^\n ] )* ) >{ m2 = p } @{ m3 = p+1 } )?
				sp* nl
			{
				// Ignore defines within conditional blocks other than the top level ifndef.
				if ifblock == 1 && strings.HasPrefix(data[m0:m1], "GL_") {
					header.Define = append(header.Define, Define{Name: data[m0:m1], Value: data[m2:m3], Heading: heading, LineBlock: lineblock})
					heading = ""
				}
			};

			# Record typedefs.
			'typedef' sp+ ( id+ ( sp id+ )* '*'? ) >{ m0 = p } sp+ >{ m1 = p } id+ >{ m2 = p } ';' >{ m3 = p; m4, m5 = 0, 0 }
				( sp* ( comment ) >{ m4 = p+2 } @{ m5 = p-2 } sp* nl )?
			{
				if (ifblock == 1) {
					header.Type = append(header.Type, Type{Name: data[m2:m3], Type: data[m0:m1], Comment: strings.TrimSpace(data[m4:m5])})
				}
			};

			# Ignore typedefs with function types.
			'typedef' [^;]+ 'APIENTRYP' [^;]+ ';';

			# Record function prototypes.
			'GLAPI' spnl+ ( 'const' spnl+ )? id >{ m0 = p; m4 = 0 } spnl+ >{ m1 = p } ( '*'+ ${ m4++ } spnl* )? 'GL'? 'APIENTRY' spnl+
				# Name
				'gl' >{ m2 = p } id ( spnl* '(' ) >{ m3 = p; f = Func{Name: data[m2:m3], Type: data[m0:m1], Addr: m4} } spnl*
				# Parameters
				( 'void' spnl* ')' | ( ( 'const' spnl+ )? id >{ m0 = p; m4 = 0 } spnl+ >{ m1 = p } ( '*'+ ${ m4++ } spnl* )? id >{ m2 = p; m5 = 0 } @{ m3 = p+1 } ( '[' [0-9]+ ${ m5 = m5*10 + (int(data[p]) - '0') } ']' )? spnl* [,)]
					>{ f.Param = append(f.Param, Param{Name: data[m2:m3], Type: data[m0:m1], Addr: m4, Array: m5}) } spnl* )+ )
				spnl* ';'
			{
				if (ifblock == 1) {
					header.Func = append(header.Func, f)
				}
			};

			# Reset relevant states on empty lines.
			sp* nl
			{
				// Reset heading comment.
				heading = ""

				// Start new line block.
				lineblock++
			};

		*|;

		skiperror := [^\n]* (';' | nl ) @{ fgoto main; };

		write init;
		write exec;
	}%%

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
