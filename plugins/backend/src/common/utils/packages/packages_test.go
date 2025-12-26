package packages

import (
	"testing"

	// testify assert
	"github.com/stretchr/testify/assert"
)

func Test_GetPackageFiles(t *testing.T) {
	pkgsName := "../../utils/types"
	pkgsFiles := GetPackageFiles(pkgsName)

	assert.NotNil(t, pkgsFiles)
	assert.Equal(t, len(pkgsFiles), 1)

	var fileNames [][]string
	for _, files := range pkgsFiles {
		var names []string
		for _, file := range files {
			names = append(names, file.Name.String())
		}
		fileNames = append(fileNames, names)
	}

	assert.Equal(t, len(fileNames), 1)
}

func Test_GetPackageFunctions(t *testing.T) {
	pkgsName := "../rpcRequest"
	pkgsFunc := GetPackageFunctions(pkgsName)

	assert.Equal(t, len(pkgsFunc), 8)
	assert.Equal(t, pkgsFunc[0], "ToString")
	assert.Equal(t, pkgsFunc[1], "ResultToString")
	assert.Equal(t, pkgsFunc[2], "ResultToJsonString")
	assert.Equal(t, pkgsFunc[3], "ResultToUint64")
	assert.Equal(t, pkgsFunc[4], "ResultToInt64")
	assert.Equal(t, pkgsFunc[5], "ErrorToString")
	assert.Equal(t, pkgsFunc[6], "RpcRequest")
	assert.Equal(t, pkgsFunc[7], "Test_RpcRequest")
}
