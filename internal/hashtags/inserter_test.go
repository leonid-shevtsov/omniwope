package hashtags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	assert.Equal(t, "A #Hugo", Insert([]string{"Hugo"}, "A Hugo"))
	assert.Equal(t, "A <a>Hugo</a> #Hugo", Insert([]string{"Hugo"}, "A <a>Hugo</a> Hugo"))
	assert.Equal(t, "#Boss\n\nA <a>Hugo</a> Hugo", Insert([]string{"Boss"}, "A <a>Hugo</a> Hugo"))
	assert.Equal(t, "#Boss\n\nA <b>#Hugo</b> Hugo", Insert([]string{"Hugo", "Boss"}, "A <b>Hugo</b> Hugo"))
}
