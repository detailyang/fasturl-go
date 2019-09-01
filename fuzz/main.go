// +build gofuzz

package url

import (
	"fmt"
	"reflect"

	"github.com/detailyang/fasturl-go/fasturl"
)

func FuzzURL(data []byte) int {
	var u1 fasturl.FastURL

	err := fasturl.Parse(&u1, data)
	if err != nil {
		return 0
	}

	d1 := u1.Encode(nil)

	var u2 fasturl.FastURL
	err = fasturl.Parse(&u2, d1)
	if err != nil {
		panic(err)
	}

	d2 := u2.Encode(nil)

	if !reflect.DeepEqual(d1, d2) {
		fmt.Printf("url1: %#v\n", string(d1))
		fmt.Printf("url2: %#v\n", string(d2))
		panic("fail")
	}
	return 1
}

func FuzzQuery(data []byte) int {
	var q1 fasturl.Query

	err := fasturl.ParseQuery(&q1, data)
	if err != nil {
		return 0
	}

	output := q1.Encode(nil)

	var q2 fasturl.Query
	err = fasturl.ParseQuery(&q2, output)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(q1, q2) {
		fmt.Printf("query1: %#v\n", q1)
		fmt.Printf("query2: %#v\n", q2)
		panic("fail")
	}
	return 1
}
