package helpers

func GetAsBitmapBoolList(data []byte, start int, length int) []bool {
	const bitsPerByte = 8
	boolList := make([]bool, 0, length*bitsPerByte)

	for nByte := start; nByte < start+length; nByte++ {
		compareByte := byte(0x80)
		curByte := data[nByte]

		for nBit := bitsPerByte - 1; nBit >= 0; nBit-- {
			curBit := (compareByte & curByte) > 0
			compareByte >>= 1
			boolList = append(boolList, curBit)
		}
	}

	return boolList
}
