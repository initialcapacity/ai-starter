package slicesupport_test

import (
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	result := slicesupport.Map([]int{1, 2, 3}, func(i int) int { return i * 2 })
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestFind(t *testing.T) {
	result, found2 := slicesupport.Find([]int{1, 2, 3}, func(i int) bool { return i == 2 })
	assert.Equal(t, 2, result)
	assert.True(t, found2)

	_, found4 := slicesupport.Find([]int{1, 2, 3}, func(i int) bool { return i == 4 })
	assert.False(t, found4)
}
