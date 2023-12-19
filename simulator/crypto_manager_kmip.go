package simulator

import (
	"fmt"
	//"log"
	//"encoding/json"

	kms "github.com/smira/go-kmip"
	"github.com/vmware/govmomi/simulator/vpx"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type CryptoManagerKmip struct {
	mo.CryptoManagerKmip
}

func (m *CryptoManagerKmip) init(r *Registry) {
	m.Enabled = true
	m.KmipServers = []types.KmipClusterInfo{
		{
			ClusterId:    types.KeyProviderId{Id: "kmipcluster"},
			Servers:      []types.KmipServerInfo{{Name: "kmipserver", Address: "localhost", Port: 5696}},
			UseAsDefault: true,
		},
	}

	root := r.content().CryptoManager // take the cryptomanager details from servicefolder
	m.CryptoManagerKmip.Self = *root
	//	m.CryptoManagerKmip =

}

func (m *CryptoManagerKmip) IsKmsClusterActive(clusterId string) bool {
	i := 0
	for i = 0; i < len(m.CryptoManagerKmip.KmipServers); i++ {
		if m.CryptoManagerKmip.KmipServers[i].ClusterId.Id == clusterId {
			return true
		}
	}

	return false
}

// Add create/get key methods here and expose them as API so that we can call these APIs from BVTs
func CreateKey() string {
	model := VPX()

	_ = New(NewServiceInstance(SpoofContext(), model.ServiceContent, model.RootFolder)) // 2nd pass panics w/o copying RoleList

	kmip := Map.Get(*vpx.ServiceContent.CryptoManager).(*CryptoManagerKmip)
	ans := kmip.IsKmsClusterActive("kmipcluster")
	fmt.Println("Is Kms CLuster Active : ", ans)

	cl, err := initClient()
	if err != nil {
		fmt.Println("Error in initializing Client :", err)
		panic(err)
	}
	err = cl.kclient.Connect()
	if err != nil {
		fmt.Println("Error in connecting Client :", err)
		panic(err)
	}
	fmt.Println("Client connected!")

	var resp interface{}
	for i := 0; i < 3; i++ {
		resp, err = cl.CreateKey()
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Println("Error in creating key (tried 3 times): ", err)
		panic(err)
	}
	fmt.Println("CreateKey: resp: , error: ", resp, err)
	response := resp.(kms.CreateResponse)

	return response.UniqueIdentifier
}

func GetKey(id string) error {

	cl, err := initClient()
	if err != nil {
		fmt.Println("Error in initializing Client :", err)
		return err
		//panic(err)
	}
	err = cl.kclient.Connect()
	if err != nil {
		fmt.Println("Error in connecting Client :", err)
		return err
		//panic(err)
	}
	fmt.Println("Client connected!")

	var resp interface{}
	for i := 0; i < 3; i++ {
		resp, err = cl.GetKey(id)
		if err == nil {
			break
		}
	}

	if err != nil {
		fmt.Println("Error in Getting key (tried 3 times): ", err)
		return err
		//panic(err)
	}
	fmt.Println("GetKey: resp: , error: ", resp, err)
	return nil
}

// create key and get key funcs here
