package dto

type UpdateUserRoleReq struct {
	Role string `json:"role" validate:"required,oneof=USER ADMIN"`
}

type AdminDeleteUserReq struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

type AdminDeletePostReq struct {
	Reason string `json:"reason" validate:"required,max=500"`
}
