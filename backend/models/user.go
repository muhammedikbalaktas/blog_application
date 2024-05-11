package models

type User struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GetUserParams struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// +----------+-------------+------+-----+---------+----------------+
// | Field    | Type        | Null | Key | Default | Extra          |
// +----------+-------------+------+-----+---------+----------------+
// | id       | int         | NO   | PRI | NULL    | auto_increment |
// | name     | varchar(20) | YES  |     | NULL    |                |
// | surname  | varchar(20) | YES  |     | NULL    |                |
// | username | varchar(20) | YES  | UNI | NULL    |                |
// | email    | varchar(50) | YES  | UNI | NULL    |                |
// | password | varchar(60) | YES  |     | NULL    |                |
// +----------+-------------+------+-----+---------+----------------+
