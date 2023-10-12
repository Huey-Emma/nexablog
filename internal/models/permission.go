package models

type Permission struct {
	PermissionID int
	Code         string
}

type Permissions []Permission

func (p Permissions) Include(code string) bool {
	for idx := range p {
		if p[idx].Code == code {
			return true
		}
	}
	return false
}
