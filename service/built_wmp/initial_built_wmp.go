package built_wmp

import (
	"fmt"

	"github.com/suitcase/butler/built_wmp/wmp_cdn"
	"github.com/suitcase/butler/built_wmp/wmp_dns"
	"github.com/suitcase/butler/built_wmp/wmp_subdomain"
	"github.com/suitcase/butler/built_wmp/wmp_whois"
	"github.com/suitcase/butler/wmpci"
	connect_builtin "github.com/suitcase/butler/wmpci/connector/connect/built-in"
	wmpsessionmanager "github.com/suitcase/butler/wmpci/manager"
)

/*
 *                            _
 *  __      ____ _ _ __ _ __ (_)_ __   __ _
 *  \ \ /\ / / _` | '__| '_ \| | '_ \ / _` |
 *   \ V  V / (_| | |  | | | | | | | | (_| |
 *    \_/\_/ \__,_|_|  |_| |_|_|_| |_|\__, |
 *                                    |___/
 *  =============== WARNING ===============
 *  Only provides internal package testing during development,
 *	Not recommended for use in production environments.
 *  =============== WARNING ===============
 *
 */

func init() {
	fmt.Println("\x1b[31m\x1b[47m")
	fmt.Println("+-------------------------------------------------------------------+")
	fmt.Println("| *                                      _                        * |")
	fmt.Println("| *            __      ____ _ _ __ _ __ (_)_ __   __ _            * |")
	fmt.Println("| *            \\ \\ /\\ / / _` | '__| '_ \\| | '_ \\ / _` |           * |")
	fmt.Println("| *             \\ V  V / (_| | |  | | | | | | | | (_| |           * |")
	fmt.Println("| *              \\_/\\_/ \\__,_|_|  |_| |_|_|_| |_|\\__, |           * |")
	fmt.Println("| *                                               |___/           * |")
	fmt.Println("| *  ========================= WARNING =========================  * |")
	fmt.Println("| *   Only provides internal package testing during development   * |")
	fmt.Println("| *   Not recommended for use in production environments.         * |")
	fmt.Println("| *  ========================= WARNING =========================  * |")
	fmt.Println("+-------------------------------------------------------------------+")
	fmt.Println("\x1b[0m")
}

var applications = []connect_builtin.BuiltinServer{
	new(wmp_cdn.WMPCDNCheck),
	new(wmp_dns.WMPDnsResolve),
	new(wmp_subdomain.WMPSubdomain),
	new(wmp_whois.WMPWhois),
}

func InitialBuiltinApplications(session_manager *wmpsessionmanager.WMPSessionManager) *wmpsessionmanager.WMPSessionManager {
	for _, app := range applications {
		if _, err := session_manager.SelectConnectorConnectionCumstom("built-in"); err != nil {
			panic(err)
		} else {
			fmt.Println("`built-in` connector custom query ok")
			custom := map[string]wmpci.WMPCustom{
				"WMPCI": {
					Value: app,
				},
			}
			if welcome, err := session_manager.ConnectorConnectionSession("built-in", custom); err != nil {
				fmt.Println("`built-in` connector session created failed")
				panic(err)
			} else {
				fmt.Println("`built-in` connector session created ok,", welcome)
			}
		}
	}

	return session_manager
}
