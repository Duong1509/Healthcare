package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}
type PermissionRW struct {
	// Key   string   `json:"Key"`
	PatientKey string   `json:"PatientKey"`
	CCCD       string   `json:"CCCD"`
	Read       []string `json:"Read"`
	Write      []string `json:"Write"`
}
type MedicalBackground struct {
	// IDMB            string `json:IDMB`
	PatientKey  string `json:"PatientKey"`
	CCCD        string `json:"CCCD"`
	Information string `json:"Information"`
	// Permission  PermissionRW `json:"Permission"`
}
type MedicalRecord struct {
	// IDRecord string `json:"IDRecord"`
	// IDPermission string `json:"IDPermission"`
	MRKey    string `json:"MRKey"`
	CCCD     string `json:"CCCD"`
	SickType string `json:"SickType"`
	Data     string `json:"Data"`
	Time     string `json:"Time"`
}
type MedicineHistory struct {
	CCCD string `json:"CCCD"`
	Data string `json:"Data"`
	Time string `json:"Time"`
}

func (c *Contract) Init(ctx contractapi.TransactionContextInterface, params string) {

}

func CheckUser(ctx contractapi.TransactionContextInterface) (string, error) {
	_, errUserId := ctx.GetClientIdentity().GetID()
	if errUserId != nil {
		return "", fmt.Errorf("Lỗi: %s", errUserId)
	}
	userRole, errUserRole := ctx.GetClientIdentity().GetMSPID()
	if errUserRole != nil {
		return "", fmt.Errorf("Lỗi: %s", errUserRole)
	}
	if userRole != "issuer" {
		return "", fmt.Errorf("Nguoi dung khong co quyen tao ban ghi")
	}
	permis, _, errGetPermis := ctx.GetClientIdentity().GetAttributeValue("permission")

	if errGetPermis != nil {
		return "", fmt.Errorf("Lỗi: %s", errGetPermis)
	}
	return permis, nil
}

//create benh an
func (c *Contract) CreateMB(ctx contractapi.TransactionContextInterface, CCCD string, information string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return errGetPermis
	}
	getRecord, errRecord := ctx.GetStub().GetState(permis)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}
	if getRecord != nil {
		return fmt.Errorf("Bản ghi đã tồn tại trên sổ cái")
	}
	//Tao so benh an
	state := new(MedicalBackground)
	state.PatientKey = permis
	state.CCCD = CCCD
	state.Information = information
	bytes, err := json.Marshal((state))
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(permis, bytes)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	return nil
}
func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func (c *Contract) DeletePermission(ctx contractapi.TransactionContextInterface, CCCD string, readPermission string, writePermission string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return fmt.Errorf("Loi: %s", errGetPermis)
	}
	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(CCCD, []string{"permis"})
	if errTaoKey != nil {
		return fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}

	var state PermissionRW
	errGetState := json.Unmarshal(getRecord, &state)
	if errRecord != nil {
		return fmt.Errorf("Loi:%s", errGetState)
	}
	if state.PatientKey != permis {
		return fmt.Errorf("Nguoi dung khong co quyen xoa")
	}

	if Contains(state.Read, readPermission) == false && readPermission != "" {
		return fmt.Errorf("Nguoi dung chua co quyen doc")
	}
	if Contains(state.Write, writePermission) == false && writePermission != "" {
		return fmt.Errorf("Nguoi dung chua co quyen ghi")
	}
	state.Read = remove(state.Read, readPermission)
	state.Write = remove(state.Write, writePermission)

	bytes, err := json.Marshal((state))
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	return nil
}

// cap quyen
func (c *Contract) GrantPermission(ctx contractapi.TransactionContextInterface, CCCD string, readPermission string, writePermission string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return fmt.Errorf("Loi: %s", errGetPermis)
	}

	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(CCCD, []string{"permis"})
	if errTaoKey != nil {
		return fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}

	var state PermissionRW
	errGetState := json.Unmarshal(getRecord, &state)
	if errRecord != nil {
		return fmt.Errorf("Loi:%s", errGetState)
	}

	if Contains(state.Read, readPermission) == true || Contains(state.Write, writePermission) == true {
		return fmt.Errorf("Nguoi dung da duoc cap quyen")
	}
	state.PatientKey = permis
	state.CCCD = CCCD
	if readPermission == "" {
		state.Write = append(state.Write, writePermission)

	} else if writePermission == "" {
		state.Read = append(state.Read, readPermission)
	} else {
		state.Read = append(state.Read, readPermission)
		state.Write = append(state.Write, writePermission)
	}

	bytes, err := json.Marshal((state))
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	return nil
}

//check element
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func CheckPermission(ctx contractapi.TransactionContextInterface, want string, permis string, CCCD string) (bool, error) {
	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(CCCD, []string{"permis"})
	if errTaoKey != nil {
		return false, fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return false, fmt.Errorf("Lỗi:%s", errRecord)
	}
	var permissionRW PermissionRW
	errGetPermission := json.Unmarshal(getRecord, &permissionRW)
	if errGetPermission != nil {
		return false, fmt.Errorf("Loi:%s", errGetPermission)
	}
	if Contains(permissionRW.Write, permis) == true && want == "write" {
		return true, nil
	}
	if Contains(permissionRW.Read, permis) == true && want == "read" {
		return true, nil
	}
	return false, fmt.Errorf("Nguoi dung khong duoc cap quyen")
}

func (c *Contract) CreateRecord(ctx contractapi.TransactionContextInterface, CCCD string, sickType string, data string, medicine string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return fmt.Errorf("Loi: %s", errGetPermis)
	}

	permissionRW, _ := CheckPermission(ctx, "write", permis, CCCD)
	if permissionRW == false {
		return fmt.Errorf("Nguoi dung khong duoc cap quyen")
	}
	currentTime := time.Now()
	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(strconv.FormatInt(currentTime.Unix(), 10), []string{CCCD, sickType, permis})
	if errTaoKey != nil {
		return fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}
	if getRecord != nil {
		return fmt.Errorf("Ban ghi da ton tai")
	}
	state := new(MedicalRecord)
	state.MRKey = permis
	state.CCCD = CCCD
	state.SickType = sickType
	state.Data = data
	state.Time = strconv.FormatInt(currentTime.Unix(), 10)
	bytes, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi:%s", err)
	}

	compositeKey, _ = ctx.GetStub().CreateCompositeKey(CCCD, []string{"Medicine"})
	medi := new(MedicineHistory)
	medi.CCCD = CCCD
	medi.Time = currentTime.String()
	medi.Data = medicine
	bytes, err = json.Marshal(medi)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi:%s", err)
	}

	return nil
}

func (c *Contract) UpdateRecord(ctx contractapi.TransactionContextInterface, CCCD string, sickType string, data string, dayOfRecord string, medicine string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return fmt.Errorf("Loi: %s", errGetPermis)
	}

	permissionRW, _ := CheckPermission(ctx, "write", permis, CCCD)
	if permissionRW == false {
		return fmt.Errorf("Nguoi dung khong duoc cap quyen")
	}
	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(dayOfRecord, []string{CCCD, sickType, permis})
	if errTaoKey != nil {
		return fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}
	if getRecord == nil {
		return fmt.Errorf("Ban ghi chua ton tai")
	}

	var state MedicalRecord
	err := json.Unmarshal(getRecord, &state)
	if err != nil {
		return fmt.Errorf("Loi:  %s", err)
	}
	if state.MRKey != permis {
		return fmt.Errorf("Nguoi dung khong co quyen cap nhat")
	}
	state.Data = data
	bytes, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi:%s", err)
	}
	currentTime := time.Now()
	compositeKey, _ = ctx.GetStub().CreateCompositeKey(CCCD, []string{"Medicine"})
	medi := new(MedicineHistory)
	medi.CCCD = CCCD
	medi.Time = currentTime.String()
	medi.Data = medicine
	bytes, err = json.Marshal(medi)
	if err != nil {
		return fmt.Errorf("Loi: %s", err)
	}
	err = ctx.GetStub().PutState(compositeKey, bytes)
	if err != nil {
		return fmt.Errorf("Loi:%s", err)
	}
	return nil
}

func (c *Contract) DeleteRecord(ctx contractapi.TransactionContextInterface, CCCD string, sickType string, dayOfRecord string) error {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return fmt.Errorf("Loi: %s", errGetPermis)
	}

	permissionRW, _ := CheckPermission(ctx, "write", permis, CCCD)
	if permissionRW == false {
		return fmt.Errorf("Nguoi dung khong duoc cap quyen xoa")
	}
	compositeKey, errTaoKey := ctx.GetStub().CreateCompositeKey(dayOfRecord, []string{CCCD, sickType, permis})
	if errTaoKey != nil {
		return fmt.Errorf("Lỗi: %s", errTaoKey)
	}
	getRecord, errRecord := ctx.GetStub().GetState(compositeKey)
	if errRecord != nil {
		return fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}
	if getRecord == nil {
		return fmt.Errorf("Bản ghi chưa tồn tại trên sổ cái")
	}
	err := ctx.GetStub().DelState(compositeKey)
	if err != nil {
		return fmt.Errorf("Loi:%s", err)
	}
	return nil
}
func (c *Contract) ReadMultiResult(ctx contractapi.TransactionContextInterface, CCCD string, sickType string) (string, error) {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return "", fmt.Errorf("Loi: %s", errGetPermis)
	}

	permissionRW, _ := CheckPermission(ctx, "read", permis, CCCD)
	if permissionRW == false {
		return "", fmt.Errorf("Nguoi dung khong duoc cap quyen")
	}
	//Lấy bản ghi
	querySelector := `{"selector":{"CCCD":"` + CCCD + `", "SickType":"` + sickType + `"},"use_index": ["_design/index1MB","indexCCCDSick"]}`
	resultIterator, errResultIterator := ctx.GetStub().GetQueryResult(querySelector)
	if errResultIterator != nil {
		return "", fmt.Errorf("Lỗi: %s", errResultIterator)
	}
	defer resultIterator.Close()
	var buffer bytes.Buffer
	//Tạo kết quả
	buffer.WriteString("[")
	next := false
	for resultIterator.HasNext() {
		item, err := resultIterator.Next()
		if err != nil {
			return "", fmt.Errorf("")
		}
		if next {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Value\":")
		buffer.WriteString(string(item.Value))
		buffer.WriteString("}")
		next = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}

// func (c *Contract) ReadRecord(ctx contractapi.TransactionContextInterface, CCCD string, sickType string) (string, error) {
// 	_, errUserRole := ctx.GetClientIdentity().GetMSPID()
// 	if errUserRole != nil {
// 		return "", fmt.Errorf("Lỗi: %s", errUserRole)
// 	}
// 	//Lấy bản ghi
// 	querySelector := `{"selector":{"CCCD":"` + CCCD + `", "SickType":"` + sickType + `"},"use_index": ["_design/index2MB","indexTime"]}`
// 	resultIterator, _, errResultIterator := ctx.GetStub().GetQueryResultWithPagination(querySelector, 3, "")
// 	if errResultIterator != nil {
// 		return "", fmt.Errorf("Lỗi: %s", errResultIterator)
// 	}
// 	defer resultIterator.Close()
// 	var buffer bytes.Buffer
// 	//Tạo kết quả
// 	buffer.WriteString("[")
// 	next := false
// 	for resultIterator.HasNext() {
// 		item, err := resultIterator.Next()
// 		if err != nil {
// 			return "", fmt.Errorf("")
// 		}
// 		if next {
// 			buffer.WriteString(",")
// 		}

// 		buffer.WriteString("{\"Value\":")
// 		buffer.WriteString(string(item.Value))
// 		buffer.WriteString("}")
// 		next = true
// 	}
// 	buffer.WriteString("]")
// 	return string(buffer.Bytes()), nil
// }

func randate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func (c *Contract) GenerateRecord(ctx contractapi.TransactionContextInterface, CCCD string) error {
	var bytes []byte
	permis, _, err := ctx.GetClientIdentity().GetAttributeValue("permission")
	if err != nil {
		return nil
	}
	state := new(MedicalRecord)
	state.MRKey = permis
	state.CCCD = CCCD
	sickList := [8]string{"Phoi", "Tim", "Chan", "Tay", "Nao", "Mat", "Tai", "Mieng"}
	currentTime := randate()
	for i := 0; i < 10; i++ {
		currentTime = randate()
		state.SickType = sickList[rand.Intn(8)]
		state.Data = `{"bacsi":"abcxyz","chan doan": "ungthu"}`
		state.Time = strconv.FormatInt(currentTime.Unix(), 10)
		bytes, _ = json.Marshal(state)
		compositeKey, _ := ctx.GetStub().CreateCompositeKey(strconv.FormatInt(currentTime.Unix(), 10), []string{CCCD, state.SickType, permis})
		ctx.GetStub().PutState(compositeKey, bytes)

	}
	return nil
}

func (c *Contract) ReadAll(ctx contractapi.TransactionContextInterface, CCCD string, hashKey string) (string, error) {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return "", fmt.Errorf("Loi: %s", errGetPermis)
	}

	if permis != hashKey {
		return "", fmt.Errorf("Nguoi dung khong duoc cap quyen")
	}
	//Lấy bản ghi
	querySelector := `{"selector":{"CCCD":"` + CCCD + `"},"use_index": ["_design/index1MB","indexCCCDSick"]}`
	resultIterator, errResultIterator := ctx.GetStub().GetQueryResult(querySelector)
	if errResultIterator != nil {
		return "", fmt.Errorf("Lỗi: %s", errResultIterator)
	}
	defer resultIterator.Close()
	var buffer bytes.Buffer
	//Tạo kết quả
	buffer.WriteString("[")
	next := false
	for resultIterator.HasNext() {
		item, err := resultIterator.Next()
		if err != nil {
			return "", fmt.Errorf("")
		}
		if next {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Value\":")
		buffer.WriteString(string(item.Value))
		buffer.WriteString("}")
		next = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}

func (c *Contract) ReadMedicine(ctx contractapi.TransactionContextInterface, CCCD string) (string, error) {
	permis, errGetPermis := CheckUser(ctx)
	if errGetPermis != nil {
		return "", fmt.Errorf("Loi: %s", errGetPermis)
	}

	permissionRW, _ := CheckPermission(ctx, "read", permis, CCCD)
	if permissionRW == false {
		return "", fmt.Errorf("Nguoi dung khong duoc cap quyen")
	}
	compositeKey, _ := ctx.GetStub().CreateCompositeKey(CCCD, []string{"Medicine"})
	getRecordIterator, errRecord := ctx.GetStub().GetHistoryForKey(compositeKey)
	if errRecord != nil {
		return "", fmt.Errorf("Không thể đọc sổ cái vì %s", errRecord)
	}
	defer getRecordIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	next := false
	for getRecordIterator.HasNext() {
		item, err := getRecordIterator.Next()

		if err != nil {
			return "", fmt.Errorf("Lỗi Query")
		}

		if next == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(item.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Value\":")
		if item.Value != nil {
			buffer.WriteString(string(item.Value))
		} else {
			buffer.WriteString(string("null"))
		}

		buffer.WriteString(",\"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(item.Timestamp.Seconds, int64(item.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(item.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		next = true
	}
	buffer.WriteString("]")
	return string(buffer.Bytes()), nil
}
