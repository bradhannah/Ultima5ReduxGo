package sprites

//func TestBordersEven(t *testing.T) {
//	borderDimensions := getCornersOfReferenceBorder(100, 50)
//
//	assert.Equal(t, 0, borderDimensions.bottomLeft.Min.X)
//	assert.Equal(t, 0, borderDimensions.topLeft.Min.X)
//
//	assert.Equal(t, 49, borderDimensions.bottomLeft.Max.X)
//	assert.Equal(t, 49, borderDimensions.topLeft.Max.X)
//
//	assert.Equal(t, 51, borderDimensions.bottomRight.Min.X)
//	assert.Equal(t, 51, borderDimensions.topRight.Min.X)
//
//	assert.Equal(t, 99, borderDimensions.bottomRight.Max.X)
//	assert.Equal(t, 99, borderDimensions.topRight.Max.X)
//
//	assert.True(t, borderDimensions.topRight.Dx() == borderDimensions.bottomRight.Dx())
//	assert.True(t, borderDimensions.topLeft.Dx() == borderDimensions.bottomLeft.Dx())
//	// assert.True(t, borderDimensions.topRight.Dy() == borderDimensions.bottomRight.Dy())
//	// assert.True(t, borderDimensions.topLeft.Dy() == borderDimensions.bottomLeft.Dy())
//}
//
//func TestBordersOdd(t *testing.T) {
//	borderDimensions := getCornersOfReferenceBorder(101, 51)
//
//	assert.Equal(t, 0, borderDimensions.bottomLeft.Min.X)
//	assert.Equal(t, 0, borderDimensions.topLeft.Min.X)
//
//	assert.Equal(t, 49, borderDimensions.bottomLeft.Max.X)
//	assert.Equal(t, 49, borderDimensions.topLeft.Max.X)
//
//	assert.Equal(t, 51, borderDimensions.bottomRight.Min.X)
//	assert.Equal(t, 51, borderDimensions.topRight.Min.X)
//
//	assert.Equal(t, 100, borderDimensions.bottomRight.Max.X)
//	assert.Equal(t, 100, borderDimensions.topRight.Max.X)
//
//	assert.True(t, borderDimensions.topRight.Dx() == borderDimensions.bottomRight.Dx())
//	assert.True(t, borderDimensions.topLeft.Dx() == borderDimensions.bottomLeft.Dx())
//	assert.True(t, borderDimensions.topRight.Dy() == borderDimensions.bottomRight.Dy())
//	assert.True(t, borderDimensions.topLeft.Dy() == borderDimensions.bottomLeft.Dy())
//}
