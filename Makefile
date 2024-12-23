BUILD=build/

all:
	go build -o emulator_gui main.go

run: all
	./emulator_gui
clean:
	rm emulator_gui