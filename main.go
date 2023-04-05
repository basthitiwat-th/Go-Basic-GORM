package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

// Logger
func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n======================================\n", sql)
}

// Declare Global variable
var db *gorm.DB

func main() {
	dsn := "root:Basbm031197#@tcp(localhost:3306)/bond?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	//Connect DB
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	//Automigrate
	db.AutoMigrate(Gender{}, Test{}, Customer{})

	// CreateGender("xxxxx")
	// CreateGender("female")
	GetGenders()
	// GetGender(10)
	// GetGenderByname("male")

	// UpdateGender2(4, "")
	// DeleteGender(4)
	// CreateTest(0, "Test2")
	// CreateTest(0, "Test3")
	// CreateTest(0, "Test4")

	// DeleteTest(2)
	// GetTest()

	// db.Migrator().CreateTable(Customer{})
	// CreateCustomer("BAS", 1)
	//UpdateGender2(1, "TEST")
	// GetCustomers()

}

type Customer struct {
	Id       uint
	Name     string
	Gender   Gender
	GenderID uint
}
type Gender struct {
	Id   uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code uint   `gorm:"comment:This is Code"`
	Name string `gorm:"column:myname;size:20;unique;default:Hello;not null"`
}

func (t Test) TableName() string {
	return "MyTest"
}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return
	}
	for _, customer := range customers {
		fmt.Printf("%v|%v|%v\n", customer.Id, customer.Name, customer.Gender.Name)
	}
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{
		Name:     name,
		GenderID: genderID,
	}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func CreateTest(code uint, name string) {
	test := Test{Code: code, Name: name}
	db.Create(&test)
}
func GetTest() {
	tests := []Test{}
	db.Find(&tests)
	for _, test := range tests {
		fmt.Println(test)

	}
}

func DeleteTest(id uint) {
	db.Unscoped().Delete(&Test{}, id)

	// db.Delete(&Test{}, id)
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id ASC").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}

func GetGender(id uint) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByname(name string) {
	gender := Gender{}
	tx := db.Where("name=?", name).First(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}

func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id = @myid", sql.Named("myid", id)).Updates(gender)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}

func CreateGender(name string) {
	gender := Gender{
		Name: name,
	}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)
}
