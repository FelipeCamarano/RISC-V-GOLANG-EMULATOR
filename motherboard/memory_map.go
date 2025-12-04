package motherboard

const (
	//para mapear a memoria RAM a partir do endereço 0x80000000 para os testes
	DefaultRAMBase = 0x00000000
	DefaultRAMSize = 16 * 1024 * 1024
	MAXADDRESS     = 0xFFFFFFFF
)
