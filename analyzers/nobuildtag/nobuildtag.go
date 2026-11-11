package nobuildtag

import (
	"go/ast"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "nobuildtag",
	Doc:  "no build constraints, GOOS/GOARCH file names, excluded files or hidden directories: all code must be visible to the linter",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		checkFile(pass, file)
	}
	for _, name := range pass.IgnoredFiles {
		reportIgnored(pass, name)
	}
	if len(pass.Files) > 0 {
		checkHiddenDirs(pass, filepath.Dir(pass.Fset.File(pass.Files[0].Pos()).Name()))
	}
	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File) {
	if hasConstraint(file) {
		kit.Report(pass, file.Package, MSG_CONSTRAINT)
	}
	if platformName(pass.Fset.File(file.Pos()).Name()) {
		kit.Report(pass, file.Package, MSG_SUFFIX)
	}
}

func hasConstraint(file *ast.File) bool {
	for _, group := range file.Comments {
		if group.Pos() >= file.Package {
			return false
		}
		if groupHasConstraint(group) {
			return true
		}
	}
	return false
}

func groupHasConstraint(group *ast.CommentGroup) bool {
	for _, comment := range group.List {
		if strings.HasPrefix(comment.Text, GO_BUILD) || strings.HasPrefix(comment.Text, PLUS_BUILD) {
			return true
		}
	}
	return false
}

// platformName is true for names like x_linux.go, x_arm64.go, x_linux_arm64.go.
func platformName(path string) bool {
	parts := strings.Split(strings.TrimSuffix(filepath.Base(path), GO_SUFFIX), NAME_SEP)
	if len(parts) < MIN_PARTS {
		return false
	}
	last := parts[len(parts)-1]
	return contains(OS_LIST, last) || contains(ARCH_LIST, last)
}

func contains(list string, word string) bool {
	for _, item := range strings.Split(list, LIST_SEP) {
		if item == word {
			return true
		}
	}
	return false
}

// reportIgnored positions the finding on a file the build excluded: it was
// never parsed, so a synthetic one-byte file gives the diagnostic a place.
func reportIgnored(pass *analysis.Pass, name string) {
	if !strings.HasSuffix(name, GO_SUFFIX) || strings.HasSuffix(name, kit.TEST_SUFFIX) {
		return
	}
	info, err := os.Stat(name)
	if err != nil {
		return
	}
	file := pass.Fset.AddFile(name, -1, int(info.Size()))
	kit.Report(pass, token.Pos(file.Base()), MSG_IGNORED)
}

func checkHiddenDirs(pass *analysis.Pass, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		sub := filepath.Join(dir, entry.Name())
		if entry.IsDir() && isHidden(sub, entry.Name()) && holdsGo(sub) {
			kit.Reportf(pass, pass.Files[0].Package, MSG_HIDDEN, entry.Name())
		}
	}
}

func isHidden(path string, name string) bool {
	if strings.HasPrefix(name, NAME_SEP) || strings.HasPrefix(name, HIDDEN_DOT) || name == VENDOR_DIR {
		return true
	}
	_, err := os.Stat(filepath.Join(path, GO_MOD))
	return err == nil
}

// holdsGo is true when root has a Go file; an unreadable tree counts as true.
func holdsGo(root string) bool {
	found := false
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		found = !entry.IsDir() && strings.HasSuffix(path, GO_SUFFIX)
		if found {
			return fs.SkipAll
		}
		return nil
	})
	return found || err != nil
}
