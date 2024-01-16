package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var db *sql.DB
var config Config

func main() {
	alo := ""
	jsonContent, err := json.Marshal(alo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(jsonContent)

	// var err error

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Đã xảy ra lỗi khi đọc tệp YAML: %v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Lỗi phân tích tệp YAML: %v", err)
	}

	dataSourceName := config.Mysql.Username + ":" + config.Mysql.Password + "@tcp(localhost:3306)/" + config.Mysql.Host

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.Default())

	//Staffs
	r.GET("/staffs/:id", getStaff)
	r.POST("/staffs", createStaff)
	r.PUT("/staffs/:id", updateStaff)
	r.DELETE("/staffs/:id", deleteStaff)
	r.GET("/staffs", getStaffsPaged)

	//Room Types
	r.GET("/room_types/:id", getRoomType)
	r.POST("/room_types", createRoomType)
	r.PUT("/room_types/:id", updateRoomType)
	r.DELETE("/room_types/:id", deleteRoomType)
	r.GET("/room_types", getRoomTypesPaged)
	r.GET("/room_types/all", getAllRoomTypes)

	//Rooms
	r.GET("/rooms/:id", getRoom)
	r.POST("/rooms", createRoom)
	r.PUT("/rooms/:id", updateRoom)
	r.DELETE("/rooms/:id", deleteRoom)
	r.GET("/rooms", getRoomsPaged)

	//Orders
	r.GET("/orders/:id", getOrder)
	r.POST("/orders", createOrder)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)
	r.GET("/orders", getOrdersPaged)

	//Posts
	r.GET("/posts/:id", getPost)
	r.POST("/posts", createPost)
	r.PUT("/posts/:id", updatePost)
	r.DELETE("/posts/:id", deletePost)
	r.GET("/posts", getPostsPaged)
	r.GET("/posts/user_id=:user_id", getPostsByUserIDPaged)

	//Blogs
	r.GET("/blogs/:id", getBlog)
	r.POST("/blogs", createBlog)
	r.PUT("/blogs/:id", updateBlog)
	r.DELETE("/blogs/:id", deleteBlog)
	r.GET("/blogs", getBlogsPaged)
	r.GET("/blogs/user_id=:user_id", getBlogsByUserIDPaged)

	//Users
	r.GET("/users/:id", getUserByID)
	r.GET("/users/email=:email", getUserByEmail)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	r.GET("/users", getUsersPaged)

	//Sales
	r.GET("/sales", getSale)

	//Images
	// r.POST("/images", uploadImage)
	// r.GET("/images", getImages)

	// Accounts
	r.GET("/accounts/:id", getAccountByID)
	r.GET("/accounts/email=:email", getAccountByEmail)
	r.POST("/accounts", createAccount)
	r.PUT("/accounts/:id", updateAccount)
	r.PUT("/accounts-without-avatar/:id", updateAccountWithoutAvatar)
	r.DELETE("/accounts/:id", deleteAccount)
	r.GET("/accounts", getAccountsPaged)

	port := ":" + config.Port
	r.Run(port)
}

type smtpServer struct {
	host string
	port string
}

// serverName URI to smtp server
func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

type Config struct {
	Port string `yaml:"PORT"`

	Mysql struct {
		Host     string `yaml:"HOST"`
		Username string `yaml:"USERNAME"`
		Password string `yaml:"PASSWORD"`
	} `yaml:"MYSQL"`

	Email struct {
		Address  string `yaml:"ADDRESS"`
		Password string `yaml:"PASSWORD"`
	} `yaml:"EMAIL"`
}

type Staff struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type RoomType struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Price          int    `json:"price"`
	BackgroundLink string `json:"background_link"`
}

type Room struct {
	ID         int `json:"id"`
	RoomTypeID int `json:"room_type_id"`
}

type Order struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`
	RoomID      int    `json:"room_id"`
	Status      string `json:"status"`
	RoomTypeID  int    `json:"room_type_id"`
	Email       string `json:"email"`
}

type Blog struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	Content   string `json:"content"`
	Avatar    string `json:"avatar"`
	UserID    int    `json:"user_id"`
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	Content   []Content `json:"content"`
	Avatar    string    `json:"avatar"`
	UserID    int       `json:"user_id"`
}

type Content struct {
	Paragraph string `json:"paragraph"`
	Image     string `json:"image"`
}

type Sale struct {
	ID         int    `json:"id"`
	RoomTypeID int    `json:"room_type_id"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
}

type SaleResult struct {
	Time      string `json:"time"`
	SaleValue int    `json:"sale_value"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Type        string `json:"type"`
	PhoneNumber string `json:"phone_number"`
}

type Account struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	AvatarName  string `json:"avatar_name"`
	AvatarData  []byte `json:"avatar_data"`
	Type        string `json:"type"`
	PhoneNumber string `json:"phone_number"`
}

type Image struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type ResultRoomTypePagenation struct {
	PageCount int        `json:"page_count"`
	Items     []RoomType `json:"items"`
}

type ResultStaffPagenation struct {
	PageCount int     `json:"page_count"`
	Items     []Staff `json:"items"`
}

type ResultRoomPagenation struct {
	PageCount int    `json:"page_count"`
	Items     []Room `json:"items"`
}

type ResultOrderPagenation struct {
	PageCount int     `json:"page_count"`
	Items     []Order `json:"items"`
}

type ResultPostPagenation struct {
	PageCount int    `json:"page_count"`
	Items     []Post `json:"items"`
}

type ResultBlogPagenation struct {
	PageCount int    `json:"page_count"`
	Items     []Blog `json:"items"`
}

type ResultUserPagenation struct {
	PageCount int    `json:"page_count"`
	Items     []User `json:"items"`
}

type ResultAccountPagenation struct {
	PageCount int       `json:"page_count"`
	Items     []Account `json:"items"`
}

//Staff

func getStaff(c *gin.Context) {
	id := c.Param("id")
	var staff Staff
	err := db.QueryRow("SELECT id, name, phone_number FROM staffs WHERE id = ?", id).Scan(&staff.ID, &staff.Name, &staff.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(http.StatusOK, staff)
}

func createStaff(c *gin.Context) {
	var staff Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO staffs (name, phone_number) VALUES (?, ?)", staff.Name, staff.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	staff.ID = int(id)
	c.JSON(http.StatusCreated, staff)
}

func updateStaff(c *gin.Context) {
	id := c.Param("id")
	var staff Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE staffs SET name = ?, phone_number = ? WHERE id = ?", staff.Name, staff.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, staff)
}

func deleteStaff(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM staffs WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff deleted"})
}

func getStaffsPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "%")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, name, phone_number FROM staffs WHERE id LIKE ? OR name LIKE ? OR phone_number LIKE ? LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var staffs []Staff
	for rows.Next() {
		var staff Staff
		err := rows.Scan(&staff.ID, &staff.Name, &staff.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		staffs = append(staffs, staff)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM staffs WHERE id LIKE ? OR name LIKE ? OR phone_number LIKE ?", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultStaffPagenation{pageCount, staffs}
	c.JSON(http.StatusOK, result)
}

//Room Types

func getRoomType(c *gin.Context) {
	id := c.Param("id")
	var roomType RoomType
	err := db.QueryRow("SELECT id, name, description, price, background_link FROM room_types WHERE id = ?", id).Scan(&roomType.ID, &roomType.Name, &roomType.Description, &roomType.Price, &roomType.BackgroundLink)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(http.StatusOK, roomType)
}

func createRoomType(c *gin.Context) {
	var roomType RoomType
	if err := c.ShouldBindJSON(&roomType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO room_types (name, description, price, background_link) VALUES (?, ?, ?, ?)", roomType.Name, roomType.Description, roomType.Price, roomType.BackgroundLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	roomType.ID = int(id)
	c.JSON(http.StatusCreated, roomType)
}

func updateRoomType(c *gin.Context) {
	id := c.Param("id")
	var roomType RoomType
	if err := c.ShouldBindJSON(&roomType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE room_types SET name = ?, description = ?, price = ?, background_link = ? WHERE id = ?", roomType.Name, roomType.Description, roomType.Price, roomType.BackgroundLink, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roomType)
}

func deleteRoomType(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM room_types WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room type deleted"})
}

func getRoomTypesPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	search := c.DefaultQuery("search", "%")

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, name, description, price, background_link FROM room_types WHERE id LIKE ? OR name LIKE ? LIMIT ? OFFSET ?", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var roomTypes []RoomType
	for rows.Next() {
		var roomType RoomType
		err := rows.Scan(&roomType.ID, &roomType.Name, &roomType.Description, &roomType.Price, &roomType.BackgroundLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		roomTypes = append(roomTypes, roomType)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM room_types WHERE id LIKE ? OR name LIKE ?", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultRoomTypePagenation{pageCount, roomTypes}
	c.JSON(http.StatusOK, result)
}

func getAllRoomTypes(c *gin.Context) {

	rows, err := db.Query("SELECT id, name, description, price, background_link FROM room_types")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var roomTypes []RoomType
	for rows.Next() {
		var roomType RoomType
		err := rows.Scan(&roomType.ID, &roomType.Name, &roomType.Description, &roomType.Price, &roomType.BackgroundLink)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		roomTypes = append(roomTypes, roomType)
	}
	c.JSON(http.StatusOK, roomTypes)
}

// Room

func getRoom(c *gin.Context) {
	id := c.Param("id")
	var room Room
	err := db.QueryRow("SELECT id, room_type_id FROM rooms WHERE id = ?", id).Scan(&room.ID, &room.RoomTypeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func createRoom(c *gin.Context) {
	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO rooms (room_type_id) VALUES (?)", room.RoomTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	room.ID = int(id)
	c.JSON(http.StatusCreated, room)
}

func updateRoom(c *gin.Context) {
	id := c.Param("id")
	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE rooms SET room_type_id = ? WHERE id = ?", room.RoomTypeID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func deleteRoom(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM rooms WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted"})
}

func getRoomsPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, room_type_id FROM rooms LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		err := rows.Scan(&room.ID, &room.RoomTypeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rooms = append(rooms, room)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultRoomPagenation{pageCount, rooms}
	c.JSON(http.StatusOK, result)
}

//Orders

func getOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	err := db.QueryRow("SELECT id, name, phone_number, email,time_start, time_end, room_id, status FROM orders WHERE id = ?", id).Scan(&order.ID, &order.Name, &order.PhoneNumber, &order.Email, &order.TimeStart, &order.TimeEnd, &order.RoomID, &order.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT id, name, phone_number, email,time_start, time_end, room_id, status, room_type_id  FROM orders WHERE room_type_id = ? AND  status = ? AND ((time_start <= ? AND time_end >= ?) || (time_start >= ? AND time_start <= ?) )", order.RoomTypeID, "confirmed", order.TimeStart, order.TimeStart, order.TimeStart, order.TimeEnd)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		var room Room
		err := db.QueryRow("SELECT id, room_type_id FROM rooms WHERE room_type_id = ? ORDER BY RAND()", order.RoomTypeID).Scan(&room.ID, &room.RoomTypeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		_, err = db.Exec("INSERT INTO orders (name, phone_number, email,time_start, time_end, room_id, status, room_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", order.Name, order.PhoneNumber, order.Email, order.TimeStart, order.TimeEnd, room.ID, "unconfirmed", order.RoomTypeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = sentEmail(order.Name, order.Email, "Liên hệ số sau (+84) 36 666 6250 hoặc đợi chúng tôi liên hệ cho bạn để xác nhận order", "Unconfirmed")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, "Create order succesfully")
		return

	}

	var orderWithSameRoomTypeIDs []Order

	var orderWithSameRoomTypeIDTemp Order
	err = rows.Scan(&orderWithSameRoomTypeIDTemp.ID, &orderWithSameRoomTypeIDTemp.Name, &orderWithSameRoomTypeIDTemp.PhoneNumber, &orderWithSameRoomTypeIDTemp.Email, &orderWithSameRoomTypeIDTemp.TimeStart, &orderWithSameRoomTypeIDTemp.TimeEnd, &orderWithSameRoomTypeIDTemp.RoomID, &orderWithSameRoomTypeIDTemp.Status, &orderWithSameRoomTypeIDTemp.RoomTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderWithSameRoomTypeIDs = append(orderWithSameRoomTypeIDs, orderWithSameRoomTypeIDTemp)

	for rows.Next() {
		var orderWithSameRoomTypeID Order
		err := rows.Scan(&orderWithSameRoomTypeID.ID, &orderWithSameRoomTypeID.Name, &orderWithSameRoomTypeID.PhoneNumber, &orderWithSameRoomTypeID.Email, &orderWithSameRoomTypeID.TimeStart, &orderWithSameRoomTypeID.TimeEnd, &orderWithSameRoomTypeID.RoomID, &orderWithSameRoomTypeID.Status, &orderWithSameRoomTypeID.RoomTypeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderWithSameRoomTypeIDs = append(orderWithSameRoomTypeIDs, orderWithSameRoomTypeID)
	}

	query := fmt.Sprintf("SELECT id, room_type_id FROM rooms WHERE room_type_id = %d AND 1 = 1", order.RoomTypeID)
	for _, value := range orderWithSameRoomTypeIDs {
		query += " AND id != " + fmt.Sprintf("%d", value.RoomID) + " AND 1 = 1"
	}

	var roomTemp Room
	err = db.QueryRow(query+" ORDER BY RAND()").Scan(&roomTemp.ID, &roomTemp.RoomTypeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("INSERT INTO orders (name, phone_number, email,time_start, time_end, room_id, status, room_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", order.Name, order.PhoneNumber, order.Email, order.TimeStart, order.TimeEnd, roomTemp.ID, "unconfirmed", order.RoomTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sentEmail(order.Name, order.Email, strconv.Itoa(order.ID), order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Create order successfully")
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE orders SET name = ?, phone_number = ?, email = ?,time_start = ?, time_end = ?, room_id = ?, status = ? WHERE id = ?", order.Name, order.PhoneNumber, order.Email, order.TimeStart, order.TimeEnd, order.RoomID, order.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sentEmail(order.Name, order.Email, id, order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	var order Order
	err := db.QueryRow("SELECT id, name, phone_number, email,time_start, time_end, room_id, status FROM orders WHERE id = ?", id).Scan(&order.ID, &order.Name, &order.PhoneNumber, &order.Email, &order.TimeStart, &order.TimeEnd, &order.RoomID, &order.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	_, err = db.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sentEmail(order.Name, order.Email, id, "Canceled")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

func getOrdersPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "%")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, name, phone_number, email, time_start, time_end, room_id, status, room_type_id FROM orders WHERE id LIKE ? OR name LIKE ? OR phone_number LIKE ? OR email LIKE ? OR time_start LIKE ?  OR time_end LIKE ? OR status LIKE ? LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", search+"%", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Name, &order.PhoneNumber, &order.Email, &order.TimeStart, &order.TimeEnd, &order.RoomID, &order.Status, &order.RoomTypeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, order)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE id LIKE ? OR name LIKE ? OR phone_number LIKE ? OR email LIKE ? OR time_start LIKE ?  OR time_end LIKE ?  OR status LIKE ?", search+"%", search+"%", search+"%", search+"%", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultOrderPagenation{pageCount, orders}
	c.JSON(http.StatusOK, result)
}

func sentEmail(name string, email string, orderID string, orderStatus string) error {
	from := config.Email.Address
	password := config.Email.Password
	to := []string{
		email,
	}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	template := "Xin chào %s.\n\nChúng tôi gửi bạn thông tin của order như sau:\n\nMã order: %s\n\nTrạng thái order: %s"
	msg := fmt.Sprintf(template, name, orderID, strings.Title(orderStatus))

	message := []byte(msg)
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

//Posts

func getPost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	var contentData []byte

	err := db.QueryRow("SELECT id, title, created_at, content, user_id FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.CreatedAt, &contentData, &post.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var contents []Content
	if err = json.Unmarshal([]byte(string(contentData)), &contents); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	post.Content = contents

	c.JSON(http.StatusOK, post)
}

func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonContent, err := json.Marshal(post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created_at := time.Now()
	post.CreatedAt = created_at.Format("2006/01/02 15:04:05")

	result, err := db.Exec("INSERT INTO posts (title, created_at, content, avatar, user_id) VALUES (?, ?, ?,?, ?)", post.Title, post.CreatedAt, jsonContent, post.Avatar, post.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	post.ID = int(id)
	c.JSON(http.StatusCreated, post)
}

func updatePost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonContent, err := json.Marshal(post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE posts SET title = ?, created_at = ?, content = ?, avatar = ? WHERE id = ?", post.Title, post.CreatedAt, jsonContent, post.Avatar, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func deletePost(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

func getPostsPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	search := c.DefaultQuery("search", "%")

	rows, err := db.Query("SELECT id, title, created_at, avatar, content, user_id FROM posts WHERE id LIKE ? OR title LIKE ? OR user_id LIKE ? ORDER BY created_at DESC  LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var contentData []byte
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.CreatedAt, &post.Avatar, &contentData, &post.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var contents []Content
		if err = json.Unmarshal([]byte(string(contentData)), &contents); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		post.Content = contents

		posts = append(posts, post)

	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE id LIKE ? OR title LIKE ? OR user_id LIKE ?", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultPostPagenation{pageCount, posts}
	c.JSON(http.StatusOK, result)
}

func getPostsByUserIDPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "%")

	user_id := c.Param("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, title, created_at, avatar, content, user_id FROM posts WHERE user_id = ? AND (id LIKE ? OR title LIKE ?) ORDER BY created_at DESC  LIMIT ? OFFSET ? ", user_id, search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var contentData []byte
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.CreatedAt, &post.Avatar, &contentData, &post.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var contents []Content
		if err = json.Unmarshal([]byte(string(contentData)), &contents); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		post.Content = contents

		posts = append(posts, post)

	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ? AND (id LIKE ? OR title LIKE ?)", user_id, search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultPostPagenation{pageCount, posts}
	c.JSON(http.StatusOK, result)
}

//Users

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user User
	err := db.QueryRow("SELECT id, email, password, name, avatar, type, phone_number FROM users WHERE id = ?", id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Avatar, &user.Type, &user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func getUserByEmail(c *gin.Context) {
	email := c.Param("email")

	fmt.Println(email)
	var user User
	err := db.QueryRow("SELECT id, email, password, name, avatar, type, phone_number FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Avatar, &user.Type, &user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO users (email, password, name, avatar, type, phone_number) VALUES (?, ?, ?, ?, ?, ?)", user.Email, user.Password, user.Name, user.Avatar, user.Type, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	sentEmailAccount(user.Name, user.Email, user.Password)

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE users SET email = ?, password = ?, name = ?, avatar = ?, type = ?, phone_number = ? WHERE id = ?", user.Email, user.Password, user.Name, user.Avatar, user.Type, user.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sentEmailAccount(user.Name, user.Email, user.Password)

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func getUsersPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "%")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, email, password, name, avatar, type, phone_number FROM users WHERE email LIKE ? OR phone_number LIKE ? OR name LIKE ?  OR id LIKE ?  OR type LIKE ? LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Avatar, &user.Type, &user.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email LIKE ? OR phone_number LIKE ? OR name LIKE ?  OR id LIKE ?  OR type LIKE ?", search+"%", search+"%", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultUserPagenation{pageCount, users}
	c.JSON(http.StatusOK, result)
}

func sentEmailAccount(name string, email string, accountPassword string) error {
	from := config.Email.Address
	password := config.Email.Password

	to := []string{
		email,
	}
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	template := "Xin chào %s.\n\nChúng tôi gửi bạn thông tin của tài khoản như sau:\n\nEmail đăng nhập: %s\n\nMật khẩu: %s"
	msg := fmt.Sprintf(template, name, email, accountPassword)

	message := []byte(msg)
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

//Blogs

func getBlog(c *gin.Context) {
	id := c.Param("id")
	var blog Blog

	err := db.QueryRow("SELECT id, title, created_at, content, user_id FROM blogs WHERE id = ?", id).Scan(&blog.ID, &blog.Title, &blog.CreatedAt, &blog.Content, &blog.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func createBlog(c *gin.Context) {
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created_at := time.Now()
	blog.CreatedAt = created_at.Format("2006/01/02 15:04:05")

	result, err := db.Exec("INSERT INTO blogs (title, created_at, content, avatar, user_id) VALUES (?, ?, ?, ?, ?)", blog.Title, blog.CreatedAt, blog.Content, blog.Avatar, blog.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	blog.ID = int(id)
	c.JSON(http.StatusCreated, blog)
}

func updateBlog(c *gin.Context) {
	id := c.Param("id")
	var blog Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE blogs SET title = ?, created_at = ?, content = ?, avatar = ? WHERE id = ?", blog.Title, blog.CreatedAt, blog.Content, blog.Avatar, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func deleteBlog(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM blogs WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

func getBlogsPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	search := c.DefaultQuery("search", "%")

	rows, err := db.Query("SELECT id, title, created_at, avatar, content, user_id FROM blogs WHERE id LIKE ? OR title LIKE ? OR user_id LIKE ? ORDER BY created_at DESC  LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.ID, &blog.Title, &blog.CreatedAt, &blog.Avatar, &blog.Content, &blog.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		blogs = append(blogs, blog)

	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM blogs WHERE id LIKE ? OR title LIKE ? OR user_id LIKE ?", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultBlogPagenation{pageCount, blogs}
	c.JSON(http.StatusOK, result)
}

func getBlogsByUserIDPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "%")

	user_id := c.Param("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	rows, err := db.Query("SELECT id, title, created_at, avatar, content, user_id FROM blogs WHERE user_id = ? AND (id LIKE ? OR title LIKE ?) ORDER BY created_at DESC  LIMIT ? OFFSET ? ", user_id, search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.ID, &blog.Title, &blog.CreatedAt, &blog.Avatar, &blog.Content, &blog.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		blogs = append(blogs, blog)

	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM blogs WHERE user_id = ? AND (id LIKE ? OR title LIKE ?)", user_id, search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultBlogPagenation{pageCount, blogs}
	c.JSON(http.StatusOK, result)
}

// SALES
func getSale(c *gin.Context) {

	month := c.DefaultQuery("month", "")
	year := c.DefaultQuery("year", "")

	rows, err := db.Query("SELECT id, room_type_id, time_start, time_end FROM sales WHERE MONTH(time_end) = ? AND YEAR(time_end) = ?", month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var saleResult SaleResult

	moneys := 0
	for rows.Next() {
		var sale Sale
		err := rows.Scan(&sale.ID, &sale.RoomTypeID, &sale.TimeStart, &sale.TimeEnd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var roomType RoomType
		err = db.QueryRow("SELECT price FROM room_types WHERE id = ?", sale.RoomTypeID).Scan(&roomType.Price)

		timeStartParse, err1 := time.Parse("2006-01-02 15:04:05", sale.TimeStart)
		timeEndParse, err2 := time.Parse("2006-01-02 15:04:05", sale.TimeEnd)

		if err1 != nil || err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		duration := timeEndParse.Sub(timeStartParse)
		days := int(duration.Hours()/24) + 1

		moneys += roomType.Price * days

	}
	saleResult.Time = month
	saleResult.SaleValue = moneys

	c.JSON(http.StatusOK, saleResult)
}

// Upload
// func uploadImage(c *gin.Context) {
// 	file, header, err := c.Request.FormFile("image")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	imageBytes := make([]byte, header.Size)
// 	_, err = file.Read(imageBytes)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	_, err = db.Exec("INSERT INTO images (name, data) VALUES (?, ?)", header.Filename, imageBytes)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})

// }

// func getImages(c *gin.Context) {
// 	rows, err := db.Query("SELECT id, name, data FROM images")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	var images []Image
// 	for rows.Next() {
// 		var image Image
// 		rows.Scan(&image.ID, &image.Name, &image.Data)
// 		images = append(images, image)
// 	}

// 	c.JSON(http.StatusOK, images)
// }

// Accounts

func getAccountByID(c *gin.Context) {
	id := c.Param("id")
	var account Account

	err := db.QueryRow("SELECT id, email, password, name, avatar_name, avatar_data, type, phone_number FROM accounts WHERE id = ?", id).Scan(&account.ID, &account.Email, &account.Password, &account.Name, &account.AvatarName, &account.AvatarData, &account.Type, &account.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func getAccountByEmail(c *gin.Context) {
	email := c.Param("email")
	var account Account

	err := db.QueryRow("SELECT id, email, password, name, avatar_name, avatar_data, type, phone_number FROM accounts WHERE email = ?", email).Scan(&account.ID, &account.Email, &account.Password, &account.Name, &account.AvatarName, &account.AvatarData, &account.Type, &account.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func createAccount(c *gin.Context) {
	var account Account

	jsonData := c.PostForm("account")
	err := json.Unmarshal([]byte(jsonData), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	imageBytes := make([]byte, header.Size)
	_, err = file.Read(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account.AvatarName = header.Filename
	account.AvatarData = imageBytes

	result, err := db.Exec("INSERT INTO accounts (email, password, name, avatar_name, avatar_data, type, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?)", account.Email, account.Password, account.Name, header.Filename, imageBytes, account.Type, account.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	account.ID = int(id)

	sentEmailAccount(account.Name, account.Email, account.Password)

	c.JSON(http.StatusCreated, account)
}

func updateAccount(c *gin.Context) {
	id := c.Param("id")
	var account Account

	jsonData := c.PostForm("account")
	err := json.Unmarshal([]byte(jsonData), &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	imageBytes := make([]byte, header.Size)
	_, err = file.Read(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account.AvatarName = header.Filename
	account.AvatarData = imageBytes

	_, err = db.Exec("UPDATE accounts SET email = ?, password = ?, name = ?, avatar_name = ?, avatar_data = ?, type = ?, phone_number = ? WHERE id = ?", account.Email, account.Password, account.Name, header.Filename, imageBytes, account.Type, account.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sentEmailAccount(account.Name, account.Email, account.Password)

	c.JSON(http.StatusOK, account)
}

func updateAccountWithoutAvatar(c *gin.Context) {
	id := c.Param("id")
	var account Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE accounts SET email = ?, password = ?, name = ?, type = ?, phone_number = ? WHERE id = ?", account.Email, account.Password, account.Name, account.Type, account.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sentEmailAccount(account.Name, account.Email, account.Password)

	c.JSON(http.StatusOK, account)
}

func deleteAccount(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM accounts WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

func getAccountsPaged(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	offset := (page - 1) * limit

	search := c.DefaultQuery("search", "%")

	rows, err := db.Query("SELECT id, email, password, name, avatar_name, avatar_data, type, phone_number FROM accounts WHERE id LIKE ? OR email LIKE ? OR name LIKE ? OR type LIKE ? OR phone_number LIKE ?  LIMIT ? OFFSET ?", search+"%", search+"%", search+"%", search+"%", search+"%", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.Email, &account.Password, &account.Name, &account.AvatarName, &account.AvatarData, &account.Type, &account.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accounts = append(accounts, account)
	}

	var pageCount int
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM accounts WHERE id LIKE ? OR email LIKE ? OR name LIKE ? OR type LIKE ? OR phone_number LIKE ?", search+"%", search+"%", search+"%", search+"%", search+"%").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count%limit != 0 || count == 0 {
		pageCount = count/limit + 1
	} else {
		pageCount = count / limit
	}

	result := ResultAccountPagenation{pageCount, accounts}
	c.JSON(http.StatusOK, result)
}
