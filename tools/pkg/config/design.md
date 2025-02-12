# llpkg.cfg 配置文件说明

## package

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | string  |        | ✅       | package name       |
| version        | string  |        | ✅       | package version    |
| versionChange  | bool    | true       | ❌     | use llpkg's auto version change |

## upstream

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | string  | "conan"       | ❌       | upstream package platform   |
| config.options | string  | ""       | ❌       | platform CLI option |

## toolchain

| key         | type    | default value   | optional | meaning           |
| -------------- | ------- | -------- | -------- | -------------- |
| name           | string  | "llcppg"       | ❌       | toolchain name  |
| version        | string  | "latest" | ❌       | toolchain version   |