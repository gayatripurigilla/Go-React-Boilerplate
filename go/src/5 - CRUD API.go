package main

import (
   "fmt"
   "github.com/gin-contrib/cors"
   "github.com/gin-gonic/gin"
   "github.com/jinzhu/gorm"
   _ "github.com/jinzhu/gorm/dialects/sqlite"           // If you want to use mysql or any other db, replace this line
)

var db *gorm.DB                                         // declaring the db globally
var err error

type Person struct {
   ID uint `json:"id"`
   FirstName string `json:"firstname"`
   LastName string `json:"lastname"`
   City string `json:"city"`
}

func main() {
   db, err = gorm.Open("sqlite3", "./gorm.db")
   // If you want to use MySQL, use the following line
   // db, _ = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
   if err != nil {
      fmt.Println(err)
   }
   defer db.Close()

   db.AutoMigrate(&Person{})
   r := gin.Default()
   r.GET("/people/", GetPeople)                             // Creating routes for each functionality
   r.GET("/people/:id", GetPerson)
   r.POST("/people", CreatePerson)
   r.PUT("/people/:id", UpdatePerson)
   r.DELETE("/people/:id", DeletePerson)
   r.Use(cors.Default())
   r.Run(":8080")                                           // Run on port 8080
}


func DeletePerson(c *gin.Context) {
   id := c.Params.ByName("id")
   var person Person
   d := db.Where("id = ?", id).Delete(&person)
   fmt.Println(d)
   c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdatePerson(c *gin.Context) {
   var person Person
   id := c.Params.ByName("id")
   if err := db.Where("id = ?", id).First(&person).Error; err != nil {
      c.AbortWithStatus(404)
      fmt.Println(err)
   }
   c.BindJSON(&person)
   db.Save(&person)
   c.JSON(200, person)
}

func CreatePerson(c *gin.Context) {
   var person Person
   c.BindJSON(&person)
   db.Create(&person)
   c.JSON(200, person)
}

func GetPerson(c *gin.Context) {
   id := c.Params.ByName("id")
   var person Person
   if err := db.Where("id = ?", id).First(&person).Error; err != nil {
      c.AbortWithStatus(404)
      fmt.Println(err)
   } else {
      c.JSON(200, person)
   }
}

func GetPeople(c *gin.Context) {
   var people []Person
   if err := db.Find(&people).Error; err != nil {
      c.AbortWithStatus(404)
      fmt.Println(err)
   } else {
      c.JSON(200, people)
   }
}
