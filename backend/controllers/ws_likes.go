package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var (
	mutex       sync.Mutex
	connections = make(map[int]*websocket.Conn)
)

type Message struct {
	DivId      string `json:"div_id" validate:"required"`
	LikeCount  int    `json:"like_count"`
	HasNotif   bool   `json:"has_notif"`
	RecieverId int
}

func CreateSocket(c *gin.Context) {

	token, valid := c.GetQuery("token")
	if !valid {
		fmt.Println("connection refused")
		return
	}
	userId, err := parseToken(token)
	if err != nil {
		fmt.Println("error on parsing token while connection websocket")
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()
	var message Message
	connections[userId] = conn
	fmt.Println("someone connected with id ", userId)

	for {

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			removeConnection(userId)
			return
		}
		err = json.Unmarshal(p, &message)
		if err != nil {
			fmt.Println(err)
			fmt.Println("invalid message")
			removeConnection(userId)
			return
		}
		validator := validator.New()
		err = validator.Struct(message)
		if err != nil {
			fmt.Println(err)
			fmt.Println("invalid message")
			removeConnection(userId)

			return
		}
		message, err := getUserId(&message)
		if err != nil {
			fmt.Println("error on getting user id")
			fmt.Println(err)
			removeConnection(userId)
			return

		}
		err = addToDatabase(userId, message)
		if err != nil {
			fmt.Println(err)
			fmt.Println("could not add to database the result set")

			continue
		}

		message.LikeCount, err = GetLikes(message.DivId)
		if err != nil {
			fmt.Println(err)
			removeConnection(userId)
			return
		}

		pushMessages(messageType, connections, *message)

	}
}
func getUserId(message *Message) (*Message, error) {

	sql, err := createDb()
	if err != nil {
		fmt.Println("error opening sql")
		return nil, err
	}
	defer sql.Close()
	query := "select user_id from blogs where div_id=?"
	rows, err := sql.Query(query, message.DivId)
	if err != nil {
		fmt.Println("error on getting rows")
		fmt.Println(err)
		return nil, err
	}
	var userId int

	if rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			fmt.Println("error on scanning rows")
			return nil, err
		}
	}
	message.RecieverId = userId

	return message, nil
}

func addToDatabase(userId int, message *Message) error {
	sql, err := createDb()

	if err != nil {
		return err
	}
	defer sql.Close()
	var blogId int
	query := "select id from blogs where div_id=?"
	rows, err := sql.Query(query, message.DivId)
	if err != nil {
		return err
	}
	if rows.Next() {
		err := rows.Scan(&blogId)
		if err != nil {
			return err
		}
	}
	//check like exist
	query = "select count(*) from notification where like_source=? and like_dest=? and blog_id=?"
	rows, err = sql.Query(query, userId, message.RecieverId, blogId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var rowCount int
	if rows.Next() {
		err := rows.Scan(&rowCount)
		if err != nil {
			fmt.Println(err)
			return err
		}

	}

	if rowCount == 1 {
		return errors.New("like exist for certain blog")
	}

	query = "insert into notification values(?,?,?,now())"
	_, err = sql.Exec(query, userId, message.RecieverId, blogId)
	if err != nil {
		return err
	}
	query = "insert into likes values(?,?)"
	_, err = sql.Exec(query, blogId, message.RecieverId)
	if err != nil {
		return err
	}
	query = "update users set has_notif= true where id=?"
	_, err = sql.Query(query, message.RecieverId)
	if err != nil {
		return err
	}

	return nil
}
func pushMessages(messageType int, connections map[int]*websocket.Conn, message Message) {

	for userId, connection := range connections {

		if userId == message.RecieverId {
			message.HasNotif = true
		} else {
			message.HasNotif = false
		}
		p, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = connection.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println(err)
			removeConnection(userId)

		}
	}
}

func removeConnection(userId int) {
	mutex.Lock()
	defer mutex.Unlock()

	for id := range connections {
		if userId == id {

			connections[id].Close()
			delete(connections, id)
			break
		}
	}

}
