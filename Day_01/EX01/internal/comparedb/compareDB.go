package comparedb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// DBReader интерфейс для публичного взаимодействия
type DBReader interface {
	dataDeserialization([]byte) (*recipes, error)
}

type recipes struct {
	Cake []cake `xml:"cake"    json:"cake"`
}

type cake struct {
	Name        string        `xml:"name"    json:"name"`
	Stovetime   string        `xml:"stovetime"    json:"time"`
	Ingredients []ingredients `xml:"ingredients>item"    json:"ingredients"`
}

type ingredients struct {
	Itemname  string `xml:"itemname"    json:"ingredient_name"`
	Itemcount string `xml:"itemcount"    json:"ingredient_count"`
	Itemunit  string `xml:"itemunit"    json:"ingredient_unit"`
}

type dbReaderXML struct{}

func (dbReaderXML) dataDeserialization(b []byte) (*recipes, error) {
	r := recipes{}
	err := xml.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type dbReaderJSON struct{}

func (dbReaderJSON) dataDeserialization(b []byte) (*recipes, error) {
	r := recipes{}
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Deserialization  публичная функция для взаимодействия, с помощью
// которой десериализуется из формата JSON или XML в структуру
func Deserialization(r [2]DBReader, b [2][]byte) ([2]interface{}, error) {
	recipesData := [2]recipes{}
	for i, v := range r {
		obj, err := v.dataDeserialization(b[i])
		if err != nil {
			return [2]interface{}{}, err
		}
		recipesData[i] = *obj
	}
	return [2]interface{}{recipesData[0], recipesData[1]}, nil
}

// ParseFlag обработывает входные данные
func ParseFlag() (*[2]DBReader, *[2][]byte, error) {
	strOld := flag.String("old", "", "file path old")
	strNew := flag.String("new", "", "file path new")
	flag.Parse()

	if len(*strOld) == 0 || len(*strNew) == 0 {
		return nil, nil, errors.New("file path not specified")
	}

	if len(os.Args) != 5 {
		return nil, nil, errors.New("incorrect number of arguments")
	}

	formatOne, dataOne, err := formatFile(strOld)
	if err != nil {
		return nil, nil, err
	}

	formatTwo, dataTwo, err := formatFile(strNew)
	if err != nil {
		return nil, nil, err
	}

	return &[2]DBReader{formatOne, formatTwo}, &[2][]byte{dataOne, dataTwo}, nil
}

func formatFile(str *string) (DBReader, []byte, error) {
	format := strings.Split(*str, ".")

	dataFromFile, err := readFileAll(str)
	if err != nil {
		return nil, nil, err
	}

	if format[len(format)-1] == "xml" {
		return dbReaderXML{}, dataFromFile, nil
	} else if format[len(format)-1] == "json" {
		return dbReaderJSON{}, dataFromFile, nil
	}

	return nil, nil, errors.New("incorrect format")
}

func readFileAll(str *string) ([]byte, error) {
	b, err := os.ReadFile(*str)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Comparison сравнивает старый и новый рецепт, в качестве пустого интерфейса должен принимать
// результат функции Deserialization, иначе ошибка
func Comparison(resRecipes [2]interface{}) {
	if reflect.TypeOf(resRecipes[0]) != reflect.TypeFor[recipes]() || reflect.TypeOf(resRecipes[1]) != reflect.TypeFor[recipes]() {
		log.Fatalln("incorrect type")
	}

	old := resRecipes[0].(recipes).Cake
	new := resRecipes[1].(recipes).Cake

	deleteOld := comapreCake(&old, &new)

	appendNewIngredients, deleteOldIngredients := compareTime(&new, &deleteOld)

	compareIngredients(&new, &old, &appendNewIngredients, &deleteOldIngredients)

}

func comapreCake(old, new *[]cake) map[string]cake {
	appendNew := make(map[string]cake)
	deleteOld := make(map[string]cake)

	for _, v := range *old {
		deleteOld[v.Name] = v
	}

	for _, v := range *new {
		appendNew[v.Name] = v
	}

	indexNew := 0
	indexOld := 0

	for {
		if removeOrAddCake(new, deleteOld, "ADDED cake", &indexNew) == 0 &&
			removeOrAddCake(old, appendNew, "REMOVED cake", &indexOld) == 0 {
			break
		}
	}

	return deleteOld
}

func removeOrAddCake(cakes *[]cake, cakeCheck map[string]cake, msg string, index *int) int {
	for _, v := range (*cakes)[*index:] {
		_, ok := cakeCheck[v.Name]
		if !ok {
			fmt.Printf("%s \"%s\"\n", msg, v.Name)
			if len(*cakes) > *index+1 {
				*cakes = append((*cakes)[:*index], (*cakes)[*index+1:]...)
			} else if len(*cakes)-1 == *index {
				*cakes = (*cakes)[:*index]
			}
			delete(cakeCheck, v.Name)
			break
		}

		*index++
	}

	return len((*cakes)[*index:])
}

func compareTime(new *[]cake, deleteOld *map[string]cake) (map[string]map[string]ingredients, map[string]map[string]ingredients) {
	appendNewIngredients := make(map[string]map[string]ingredients)
	deleteOldIngredients := make(map[string]map[string]ingredients)

	for _, v := range *new {
		if (*deleteOld)[v.Name].Stovetime != v.Stovetime {
			fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", v.Name, v.Stovetime, (*deleteOld)[v.Name].Stovetime)
		}

		appendNewIngredients[v.Name] = make(map[string]ingredients)
		for _, val := range v.Ingredients {
			appendNewIngredients[v.Name][val.Itemname] = val
		}

		deleteOldIngredients[v.Name] = make(map[string]ingredients)
		for _, val := range (*deleteOld)[v.Name].Ingredients {
			deleteOldIngredients[v.Name][val.Itemname] = val
		}
	}

	return appendNewIngredients, deleteOldIngredients
}

func compareIngredients(new, old *[]cake, appendNewIngredients, deleteOldIngredients *map[string]map[string]ingredients) {
	resAddOrRemovedIngredient := []string{}
	resChangedIngredient := []string{}
	resAddOrRemovedUnit := []string{}
	resChangedUnit := []string{}

	for _, val := range *new {
		resNew := (*appendNewIngredients)[val.Name]
		for _, v := range resNew {
			res, ok := (*deleteOldIngredients)[val.Name][v.Itemname]
			if !ok {
				resAddOrRemovedIngredient = append(resAddOrRemovedIngredient, fmt.Sprintf("ADDED ingredient \"%s\" for cake  \"%s\"\n", v.Itemname, val.Name))
				delete(*appendNewIngredients, v.Itemname)
			} else {
				if v.Itemcount != res.Itemcount {
					resChangedIngredient = append(resChangedIngredient, fmt.Sprintf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", v.Itemname, val.Name, v.Itemcount, res.Itemcount))
				}
				if res.Itemunit == "" && v.Itemunit != "" {
					resAddOrRemovedUnit = append(resAddOrRemovedUnit, fmt.Sprintf("ADDED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n", v.Itemunit, v.Itemname, val.Name))
				} else if v.Itemunit != res.Itemunit && v.Itemunit != "" && res.Itemunit != "" {
					resChangedUnit = append(resChangedUnit, fmt.Sprintf("CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", v.Itemname, val.Name, v.Itemunit, res.Itemunit))
				}

			}
		}
	}

	for _, val := range *old {
		resOld := (*deleteOldIngredients)[val.Name]
		for _, v := range resOld {
			res, ok := (*appendNewIngredients)[val.Name][v.Itemname]
			if !ok {
				resAddOrRemovedIngredient = append(resAddOrRemovedIngredient, fmt.Sprintf("REMOVED ingredient \"%s\" for cake  \"%s\"\n", v.Itemname, val.Name))
				delete((*appendNewIngredients), v.Itemname)
			} else {
				if res.Itemunit == "" && v.Itemunit != "" {
					resAddOrRemovedUnit = append(resAddOrRemovedUnit, fmt.Sprintf("REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n", v.Itemunit, v.Itemname, val.Name))
				}
			}
		}
	}

	printIngredientStatus(resAddOrRemovedIngredient, resChangedIngredient, resChangedUnit, resAddOrRemovedUnit)
}

func printIngredientStatus(resAddOrRemovedIngredient, resChangedIngredient, resChangedUnit, resAddOrRemovedUnit []string) {
	for _, v := range resAddOrRemovedIngredient {
		fmt.Print(v)
	}

	for _, v := range resChangedIngredient {
		fmt.Print(v)
	}

	for _, v := range resChangedUnit {
		fmt.Print(v)
	}

	for _, v := range resAddOrRemovedUnit {
		fmt.Print(v)
	}
}
