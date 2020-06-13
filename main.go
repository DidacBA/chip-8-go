package main

import (
	"fmt"
	"os"
)

var FontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type CPU struct {
	Memory    [4096]byte
	Registers [16]byte
	Stack     [16]uint16
	I         uint16 //
	PC        uint16 // Program Counter
	SP        byte   // Stack Pointer
	ST        byte   // Sound Timer
	DT        byte   // Delay Timer
}

func (cpu *CPU) Load(rom []byte) error {
	// Load memory with rom. Starts loading at the initial position of the program counter: 0x200 (512)

	//Check if rom is bigger than available memory
	if int32(len(cpu.Memory)-512) < int32(len(rom)) {
		return fmt.Errorf("Rom size (%d) bigger than available memory (%d)", int32(len(cpu.Memory)-512), int32(len(rom)))
	}

	for i := 0; i < len(rom); i++ {
		cpu.Memory[i+512] = rom[i]
	}

	return nil
}

func (cpu *CPU) Step() {
	// Instruction step (fetch, decode, execute)
}

func (cpu *CPU) Fetch() uint16 {
	// Return current opcode
	// TODO: add check for out of bounds memory
	return uint16(cpu.Memory[cpu.PC])<<8 | uint16(cpu.Memory[cpu.PC+1])
}

func (cpu *CPU) DecodeOpCode() {
	// Disassemble opCode
}

func (cpu *CPU) ExecuteOP() {
	// Execute Instruction
}

func (cpu *CPU) LoadFontSet() {
	for i := 0; i < len(FontSet); i++ {
		cpu.Memory[i] = FontSet[i]
	}
}

func (cpu *CPU) Reset() {
	// New CPU and set it through the pointer cpu
	newCpu := &CPU{
		PC: 0x200,
	}
	*cpu = *newCpu
}

func openRomFile(path string) ([]byte, error) {
	// Opens rom file and returns a byte slice
	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("Error opening the file", path)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error reading stats of the file", err)
		return nil, err
	}

	romBuffer := make([]byte, fileInfo.Size())
	if _, err := file.Read(romBuffer); err != nil {
		return nil, err
	}

	return romBuffer, nil
}

func main() {
	fmt.Println("Init Chip 8 interpreter")

	cpu := &CPU{
		PC: 0x200,
	}

	rompath := os.Args[1]

	cpu.LoadFontSet()
	rom, err := openRomFile(rompath)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	cpu.Load(rom)
	cpu.Reset()
}
