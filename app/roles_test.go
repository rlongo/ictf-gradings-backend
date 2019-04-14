package app

import (
	"testing"
)

func TestSingleRole(t *testing.T) {
	expectedRole := RoleAdmin

	inputs := []struct {
		role  Role
		match bool
	}{
		{(RoleNone | RoleSupervisor), false},
		{(RoleAdmin | RoleInstructor), true},
		{RoleSupervisor, false},
		{(Role)(0), false},
	}

	for i, input := range inputs {
		if input.match != expectedRole.Matches(input.role) {
			t.Errorf("Input %d failed. Expecting: %t", i, input.match)
		}
	}
}

func TestMultiRole(t *testing.T) {
	expectedRole := (RoleAdmin | RoleSupervisor)

	inputs := []struct {
		role  Role
		match bool
	}{
		{(RoleNone | RoleSupervisor), true},
		{(RoleAdmin | RoleInstructor), true},
		{RoleSupervisor, true},
		{(Role)(0), false},
		{(RoleAdmin | RoleSupervisor), true},
		{(RoleAdmin | RoleSupervisor | RoleInstructor), true},
	}

	for i, input := range inputs {
		if input.match != expectedRole.Matches(input.role) {
			t.Errorf("Input %d failed. Expecting: %t", i, input.match)
		}
	}
}
