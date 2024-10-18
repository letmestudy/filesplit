# filesplit
一个文件分割工具，使用 golang 编写，支持多架构编译

# Uasge
## 文件分割
```shell
filesplit split input.txt 10M output_dir
```

## 文件合并
```shell
filesplit merge merged_input.txt output_dir
```

## 文件校验
```shell
filesplit verify merged_input.txt input.txt # 文件之间校验
filesplit verify split_dir input.txt # 文件夹与文件之间校验
```
