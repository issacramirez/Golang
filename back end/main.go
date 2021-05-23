package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/denisenkom/go-mssqldb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const SecretKey = "secret"

// modelo de la tabla Student
type (
	Student struct {
		Id     int    `json:"id"`
		Nombre string `json:"nombre"`
		Edad   string `json:"edad"`
	}

	JwtClaims struct {
		Name string `json:"name"`
		jwt.StandardClaims
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

// funcion para "iniciar sesion" y generar el token
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// autorizacion de usuario
	if username != "sa" || password != "1234" {
		return echo.ErrUnauthorized
	}
	token, err := createJwtToken()
	if err != nil {
		log.Println("error creando el jwt token", err)
		return c.String(http.StatusInternalServerError, "Algo salió mal al intentar crear el token")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "You are logged in",
		"token":   token,
	})
}

// genera el token
func createJwtToken() (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["name"] = "José Issac"
	claims["type"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err := t.SignedString([]byte(SecretKey))
	if err != nil {
		return "Error creando el token", err
	}
	return token, nil
}

func main() {
	e := echo.New()
	// grupo para los end points que necesitan autorizacion jwt
	jwtGroups := e.Group("/jwt")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323"},
		AllowHeaders: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// prueba de que funciona la conexion a DB
	connectionSql()

	// restriccion que tendrá autorizacion JWT
	jwtGroups.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(SecretKey),
	}))

	// este endpoint no tiene seguridad con el fin de que entregue el token
	e.POST("/login", login)

	// peticiones para url
	jwtGroups.GET("/Students", AllStudents)
	jwtGroups.POST("/Students", NewStudent)
	jwtGroups.DELETE("/Students/:id", DeleteStudent)
	jwtGroups.GET("/Students/:id", GetStudent)
	jwtGroups.PUT("/Students/:id", UpdateStudent)

	// inicio del servidor
	e.Logger.Fatal(e.Start("localhost:1323"))
}
