package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToMap(t *testing.T) {
	p := PermissionList{
		{
			Subject: User,
			Action:  Write,
		},
	}

	m := p.ToMap()
	assert.Equal(t, true, m[User][Write])
}
