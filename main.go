package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func getNetworkInterfaces() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var result string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			result += fmt.Sprintf("%s: %s\n", iface.Name, addr.String())
		}
	}
	return result, nil
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getUsers() string {

	hostname, _ := os.Hostname()
	usersTxt := fmt.Sprintf("%s\n", hostname)
	users, _ := host.Users()

	for _, user := range users {
		// data, _ := json.MarshalIndent(user, "", " ")
		// usersTxt += string(data)

		usersTxt += fmt.Sprintf("%s@%s %s %s\n", user.User, user.Host, user.Terminal, time.Unix(int64(user.Started), 0).Format("2006-01-02 15:04:05"))
	}
	return usersTxt
}

func getCpuInfo() string {
	percent, err := cpu.Percent(3*time.Millisecond, false)
	// perPercents, _ := cpu.Percent(3*time.Millisecond, true)
	if err != nil {
		return ""
	}

	infos, _ := cpu.Info()

	cpu_txt := ""
	for _, info := range infos {
		cpu_txt = info.ModelName
		break
	}

	physicalCount, _ := cpu.Counts(false)
	logicalCount, _ := cpu.Counts(true)

	return fmt.Sprintf("syst: %s %s\ncpus: %s (%d/%d %.2f%%)",
		runtime.GOOS, runtime.GOARCH, cpu_txt, physicalCount, logicalCount, percent[0])
}

func getMemory() string {

	v, _ := mem.VirtualMemory()

	return fmt.Sprintf("%v Bytes (%f%%)", v.Total, v.UsedPercent)
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

	return fmt.Sprintf("boot: %s\n%s\nmemo: %s\n",
		getBootTime(), getCpuInfo(), memTotal)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	networkInfo, err := getNetworkInterfaces()
	if err != nil {
		log.Fatalf("failed to get network interfaces: %v", err)
	}

	// Create a paragraph for the network interfaces
	pNetwork := widgets.NewParagraph()
	pNetwork.Title = "Network Interfaces"
	pNetwork.BorderStyle.Fg = ui.ColorRed
	pNetwork.TitleStyle.Fg = ui.ColorRed
	pNetwork.TextStyle.Fg = ui.ColorRed
	pNetwork.Text = networkInfo

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
			case e := <-uiEvents:
				if e.Type == ui.KeyboardEvent {
					fmt.Println("__DEBUG_", ui.KeyboardEvent)
					return
				}
			}
		}
	}()

	// Initial rendering of the widgets
	width, height := ui.TerminalDimensions()
	pClock.SetRect(0, 0, width*2/3, height/2)
	pUsers.SetRect(width*2/3, 0, width, height/2)
	pNetwork.SetRect(0, height/2, width, height)

	ui.Render(pNetwork, pClock, pUsers)

	// Block main goroutine until a keyboard event is received
	<-uiEvents
}
