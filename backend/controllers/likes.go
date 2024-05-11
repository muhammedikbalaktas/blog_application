package controllers

import "fmt"

func GetLikes(divId string) (int, error) {

	db, err := createDb()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	query := "select count(*) from blogs inner join likes on id=blog_id where div_id=?"
	rows, err := db.Query(query, divId)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	var likeCount int
	if rows.Next() {
		err := rows.Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			return -1, err
		}
	}
	db.Close()
	return likeCount, nil
}
