# gendiff

Utility for calculating diff between two files

[![Actions Status](https://github.com/i1yas/go-from-scratch-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/i1yas/go-from-scratch-project-244/actions)
[![Quality gate status](https://sonarcloud.io/api/project_badges/measure?project=i1yas_go-from-scratch-project-244&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=i1yas_go-from-scratch-project-244)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=i1yas_go-from-scratch-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=i1yas_go-from-scratch-project-244)


## Usage in CLI

Command expects two file paths

```bash
make build # build binary
./bin/gendiff file1.json file2.json # run binary from bin
```

Output
```
{
  + email: myemail@mydomain.com
    id: 101
  - name: john
  + name: John
  - phone: 123
}
```

### Format of input files
Format of input files detected automatically from file extension.
Supported __JSON__ and __YAML__ formats.

### Format of ouput

Format can be specified by `--format` or `-f` flags.

There are 3 output formats:
- __stylish__ (default) - human-readable, json-like format
- __plain__ - flat diff representation, only changes are shown
- __json__ - format that outputs valid json representation of diff

## stylish
Outputs human-readable format that resembles json. Keys are sorted in alphabetical order.

Note: output is not valid json, if you need one, use `json` format.
```
./bin/gendiff -f stylish file1.json file2.json
```

Output:
```
{
  + email: myemail@mydomain.com
    id: 101
  - name: john
  + name: John
  - phone: 123
}
```

## plain
Flat representation of diff. Only changes are shown. Each line represents one change. For nested changes shows path to the key where each key separated by `.` (`somekey.nested.key`).

```
./bin/gendiff -f plain file1.json file2.json
```

Output:
```
Property 'email' was added with value: 'myemail@mydomain.com'
Property 'name' was updated. From 'john' to 'John'
Property 'phone' was removed
```

## json
Outputs diff representation as valid json. Diff represented as tree of nodes. Root node has key `""` and type `"root"`.

Node types:
- `"root"` - root node, has `"children"` array of nested nodes
- `"nested"` - node with nested changes, has `"children"` array of nested nodes
- `"added"` - node with added value in `"value2"` key
- `"deleted"` - node with removed value in `"value1"` key
- `"changed"` - node with changed value, old value in `"value1"` key, new value in `"value2"` key
- `"unchanged"` - node with unchanged value, value in `"value2"` key


```
./bin/gendiff -f json file1.json file2.json
```

Output (prettified):
```
{
    "key": "",
    "type": "root",
    "children": [
        {
            "key": "email",
            "type": "added",
            "value2": "myemail@maydomain.com"
        },
        {
            "key": "id",
            "type": "unchanged",
            "value1": 101
        },
        {
            "key": "name",
            "type": "changed",
            "value1": "john",
            "value2": "John"
        },
        {
            "key": "phone",
            "type": "deleted",
            "value1": 123
        }
    ]
}
```