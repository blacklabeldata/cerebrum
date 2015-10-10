package schema

import (
    // "fmt"

    "io/ioutil"
    "log"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
    // "time"
)

func TestComment(t *testing.T) {
    text := `
    // this is a comment
    `

    var token Token
    l := NewLexer("tuple", text, func(t Token) {
        token = t
    })

    // lex comment
    lexComment(l)

    // consume token
    // token := <-l.tokens

    // expecting comment token
    assert.Equal(t, TokenComment, token.Type)

    // validate text
    assert.Equal(t, "// this is a comment", token.Value)
}

func TestMultiLineComment(t *testing.T) {
    text := `
    // this is a comment
    // This is also a comment
    // This is one too
    `

    var token Token
    l := NewLexer("tuple", text, func(t Token) {
        // fmt.Println("handler: ", t)
        token = t
    })

    // lex comment
    lexComment(l)

    // expecting comment token and validating text
    assert.Equal(t, TokenComment, token.Type)
    assert.Equal(t, "// this is a comment", token.Value)

    // lex second comment
    lexComment(l)

    // expecting comment token and validating text
    assert.Equal(t, TokenComment, token.Type)
    assert.Equal(t, "// This is also a comment", token.Value)

    // lex third comment
    lexComment(l)

    // expecting comment token and validating text
    assert.Equal(t, TokenComment, token.Type)
    assert.Equal(t, "// This is one too", token.Value)
}

func TestPackageParsing(t *testing.T) {

    text := `package users`
    var tokens []Token
    l := NewLexer("tuple", text, func(t Token) {
        // fmt.Println("handler: ", t)
        tokens = append(tokens, t)
    })

    // lex content
    l.run()

    // there should be 2 tokens
    assert.Equal(t, len(tokens), 2)

    // expecting package token and validating text
    assert.Equal(t, TokenPackage, tokens[0].Type)
    assert.Equal(t, "package", tokens[0].Value)

    // expecting package name token and validating text
    assert.Equal(t, TokenPackageName, tokens[1].Type)
    assert.Equal(t, "users", tokens[1].Value)
}

func TestImportParsing(t *testing.T) {

    text := `from project.users import User, Gadget, Widget`
    // var token Token
    var tokens []Token
    l := NewLexer("tuple", text, func(t Token) {
        // fmt.Println("handler: ", t)
        // token = t
        tokens = append(tokens, t)
    })

    // lex content
    l.run()

    // there should be 6 tokens
    assert.Equal(t, len(tokens), 8)

    // expecting from token and validating text
    assert.Equal(t, TokenFrom, tokens[0].Type)
    assert.Equal(t, "from", tokens[0].Value)

    // expecting package name token and validating text
    assert.Equal(t, TokenPackageName, tokens[1].Type)
    assert.Equal(t, "project.users", tokens[1].Value)

    // expecting import token and validating text
    assert.Equal(t, TokenImport, tokens[2].Type)
    assert.Equal(t, "import", tokens[2].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[3].Type)
    assert.Equal(t, "User", tokens[3].Value)

    // expecting comma token and validating text
    assert.Equal(t, TokenComma, tokens[4].Type)
    assert.Equal(t, ",", tokens[4].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[5].Type)
    assert.Equal(t, "Gadget", tokens[5].Value)

    // expecting comma token and validating text
    assert.Equal(t, TokenComma, tokens[6].Type)
    assert.Equal(t, ",", tokens[6].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[7].Type)
    assert.Equal(t, "Widget", tokens[7].Value)
}

func TestImportAllParsing(t *testing.T) {

    text := `from project.users import *`
    // var token Token
    var tokens []Token
    l := NewLexer("tuple", text, func(t Token) {
        // fmt.Println("handler: ", t)
        // token = t
        tokens = append(tokens, t)
    })

    // lex content
    l.run()

    // there should be 6 tokens
    assert.Equal(t, len(tokens), 4)

    // expecting from token and validating text
    assert.Equal(t, TokenFrom, tokens[0].Type)
    assert.Equal(t, "from", tokens[0].Value)

    // expecting package name token and validating text
    assert.Equal(t, TokenPackageName, tokens[1].Type)
    assert.Equal(t, "project.users", tokens[1].Value)

    // expecting import token and validating text
    assert.Equal(t, TokenImport, tokens[2].Type)
    assert.Equal(t, "import", tokens[2].Value)

    // expecting asterisk token and validating text
    assert.Equal(t, TokenAsterisk, tokens[3].Type)
    assert.Equal(t, "*", tokens[3].Value)
}

func TestTypeDef(t *testing.T) {

    text := `type User {}`
    var tokens []Token
    l := NewLexer("TypeDef", text, func(t Token) {
        // fmt.Println("handler: ", t)
        // token = t
        tokens = append(tokens, t)
    })

    // lex content
    l.run()
    // t.Log(tokens)

    // there should be 4 tokens
    assert.Equal(t, len(tokens), 4)

    // expecting type token and validating text
    assert.Equal(t, TokenTypeDef, tokens[0].Type)
    assert.Equal(t, "type", tokens[0].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[1].Type)
    assert.Equal(t, "User", tokens[1].Value)

    // expecting openScope token and validating text
    assert.Equal(t, TokenOpenCurlyBracket, tokens[2].Type)
    assert.Equal(t, "{", tokens[2].Value)

    // expecting closeScope token and validating text
    assert.Equal(t, TokenCloseCurlyBracket, tokens[3].Type)
    assert.Equal(t, "}", tokens[3].Value)
}

func TestIdentifier(t *testing.T) {

    text := `User.`
    var tokens []Token
    l := NewLexer("TestIdentifier", text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    lexIdentifier(l, nil, false)
    // t.Log(tokens)

    // there should be 1 token
    assert.Equal(t, len(tokens), 1)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[0].Type)
    assert.Equal(t, "User", tokens[0].Value)
}

func TestVersion(t *testing.T) {

    text := `version 1`
    var tokens []Token
    l := NewLexer("TestVersion", text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    lexVersion(l)
    // t.Log(tokens)

    // there should be 2 tokens
    assert.Equal(t, len(tokens), 2)

    // expecting version token and validating text
    assert.Equal(t, TokenVersion, tokens[0].Type)
    assert.Equal(t, "version", tokens[0].Value)

    // expecting version number token and validating text
    assert.Equal(t, TokenVersionNumber, tokens[1].Type)
    assert.Equal(t, "1", tokens[1].Value)
}

func TestVersionFail(t *testing.T) {

    text := `version abc`
    var tokens []Token
    l := NewLexer("TestVersionFail", text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    l.run()
    // lexVersion(l)
    // t.Log(tokens)

    // there should be 4 tokens
    assert.Equal(t, len(tokens), 4)

    // expecting version token and validating text
    assert.Equal(t, TokenVersion, tokens[0].Type)
    assert.Equal(t, "version", tokens[0].Value)

    // expecting error token and validating text
    assert.Equal(t, TokenError, tokens[1].Type)
    assert.Equal(t, "TestVersionFail[0:9] unknown token: \"a\"", tokens[1].Value)

    // expecting error token and validating text
    assert.Equal(t, TokenError, tokens[2].Type)
    assert.Equal(t, "TestVersionFail[0:10] unknown token: \"b\"", tokens[2].Value)

    // expecting error token and validating text
    assert.Equal(t, TokenError, tokens[3].Type)
    assert.Equal(t, "TestVersionFail[0:11] unknown token: \"c\"", tokens[3].Value)
}

func TestOpenScope(t *testing.T) {

    text := `{`
    var tokens []Token
    l := NewLexer("TestOpenScope", text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    l.run()

    // there should be 1 token
    assert.Equal(t, len(tokens), 1)

    // expecting openScope token and validating text
    assert.Equal(t, TokenOpenCurlyBracket, tokens[0].Type)
    assert.Equal(t, "{", tokens[0].Value)
}

func TestCloseScope(t *testing.T) {

    text := `}`
    var tokens []Token
    l := NewLexer("TestCloseScope", text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    l.run()

    // there should be 1 token
    assert.Equal(t, len(tokens), 1)

    // expecting closeScope token and validating text
    assert.Equal(t, TokenCloseCurlyBracket, tokens[0].Type)
    assert.Equal(t, "}", tokens[0].Value)
}

func TestFieldBasicType(t *testing.T) {

    for _, text := range TypeNames {

        var tokens []Token
        l := NewLexer("Test: "+text, text+" fieldname", func(t Token) {
            tokens = append(tokens, t)
        })

        // lex content
        lexType(l)
        // t.Log(tokens)

        // there should be 2 token
        assert.Equal(t, len(tokens), 2)

        // expecting value type token and validating text
        assert.Equal(t, TokenValueType, tokens[0].Type)
        assert.Equal(t, text, tokens[0].Value)

        // expecting identifier token and validating text
        assert.Equal(t, TokenIdentifier, tokens[1].Type)
        assert.Equal(t, "fieldname", tokens[1].Value)
    }
}

func TestFieldBasicArrayType(t *testing.T) {

    for _, text := range TypeNames {

        var tokens []Token
        l := NewLexer("Test: "+text, "[]"+text+" fieldname", func(t Token) {
            tokens = append(tokens, t)
        })

        // lex content
        lexType(l)
        // t.Log(tokens)

        // there should be 4 token
        assert.Equal(t, len(tokens), 4)

        // expecting openArray token and validating text
        assert.Equal(t, TokenOpenArrayBracket, tokens[0].Type)
        assert.Equal(t, "[", tokens[0].Value)

        // expecting closeArray token and validating text
        assert.Equal(t, TokenCloseArrayBracket, tokens[1].Type)
        assert.Equal(t, "]", tokens[1].Value)

        // expecting value type token and validating text
        assert.Equal(t, TokenValueType, tokens[2].Type)
        assert.Equal(t, text, tokens[2].Value)

        // expecting identifier token and validating text
        assert.Equal(t, TokenIdentifier, tokens[3].Type)
        assert.Equal(t, "fieldname", tokens[3].Value)
    }
}

func TestRequiredFieldType(t *testing.T) {

    text := `required tuple data`
    var tokens []Token
    l := NewLexer("Test: "+text, text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    l.run()
    // t.Log(tokens)

    // there should be 3 token
    assert.Equal(t, len(tokens), 3)

    // expecting required token and validating text
    assert.Equal(t, TokenRequired, tokens[0].Type)
    assert.Equal(t, "required", tokens[0].Value)

    // expecting value type token and validating text
    assert.Equal(t, TokenValueType, tokens[1].Type)
    assert.Equal(t, "tuple", tokens[1].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[2].Type)
    assert.Equal(t, "data", tokens[2].Value)
}

func TestOptionalFieldType(t *testing.T) {

    text := `optional tuple data`
    var tokens []Token
    l := NewLexer("Test: "+text, text, func(t Token) {
        tokens = append(tokens, t)
    })

    // lex content
    l.run()
    // t.Log(tokens)

    // there should be 3 token
    assert.Equal(t, len(tokens), 3)

    // expecting optional token and validating text
    assert.Equal(t, TokenOptional, tokens[0].Type)
    assert.Equal(t, "optional", tokens[0].Value)

    // expecting value type token and validating text
    assert.Equal(t, TokenValueType, tokens[1].Type)
    assert.Equal(t, "tuple", tokens[1].Value)

    // expecting identifier token and validating text
    assert.Equal(t, TokenIdentifier, tokens[2].Type)
    assert.Equal(t, "data", tokens[2].Value)
}

func TestErrorf(t *testing.T) {

    text := `package sys

    from users import User
    `

    l := NewLexer("TestErrorf", text, func(tok Token) {

        if tok.Type == TokenError {
            t.Log(tok)

            // expecting error token and validating text
            assert.Equal(t, "TestErrorf[4:1] testing error", tok.Value)
        }
    })

    l.run()
    l.errorf("testing error")
}

func TestLoop(t *testing.T) {
    text := `
    // this is a comment
    // This is also a comment
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

    var token Token
    l := NewLexer("tuple", text, func(t Token) {
        // fmt.Println("handler: ", t.Type, t)
        token = t
    })
    // lexText(l)
    //
    // var start = time.Now()
    l.run()
    // fmt.Println(time.Now().Sub(start).Seconds())
}

func TestComplexFile(t *testing.T) {
    // filename, err := filepath.Abs("./examples/complex.ent")
    file, err := os.Open("./examples/complex.ent") // For read access.
    if err != nil {
        log.Fatal(err)
    }
    bytes, err := ioutil.ReadAll(file)
    text := string(bytes)

    l := NewLexer("complex file", text, func(t Token) {
        // fmt.Printf("%#v\n", t)
    })
    l.run()
}

func BenchmarkLargePackage(b *testing.B) {

    text := `package users.package

    from some.package import Something

    // Location type
    type Location {

      // Coordinates
      version 1 {
        required float64 latitude, longitude, altitude
      }
    }

    // User object
    type User {

      // base user
      version 1 {
        required string uuid, username
        optional uint8 age
      }

      // 11/15/14
      version 2 {
          optional Location location
      }

      version 3 {
        optional []User friends
      }
    }`

    for i := 0; i < b.N; i++ {
        l := NewLexer("large package", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkLexPackageDecl(b *testing.B) {
    text := `package users`

    for i := 0; i < b.N; i++ {
        l := NewLexer("package decl", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkImportStmt(b *testing.B) {
    text := `from some.package import Something`

    for i := 0; i < b.N; i++ {
        l := NewLexer("import stmt", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkVersionStmt(b *testing.B) {
    text := `version 1 {}`

    for i := 0; i < b.N; i++ {
        l := NewLexer("version stmt", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkArrayFieldStmt(b *testing.B) {
    text := `required []User friends`

    for i := 0; i < b.N; i++ {
        l := NewLexer("version stmt", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkFieldStmt(b *testing.B) {
    text := `required uint8 age`

    for i := 0; i < b.N; i++ {
        l := NewLexer("version stmt", text, func(t Token) {
        })
        l.run()
    }
}

func BenchmarkMultiFieldStmt(b *testing.B) {
    text := `required string uuid, username, email`

    for i := 0; i < b.N; i++ {
        l := NewLexer("version stmt", text, func(t Token) {
        })
        l.run()
    }
}
