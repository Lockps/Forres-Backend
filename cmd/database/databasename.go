package database

type DBName struct {
	Admin    string
	Staff    string
	Customer string
	Users    string
}

var dbname DBName

func Init() {
	dbname = DBName{
		Admin:    "Admin",
		Staff:    "Staff",
		Customer: "Customer",
		Users:    "Users"}
}

func GetLocation(permission int) string {
	Init()
	if permission == 0 {
		return dbname.Users
	}
	if permission == 1 {
		return dbname.Customer
	}
	if permission == 2 {
		return dbname.Staff
	}
	if permission == 3 {
		return dbname.Admin
	}

	return ""
}
