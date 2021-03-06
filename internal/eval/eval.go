package eval

import (
	"context"
	"errors"
	"fmt"
	"git.internal.yunify.com/qxp/misc/logger"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//!+env

// Env environment variable
type Env map[Var]Substance

//!-env

//!+Eval1

// Eval var eval
func (v Var) Eval(env Env) Substance {
	s := string(v)
	if strings.HasSuffix(s,"'") && strings.HasPrefix(s,"'"){
		return &String{s[1:len(s)-1]}
	}
	return env[v]
}

// Eval literal eval
func (l literal) Eval(_ Env) Substance {
	return &Float64{float64(l)}
}

//!-Eval1

//!+Eval2

// Eval unary eval
func (u unary) Eval(env Env) Substance {
	switch u.op {
	case '+':
		return &Float64{+u.x.Eval(env).Float64()}
	case '-':
		return &Float64{-u.x.Eval(env).Float64()}
	}
	return &Float64{}
}

// Eval binary eval
func (b binary) Eval(env Env) Substance {
	err := b.Check(GenMapFromEnv(env))
	if err != nil {
		return &String{err.Error()}
	}
	switch b.op {
	case '+':
		return &Float64{b.x.Eval(env).Float64() + b.y.Eval(env).Float64()}
	case '-':
		return &Float64{b.x.Eval(env).Float64() - b.y.Eval(env).Float64()}
	case '*':
		return &Float64{b.x.Eval(env).Float64() * b.y.Eval(env).Float64()}
	case '/':
		return &Float64{b.x.Eval(env).Float64() / b.y.Eval(env).Float64()}
	case '%':
		return &Float64{math.Mod(b.x.Eval(env).Float64(),b.y.Eval(env).Float64())}
	case '>':
		return &Boolean{b.x.Eval(env).Float64() > b.y.Eval(env).Float64()}
	case '<':
		return &Boolean{b.x.Eval(env).Float64() < b.y.Eval(env).Float64()}
	case '≥':
		return &Boolean{b.x.Eval(env).Float64() >= b.y.Eval(env).Float64()}
	case '≤':
		return &Boolean{b.x.Eval(env).Float64() <= b.y.Eval(env).Float64()}
	case '≡':
		return &Boolean{b.x.Eval(env).String() == b.y.Eval(env).String()}
	case '≠':
		return &Boolean{b.x.Eval(env).String() != b.y.Eval(env).String()}
	case '∩':
		return &Boolean{b.x.Eval(env).Boolean() || b.y.Eval(env).Boolean()}
	case '∪':
		return &Boolean{b.x.Eval(env).Boolean() && b.y.Eval(env).Boolean()}
	}
	return &Float64{}
}

// Eval section eval
func (s section) Eval(env Env) Substance {
	err := s.Check(GenMapFromEnv(env))
	if err != nil {
		return &String{err.Error()}
	}
	var ans []string
	for _, arg := range s.y {
		ans = append(ans,arg.Eval(env).String())
	}
	switch s.op {
	case '∈':
		return &Boolean{in(s.x.Eval(env).String(),&ans)}
	case '∉':
		return &Boolean{!in(s.x.Eval(env).String(),&ans)}
	}
	return &Float64{}
}
// Result response
type Result interface{}

// Handler Handler
func Handler(c context.Context,expr string,param map[string]interface{}) (Result,error) {
	expr = symbolReplace(expr)
	er, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	env := GenEnv(param)
	r := er.Eval(env)
	if r.String() == errCodeMissParam{
		return nil,errors.New("miss params")
	}
	if r.String() == errCodeWithoutFun{
		return nil,errors.New("expression error")
	}
	return r,nil
}

func symbolReplace(expr string)string  {
	for k,v := range op{
		if strings.Contains(expr,k){
			expr = strings.ReplaceAll(expr,k,v)
		}
	}
	return expr
}

// GenEnv GenEnv
func GenEnv(param map[string]interface{}) Env  {
	if param == nil {
		return map[Var]Substance{}
	}
	env := make(map[Var]Substance,len(param))
	for k,v := range param{
		if strings.Contains(k,":"){
			k = strings.Replace(k,":","_",1)
		}
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			env[Var(k)] = &String{reflect.ValueOf(v).String()}
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			env[Var(k)] = &Float64{float64(reflect.ValueOf(v).Int())}
		case reflect.Float32, reflect.Float64:
			env[Var(k)] = &Float64{reflect.ValueOf(v).Float()}
		}
	}
	return env
}

// GenMapFromEnv GenMapFromEnv
func GenMapFromEnv(env Env)map[Var]bool  {
	r := make(map[Var]bool,0)
	for k := range env{
		r[k] = true
	}
	return r
}
// Eval call eval
func (c call) Eval(env Env) Substance {
	if f, ok := formulas[c.fn];ok{
		err := c.Check(GenMapFromEnv(env))
		if err != nil {
			fmt.Println(err)
			return &String{err.Error()}
		}
		return f(c.args,env)
	}
	logger.Logger.Errorw(ErrNoFormula.Error())
	return &String{errCodeWithoutFun}
}

//!-Eval2

// T type
type T string

const (
	stringT  T = "string"
	float64T T = "float64"
	boolT    T = "bool"
)

// Substance substance
type Substance interface {
	GetType() T
	String() string
	Float64() float64
	Boolean() bool
}

// String string
type String struct {
	Val string  `json:"result"`
}

// String to string
func (s *String) String() string {
	return s.Val
}

// Float64 to float64
func (s *String) Float64() float64 {
	value, err := strconv.ParseFloat(s.Val, 64)
	if err != nil {
		return 0
	}
	return value
}
// Boolean to bool
func (s *String) Boolean() bool {
	value, err := strconv.ParseBool(s.Val)
	if err != nil {
		return false
	}
	return value
}


// GetType get string type
func (s *String) GetType() T {
	return stringT
}

// Float64 float64
type Float64 struct {
	Val float64 `json:"result"`
}

// String to string
func (f *Float64) String() string {
	return strconv.FormatFloat(f.Val, 'f', -1, 64)
}

// Float64 to float64
func (f *Float64) Float64() float64 {
	return f.Val
}

// Boolean to bool
func (f *Float64) Boolean() bool {
	s := strconv.FormatFloat(f.Val, 'f', -1, 64)
	value, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return value
}

// GetType get float64 type
func (f *Float64) GetType() T {
	return float64T
}

// Boolean bool
type Boolean struct {
	Val bool  `json:"result"`
}

// String to string
func (b *Boolean) String() string {
	return strconv.FormatBool(b.Val)
}

// Float64 to float64
func (b *Boolean) Float64() float64 {
	value, err := strconv.ParseFloat(strconv.FormatBool(b.Val), 64)
	if err != nil {
		return 0
	}
	return value
}

// Boolean to bool
func (b *Boolean) Boolean() bool {
	return b.Val
}

// GetType get float64 type
func (b *Boolean) GetType() T {
	return boolT
}

// QuickSort QuickSort
func QuickSort(arr []float64) []float64 {
	_sort(arr, 0, len(arr)-1)
	return arr
}

func _sort(arr []float64, left int, right int){
	if left >= right {
		return
	}
	temp := arr[left]
	start := left
	stop := right
	for right != left {
		for right > left && arr[right] >= temp  {
			right --
		}
		for left < right && arr[left] <= temp  {
			left ++
		}
		if right > left {
			arr[right], arr[left] = arr[left], arr[right]
		}
	}
	arr[right], arr[start] = temp, arr[right]
	_sort(arr, start, left)
	_sort(arr, right+1, stop)
}

func in(target string, filter *[]string) bool {
	index := sort.SearchStrings(*filter, target)
	if index < len(*filter) && (*filter)[index] == target {
		return true
	}
	return false
}
