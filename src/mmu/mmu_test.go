package mmu

import (
	"cartridge"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"types"
)

func TestImplementsInterface(t *testing.T) {
	gbc := NewGbcMMU()
	assert.Implements(t, (*MemoryMappedUnit)(nil), gbc)
}

func TestLoadBiosROM(t *testing.T) {
	ROM := GenerateDummyBIOSROM()
	mmu := NewGbcMMU()
	mmu.LoadBIOS(ROM)
	assert.Equal(t, mmu.bios, ROM)
}

func TestLoadBiosROMForSizeError(t *testing.T) {
	mmu := NewGbcMMU()
	ok, err := mmu.LoadBIOS(make([]byte, 300))

	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ROMIsBiggerThanRegion, err)
}

func TestReadByteFromBIOS(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(true)

	BIOS := GenerateDummyBIOSROM()
	_, err := mmu.LoadBIOS(BIOS)
	if err != nil {
		t.FailNow()
	}

	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	cartRom := make([]byte, 32768)
	var counter int = 255
	for i := 0; i < 32768; i++ {
		cartRom[i] = byte(counter)
		counter--
	}
	cart.Init(cartRom)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0x0000; i < 0x0100; i++ {
		memValue := mmu.ReadByte(i)
		assert.Equal(t, mmu.bios[i], memValue)
		//shouldn't get bytes from cartridge address space
		assert.NotEqual(t, cartRom[i], memValue)
	}
}

func TestReadByteFromCartridgeFromBIOSAddressSpace(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	BIOS := GenerateDummyBIOSROM()
	_, err := mmu.LoadBIOS(BIOS)
	if err != nil {
		t.FailNow()
	}

	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	cartRom := make([]byte, 32768)
	var counter int = 255
	for i := 0; i < 32768; i++ {
		cartRom[i] = byte(counter)
		counter--
	}
	cartRom[0x0147] = cartridge.MBC0

	//reinitialise cartridge with new rom
	cart.Init(cartRom)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0x0000; i < 0x0100; i++ {
		memValue := mmu.ReadByte(i)
		//shouldn't get bytes from BIOS address space
		assert.NotEqual(t, mmu.bios[i], memValue)
		assert.Equal(t, cartRom[i], memValue)
	}
}

func TestReadByteFromROMBank0ForMBC0(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0x0000; i <= 0x3FFF; i++ {
		memValue := mmu.ReadByte(i)
		assert.Equal(t, cart.ReadByteFromROM(i), memValue)
	}
}

func TestReadByteFromSwitchableROMBankForMBC0(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0x4000; i <= 0x7FFF; i++ {
		memValue := mmu.ReadByte(i)
		assert.Equal(t, cart.ReadByteFromSwitchableROM(0, i), memValue)
	}
}

func TestReadByteFromGraphicsRAM(t *testing.T) {
	mmu := NewGbcMMU()
	gpu := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(gpu, 0x8000, 0x9FFF)

	var addr types.Word
	for addr = 0x8000; addr <= 0x9FFF; addr++ {
		expected := byte(addr & 0xFF)
		gpu.internalMemory[addr] = expected
		memValue := mmu.ReadByte(addr)
		assert.Equal(t, expected, memValue)
	}
}

func TestReadByteFromSwitchableRAMBankForMBC0(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0xA000; i <= 0xBFFF; i++ {
		memValue := mmu.ReadByte(i)
		assert.Equal(t, cart.ReadByteFromSwitchableRAM(0, i), memValue)
	}
}

func TestReadByteFromInternalRAM(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	var i types.Word
	for i = 0xC000; i <= 0xDFFF; i++ {
		addr := i & (0xDFFF - 0xC000)
		expected := byte(i & 0x00FF)
		mmu.internalRAM[addr] = expected

		actual := mmu.ReadByte(i)
		assert.Equal(t, expected, actual)
	}
}

func TestReadByteFromInternalRAMShadow(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	var i types.Word
	for i = 0xE000; i <= 0xFDFF; i++ {
		addr := i & (0xFDFF - 0xE000)
		expected := byte(i & 0x00FF)
		mmu.internalRAMShadow[addr] = expected

		actual := mmu.ReadByte(i)
		assert.Equal(t, expected, actual)
	}
}

func TestReadByteFromGraphicsSpriteRAM(t *testing.T) {
	mmu := NewGbcMMU()
	gpu := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(gpu, 0xFE00, 0xFE9F)

	var addr types.Word
	for addr = 0xFE00; addr <= 0xFE9F; addr++ {
		expected := byte(addr & 0xFF)
		gpu.internalMemory[addr] = expected
		memValue := mmu.ReadByte(addr)
		assert.Equal(t, expected, memValue)
	}
}

func TestReadByteFromDMGRegister(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	var addr types.Word = 0xFF50
	expected := byte(addr & 0xFF)
	mmu.dmgStatusRegister = expected

	memValue := mmu.ReadByte(addr)
	assert.Equal(t, expected, memValue)
}

func TestReadByteFromGraphicsRegisters(t *testing.T) {
	mmu := NewGbcMMU()
	peripheral := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(peripheral, 0xFF40, 0xFF70)

	var addr types.Word
	for addr = 0xFF40; addr <= 0xFF70; addr++ {
		expected := byte(addr & 0xFF)
		peripheral.internalMemory[addr] = expected
		memValue := mmu.ReadByte(addr)
		assert.Equal(t, expected, memValue)
	}
}

func TestReadByteFromZeroPageRAM(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	var i types.Word
	for i = 0xFF80; i < 0xFFFF; i++ {
		addr := i & (0xFFFF - 0xFF80)
		expected := byte(i & 0x00FF)
		mmu.zeroPageRAM[addr] = expected

		actual := mmu.ReadByte(i)
		assert.Equal(t, expected, actual)
	}
}

func TestWriteByteToGraphicsMemory(t *testing.T) {
	mmu := NewGbcMMU()
	peripheral := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(peripheral, 0x8000, 0x9FFF)

	var addr types.Word
	for addr = 0x8000; addr <= 0x9FFF; addr++ {
		expected := byte(addr & 0xFF)
		mmu.WriteByte(addr, expected)
		memValue := peripheral.internalMemory[addr]
		assert.Equal(t, expected, memValue)
	}
}

func TestWriteByteToSwitchableRAMBankForMBC0(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	cart := GenerateDummyCartridge(cartridge.CartridgeTypes[cartridge.MBC0], 32768)
	mmu.LoadCartridge(cart)

	var i types.Word
	for i = 0xA000; i <= 0xBFFF; i++ {
		expected := byte(i & 0xFF)
		mmu.WriteByte(i, expected)
		assert.Equal(t, cart.ReadByteFromSwitchableRAM(0, i), expected)
	}
}

func TestWriteByteToInternalRAM(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	var i types.Word
	for i = 0xC000; i <= 0xDFFF; i++ {
		addr := i & (0xDFFF - 0xC000)
		expected := byte(i & 0x00FF)
		mmu.WriteByte(i, expected)
		actual := mmu.internalRAM[addr]
		assert.Equal(t, expected, actual)
	}
}

func TestWriteByteToInternalRAMAndRAMIsShadowed(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	var expected byte = 0xAA

	var i types.Word
	for i = 0xC000; i <= 0xDFFF; i++ {
		mmu.WriteByte(i, expected)
	}

	for i := 0xE000; i <= 0xFDFF; i++ {
		ramAddr := i & (0xDFFF - 0xC000)
		ramShadowAddr := i & (0xFDFF - 0xE000)
		assert.Equal(t, mmu.internalRAM[ramAddr], mmu.internalRAMShadow[ramShadowAddr])
	}
}

func TestWriteByteToGraphicsSpriteRAM(t *testing.T) {
	mmu := NewGbcMMU()
	gpu := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(gpu, 0xFE00, 0xFE9F)

	var addr types.Word
	for addr = 0xFE00; addr <= 0xFE9F; addr++ {
		expected := byte(addr & 0xFF)
		mmu.WriteByte(addr, expected)
		memValue := gpu.internalMemory[addr]
		assert.Equal(t, expected, memValue)
	}
}

func TestWriteByteToDMGRegister(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)
	var addr types.Word = 0xFF50
	var expected byte = 0xDF

	mmu.WriteByte(addr, expected)
	memValue := mmu.dmgStatusRegister
	assert.Equal(t, expected, memValue)
}

func TestWriteByteToGraphicsRegisters(t *testing.T) {
	mmu := NewGbcMMU()
	peripheral := new(MockPeripheral)
	mmu.SetInBootMode(false)
	mmu.ConnectPeripheral(peripheral, 0xFF40, 0xFF70)

	var addr types.Word
	for addr = 0xFF40; addr <= 0xFF70; addr++ {
		expected := byte(addr & 0xFF)
		mmu.WriteByte(addr, expected)
		memValue := peripheral.internalMemory[addr]
		assert.Equal(t, expected, memValue)
	}
}

func TestWriteByteToZeroPageRAM(t *testing.T) {
	mmu := NewGbcMMU()
	mmu.SetInBootMode(false)

	var i types.Word
	for i = 0xFF80; i < 0xFFFF; i++ {
		addr := i & (0xFFFF - 0xFF80)
		expected := byte(i & 0x00FF)
		mmu.WriteByte(i, expected)
		actual := mmu.zeroPageRAM[addr]
		assert.Equal(t, expected, actual)
	}
}

func GenerateDummyBIOSROM() []byte {
	var BIOSROM []byte
	for i := 0; i < 256; i++ {
		BIOSROM = append(BIOSROM, byte(i))
	}
	return BIOSROM
}

func GenerateDummyCartridge(cartType cartridge.CartridgeType, size int) *cartridge.Cartridge {
	cartRom := GenerateCartridgeROM(size)
	cartRom[0x0147] = 0x00
	title := "testing         "
	for i := 0x0134; i < 0x0143; i++ {
		cartRom[i] = title[i-0x0134]
	}

	cart, _ := cartridge.NewCartridge(cartRom)
	cart.Type = cartType
	return cart
}

func GenerateCartridgeROM(no int) []byte {
	var bytes []byte
	for i := 0; i < no; i++ {
		bytes = append(bytes, byte(i))
	}
	return bytes
}

// MOCK
type MockPeripheral struct {
	internalMemory [65535]byte
}

func (m *MockPeripheral) Name() string {
	return "MockPeripheral"
}

func (m *MockPeripheral) Read(addr types.Word) byte {
	return m.internalMemory[addr]
}

func (m *MockPeripheral) Write(addr types.Word, value byte) {
	m.internalMemory[addr] = value
}
