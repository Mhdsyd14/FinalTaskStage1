package main

import (
	"context"
	"fmt"
	"golang/connection"
	"golang/middleware"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Projek struct{
	Id int
	Name string
	StartDate time.Time
	EndDate time.Time
	Duration string
	Description string
	Technologies []string
	ReactJs bool
	NodeJs bool
	JavaScript bool
	Golang bool
	Image string
	Author int
	Email string
}

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
}

var ProjectData = []Projek{

}

type UserLoginSession struct {
	IsLogin bool
	Name    string
}

var userLoginSession = UserLoginSession{}


func main() {
	e := echo.New()
	e.Static("/Aset","Aset")
	e.Static("/Gambar", "uploads")

	connection.DatabaseConnect()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("cruddatabase"))))

	e.GET("/",Hello)
	e.GET("/Home",Home)
	e.GET("/ContactMe",Kontak)
	e.GET("/Project",Project)
	e.GET("/Edit/:id",EditProject)
	e.GET("/Register", Register)
	e.GET("/Login", Login)
	e.GET("/Welcome", Welcome)

	e.POST("/AddProjek",middleware.UploadFile(AddProjek))
	e.POST("/RegisterAkun",RegisterAkun)
	e.POST("/LoginAkun",LoginAkun)
	e.POST("/EditProjek",middleware.UploadFile(ProjectEdit))
	e.POST("/Delete/:id",DeleteProject)
	e.POST("/logout", logout)




	e.Logger.Fatal(e.Start("localhost:5000"))
}

func Hello(c echo.Context)error {
	return c.JSON(http.StatusOK, map[string]string{
		"Nama":"Muhammad Irsyad",
		"Alamat" : "JL.Mampang Prapatan 5",
	})
	}

func Home(c echo.Context)error {
	tmpl,err := template.ParseFiles("Views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}
sess, errSess := session.Get("irsyad", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = true
		userLoginSession.Name = sess.Values["name"].(string)
	}


	
 dataBlogs, errBlogs := connection.Conn.Query(context.Background(),"SELECT p.id, p.name, p.start_date, p.end_date, p.technologies, p.description, p.image, p.author,u.email FROM tb_projects p JOIN tb_user u ON p.author = u.id WHERE author = $1",sess.Values["id"])

	if errBlogs != nil {
		return c.JSON(http.StatusInternalServerError, errBlogs.Error())
	}

    var resultProjeks []Projek
	for dataBlogs.Next() {
		var each  = Projek{}

		err := dataBlogs.Scan(&each.Id, &each.Name, &each.StartDate, &each.EndDate, &each.Technologies, &each.Description, &each.Image,&each.Author,&each.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		each.Duration = CountDuration(each.StartDate,each.EndDate)
	
		

		resultProjeks = append(resultProjeks, each)
	}
	

	// data := map[string]interface{}{
	// 	"Projek": resultProjeks,
	// }

	

	data := map[string]interface{}{
		"Projek" : resultProjeks,
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],  
		"FlashNama" : sess.Values["name"],
		"UserLoginSession": userLoginSession,
	}
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	

	
	return tmpl.Execute(c.Response(),data)
}

func Kontak(c echo.Context)error {
	tmpl,err := template.ParseFiles("Views/myproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}
	
	return tmpl.Execute(c.Response(),nil)
}

func EditProject(c echo.Context)error {
	id := c.Param("id")


	tmpl,err := template.ParseFiles("Views/EditProject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}

	idToInt,_ := strconv.Atoi(id)

	ProjectDetail := Projek{}

 	errQuery := connection.Conn.QueryRow(context.Background(),"SELECT * FROM tb_projects WHERE id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description,  &ProjectDetail.Technologies, &ProjectDetail.Image, &ProjectDetail.Author)
	
	if errQuery != nil{
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ProjectDetail.Duration = CountDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	data := map[string]interface{}{ 
		"Id":   id,
		"Project": ProjectDetail,
		"startDateString": ProjectDetail.StartDate.Format("2006-01-02"),
		"endDateString" : ProjectDetail.EndDate.Format("2006-01-02"),

	}

	fmt.Println(ProjectDetail)

	return tmpl.Execute(c.Response(),data)
}

func ProjectEdit(c echo.Context)error {
	Id := c.FormValue("id") 
	NamaProjek := c.FormValue("NamaProject")
	TanggalMulai := c.FormValue("TanggalMulai")
	TanggalAkhir := c.FormValue("TanggalAkhir")
	Deskripsi := c.FormValue("Deskripsi")
	Gambar := c.Get("FileNama").(string)
	ReactJs := c.FormValue("ReactJs")
	JavaScript := c.FormValue("JavaScript")
	NodeJs := c.FormValue("NodeJs")
	Golang := c.FormValue("Golang")

	// _, err := strconv.Atoi(id)
	idToInt,_ := strconv.Atoi(Id)

	fmt.Println(idToInt)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }


	data, err := connection.Conn.Exec(context.Background(),"UPDATE tb_projects SET name=$1, start_date=$2, end_date=$3, description=$4, technologies[1]=$5,  technologies[2]=$6,  technologies[3]=$7,  technologies[4]=$8, image=$9 WHERE id=$10", NamaProjek, TanggalMulai, TanggalAkhir, Deskripsi, ReactJs, NodeJs, JavaScript, Golang, Gambar, Id)
 
	fmt.Println(data)

	if err != nil {
	c.JSON(http.StatusInternalServerError,err.Error())
	
	}



return c.Redirect(http.StatusMovedPermanently, "/Home")
}




func Project(c echo.Context)error {
	tmpl,err := template.ParseFiles("Views/AddProject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}

		sess, errSess := session.Get("irsyad", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = true
		userLoginSession.Name = sess.Values["name"].(string)
	}

	data := map[string]interface{}{
	
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],  
		"FlashNama" : sess.Values["name"],
		"IdUser" : sess.Values["id"],
		"UserLoginSession": userLoginSession,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	

	
	return tmpl.Execute(c.Response(),data)
}

func Welcome(c echo.Context)error {
	tmpl,err := template.ParseFiles("Views/selesai.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}
	
	return tmpl.Execute(c.Response(),nil)
}

func Register(c echo.Context)error {
	tmpl,err := template.ParseFiles("Views/Register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}
	
	return tmpl.Execute(c.Response(),nil)
}


func RegisterAkun(c echo.Context)error {
	inputname := c.FormValue("inputnama")
	inputemail := c.FormValue("inputemail")
	inputpassword := c.FormValue("inputpassword")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputpassword),10)

	if err != nil{
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	query, err := connection.Conn.Exec(context.Background(),"INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)",inputname, inputemail, hashedPassword)
	
	if err != nil{
		return c.JSON(http.StatusInternalServerError, err.Error())
	}


	fmt.Println("affected row : ", query.RowsAffected())

	return Irsyad(c, "Register berhasil!", true, "/Login")
	
	// return c.Redirect(http.StatusMovedPermanently,"/Login")
	
}

func LoginAkun(c echo.Context)error {
	inputemail := c.FormValue("emailinput")
	inputpassword := c.FormValue("passwordinput")


	user := User{}

	err := connection.Conn.QueryRow(context.Background(),"SELECT id, name, email, password FROM tb_user WHERE email=$1", inputemail).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword)

	if err != nil{
			return Irsyad(c, "Login gagal!", false, "/Login")
		// return c.JSON(http.StatusInternalServerError,err.Error())
	}

	
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(inputpassword))

	if errPassword != nil {
		// return c.JSON(http.StatusInternalServerError,err.Error())
		return Irsyad(c, "Login gagal!", false, "/Login")
	}

	sess, _ := session.Get("irsyad", c)
	sess.Options.MaxAge = 10800 // 
	sess.Values["message"] = "Login Berhasil!"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["email"] = user.Email
	sess.Values["id"] = user.Id
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	



return c.Redirect(http.StatusMovedPermanently,"/Home")
}

func Login(c echo.Context)error {


	tmpl,err := template.ParseFiles("Views/Login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
		
	}

	sess, errSess := session.Get("irsyad", c)
	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	

	flash := map[string]interface{}{
		"FlashMessage": sess.Values["message"],
		"FlashStatus":  sess.Values["status"],  
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	
	


	return tmpl.Execute(c.Response(),flash)

}



func AddProjek(c echo.Context)error {
	NamaProjek := c.FormValue("NamaProject")
	TanggalMulai := c.FormValue("TanggalMulai")
	TanggalAkhir := c.FormValue("TanggalAkhir")
	Deskripsi := c.FormValue("Deskripsi")
	ReactJs := c.FormValue("ReactJs")
	JavaScript := c.FormValue("JavaScript")
	NodeJs := c.FormValue("NodeJs")
	Golang := c.FormValue("Golang")
	Author := c.FormValue("author")

	Gambar := c.Get("FileNama").(string)

	// fmt.Println(NamaProjek)
	// fmt.Println(TanggalMulai)
	// fmt.Println(TanggalAkhir)
	// fmt.Println(Deskripsi)
	// fmt.Println(ReactJs)
	// fmt.Println(JavaScript)
	// fmt.Println(NodeJs)
	// fmt.Println(Golang)

_, err := connection.Conn.Exec(context.Background(),"INSERT INTO tb_projects(name, start_date, end_date, description, technologies[1], technologies[2], technologies[3], technologies[4], image, author) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", NamaProjek, TanggalMulai, TanggalAkhir, Deskripsi, ReactJs, NodeJs, JavaScript, Golang, Gambar, Author)
 
if err != nil {
	c.JSON(http.StatusInternalServerError,err.Error())
	
}



return c.Redirect(http.StatusMovedPermanently, "/Home")
}

func DeleteProject(c echo.Context) error {
	id := c.Param("id") // ID : 1

	idToInt, _ := strconv.Atoi(id)
 
	connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", idToInt)




	return c.Redirect(http.StatusMovedPermanently, "/Home")
}


func Irsyad(c echo.Context, message string, status bool, redirectPath string) error {
	sess, errSess := session.Get("irsyad", c)

	if errSess != nil {
		return c.JSON(http.StatusInternalServerError, errSess.Error())
	}

	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, redirectPath)
}


func logout(c echo.Context) error {
	sess, _ := session.Get("irsyad", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return Irsyad(c, "Logout berhasil!", true, "/Welcome")
}



func CountDuration(d1 time.Time, d2 time.Time)string  {
	diff := d2.Sub(d1)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	month := days / 30 

	if month > 12 {
		return strconv.Itoa(month/12) + "Tahun"
	}
	if month > 0 {
		return strconv.Itoa(month) + "Bulan"
	}
	if weeks > 0 {
		return strconv.Itoa(weeks) + "Minggu"
	}

	return strconv.Itoa(days) + "Hari"
	
}


