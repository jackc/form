package form

// import (
// 	"time"
// )

type NewGameForm struct {
	Date     string
	LastName string
	Age      int32
}

/*
p.FirstName = form.values["FirstName"]
p.LastName = form.values["LastName"]
p.Age, err = strconv.ParseInt(form.values["Age"], 10, 32)
if err != nil {

}

presence

string
length min/max
regex pattern
options

integer
min max

float
min max

uniqueness

goals:
map form to struct
 * type conversions
 * differentiate between none / 0 / ""
declarative validations? (for JS client side)


*/
