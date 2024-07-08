package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getNetworkInterfaces() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	var result string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			result += fmt.Sprintf("%s: %s\n", iface.Name, addr.String())
		}
	}
	return result
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func IntToString(orig []int8) string {
	ret := make([]byte, len(orig))
	size := -1
	for i, o := range orig {
		if o == 0 {
			size = i
			break
		}
		ret[i] = byte(o)
	}
	if size == -1 {
		size = len(orig)
	}

	return string(ret[0:size])
}

func getUsers() string {

	hostname, _ := os.Hostname()
	result := fmt.Sprintf("%s\n", hostname)
	h, _ := host.Info()

	if h.KernelArch == "aarch64" {
		file, _ := os.Open("/var/run/utmp")
		defer file.Close()

		stat, _ := file.Stat()

		buf := make([]byte, stat.Size())
		file.Read(buf)
		count := len(buf) / sizeOfUtmp

		for i := 0; i < count; i++ {
			b := buf[i*sizeOfUtmp : (i+1)*sizeOfUtmp]

			var u utmp
			br := bytes.NewReader(b)
			err := binary.Read(br, binary.LittleEndian, &u)
			if err != nil {
				continue
			}
			if u.Type != 7 {
				continue
			}

			result += fmt.Sprintf("%s@%s %s\n", IntToString(u.User[:]), IntToString(u.Host[:]), IntToString(u.Line[:]))
		}

	} else {
		users, _ := host.Users()
		for _, user := range users {
			// data, _ := json.MarshalIndent(user, "", " ")
			// result += fmt.Sprintf("%s@%s %s %s\n", user.User, user.Host, user.Terminal, time.Unix(int64(user.Started), 0).Format("0102 15:04:05"))
			result += fmt.Sprintf("%s@%s %s %s\n", user.User, user.Host, user.Terminal)
		}
	}

	return result
}

func getCpuInfo() string {

	percent, _ := cpu.Percent(3*time.Millisecond, false)

	h, _ := host.Info()

	infos, _ := cpu.Info()
	u := ""
	for _, info := range infos {
		if len(info.ModelName) > 0 {
			u = fmt.Sprintf("%s", info.ModelName)
		} else {
			u = fmt.Sprintf("%s", h.KernelArch)
		}
		break
	}

	physicalCount, _ := cpu.Counts(false)
	logicalCount, _ := cpu.Counts(true)
	sensors, _ := host.SensorsTemperatures()

	// raspberry pi cpu Temperature
	t := ""
	for _, sensor := range sensors {
		if sensor.SensorKey == "cpu_thermal_input" {
			t = fmt.Sprintf("%.2f°C ", sensor.Temperature)
		}
	}

	return fmt.Sprintf("name: %s %s\nos  : %s-%s\ncpu : %s\ncore: %d/%d (%s%.2f%%)",
		h.Platform, h.PlatformVersion, h.OS, h.KernelVersion, u, physicalCount, logicalCount, t, percent[0])
}

func getMemory() string {

	v, _ := mem.VirtualMemory()

	return fmt.Sprintf("%d/%d MB (%.2f%%)", v.Active/1024/1024, v.Total/1024/1024, v.UsedPercent)
}

func getBootTime() string {

	info, err := host.Info()
	if err != nil {
		fmt.Println("Failed to get host info:", err)
		return ""
	}

	return time.Unix(int64(info.BootTime), 0).Format("2006-01-02 15:04:05")
}

func getSysInfo() string {

	memTotal := getMemory()

	return fmt.Sprintf("boot: %s\n%s\nram : %s\n",
		getBootTime(), getCpuInfo(), memTotal)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create a paragraph for the network interfaces
	pNetwork := widgets.NewParagraph()
	pNetwork.Title = "Network Interfaces"
	pNetwork.BorderStyle.Fg = ui.ColorRed
	pNetwork.TitleStyle.Fg = ui.ColorRed
	pNetwork.TextStyle.Fg = ui.ColorRed
	pNetwork.Text = getNetworkInterfaces()

	// Create a paragraph for the clock
	pClock := widgets.NewParagraph()
	pClock.Title = "System Info"
	pClock.BorderStyle.Fg = ui.ColorYellow
	pClock.TitleStyle.Fg = ui.ColorYellow
	pClock.TextStyle.Fg = ui.ColorYellow
	pClock.Text = fmt.Sprintf("date: %s\n%s", getCurrentTime(), getSysInfo())

	// Create a paragraph for the Users
	pUsers := widgets.NewParagraph()
	pUsers.Title = "Users"
	pUsers.BorderStyle.Fg = ui.ColorGreen
	pUsers.TitleStyle.Fg = ui.ColorGreen
	pUsers.TextStyle.Fg = ui.ColorGreen
	pUsers.Text = getUsers()

	// Function to update the clock
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-ticker:
				pClock.Text = fmt.Sprintf("date: %s\n%s", getCurrentTime(), getSysInfo())
				ui.Render(pClock)
				pUsers.Text = getUsers()
				ui.Render(pUsers)
			// pNetwork.Text = getNetworkInterfaces()
			// ui.Render(pNetwork)
			case e := <-uiEvents:
				if e.Type == ui.KeyboardEvent {
					return
				}
			}
		}
	}()

	// Initial rendering of the widgets
	width, height := ui.TerminalDimensions()
	pClock.SetRect(0, 0, width*1/2, height/2)
	pUsers.SetRect(width*1/2, 0, width, height/2)
	pNetwork.SetRect(0, height/2, width, height)

	ui.Render(pNetwork, pClock, pUsers)

	// Block main goroutine until a keyboard event is received
	<-uiEvents
}
