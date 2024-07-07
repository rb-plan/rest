package main

import (
	"encoding/json"
	"fmt"

	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

func TestMain(m *testing.T) {

}

func TestGetSysInfo(m *testing.T) {
	fmt.Println(getSysInfo())
	infos, _ := cpu.Info()
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", " ")
		fmt.Print(string(data))
	}
}

func TestGetUsers(t *testing.T) {
	// result := getUsers()
	// if len(result) == 0 {
	// 	t.Errorf("Sum(1, 2) expected 3, got %s", result)
	// }
	users, _ := host.Users()

	for _, user := range users {
		data, _ := json.MarshalIndent(user, "", " ")
		fmt.Println(data)
	}
}
