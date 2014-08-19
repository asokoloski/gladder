package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// bindata_read reads the given file from disk. It returns an error on failure.
func bindata_read(path, name string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading asset %s at %s: %v", name, path, err)
	}
	return buf, err
}

// resources_templates_create_user_html reads file data from disk. It returns an error on failure.
func resources_templates_create_user_html() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/templates/create_user.html",
		"resources/templates/create_user.html",
	)
}

// resources_templates_index_html reads file data from disk. It returns an error on failure.
func resources_templates_index_html() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/templates/index.html",
		"resources/templates/index.html",
	)
}

// resources_templates_create_user_html_ reads file data from disk. It returns an error on failure.
func resources_templates_create_user_html_() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/templates/create_user.html~",
		"resources/templates/create_user.html~",
	)
}

// resources_templates_index_html_ reads file data from disk. It returns an error on failure.
func resources_templates_index_html_() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/templates/index.html~",
		"resources/templates/index.html~",
	)
}

// resources_css_site_css reads file data from disk. It returns an error on failure.
func resources_css_site_css() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/css/site.css",
		"resources/css/site.css",
	)
}

// resources_css_site_css_ reads file data from disk. It returns an error on failure.
func resources_css_site_css_() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/css/site.css~",
		"resources/css/site.css~",
	)
}

// resources_js_gladder_js_ reads file data from disk. It returns an error on failure.
func resources_js_gladder_js_() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/js/gladder.js~",
		"resources/js/gladder.js~",
	)
}

// resources_js_gladder_js reads file data from disk. It returns an error on failure.
func resources_js_gladder_js() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/js/gladder.js",
		"resources/js/gladder.js",
	)
}

// resources_js_jquery_2_1_1_min_js reads file data from disk. It returns an error on failure.
func resources_js_jquery_2_1_1_min_js() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/js/jquery-2.1.1.min.js",
		"resources/js/jquery-2.1.1.min.js",
	)
}

// resources_js_mustache_js reads file data from disk. It returns an error on failure.
func resources_js_mustache_js() ([]byte, error) {
	return bindata_read(
		"/home/aaron/stuff/go/src/gladder/resources/js/mustache.js",
		"resources/js/mustache.js",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"resources/templates/create_user.html":  resources_templates_create_user_html,
	"resources/templates/index.html":        resources_templates_index_html,
	"resources/templates/create_user.html~": resources_templates_create_user_html_,
	"resources/templates/index.html~":       resources_templates_index_html_,
	"resources/css/site.css":                resources_css_site_css,
	"resources/css/site.css~":               resources_css_site_css_,
	"resources/js/gladder.js~":              resources_js_gladder_js_,
	"resources/js/gladder.js":               resources_js_gladder_js,
	"resources/js/jquery-2.1.1.min.js":      resources_js_jquery_2_1_1_min_js,
	"resources/js/mustache.js":              resources_js_mustache_js,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"resources": &_bintree_t{nil, map[string]*_bintree_t{
		"templates": &_bintree_t{nil, map[string]*_bintree_t{
			"create_user.html":  &_bintree_t{resources_templates_create_user_html, map[string]*_bintree_t{}},
			"index.html":        &_bintree_t{resources_templates_index_html, map[string]*_bintree_t{}},
			"create_user.html~": &_bintree_t{resources_templates_create_user_html_, map[string]*_bintree_t{}},
			"index.html~":       &_bintree_t{resources_templates_index_html_, map[string]*_bintree_t{}},
		}},
		"css": &_bintree_t{nil, map[string]*_bintree_t{
			"site.css":  &_bintree_t{resources_css_site_css, map[string]*_bintree_t{}},
			"site.css~": &_bintree_t{resources_css_site_css_, map[string]*_bintree_t{}},
		}},
		"js": &_bintree_t{nil, map[string]*_bintree_t{
			"gladder.js~":         &_bintree_t{resources_js_gladder_js_, map[string]*_bintree_t{}},
			"gladder.js":          &_bintree_t{resources_js_gladder_js, map[string]*_bintree_t{}},
			"jquery-2.1.1.min.js": &_bintree_t{resources_js_jquery_2_1_1_min_js, map[string]*_bintree_t{}},
			"mustache.js":         &_bintree_t{resources_js_mustache_js, map[string]*_bintree_t{}},
		}},
	}},
}}
