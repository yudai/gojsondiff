# Go JSON Diff

## How to use

### Getting code

```sh
go get github.com/yudai/gojsondiff
```

### Comparing two JSON strings

See `jd/main.go` for how to use this library.


## CLI tool

This repository contains a package that you can use as a CLI tool.

### How to install

```sh
go get github.com/yudai/gojsondiff/jd
go install github.com/yudai/gojsondiff/jd
```

### Usage

Just give two json files:

```sh
jd one.json another.json
```

Outputs would be something like:

```diff
 {
   "arr": [
     0: "arr0",
     1: 21,
     2: {
       "num": 1,
-      "str": "pek3f"
+      "str": "changed"
     },
     3: [
       0: 0,
-      1: "1"
+      1: "changed"
     ]
   ],
   "bool": true,
   "num_float": 39.39,
   "num_int": 13,
   "obj": {
     "arr": [
       0: 17,
       1: "str",
       2: {
-        "str": "eafeb"
+        "str": "changed"
       }
     ],
+    "new": "added",
-    "num": 19,
     "obj": {
-      "num": 14,
+      "num": 9999
-      "str": "efj3"
+      "str": "changed"
     },
     "str": "bcded"
   },
   "str": "abcde"
 }
```

## License

MIT License (see `LICENSE` for detail)
