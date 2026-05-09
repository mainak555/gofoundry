package util

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func GetAzCredential() (azcore.TokenCredential, error) {
	if os.Getenv("NODE_ENV") == "local" {
		return azidentity.NewClientSecretCredential(
			os.Getenv("AZ_TENENT_ID"),
			os.Getenv("AZ_CLIENT_ID"),
			os.Getenv("AZ_CLIENT_SECRET"),
			nil)
	}
	return azidentity.NewDefaultAzureCredential(nil)
}

func GetVault(vaultName string) (*azsecrets.Client, error) {
	credential, err := GetAzCredential()
	if err != nil {
		return nil, err
	}

	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net", vaultName)
	return azsecrets.NewClient(keyVaultUrl, credential, nil)
}

func GetFromVault(client *azsecrets.Client, secretKey string) (string, error) {
	secret, err := client.GetSecret(context.TODO(), secretKey, "", nil)
	if err != nil {
		return secretKey, err
	}
	return *secret.Value, nil
}
