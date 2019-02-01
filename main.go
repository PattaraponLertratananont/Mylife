package main

//! Golang(API) =>Echo + Mondgodb

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Profile struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	Passwaord string        `json:"password" bson:"password"`
}

var db *mgo.Session

type (
	Handler struct {
		DB *mgo.Session
	}
)

const (
	mongo_host = "mongodb://admin:muyon@mylife-shard-00-00-3kask.gcp.mongodb.net:27017,mylife-shard-00-01-3kask.gcp.mongodb.net:27017,mylife-shard-00-02-3kask.gcp.mongodb.net:27017/test?&replicaSet=MyLife-shard-0&authSource=admin"
)

func main() {
	h := &Handler{DB: db}
	//+ Echo instance
	e := echo.New()

	//+ Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	//+ Route =>handler
	//* Hi!
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hi!, from mongod+API(echo)")
	})
	//* Query all data
	e.GET("/read", h.Getdata)
	e.POST("/post", h.Postdata)
	e.DELETE("/delete", h.Deletedata)
	e.PUT("/update", h.Updatedata)

	//+ Start server
	e.Logger.Fatal(e.Start(getPort()))
}
func getPort() string {
	var port = os.Getenv("PORT") // ----> (A)
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port // ----> (B)
}

func (h *Handler) Getdata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongo_host)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()

	s := session.DB("MyLife").C("userpass")
	var profiles []Profile
	err = s.Find(bson.M{}).Limit(100).All(&profiles)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}

	return c.JSON(http.StatusOK, profiles)

}

func (h *Handler) Deletedata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongo_host)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()
	var messege Profile
	if err = c.Bind(&messege); err != nil {
		return
	}
	s := session.DB("MyLife").C("userpass")
	fmt.Println(messege)
	err = s.Remove(messege)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}
	return c.JSON(http.StatusCreated, "success")

}

func (h *Handler) Postdata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongo_host)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()
	var messege Profile
	if err = c.Bind(&messege); err != nil {
		return
	}
	//strmsg, err := json.Marshal(messege)
	messege.ID = bson.NewObjectId()
	s := session.DB("MyLife").C("userpass")
	fmt.Println(messege)
	err = s.Insert(messege)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}
	return c.JSON(http.StatusCreated, "success")

}

func (h *Handler) Updatedata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongo_host)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()
	var messege Profile
	if err = c.Bind(&messege); err != nil {
		return
	}
	s := session.DB("MyLife").C("userpass")
	fmt.Println(messege)
	err = s.UpdateId(messege.ID, &messege)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}
	return c.JSON(http.StatusCreated, "success")

}
