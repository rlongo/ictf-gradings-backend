package app

// Role is what type of signed-in user are we
type Role int

// Types of roles within our application
// These are flags. You can perform bitwise ops on them
const (
	RoleNone       Role = 0
	RoleInstructor Role = 1 << 0
	RoleSupervisor Role = 1 << 1
	RoleAdmin      Role = 1 << 2
)

// Matches determines whether the input role is a superset of
// the role being matched against
func (r *Role) Matches(r2 Role) bool {
	return *r == RoleNone || ((*r)&r2) > 0
}
