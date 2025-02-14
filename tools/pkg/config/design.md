### package

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | `string`  |        | ❌       | package name       |
| version        | `string`  | \<latest\>       | ✅       | package version    |

### upstream

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | `string`  | "conan"       | ✅       | upstream package platform   |
| config        | `map[string]string`  | []       | ✅       | platform CLI option |

### toolchain

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | `string`  | "llcppg"       | ✅       | toolchain name  |
| version        | `string`  | "latest" | ✅       | toolchain version   |