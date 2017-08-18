package test

import (
	"fmt"

	xguid "github.com/beevik/guid"
)

func guid_foo() {
	guid := xguid.NewString()
	fmt.Println(guid)
}
