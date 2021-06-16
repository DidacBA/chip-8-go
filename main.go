package main

import (
	"fmt"
	"os"
	"strconv"
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
	I         uint16 // Index register
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

func (cpu *CPU) Step() uint16 {
	// Instruction step (fetch, execute)
	opcode := cpu.Fetch()
	cpu.Execute(opcode)
	return opcode
}

func (cpu *CPU) Fetch() uint16 {
	// Return current opcode
	// TODO: add check for out of bounds memory
	return uint16(cpu.Memory[cpu.PC])<<8 | uint16(cpu.Memory[cpu.PC+1])
}

func (cpu *CPU) Execute(opcode uint16) error {
	// Execute Instruction
	// fmt.Println(strconv.FormatUint(uint64(opcode&0xF000), 16))
	switch opcode & 0xF000 {
	case 0x000:
		switch opcode {
		case 0x00E0:
			// CLS
			// Clear the display
			cpu.PC += 2
			break
		case 0x00EE:
			// RET
			// Return from a subroutine
			cpu.PC = cpu.Stack[15]
			cpu.SP -= 1
			cpu.PC += 2
			break
		}
	case 0x1000:
		// JP addr
		// Jump to location nnn
		cpu.PC = opcode & 0x0FFF
		break
	case 0x2000:
		// Call addr
		// Call subroutine at nnn
		cpu.SP += 1
		cpu.Stack[15] = cpu.PC
		cpu.PC = opcode & 0x0FFF
		break
	case 0x3000:
		// SE Vx, byte
		// Skip next instruction if Vx = kk
		cpu.PC += 2
		if cpu.Registers[(opcode&0x0F00)>>8] == byte(opcode&0x00FF) {
			cpu.PC += 2
		}
		break
	case 0x4000:
		// SNE Vx, byte
		// Skip next instruction if Vx != kk
		cpu.PC += 2
		if cpu.Registers[(opcode&0x0F00)>>8] != byte(opcode&0x00FF) {
			cpu.PC += 2
		}
		break
	case 0x5000:
		// SE Vx, Vy
		// Skip next instruction if Vx = Vy
		cpu.PC += 2
		if cpu.Registers[(opcode&0x0F00)>>8] == cpu.Registers[(opcode&0x00F0)>>4] {
			cpu.PC += 2
		}
		break
	case 0x6000:
		// LD Vx, byte
		// Set Vx == kk
		cpu.Registers[(opcode&0x0F00)>>8] = byte(opcode & 0x00FF)
		cpu.PC += 2
		break
	case 0x7000:
		// Add Vx, byte
		cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] + byte(opcode&0x00FF)
		cpu.PC += 2
		break
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// Set Vx = Vy
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x00F0)>>4]
			cpu.PC += 2
			break
		case 0x0001:
			// OR Vx, Vy
			// Set Vx = Vx OR Vy
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] | cpu.Registers[(opcode&0x00F0)>>4]
			cpu.PC += 2
			break
		case 0x0002:
			// AND Vx, Vy
			// Set Vx = Vx AND Vy
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] & cpu.Registers[(opcode&0x00F0)>>4]
			break
		case 0x0003:
			// XOR Vx, Vy
			// Set Vx = Vx XOR Vy
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] ^ cpu.Registers[(opcode&0x00F0)>>4]
			cpu.PC += 2
			break
		case 0x0004:
			// ADD Vx, Vy
			sum := cpu.Registers[(opcode&0x0F00)>>8] + cpu.Registers[(opcode&0x00F0)>>4]
			if sum > 255 {
				cpu.Registers[0xF] = 1
			} else {
				cpu.Registers[0xF] = 0
			}
			cpu.Registers[(opcode&0x0F00)>>8] = sum & 0xFF
			cpu.PC += 2
			break
		case 0x0005:
			// SUB Vx, Vy
			if cpu.Registers[(opcode&0x0F00)>>8] > cpu.Registers[(opcode&0x00F0)>>4] {
				cpu.Registers[0xF] = 1
			} else {
				cpu.Registers[0xF] = 0
			}
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] - cpu.Registers[(opcode&0x00F0)>>4]
			cpu.PC += 2
			break
		case 0x0006:
			// SHR Vx {, Vy}
			if cpu.Registers[(opcode&0x0F00)>>8]&0x1 == 1 {
				cpu.Registers[0xF] = 1
			} else {
				cpu.Registers[0xF] = 0
			}
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] >> 1
			break
		case 0x0007:
			// SUBN Vx, Vy
			// Set Vx = Vy - Vx, set VF = NOT borrow
			if cpu.Registers[(opcode&0x00F0)>>4] > cpu.Registers[(opcode&0x0F00)>>8] {
				cpu.Registers[0xF] = 1
			} else {
				cpu.Registers[0xF] = 0
			}
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x00F0)>>4] - cpu.Registers[(opcode&0x0F00)>>8]
			cpu.PC += 2
			break
		case 0x000E:
			// SHL Vx {, Vy}
			// Set Vx = Vx SHL 1
			if cpu.Registers[(opcode&0x0F00)>>8]&0x8 == 0x8 {
				cpu.Registers[0xF] = 1
			} else {
				cpu.Registers[0xF] = 0
			}
			cpu.Registers[(opcode&0x0F00)>>8] = cpu.Registers[(opcode&0x0F00)>>8] * 2
			break
		}
	case 0x9000:
		// Skip next instruction if Vx != Vy
		if cpu.Registers[(opcode&0x0F00)>>8] != cpu.Registers[(opcode&0x00F0)>>4] {
			cpu.PC += 2
		}
		break
	case 0xA000:
		// Set I = nnn
		cpu.I = opcode & 0x0FFF
		break
	case 0xB000:
		// JP V0, addr
		// Jump to location nnn + V0
		cpu.PC = (opcode & 0x0FFF) + uint16(cpu.Registers[0])
		break
	case 0xC000:
		fmt.Printf("Unimplemented instruction %d", opcode)
		break
	case 0xD000:
		fmt.Printf("Unimplemented instruction %d", opcode)
		break
	case 0xE000:
		switch opcode & 0x000F {
		case 0x000E:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0001:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x000A:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0015:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0018:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x001E:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0029:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0033:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0055:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		case 0x0065:
			fmt.Printf("Unimplemented instruction %d", opcode)
			break
		}
	}
	return nil
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
	for {
		opcode := cpu.Step()
		fmt.Println(strconv.FormatUint(uint64(opcode), 16))
	}
}
