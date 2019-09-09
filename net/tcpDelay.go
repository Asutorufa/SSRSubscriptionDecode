package getdelay

import (
	"fmt"
	"log"
	"net"
	"time"

	configJSON "../config/configjson"
)

// TCPDelay get once delay by tcp
func TCPDelay(address, port string) (time.Duration, bool, error) {
	//fmt.Print("tcp connecting")
	timeNow := time.Now()
	conn, err := net.DialTimeout("tcp", address+":"+port, 3*time.Second)
	if err != nil {
		if time.Since(timeNow) > 3*time.Second {
			log.Println("tcp timeout,tcp connect time over 5s")
			return 999 * time.Hour, false, err
		}
		log.Println("tcp connect error")
		return 999 * time.Hour, false, err
	}
	defer conn.Close()
	delay := time.Since(timeNow)
	//fmt.Print(delay, " ")
	return delay, true, nil
}

func getTCPDelayAverage(server, serverPort string) time.Duration {
	var delay [3]time.Duration
	var err error
	for i := 0; i < 3; i++ {
		delay[i], _, err = TCPDelay(server, serverPort)
		if err != nil {
			// log.Println("tcp connect error")
			// log.Println(err)
			continue
		}
	}
	return (delay[0] + delay[1] + delay[2]) / 3
}

// GetTCPDelayJSON get delay by tcp
func GetTCPDelayJSON(configPath string) {
	for {
		node, err := configJSON.SelectNode(configPath)
		if err != nil {
			return
		}
		if node.Server == "" {
			break
		}
		fmt.Print(node.Remarks + "delay(3 times): ")
		fmt.Println("average:", getTCPDelayAverage(node.Server, node.ServerPort))
	}
}
