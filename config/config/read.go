package config

import (
	"../configJson"
	"log"
	"os"
)

func GetConfig(configPath string) map[string]string {
	argument := map[string]string{}
	settingDecodeJSON, err := configJSON.SettingDecodeJSON(configPath)
	if err != nil {
		log.Println(err)
		return argument
	}
	argument["pidFile"] = settingDecodeJSON.PidFile
	argument["cidrFile"] = settingDecodeJSON.BypassFile
	argument["logFile"] = os.DevNull
	argument["pythonPath"] = settingDecodeJSON.PythonPath
	argument["httpProxy"] = settingDecodeJSON.HttpProxyAddressAndPort
	argument["dnsServer"] = settingDecodeJSON.DnsServer

	// if argument["Workers"] == "" {
	// 	argument["Workers"] = "--workers " + "1 "
	// }

	argument["ssrPath"] = settingDecodeJSON.SsrPath
	argument["localAddress"] = settingDecodeJSON.LocalAddress
	argument["localPort"] = settingDecodeJSON.LocalPort

	return argument
}

//func argumentMatch(argument map[string]string, configTemp2 []string) {
//	switch configTemp2[0] {
//	case "python_path":
//		argument["pythonPath"] = configTemp2[1]
//	case "-python_path":
//		argument["pythonPath"] = ""
//	case "ssr_path":
//		argument["ssrPath"] = configTemp2[1]
//	case "-ssr_path":
//		argument["ssrPath"] = ""
//	case "config_path":
//		argument["configPath"] = configTemp2[1]
//	case "connect-verbose-info":
//		argument["connectVerboseInfo"] = "--connect-verbose-info"
//	case "workers":
//		argument["workers"] = configTemp2[1]
//	case "fast-open":
//		argument["fastOpen"] = "fast-open"
//	case "pid-file":
//		argument["pidFile"] = configTemp2[1]
//	case "-pid-file":
//		argument["pidFile"] = ""
//	case "log-file":
//		argument["logFile"] = configTemp2[1]
//	case "-log-file":
//		argument["logFile"] = ""
//	case "local_address":
//		argument["localAddress"] = configTemp2[1]
//	case "local_port":
//		argument["localPort"] = configTemp2[1]
//	case "acl":
//		argument["acl"] = configTemp2[1]
//	case "timeout":
//		argument["timeout"] = configTemp2[1]
//	case "httpProxy":
//		argument["httpProxy"] = configTemp2[1]
//	case "cidrFile":
//		argument["cidrFile"] = configTemp2[1]
//	case "dnsServer":
//		argument["dnsServer"] = configTemp2[1]
//		// case "daemon":
//		// 	argument["daemon"] = "-d start"
//	}
//}