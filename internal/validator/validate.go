package validator

import (
	"cmp"
	"regexp"
	"slices"
	"sort"
)

/*
Handle user input validation... data
*/

var EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

// validate struct to hold user errors
type ValidateData struct {
	Errors map[string]interface{}
}

// initialize the validate object
func NewValidator() *ValidateData {
	return &ValidateData{Errors: make(map[string]interface{})}
}

func (v *ValidateData) CheckErrorExists() bool {
	return len(v.Errors) == 0
}

// add user error to validate object
func (v *ValidateData) AddError(key, message string) {
	_, ok := v.Errors[key]
	if !ok {
		// add errors to object
		v.Errors[key] = message
	}
}

// check if the data inputed is an error
func (v *ValidateData) CheckIsError(ok bool, key, message string) {
	if ok {
		v.AddError(key, message)
	}
}

// check if a value exist in a slice
func In[T cmp.Ordered](val T, item ...T) bool {

	sort.Slice(item, func(i, j int) bool {
		return item[i] < item[j]
	})
	_, found := slices.BinarySearch(item, val)
	return found
}

type ItemInterface interface {
	~int | ~string | ~float64
}

// check value is unique
func CheckDuplicate[T ItemInterface](val []string) bool {
	notUnique := make(map[string]bool)
	for _, val := range val {
		if _, ok := notUnique[val]; ok {
			return true
		} else {
			notUnique[val] = true
		}
	}
	return false
}

var movieStatus = []string{"released", "upcoming"}

func (v *ValidateData) CheckMovieStatus(val string) bool {
	return In(val, movieStatus...)
}

func Isimage(mime string) bool {
	mimeType := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	return In[string](mime, mimeType...)
}

func MatchPattern(pattern, input string) bool {
	rx := regexp.MustCompile(pattern)
	return rx.MatchString(input)
}
