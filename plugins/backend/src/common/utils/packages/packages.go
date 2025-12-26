package packages

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// get packages info
// _packageName
// ê°€?¸ì˜¤?¤ëŠ” package???´ë¦„
// import ?€ ?¬ë¦¬ ?¤í–‰?˜ë ¤???Œì¼??ê¸°ì??¼ë¡œ ?ë?ê²½ë¡œë¥??…ë ¥?´ì•¼??
// ex)
//  1. common ?´ë” ?„ëž˜??main.go ?Œì¼?ì„œ
//     "cia/common/utils/types" package ?•ë³´ë¥?ê°€?¸ì˜¤?¤ê³  ????
//     _packageName -> "/utils/types"
//  2. common ?´ë” ?„ëž˜??main.go ?Œì¼?ì„œ
//     "apiServer/logHandler" package ?•ë³´ë¥?ê°€?¸ì˜¤??????
//     _packageName -> "../../apiServer/logHandler"
func GetPackages(_packageName string) map[string]*ast.Package {
	pkgs, err := parser.ParseDir(
		token.NewFileSet(),
		_packageName,
		nil,
		0,
	)
	if err != nil {
		panic(err)
	}

	return pkgs
}

// get info of all files on packages
func GetPackageFiles(_packageName string) []map[string]*ast.File {
	pkgs := GetPackages(_packageName)

	mapFiles := []map[string]*ast.File{}
	for _, pkg := range pkgs {
		mapFiles = append(mapFiles, pkg.Files)
	}

	return mapFiles
}

// get all function names of packages
func GetPackageFunctions(_packageName string) []string {
	mapFiles := GetPackageFiles(_packageName)

	var funcNames []string
	for _, files := range mapFiles {
		for _, file := range files {
			// function names on file
			for _, decl := range file.Decls {
				if function, ok := decl.(*ast.FuncDecl); ok {
					funcNames = append(funcNames, function.Name.String())
				}
			}
		}
	}

	return funcNames
}
