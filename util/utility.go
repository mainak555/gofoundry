package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	YYYYMMDDhhmmss = "2006-01-02 15:04:05"
	YYYYMMDD       = "2006-01-02"
	MAX_PAGE_SIZE  = 25000
)

func GetEntityName[T any]() (string, error) {
	var entity T
	if reflect.TypeOf(entity).Kind() != reflect.Struct {
		return "", errors.New("entity type must a struct")
	}
	return reflect.TypeOf(entity).Name(), nil
}

func ToJson[T interface{}](obj *T) string {
	json, err := json.Marshal(*obj)
	if err != nil {
		return `{"err": "Unsupported value type"}`
	}
	return string(json)
}

func Unmarshal[T any](data any) (*T, error) {
	var tmp T
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(bytes, &tmp); err != nil {
		return nil, err
	}
	return &tmp, nil
}

func ToTimePtr(t time.Time) *time.Time {
	return &t
}

func MapSlice[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

func PanicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func StringArrayToObjectIDArray(ids ...string) ([]primitive.ObjectID, error) {
	var arrayOfIds []primitive.ObjectID
	for _, num := range ids {
		objID, err := primitive.ObjectIDFromHex(num)
		if err != nil {
			return nil, err
		} else if !objID.IsZero() {
			arrayOfIds = append(arrayOfIds, objID)
		}
	}
	return arrayOfIds, nil
}

func ObjectIDArrayToStringArray(ids ...primitive.ObjectID) []string {
	var arrayOfIds []string
	for _, num := range ids {
		objID := num.Hex()
		arrayOfIds = append(arrayOfIds, objID)
	}
	return arrayOfIds
}

func FindDuplicates(arr []string) []string {
	seen := make(map[string]bool)
	duplicates := []string{}

	for _, num := range arr {
		if seen[num] {
			duplicates = append(duplicates, num)
		}
		seen[num] = true
	}
	return duplicates
}

func StringDateArrayToTimeDateArray(format string, dates ...string) ([]time.Time, error) {
	var arrayOfDates []time.Time
	if format == "" {
		format = YYYYMMDD
	}
	for _, date := range dates {
		newDate, err := time.Parse(format, date)
		if err != nil {
			return nil, err
		}
		arrayOfDates = append(arrayOfDates, newDate)
	}
	return arrayOfDates, nil
}

func StringToDate(date *string, format string) (*time.Time, error) {
	if date != nil && *date != "" {
		if date, err := StringDateArrayToTimeDateArray(format, *date); err == nil {
			return &date[0], nil
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func StructMap[T any](callBack func(key string) string, entity *T) error {
	if reflect.TypeOf(*entity).Kind() != reflect.Struct {
		return errors.New("unknown struct type")
	}

	dType := reflect.TypeOf(entity)
	for i := 0; i < dType.Elem().NumField(); i++ {
		filed := dType.Elem().Field(i)
		var value string
		if filed.Tag != "" && filed.Tag.Get("json") != "" {
			value = callBack(filed.Tag.Get("json"))
		} else {
			value = callBack(filed.Name)
		}
		if value == "" {
			continue
		}
		switch filed.Type.Kind() {
		case reflect.Bool:
			if v, e := strconv.ParseBool(value); e == nil {
				reflect.ValueOf(entity).Elem().Field(i).Set(reflect.ValueOf(v))
			}
		case reflect.Int64:
			if v, e := strconv.ParseInt(value, 0, 64); e == nil {
				reflect.ValueOf(entity).Elem().Field(i).Set(reflect.ValueOf(v))
			}
		case reflect.Float64:
			if v, e := strconv.ParseFloat(value, 64); e != nil {
				reflect.ValueOf(entity).Elem().Field(i).Set(reflect.ValueOf(v))
			}
		case reflect.String:
			reflect.ValueOf(entity).Elem().Field(i).Set(reflect.ValueOf(value))
		}
	}
	return nil
}

func ChunkBySplitSize[T any](array []T, size int) (chunks [][]T) {
	for size < len(array) {
		array, chunks = array[size:], append(chunks, array[0:size:size])
	}
	return append(chunks, array)
}

func ChunkBySplitNumber[T any](array []T, number int) (chunks [][]T) {
	length := len(array)
	size := (length + number - 1) / number
	for i := 0; i < length; i++ {
		end := i + size
		if end > length {
			end = length
		}
		chunks = append(chunks, array[i:end])
	}
	return chunks
}

func MD5Hash(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func AESEncryptBase64(value []byte, passPhrase string) (base64String string, err error) {
	aesBlock, err := aes.NewCipher([]byte(MD5Hash(passPhrase)))
	if err != nil {
		return "", err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipheredByte := gcmInstance.Seal(nonce, nonce, value, nil)
	return base64.StdEncoding.EncodeToString(cipheredByte), nil
}

func AESDecryptBase64(base64String, passPhrase string) (value []byte, err error) {
	ciphered, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	aesBlock, err := aes.NewCipher([]byte(MD5Hash(passPhrase)))
	if err != nil {
		return nil, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	return gcmInstance.Open(nil, nonce, cipheredText, nil)
}

func RandomString(length int, chars ...rune) (*string, error) {
	var table []byte
	if len(chars) > 0 {
		for _, c := range chars {
			table = append(table, byte(c))
		}
	} else {
		table = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return nil, err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	str := string(b)
	return &str, nil
}

func RandomNumbers(count, length int) ([]uint64, error) {
	var result []uint64
	for i := 0; i < count; i++ {
		str, err := RandomString(length)
		if err != nil {
			return result, err
		}
		u64, err := strconv.ParseUint(*str, 10, 64)
		if err != nil {
			return result, err
		}
		result = append(result, u64)
	}
	return result, nil
}

func JsonEntityMapper[T any](resp any) (*T, error) {
	jsonByte, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var data T
	if err := json.Unmarshal(jsonByte, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func ObjectIDNotEqual(id1 *primitive.ObjectID, id2 *string) bool { //returns true if not matches
	return (id1 != nil && id2 != nil && *id2 != "" && id1.Hex() != *id2) ||
		((id1 == nil || id1.IsZero()) && id2 != nil && *id2 != "") ||
		(id1 != nil && !id1.IsZero() && (id2 == nil || *id2 == ""))
}
