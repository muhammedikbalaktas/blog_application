package models

import "mime/multipart"

type Blog struct {
	Id        int                   `form:"id,omitempty"`
	Username  string                `form:"username,omitempty"`
	Title     string                `form:"title" validate:"required"`
	Content   string                `form:"content" validate:"required"`
	UserToken string                `form:"user_token,omitempty"`
	ImageName string                `form:"image_name,omitempty"`
	DivId     string                `form:"div_id,omitempty"`
	LikeCount int                   `form:"like_count"`
	Image     *multipart.FileHeader `form:"image,omitempty" binding:"-"`
}
