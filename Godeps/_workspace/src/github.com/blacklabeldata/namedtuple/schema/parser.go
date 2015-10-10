package schema

import (
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
)

// Config simply stores the parsing configuration.
type Config struct {
    PackageRootDir string
}

// SyntaxError represents an error while parsing the schema
type SyntaxError struct {
    Message string
}

func (s SyntaxError) Error() string {
    return s.Message
}

// LoadDirectory reads all the schema files from a directory.
func LoadDirectory(dir string, parser Parser) (err error) {

    // Open dir for reading
    d, err := os.Open(dir)
    if err != nil {
        return
    }

    // Iterate over all the files in the directory.
    for {

        // Only read 128 files at a time.
        if fis, err := d.Readdir(128); err == nil {

            // Read each entry
            for _, fi := range fis {
                // fmt.Println("%#v", fi)

                // If the FileInfo is a directory, read the directory.
                // Otherwise, read the file.
                switch fi.IsDir() {
                case true:

                    // return error if there is one
                    if err := LoadDirectory(fi.Name(), parser); err != nil {
                        return err
                    }
                case false:

                    // All schema files should end with .nt
                    if !strings.HasSuffix(fi.Name(), ".ent") {
                        break
                    }

                    // Read the file
                    if _, err := LoadFile(filepath.Join(dir, fi.Name()), parser); err != nil {
                        return err
                    }
                }
            }
        } else if err == io.EOF {
            // If there are no more files in the directory, break.
            break
        } else {
            // If there is any other error, return it.
            return err
        }
    }

    // If you have reached this far, you are done.
    return nil
}

// LoadFile reads a schema document from a file.
func LoadFile(filename string, parser Parser) (Package, error) {
    file, err := os.Open(filename)
    if err != nil {
        return Package{}, err
    }
    defer file.Close()

    // read file
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return Package{}, err
    }

    // convert to string and load
    return parser.Parse(file.Name(), string(bytes))
}

// LoadPackage parses a text string.
func LoadPackage(parser Parser, name, text string) (Package, error) {
    return parser.Parse(name, text)
}

func NewParser(pkgList PackageList) Parser {
    var lock sync.Mutex
    return &parser{pkgList, []Token{}, 0, lock, ""}
}

type Parser interface {
    Parse(name, text string) (Package, error)
}

type parser struct {
    pkgList PackageList
    tokens  []Token
    pos     int
    lock    sync.Mutex
    name    string
}

func (p *parser) Parse(name string, text string) (pkg Package, err error) {
    p.lock.Lock()
    defer p.lock.Unlock()
    p.name = name

    l := NewLexer(name, text, func(tok Token) {
        p.tokens = append(p.tokens, tok)
    })
    l.run()

    pkg, err = p.parsePackage()

    // if no error, add to package list
    if err == nil {
        p.pkgList.Add(pkg)
    }
    return
}

func (p *parser) advance(skip int) {
    p.pos += skip
}

func (p *parser) current() (tok Token) {
    if p.pos >= len(p.tokens) {
        tok = Token{TokenError, "end of input"}
    } else {
        tok = p.tokens[p.pos]
        if tok.Type == TokenComment {
            p.advance(1)
            return p.current()
        }
    }
    return
}

func (p *parser) next() (tok Token) {
    tok = p.current()
    p.advance(1)

    // ignore comments
    if tok.Type == TokenComment {
        tok = p.next()
    }

    return
}

func (p *parser) backup() {
    if p.pos > 0 {
        p.pos--
    }
}

func (p *parser) typeCheck(t TokenType, errMsg string) (tok Token, err error) {

    // next token
    tok = p.next()

    // is it an error
    if tok.Type == TokenError {
        return tok, SyntaxError{p.name + ": " + tok.Value}
    }

    // is it the correct type
    if tok.Type != t {
        // fmt.Println(p.tokens[p.pos:])
        return tok, SyntaxError{p.name + ": " + errMsg}
    }

    // currect token type
    return
}

func (p *parser) parsePackage() (pkg Package, err error) {
    if len(p.tokens) == 0 {
        return pkg, SyntaxError{"empty input string"}
    }

    // consume package decl
    _, err = p.typeCheck(TokenPackage, "expected package declaration")
    if err != nil {
        return
    }

    // consume package name
    tok, err := p.typeCheck(TokenPackageName, "expected package name")
    if err != nil {
        return
    }
    pkg.Name = tok.Value

    // parse imports
    if err = p.parseImports(&pkg); err != nil {
        return
    }

    // parse types
    if err = p.parseTypes(&pkg); err != nil {
        return
    }

    return
}

func (p *parser) parseImports(pkg *Package) (err error) {

    for p.current().Type == TokenFrom {

        var imp Import

        // consume 'from' keyword
        if _, err := p.typeCheck(TokenFrom, "expected 'from' keyword"); err != nil {
            return err
        }

        // consume package name
        tok, err := p.typeCheck(TokenPackageName, "expected package name")
        if err != nil {
            return err
        }

        // set import package name
        imp.PackageName = tok.Value

        // consume 'import' keyword
        if _, err := p.typeCheck(TokenImport, "expected 'import' keyword"); err != nil {
            return err
        }

        // consume type name
        tok, err = p.typeCheck(TokenIdentifier, "expected type name")
        if err != nil {
            return err
        }
        imp.TypeNames = append(imp.TypeNames, tok.Value)

        // consume multiple type names
        for p.current().Type == TokenComma {
            // skip comma token
            p.advance(1)

            // consume type name
            tok, err = p.typeCheck(TokenIdentifier, "expected type name")
            if err != nil {
                return err
            }
            imp.TypeNames = append(imp.TypeNames, tok.Value)
        }

        // add import to package
        pkg.Imports = append(pkg.Imports, imp)
    }

    return nil
}

func (p *parser) parseTypes(pkg *Package) (err error) {

    // iterate over type defs
    for p.current().Type == TokenTypeDef {

        var t Type

        // consume 'type' keyword
        if _, err := p.typeCheck(TokenTypeDef, "expected 'type' keyword"); err != nil {
            return err
        }

        // consume type name
        tok, err := p.typeCheck(TokenIdentifier, "expected type name")
        if err != nil {
            return err
        }

        // set type name
        t.Name = tok.Value

        // consume open scope
        if _, err := p.typeCheck(TokenOpenCurlyBracket, "expected open bracket"); err != nil {
            return err
        }

        // parse versions
        if err = p.parseVersions(pkg, &t); err != nil {
            return err
        }

        // consume close scope
        _, err = p.typeCheck(TokenCloseCurlyBracket, "expected close bracket")
        if err != nil {
            return err
        }

        pkg.Types = append(pkg.Types, t)
    }

    return nil
}

func (p *parser) parseVersions(pkg *Package, t *Type) (err error) {

    // iterate over versions
    for p.current().Type == TokenVersion {

        var ver Version

        // consume 'version' keyword
        if _, err := p.typeCheck(TokenVersion, "expected 'version' keyword"); err != nil {
            return err
        }

        // consume version number
        tok, err := p.typeCheck(TokenVersionNumber, "expected version number")
        if err != nil {
            return err
        }

        num, err := strconv.Atoi(tok.Value)
        if err != nil {
            return err
        }

        // set version num
        ver.Number = num

        // consume open scope
        if _, err := p.typeCheck(TokenOpenCurlyBracket, "expected open bracket"); err != nil {
            return err
        }

        // parse fields
    OUTER:
        for {
            switch p.current().Type {
            case TokenRequired, TokenOptional:
                if err = p.parseField(pkg, &ver); err != nil {
                    return err
                }
            default:
                break OUTER
            }
        }

        // consume close scope
        _, err = p.typeCheck(TokenCloseCurlyBracket, "expected close bracket")
        if err != nil {
            return err
        }

        t.Versions = append(t.Versions, ver)
    }
    return nil
}

func (p *parser) parseField(pkg *Package, ver *Version) (err error) {

    var field Field
    switch p.current().Type {
    case TokenRequired:
        field.IsRequired = true
    case TokenOptional:
        field.IsRequired = false
    default:
        return SyntaxError{"expected 'required' or 'optional' keyword"}
    }
    p.advance(1)

    // consume optional array bracket
    tok := p.next()
    if tok.Type == TokenOpenArrayBracket {
        field.IsArray = true

        _, err = p.typeCheck(TokenCloseArrayBracket, "expected array close bracket")
        if err != nil {
            return err
        }
    } else if tok.Type == TokenError {
        return SyntaxError{tok.Value}
    } else {
        p.backup()
    }

    // field type should be next
    if p.current().Type != TokenValueType {
        return SyntaxError{"expected field type, not '" + p.current().Value + "'"}
    }

    typeName := p.current().Value
    var found bool
    for _, t := range TypeNames {
        if typeName == t {
            field.Type = typeName
            found = true
            break
        }
    }

    // eval import stmts
    if !found {
    OUTER:
        for _, imp := range pkg.Imports {
            for _, typ := range imp.TypeNames {
                if typ == typeName {
                    found = true
                    field.Type = typeName
                    break OUTER
                }
            }
        }
    }

    // eval previous package types
    for _, typ := range pkg.Types {
        if typ.Name == typeName {
            found = true
            field.Type = typeName
            break
        }
    }

    // If type has still not been found, return error
    if !found {
        return SyntaxError{"unknown type '" + typeName + "'"}
    }

    // Skip type name
    p.advance(1)

    // Consume field names
    for {

        // Perform type check
        tok, err = p.typeCheck(TokenIdentifier, "expected field name")
        if err != nil {
            return err
        }

        // Set field name
        field.Name = tok.Value

        // Add field to version
        ver.Fields = append(ver.Fields, field)

        // Determine if next token is a comma
        tok, err = p.typeCheck(TokenComma, "")

        // If the next token is not a comma, backup
        if err != nil {
            break
        } else {
            continue
        }
    }
    p.backup()

    return nil
}
