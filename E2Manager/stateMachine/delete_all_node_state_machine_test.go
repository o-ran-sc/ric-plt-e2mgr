package stateMachine

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeNextStateDeleteAll(t *testing.T) {

	_, result := NodeNextStateDeleteAll(entities.ConnectionStatus_UNKNOWN_CONNECTION_STATUS)

	assert.False(t, result)
}
