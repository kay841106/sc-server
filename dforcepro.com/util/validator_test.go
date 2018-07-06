package util

import (
	//"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type gameInfo struct {
	Name      string `json:"name,omitempty" bson:"name,omitempty" valid:"alphanum"`
	Developer string `json:"develop,omitempty" bson:"develop,omitempty" valid:"alphanum,required"`
	Type      string `json:"type,omitempty" bson:"type,omitempty" valid:"alphanum"`
	platform  platform
	Rated     string `json:"rated,omitempty" bson:"rated,omitempty" valid:"email"`
	Sale      int    `json:"sale,omitempty" bson:"sale,omitempty" valid:"int"`
	test      []nestedNest
}

type platform struct {
	PlatformName string ` valid:"uppercase"`
	nestedNest
}

type nestedNest struct {
	Name string ` valid:"alphanum"`
}

func gameInfoform() *gameInfo {
	return &gameInfo{}
}

func platformform() *platform {
	return &platform{}
}

func ExampleCheckParam() {
	//initialize validator
	InitValidator()
	//Validate mode
	target := map[string][]string{
		"a": []string{"Alphanumeric", "Alpha", "Required"},
		"b": []string{"Alphanumeric", "Alpha"},
		"c": []string{"Alphanumeric", "Alpha"},
		"d": []string{"Numeric"},
	}

	x := map[string][]string{
		"a": []string{"", "applause"},
		"b": []string{"banana", "balista"},
		"c": []string{"catherine", "catal#yst123"},
		"d": []string{"432", "301", "81927l39824"},
	}

	keys := []string{"a", "b", "c", "d"}

	//Simulate input []interface with getQueryValue().
	//This is the param we get when parsing the GET
	input := getQueryValue(x, keys)

	_, detail := CheckParam(*input, target)
	fmt.Println("Result: ", detail)
}

func ExampleCheckStruct() {
	InitValidator()
	people := struct {
		Name string `valid:"required,alpha"`
		Age  int    `valid:"required,numeric"`
	}{"Gorden", 22}
	_, err := CheckStruct(people)
	fmt.Println("result: ", err)
}

func ExampleCheckRangeInt() {
	input := "33625"
	//Leave the array empty if the range test is not necessary
	bound := []int{500, 60000}
	result, err := IsInt(input, bound)
	fmt.Println("result: ", result, "\n", "err: ", err)
}

func ExampleCheckRangeFloat64() {
	input := "23569.872340"
	//Leave the array empty if the range test is not necessary
	bound := []float64{500.01, 30000.588}
	result, err := IsFloat64(input, bound)
	fmt.Println("result: ", result, "\n", "err: ", err)
}

//This test includes two sets
func Test_CheckParam(t *testing.T) {

	//Validate mode

	at := []string{"Alphanumeric", "Alpha"}
	bt := []string{"Alphanumeric", "Alpha"}
	ct := []string{"Alphanumeric", "Alpha"}
	dt := []string{"Numeric"}
	target := map[string][]string{
		"a": at,
		"b": bt,
		"c": ct,
		"d": dt,
	}

	//Test set 1
	standardOutput := make(map[string]string)

	standardOutput["c"] = "[Check catal#yst123 with Alphanumeric failed][Check catal#yst123 with Alpha failed]"
	standardOutput["d"] = "[Check 81927l39824 with Numeric failed]"

	a := []string{"apple", "applause"}
	b := []string{"banana", "balista"}
	c := []string{"catherine", "catal#yst123"}
	d := []string{"432", "301", "81927l39824"}

	x := map[string][]string{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
	}

	keys := []string{"a", "b", "c", "d"}

	//Simulate input []interface with getQueryValue().
	//This is the param we get when parsing the GET

	input := getQueryValue(x, keys)

	//fmt.Println("input: ", x)
	//fmt.Println("filter: ", target)
	_, detail := CheckParam(*input, target)

	fmt.Println("Result: ", detail)

	assert.Equal(t, detail, standardOutput, "The two words should be the same.")

	//Test set 2

	standardOutput = make(map[string]string)

	standardOutput["d"] = "[Check monopolosomeplace.com with Email failed]"

	a = []string{"op.gg", "www.yahoo.com.tw"}
	b = []string{"banana", "balista"}
	c = []string{"catherine", "catalyst"}
	d = []string{"tig4605246@gmail.com", "monopolosomeplace.com"}

	x = map[string][]string{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
	}

	keys = []string{"a", "b", "c", "d"}

	//Simulate input []interface with getQueryValue().
	//This is the param we get when parsing the GET

	input = getQueryValue(x, keys)

	dt2 := []string{"Email"}
	at2 := []string{"DNS"}
	target["a"] = at2
	target["d"] = dt2
	//fmt.Println("input: ", x)
	//fmt.Println("filter: ", target)
	_, detail = CheckParam(*input, target)
	//fmt.Println("Result: ",detail,"\nexpected: ",standardOutput)
	assert.Equal(t, detail, standardOutput, "The two words should be the same.")

}

func Test_IsInt(t *testing.T) {
	intRange := []int{10, 100}
	input := "60"
	//fmt.Println("input: ", input)
	//fmt.Println("range: ", intRange)
	result, _ := IsInt(input, intRange)
	assert.Equal(t, result, 60, "The two integer should be the same.")

}

func Test_IsFloat64(t *testing.T) {
	float64Range := []float64{10.0, 1000.0}
	input := "650.356541234"
	//fmt.Println("input: ", input)
	//fmt.Println("range: ", float64Range)
	result, _ := IsFloat64(input, float64Range)
	assert.Equal(t, result, 650.356541234, "The two float64 should be the same.")

}

func Test_CheckStruct(t *testing.T) {
	//standardOutput := errors.New("platform: Xbox does not validate as uppercase")
	InitValidator()
	testStruct := generateGameInfo(true)
	//fmt.Println("input: ", testStruct)
	_, err := CheckStruct(testStruct)
	//fmt.Println(standardOutput, testStruct, err)
	assert.Equal(t, err, nil, "The two results should be the same.")

}

func generateGameInfo(isLegal bool) *gameInfo {
	newGame := gameInfoform()
	newGame.Name = "TheElderScrollVSkyrim"
	newGame.Developer = "Bethesda"
	newGame.Type = "RPG"
	newGame.platform.PlatformName = "PS4"
	newGame.Rated = "konawiba@gmail.com"
	newGame.Sale = 0
	newGame.test = []nestedNest{nestedNest{"&&&"}}

	return newGame
}

func getQueryValue(queries map[string][]string, keys []string) *map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range keys {
		value, ok := queries[key]
		if !ok {
			continue
		}
		if len(value) == 1 {
			result[key] = value[0]
		} else {
			result[key] = value
		}
	}
	return &result
}
