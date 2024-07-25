package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	AppraisedValue int     `json:"AppraisedValue"`
	Color          string  `json:"Color"`
	ID             string  `json:"ID"`
	Manufacter     string  `json:"Manufacter"`
	Material       string  `json:"Material"`
	Model          string  `json:"Model"`
	ProductionDate string  `json:"ProductionDate"` // Ver si se puede introducir con tipo DATE/TIME
	Recyclability  string  `json:"Recyclability"`
	SerialNumber   string  `json:"SerialNumber"` // Ver en qué formato implantarlo
	Size           int     `json:"Size"`
	Weight         float32 `json:"Weight"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Manufacter: "Audi", Model: "A4", Material: "Plastic", Color: "black", ProductionDate: "12-07-2023",
			SerialNumber: "SN123456789", Recyclability: "low", Size: 135, Weight: 50.0, AppraisedValue: 300},
		{ID: "asset2", Manufacter: "BMW", Model: "X5", Material: "Metal", Color: "red", ProductionDate: "24-06-21",
			SerialNumber: "SN987654321", Recyclability: "medium", Size: 150, Weight: 75.5, AppraisedValue: 400},
		{ID: "asset3", Manufacter: "Tesla", Model: "Model 3", Material: "Aluminum", Color: "white", ProductionDate: "15-11-2022",
			SerialNumber: "SN192837465", Recyclability: "high", Size: 180, Weight: 68.7, AppraisedValue: 500},
		{ID: "asset4", Manufacter: "Ford", Model: "Mustang", Material: "Steel", Color: "blue", ProductionDate: "04-02-2024",
			SerialNumber: "SN543216789", Recyclability: "low", Size: 190, Weight: 80.0, AppraisedValue: 600},
		{ID: "asset5", Manufacter: "Chevrolet", Model: "Camaro", Material: "Carbon Fiber", Color: "yellow", ProductionDate: "30-09-2019",
			SerialNumber: "SN987321654", Recyclability: "medium", Size: 185, Weight: 78.3, AppraisedValue: 550},
	}

	for _, asset := range assets { // El uso de _ en el bucle for indica que el índice de cada iteración no es relevante para el procesamiento actual
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil // return nil en una función que devuelve un error indica que la función ha terminado exitosamente sin encontrar ningún error
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, appraisedValue int, color string, manufacter string, material string, model string, productionDate string, recyclability string, serialNumber string, size int, weight float32) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		AppraisedValue: appraisedValue,
		Color:          color,
		ID:             id,
		Manufacter:     manufacter,
		Material:       material,
		Model:          model,
		ProductionDate: productionDate,
		Recyclability:  recyclability,
		SerialNumber:   serialNumber,
		Size:           size,
		Weight:         weight,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON) // Si PutState se ejecuta correctamente, retorna nil, indicando éxito. Si hay un error, se retorna el error
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil { //  If the key does not exist in the state database, (nil, nil) is returned
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset) // El segundo parámetro es una referencia a la variable donde se almacenarán los datos deserializados
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, appraisedValue int, color string, manufacter string, material string, model string, productionDate string, recyclability string, serialNumber string, size int, weight float32) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		AppraisedValue: appraisedValue,
		Color:          color,
		ID:             id,
		Manufacter:     manufacter,
		Material:       material,
		Model:          model,
		ProductionDate: productionDate,
		Recyclability:  recyclability,
		SerialNumber:   serialNumber,
		Size:           size,
		Weight:         weight,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newManufacter string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldManufacter := asset.Manufacter
	asset.Manufacter = newManufacter

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldManufacter, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "") // resultsIterator: Es un iterador que permite recorrer todos los registros que coinciden con la consulta.
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close() // Asegura que el iterador se cerrará cuando la función GetAllAssets termine su ejecución

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset) // Deserializa los datos JSON de queryResponse.Value en la variable asset
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset) // Agrega un puntero al asset deserializado al slice assets.
	}

	return assets, nil
}
