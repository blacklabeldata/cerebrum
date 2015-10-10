package schema

import "testing"
import "github.com/stretchr/testify/assert"

func TestParse(t *testing.T) {

    text := `
    package users

    from locale import Location, Country

    // This is a comment
    type User {
        // version comment
        version 1 {
            required string uuid, username
            optional uint8 age, num
        }

        // 11/15/14
        version 2 {
            optional Location location
        }
    }
    `

    pkgList := NewPackageList()

    // create parser
    parser := NewParser(pkgList)
    pkg, err := parser.Parse("TestParse", text)
    assert.Nil(t, err)
    assert.NotNil(t, pkg)

    // test pkg name
    assert.Equal(t, pkg.Name, "users")

    // test imports
    assert.Equal(t, len(pkg.Imports), 1)
    assert.Equal(t, pkg.Imports[0].PackageName, "locale")
    assert.Equal(t, len(pkg.Imports[0].TypeNames), 2)
    assert.Equal(t, pkg.Imports[0].TypeNames[0], "Location")
    assert.Equal(t, pkg.Imports[0].TypeNames[1], "Country")

    // test types
    assert.Equal(t, len(pkg.Types), 1)
    assert.Equal(t, pkg.Types[0].Name, "User")

    // test versions
    assert.Equal(t, len(pkg.Types[0].Versions), 2)
    assert.Equal(t, pkg.Types[0].Versions[0].Number, 1)
    assert.Equal(t, len(pkg.Types[0].Versions[0].Fields), 4)

    assert.Equal(t, pkg.Types[0].Versions[0].Fields[0].IsRequired, true)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[0].IsArray, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[0].Type, "string")
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[0].Name, "uuid")

    assert.Equal(t, pkg.Types[0].Versions[0].Fields[1].IsRequired, true)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[1].IsArray, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[1].Type, "string")
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[1].Name, "username")

    assert.Equal(t, pkg.Types[0].Versions[0].Fields[2].IsRequired, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[2].IsArray, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[2].Type, "uint8")
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[2].Name, "age")

    assert.Equal(t, pkg.Types[0].Versions[0].Fields[3].IsRequired, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[3].IsArray, false)
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[3].Type, "uint8")
    assert.Equal(t, pkg.Types[0].Versions[0].Fields[3].Name, "num")

    assert.Equal(t, pkg.Types[0].Versions[1].Number, 2)
    assert.Equal(t, pkg.Types[0].Versions[1].Fields[0].IsRequired, false)
    assert.Equal(t, pkg.Types[0].Versions[1].Fields[0].IsArray, false)
    assert.Equal(t, pkg.Types[0].Versions[1].Fields[0].Type, "Location")
    assert.Equal(t, pkg.Types[0].Versions[1].Fields[0].Name, "location")

    // t.Logf("%#v\n", pkg)
    // t.Log(err)
}

func TestLoadDirectory(t *testing.T) {

    pkgList := NewPackageList()
    p := NewParser(pkgList)
    err := LoadDirectory("examples", p)
    t.Log(err)

    t.Log(pkgList)
    t.Log(pkgList.Get("users"))

}

func BenchmarkParse(b *testing.B) {

    text := `
    package users

    from locale import Location

    // This is one too
    type User {
        // version comment
        version 1 {
            required string uuid
            required string username
            optional uint8 age
        }

        // 11/15/14
        version 2 {
            optional Location location
        }
    }
    `

    pkgList := NewPackageList()

    // create parser
    parser := NewParser(pkgList)

    for i := 0; i < b.N; i++ {
        parser.Parse("TestParse", text)
    }
    // t.Logf("%#v\n", pkg)
    // t.Log(err)
}
