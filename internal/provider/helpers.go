package provider

import (
	"context"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
	"github.com/infrahq/infra/uid"
)

func ParseID(d *schema.ResourceData, key string) (uid.ID, error) {
	return uid.Parse([]byte(d.Get(key).(string)))
}

func GrantResource(d *schema.ResourceData) string {
	resource := d.Get("cluster").(string)
	if namespace := d.Get("namespace").(string); namespace != "" {
		resource = fmt.Sprintf("%s.%s", resource, namespace)
	}

	return resource
}

func DurationDiffSuppressFunc() schema.SchemaDiffSuppressFunc {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		oldDuration, err := time.ParseDuration(oldValue)
		if err != nil {
			return false
		}

		newDuration, err := time.ParseDuration(newValue)
		if err != nil {
			return false
		}

		return oldDuration == newDuration
	}
}

func DecodePEM(data []byte, keytype string) ([]byte, error) {
	blocks, _ := pem.Decode(data)
	if blocks == nil || blocks.Type != keytype {
		return nil, fmt.Errorf("failed to decode PEM block containing %s", strings.ToLower(keytype))
	}

	return pem.EncodeToMemory(blocks), nil
}

func DecodePEMFile(filepath string, keytype string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return DecodePEM(data, keytype)
}

func requireMinimumServerVersion(ctx context.Context, client *api.Client, minimum string) error {
	minimumSemVer, err := semver.NewVersion(minimum)
	if err != nil {
		return err
	}

	serverVersion, err := client.GetServerVersion(ctx)
	if err != nil {
		return err
	}

	serverSemVer, err := semver.NewVersion(serverVersion.Version)
	if err != nil {
		return err
	}

	if serverSemVer.LessThan(minimumSemVer) {
		return fmt.Errorf("server must be at least version %s or higher. currently %s", minimum, serverVersion.Version)
	}

	return nil
}
