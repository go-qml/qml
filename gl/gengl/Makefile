
gengl: main.go parsegl.go
	go build

parsegl.go: parsegl.rl
	ragel -Z -G2 -o parsegl.go parsegl.rl

run: gengl
	./gengl
