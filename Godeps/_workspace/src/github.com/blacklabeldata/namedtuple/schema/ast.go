package schema

import "sync"

// PackageList is an interface for a package registry.
type PackageList interface {
    Add(pkg Package)
    Remove(pkg string)
    Get(name string) (Package, bool)
}

// packageList contains a registry of known packages
type packageList struct {
    pkgList map[string]Package
    lock    sync.Mutex
}

func (p *packageList) Add(pkg Package) {
    p.lock.Lock()
    p.pkgList[pkg.Name] = pkg
    p.lock.Unlock()
    return
}

func (p *packageList) Remove(pkg string) {
    p.lock.Lock()
    delete(p.pkgList, pkg)
    p.lock.Unlock()
    return
}

func (p *packageList) Get(name string) (pkg Package, ok bool) {
    p.lock.Lock()
    pkg, ok = p.pkgList[name]
    p.lock.Unlock()
    return
}

// NewPackageList creates a new package registry
func NewPackageList() PackageList {
    var lock sync.Mutex
    return &packageList{make(map[string]Package), lock}
}

// Package contains an entire schema document.
type Package struct {
    Name    string
    Imports []Import
    Types   []Type
}

// Import references one or more Types from another Package
type Import struct {
    PackageName string
    TypeNames   []string
}

// Type represents a data type. It encapsulates several versions, each with their own fields.
type Type struct {
    Name     string
    Versions []Version
}

// Version is the only construct for adding one or more Fields to a Type.
type Version struct {
    Number int
    Fields []Field
}

// Field is the lowest level of granularity in a schema. Fields belong to a single Version within a Type. They are effectively immutable and should not be changed.
type Field struct {
    IsRequired bool
    IsArray    bool
    Type       string
    Name       string
}
