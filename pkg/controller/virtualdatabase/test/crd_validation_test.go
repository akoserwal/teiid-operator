package test

import (
	"strings"
	"testing"

	"github.com/RHsyseng/operator-utils/pkg/validation"
	"github.com/ghodss/yaml"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/assert"
	"github.com/teiid/teiid-operator/pkg/apis/vdb/v1alpha1"
)

func TestSampleCustomResources(t *testing.T) {
	schema := getSchema(t)
	box := packr.New("deploy/crs", "../../../../deploy/crs")
	for _, file := range box.List() {
		yamlString, err := box.FindString(file)
		assert.NoError(t, err, "Error reading %v CR yaml", file)
		var input map[string]interface{}
		assert.NoError(t, yaml.Unmarshal([]byte(yamlString), &input))
		assert.NoError(t, schema.Validate(input), "File %v does not validate against the CRD schema", file)
	}
}

func TestExampleCustomResources(t *testing.T) {
	schema := getSchema(t)
	box := packr.New("deploy/examples", "../../../../deploy/examples")
	for _, file := range box.List() {
		yamlString, err := box.FindString(file)
		assert.NoError(t, err, "Error reading %v CR yaml", file)
		var input map[string]interface{}
		assert.NoError(t, yaml.Unmarshal([]byte(yamlString), &input))
		assert.NoError(t, schema.Validate(input), "File %v does not validate against the CRD schema", file)
	}
}

func TestTrialEnvMinimum(t *testing.T) {
	var inputYaml = `
apiVersion: vdb.teiid.io/v1alpha1
kind: VirtualDatabase
metadata:
  name: trial
spec:
  build:
    gitSource:
      uri: https://github.com/teiid/teiid-openshift-examples
`
	var input map[string]interface{}
	assert.NoError(t, yaml.Unmarshal([]byte(inputYaml), &input))

	schema := getSchema(t)
	assert.NoError(t, schema.Validate(input))

	//	deleteNestedMapEntry(input, "spec", "environment")
	//	assert.Error(t, schema.Validate(input))
}

func TestCompleteCRD(t *testing.T) {
	schema := getSchema(t)
	missingEntries := schema.GetMissingEntries(&v1alpha1.VirtualDatabase{})
	for _, missing := range missingEntries {
		if strings.HasPrefix(missing.Path, "/status") {
			//Not using subresources, so status is not expected to appear in CRD
		} else if strings.Contains(missing.Path, "/env/valueFrom/") {
			//The valueFrom is not expected to be used and is not fully defined TODO: verify
		} else if strings.HasSuffix(missing.Path, "/from/uid") {
			//The ObjectReference in From is not expected to be used and is not fully defined TODO: verify
		} else if strings.HasSuffix(missing.Path, "/from/apiVersion") {
			//The ObjectReference in From is not expected to be used and is not fully defined TODO: verify
		} else if strings.HasSuffix(missing.Path, "/from/resourceVersion") {
			//The ObjectReference in From is not expected to be used and is not fully defined TODO: verify
		} else if strings.HasSuffix(missing.Path, "/from/fieldPath") {
			//The ObjectReference in From is not expected to be used and is not fully defined TODO: verify
		} else {
			assert.Fail(t, "Discrepancy between CRD and Struct", "Missing or incorrect schema validation at %v, expected type %v", missing.Path, missing.Type)
		}
	}
}

func deleteNestedMapEntry(object map[string]interface{}, keys ...string) {
	for index := 0; index < len(keys)-1; index++ {
		object = object[keys[index]].(map[string]interface{})
	}
	delete(object, keys[len(keys)-1])
}

func getSchema(t *testing.T) validation.Schema {
	box := packr.New("deploy/crds", "../../../../deploy/crds")
	crdFile := "virtualdatabase.crd.yaml"
	assert.True(t, box.Has(crdFile))
	yamlString, err := box.FindString(crdFile)
	assert.NoError(t, err, "Error reading CRD yaml %v", yamlString)
	schema, err := validation.New([]byte(yamlString))
	assert.NoError(t, err)
	return schema
}
