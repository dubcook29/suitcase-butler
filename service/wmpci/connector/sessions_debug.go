package connector

import (
	"context"
	"fmt"

	"github.com/suitcase/butler/wmpci"
)

func WMPCISessionCompleteTesting(connector WMPConnector, conn map[string]wmpci.WMPCustom) {

	session := NewWMPCISession("./.wmpci")

	if data, ok, err := session.Connect(context.TODO(),
		connector,
		conn); err != nil {
		fmt.Println("wmp application connection errors:")
		panic(err)
	} else if ok {
		fmt.Println("wmp application connection successfully:", string(data))
	} else {
		fmt.Println("wmp application connection failed")
	}

	if response, err := session.Service(
		context.TODO(),
		session.WMPRequest(),
	); err != nil {
		panic(err)
	} else {
		fmt.Println("wmp application service request successfully")
		fmt.Printf("%+v\n", response)
	}

	if ok, err := session.Config(
		context.TODO(),
		map[string]wmpci.WMPCustom{
			"test": {Value: "Hello, this is a test"},
		}); err != nil {
		panic(err)
	} else {
		fmt.Println("wmp application config refresh:", ok)
	}

	fmt.Printf("====== prepare for testing <%s> fucntion interface\n", "connection")
	if msg, ok, err := session.Connect(context.TODO(), connector, conn); err != nil {
		panic(err)
	} else {
		fmt.Println(string(msg), ok, "\nwmp application registration, the registration information is as follows:")
		session.Registration.PrintAllWmpregistration()
	}
	fmt.Printf("------ test end <%s> fucntion interface\n", "connection")

}
