package main

import (
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// modelo de la tabla Student
type (
	Student struct {
		Id     int    `json:"id"`
		Nombre string `json:"nombre"`
		Edad   string `json:"edad"`
	}
)

// conexion a DB en SQL Server
func connectionSql() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1:3306)/goDB?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal("No se pudo establecer conexion a DB: " + err.Error())
		return db, err
	} else {
		log.Println("Conectado a DB! " + db.Name())
		return db, err
	}
}

// trae todos los estudiantes en la DB
func AllStudents(c echo.Context) error {
	db, err := connectionSql()
	if db == nil {
		log.Fatal("FALLÓ CONEXIÓN A DB: " + err.Error())
	}
	var Students []Student
	db.Find(&Students)
	return c.JSON(http.StatusOK, Students)
}

// agrega un nuevo estudiante
func NewStudent(c echo.Context) error {
	student := new(Student)
	student.Nombre = c.FormValue("Nombre")
	student.Edad = c.FormValue("Edad")

	db, err := connectionSql()
	if db == nil {
		log.Fatal("FALLÓ CONEXIÓN A DB: " + err.Error())
	}
	db.Select("Nombre", "Edad").Create(&student)

	return c.JSON(http.StatusCreated, student)
}

// borra un estudiante existente
func DeleteStudent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	db, err := connectionSql()
	if err != nil {
		log.Fatal("FALLÓ CONEXIÓN A DB: " + err.Error())
	}
	var student Student
	db.Where("id = ?", id).Delete(&student)

	return c.NoContent(http.StatusNoContent)
}

// trae un estudiante
func GetStudent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	db, err := connectionSql()
	if err != nil {
		log.Fatal("FALLÓ CONEXIÓN A DB: " + err.Error())
	}
	var student Student

	db.First(&student, id)
	return c.JSON(http.StatusOK, student)
}

// actualiza la edad de un estudiante, se envia el id en la url y la edad nueva en el body
func UpdateStudent(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	db, err := connectionSql()
	if err != nil {
		log.Fatal("FALLÓ CONEXIÓN A DB: " + err.Error())
	}
	var student Student
	db.First(&student, id)

	var nuevaEdad = c.FormValue("Edad")
	student.Edad = nuevaEdad

	db.Save(&student)
	return c.JSON(http.StatusOK, student)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:4200", "http://localhost:3100"},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	// prueba de que funciona la conexion a DB
	connectionSql()

	// peticiones para url
	e.GET("/Students", AllStudents)
	e.POST("/Students", NewStudent)
	e.DELETE("/Students/:id", DeleteStudent)
	e.GET("/Students/:id", GetStudent)
	e.PUT("/Students/:id", UpdateStudent)

	e.Logger.Fatal(e.Start("localhost:1323"))
}
