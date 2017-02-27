package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var containerIndexStr = "_containerindex"    //This will be used as key and a value will be an array of Container IDs	

var openOrdersStr = "_openorders"	  // This will be the key, value will be a list of orders(technically - array of order structs)

//var coinIndexStr = "_coinindex"          //key - "_coinindex", value - an array of Coin Ids, no matter whom they are tagged to

type MilkContainer struct{

        ContainerID string `json:"containerid"`
        User string        `json:"user"`

        Litres int       `json:"litres"`

}


type SupplyCoin struct{

        CoinID string `json:"coinid"`
        User string        `json:"user"`
}

type Order struct{
        OrderID string `json:"orderid"`
       User string `json:"user"`
       Status string `json:"status"`
       Litres int   `json:"litres"`
}

type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}


type Asset struct{
	  User string        `json:"user"`
	containerIDs []string `json:"containerids"`
	LitresofMilk int `json:"litresofmilk"`
	coinIDs []string `json:"coinids"`
	Supplycoins int `json:"supplycoins"`
}




func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	
	var err error
	
	fmt.Println("Welcome to Supplychain management Phase 1, Deployment is on the go")
 
       if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
       }

       err = stub.PutState("hello world",[]byte(args[0]))  //Just to check the network whether we can read and write
       if err != nil {
		return nil, err
       }
	
        /* Reset container index list - Making sure the value corresponding to containerIndexStr  is empty */

       var empty []string
       jsonAsBytes, _ := json.Marshal(empty)                                   //create an empty array of string
       err = stub.PutState(containerIndexStr, jsonAsBytes)                     //Resetting - Making milk container list as empty 
       if err != nil {
		return nil, err
        }  
	
	
	/* Resetting the order list - Making sure the value corresponding to openOrdersStr is empty */
       var orders AllOrders                                            // new instance of Orderlist 
	jsonAsBytes, _ = json.Marshal(orders)				//  it will be null initially
	err = stub.PutState(openOrdersStr, jsonAsBytes)                 //So the value for key is null
	if err != nil {       
		return nil, err
}
	// Resetting the Assets of Supplier,Market, Logistiscs.
	
	var emptyasset Asset
	
	
	jsonAsBytes, _ = json.Marshal(emptyasset)                // this is the byte format format of empty Asset structure
	err = stub.PutState("SupplierAssets",jsonAsBytes)        // key -Supplier assets and value is empty now --> Supplier has no assets
	err = stub.PutState("MarketAssets", jsonAsBytes)         // key -Market assets and value is empty now --> Market has no assets
	err = stub.PutState("LogisticsAssets", jsonAsBytes)      // key - Logistics assets and value is empty now --> Logistic has no assets
	
	fmt.Println("Successfully deployed the code and orders and assets are reset")
	
	 return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		
		return t.Create_milkcontainer(stub, args)
		
	}else if function == "Create_coin" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Create_coin(stub, args)	
        }else if function == "Buy_milk" {		         //creates a coin - invoked by market /logistics - params - coin id, entity name
		return t.Buy_milk(stub, args)	
        }/*else if function == "Order_milk"{                      // To order something - invoked by market - params - litres
		res,err :=  t.Order_milk(stub,args)
		//jsonresp,_ := View_order(stub,args)
		//fmt.Println(jsonresp)
		
		//jsonresp,_ := init_logistics(stub,args)
		//fmt.Println(jsonresp)
		//jsonresp,_ = set_user(stub,args)
		
		//jsonresp,_ = checktheproduct(stub,args)
		
		//jsonresp,_ = cointransfer(stub,args)
		
               
		//printdetails(stub, 3)
		return res,err
	}*/
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}






func (t *SimpleChaincode) Create_milkcontainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error

// "1x22" "supplier" 20 
// args[0] args[1] args[2] 
	fmt.Println("Creating milkcontainer asset")
id := args[0]
user := args[1]
litres,_ :=strconv.Atoi(args[2])
	
// Checking if the container already exists in the network
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details og given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

//If not present, create it and Update ledger, containerIndexStr, Assets of Supplier
//Creation
res.ContainerID = id
res.User = user
res.Litres = litres
milkAsBytes, _ =json.Marshal(res)

stub.PutState(res.ContainerID,milkAsBytes)
	
	 fmt.Println("Container created successfully, details are")
	fmt.Printf("%+v\n", res)
	
//Update containerIndexStr	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get container index")
	}
	var containerIndex []string                                        //an array to store container indices - later this wil be the value for containerIndexStr
	json.Unmarshal(containerAsBytes, &containerIndex)	
	
	
	containerIndex = append(containerIndex, res.ContainerID)          //append the newly created container to the global container list									//add marble name to index list
	fmt.Println("! container index: ", containerIndex)
	jsonAsBytes, _ := json.Marshal(containerIndex)
        err = stub.PutState(containerIndexStr, jsonAsBytes)
	
// append the container ID to the existing assets of the Supplier
	
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
	supplierasset.User = "Supplier"
	supplierasset.containerIDs = append(supplierasset.containerIDs, res.ContainerID)
	supplierasset.LitresofMilk += res.Litres
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
	fmt.Println("Balance of Supplier")
       fmt.Printf("%+v\n", supplierasset)

	return nil,nil

}





func (t *SimpleChaincode) Create_coin(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"1x245" "Market/Logistics"
id := args[0]
user:= args[1]
//Check if coin already exists in network
	fmt.Println(" Inside Create coin function")
coinAsBytes , err := stub.GetState(id)
if err != nil{
              return nil, errors.New("Failed to get details of given id")
} 

res :=SupplyCoin{}

json.Unmarshal(coinAsBytes, &res)

if res.CoinID == id{

          fmt.Println("Coin already exists")
          fmt.Println(res)
          return nil,errors.New("This coin already exists")
}
// Proceed to create if not der in ntwrk
res.CoinID = id
res.User = user

coinAsBytes, _ = json.Marshal(res)
stub.PutState(id,coinAsBytes)
	
	 fmt.Println("Coin created successfully, details are")
	fmt.Printf("%+v\n", res)
//t.read(stub,"res.CoinID")
	
	// append the container ID to the existing assets of the Supplier
	
	userAssets := user +"Assets"
	fmt.Println("Updating ",userAssets)
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
	asset.User = user
	asset.coinIDs = append(asset.coinIDs, res.CoinID)
	asset.Supplycoins += 1
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	fmt.Println("Balance of " , user)
       fmt.Printf("%+v\n", asset)


return nil,nil
}


/*********************Buy milk - Customer interactio*******************/

func (t *SimpleChaincode) Buy_milk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//args[0]
//"10"
// customer asks for a qty, check if market has that much quantity, if there-create a container for customer with qty he asked, and subtract the same from Market
	fmt.Println("Hi, welcome to Supply chain, we are inside Buy_milk")
	quantity,_ := strconv.Atoi(args[0])
	marketassetAsBytes, err := stub.GetState("MarketAssets")
	if err != nil {
		return nil, errors.New("Failed to get market asset")
	}
	Marketasset := Asset{}             
	json.Unmarshal(marketassetAsBytes, &Marketasset )
	Threshold := 15
	if (Marketasset.LitresofMilk - quantity <= Threshold){
		a,b := Order_milk(stub,"20")
		fmt.Println(a,b)
	}
	
	return nil,nil
}




func Order_milk(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
//"20"
//litres
var err error
Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "pending"
Openorder.OrderID = "abcd"
Openorder.Litres,err = strconv.Atoi(args)
orderAsBytes,_ := json.Marshal(Openorder)
	
err = stub.PutState(Openorder.OrderID,orderAsBytes)
	
	 fmt.Println("your Order has been generated successfully")
	fmt.Printf("%+v\n", Openorder)
	
if err != nil {
		return nil, err
}

//Add the new order to the orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)         // note this is ordersAsBytes - plural, above one is orderAsBytes-Singular
	if err != nil {
		return nil, errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)				
	
	orders.OpenOrders = append(orders.OpenOrders , Openorder);		//append the new order - Openorder
	fmt.Println("! appended %q to existing orders", Openorder.OrderID)
	jsonAsBytes, _ := json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)		  // Update the value of the key openOrdersStr
	if err != nil {
		return nil, err
}
	//t.read(stub,"openOrdersStr")
	 View_order(stub)
	//fmt.Println("After View order invocation from order milk")
	//fmt.Println("Data returned from View order")
	//fmt.Println(a)
	//fmt.Println("errors from view order if any")
	//fmt.Println(b)
return nil,nil
}

/*
func printdetails(stub  shim.ChaincodeStubInterface, a int)(err error) {
	
	fmt.Println(a)
	err = stub.PutState("hello cts",[]byte("inside printdetails")) //Just to check the network whether we can read and write
       if err != nil {
		return  err
       }
	return err
}
*/
func  View_order(stub shim.ChaincodeStubInterface) ( error) {
	// This will be invoked by Supplier- think of UI-View orders- does he pass any parameter there...
	// so here also no need to pass any arguments. args will be empty-but just for syntax-pass something as parameter
       // a := args[0]
//	fmt.Println(a)
	
	/* fetching the Orders*/
	fmt.Printf("Inside View order")
	
	
	ordersAsBytes, _ := stub.GetState(openOrdersStr)
	
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	
	
/*fetching the containers*/	
	
	containerAsBytes, err := stub.GetState(containerIndexStr)
	if err != nil {
		return errors.New("Failed to get container index")
	}
	var containerIndex []string             //an array to clone container indices
	json.Unmarshal(containerAsBytes, &containerIndex)
	
// From the list of Id's , picking up one Id and fetching its details
	
	containerAsBytes,_ = stub.GetState(containerIndex[0])
	
	res := MilkContainer{} 
        json.Unmarshal(containerAsBytes, &res)

// If ordered quantity and container quantity , then proceed and trigger logistics(How ever this is not automated here,we will do it 
	
	if (res.Litres == orders.OpenOrders[0].Litres) {
		fmt.Println("Found a suitable container and about to ship")
		
		orders.OpenOrders[0].Status = "Ready to be Shipped"
		//t.init_logistics(stub,orders.OpenOrders[0].OrderId, containerIndex[0])
		ordersAsBytes,_ = json.Marshal(orders)
		fmt.Printf("%+v\n", orders.OpenOrders[0])
		stub.PutState(openOrdersStr,ordersAsBytes)
		
		OrderID := orders.OpenOrders[0].OrderID
		orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ShipOrder.Status = "Ready to be Shipped"
	 
	orderAsBytes,err = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,orderAsBytes)
		
	a := []string{ShipOrder.OrderID,res.ContainerID}
		init_logistics(stub,a)
	return nil	
		
		//t.read(stub,openOrdersStr)
	}else{
                stub.PutState("sorry",[]byte("we couldn't find a product for your choice of requirements"))
		return nil
        }


}

func init_logistics(stub shim.ChaincodeStubInterface, args []string) ( error) {
	
	
	
	//args[0] args[1]
	// OrderId, ContainerID
	fmt.Println("we are moving, item is on Go, inside init_logistics")
	fmt.Println("Inside Init logistics function")
	OrderID := args[0]
	//ContainerID := args[1]
	
	// fetch the order details and update status as "in transit"
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
	
	ShipOrder.Status = "In transit"
	 
	orderAsBytes,err = json.Marshal(ShipOrder)
	
	stub.PutState(OrderID,orderAsBytes)
	
	fmt.Printf("%+v\n", ShipOrder)
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
	
	orders.OpenOrders[0].Status = "In transit"
	ordersAsBytes,_ = json.Marshal(orders)
	stub.PutState(openOrdersStr,ordersAsBytes)
	set_user(stub,args)
	
	//t.read(stub, OrderID)
	//t.set_user(stub,OrderID,ContainerID)
	
	
	
return nil
}

func  set_user(stub shim.ChaincodeStubInterface, args []string) ( error) {
	
// OrderId  ContainerID
//args[0] args[1]
	
//So here we will set the user name in container ID to the one in Order ID and Status to Delivered - Asset Transfer
	fmt.Println("Transferring asset owner ship")
	OrderID := args[0]
	ContainerID := args[1]
//fetch order details
       orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details	
	assetAsBytes,err := stub.GetState(ContainerID)
	container := MilkContainer{}
	json.Unmarshal(assetAsBytes, &container)

	if (container.User == "Supplier"){
	
	container.User = ShipOrder.User             //ASSET TRANSFER
	fmt.Printf("%+v\n", container)
	fmt.Println("pushing the updated container back to ledger")
	assetAsBytes,err = json.Marshal(container)
	stub.PutState(ContainerID, assetAsBytes)    //Pushing the updated container  back to the ledger
	
	fmt.Println("Updating Supplier assets..")
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
		
	userAssets := container.User +"Assets"
	fmt.Println("Updating ",userAssets)
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
		
	asset.LitresofMilk += container.Litres
	supplierasset.LitresofMilk -= container.Litres
		
        supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
		
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	
	ShipOrder.Status = "Delivered"
	orderAsBytes,err = json.Marshal(ShipOrder) 
	
	stub.PutState(OrderID,orderAsBytes)      //pushing the updated Order back to ledger
		fmt.Printf("%+v\n", ShipOrder)
	//Updating the orders list 
		ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return errors.New("Failed to get openorders")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)	
		
		orders.OpenOrders[0].Status = ShipOrder.Status
		ordersAsBytes,_ = json.Marshal(orders)
		stub.PutState(openOrdersStr,ordersAsBytes)
		checktheproduct(stub,args)
		
	}else
        {
                stub.PutState("setuser",[]byte("failure in this function"))
                //t.read(stub,"setuser")
                return nil
        }


return nil
	//t.checktheproduct(stub,OrderID,ContainerID)
}


func  checktheproduct(stub shim.ChaincodeStubInterface, args []string) ( error) {

// args[0] args[1]
// OrderID, ContainerID
	fmt.Println("Let us check the product")
	OrderID := args[0]
	ContainerID := args[1]
//fetch order details
	orderAsBytes, err := stub.GetState(OrderID)
	if err != nil {
		return  errors.New("Failed to get openorders")
	}
	ShipOrder := Order{} 
	json.Unmarshal(orderAsBytes, &ShipOrder)
//fetch container details
       assetAsBytes,_ := stub.GetState(ContainerID)
	Deliveredcontainer := MilkContainer{}
	json.Unmarshal(assetAsBytes, &Deliveredcontainer)

//check and transfer coin
	if (Deliveredcontainer.User == "Market" && Deliveredcontainer.Litres == ShipOrder.Litres) {
		
		fmt.Println("Thanks, I got  the right product")
		stub.PutState("Market Response",[]byte("Product received"))
		var b [3]string
		b[0]= "1x245"
		b[1] = "Market"
		b[2] = "Supplier"
		cointransfer(stub,b)
		//t.cointransfer(stub,coinid) coinid -hard code it and send the coin id created by market
		return nil
       }else{
                stub.PutState("checktheproduct",[]byte("failure"))
		fmt.Println("I didn't get the right product")
               // t.read(stub,"checktheproduct")
                return nil
        }

	
return nil


}


func cointransfer( stub shim.ChaincodeStubInterface, args [3]string) ( error) {
	
//args[0] 
//coinID  
	//lets keep it simple for now, just fetch the coin from ledger, change username to Supplier and End of Story
	CoinID := args[0]
	sender := args[1]
	receiver := args[2]
	
	fmt.Println("Payment time, inside moneytransfer")
	
	assetAsBytes,err := stub.GetState(CoinID)
	if err != nil{
		fmt.Println("Something wrong happened while fetching coin details")
	}
	
	Transfercoin := SupplyCoin{}
	json.Unmarshal(assetAsBytes, &Transfercoin)
	
	if (Transfercoin.User == sender) {   // check if the market guy actually holds coin in his name
	
		Transfercoin.User = receiver
		assetAsBytes,err = json.Marshal(Transfercoin)
		stub.PutState(CoinID, assetAsBytes)
		fmt.Printf("%+v\n", Transfercoin)
	fmt.Println("Updating Supplier assets..")
	supplierassetAsBytes,_ := stub.GetState("SupplierAssets")        // The same key which we used in Init function 
	supplierasset := Asset{}
	json.Unmarshal( supplierassetAsBytes, &supplierasset)
		
	userAssets := container.User +"Assets"
	fmt.Println("Updating ",userAssets)
	assetAsBytes,_ := stub.GetState(userAssets)        // The same key which we used in Init function 
	asset := Asset{}
	json.Unmarshal( assetAsBytes, &asset)
		
		asset.Supplycoins  -= 1
		supplierasset.Supplycoins += 1
	
	supplierassetAsBytes,_=  json.Marshal(supplierasset)
	stub.PutState("SupplierAssets",supplierassetAsBytes)
		
	assetAsBytes,_=  json.Marshal(asset)
	stub.PutState(userAssets,assetAsBytes)
	
		
		fmt.Println("END OF PHASE 1")
		return nil
		
	}else{
	
		fmt.Println("There was some issue in transferring")
		//stub.PutState("cointransfer",[]byte("problem in coin transfer"))
		//t.read(stub,"cointransfer")
		return nil
	}

	
return nil
	
}



// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
