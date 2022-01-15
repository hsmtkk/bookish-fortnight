# bookish-fortnight
Azure Functions timer trigger sample

## memo

* リソースグループを作成する

az group create --name sample-azure-functions-test-01 --location japaneast

* Cosmos DBを作成する

az cosmosdb create --name sample-azfunc-cosmosdb-01 --resource-group sample-azure-functions-test-01
